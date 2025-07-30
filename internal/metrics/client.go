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

var clientLabels = append(siteLabels,
	"client_mac",
	"client_username",
	"client_hostname",
	"client_os",
	"client_manufacture",
	"client_family",
	"client_model",
	"device_id",
	"ssid",
	"radio",
)

func ClientLabels(s mistclient.Site, c mistclient.Client) []string {
	return append(SiteLabels(s),
		c.Mac,
		c.Username,
		c.Hostname,
		c.OS,
		c.Manufacture,
		c.Family,
		c.Model,
		c.APID,
		c.SSID,
		c.Band.String(),
	)
}

func NewClientMetrics() *ClientMetrics {
	return &ClientMetrics{
		LastSeen: prometheus.NewDesc(
			"mist_client_last_seen",
			"Client last seen time",
			clientLabels,
			nil,
		),
		Uptime: prometheus.NewDesc(
			"mist_client_uptime",
			"Client connected time (s)",
			clientLabels,
			nil,
		),
		Idletime: prometheus.NewDesc(
			"mist_client_idletime",
			"Client idle time (s), since the last RX packet",
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
			"Client's current channel",
			clientLabels,
			nil,
		),
		RSSI: prometheus.NewDesc(
			"mist_client_rssi",
			"Client's received signal strength indicator (dBm)",
			clientLabels,
			nil,
		),
		SNR: prometheus.NewDesc(
			"mist_client_snr",
			"Client's signal to noise ratio",
			clientLabels,
			nil,
		),
		TxRate: prometheus.NewDesc(
			"mist_client_tx_rate",
			"Transmit rate to client (Mbps)",
			clientLabels,
			nil,
		),
		RxRate: prometheus.NewDesc(
			"mist_client_rx_rate",
			"Receive rate from client (Mbps)",
			clientLabels,
			nil,
		),
		TxBytes: prometheus.NewDesc(
			"mist_client_tx_bytes",
			"Bytes transmitted to client since connect",
			clientLabels,
			nil,
		),
		RxBytes: prometheus.NewDesc(
			"mist_client_rx_bytes",
			"Bytes received from client since connect",
			clientLabels,
			nil,
		),
		TxBps: prometheus.NewDesc(
			"mist_client_tx_bps",
			"Transmit rate to client (bps) last known",
			clientLabels,
			nil,
		),
		RxBps: prometheus.NewDesc(
			"mist_client_rx_bps",
			"Receive rate from client (bps) last known",
			clientLabels,
			nil,
		),
		TxPackets: prometheus.NewDesc(
			"mist_client_tx_packets",
			"Packets transmitted to client since connect",
			clientLabels,
			nil,
		),
		RxPackets: prometheus.NewDesc(
			"mist_client_rx_packets",
			"Packets received from client since connect",
			clientLabels,
			nil,
		),
		TxRetries: prometheus.NewDesc(
			"mist_client_tx_retries",
			"Number of transmit retries to client since connect",
			clientLabels,
			nil,
		),
		RxRetries: prometheus.NewDesc(
			"mist_client_rx_retries",
			"Number of receive retries from client since connect",
			clientLabels,
			nil,
		),
	}
}
