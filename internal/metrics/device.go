package metrics

import (
	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

type DeviceMetrics struct {
	LastSeen            *prometheus.Desc
	Uptime              *prometheus.Desc
	WLANs               *prometheus.Desc
	Clients             *prometheus.Desc
	TxBps               *prometheus.Desc
	RxBps               *prometheus.Desc
	TxBytes             *prometheus.Desc
	RxBytes             *prometheus.Desc
	TxPackets           *prometheus.Desc
	RxPackets           *prometheus.Desc
	Power               *prometheus.Desc
	Channel             *prometheus.Desc
	Bandwidth           *prometheus.Desc
	UtilAll             *prometheus.Desc
	UtilTx              *prometheus.Desc
	UtilRxInBSS         *prometheus.Desc
	UtilRxOtherBSS      *prometheus.Desc
	UtilUnknownWiFi     *prometheus.Desc
	UtilNonWiFi         *prometheus.Desc
	UtilUndecodableWiFi *prometheus.Desc
}

var deviceLabels = append(siteLabels,
	"device_id",
	"device_name",
	"device_type",
	"device_model",
	"device_hw_rev",
)

var deviceLabelsWithRadio = append(deviceLabels,
	"radio",
)

func DeviceStatLabels(s mistclient.Site, ds mistclient.DeviceStat) []string {
	return append(SiteLabels(s),
		ds.ID,
		ds.Name,
		ds.Type.String(),
		ds.Model,
		ds.HwRev,
	)
}

func DeviceStatLabelsWithRadio(s mistclient.Site, ds mistclient.DeviceStat, radio string) []string {
	return append(DeviceStatLabels(s, ds), radio)
}

func NewDeviceMetrics() *DeviceMetrics {
	return &DeviceMetrics{
		LastSeen: prometheus.NewDesc(
			"mist_device_last_seen",
			"Device last seen time",
			deviceLabels,
			nil,
		),
		Uptime: prometheus.NewDesc(
			"mist_device_uptime",
			"Device uptime (s)",
			deviceLabels,
			nil,
		),
		WLANs: prometheus.NewDesc(
			"mist_device_wlans",
			"Number of WLANs assigned to the device",
			deviceLabels,
			nil,
		),
		TxBps: prometheus.NewDesc(
			"mist_device_tx_bps",
			"Device's transmit rate (bps)",
			deviceLabels,
			nil,
		),
		RxBps: prometheus.NewDesc(
			"mist_device_rx_bps",
			"Device's receive rate (bps)",
			deviceLabels,
			nil,
		),
		Clients: prometheus.NewDesc(
			"mist_device_clients",
			"Number of clients connected to the device",
			deviceLabelsWithRadio,
			nil,
		),
		TxBytes: prometheus.NewDesc(
			"mist_device_tx_bytes",
			"Device's transmitted bytes",
			deviceLabelsWithRadio,
			nil,
		),
		RxBytes: prometheus.NewDesc(
			"mist_device_rx_bytes",
			"Device's received bytes",
			deviceLabelsWithRadio,
			nil,
		),
		TxPackets: prometheus.NewDesc(
			"mist_device_tx_packets",
			"Device's transmitted packets",
			deviceLabelsWithRadio,
			nil,
		),
		RxPackets: prometheus.NewDesc(
			"mist_device_rx_packets",
			"Device's received packets",
			deviceLabelsWithRadio,
			nil,
		),
		Power: prometheus.NewDesc(
			"mist_device_power",
			"Device's transmit power (dBm)",
			deviceLabelsWithRadio,
			nil,
		),
		Channel: prometheus.NewDesc(
			"mist_device_channel",
			"Device's current channel",
			deviceLabelsWithRadio,
			nil,
		),
		Bandwidth: prometheus.NewDesc(
			"mist_device_bandwidth",
			"Device's current channel bandwidth, 20/40/80/160 MHz",
			deviceLabelsWithRadio,
			nil,
		),
		UtilAll: prometheus.NewDesc(
			"mist_device_util_all",
			"Device's all utilization (%)",
			deviceLabelsWithRadio,
			nil,
		),
		UtilTx: prometheus.NewDesc(
			"mist_device_util_tx",
			"Device's transmit utilization (%)",
			deviceLabelsWithRadio,
			nil,
		),
		UtilRxInBSS: prometheus.NewDesc(
			"mist_device_util_rx_in_bss",
			"Device's reception of “In BSS” utilization (%), only frames that are received from AP/STAs within the BSS",
			deviceLabelsWithRadio,
			nil,
		),
		UtilRxOtherBSS: prometheus.NewDesc(
			"mist_device_util_rx_other_bss",
			"Device's eeception of “Other BSS” utilization (%), all frames received from AP/STAs that are outside the BSS",
			deviceLabelsWithRadio,
			nil,
		),
		UtilUnknownWiFi: prometheus.NewDesc(
			"mist_device_util_unknown_wifi",
			"Device's reception of “No Category” utilization (%), all 802.11 frames that are corrupted at the receiver",
			deviceLabelsWithRadio,
			nil,
		),
		UtilNonWiFi: prometheus.NewDesc(
			"mist_device_util_non_wifi",
			"Device's reception of “No Packets” utilization (%), received frames with invalid PLCPs and CRC glitches as noise",
			deviceLabelsWithRadio,
			nil,
		),
		UtilUndecodableWiFi: prometheus.NewDesc(
			"mist_device_util_undecodable_wifi",
			"Device's reception of “UnDecodable Wifi” utilization (%), only Preamble, PLCP header is decoded, rest is undecodable in this radio",
			deviceLabelsWithRadio,
			nil,
		),
	}
}
