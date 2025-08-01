package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Setup: Create a temporary config file
	content := `
org_id: "my-org-id"
mist_api:
  base_url: "https://api.mist.com"
  api_key: "my-api-key"
  timeout: 15s
exporter:
  address: "127.0.0.1"
  port: 9999
collector:
  timeout: 25s
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	// Test loading the config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() returned an unexpected error: %v", err)
	}

	if cfg.OrgId != "my-org-id" {
		t.Errorf("expected OrgId to be 'my-org-id', got %q", cfg.OrgId)
	}
	if cfg.MistClient.BaseURL != "https://api.mist.com" {
		t.Errorf("expected BaseURL to be 'https://api.mist.com', got %q", cfg.MistClient.BaseURL)
	}
	if cfg.MistClient.APIKey != "my-api-key" {
		t.Errorf("expected APIKey to be 'my-api-key', got %q", cfg.MistClient.APIKey)
	}
	if cfg.MistClient.Timeout != (15 * time.Second) {
		t.Errorf("expected Timeout to be 15s, got %v", cfg.MistClient.Timeout)
	}
	if cfg.Exporter.Address != "127.0.0.1" {
		t.Errorf("expected Address to be '127.0.0.1', got %q", cfg.Exporter.Address)
	}
	if cfg.Exporter.Port != 9999 {
		t.Errorf("expected Port to be 9999, got %d", cfg.Exporter.Port)
	}
	if cfg.Collector.Timeout != 25*time.Second {
		t.Errorf("expected Collector.Timeout to be 25s, got %v", cfg.Collector.Timeout)
	}
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Setup: Set environment variables
	t.Setenv("TEST_MIST_ORG_ID", "env-org-id")
	t.Setenv("TEST_MIST_API_KEY", "env-api-key")

	content := `
org_id: "${TEST_MIST_ORG_ID}"
mist_api:
  api_key: "${TEST_MIST_API_KEY}"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	// Test loading the config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() returned an unexpected error: %v", err)
	}

	if cfg.OrgId != "env-org-id" {
		t.Errorf("expected OrgId to be 'env-org-id', got %q", cfg.OrgId)
	}
	if cfg.MistClient.APIKey != "env-api-key" {
		t.Errorf("expected APIKey to be 'env-api-key', got %q", cfg.MistClient.APIKey)
	}
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Setup: Create a minimal config file
	content := `
mist_api:
  api_key: "my-api-key"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config file: %v", err)
	}

	// Test loading the config
	cfg, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig() returned an unexpected error: %v", err)
	}

	// Check that defaults are applied
	if cfg.Exporter.Address != defaultExporterAddress {
		t.Errorf("expected default Address to be %q, got %q", defaultExporterAddress, cfg.Exporter.Address)
	}
	if cfg.Exporter.Port != defaultExporterPort {
		t.Errorf("expected default Port to be %d, got %d", defaultExporterPort, cfg.Exporter.Port)
	}
	if cfg.Collector.Timeout != defaultCollectTimeout {
		t.Errorf("expected default Collector.Timeout to be %v, got %v", defaultCollectTimeout, cfg.Collector.Timeout)
	}
}

func TestLoadConfig_FileNotExist(t *testing.T) {
	_, err := LoadConfig("non-existent-file.yaml")
	if err == nil {
		t.Error("expected an error for non-existent file, but got nil")
	}
}
