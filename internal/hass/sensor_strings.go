// Code generated by "stringer -type=SensorType,SensorDeviceClass,SensorStateClass -output sensor_strings.go -trimprefix Sensor"; DO NOT EDIT.

package hass

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeSensor-1]
	_ = x[TypeBinary-2]
}

const _SensorType_name = "TypeSensorTypeBinary"

var _SensorType_index = [...]uint8{0, 10, 20}

func (i SensorType) String() string {
	i -= 1
	if i < 0 || i >= SensorType(len(_SensorType_index)-1) {
		return "SensorType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _SensorType_name[_SensorType_index[i]:_SensorType_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Apparent_power-3]
	_ = x[Aqi-4]
	_ = x[Atmospheric_pressure-5]
	_ = x[SensorBattery-6]
	_ = x[Carbon_dioxide-7]
	_ = x[Carbon_monoxide-8]
	_ = x[Current-9]
	_ = x[Data_rate-10]
	_ = x[Data_size-11]
	_ = x[Date-12]
	_ = x[Distance-13]
	_ = x[Duration-14]
	_ = x[Energy-15]
	_ = x[Enum-16]
	_ = x[Frequency-17]
	_ = x[Gas-18]
	_ = x[Humidity-19]
	_ = x[Illuminance-20]
	_ = x[Irradiance-21]
	_ = x[Moisture-22]
	_ = x[Monetary-23]
	_ = x[Nitrogen_dioxide-24]
	_ = x[Nitrogen_monoxide-25]
	_ = x[Nitrous_oxide-26]
	_ = x[Ozone-27]
	_ = x[Pm1-28]
	_ = x[Pm25-29]
	_ = x[Pm10-30]
	_ = x[Power_factor-31]
	_ = x[SensorPower-32]
	_ = x[Precipitation-33]
	_ = x[Precipitation_intensity-34]
	_ = x[Pressure-35]
	_ = x[Reactive_power-36]
	_ = x[Signal_strength-37]
	_ = x[Sound_pressure-38]
	_ = x[Speed-39]
	_ = x[Sulphur_dioxide-40]
	_ = x[SensorTemperature-41]
	_ = x[Timestamp-42]
	_ = x[Volatile_organic_compounds-43]
	_ = x[Voltage-44]
	_ = x[Volume-45]
	_ = x[Water-46]
	_ = x[Weight-47]
	_ = x[Wind_speed-48]
}

const _SensorDeviceClass_name = "Apparent_powerAqiAtmospheric_pressureBatteryCarbon_dioxideCarbon_monoxideCurrentData_rateData_sizeDateDistanceDurationEnergyEnumFrequencyGasHumidityIlluminanceIrradianceMoistureMonetaryNitrogen_dioxideNitrogen_monoxideNitrous_oxideOzonePm1Pm25Pm10Power_factorPowerPrecipitationPrecipitation_intensityPressureReactive_powerSignal_strengthSound_pressureSpeedSulphur_dioxideTemperatureTimestampVolatile_organic_compoundsVoltageVolumeWaterWeightWind_speed"

var _SensorDeviceClass_index = [...]uint16{0, 14, 17, 37, 44, 58, 73, 80, 89, 98, 102, 110, 118, 124, 128, 137, 140, 148, 159, 169, 177, 185, 201, 218, 231, 236, 239, 243, 247, 259, 264, 277, 300, 308, 322, 337, 351, 356, 371, 382, 391, 417, 424, 430, 435, 441, 451}

func (i SensorDeviceClass) String() string {
	i -= 3
	if i < 0 || i >= SensorDeviceClass(len(_SensorDeviceClass_index)-1) {
		return "SensorDeviceClass(" + strconv.FormatInt(int64(i+3), 10) + ")"
	}
	return _SensorDeviceClass_name[_SensorDeviceClass_index[i]:_SensorDeviceClass_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Measurement-48]
	_ = x[Total-49]
	_ = x[Total_increasing-50]
}

const _SensorStateClass_name = "MeasurementTotalTotal_increasing"

var _SensorStateClass_index = [...]uint8{0, 11, 16, 32}

func (i SensorStateClass) String() string {
	i -= 48
	if i < 0 || i >= SensorStateClass(len(_SensorStateClass_index)-1) {
		return "SensorStateClass(" + strconv.FormatInt(int64(i+48), 10) + ")"
	}
	return _SensorStateClass_name[_SensorStateClass_index[i]:_SensorStateClass_index[i+1]]
}
