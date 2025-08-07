package collector

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// testAPIServerHandler serves mock API responses from the testdata directory.
func testAPIServerHandler(t *testing.T, dataDir string) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token test-api-key" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		path := filepath.Join(dataDir, strings.TrimPrefix(r.URL.Path, "/"))
		body, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, r)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}

func TestNew(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	c := New(nil, "test-org", logger)

	if c == nil {
		t.Fatal("New() returned a nil collector")
	}
	if c.orgID != "test-org" {
		t.Errorf("New() orgID = %q, want %q", c.orgID, "test-org")
	}
	if c.logger == nil {
		t.Error("New() logger is nil")
	}
}

func TestCollect(t *testing.T) {
	testCases := []struct {
		name         string
		handler      http.HandlerFunc
		expectedFile string
		lint         bool
	}{
		{
			name:         "success",
			handler:      testAPIServerHandler(t, "testdata"),
			expectedFile: "testdata/success.prom",
			lint:         true,
		},
		{
			name: "api error on org endpoints",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/api/v1/orgs/") {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				testAPIServerHandler(t, "testdata")(w, r)
			}),
			expectedFile: "testdata/org_error.prom",
		},
		{
			name: "api error on site stats",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if strings.HasPrefix(r.URL.Path, "/api/v1/sites/") && strings.HasSuffix(r.URL.Path, "/stats") {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}
				testAPIServerHandler(t, "testdata")(w, r)
			}),
			expectedFile: "testdata/site_stats_error.prom",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(tc.handler)

			client, err := mistclient.New(&mistclient.Config{BaseURL: server.URL, APIKey: "test-api-key"}, nil)
			if err != nil {
				t.Fatalf("failed to create mist client: %v", err)
			}
			// We close the server inside the test run to ensure all client connections are terminated
			// before the next test case starts.
			t.Cleanup(server.Close)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			collector := New(client, "test-org-id", logger)

			expected, err := os.Open(tc.expectedFile)
			if err != nil {
				t.Fatalf("failed to open expected metrics file %s: %v", tc.expectedFile, err)
			}
			defer expected.Close()

			if err := testutil.CollectAndCompare(collector, expected); err != nil {
				t.Errorf("unexpected metrics collected:\n%v", err)
			}

			if tc.lint {
				problems, err := testutil.CollectAndLint(collector)
				if err != nil {
					t.Errorf("metric linting failed with an error: %v", err)
				}
				if len(problems) > 0 {
					t.Errorf("metric linting failed with problems:\n%v", problems)
				}
			}
		})
	}
}
