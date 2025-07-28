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
		ch <- prometheus.MustNewConstMetric(
			c.deviceMetrics.LastSeen,
			prometheus.GaugeValue,
			float64(deviceStat.LastSeen.Unix()),
			metrics.DeviceStatLabels(deviceStat)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.deviceMetrics.Uptime,
			prometheus.GaugeValue,
			float64(deviceStat.Uptime),
			metrics.DeviceStatLabels(deviceStat)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.deviceMetrics.WLANs,
			prometheus.GaugeValue,
			float64(deviceStat.NumWLANs),
			metrics.DeviceStatLabels(deviceStat)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.deviceMetrics.TxBps,
			prometheus.GaugeValue,
			float64(deviceStat.TxBps),
			metrics.DeviceStatLabels(deviceStat)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.deviceMetrics.RxBps,
			prometheus.GaugeValue,
			float64(deviceStat.RxBps),
			metrics.DeviceStatLabels(deviceStat)...,
		)

		for radioConfig, radioStat := range deviceStat.RadioStats {
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.Clients,
				prometheus.GaugeValue,
				float64(radioStat.NumClients),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.TxBytes,
				prometheus.GaugeValue,
				float64(radioStat.TxBytes),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.RxBytes,
				prometheus.GaugeValue,
				float64(radioStat.RxBytes),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.TxPackets,
				prometheus.GaugeValue,
				float64(radioStat.TxPkts),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.RxPackets,
				prometheus.GaugeValue,
				float64(radioStat.RxPkts),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.Power,
				prometheus.GaugeValue,
				float64(radioStat.Power),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.Channel,
				prometheus.GaugeValue,
				float64(radioStat.Channel),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.Bandwidth,
				prometheus.GaugeValue,
				float64(radioStat.Bandwidth),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilAll,
				prometheus.GaugeValue,
				float64(radioStat.UtilAll),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilTx,
				prometheus.GaugeValue,
				float64(radioStat.UtilTx),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilRxInBSS,
				prometheus.GaugeValue,
				float64(radioStat.UtilRxInBSS),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilRxOtherBSS,
				prometheus.GaugeValue,
				float64(radioStat.UtilRxOtherBSS),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilUnknownWiFi,
				prometheus.GaugeValue,
				float64(radioStat.UtilUnknownWiFi),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilNonWiFi,
				prometheus.GaugeValue,
				float64(radioStat.UtilNonWiFi),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.deviceMetrics.UtilUndecodableWiFi,
				prometheus.GaugeValue,
				float64(radioStat.UtilUndecodableWiFi),
				metrics.DeviceStatLabelsWithRadio(deviceStat, radioConfig.String())...,
			)
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
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.LastSeen,
			prometheus.GaugeValue,
			float64(client.LastSeen.Unix()),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.Uptime,
			prometheus.GaugeValue,
			float64(client.Uptime),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.Idletime,
			prometheus.GaugeValue,
			float64(client.Idletime),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.PowerSaving,
			prometheus.GaugeValue,
			boolToFloat64(client.PowerSaving),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.DualBand,
			prometheus.GaugeValue,
			boolToFloat64(client.DualBand),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.Channel,
			prometheus.GaugeValue,
			float64(client.Channel),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.RSSI,
			prometheus.GaugeValue,
			float64(client.RSSI),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.SNR,
			prometheus.GaugeValue,
			float64(client.SNR),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.TxRate,
			prometheus.GaugeValue,
			client.TxRate,
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.RxRate,
			prometheus.GaugeValue,
			client.RxRate,
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.TxBytes,
			prometheus.GaugeValue,
			float64(client.TxBytes),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.RxBytes,
			prometheus.GaugeValue,
			float64(client.RxBytes),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.TxBps,
			prometheus.GaugeValue,
			float64(client.TxBps),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.RxBps,
			prometheus.GaugeValue,
			float64(client.RxBps),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.TxPackets,
			prometheus.GaugeValue,
			float64(client.TxPackets),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.RxPackets,
			prometheus.GaugeValue,
			float64(client.RxPackets),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.TxRetries,
			prometheus.GaugeValue,
			float64(client.TxRetries),
			metrics.ClientLabels(client)...,
		)
		ch <- prometheus.MustNewConstMetric(
			c.clientMetrics.RxRetries,
			prometheus.GaugeValue,
			float64(client.RxRetries),
			metrics.ClientLabels(client)...,
		)
	}

	return nil
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
