package metrics

import "github.com/gregwight/mistclient"

var SiteLabelNames = []string{
	"site_id",
	"site_name",
	"country_code",
	"timezone",
}

func SiteLabelValues(s mistclient.Site) []string {
	return []string{
		s.ID,
		s.Name,
		s.CountryCode,
		s.Timezone,
	}
}
