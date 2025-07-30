package collector

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"log/slog"

	"github.com/gregwight/mistclient"
	"github.com/gregwight/mistexporter/internal/config"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

// setupTestCollector creates a mock API server, a client, and a collector for testing.
func setupTestCollector(t *testing.T, mux *http.ServeMux) *MistCollector {
	t.Helper()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := mistclient.New(&mistclient.Config{
		BaseURL: server.URL,
		APIKey:  "test-api-key",
	}, slog.Default())

	collectorConfig := &config.Collector{
		Timeout: 30 * time.Second,
		Workers: 10,
	}

	return New(collectorConfig, client, "test-org-id", slog.Default())
}

func newTestMux(t *testing.T) *http.ServeMux {
	t.Helper()
	// Setup mock server
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/orgs/test-org-id/alarms/count", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/org_alarms.json")
	})
	mux.HandleFunc("/api/v1/orgs/test-org-id/tickets/count", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/org_tickets.json")
	})
	mux.HandleFunc("/api/v1/orgs/test-org-id/sites", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/org_sites.json")
	})
	mux.HandleFunc("/api/v1/sites/site-1/stats/devices", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/site_device_stats.json")
	})
	mux.HandleFunc("/api/v1/sites/site-1/stats/clients", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/site_clients.json")
	})
	return mux
}

func TestMistCollector_Collect(t *testing.T) {
	// Check if testdata files exist
	if _, err := os.Stat("testdata"); os.IsNotExist(err) {
		t.Skip("testdata directory not found, skipping integration test")
	}

	mux := newTestMux(t)

	collector := setupTestCollector(t, mux)

	t.Run("OrgMetrics", func(t *testing.T) {
		expected := `
			# HELP mist_org_alarms Number of alarms in the organization
			# TYPE mist_org_alarms gauge
			mist_org_alarms{alarm_type="AP_DISCONNECTED"} 2
			# HELP mist_org_site A site confured in the organization
			# TYPE mist_org_site gauge
			mist_org_site{country_code="US",site_id="site-1",site_name="Test Site"} 1
			# HELP mist_org_tickets Number of tickets in the organization
			# TYPE mist_org_tickets gauge
			mist_org_tickets{ticket_status="open"} 1
		`
		err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
			"mist_org_alarms", "mist_org_tickets", "mist_org_site")
		if err != nil {
			t.Errorf("unexpected collecting result for org metrics:\n%s", err)
		}
	})

	t.Run("DeviceMetrics", func(t *testing.T) {
		expected := `
			# HELP mist_device_channel Device's current channel
			# TYPE mist_device_channel gauge
			mist_device_channel{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_24",site_id="site-1",site_name="Test Site"} 6
			mist_device_channel{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_5",site_id="site-1",site_name="Test Site"} 36
			# HELP mist_device_clients Number of clients connected to the device
			# TYPE mist_device_clients gauge
			mist_device_clients{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_24",site_id="site-1",site_name="Test Site"} 5
			mist_device_clients{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_5",site_id="site-1",site_name="Test Site"} 10
			# HELP mist_device_last_seen Device last seen time
			# TYPE mist_device_last_seen gauge
			mist_device_last_seen{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",site_id="site-1",site_name="Test Site"} 1.6725312e+09
			# HELP mist_device_power Device's transmit power (dBm)
			# TYPE mist_device_power gauge
			mist_device_power{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_24",site_id="site-1",site_name="Test Site"} 17
			mist_device_power{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_5",site_id="site-1",site_name="Test Site"} 20
			# HELP mist_device_rx_bps Device's receive rate (bps)
			# TYPE mist_device_rx_bps gauge
			mist_device_rx_bps{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",site_id="site-1",site_name="Test Site"} 2000
			# HELP mist_device_rx_bytes Device's received bytes
			# TYPE mist_device_rx_bytes gauge
			mist_device_rx_bytes{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_24",site_id="site-1",site_name="Test Site"} 2048
			mist_device_rx_bytes{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_5",site_id="site-1",site_name="Test Site"} 8192
			# HELP mist_device_tx_bps Device's transmit rate (bps)
			# TYPE mist_device_tx_bps gauge
			mist_device_tx_bps{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",site_id="site-1",site_name="Test Site"} 1000
			# HELP mist_device_tx_bytes Device's transmitted bytes
			# TYPE mist_device_tx_bytes gauge
			mist_device_tx_bytes{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_24",site_id="site-1",site_name="Test Site"} 1024
			mist_device_tx_bytes{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",radio="band_5",site_id="site-1",site_name="Test Site"} 4096
			# HELP mist_device_uptime Device uptime (s)
			# TYPE mist_device_uptime gauge
			mist_device_uptime{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",site_id="site-1",site_name="Test Site"} 86400
			# HELP mist_device_wlans Number of WLANs assigned to the device
			# TYPE mist_device_wlans gauge
			mist_device_wlans{country_code="US",device_hw_rev="1.0",device_id="device-1",device_model="AP43",device_name="AP One",device_type="ap",site_id="site-1",site_name="Test Site"} 2
		`
		err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
			"mist_device_last_seen", "mist_device_uptime", "mist_device_wlans", "mist_device_tx_bps", "mist_device_rx_bps",
			"mist_device_clients", "mist_device_tx_bytes", "mist_device_rx_bytes", "mist_device_power", "mist_device_channel")
		if err != nil {
			t.Errorf("unexpected collecting result for device metrics:\n%s", err)
		}
	})

	t.Run("ClientMetrics", func(t *testing.T) {
		expected := `
			# HELP mist_client_last_seen Client last seen time
			# TYPE mist_client_last_seen gauge
			mist_client_last_seen{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} 1.672531201e+09
			# HELP mist_client_idletime Client idle time (s), since the last RX packet
			# TYPE mist_client_idletime gauge
			mist_client_idletime{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} 10
			# HELP mist_client_rssi Client's received signal strength indicator (dBm)
			# TYPE mist_client_rssi gauge
			mist_client_rssi{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} -60
			# HELP mist_client_rx_bytes Bytes received from client since connect
			# TYPE mist_client_rx_bytes gauge
			mist_client_rx_bytes{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} 1024
			# HELP mist_client_snr Client's signal to noise ratio
			# TYPE mist_client_snr gauge
			mist_client_snr{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} 30
			# HELP mist_client_tx_bytes Bytes transmitted to client since connect
			# TYPE mist_client_tx_bytes gauge
			mist_client_tx_bytes{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} 512
			# HELP mist_client_uptime Client connected time (s)
			# TYPE mist_client_uptime gauge
			mist_client_uptime{client_family="",client_hostname="test-client",client_mac="aa:bb:cc:dd:ee:ff",client_manufacture="",client_model="",client_os="macOS",client_username="",country_code="US",device_id="device-1",radio="2.4",site_id="site-1",site_name="Test Site",ssid="TestSSID"} 3600
		`
		err := testutil.CollectAndCompare(collector, strings.NewReader(expected),
			"mist_client_last_seen", "mist_client_uptime", "mist_client_idletime", "mist_client_rssi", "mist_client_snr", "mist_client_tx_bytes", "mist_client_rx_bytes")
		if err != nil {
			t.Errorf("unexpected collecting result for client metrics:\n%s", err)
		}
	})
}
