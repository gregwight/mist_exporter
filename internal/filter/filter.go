package filter

import (
	"fmt"
	"path/filepath"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/config"
)

// Filter holds the patterns for including and excluding sites
type Filter struct {
	include []string
	exclude []string
}

// New creates a new site filter from the configuration
func New(cfg *config.SiteFilter) (*Filter, error) {
	if cfg == nil {
		return &Filter{
			include: []string{},
			exclude: []string{},
		}, nil
	}

	// Validate patterns
	for _, p := range cfg.Include {
		if _, err := filepath.Match(p, ""); err != nil {
			return nil, fmt.Errorf("invalid include glob pattern %q: %w", p, err)
		}
	}
	for _, p := range cfg.Exclude {
		if _, err := filepath.Match(p, ""); err != nil {
			return nil, fmt.Errorf("invalid exclude glob pattern %q: %w", p, err)
		}
	}

	return &Filter{
		include: cfg.Include,
		exclude: cfg.Exclude,
	}, nil
}

// IsFiltered determines if a site should be filtered out based on the rules.
func (f *Filter) IsFiltered(site mistclient.Site) (bool, error) {
	// Exclusion takes precedence.
	if excluded, err := f.matches(site.Name, f.exclude); err != nil {
		return false, fmt.Errorf("unable to match site name %q against exclude patterns: %w", site.Name, err)
	} else if excluded {
		return true, nil
	}

	// If the site isn't excluded and there are
	// no explicit includes we can shortcut
	if len(f.include) == 0 {
		return false, nil
	}

	if included, err := f.matches(site.Name, f.include); err != nil {
		return false, fmt.Errorf("unable to match site name %q against include patterns: %w", site.Name, err)
	} else if !included {
		return true, nil
	}

	return false, nil
}

// matches checks if a name matches any of the glob patterns.
func (f *Filter) matches(name string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		if matched, err := filepath.Match(pattern, name); err != nil {
			return false, err
		} else if matched {
			return true, nil
		}
	}

	return false, nil
}
