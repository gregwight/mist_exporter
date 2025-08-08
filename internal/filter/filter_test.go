package filter

import (
	"testing"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/config"
)

func TestApplySiteFilter(t *testing.T) {
	sites := map[string]mistclient.Site{
		"HQ":      {Name: "Headquarters"},
		"BO1":     {Name: "Branch Office 1"},
		"BO2":     {Name: "Branch Office 2"},
		"WH":      {Name: "Warehouse"},
		"Lab":     {Name: "Test Lab"},
		"EU-Prod": {Name: "EU-Prod-Site"},
		"EU-Test": {Name: "EU-Test-Site"},
	}

	testCases := []struct {
		name           string
		filterCfg      *config.SiteFilter
		siteKey        string
		expectFiltered bool
		expectErr      bool
	}{
		{
			name:           "no filter, should not be filtered",
			filterCfg:      &config.SiteFilter{},
			siteKey:        "HQ",
			expectFiltered: false,
		},
		{
			name: "include only, match",
			filterCfg: &config.SiteFilter{
				Include: []string{"Branch*"},
			},
			siteKey:        "BO1",
			expectFiltered: false,
		},
		{
			name: "include only, no match",
			filterCfg: &config.SiteFilter{
				Include: []string{"Branch*"},
			},
			siteKey:        "HQ",
			expectFiltered: true,
		},
		{
			name: "exclude only, match",
			filterCfg: &config.SiteFilter{
				Exclude: []string{"*Lab"},
			},
			siteKey:        "Lab",
			expectFiltered: true,
		},
		{
			name: "exclude only, no match",
			filterCfg: &config.SiteFilter{
				Exclude: []string{"*Lab"},
			},
			siteKey:        "HQ",
			expectFiltered: false,
		},
		{
			name: "include and exclude, matches include, not exclude",
			filterCfg: &config.SiteFilter{
				Include: []string{"*Office*"},
				Exclude: []string{"*Test*"},
			},
			siteKey:        "BO1",
			expectFiltered: false,
		},
		{
			name: "include and exclude, matches exclude",
			filterCfg: &config.SiteFilter{
				Include: []string{"EU-*"},
				Exclude: []string{"*-Test-Site"},
			},
			siteKey:        "EU-Test",
			expectFiltered: true,
		},
		{
			name: "include and exclude, matches both (exclude wins)",
			filterCfg: &config.SiteFilter{
				Include: []string{"EU-Test-Site"},
				Exclude: []string{"EU-Test-Site"},
			},
			siteKey:        "EU-Test",
			expectFiltered: true,
		},
		{
			name: "include and exclude, does not match include",
			filterCfg: &config.SiteFilter{
				Include: []string{"EU-*"},
				Exclude: []string{"*-Staging-*"},
			},
			siteKey:        "HQ",
			expectFiltered: true,
		},
		{
			name:      "invalid include pattern",
			filterCfg: &config.SiteFilter{Include: []string{"["}},
			siteKey:   "HQ",
			expectErr: true,
		},
		{
			name:      "invalid exclude pattern",
			filterCfg: &config.SiteFilter{Exclude: []string{"["}},
			siteKey:   "HQ",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := New(tc.filterCfg)

			if tc.expectErr {
				if err == nil {
					t.Fatal("New() expected an error, but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("New() returned an unexpected error: %v", err)
			}

			site := sites[tc.siteKey]
			isFiltered, err := f.IsFiltered(site)
			if err != nil {
				t.Fatalf("IsFiltered() returned an unexpected error: %v", err)
			}

			if isFiltered != tc.expectFiltered {
				t.Errorf("IsFiltered() for site %q returned %v, want %v", site.Name, isFiltered, tc.expectFiltered)
			}
		})
	}
}
