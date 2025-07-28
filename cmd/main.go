package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/collector"
	"github.com/gregwight/mistexporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to the configuration file")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	// Initialize logger
	loggerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	if *debug {
		loggerOpts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, loggerOpts))

	// Load configuration
	config, err := config.LoadConfig(*configFile)
	if err != nil {
		logger.Error("unable to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize Mist API client
	client := mistclient.New(config.MistClient, logger)

	// Determine Mist OrgID
	var orgID string
	if config.OrgId != "" {
		orgID = config.OrgId
	} else {
		self, err := client.GetSelf()
		if err != nil {
			logger.Error("unable to retrieve self", "error", err)
			os.Exit(1)
		}

		for _, priv := range self.Privileges {
			if priv.Scope != "org" {
				continue
			}
			if orgID != "" {
				logger.Error("api key has access to multiple Mist organizations - please specify desired orgID using 'org_id' configurtaion key")
				os.Exit(1)
			}
			orgID = priv.OrgID
		}
	}

	// Create a pedantic registry
	registry := prometheus.NewPedanticRegistry()

	// Create and register the collector
	collector := collector.New(config.Collector, client, orgID, logger)
	registry.MustRegister(collector)

	// Add Go runtime metrics
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc(("/"), func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/metrics", http.StatusFound)
	})
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
		Registry:          registry,
		Timeout:           config.Collector.Timeout,
	}))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		configBytes, err := yaml.Marshal(config)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error marshaling config: %v", err)))
			return
		}
		w.Header().Set("Content-Type", "application/yaml")
		w.WriteHeader(http.StatusOK)
		w.Write(configBytes)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Exporter.Address, config.Exporter.Port),
		Handler: mux,
	}

	// Create context with signal handling
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	defer cancel()

	// Use errgroup for managing goroutines
	g, gCtx := errgroup.WithContext(ctx)

	// Start HTTP server
	g.Go(func() error {
		logger.Info("starting Mist Prometheus exporter", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTP server error: %w", err)
		}
		return nil
	})

	// Handle graceful shutdown
	g.Go(func() error {
		<-gCtx.Done()
		logger.Info("shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		return server.Shutdown(shutdownCtx)
	})

	// Wait for all goroutines to complete
	if err := g.Wait(); err != nil {
		logger.Error("server shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server shutdown complete")
}
