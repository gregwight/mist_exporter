package metrics

import (
	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

// DeviceLabelNames defines the labels attached to device metrics.
var DeviceLabelNames = append(SiteLabelNames,
	"device_id",
	"device_name",
	"device_type",
	"device_model",
	"device_hw_rev",
)

// DeviceWithRadioLabelNames defines the labels attached to radio-specific device metrics.
var DeviceWithRadioLabelNames = append(DeviceLabelNames,
	"radio",
)

// DeviceLabelValues generates label values for device metrics.
func DeviceLabelValues(s mistclient.Site, ds mistclient.DeviceStat) []string {
	return append(SiteLabelValues(s),
		ds.ID,
		ds.Name,
		ds.Type.String(),
		ds.Model,
		ds.HwRev,
	)
}

// DeviceWithRadioLabelValues generates label values for radio-specific device metrics.
func DeviceWithRadioLabelValues(s mistclient.Site, ds mistclient.DeviceStat, radio string) []string {
	return append(DeviceLabelValues(s, ds), radio)
}

var deviceMetrics *DeviceMetrics

// DeviceMetrics holds metrics related to devices.
type DeviceMetrics struct {
	accelerationX      *prometheus.GaugeVec
	accelerationY      *prometheus.GaugeVec
	accelerationZ      *prometheus.GaugeVec
	ambientTemperature *prometheus.GaugeVec
	attitude           *prometheus.GaugeVec
	cpuTemperature     *prometheus.GaugeVec
	cpuUtilization     *prometheus.GaugeVec
	createdTimestamp   *prometheus.GaugeVec
	humidityPercent    *prometheus.GaugeVec
	lastSeenTimestamp  *prometheus.GaugeVec
	magnetometerX      *prometheus.GaugeVec
	magnetometerY      *prometheus.GaugeVec
	magnetometerZ      *prometheus.GaugeVec
	modifiedTimestamp  *prometheus.GaugeVec
	powerBudget        *prometheus.GaugeVec
	powerConstrained   *prometheus.GaugeVec
	pressure           *prometheus.GaugeVec
	receiveBps         *prometheus.GaugeVec
	status             *prometheus.GaugeVec
	transmitBps        *prometheus.GaugeVec
	uptimeSeconds      *prometheus.GaugeVec
	vcoreVoltage       *prometheus.GaugeVec

	// Radio metrics
	radioBandwidthMhz                      *prometheus.GaugeVec
	radioChannel                           *prometheus.GaugeVec
	radioDynamicChainingEnabled            *prometheus.GaugeVec
	radioNoiseFloorDbm                     *prometheus.GaugeVec
	radioClients                           *prometheus.GaugeVec
	radioWlans                             *prometheus.GaugeVec
	radioTransmitPowerDbm                  *prometheus.GaugeVec
	radioReceiveBytesTotal                 *prometheus.GaugeVec
	radioReceivePacketsTotal               *prometheus.GaugeVec
	radioTransmitBytesTotal                *prometheus.GaugeVec
	radioTransmitPacketsTotal              *prometheus.GaugeVec
	radioUtilizationAllPercent             *prometheus.GaugeVec
	radioUtilizationNonWifiPercent         *prometheus.GaugeVec
	radioUtilizationReceiveInBssPercent    *prometheus.GaugeVec
	radioUtilizationReceiveOtherBssPercent *prometheus.GaugeVec
	radioUtilizationTransmitPercent        *prometheus.GaugeVec
	radioUtilizationUndecodableWifiPercent *prometheus.GaugeVec
	radioUtilizationUnknownWifiPercent     *prometheus.GaugeVec
}

func newDeviceMetrics(reg *prometheus.Registry) *DeviceMetrics {
	m := &DeviceMetrics{
		accelerationX: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "acceleration_x_g",
				Help:      "Accelerometer reading on the X-axis in G-force.",
			}, DeviceLabelNames,
		),
		accelerationZ: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "acceleration_z_g",
				Help:      "Accelerometer reading on the Z-axis in G-force.",
			}, DeviceLabelNames,
		),
		accelerationY: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "acceleration_y_g",
				Help:      "Accelerometer reading on the Y-axis in G-force.",
			}, DeviceLabelNames,
		),
		ambientTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "ambient_temperature_celsius",
				Help:      "Ambient temperature measured by the device in Celsius.",
			}, DeviceLabelNames,
		),
		attitude: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "attitude_degrees",
				Help:      "Device attitude or orientation in degrees.",
			}, DeviceLabelNames,
		),
		cpuTemperature: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "cpu_temperature_celsius",
				Help:      "CPU temperature of the device in Celsius.",
			}, DeviceLabelNames,
		),
		cpuUtilization: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "cpu_utilization_percent",
				Help:      "Current CPU utilization of the device.",
			}, DeviceLabelNames,
		),
		createdTimestamp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "created_timestamp_seconds",
				Help:      "The time the device was created, as a Unix timestamp.",
			}, DeviceLabelNames,
		),
		humidityPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "humidity_percent",
				Help:      "Relative humidity percentage.",
			}, DeviceLabelNames,
		),
		lastSeenTimestamp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "last_seen_timestamp_seconds",
				Help:      "The last time the device was seen, as a Unix timestamp.",
			}, DeviceLabelNames,
		),
		magnetometerX: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "magnetometer_x_tesla",
				Help:      "Magnetometer reading on the X-axis in micro-Teslas.",
			}, DeviceLabelNames,
		),
		magnetometerY: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "magnetometer_y_tesla",
				Help:      "Magnetometer reading on the Y-axis in micro-Teslas.",
			}, DeviceLabelNames,
		),
		magnetometerZ: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "magnetometer_z_tesla",
				Help:      "Magnetometer reading on the Z-axis in micro-Teslas.",
			}, DeviceLabelNames,
		),
		modifiedTimestamp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "modified_timestamp_seconds",
				Help:      "The last time the device stats were modified, as a Unix timestamp.",
			}, DeviceLabelNames,
		),
		powerBudget: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "power_budget_watts",
				Help:      "The power budget of the device in watts.",
			}, DeviceLabelNames,
		),
		powerConstrained: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "power_constrained_status",
				Help:      "Whether the device is power constrained (1 for true, 0 for false).",
			}, DeviceLabelNames,
		),
		pressure: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "pressure_pascals",
				Help:      "Atmospheric pressure in Pascals.",
			}, DeviceLabelNames,
		),
		receiveBps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "receive_bits_per_second",
				Help:      "Bits per second received by the device.",
			}, DeviceLabelNames,
		),
		status: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "status_code",
				Help:      "The operational status of the device (e.g., 1 for connected, 0 for disconnected).",
			}, DeviceLabelNames,
		),
		transmitBps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "transmit_bits_per_second",
				Help:      "Bits per second transmitted by the device.",
			}, DeviceLabelNames,
		),
		uptimeSeconds: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "uptime_seconds",
				Help:      "Device uptime in seconds.",
			}, DeviceLabelNames,
		),
		vcoreVoltage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "vcore_voltage_volts",
				Help:      "VCore voltage of the device.",
			}, DeviceLabelNames,
		),

		// Radio metrics
		radioBandwidthMhz: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_bandwidth_mhz",
				Help:      "Radio channel bandwidth in MHz.",
			}, DeviceWithRadioLabelNames,
		),
		radioChannel: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_channel",
				Help:      "The current radio channel.",
			}, DeviceWithRadioLabelNames,
		),
		radioDynamicChainingEnabled: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_dynamic_chaining_enabled",
				Help:      "Whether dynamic chaining is enabled for the radio (1 for true, 0 for false).",
			}, DeviceWithRadioLabelNames,
		),
		radioNoiseFloorDbm: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_noise_floor_dbm",
				Help:      "The radio noise floor in dBm.",
			}, DeviceWithRadioLabelNames,
		),
		radioClients: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_clients_total",
				Help:      "Number of clients connected to this radio.",
			}, DeviceWithRadioLabelNames,
		),
		radioWlans: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_wlans_total",
				Help:      "Number of WLANs served by this radio.",
			}, DeviceWithRadioLabelNames,
		),
		radioTransmitPowerDbm: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_transmit_power_dbm",
				Help:      "The radio's transmit power in dBm.",
			}, DeviceWithRadioLabelNames,
		),
		radioReceiveBytesTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_receive_bytes_total",
				Help:      "Total bytes received by the radio.",
			}, DeviceWithRadioLabelNames,
		),
		radioReceivePacketsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_receive_packets_total",
				Help:      "Total packets received by the radio.",
			}, DeviceWithRadioLabelNames,
		),
		radioTransmitBytesTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_transmit_bytes_total",
				Help:      "Total bytes transmitted by the radio.",
			}, DeviceWithRadioLabelNames,
		),
		radioTransmitPacketsTotal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_transmit_packets_total",
				Help:      "Total packets transmitted by the radio.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationAllPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_all_percent",
				Help:      "Total radio channel utilization percentage.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationNonWifiPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_non_wifi_percent",
				Help:      "Radio channel utilization percentage by non-WiFi sources.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationReceiveInBssPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_receive_in_bss_percent",
				Help:      "Radio channel utilization percentage by receiving data in the same BSS.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationReceiveOtherBssPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_receive_other_bss_percent",
				Help:      "Radio channel utilization percentage by receiving data from other BSS.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationTransmitPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_transmit_percent",
				Help:      "Radio channel utilization percentage by transmitting data.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationUndecodableWifiPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_undecodable_wifi_percent",
				Help:      "Radio channel utilization percentage by undecodable WiFi sources.",
			}, DeviceWithRadioLabelNames,
		),
		radioUtilizationUnknownWifiPercent: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "radio_utilization_unknown_wifi_percent",
				Help:      "Radio channel utilization percentage by unknown WiFi sources.",
			}, DeviceWithRadioLabelNames,
		),
	}

	reg.MustRegister(
		m.accelerationX,
		m.accelerationZ,
		m.accelerationY,
		m.ambientTemperature,
		m.attitude,
		m.radioBandwidthMhz,
		m.radioChannel,
		m.cpuTemperature,
		m.cpuUtilization,
		m.createdTimestamp,
		m.radioDynamicChainingEnabled,
		m.humidityPercent,
		m.lastSeenTimestamp,
		m.magnetometerX,
		m.magnetometerY,
		m.magnetometerZ,
		m.modifiedTimestamp,
		m.radioNoiseFloorDbm,
		m.radioClients,
		m.radioWlans,
		m.radioTransmitPowerDbm,
		m.powerBudget,
		m.powerConstrained,
		m.pressure,
		m.receiveBps,
		m.radioReceiveBytesTotal,
		m.radioReceivePacketsTotal,
		m.status,
		m.transmitBps,
		m.radioTransmitBytesTotal,
		m.radioTransmitPacketsTotal,
		m.uptimeSeconds,
		m.radioUtilizationAllPercent,
		m.radioUtilizationNonWifiPercent,
		m.radioUtilizationReceiveInBssPercent,
		m.radioUtilizationReceiveOtherBssPercent,
		m.radioUtilizationTransmitPercent,
		m.radioUtilizationUndecodableWifiPercent,
		m.radioUtilizationUnknownWifiPercent,
		m.vcoreVoltage,
	)

	return m
}

func handleSiteDeviceStat(site mistclient.Site, stat mistclient.DeviceStat) {
	labels := DeviceLabelValues(site, stat)

	deviceMetrics.accelerationX.WithLabelValues(labels...).Set(float64(stat.EnvStat.AccelX))
	deviceMetrics.accelerationY.WithLabelValues(labels...).Set(float64(stat.EnvStat.AccelY))
	deviceMetrics.accelerationZ.WithLabelValues(labels...).Set(float64(stat.EnvStat.AccelZ))
	deviceMetrics.ambientTemperature.WithLabelValues(labels...).Set(float64(stat.EnvStat.AmbientTemp))
	deviceMetrics.attitude.WithLabelValues(labels...).Set(float64(stat.EnvStat.Attitude))
	deviceMetrics.cpuTemperature.WithLabelValues(labels...).Set(float64(stat.EnvStat.CPUTemp))
	deviceMetrics.cpuUtilization.WithLabelValues(labels...).Set(float64(stat.CPUUtil))
	deviceMetrics.createdTimestamp.WithLabelValues(labels...).Set(float64(stat.CreatedTime.Unix()))
	deviceMetrics.humidityPercent.WithLabelValues(labels...).Set(float64(stat.EnvStat.Humidity))
	deviceMetrics.lastSeenTimestamp.WithLabelValues(labels...).Set(float64(stat.LastSeen.Unix()))
	deviceMetrics.magnetometerX.WithLabelValues(labels...).Set(float64(stat.EnvStat.MagneX))
	deviceMetrics.magnetometerY.WithLabelValues(labels...).Set(float64(stat.EnvStat.MagneY))
	deviceMetrics.magnetometerZ.WithLabelValues(labels...).Set(float64(stat.EnvStat.MagneZ))
	deviceMetrics.modifiedTimestamp.WithLabelValues(labels...).Set(float64(stat.ModifiedTime.Unix()))
	deviceMetrics.powerBudget.WithLabelValues(labels...).Set(float64(stat.PowerBudget))
	deviceMetrics.powerConstrained.WithLabelValues(labels...).Set(boolToFloat64(stat.PowerConstrained))
	deviceMetrics.pressure.WithLabelValues(labels...).Set(float64(stat.EnvStat.Pressure))
	deviceMetrics.receiveBps.WithLabelValues(labels...).Set(float64(stat.RxBps))
	deviceMetrics.status.WithLabelValues(labels...).Set(float64(stat.Status))
	deviceMetrics.transmitBps.WithLabelValues(labels...).Set(float64(stat.TxBps))
	deviceMetrics.uptimeSeconds.WithLabelValues(labels...).Set(float64(stat.Uptime))
	deviceMetrics.vcoreVoltage.WithLabelValues(labels...).Set(float64(stat.EnvStat.VcoreVoltage))

	// Radio metrics
	for radioConfig, radioStat := range stat.RadioStats {
		labels := DeviceWithRadioLabelValues(site, stat, radioConfig.String())

		deviceMetrics.radioBandwidthMhz.WithLabelValues(labels...).Set(float64(radioStat.Bandwidth))
		deviceMetrics.radioChannel.WithLabelValues(labels...).Set(float64(radioStat.Channel))
		deviceMetrics.radioDynamicChainingEnabled.WithLabelValues(labels...).Set(boolToFloat64(radioStat.DynamicChainingEnabled))
		deviceMetrics.radioNoiseFloorDbm.WithLabelValues(labels...).Set(float64(radioStat.NoiseFloor))
		deviceMetrics.radioClients.WithLabelValues(labels...).Set(float64(radioStat.NumClients))
		deviceMetrics.radioWlans.WithLabelValues(labels...).Set(float64(radioStat.NumWLANs))
		deviceMetrics.radioTransmitPowerDbm.WithLabelValues(labels...).Set(float64(radioStat.Power))
		deviceMetrics.radioReceiveBytesTotal.WithLabelValues(labels...).Set(float64(radioStat.RxBytes))
		deviceMetrics.radioReceivePacketsTotal.WithLabelValues(labels...).Set(float64(radioStat.RxPkts))
		deviceMetrics.radioTransmitBytesTotal.WithLabelValues(labels...).Set(float64(radioStat.TxBytes))
		deviceMetrics.radioTransmitPacketsTotal.WithLabelValues(labels...).Set(float64(radioStat.TxPkts))
		deviceMetrics.radioUtilizationAllPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilAll))
		deviceMetrics.radioUtilizationNonWifiPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilNonWiFi))
		deviceMetrics.radioUtilizationReceiveInBssPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilRxInBSS))
		deviceMetrics.radioUtilizationReceiveOtherBssPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilRxOtherBSS))
		deviceMetrics.radioUtilizationTransmitPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilTx))
		deviceMetrics.radioUtilizationUndecodableWifiPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilUndecodableWiFi))
		deviceMetrics.radioUtilizationUnknownWifiPercent.WithLabelValues(labels...).Set(float64(radioStat.UtilUnknownWiFi))
	}
}
