package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/yaml.v3"
)

func TestNew(t *testing.T) {
	cfg := &config.Config{
		Exporter: &config.Exporter{
			Address: "127.0.0.1",
			Port:    9876,
		},
		Collector: &config.Collector{
			CollectTimeout: 10 * time.Second,
		},
	}
	reg := prometheus.NewRegistry()

	srv, err := New(cfg, reg)
	if err != nil {
		t.Fatalf("New() returned an unexpected error: %v", err)
	}

	if srv == nil {
		t.Fatal("New() returned a nil server")
	}

	expectedAddr := "127.0.0.1:9876"
	if srv.Addr != expectedAddr {
		t.Errorf("New() server address = %q, want %q", srv.Addr, expectedAddr)
	}

	if srv.Handler == nil {
		t.Error("New() server handler is nil")
	}
}

func TestServerHandlers(t *testing.T) {
	cfg := &config.Config{
		Exporter: &config.Exporter{
			Address: "localhost",
			Port:    9090,
		},
		Collector: &config.Collector{
			CollectTimeout: 5 * time.Second,
		},
		MistClient: &mistclient.Config{
			BaseURL: "https://test.api.com",
			APIKey:  "supersecretapikey",
		},
	}
	reg := prometheus.NewRegistry()

	srv, err := New(cfg, reg)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	expectedConfig := *cfg
	expectedConfig.MistClient.APIKey = "*****"
	configBytes, err := yaml.Marshal(expectedConfig)
	if err != nil {
		t.Fatalf("failed to marshal expected config: %v", err)
	}
	expectedConfigYAML := "---\n" + string(configBytes)

	testCases := []struct {
		name           string
		path           string
		wantStatusCode int
		wantBody       string
		wantHeaders    map[string]string
		skipBodyCheck  bool
	}{
		{
			name:           "Root",
			path:           "/",
			wantStatusCode: http.StatusOK,
			wantBody:       string(indexHTML),
			wantHeaders:    map[string]string{"Content-Type": "text/html; charset=utf-8"},
		},
		{
			name:           "Health",
			path:           "/health",
			wantStatusCode: http.StatusOK,
			wantBody:       "OK",
		},
		{
			name:           "Config",
			path:           "/config",
			wantStatusCode: http.StatusOK,
			wantBody:       expectedConfigYAML,
			wantHeaders:    map[string]string{"Content-Type": "text/plain"},
		},
		{
			name:           "Metrics",
			path:           "/metrics",
			wantStatusCode: http.StatusOK,
			skipBodyCheck:  true, // Body is dynamic, just check it's not empty
			wantHeaders:    map[string]string{"Content-Type": "text/plain; version=0.0.4; charset=utf-8; escaping=underscores"},
		},
		{
			name:           "Not Found",
			path:           "/not-a-real-path",
			wantStatusCode: http.StatusOK,
			wantBody:       string(indexHTML),
			wantHeaders:    map[string]string{"Content-Type": "text/html; charset=utf-8"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()

			srv.Handler.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatusCode {
				t.Errorf("handler for %q returned wrong status code: got %v want %v", tc.path, rr.Code, tc.wantStatusCode)
			}

			for key, value := range tc.wantHeaders {
				if got := rr.Header().Get(key); got != value {
					t.Errorf("handler for %q returned wrong header %q: got %q want %q", tc.path, key, got, value)
				}
			}

			body := rr.Body.String()
			if tc.skipBodyCheck {
				if body == "" {
					t.Errorf("handler for %q returned an empty body", tc.path)
				}
				return
			}

			if tc.wantBody != "" && strings.TrimSpace(body) != strings.TrimSpace(tc.wantBody) {
				t.Errorf("handler for %q returned unexpected body:\ngot:\n%v\nwant:\n%v", tc.path, body, tc.wantBody)
			}
		})
	}
}
