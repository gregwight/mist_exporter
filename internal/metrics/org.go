package metrics

import (
	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

type OrgMetrics struct {
	Alarms  *prometheus.Desc
	Tickets *prometheus.Desc
	Site    *prometheus.Desc
}

var siteLabels = []string{
	"site_id",
	"site_name",
	"country_code",
}

func SiteLabels(s mistclient.Site) []string {
	return []string{
		s.ID,
		s.Name,
		s.CountryCode,
	}
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
		Site: prometheus.NewDesc(
			"mist_org_site",
			"A site confured in the organization",
			siteLabels,
			nil,
		),
	}
}
