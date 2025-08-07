package metrics

import (
	"github.com/gregwight/mistclient"
	"github.com/prometheus/client_golang/prometheus"
)

var DeviceLabelNames = append(SiteLabelNames,
	"device_id",
	"device_name",
	"device_type",
	"device_model",
	"device_hw_rev",
)

var DeviceWithRadioLabelNames = append(DeviceLabelNames,
	"radio",
)

func DeviceLabelValues(s mistclient.Site, ds mistclient.DeviceStat) []string {
	return append(SiteLabelValues(s),
		ds.ID,
		ds.Name,
		ds.Type.String(),
		ds.Model,
		ds.HwRev,
	)
}

func DeviceWithRadioLabelValues(s mistclient.Site, ds mistclient.DeviceStat, radio string) []string {
	return append(DeviceLabelValues(s, ds), radio)
}

var deviceMetrics *DeviceMetrics

type DeviceMetrics struct {
	accelX           *prometheus.GaugeVec
	accelY           *prometheus.GaugeVec
	accelZ           *prometheus.GaugeVec
	ambientTemp      *prometheus.GaugeVec
	attitude         *prometheus.GaugeVec
	cpuTemp          *prometheus.GaugeVec
	cpuUtil          *prometheus.GaugeVec
	createdTime      *prometheus.GaugeVec
	humidity         *prometheus.GaugeVec
	lastSeen         *prometheus.GaugeVec
	magneX           *prometheus.GaugeVec
	magneY           *prometheus.GaugeVec
	magneZ           *prometheus.GaugeVec
	modifiedTime     *prometheus.GaugeVec
	powerBudget      *prometheus.GaugeVec
	powerConstrained *prometheus.GaugeVec
	pressure         *prometheus.GaugeVec
	rxBps            *prometheus.GaugeVec
	status           *prometheus.GaugeVec
	txBps            *prometheus.GaugeVec
	uptime           *prometheus.GaugeVec
	vcoreVoltage     *prometheus.GaugeVec

	// Radio metrics
	bandwidth              *prometheus.GaugeVec
	channel                *prometheus.GaugeVec
	dynamicChainingEnabled *prometheus.GaugeVec
	noiseFloor             *prometheus.GaugeVec
	numClients             *prometheus.GaugeVec
	numWLANs               *prometheus.GaugeVec
	power                  *prometheus.GaugeVec
	rxBytes                *prometheus.GaugeVec
	rxPackets              *prometheus.GaugeVec
	txBytes                *prometheus.GaugeVec
	txPackets              *prometheus.GaugeVec
	utilAll                *prometheus.GaugeVec
	utilNonWiFi            *prometheus.GaugeVec
	utilRxInBSS            *prometheus.GaugeVec
	utilRxOtherBSS         *prometheus.GaugeVec
	utilTx                 *prometheus.GaugeVec
	utilUndecodableWiFi    *prometheus.GaugeVec
	utilUnknownWiFi        *prometheus.GaugeVec
}

func newDeviceMetrics(reg *prometheus.Registry) *DeviceMetrics {
	m := &DeviceMetrics{
		accelX: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "accel_x",
				Help:      "",
			}, DeviceLabelNames,
		),
		accelZ: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "accel_z",
				Help:      "",
			}, DeviceLabelNames,
		),
		accelY: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "accel_y",
				Help:      "",
			}, DeviceLabelNames,
		),
		ambientTemp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "ambient_temp",
				Help:      "",
			}, DeviceLabelNames,
		),
		attitude: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "attitude",
				Help:      "",
			}, DeviceLabelNames,
		),
		cpuTemp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "cpu_temp",
				Help:      "",
			}, DeviceLabelNames,
		),
		cpuUtil: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "cpu_util",
				Help:      "",
			}, DeviceLabelNames,
		),
		createdTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "created_time",
				Help:      "",
			}, DeviceLabelNames,
		),
		humidity: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "humidity",
				Help:      "",
			}, DeviceLabelNames,
		),
		lastSeen: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "last_seen",
				Help:      "",
			}, DeviceLabelNames,
		),
		magneX: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "magne_x",
				Help:      "",
			}, DeviceLabelNames,
		),
		magneY: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "magne_y",
				Help:      "",
			}, DeviceLabelNames,
		),
		magneZ: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "magne_z",
				Help:      "",
			}, DeviceLabelNames,
		),
		modifiedTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "modified_time",
				Help:      "",
			}, DeviceLabelNames,
		),
		powerBudget: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "power_budget",
				Help:      "",
			}, DeviceLabelNames,
		),
		powerConstrained: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "power_constrained",
				Help:      "",
			}, DeviceLabelNames,
		),
		pressure: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "pressure",
				Help:      "",
			}, DeviceLabelNames,
		),
		rxBps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "rx_bps",
				Help:      "",
			}, DeviceLabelNames,
		),
		status: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "status",
				Help:      "",
			}, DeviceLabelNames,
		),
		txBps: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "tx_bps",
				Help:      "",
			}, DeviceLabelNames,
		),
		uptime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "uptime",
				Help:      "",
			}, DeviceLabelNames,
		),
		vcoreVoltage: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "vcore_voltage",
				Help:      "",
			}, DeviceLabelNames,
		),

		// Radio metrics
		bandwidth: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "bandwidth",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		channel: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "channel",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		dynamicChainingEnabled: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "dynamic_chaining_enabled",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		noiseFloor: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "noise_floor",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		numClients: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "num_clients",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		numWLANs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "num_wlans",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		power: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "power",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		rxBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "rx_bytes",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		rxPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "rx_packets",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		txBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "tx_bytes",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		txPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "tx_packets",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilAll: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_all",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilNonWiFi: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_non_wifi",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilRxInBSS: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_rx_in_bss",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilRxOtherBSS: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_rx_other_bss",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilTx: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_tx",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilUndecodableWiFi: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_undecodable_wifi",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
		utilUnknownWiFi: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "mist",
				Subsystem: "device",
				Name:      "util_unknown_wifi",
				Help:      "",
			}, DeviceWithRadioLabelNames,
		),
	}

	reg.MustRegister(
		m.accelX,
		m.accelZ,
		m.accelY,
		m.ambientTemp,
		m.attitude,
		m.bandwidth,
		m.channel,
		m.cpuTemp,
		m.cpuUtil,
		m.createdTime,
		m.dynamicChainingEnabled,
		m.humidity,
		m.lastSeen,
		m.magneX,
		m.magneY,
		m.magneZ,
		m.modifiedTime,
		m.noiseFloor,
		m.numClients,
		m.numWLANs,
		m.power,
		m.powerBudget,
		m.powerConstrained,
		m.pressure,
		m.rxBps,
		m.rxBytes,
		m.rxPackets,
		m.status,
		m.txBps,
		m.txBytes,
		m.txPackets,
		m.uptime,
		m.utilAll,
		m.utilNonWiFi,
		m.utilRxInBSS,
		m.utilRxOtherBSS,
		m.utilTx,
		m.utilUndecodableWiFi,
		m.utilUnknownWiFi,
		m.vcoreVoltage,
	)

	return m
}

func handleSiteDeviceStat(site mistclient.Site, stat mistclient.DeviceStat) {
	labels := DeviceLabelValues(site, stat)

	deviceMetrics.accelX.WithLabelValues(labels...).Set(float64(stat.EnvStat.AccelX))
	deviceMetrics.accelY.WithLabelValues(labels...).Set(float64(stat.EnvStat.AccelY))
	deviceMetrics.accelZ.WithLabelValues(labels...).Set(float64(stat.EnvStat.AccelZ))
	deviceMetrics.ambientTemp.WithLabelValues(labels...).Set(float64(stat.EnvStat.AmbientTemp))
	deviceMetrics.attitude.WithLabelValues(labels...).Set(float64(stat.EnvStat.Attitude))
	deviceMetrics.cpuTemp.WithLabelValues(labels...).Set(float64(stat.EnvStat.CPUTemp))
	deviceMetrics.cpuUtil.WithLabelValues(labels...).Set(float64(stat.CPUUtil))
	deviceMetrics.createdTime.WithLabelValues(labels...).Set(float64(stat.CreatedTime.Unix()))
	deviceMetrics.humidity.WithLabelValues(labels...).Set(float64(stat.EnvStat.Humidity))
	deviceMetrics.lastSeen.WithLabelValues(labels...).Set(float64(stat.LastSeen.Unix()))
	deviceMetrics.magneX.WithLabelValues(labels...).Set(float64(stat.EnvStat.MagneX))
	deviceMetrics.magneY.WithLabelValues(labels...).Set(float64(stat.EnvStat.MagneY))
	deviceMetrics.magneZ.WithLabelValues(labels...).Set(float64(stat.EnvStat.MagneZ))
	deviceMetrics.modifiedTime.WithLabelValues(labels...).Set(float64(stat.ModifiedTime.Unix()))
	deviceMetrics.powerBudget.WithLabelValues(labels...).Set(float64(stat.PowerBudget))
	deviceMetrics.powerConstrained.WithLabelValues(labels...).Set(boolToFloat64(stat.PowerConstrained))
	deviceMetrics.pressure.WithLabelValues(labels...).Set(float64(stat.EnvStat.Pressure))
	deviceMetrics.rxBps.WithLabelValues(labels...).Set(float64(stat.RxBps))
	deviceMetrics.status.WithLabelValues(labels...).Set(float64(stat.Status))
	deviceMetrics.txBps.WithLabelValues(labels...).Set(float64(stat.TxBps))
	deviceMetrics.uptime.WithLabelValues(labels...).Set(float64(stat.Uptime))
	deviceMetrics.vcoreVoltage.WithLabelValues(labels...).Set(float64(stat.EnvStat.VcoreVoltage))

	// Radio metrics
	for radioConfig, radioStat := range stat.RadioStats {
		labels := DeviceWithRadioLabelValues(site, stat, radioConfig.String())

		deviceMetrics.bandwidth.WithLabelValues(labels...).Set(float64(radioStat.Bandwidth))
		deviceMetrics.channel.WithLabelValues(labels...).Set(float64(radioStat.Channel))
		deviceMetrics.dynamicChainingEnabled.WithLabelValues(labels...).Set(boolToFloat64(radioStat.DynamicChainingEnabled))
		deviceMetrics.noiseFloor.WithLabelValues(labels...).Set(float64(radioStat.NoiseFloor))
		deviceMetrics.numClients.WithLabelValues(labels...).Set(float64(radioStat.NumClients))
		deviceMetrics.numWLANs.WithLabelValues(labels...).Set(float64(radioStat.NumWLANs))
		deviceMetrics.power.WithLabelValues(labels...).Set(float64(radioStat.Power))
		deviceMetrics.rxBytes.WithLabelValues(labels...).Set(float64(radioStat.RxBytes))
		deviceMetrics.rxPackets.WithLabelValues(labels...).Set(float64(radioStat.RxPkts))
		deviceMetrics.txBytes.WithLabelValues(labels...).Set(float64(radioStat.TxBytes))
		deviceMetrics.txPackets.WithLabelValues(labels...).Set(float64(radioStat.TxPkts))
		deviceMetrics.utilAll.WithLabelValues(labels...).Set(float64(radioStat.UtilAll))
		deviceMetrics.utilNonWiFi.WithLabelValues(labels...).Set(float64(radioStat.UtilNonWiFi))
		deviceMetrics.utilRxInBSS.WithLabelValues(labels...).Set(float64(radioStat.UtilRxInBSS))
		deviceMetrics.utilRxOtherBSS.WithLabelValues(labels...).Set(float64(radioStat.UtilRxOtherBSS))
		deviceMetrics.utilTx.WithLabelValues(labels...).Set(float64(radioStat.UtilTx))
		deviceMetrics.utilUndecodableWiFi.WithLabelValues(labels...).Set(float64(radioStat.UtilUndecodableWiFi))
		deviceMetrics.utilUnknownWiFi.WithLabelValues(labels...).Set(float64(radioStat.UtilUnknownWiFi))
	}
}
