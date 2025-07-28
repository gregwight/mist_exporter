package metrics

import (
	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

type ClientMetrics struct {
	LastSeen    *prometheus.Desc
	Uptime      *prometheus.Desc
	Idletime    *prometheus.Desc
	PowerSaving *prometheus.Desc
	DualBand    *prometheus.Desc
	Channel     *prometheus.Desc
	RSSI        *prometheus.Desc
	SNR         *prometheus.Desc
	TxRate      *prometheus.Desc
	RxRate      *prometheus.Desc
	TxBytes     *prometheus.Desc
	RxBytes     *prometheus.Desc
	TxBps       *prometheus.Desc
	RxBps       *prometheus.Desc
	TxPackets   *prometheus.Desc
	RxPackets   *prometheus.Desc
	TxRetries   *prometheus.Desc
	RxRetries   *prometheus.Desc
}

var clientLabels = []string{
	"client_mac",
	"ap_id",
	"ap_mac",
	"username",
	"hostname",
	"os",
	"manufacture",
	"family",
	"model",
	"ssid",
}

func ClientLabels(c mistclient.Client) []string {
	return []string{
		c.Mac,
		c.APID,
		c.APMac,
		c.Username,
		c.Hostname,
		c.OS,
		c.Manufacture,
		c.Family,
		c.Model,
		c.SSID,
	}
}

func NewClientMetrics() *ClientMetrics {
	return &ClientMetrics{
		LastSeen: prometheus.NewDesc(
			"mist_client_last_seen",
			"Last time the client was seen",
			clientLabels,
			nil,
		),
		Uptime: prometheus.NewDesc(
			"mist_client_uptime",
			"Client connected time in seconds",
			clientLabels,
			nil,
		),
		Idletime: prometheus.NewDesc(
			"mist_client_idletime",
			"Client idle time in seconds (since the last RX packet)",
			clientLabels,
			nil,
		),
		PowerSaving: prometheus.NewDesc(
			"mist_client_power_saving",
			"Client in power-save mode",
			clientLabels,
			nil,
		),
		DualBand: prometheus.NewDesc(
			"mist_client_dual_band",
			"Client is dual-band capable",
			clientLabels,
			nil,
		),
		Channel: prometheus.NewDesc(
			"mist_client_channel",
			"Current channel",
			clientLabels,
			nil,
		),
		RSSI: prometheus.NewDesc(
			"mist_client_rssi",
			"Received Signal Strength Indicator",
			clientLabels,
			nil,
		),
		SNR: prometheus.NewDesc(
			"mist_client_snr",
			"Signal to Noise Ratio",
			clientLabels,
			nil,
		),
		TxRate: prometheus.NewDesc(
			"mist_client_tx_rate",
			"Transmit rate (Mbps)",
			clientLabels,
			nil,
		),
		RxRate: prometheus.NewDesc(
			"mist_client_rx_rate",
			"Receive rate (Mbps)",
			clientLabels,
			nil,
		),
		TxBytes: prometheus.NewDesc(
			"mist_client_tx_bytes",
			"Transmitted bytes",
			clientLabels,
			nil,
		),
		RxBytes: prometheus.NewDesc(
			"mist_client_rx_bytes",
			"Received bytes",
			clientLabels,
			nil,
		),
		TxBps: prometheus.NewDesc(
			"mist_client_tx_bps",
			"Transmit rate bits per second",
			clientLabels,
			nil,
		),
		RxBps: prometheus.NewDesc(
			"mist_client_rx_bps",
			"Receive rate bits per second",
			clientLabels,
			nil,
		),
		TxPackets: prometheus.NewDesc(
			"mist_client_tx_packets",
			"Transmitted packets",
			clientLabels,
			nil,
		),
		RxPackets: prometheus.NewDesc(
			"mist_client_rx_packets",
			"Received packets",
			clientLabels,
			nil,
		),
		TxRetries: prometheus.NewDesc(
			"mist_client_tx_retries",
			"Transmit retries",
			clientLabels,
			nil,
		),
		RxRetries: prometheus.NewDesc(
			"mist_client_rx_retries",
			"Received retries",
			clientLabels,
			nil,
		),
	}
}
