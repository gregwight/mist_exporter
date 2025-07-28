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

var deviceLabels = []string{
	"device_id",
	"device_name",
	"device_type",
	"device_model",
	"device_hw_rev",
	"site_id",
}

var deviceLabelsWithRadio = append(deviceLabels, "radio")

func DeviceStatLabels(ds mistclient.DeviceStat) []string {
	return []string{
		ds.ID,
		ds.Name,
		ds.Type.String(),
		ds.Model,
		ds.HwRev,
		ds.SiteID,
	}
}

func DeviceStatLabelsWithRadio(ds mistclient.DeviceStat, radio string) []string {
	return append(DeviceStatLabels(ds), radio)
}

func NewDeviceMetrics() *DeviceMetrics {
	return &DeviceMetrics{
		LastSeen: prometheus.NewDesc(
			"mist_device_last_seen",
			"Last time the device was seen",
			deviceLabels,
			nil,
		),
		Uptime: prometheus.NewDesc(
			"mist_device_uptime",
			"Device uptime in seconds",
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
			"Transmit rate bits per second",
			deviceLabels,
			nil,
		),
		RxBps: prometheus.NewDesc(
			"mist_device_rx_bps",
			"Receive rate bits per second",
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
			"Transmitted bytes",
			deviceLabelsWithRadio,
			nil,
		),
		RxBytes: prometheus.NewDesc(
			"mist_device_rx_bytes",
			"Received bytes",
			deviceLabelsWithRadio,
			nil,
		),
		TxPackets: prometheus.NewDesc(
			"mist_device_tx_packets",
			"Transmitted packets",
			deviceLabelsWithRadio,
			nil,
		),
		RxPackets: prometheus.NewDesc(
			"mist_device_rx_packets",
			"Received packets",
			deviceLabelsWithRadio,
			nil,
		),
		Power: prometheus.NewDesc(
			"mist_device_power",
			"Transmit power (in dBm)",
			deviceLabelsWithRadio,
			nil,
		),
		Channel: prometheus.NewDesc(
			"mist_device_channel",
			"Current channel",
			deviceLabelsWithRadio,
			nil,
		),
		Bandwidth: prometheus.NewDesc(
			"mist_device_bandwidth",
			"Current channel bandwidth, 20/40/80/160 MHz",
			deviceLabelsWithRadio,
			nil,
		),
		UtilAll: prometheus.NewDesc(
			"mist_device_util_all",
			"All utilization percent",
			deviceLabelsWithRadio,
			nil,
		),
		UtilTx: prometheus.NewDesc(
			"mist_device_util_tx",
			"Transmit utilization percent",
			deviceLabelsWithRadio,
			nil,
		),
		UtilRxInBSS: prometheus.NewDesc(
			"mist_device_util_rx_in_bss",
			"Reception of “In BSS” utilization in percentage, only frames that are received from AP/STAs within the BSS",
			deviceLabelsWithRadio,
			nil,
		),
		UtilRxOtherBSS: prometheus.NewDesc(
			"mist_device_util_rx_other_bss",
			"Eeception of “Other BSS” utilization in percentage, all frames received from AP/STAs that are outside the BSS",
			deviceLabelsWithRadio,
			nil,
		),
		UtilUnknownWiFi: prometheus.NewDesc(
			"mist_device_util_unknown_wifi",
			"Reception of “No Category” utilization in percentage, all 802.11 frames that are corrupted at the receiver",
			deviceLabelsWithRadio,
			nil,
		),
		UtilNonWiFi: prometheus.NewDesc(
			"mist_device_util_non_wifi",
			"Reception of “No Packets” utilization in percentage, received frames with invalid PLCPs and CRC glitches as noise",
			deviceLabelsWithRadio,
			nil,
		),
		UtilUndecodableWiFi: prometheus.NewDesc(
			"mist_device_util_undecodable_wifi",
			"Reception of “UnDecodable Wifi” utilization in percentage, only Preamble, PLCP header is decoded, Rest is undecodable in this radio",
			deviceLabelsWithRadio,
			nil,
		),
	}
}
