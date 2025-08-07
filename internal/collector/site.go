package collector

import (
	"github.com/gregwight/mistexporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	latDesc = prometheus.NewDesc(
		"mist_site_lat",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	lngDesc = prometheus.NewDesc(
		"mist_site_lng",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	modifiedTimeDesc = prometheus.NewDesc(
		"mist_site_modified_time",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numAPDesc = prometheus.NewDesc(
		"mist_site_num_ap",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numAPConnectedDesc = prometheus.NewDesc(
		"mist_site_num_ap_connected",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numClientsDesc = prometheus.NewDesc(
		"mist_site_num_clients",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numDevicesDesc = prometheus.NewDesc(
		"mist_site_num_devices",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numDevicesConnectedDesc = prometheus.NewDesc(
		"mist_site_num_devices_connected",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numGatewayDesc = prometheus.NewDesc(
		"mist_site_num_gateway",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numGatewayConnectedDesc = prometheus.NewDesc(
		"mist_site_num_gateway_connected",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numSwitchDesc = prometheus.NewDesc(
		"mist_site_num_switch",
		"",
		metrics.SiteLabelNames,
		nil,
	)
	numSwitchConnectedDesc = prometheus.NewDesc(
		"mist_site_num_switch_connected",
		"",
		metrics.SiteLabelNames,
		nil,
	)
)

func (c *MistCollector) collectSiteStats(ch chan<- prometheus.Metric) {
	defer c.wg.Done()

	sites, err := c.client.GetOrgSites(c.orgID)
	if err != nil {
		c.logger.Error("unable to fetch sites", "error", err)
		return
	}

	for _, site := range sites {
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()

			stat, err := c.client.GetSiteStats(site.ID)
			if err != nil {
				c.logger.Error("unable to fetch site stats", "error", err)
				return
			}

			labels := metrics.SiteLabelValues(site)

			c.sendMetric(ch, latDesc, prometheus.GaugeValue, float64(stat.Lat), labels...)
			c.sendMetric(ch, lngDesc, prometheus.GaugeValue, float64(stat.Lng), labels...)
			c.sendMetric(ch, modifiedTimeDesc, prometheus.GaugeValue, float64(stat.ModifiedTime.Unix()), labels...)
			c.sendMetric(ch, numAPDesc, prometheus.GaugeValue, float64(stat.NumAP), labels...)
			c.sendMetric(ch, numAPConnectedDesc, prometheus.GaugeValue, float64(stat.NumAPConnected), labels...)
			c.sendMetric(ch, numClientsDesc, prometheus.GaugeValue, float64(stat.NumClients), labels...)
			c.sendMetric(ch, numDevicesDesc, prometheus.GaugeValue, float64(stat.NumDevices), labels...)
			c.sendMetric(ch, numDevicesConnectedDesc, prometheus.GaugeValue, float64(stat.NumDevicesConnected), labels...)
			c.sendMetric(ch, numGatewayDesc, prometheus.GaugeValue, float64(stat.NumGateway), labels...)
			c.sendMetric(ch, numGatewayConnectedDesc, prometheus.GaugeValue, float64(stat.NumGatewayConnected), labels...)
			c.sendMetric(ch, numSwitchDesc, prometheus.GaugeValue, float64(stat.NumSwitch), labels...)
			c.sendMetric(ch, numSwitchConnectedDesc, prometheus.GaugeValue, float64(stat.NumSwitchConnected), labels...)

		}()

	}

}
