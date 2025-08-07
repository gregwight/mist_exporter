package metrics

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
)

// MistMetrics
type MistMetrics struct {
	config *config.Collector
	client *mistclient.APIClient
	orgID  string
	ready  chan struct{}
	reg    *prometheus.Registry
	logger *slog.Logger

	mu    sync.RWMutex
	sites map[string]*StreamCollector
}

// New creates a new MistMetrics
func New(config *config.Collector, client *mistclient.APIClient, orgID string, reg *prometheus.Registry, logger *slog.Logger) *MistMetrics {
	deviceMetrics = newDeviceMetrics(reg)
	clientMetrics = newClientMetrics(reg)

	return &MistMetrics{
		config: config,
		client: client,
		orgID:  orgID,
		ready:  make(chan struct{}),
		reg:    reg,
		logger: logger.With(slog.String("component", "metrics")),
		sites:  make(map[string]*StreamCollector),
	}
}

func (c *MistMetrics) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	if err := c.manageSiteStreams(ctx, wg); err != nil {
		return fmt.Errorf("unable to initialize site metric streams: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(c.config.SiteRefreshInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.manageSiteStreams(ctx, wg); err != nil {
					c.logger.Error("unable to refresh site metric streams", "error", err)
				}
			}
		}
	}()

	close(c.ready)
	wg.Wait()
	return nil
}

func (c *MistMetrics) manageSiteStreams(ctx context.Context, wg *sync.WaitGroup) error {
	c.logger.Debug("running site metric stream manager...")
	defer c.logger.Debug("site metric stream manager finished")

	c.mu.Lock()
	defer c.mu.Unlock()

	sites, err := c.client.GetOrgSites(c.orgID)
	if err != nil {
		return fmt.Errorf("unable to fetch site list: %w", err)
	}

	activeSites := make(map[string]struct{})
	for _, site := range sites {
		activeSites[site.ID] = struct{}{}
		streamer, ok := c.sites[site.ID]
		if !ok {
			streamer = newStreamCollector(c.client, site, c.logger)
			c.sites[site.ID] = streamer
		}

		streamer.mu.RLock()
		if !streamer.running {
			wg.Add(1)
			go streamer.run(ctx, wg)
		}
		streamer.mu.RUnlock()
	}

	for siteID, streamer := range c.sites {
		if _, ok := activeSites[siteID]; !ok {
			streamer.stop()
			delete(c.sites, siteID)
		}
	}

	return nil
}

func (c *MistMetrics) Ready() <-chan struct{} {
	return c.ready
}

type StreamCollector struct {
	client *mistclient.APIClient
	site   mistclient.Site
	logger *slog.Logger

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
}

func (c *StreamCollector) stop() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cancel != nil {
		c.cancel()
	}
}

func newStreamCollector(client *mistclient.APIClient, site mistclient.Site, logger *slog.Logger) *StreamCollector {
	return &StreamCollector{
		client: client,
		site:   site,
		logger: logger.With(slog.String("site", site.Name)),
	}
}

func (c *StreamCollector) run(ctx context.Context, wg *sync.WaitGroup) {
	runCtx, cancel := context.WithCancel(ctx)

	c.mu.Lock()
	c.running = true
	c.cancel = cancel
	c.mu.Unlock()

	c.logger.Info("starting site metrics stream...")
	defer func() {
		cancel()
		c.mu.Lock()
		c.running = false
		c.cancel = nil
		c.mu.Unlock()

		c.logger.Info("site metrics stream stopped")
		wg.Done()
	}()

	deviceStats, err := c.client.StreamSiteDeviceStats(runCtx, c.site.ID)
	if err != nil {
		c.logger.Error("unable to start site device stats stream", "error", err)
		return
	}
	c.logger.Debug("site device stats stream started")

	clientStats, err := c.client.StreamSiteClientStats(runCtx, c.site.ID)
	if err != nil {
		c.logger.Error("unable to start site client stats stream", "error", err)
		return
	}
	c.logger.Debug("site client stats stream started")

	// WaitGroup to ensure all subscriptions are closed before we exit.
	// If we get a failure on one channel we cancel the context to
	// force the other channels to disconnect. We will be restarted by
	// the stream manager unless the parent context is done.
	hwg := &sync.WaitGroup{}
	hwg.Add(1)
	go func() {
		defer hwg.Done()
		defer cancel()

		for stat := range deviceStats {
			handleSiteDeviceStat(c.site, stat)
		}
	}()

	hwg.Add(1)
	go func() {
		defer hwg.Done()
		defer cancel()

		for stat := range clientStats {
			handleSiteClientStat(c.site, stat)
		}
	}()

	hwg.Wait()
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
