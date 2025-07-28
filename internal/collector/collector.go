package collector

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/config"
	"github.com/gregwight/mistexporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var _ prometheus.Collector = &MistCollector{}

// MistCollector implements the prometheus.Collector interface
type MistCollector struct {
	sync.RWMutex

	config *config.Collector
	client *mistclient.APIClient
	orgID  string
	logger *slog.Logger

	orgMetrics    *metrics.OrgMetrics
	deviceMetrics *metrics.DeviceMetrics
	clientMetrics *metrics.ClientMetrics
}

// New creates a new MistCollector
func New(config *config.Collector, client *mistclient.APIClient, orgID string, logger *slog.Logger) *MistCollector {
	return &MistCollector{
		config: config,
		client: client,
		orgID:  orgID,
		logger: logger,

		orgMetrics:    metrics.NewOrgMetrics(),
		deviceMetrics: metrics.NewDeviceMetrics(),
		clientMetrics: metrics.NewClientMetrics(),
	}
}

// Describe implements the Collector interface
func (c *MistCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect implements the Collector interface
func (c *MistCollector) Collect(ch chan<- prometheus.Metric) {
	c.Lock()
	defer c.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	// Get alarms for the organization
	alarms, err := c.client.CountOrgAlarms(c.orgID)
	if err != nil {
		c.logger.Error("unable to fetch alarms", "error", err)
		return
	}

	for alarmType, count := range alarms {
		ch <- prometheus.MustNewConstMetric(
			c.orgMetrics.Alarms,
			prometheus.GaugeValue,
			float64(count),
			alarmType,
		)
	}

	// Get tickets for the organization
	tickets, err := c.client.CountOrgTickets(c.orgID)
	if err != nil {
		c.logger.Error("unable to fetch tickets", "error", err)
		return
	}

	for status, count := range tickets {
		ch <- prometheus.MustNewConstMetric(
			c.orgMetrics.Tickets,
			prometheus.GaugeValue,
			float64(count),
			status.String(),
		)
	}

	// Get sites for the organization
	sites, err := c.client.GetOrgSites(c.orgID)
	if err != nil {
		c.logger.Error("unable to fetch sites", "error", err)
		return
	}

	// Use errgroup for concurrent collection
	g, gCtx := errgroup.WithContext(ctx)

	// Limit concurrency to avoid overwhelming the API
	semaphore := make(chan struct{}, 10)

	for _, site := range sites {
		ch <- prometheus.MustNewConstMetric(
			c.orgMetrics.Sites,
			prometheus.GaugeValue,
			1,
			site.ID, site.Name, site.CountryCode,
		)
		g.Go(func() error {
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
				return c.collectSiteMetrics(gCtx, ch, site)
			case <-gCtx.Done():
				return gCtx.Err()
			}
		})
	}

	if err := g.Wait(); err != nil {
		c.logger.Error("unable to collect site metrics", "error", err)
	}
}

func (c *MistCollector) collectSiteMetrics(ctx context.Context, ch chan<- prometheus.Metric, site mistclient.Site) error {
	// Collect devices and clients concurrently
	g, gCtx := errgroup.WithContext(ctx)

	// Collect devices
	g.Go(func() error {
		return c.collectDeviceMetrics(gCtx, ch, site)
	})

	// Collect clients
	g.Go(func() error {
		return c.collectClientMetrics(gCtx, ch, site)
	})

	return g.Wait()
}

func (c *MistCollector) collectDeviceMetrics(ctx context.Context, ch chan<- prometheus.Metric, site mistclient.Site) error {
	deviceStats, err := c.client.GetSiteDeviceStats(site.ID)
	if err != nil {
		return fmt.Errorf("unable to fetch device stats for site %s: %w", site.Name, err)
	}

	for _, deviceStat := range deviceStats {
		deviceLabels := metrics.DeviceStatLabels(deviceStat)
		c.sendMetric(ch, c.deviceMetrics.LastSeen, prometheus.GaugeValue, float64(deviceStat.LastSeen.Unix()), deviceLabels...)
		c.sendMetric(ch, c.deviceMetrics.Uptime, prometheus.GaugeValue, float64(deviceStat.Uptime), deviceLabels...)
		c.sendMetric(ch, c.deviceMetrics.WLANs, prometheus.GaugeValue, float64(deviceStat.NumWLANs), deviceLabels...)
		c.sendMetric(ch, c.deviceMetrics.TxBps, prometheus.GaugeValue, float64(deviceStat.TxBps), deviceLabels...)
		c.sendMetric(ch, c.deviceMetrics.RxBps, prometheus.GaugeValue, float64(deviceStat.RxBps), deviceLabels...)

		for radioConfig, radioStat := range deviceStat.RadioStats {
			radioLabels := metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())
			c.sendMetric(ch, c.deviceMetrics.Clients, prometheus.GaugeValue, float64(radioStat.NumClients), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.TxBytes, prometheus.GaugeValue, float64(radioStat.TxBytes), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.RxBytes, prometheus.GaugeValue, float64(radioStat.RxBytes), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.TxPackets, prometheus.GaugeValue, float64(radioStat.TxPkts), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.RxPackets, prometheus.GaugeValue, float64(radioStat.RxPkts), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.Power, prometheus.GaugeValue, float64(radioStat.Power), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.Channel, prometheus.GaugeValue, float64(radioStat.Channel), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.Bandwidth, prometheus.GaugeValue, float64(radioStat.Bandwidth), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilAll, prometheus.GaugeValue, float64(radioStat.UtilAll), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilTx, prometheus.GaugeValue, float64(radioStat.UtilTx), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilRxInBSS, prometheus.GaugeValue, float64(radioStat.UtilRxInBSS), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilRxOtherBSS, prometheus.GaugeValue, float64(radioStat.UtilRxOtherBSS), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilUnknownWiFi, prometheus.GaugeValue, float64(radioStat.UtilUnknownWiFi), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilNonWiFi, prometheus.GaugeValue, float64(radioStat.UtilNonWiFi), radioLabels...)
			c.sendMetric(ch, c.deviceMetrics.UtilUndecodableWiFi, prometheus.GaugeValue, float64(radioStat.UtilUndecodableWiFi), radioLabels...)
		}
	}

	return nil
}

func (c *MistCollector) collectClientMetrics(ctx context.Context, ch chan<- prometheus.Metric, site mistclient.Site) error {
	clients, err := c.client.GetSiteClients(site.ID)
	if err != nil {
		return fmt.Errorf("unable to fetch client stats for site %s: %w", site.Name, err)
	}

	for _, client := range clients {
		clientLabels := metrics.ClientLabels(client)
		c.sendMetric(ch, c.clientMetrics.LastSeen, prometheus.GaugeValue, float64(client.LastSeen.Unix()), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.Uptime, prometheus.GaugeValue, float64(client.Uptime), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.Idletime, prometheus.GaugeValue, float64(client.Idletime), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.PowerSaving, prometheus.GaugeValue, boolToFloat64(client.PowerSaving), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.DualBand, prometheus.GaugeValue, boolToFloat64(client.DualBand), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.Channel, prometheus.GaugeValue, float64(client.Channel), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.RSSI, prometheus.GaugeValue, float64(client.RSSI), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.SNR, prometheus.GaugeValue, float64(client.SNR), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.TxRate, prometheus.GaugeValue, client.TxRate, clientLabels...)
		c.sendMetric(ch, c.clientMetrics.RxRate, prometheus.GaugeValue, client.RxRate, clientLabels...)
		c.sendMetric(ch, c.clientMetrics.TxBytes, prometheus.GaugeValue, float64(client.TxBytes), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.RxBytes, prometheus.GaugeValue, float64(client.RxBytes), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.TxBps, prometheus.GaugeValue, float64(client.TxBps), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.RxBps, prometheus.GaugeValue, float64(client.RxBps), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.TxPackets, prometheus.GaugeValue, float64(client.TxPackets), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.RxPackets, prometheus.GaugeValue, float64(client.RxPackets), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.TxRetries, prometheus.GaugeValue, float64(client.TxRetries), clientLabels...)
		c.sendMetric(ch, c.clientMetrics.RxRetries, prometheus.GaugeValue, float64(client.RxRetries), clientLabels...)
	}

	return nil
}

func (c *MistCollector) sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labels ...string) {
	ch <- prometheus.MustNewConstMetric(desc, valueType, value, labels...)
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
