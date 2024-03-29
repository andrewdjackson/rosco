package rosco

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"reflect"
	"strings"
	"time"
)

// ECUReaderInstance communication structure for MEMS
type ECUReaderInstance struct {
	EcuReader   ECUReader
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

	ecu.EcuReader = NewECUReader(port)

	if connected, err = ecu.connectToECU(); err == nil {
		if connected {
			ecu.Status.Connected = true
			// get the ecu serial number and iac position
			ecu.Status.ECUSerial, err = ecu.getECUSerial()
			ecu.Status.IACPosition, err = ecu.GetIACPosition()
			ecu.Status.ECUID, err = ecu.getECUID()

			// assign the responder
			if reflect.TypeOf(ecu.EcuReader) == reflect.TypeOf(&ScenarioReader{}) {
				ecu.Responder = ecu.EcuReader.(*ScenarioReader).Responder
				log.Infof("ECUReader is a Scenario reader, assigning responder")
			}

			ecu.openLog()
		}
	}

	return ecu.Status.Connected, err
}

func (ecu *ECUReaderInstance) Disconnect() error {
	var err error

	if err = ecu.EcuReader.Disconnect(); err == nil {
		log.Info("disconnected ecu")
	} else {
		log.Warnf("error disconnecting (%s)", err)
	}

	ecu.Status.Connected = false
	ecu.closeLog()

	// save a scenario file
	_ = ecu.saveScenario()
	// reset status
	ecu.resetStatus()

	return err
}

// ResetDiagnostics clears and resets the diagnostic data
func (ecu *ECUReaderInstance) ResetDiagnostics() {
	// update the status
	log.Info("resetting ecu diagnostics")
	ecu.Diagnostics = NewDataframeAnalysis(20)
}

func (ecu *ECUReaderInstance) GetDataframes() (MemsData, error) {
	var err error
	var d80, d7d []byte
	var df80 DataFrame80
	var df7d DataFrame7d

	df := MemsData{}

	// read the raw dataframes
	log.Info("getting 0x7d and 0x80 dataframes")

	if d80, d7d, err = ecu.readRawDataFrames(); err == nil {
		// create the dataframes from the raw binary df
		if df80, err = ecu.createDataframe80(d80); err == nil {
			if df7d, err = ecu.createDataframe7D(d7d); err == nil {
				// build the Mems Dataframe using the raw df and applying the relevant adjustments and calculations
				df = ecu.createMemsDataframe(df80, df7d)
				// include the raw df converted into string format
				df.Dataframe80 = hex.EncodeToString(d80)

				if len(df.Dataframe80) != 58 {
					log.Warnf("dataframe 0x80 length exception, expected 29 (%s)", df.Dataframe80)
				}

				df.Dataframe7d = hex.EncodeToString(d7d)
				if len(df.Dataframe7d) != 66 {
					log.Warnf("dataframe 0x7D length exception, expected 33 (%s)", df.Dataframe7d)
				}

				df.Analytics = ecu.analyseMemsData(df)
				log.Infof("generated ecu df from dataframe (%+v)", df)
			}
		}

		ecu.writeToLog(df)
	}

	return df, err

}

func (ecu *ECUReaderInstance) analyseMemsData(df MemsData) AnalysisReport {
	ecu.Diagnostics.Analyse(df)
	return ecu.Diagnostics.Analysis
}

func (ecu *ECUReaderInstance) connectToECU() (bool, error) {
	return ecu.EcuReader.Connect()
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

func (ecu *ECUReaderInstance) readRawDataFrames() ([]byte, []byte, error) {
	var dferr error
	var err error
	var dataframe7d, dataframe80 []byte

	if dataframe80, err = ecu.EcuReader.SendAndReceive(MEMSReqData80); err != nil {
		dferr = fmt.Errorf("error recieving dataframe 0x80 (%s)", err)
		log.Errorf("%s", dferr)
	}

	// MEMS 1.3 doesn't support 7D command, so assuming we received the 80 dataframe
	// correctly, we'll supply default data
	if dataframe7d, err = ecu.EcuReader.SendAndReceive(MEMSReqData7D); err != nil {
		if dferr != nil {
			// 80 failed, so assume the worst
			dferr = fmt.Errorf("error recieving dataframe 0x7d (%s)", err)
			log.Errorf("%s", dferr)
		} else {
			// set 7dx03 to 13 (0d) to indicate this is MEMS 1.3
			// this is an unknown parameter and is always 0 from a MEMS 1.6
			log.Errorf("error recieving dataframe 0x7d with valid 0x80, assuming MEMS 1.3")
			dataframe7d = []byte{0x7d, 0x20, 0x10, 0x00, 0x0D, 0x92, 0x00, 0x57, 0xff, 0xff, 0x01, 0x00, 0x64, 0x64, 0x00, 0xff, 0x00, 0xff, 0xff, 0x30, 0x00, 0x7b, 0x69, 0xff, 0x19, 0x40, 0x1e, 0xc0, 0x26, 0x40, 0x34, 0xc0, 0x00}
		}
	}

	return dataframe80, dataframe7d, dferr
}

func (ecu *ECUReaderInstance) openLog() {
	// initialise logging
	if ecu.isMEMSReader() {
		ecu.dataLogger = NewMemsDataLogger(GetLogFolder(), ecu.Status.ECUID)
	}
}

func (ecu *ECUReaderInstance) closeLog() {
	if ecu.isMEMSReader() {
		if ecu.dataLogger != nil {
			ecu.dataLogger.Close()
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

func (ecu *ECUReaderInstance) saveScenario() error {
	var err error

	if ecu.isMEMSReader() {
		if ecu.dataLogger != nil {
			csvFile := ecu.dataLogger.Filename
			// save the log file as a scenario file
			fcrFile := strings.Replace(csvFile, ".csv", ".fcr", 1)

			// create a new scenario file
			s := NewScenarioFile(fcrFile)
			s.ECUID = ecu.Status.ECUID
			s.ECUSerial = ecu.Status.ECUSerial

			// convert the csv
			if err := s.ConvertLogToScenario(csvFile); err == nil {
				if err = s.Write(); err == nil {
					log.Infof("successfully saved scenario %s", fcrFile)
				} else {
					err := fmt.Errorf("error writing scenario file (%s)", err)
					log.Errorf("%s", err)
				}
			} else {
				err := fmt.Errorf("error converting scenario file (%s)", err)
				log.Errorf("%s", err)
			}
		} else {
			err := fmt.Errorf("error saving scenario file, data log not initialised")
			log.Warnf("%s", err)
		}
	}

	return err
}

func (ecu *ECUReaderInstance) isMEMSReader() bool {
	return reflect.TypeOf(ecu.EcuReader) == reflect.TypeOf(&MEMSReader{})
}
