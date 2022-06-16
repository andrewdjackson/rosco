package rosco

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
	"time"
)

func Test_analyseOperationalFaults(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		EngineRPM:                engineRunning,
		BatteryVoltage:           lowBattery,
		ManifoldAbsolutePressure: highIdleMAP,
		CoilTime:                 highCoilTime,
		LambdaStatus:             inactiveLambdaStatus,
		CoolantTemp:              coldEngineTemperature,
		IdleBasePosition:         goodIdleBasePosition,
		LambdaVoltage:            highestLambdaValue + 1,
		JackCount:                highestJackCount - 1,
		CrankshaftPositionSensor: goodCASPosition,
	}

	d.analyseOperationalFaults(data)
	then.AssertThat(t, d.Analysis.BatteryFault, is.True())
	then.AssertThat(t, d.Analysis.MapFault, is.True())
	then.AssertThat(t, d.Analysis.CoilFault, is.False())
	then.AssertThat(t, d.Analysis.O2SystemFault, is.True())
	then.AssertThat(t, d.Analysis.IsEngineIdleFault, is.True())
	then.AssertThat(t, d.Analysis.VacuumFault, is.True())
	then.AssertThat(t, d.Analysis.LambdaRangeFault, is.True())
	then.AssertThat(t, d.Analysis.IdleAirControlJackFault, is.False())
	then.AssertThat(t, d.Analysis.CrankshaftSensorFault, is.False())

	data = MemsData{
		EngineRPM:                engineRunning,
		BatteryVoltage:           goodBattery,
		ManifoldAbsolutePressure: goodIdleMap,
		CoilTime:                 highCoilTime,
		CoolantTemp:              warmEngineTemperature,
		IdleHot:                  lowIdleHot,
		IdleSpeedOffset:          maximumIdleOffset + 1,
		IACPosition:              invalidIACPosition,
		LambdaVoltage:            goodLambdaValue,
		JackCount:                highestJackCount,
		CrankshaftPositionSensor: invalidCASPosition,
	}

	d.analyseOperationalFaults(data)
	then.AssertThat(t, d.Analysis.BatteryFault, is.False())
	then.AssertThat(t, d.Analysis.MapFault, is.False())
	then.AssertThat(t, d.Analysis.CoilFault, is.True())
	then.AssertThat(t, d.Analysis.IdleHotFault, is.True())
	then.AssertThat(t, d.Analysis.IdleAirControlFault, is.True())
	then.AssertThat(t, d.Analysis.VacuumFault, is.False())
	then.AssertThat(t, d.Analysis.LambdaRangeFault, is.False())
	then.AssertThat(t, d.Analysis.IdleAirControlJackFault, is.True())
	then.AssertThat(t, d.Analysis.CrankshaftSensorFault, is.True())

	data = MemsData{
		EngineRPM:                engineRunning,
		BatteryVoltage:           goodBattery,
		ManifoldAbsolutePressure: goodIdleMap,
		CoilTime:                 goodCoilTime,
		LambdaStatus:             activeLambdaStatus,
	}

	d.analyseOperationalFaults(data)
	then.AssertThat(t, d.Analysis.BatteryFault, is.False())
	then.AssertThat(t, d.Analysis.MapFault, is.False())
	then.AssertThat(t, d.Analysis.CoilFault, is.False())
	then.AssertThat(t, d.Analysis.O2SystemFault, is.False())
}

func Test_isBatteryVoltageLow(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		BatteryVoltage: lowBattery,
	}

	result := d.isBatteryVoltageLow(data)
	then.AssertThat(t, result, is.True())

	data = MemsData{
		BatteryVoltage: goodBattery,
	}

	result = d.isBatteryVoltageLow(data)
	then.AssertThat(t, result, is.False())
}

func Test_isCoilFaulty(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		EngineRPM:      engineRunning,
		BatteryVoltage: goodBattery,
		CoilTime:       highCoilTime,
	}

	result := d.isCoilFaulty(data)
	then.AssertThat(t, result, is.True())

	data = MemsData{
		EngineRPM:      engineRunning,
		BatteryVoltage: goodBattery,
		CoilTime:       goodCoilTime,
	}

	result = d.isCoilFaulty(data)
	then.AssertThat(t, result, is.False())

	// battery low, high coil time ignored
	data = MemsData{
		EngineRPM:      engineRunning,
		BatteryVoltage: lowBattery,
		CoilTime:       highCoilTime,
	}

	result = d.isCoilFaulty(data)
	then.AssertThat(t, result, is.False())

	// faulty
	data = convertDataframeStringToMemsData("7D201014FF92003CFFFF01017A6300FF56FFFF30807FF9FF19401EC0264034C008", "801C04825AFF47FF278625001001000000208C6C000047069A10000080")

	result = d.isCoilFaulty(data)
	then.AssertThat(t, result, is.True())
}

func Test_isMAPHigh(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		EngineRPM:                engineRunning,
		ManifoldAbsolutePressure: goodIdleMap,
	}

	result := d.isMAPHigh(data)
	then.AssertThat(t, result, is.False())

	data = MemsData{
		EngineRPM:                engineRunning,
		ManifoldAbsolutePressure: highIdleMAP,
	}

	result = d.isMAPHigh(data)
	then.AssertThat(t, result, is.True())
}

func Test_isO2SystemActive(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		LambdaStatus: inactiveLambdaStatus,
	}

	result := d.isO2SystemActive(data)
	then.AssertThat(t, result, is.False())

	data = MemsData{
		LambdaStatus: activeLambdaStatus,
	}

	result = d.isO2SystemActive(data)
	then.AssertThat(t, result, is.True())
}

func Test_isEngineIdleFaulty(t *testing.T) {
	d := NewDataframeAnalysis(1)

	// engine off, no fault
	data := MemsData{
		EngineRPM:        engineStopped,
		CoolantTemp:      coldEngineTemperature,
		IdleBasePosition: highIdleBasePosition,
	}

	result := d.isEngineIdleFaulty(data)
	then.AssertThat(t, result, is.False())

	// idle below operating temp., no fault
	data = MemsData{
		EngineRPM:        engineRunning,
		CoolantTemp:      coldEngineTemperature,
		IdleBasePosition: highIdleBasePosition,
	}

	result = d.isEngineIdleFaulty(data)
	then.AssertThat(t, result, is.False())

	// idle below operating temp., faulty
	data = MemsData{
		EngineRPM:        engineRunning,
		CoolantTemp:      coldEngineTemperature,
		IdleBasePosition: goodIdleBasePosition,
	}

	result = d.isEngineIdleFaulty(data)
	then.AssertThat(t, result, is.True())

	// idle at operating temp., no fault
	data = MemsData{
		EngineRPM:        engineRunning,
		CoolantTemp:      warmEngineTemperature,
		IdleBasePosition: goodIdleBasePosition,
	}

	result = d.isEngineIdleFaulty(data)
	then.AssertThat(t, result, is.False())

	// idle at operating temp., fault
	data = MemsData{
		EngineRPM:        engineRunning,
		CoolantTemp:      warmEngineTemperature,
		IdleBasePosition: highIdleHot,
	}

	result = d.isEngineIdleFaulty(data)
	then.AssertThat(t, result, is.True())
}

func Test_isHotIdleFaulty(t *testing.T) {
	d := NewDataframeAnalysis(1)

	// engine cold, no fault
	data := MemsData{
		EngineRPM:   engineStopped,
		CoolantTemp: coldEngineTemperature,
	}

	result := d.isHotIdleFaulty(data)
	then.AssertThat(t, result, is.False())

	// engine warm, hot idle low, fault
	data = MemsData{
		EngineRPM:   1,
		CoolantTemp: 80,
		IdleHot:     5,
	}

	result = d.isHotIdleFaulty(data)
	then.AssertThat(t, result, is.True())

	// engine warm, hot idle high, fault
	data = MemsData{
		EngineRPM:   engineRunning,
		CoolantTemp: warmEngineTemperature,
		IdleHot:     highIdleHot,
	}

	result = d.isHotIdleFaulty(data)
	then.AssertThat(t, result, is.True())

	// engine warm, hot idle normal, no fault
	data = MemsData{
		EngineRPM:   engineRunning,
		CoolantTemp: warmEngineTemperature,
		IdleHot:     goodIdleHot,
	}

	result = d.isHotIdleFaulty(data)
	then.AssertThat(t, result, is.False())
}

func Test_isVacuumFaulty(t *testing.T) {
	d := NewDataframeAnalysis(1)

	// engine idle, high MAP, fault
	data := MemsData{
		EngineRPM:                engineRunning,
		ManifoldAbsolutePressure: highIdleMAP,
	}

	result := d.isVacuumFaulty(data)
	then.AssertThat(t, result, is.True())

	// engine idle, good MAP, no fault
	data = MemsData{
		EngineRPM:                engineRunning,
		ManifoldAbsolutePressure: goodIdleMap,
	}

	result = d.isVacuumFaulty(data)
	then.AssertThat(t, result, is.False())
}

func Test_isLambdaOutofRange(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: highestLambdaValue + 1,
	}

	result := d.isLambdaOutOfRange(data)
	then.AssertThat(t, result, is.True())

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: lowestLambdaValue - 1,
	}

	result = d.isLambdaOutOfRange(data)
	then.AssertThat(t, result, is.True())

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: goodLambdaValue,
	}

	result = d.isLambdaOutOfRange(data)
	then.AssertThat(t, result, is.False())
}

func Test_isJackCountHigh(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		JackCount: highestJackCount - 1,
	}

	result := d.isJackCountHigh(data)
	then.AssertThat(t, result, is.False())

	data = MemsData{
		JackCount: highestJackCount,
	}

	result = d.isJackCountHigh(data)
	then.AssertThat(t, result, is.True())
}

func Test_isCrankshaftPositionSensorFaulty(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		CrankshaftPositionSensor: invalidCASPosition,
	}

	result := d.isCrankshaftSensorFaulty(data)
	then.AssertThat(t, result, is.True())

	data = MemsData{
		CrankshaftPositionSensor: goodCASPosition,
	}

	result = d.isCrankshaftSensorFaulty(data)
	then.AssertThat(t, result, is.False())

	// faulty CPS
	data = convertDataframeStringToMemsData("7D201014FF92003CFFFF01017A6300FF56FFFF30807FF9FF19401EC0264034C008", "801C04825AFF47FF278625001001000000208C6C000047069A10000080")

	result = d.isCrankshaftSensorFaulty(data)
	then.AssertThat(t, result, is.True())
}

func Test_isLambdaFaulty(t *testing.T) {
	d := NewDataframeAnalysis(3)

	data := MemsData{
		Time:             "12:00:00.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    lowestLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		JackCount:        highestJackCount - 1,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:00:01.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		JackCount:        highestJackCount - 1,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:00:02.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    highestLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	result := d.isLambdaFaulty(data)
	then.AssertThat(t, result, is.False())

	d = NewDataframeAnalysis(3)

	data = MemsData{
		Time:             "12:00:00.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue - 50,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:00:01.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:00:02.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue + 50,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)
	result = d.isLambdaFaulty(data)
	then.AssertThat(t, result, is.False())

	d = NewDataframeAnalysis(3)

	data = MemsData{
		Time:             "12:00:00.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue - 50,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:01:00.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:01:31.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue + 50,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		JackCount:        highestJackCount - 1,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)
	result = d.isLambdaFaulty(data)
	then.AssertThat(t, result, is.True())

	d = NewDataframeAnalysis(3)

	data = MemsData{
		Time:             "12:00:00.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    lowestLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		JackCount:        highestJackCount - 1,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:01:00.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    goodLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		JackCount:        highestJackCount - 1,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)

	data = MemsData{
		Time:             "12:01:31.000",
		EngineRPM:        engineRunning,
		LambdaVoltage:    highestLambdaValue,
		CoolantTemp:      warmEngineTemperature,
		IntakeAirTemp:    goodIntakeTemperature,
		IdleBasePosition: goodIdleBasePosition,
		DTC5:             expectedDTC5,
		JackCount:        highestJackCount - 1,
		BatteryVoltage:   goodBattery,
	}

	d.Analyse(data)
	result = d.isLambdaFaulty(data)
	then.AssertThat(t, result, is.False())
}

func Test_isLambdaOscillating(t *testing.T) {
	d := NewDataframeAnalysis(3)

	data := MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: lowestLambdaValue,
	}

	d.addToDataset(data)

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: goodLambdaValue,
	}

	d.addToDataset(data)

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: highestLambdaValue,
	}

	d.addToDataset(data)

	result := d.isLambdaOscillating(data)
	then.AssertThat(t, result, is.True())

	d = NewDataframeAnalysis(3)

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: goodLambdaValue - 50,
	}

	d.addToDataset(data)

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: goodLambdaValue,
	}

	d.addToDataset(data)

	data = MemsData{
		EngineRPM:     engineRunning,
		LambdaVoltage: goodLambdaValue + 50,
	}

	d.addToDataset(data)
	result = d.isLambdaOscillating(data)
	then.AssertThat(t, result, is.False())
}

func Test_isThermostatFaulty(t *testing.T) {
	d := NewDataframeAnalysis(2)
	d.expectedTimeEngineWarm, _ = time.Parse("15:04:05.000", "12:01:50.000")

	// not yet reached warm, no fault
	data := MemsData{
		Time:        "12:00:00.000",
		EngineRPM:   engineRunning,
		CoolantTemp: 70,
	}

	d.addToDataset(data)

	data = MemsData{
		Time:        "12:00:11.000",
		EngineRPM:   engineRunning,
		CoolantTemp: 71,
	}

	result := d.isThermostatFaulty(data)
	then.AssertThat(t, result, is.False())

	d = NewDataframeAnalysis(2)
	d.expectedTimeEngineWarm, _ = time.Parse("15:04:05.000", "12:00:11.000")

	// reached warm, no fault
	data = MemsData{
		Time:        "12:00:00.000",
		EngineRPM:   engineRunning,
		CoolantTemp: lowestEngineWarmTemperature - 1,
	}

	d.addToDataset(data)

	data = MemsData{
		Time:        "12:00:11.000",
		EngineRPM:   engineRunning,
		CoolantTemp: warmEngineTemperature,
	}

	result = d.isThermostatFaulty(data)
	then.AssertThat(t, result, is.False())

	d = NewDataframeAnalysis(2)
	d.expectedTimeEngineWarm, _ = time.Parse("15:04:05.000", "12:01:50.000")

	// not reached warm,  fault
	data = MemsData{
		Time:        "12:00:00.000",
		EngineRPM:   engineRunning,
		CoolantTemp: warmEngineTemperature - 10,
	}

	d.addToDataset(data)

	data = MemsData{
		Time:        "12:01:51.000",
		EngineRPM:   engineRunning,
		CoolantTemp: lowestEngineWarmTemperature - 1,
	}

	result = d.isThermostatFaulty(data)
	then.AssertThat(t, result, is.True())
}

func Test_isIdleSpeedFaulty(t *testing.T) {
	d := NewDataframeAnalysis(1)

	data := MemsData{
		EngineRPM:        engineRunning,
		IdleBasePosition: 200,
	}

	d.addToDataset(data)
	result := d.isIdleSpeedFaulty(data)
	then.AssertThat(t, result, is.True())

	data = MemsData{
		EngineRPM:        engineRunning,
		IdleBasePosition: 50,
	}

	d.addToDataset(data)
	result = d.isIdleSpeedFaulty(data)
	then.AssertThat(t, result, is.False())
}

func convertDataframeStringToMemsData(dataframe7d string, dataframe80 string) MemsData {
	r := NewResponder()
	d7d := r.convertHexStringToByteArray(dataframe7d)
	d80 := r.convertHexStringToByteArray(dataframe80)
	ecu := NewECUReaderInstance()
	df80, _ := ecu.createDataframe80(d80)
	df7d, _ := ecu.createDataframe7D(d7d)
	return ecu.createMemsData(df80, df7d)
}
