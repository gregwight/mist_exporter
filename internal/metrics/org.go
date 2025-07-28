package metrics

import "github.com/prometheus/client_golang/prometheus"

type OrgMetrics struct {
	Alarms  *prometheus.Desc
	Tickets *prometheus.Desc
	Sites   *prometheus.Desc
}

func NewOrgMetrics() *OrgMetrics {
	return &OrgMetrics{
		Alarms: prometheus.NewDesc(
			"mist_org_alarms",
			"Number of alarms in the organization",
			[]string{"alarm_type"},
			nil,
		),
		Tickets: prometheus.NewDesc(
			"mist_org_tickets",
			"Number of tickets in the organization",
			[]string{"ticket_status"},
			nil,
		),
		Sites: prometheus.NewDesc(
			"mist_org_sites",
			"Number of sites in the organization",
			[]string{"site_id", "site_name", "country_code"},
			nil,
		),
	}
}
