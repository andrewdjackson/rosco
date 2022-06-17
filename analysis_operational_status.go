package rosco

import "time"

func (df *DataframeAnalysis) analyseOperationalStatus(data MemsData) {
	df.Analysis.IsEngineRunning = df.isEngineRunning(data)
	df.Analysis.IsEngineWarming = df.isEngineWarming(data)
	df.Analysis.IsAtOperatingTemp = df.isEngineWarm(data)
	df.Analysis.IsEngineIdle = df.isEngineIdle(data)
	df.Analysis.IsClosedLoop = df.isLoopClosed(data)
	df.Analysis.IsThrottleActive = df.isThrottleActive(data)
}

func (df *DataframeAnalysis) isEngineRunning(data MemsData) bool {
	engineRunning := data.EngineRPM > engineNotRunningRPM

	// set the engine start time
	if engineRunning && df.engineStartedAt.IsZero() {
		df.engineStartedAt, _ = time.Parse(timeFormat, data.Time)
	}

	return engineRunning
}

func (df *DataframeAnalysis) isEngineWarming(data MemsData) bool {
	return data.CoolantTemp < lowestEngineWarmTemperature
}

func (df *DataframeAnalysis) isEngineWarm(data MemsData) bool {
	return data.CoolantTemp >= lowestEngineWarmTemperature
}

func (df *DataframeAnalysis) isEngineIdle(data MemsData) bool {
	// engine is deemed to be at idle if the engine is running
	// and the angle of the throttle pot indicates the throttle is off
	// later MEMS ECUs use the throttle pot to determine the idle position

	return df.isEngineRunning(data) && data.ThrottleAngle <= defaultIdleThrottleAngle
}

func (df *DataframeAnalysis) isLoopClosed(data MemsData) bool {
	return data.ClosedLoop
}

func (df *DataframeAnalysis) isThrottleActive(data MemsData) bool {
	// throttle angle is not available on MEMS1.3, so swapped to using the pot value
	return data.ThrottleAngle > defaultIdleThrottleAngle || data.EngineRPM > highestIdleRPM
}

func (df *DataframeAnalysis) getThrottleAngle(data MemsData) int {
	return data.EngineRPM / throttleAngleFactor
}
