package rosco

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strings"
	"time"
)

// ECUReaderInstance communication structure for MEMS
type ECUReaderInstance struct {
	ecuReader   ECUReader
	dataLogger  *MemsDataLogger
	Status      *ECUStatus
	Diagnostics *DataframeAnalysis
	Responder   *ScenarioResponder
}

// NewECUReaderInstance creates a new mems structure
func NewECUReaderInstance() *ECUReaderInstance {
	m := &ECUReaderInstance{}
	m.Status = &ECUStatus{}
	m.Diagnostics = NewDataframeAnalysis(20)
	m.resetStatus()

	return m
}

func (ecu *ECUReaderInstance) ConnectAndInitialiseECU(port string) (bool, error) {
	var err error
	var connected bool

	ecu.ecuReader = NewECUReader(port)

	if reflect.TypeOf(ecu.ecuReader) == reflect.TypeOf(&ScenarioReader{}) {
		ecu.Responder = ecu.ecuReader.(*ScenarioReader).Responder
	}

	if connected, err = ecu.connectToECU(); err == nil {
		if connected {
			ecu.Status.Connected = true
			// get the ecu id, serial number and iac position
			ecu.Status.ECUID, err = ecu.getECUID()
			ecu.Status.ECUSerial, err = ecu.getECUSerial()
			ecu.Status.IACPosition, err = ecu.GetIACPosition()

			ecu.openLog()
		}
	}

	return ecu.Status.Connected, err
}

func (ecu *ECUReaderInstance) Disconnect() error {
	var err error

	if err = ecu.ecuReader.Disconnect(); err == nil {
		log.Info("disconnected ecu")
	} else {
		log.Warnf("error disconnecting (%s)", err)
	}

	ecu.resetStatus()
	ecu.closeLog()

	return err
}

// ResetDiagnostics clears and resets the diagnostic data
func (ecu *ECUReaderInstance) ResetDiagnostics() {
	// update the status
	log.Info("resetting ecu diagnostics")
	ecu.Diagnostics = NewDataframeAnalysis(20)
}

func (ecu *ECUReaderInstance) GetDataframes() MemsData {
	df := MemsData{}

	// read the raw dataframes
	log.Info("getting 0x7d and 0x80 dataframes")
	d80, d7d := ecu.readRawDataFrames()

	// create the dataframes from the raw binary df
	if df80, err := ecu.createDataframe80(d80); err == nil {
		if df7d, err := ecu.createDataframe7D(d7d); err == nil {
			// build the Mems Dataframe using the raw df and applying the relevant adjustments and calculations
			df = ecu.createMemsDataframe(df80, df7d)
			// include the raw df converted into string format
			df.Dataframe80 = hex.EncodeToString(d80)
			df.Dataframe7d = hex.EncodeToString(d7d)

			log.Infof("generated ecu df from dataframe (%+v)", df)

			ecu.Diagnostics.Analyse(df)
			df.Analytics = ecu.Diagnostics.Analysis
		}
	}

	ecu.writeToLog(df)

	return df

}

func (ecu *ECUReaderInstance) connectToECU() (bool, error) {
	return ecu.ecuReader.Connect()
}

func (ecu *ECUReaderInstance) createMemsDataframe(df80 DataFrame80, df7d DataFrame7d) MemsData {
	t := time.Now()

	memsdata := MemsData{
		Time:                     t.Format("2006-01-02 15:04:05.000"),
		EngineRPM:                int(df80.EngineRpm),
		CoolantTemp:              int(df80.CoolantTemp) - 55,
		AmbientTemp:              int(df80.AmbientTemp) - 55,
		IntakeAirTemp:            int(df80.IntakeAirTemp) - 55,
		FuelTemp:                 int(df80.FuelTemp) - 55,
		ManifoldAbsolutePressure: float32(df80.ManifoldAbsolutePressure),
		BatteryVoltage:           float32(df80.BatteryVoltage) / 10,
		ThrottlePotSensor:        roundTo2DecimalPoints(float32(df80.ThrottlePotSensor) * 0.02),
		IdleSwitch:               bool(df80.IdleSwitch&IdleSwitchActive != 0),
		AirconSwitch:             bool(df80.AirconSwitch != 0),
		ParkNeutralSwitch:        bool(df80.ParkNeutralSwitch != 0),
		DTC0:                     df80.Dtc0,
		DTC1:                     df80.Dtc1,
		IdleSetPoint:             int(df80.IdleSetPoint),
		IdleHot:                  int(df80.IdleHot), // was (hotidle - 35) but don't understand why this offset
		IACPosition:              int(df80.IacPosition),
		IdleSpeedDeviation:       int(df80.IdleSpeedDeviation),
		IgnitionAdvanceOffset80:  int(df80.IgnitionAdvanceOffset80),
		IgnitionAdvance:          (float32(df80.IgnitionAdvance) / 2) - 24,
		CoilTime:                 roundTo2DecimalPoints(float32(df80.CoilTime) * 0.002),
		CrankshaftPositionSensor: int(df80.CrankshaftPositionSensor),
		CoolantTempSensorFault:   bool(df80.Dtc0&CoolantSensorFaultCode != 0),
		IntakeAirTempSensorFault: bool(df80.Dtc0&AirSensorFaultCode != 0),
		FuelPumpCircuitFault:     bool(df80.Dtc1&FuelPumpFaultCode != 0),
		ThrottlePotCircuitFault:  bool(df80.Dtc1&ThrottlePotFaultCode != 0),
		IgnitionSwitch:           bool(df7d.IgnitionSwitch != 0),
		ThrottleAngle:            int(math.Round(float64(df7d.ThrottleAngle) * 6 / 10)),
		AirFuelRatio:             float32(df7d.AirFuelRatio) / 10,
		DTC2:                     df7d.Dtc2,
		LambdaVoltage:            int(df7d.LambdaVoltage) * 5,
		LambdaFrequency:          int(df7d.LambdaFrequency),
		LambdaDutycycle:          int(df7d.LambdaDutyCycle),
		LambdaStatus:             int(df7d.LambdaStatus),
		ClosedLoop:               bool(df7d.LoopIndicator != 0),
		LongTermFuelTrim:         int(df7d.LongTermFuelTrim) - 128,
		ShortTermFuelTrim:        int(df7d.ShortTermFuelTrim) - 100,
		FuelTrimCorrection:       int(df7d.ShortTermFuelTrim) - 100,
		CarbonCanisterPurgeValve: int(df7d.CarbonCanisterPurgeValve),
		DTC3:                     df7d.Dtc3,
		IdleBasePosition:         int(df7d.IdleBasePos),
		DTC4:                     df7d.Dtc4,
		IgnitionAdvanceOffset7d:  int(df7d.IgnitionAdvanceOffset7d) - 48,
		IdleSpeedOffset:          int(df7d.IdleSpeedOffset), // - 128) * 25,
		DTC5:                     df7d.Dtc5,
		JackCount:                int(df7d.JackCount),
	}

	return memsdata
}

func (ecu *ECUReaderInstance) createDataframe7D(d7d []byte) (DataFrame7d, error) {
	var err error
	var df7d DataFrame7d

	defer func() {
		if err := recover(); err != nil {
			log.Warnf("dataframe conversion panic occurred %s", err)
		}
	}()

	// populate the DataFrame structure for command 0x7d
	byteReader := bytes.NewReader(d7d)

	if err = binary.Read(byteReader, binary.BigEndian, &df7d); err != nil {
		log.Errorf("error reading dataframe x7d (%s)", err)
	} else {
		log.Infof("dataframe x7d received (data: %X dataframe: %+v)", byteReader, df7d)
	}

	return df7d, err
}

func (ecu *ECUReaderInstance) createDataframe80(d80 []byte) (DataFrame80, error) {
	var err error
	var df80 DataFrame80

	defer func() {
		if err := recover(); err != nil {
			log.Warnf("dataframe conversion panic occurred %s", err)
		}
	}()

	// populate the DataFrame structure for command 0x80
	byteReader := bytes.NewReader(d80)

	if err = binary.Read(byteReader, binary.BigEndian, &df80); err != nil {
		log.Errorf("error reading dataframe x80 (%s)", err)
	} else {
		log.Infof("dataframe x80 received (data: %X dataframe: %+v)", byteReader, df80)
	}

	return df80, err
}

func (ecu *ECUReaderInstance) readRawDataFrames() ([]byte, []byte) {
	var err error
	var dataframe7d, dataframe80 []byte

	if dataframe80, err = ecu.ecuReader.SendAndReceive(MEMSReqData80); err != nil {
		log.Errorf("error recieving dataframe 0x80 (%s)", err)
	}
	if dataframe7d, err = ecu.ecuReader.SendAndReceive(MEMSReqData7D); err != nil {
		log.Errorf("error recieving dataframe 0x7D (%s)", err)
	}

	return dataframe80, dataframe7d
}

func (ecu *ECUReaderInstance) openLog() {
	// initialise logging
	if ecu.isMEMSReader() {
		ecu.dataLogger = NewMemsDataLogger(GetLogFolder(), ecu.Status.ECUSerial)
	}
}

func (ecu *ECUReaderInstance) closeLog() {
	if ecu.isMEMSReader() {
		if ecu.dataLogger != nil {
			ecu.dataLogger.Close()
			// save a scenario file
			ecu.SaveScenario(ecu.dataLogger.Filename)
		}
	}
}

func (ecu *ECUReaderInstance) writeToLog(df MemsData) {
	if ecu.dataLogger != nil {
		if ecu.isMEMSReader() {
			// write to a logfile if the ecu reader is a real (or virtual) ECU
			go ecu.dataLogger.WriteMemsDataToFile(df)
		}
	}
}

func (ecu *ECUReaderInstance) SaveScenario(csvFile string) {
	// save the log file as a scenario file
	fcrFile := strings.Replace(csvFile, ".csv", ".fcr", 1)

	// create a new scenario file
	s := NewScenarioFile(fcrFile)
	s.ECUID = ecu.Status.ECUID
	s.ECUSerial = ecu.Status.ECUSerial

	// convert the csv
	err := s.ConvertLogToScenario(csvFile)
	if err == nil {
		err = s.Write()
		log.Infof("saved scenario %s", fcrFile)
	}
}

func (ecu *ECUReaderInstance) isMEMSReader() bool {
	return reflect.TypeOf(ecu.ecuReader) == reflect.TypeOf(&MEMSReader{})
}
