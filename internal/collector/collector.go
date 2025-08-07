package collector

import (
	"log/slog"
	"sync"

	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

type MistCollector struct {
	client *mistclient.APIClient
	orgID  string
	wg     *sync.WaitGroup
	logger *slog.Logger
}

func New(client *mistclient.APIClient, orgID string, logger *slog.Logger) *MistCollector {
	return &MistCollector{
		client: client,
		orgID:  orgID,
		wg:     &sync.WaitGroup{},
		logger: logger.With(slog.String("component", "collector")),
	}
}

// Describe implements the prometheus.Collector interface
func (c *MistCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect implements the prometheus.Collector interface
func (c *MistCollector) Collect(ch chan<- prometheus.Metric) {
	// Get alarms for the organization
	c.wg.Add(1)
	go c.collectOrgAlarms(ch)

	// Get tickets for the organization
	c.wg.Add(1)
	go c.collectOrgTickets(ch)

	// Get sites for the organization
	c.wg.Add(1)
	go c.collectSiteStats(ch)

	c.wg.Wait()
}

func (c *MistCollector) sendMetric(ch chan<- prometheus.Metric, desc *prometheus.Desc, valueType prometheus.ValueType, value float64, labels ...string) {
	ch <- prometheus.MustNewConstMetric(desc, valueType, value, labels...)
}
