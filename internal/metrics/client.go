package metrics

import (
	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

var ClientLabelNames = append(SiteLabelNames,
	"client_mac",
	"client_username",
	"client_hostname",
	"client_os",
	"client_manufacture",
	"client_family",
	"client_model",
	"device_id",
	"proto",
	"radio",
	"ssid",
)

func ClientLabelValues(s mistclient.Site, c mistclient.Client) []string {
	return append(SiteLabelValues(s),
		c.Mac,
		c.Username,
		c.Hostname,
		c.OS,
		c.Manufacture,
		c.Family,
		c.Model,
		c.APID,
		c.Proto.String(),
		c.Band.String(),
		c.SSID,
	)
}

var clientMetrics *ClientMetrics

type ClientMetrics struct {
	channel        *prometheus.GaugeVec
	dualBand       *prometheus.GaugeVec
	idletime       *prometheus.GaugeVec
	isGuest        *prometheus.GaugeVec
	lastSeen       *prometheus.GaugeVec
	numLocatingAPs *prometheus.GaugeVec
	powerSaving    *prometheus.GaugeVec
	rssi           *prometheus.GaugeVec
	rxBps          *prometheus.GaugeVec
	rxBytes        *prometheus.GaugeVec
	rxPackets      *prometheus.GaugeVec
	rxRate         *prometheus.GaugeVec
	rxRetries      *prometheus.GaugeVec
	snr            *prometheus.GaugeVec
	txBps          *prometheus.GaugeVec
	txBytes        *prometheus.GaugeVec
	txPackets      *prometheus.GaugeVec
	txRate         *prometheus.GaugeVec
	txRetries      *prometheus.GaugeVec
	uptime         *prometheus.GaugeVec
}

func newClientMetrics(reg *prometheus.Registry) *ClientMetrics {
	m := &ClientMetrics{
		channel: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "channel",
				Help:      "",
			}, ClientLabelNames,
		),
		dualBand: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "dual_band",
				Help:      "",
			}, ClientLabelNames,
		),
		idletime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "idletime",
				Help:      "",
			}, ClientLabelNames,
		),
		isGuest: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "is_guest",
				Help:      "",
			}, ClientLabelNames,
		),
		lastSeen: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "last_seen",
				Help:      "",
			}, ClientLabelNames,
		),
		numLocatingAPs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "num_locating_aps",
				Help:      "",
			}, ClientLabelNames,
		),
		powerSaving: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "power_saving",
				Help:      "",
			}, ClientLabelNames,
		),
		rssi: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "rssi",
				Help:      "",
			}, ClientLabelNames,
		),
		rxBps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "rx_bps",
				Help:      "",
			}, ClientLabelNames,
		),
		rxBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "rx_bytes",
				Help:      "",
			}, ClientLabelNames,
		),
		rxPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "rx_packets",
				Help:      "",
			}, ClientLabelNames,
		),
		rxRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "rx_rate",
				Help:      "",
			}, ClientLabelNames,
		),
		rxRetries: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "rx_retries",
				Help:      "",
			}, ClientLabelNames,
		),
		snr: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "snr",
				Help:      "",
			}, ClientLabelNames,
		),
		txBps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "tx_bps",
				Help:      "",
			}, ClientLabelNames,
		),
		txBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "tx_bytes",
				Help:      "",
			}, ClientLabelNames,
		),
		txPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "tx_packets",
				Help:      "",
			}, ClientLabelNames,
		),
		txRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "tx_rate",
				Help:      "",
			}, ClientLabelNames,
		),
		txRetries: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "tx_retries",
				Help:      "",
			}, ClientLabelNames,
		),
		uptime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "client",
				Name:      "uptime",
				Help:      "",
			}, ClientLabelNames,
		),
	}

	reg.MustRegister(
		m.channel,
		m.dualBand,
		m.idletime,
		m.isGuest,
		m.lastSeen,
		m.numLocatingAPs,
		m.powerSaving,
		m.rssi,
		m.rxBps,
		m.rxBytes,
		m.rxPackets,
		m.rxRate,
		m.rxRetries,
		m.snr,
		m.txBps,
		m.txBytes,
		m.txPackets,
		m.txRate,
		m.txRetries,
		m.uptime,
	)

	return m
}

func handleSiteClientStat(site mistclient.Site, stat mistclient.Client) {
	labels := ClientLabelValues(site, stat)

	clientMetrics.channel.WithLabelValues(labels...).Set(float64(stat.Channel))
	clientMetrics.dualBand.WithLabelValues(labels...).Set(boolToFloat64(stat.DualBand))
	clientMetrics.idletime.WithLabelValues(labels...).Set(float64(stat.Idletime))
	clientMetrics.isGuest.WithLabelValues(labels...).Set(boolToFloat64(stat.IsGuest))
	clientMetrics.lastSeen.WithLabelValues(labels...).Set(float64(stat.LastSeen.Unix()))
	clientMetrics.numLocatingAPs.WithLabelValues(labels...).Set(float64(stat.NumLocatingAPs))
	clientMetrics.powerSaving.WithLabelValues(labels...).Set(boolToFloat64(stat.PowerSaving))
	clientMetrics.rssi.WithLabelValues(labels...).Set(float64(stat.RSSI))
	clientMetrics.rxBps.WithLabelValues(labels...).Set(float64(stat.RxBps))
	clientMetrics.rxBytes.WithLabelValues(labels...).Set(float64(stat.RxBytes))
	clientMetrics.rxPackets.WithLabelValues(labels...).Set(float64(stat.RxPackets))
	clientMetrics.rxRate.WithLabelValues(labels...).Set(float64(stat.RxRate))
	clientMetrics.rxRetries.WithLabelValues(labels...).Set(float64(stat.RxRetries))
	clientMetrics.snr.WithLabelValues(labels...).Set(float64(stat.SNR))
	clientMetrics.txBps.WithLabelValues(labels...).Set(float64(stat.TxBps))
	clientMetrics.txBytes.WithLabelValues(labels...).Set(float64(stat.TxBytes))
	clientMetrics.txPackets.WithLabelValues(labels...).Set(float64(stat.TxPackets))
	clientMetrics.txRate.WithLabelValues(labels...).Set(float64(stat.TxRate))
	clientMetrics.txRetries.WithLabelValues(labels...).Set(float64(stat.TxRetries))
	clientMetrics.uptime.WithLabelValues(labels...).Set(float64(stat.Uptime))
}
