package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	alarmsDesc = prometheus.NewDesc(
		"mist_org_alarms",
		"Number of alarms in the organization",
		[]string{"alarm_type"},
		nil,
	)
	ticketsDesc = prometheus.NewDesc(
		"mist_org_tickets",
		"Number of tickets in the organization",
		[]string{"ticket_status"},
		nil,
	)
)

func (c *MistCollector) collectOrgAlarms(ch chan<- prometheus.Metric) {
	defer c.wg.Done()

	alarms, err := c.client.CountOrgAlarms(c.orgID)
	if err != nil {
		c.logger.Error("unable to fetch org alarms", "error", err)
		return
	}

	for alarmType, count := range alarms {
		c.sendMetric(ch, alarmsDesc, prometheus.GaugeValue, float64(count), alarmType)
	}
}

func (c *MistCollector) collectOrgTickets(ch chan<- prometheus.Metric) {
	defer c.wg.Done()

	tickets, err := c.client.CountOrgTickets(c.orgID)
	if err != nil {
		c.logger.Error("unable to fetch org tickets", "error", err)
		return
	}

	for status, count := range tickets {
		c.sendMetric(ch, ticketsDesc, prometheus.GaugeValue, float64(count), status.String())

	}
}
