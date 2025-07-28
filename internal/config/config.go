package config

import (
	"fmt"
	"os"
	"time"

	"github.com/gregwight/mistclient"
	"gopkg.in/yaml.v3"
)

const (
	defaultExporterAddress string        = "0.0.0.0"
	defaultExporterPort    int           = 9200
	defaultCollectTimeout  time.Duration = 30 * time.Second
)

type Config struct {
	OrgId      string             `yaml:"org_id"`
	MistClient *mistclient.Config `yaml:"mist_api"`
	Exporter   *Exporter          `yaml:"exporter"`
	Collector  *Collector         `yaml:"collector"`
}

type Exporter struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Collector struct {
	Timeout time.Duration `yaml:"timeout"`
}

// loadConfig loads and processes the YAML configuration with environment variable substitution
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Replace environment variables in the format ${VAR_NAME}
	configStr := string(data)
	configStr = os.ExpandEnv(configStr)

	config := newDefaultConfig()
	if err := yaml.Unmarshal([]byte(configStr), config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	return config, nil
}

func newDefaultConfig() *Config {
	return &Config{
		Exporter: &Exporter{
			Address: defaultExporterAddress,
			Port:    defaultExporterPort,
		},
		Collector: &Collector{
			Timeout: defaultCollectTimeout,
		},
	}
}
