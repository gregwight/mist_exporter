package server

import (
	"fmt"
	"net/http"

	"github.com/gregwight/mistexporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"

	_ "embed"
)

//go:embed index.html
var indexHTML []byte

// New creates a new HTTP server for the main exporter API.
func New(cfg *config.Config, reg *prometheus.Registry) (*http.Server, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
		Registry:          reg,
		Timeout:           cfg.Collector.CollectTimeout,
	}))
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/config", handleConfig(cfg))

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Exporter.Address, cfg.Exporter.Port),
		Handler: mux,
	}, nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(indexHTML)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func handleConfig(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localConfig := *cfg
		localConfig.MistClient.APIKey = "*****"

		configBytes, err := yaml.Marshal(localConfig)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error marshaling config: %v", err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(append([]byte("---\n"), configBytes...))
	}
}
