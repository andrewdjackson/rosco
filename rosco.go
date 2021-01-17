package rosco

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// MemsCommandResponse communication pair
type MemsCommandResponse struct {
	Command       []byte   `json:"Command"`
	Response      []byte   `json:"Response"`
	MemsDataFrame MemsData `json:"MemsData"`
}

// MemsConnection communication structure for MEMS
type MemsConnection struct {
	SerialPort      *serial.Port
	CommandResponse *MemsCommandResponse
	Diagnostics     *MemsDiagnostics
	responder       *Responder
	Status          *MemsConnectionStatus
}

// MemsConnectionStatus are we?
type MemsConnectionStatus struct {
	Emulated    bool   `json:"Emulated"`
	Connected   bool   `json:"Connected"`
	Initialised bool   `json:"Initialised"`
	ECUID       string `json:"ECUID"`
	IACPosition int    `json:"IACPosition"`
}

// NewMemsConnection creates a new mems structure
func NewMemsConnection() *MemsConnection {
	m := &MemsConnection{}
	//m.TxECU = make(chan MemsCommandResponse)
	//m.RxECU = make(chan MemsCommandResponse)
	m.CommandResponse = &MemsCommandResponse{}
	// engine diagnostics
	m.Diagnostics = NewMemsDiagnostics()
	// set status
	m.Status = &MemsConnectionStatus{}
	m.Status.Connected = false
	m.Status.Initialised = false
	m.Status.Emulated = false
	m.Status.ECUID = ""
	m.Status.IACPosition = m.Diagnostics.Analysis.IACPosition

	return m
}

// ConnectAndInitialiseECU connect and initialise the ECU
func (mems *MemsConnection) ConnectAndInitialiseECU(port string) {
	if mems.isScenario(port) {
		// emulate ECU if scenario file is supplied
		mems.Status.Emulated = true
		mems.responder = NewResponder()
	}

	if !mems.Status.Connected {
		mems.connect(port)
		if mems.Status.Connected {
			mems.initialise()
		}
	}

	// update status
	mems.Status.IACPosition = mems.Diagnostics.Analysis.IACPosition
}

// Disconnect from the ECU
func (mems *MemsConnection) Disconnect() MemsConnectionStatus {
	// close the connection
	_ = mems.SerialPort.Flush()
	_ = mems.SerialPort.Close()

	// update the status
	mems.Status.Connected = false
	mems.Status.Initialised = false
	mems.Status.Emulated = false
	mems.Status.ECUID = ""
	mems.Status.IACPosition = 0

	return *mems.Status
}

// ResetDiagnostics clears and resets the diagnostic data
func (mems *MemsConnection) ResetDiagnostics() {
	// update the status
	mems.Diagnostics = NewMemsDiagnostics()
}

// GetStatus returns the connection and ECU status
func (mems *MemsConnection) GetStatus() MemsConnectionStatus {
	return *mems.Status
}

// sendCommandAndWaitResponse sends a command and returns the response
func (mems *MemsConnection) sendCommandAndWaitResponse(cmd []byte) ([]byte, error) {
	mems.writeSerial(cmd)
	response, e := mems.readSerial()

	if e != nil {
		LogW.Printf("%s command send/receive fault %v", ECUResponseTrace, e)
	}

	return response, e
}

func (mems *MemsConnection) GetDataframes() MemsData {
	LogI.Printf("%s getting x7d and x80 dataframes", ECUCommandTrace)

	// read the raw dataframes
	d80, d7d, e := mems.readRawDataFrames()

	if e != nil {
		LogE.Printf("%s Unable to create memsdata, corrupt dataframes", ECUResponseTrace)
	}

	// populate the DataFrame structure for command 0x80
	r := bytes.NewReader(d80)
	var df80 DataFrame80

	if err := binary.Read(r, binary.BigEndian, &df80); err != nil {
		LogE.Printf("%s dataframe x80 binary.Read failed: %v", ECUCommandTrace, err)
	}

	// populate the DataFrame structure for command 0x7d
	r = bytes.NewReader(d7d)
	var df7d DataFrame7d

	if err := binary.Read(r, binary.BigEndian, &df7d); err != nil {
		LogE.Printf("%s dataframe x7d binary.Read failed: %v", ECUCommandTrace, err)
	}

	t := time.Now()

	// build the Mems Data frame using the raw data and applying the relevant
	// adjustments and calculations
	memsdata := MemsData{
		Time:                     t.Format("15:04:05.000"),
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
		IdleHot:                  int(df80.IdleHot) - 35,
		IACPosition:              int(df80.IacPosition),
		IdleSpeedDeviation:       int(df80.IdleSpeedDeviation),
		IgnitionAdvanceOffset80:  int(df80.IgnitionAdvanceOffset80),
		IgnitionAdvance:          (float32(df80.IgnitionAdvance) / 2) - 24,
		CoilTime:                 roundTo2DecimalPoints(float32(df80.CoilTime) * 0.002),
		CrankshaftPositionSensor: bool(df80.CrankshaftPositionSensor != 0),
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
		ShortTermFuelTrim:        int(df7d.ShortTermFuelTrim),
		FuelTrimCorrection:       int(df7d.ShortTermFuelTrim) - 100,
		CarbonCanisterPurgeValve: int(df7d.CarbonCanisterPurgeValve),
		DTC3:                     df7d.Dtc3,
		IdleBasePosition:         int(df7d.IdleBasePos),
		DTC4:                     df7d.Dtc4,
		IgnitionAdvanceOffset7d:  int(df7d.IgnitionAdvanceOffset7d) - 48,
		IdleSpeedOffset:          (int(df7d.IdleSpeedOffset) - 128) * 25,
		DTC5:                     df7d.Dtc5,
		JackCount:                int(df7d.JackCount),
		Dataframe80:              hex.EncodeToString(d80),
		Dataframe7d:              hex.EncodeToString(d7d),
	}

	// add the data for diagnostics
	mems.Diagnostics.Add(memsdata)
	mems.Diagnostics.Analyse()

	LogI.Printf("%s built mems dataframe %v", ECUCommandTrace, memsdata)

	return memsdata
}

func (mems *MemsConnection) SendHeartbeat() bool {
	return mems.clearState(MEMSHeartbeat)
}

// ResetAdjustments resets the adjustable values
func (mems *MemsConnection) ResetAdjustments() bool {
	return mems.clearState(MEMSResetAdj)
}

// ResetECU clears fault codes. resets adjustable values and learnt values
func (mems *MemsConnection) ResetECU() bool {
	return mems.clearState(MEMSResetECU)
}

// ClearFaults clears fault codes
func (mems *MemsConnection) ClearFaults() bool {
	return mems.clearState(MEMSClearFaults)
}

// GetIACPosition returns the current IAC position
func (mems *MemsConnection) GetIACPosition() int {
	var data []byte

	data, _ = mems.sendCommandAndWaitResponse(MEMSGetIACPosition)

	if len(data) > 1 {
		return int(data[1])
	} else {
		return MEMSIACPositionDefault
	}
}

// AdjustShortTermFuelTrim increments or decrements by the number of steps
func (mems *MemsConnection) AdjustShortTermFuelTrim(steps int) int {
	return mems.applyAdjustment(MEMSSTFTIncrement, MEMSSTFTDecrement, MEMSFuelTrimDefault, steps)
}

// AdjustLongTermFuelTrim increments or decrements by the number of steps
func (mems *MemsConnection) AdjustLongTermFuelTrim(steps int) int {
	return mems.applyAdjustment(MEMSLTFTIncrement, MEMSLTFTDecrement, MEMSFuelTrimDefault, steps)
}

// AdjustIdleDecay increments or decrements by the number  of steps
func (mems *MemsConnection) AdjustIdleDecay(steps int) int {
	return mems.applyAdjustment(MEMSIdleDecayIncrement, MEMSIdleDecayDecrement, MEMSIdleDecayDefault, steps)
}

// AdjustIdleSpeed increments or decrements by the number of steps
func (mems *MemsConnection) AdjustIdleSpeed(steps int) int {
	return mems.applyAdjustment(MEMSIdleSpeedIncrement, MEMSIdleSpeedDecrement, MEMSIdleSpeedDefault, steps)
}

// AdjustIgnitionAdvanceOffset increments or decrements by the number of steps
func (mems *MemsConnection) AdjustIgnitionAdvanceOffset(steps int) int {
	return mems.applyAdjustment(MEMSIgnitionAdvanceOffsetIncrement, MEMSIgnitionAdvanceOffsetDecrement, MEMSIgnitionAdvanceOffsetDefault, steps)
}

// AdjustIACPosition increments or decrements by the number of steps
func (mems *MemsConnection) AdjustIACPosition(steps int) int {
	return mems.applyAdjustment(MEMSIACIncrement, MEMSIACDecrement, MEMSIACPositionDefault, steps)
}

// TestFuelPump test
func (mems *MemsConnection) TestFuelPump(activate bool) bool {
	return mems.activateActuator(MEMSFuelPumpOn, MEMSFuelPumpOff, activate)
}

// PTCRelay test
func (mems *MemsConnection) TestPTCRelay(activate bool) bool {
	return mems.activateActuator(MEMSPTCRelayOn, MEMSPTCRelayOff, activate)
}

// ACRelay test
func (mems *MemsConnection) TestACRelay(activate bool) bool {
	return mems.activateActuator(MEMSACRelayOn, MEMSACRelayOff, activate)
}

// TestPurgeValve test
func (mems *MemsConnection) TestPurgeValve(activate bool) bool {
	return mems.activateActuator(MEMSPurgeValveOn, MEMSPurgeValveOff, activate)
}

// TestO2Heater test
func (mems *MemsConnection) TestO2Heater(activate bool) bool {
	return mems.activateActuator(MEMSO2HeaterOn, MEMSO2HeaterOff, activate)
}

// TestBoostValve test
func (mems *MemsConnection) TestBoostValve(activate bool) bool {
	return mems.activateActuator(MEMSBoostValveOn, MEMSBoostValveOff, activate)
}

// TestFan1 test
func (mems *MemsConnection) TestFan1(activate bool) bool {
	return mems.activateActuator(MEMSFan1On, MEMSFan1Off, activate)
}

// TestFan2 test
func (mems *MemsConnection) TestFan2(activate bool) bool {
	return mems.activateActuator(MEMSFan2On, MEMSFan2Off, activate)
}

// TestInjectors test, the activate state is ignored on this test
func (mems *MemsConnection) TestInjectors(activate bool) bool {
	return mems.activateActuator(MEMSTestInjectors, MEMSTestInjectors, activate)
}

// TestCoil test, the activate state is ignored on this test
func (mems *MemsConnection) TestCoil(activate bool) bool {
	return mems.activateActuator(MEMSFireCoil, MEMSFireCoil, activate)
}

//
// Private functions
//

// Increment or Decrement the adjustment by n steps
// Returns the final value of the adjustment
func (mems *MemsConnection) applyAdjustment(incrementCommand []byte, decrementCommand []byte, defaultValue int, steps int) int {
	var data []byte
	var cmd []byte

	// if the steps are positive then increment the adjustment
	// by n steps.
	// ignore all but the last value reading
	if steps > 0 {
		for step := 0; step < steps; step++ {
			cmd = incrementCommand
			data, _ = mems.sendCommandAndWaitResponse(cmd)
		}
	}

	// if the steps are negative then decrement the adjustment
	// by n steps.
	// ignore all but the last value reading
	if steps < 0 {
		for step := steps; step < 0; step++ {
			cmd = decrementCommand
			data, _ = mems.sendCommandAndWaitResponse(cmd)
		}
	}

	// ensure we have at least 1 byte returned
	// before returning the value
	if data != nil {
		if len(data) > 1 {
			if cmd[0] == data[0] {
				return int(data[1])
			}
		}
	}

	// the data returned was either invalid or the command echo (byte 0) in the response
	// didn't match the command sent.
	// return the default value for the adjuster
	return defaultValue
}

// Switches on or off the actuator
// Returns the success of the operation
func (mems *MemsConnection) activateActuator(activateCommand []byte, deactivateCommand []byte, activate bool) bool {
	var cmd []byte
	var data []byte

	if activate {
		cmd = activateCommand
	} else {
		cmd = deactivateCommand
	}

	data, _ = mems.sendCommandAndWaitResponse(cmd)

	if data != nil {
		if len(data) > 0 {
			LogI.Printf("--> %d == %d", data[0], cmd[0])
			return data[0] == cmd[0]
		}
	}

	return false
}

// Clears the state for the reset command
// Returns success of the operation
func (mems *MemsConnection) clearState(resetCommand []byte) bool {
	var data []byte

	data, _ = mems.sendCommandAndWaitResponse(resetCommand)

	if data != nil {
		if len(data) > 1 {
			return data[0] == resetCommand[0]
		}
	}

	return false
}

// connect to MEMS via serial port
func (mems *MemsConnection) connect(port string) {
	var err error
	var s *serial.Port

	// assume not connected or initialised
	mems.Status.Connected = false
	mems.Status.Initialised = false

	if mems.Status.Emulated {
		err = mems.responder.LoadScenario(port)
	} else {
		// connect to the ecu, timeout if we don't get data after a couple of seconds
		c := &serial.Config{Name: port, Baud: 9600, ReadTimeout: time.Millisecond * 2000}

		LogI.Println("opening ", port)

		s, err = serial.OpenPort(c)
		if s != nil {
			mems.SerialPort = s
		}
	}

	if err == nil {
		LogI.Println("connected to ", port)
		mems.Status.Connected = true
	} else {
		LogE.Printf("error opening port (%s)", err)
		mems.Status.Connected = false
		mems.Status.Initialised = false
	}
}

// check if the port is a CSV file, if so then a scenario emulation
// has been requested rather than a real serial connection
func (mems *MemsConnection) isScenario(port string) bool {
	return strings.HasSuffix(port, ".csv")
}

// checks the first byte of the response against the sent command
func (mems *MemsConnection) isCommandEcho() bool {
	return mems.CommandResponse.Command[0] == mems.CommandResponse.Response[0]
}

// initialises the connection to the ECU
// The initialisation sequence is as follows:
//
// 1. Send command CA (MEMS_InitCommandA)
// 2. Recieve response CA
// 3. Send command 75 (MEMS_InitCommandB)
// 4. Recieve response 75
// 5. Send request ECU ID command D0 (MEMS_InitECUID)
// 6. Recieve response D0 XX XX XX XX
//
func (mems *MemsConnection) initialise() {
	// assume not initialised
	mems.Status.Initialised = false

	if mems.Status.Emulated {
		mems.Status.Initialised = true
	} else {
		if mems.SerialPort != nil {
			_ = mems.SerialPort.Flush()

			mems.writeSerial(MEMSInitCommandA)
			_, _ = mems.readSerial()

			mems.writeSerial(MEMSInitCommandB)
			_, _ = mems.readSerial()

			mems.writeSerial(MEMSInitECUID)
			ECUID, _ := mems.readSerial()
			mems.Status.ECUID = fmt.Sprintf("%X", ECUID)

			// get the IAC position
			mems.writeSerial(MEMSGetIACPosition)
			response, _ := mems.readSerial()
			iac, _ := binary.Uvarint(response)
			mems.Diagnostics.Analysis.IACPosition = int(iac)

			mems.Status.Initialised = true
		}
	}
}

// readSerial read from MEMS
// read 1 byte at a time until we have all the expected bytes
func (mems *MemsConnection) readSerial() ([]byte, error) {
	var n int
	var e error

	size := mems.getResponseSize(mems.CommandResponse.Command)

	// serial read buffer
	b := make([]byte, size)

	//  data frame buffer
	data := make([]byte, 0)

	if mems.Status.Emulated {
		// emulate the response
		data = mems.responder.GetECUResponse(mems.CommandResponse.Command)
		LogI.Printf("%s data read for emulation %x", EmulatorTrace, data)
	} else {
		if mems.SerialPort != nil {
			// read all the expected bytes before returning the data
			for count := 0; count < size; {
				// wait for a response from MEMS
				n, _ = mems.SerialPort.Read(b)

				if n == 0 {
					LogW.Printf("serial port read error, timeout?")
					// drop out of loop, send back a 0x00 byte array response
					// this prevents the loop getting blocked on a read error
					count = size
					data = append(data, b...)
					e = errors.New("serial port read error")
				} else {
					// append the read bytes to the data frame
					data = append(data, b[:n]...)
				}

				// increment by the number of bytes read
				count = count + n
				if count > size {
					LogW.Printf("%s dataframe size mismatch (received %d, expected %d)", ECUResponseTrace, count, size)
					e = errors.New("size mismatch")
				}
			}
		}
	}

	LogI.Printf("%s received data from ECU [%d] < %x", ECUResponseTrace, n, data)
	mems.CommandResponse.Response = data

	if !mems.isCommandEcho() {
		LogW.Printf("%s expecting command echo (%x)\n", ECUResponseTrace, mems.CommandResponse.Command)
		e = errors.New("command mismatch")
	}

	return data, e
}

// writeSerial write to MEMS
func (mems *MemsConnection) writeSerial(data []byte) {
	if mems.Status.Emulated {
		LogI.Printf("%s data stored for emulation %x", EmulatorTrace, data)
		mems.CommandResponse.Command = data
	} else {
		if mems.SerialPort != nil {
			// save the sent command
			mems.CommandResponse.Command = data

			// write the response to the code reader
			n, e := mems.SerialPort.Write(data)

			if e != nil {
				LogE.Printf("%s error sending data to serial port (%s)", ECUCommandTrace, e)
			}

			if n > 0 {
				LogI.Printf("%s data sent to serial port %x", ECUCommandTrace, data)
			}
		}
	}
}

func roundTo2DecimalPoints(x float32) float32 {
	return float32(math.Round(float64(x)*100) / 100)
}

// readRawDataFrames reads dataframe 80 and then dataframe 7d as raw byte arrays
func (mems *MemsConnection) readRawDataFrames() ([]byte, []byte, error) {
	mems.writeSerial(MEMSReqData80)
	dataframe80, e := mems.readSerial()

	if e != nil {
		LogW.Printf("%s dataframe80 command send/receive fault %v", ECUResponseTrace, e)
	}

	mems.writeSerial(MEMSReqData7D)
	dataframe7d, e := mems.readSerial()

	if e != nil {
		LogW.Printf("%s dataframe7d command send/receive fault %v", ECUResponseTrace, e)
	}

	return dataframe80, dataframe7d, e
}

// getResponseSize returns the expected number of bytes for a given command
// The 'response' variable contains the formats for each command response pattern
// by default the response size is 2 bytes unless the command has a special format.
func (mems *MemsConnection) getResponseSize(command []byte) int {
	size := 2

	c := hex.EncodeToString(command)
	c = strings.ToUpper(c)
	r := responseMap[c]

	if r != nil {
		size = len(r)
	} else {
		r = responseMap["00"]
		copy(r[0:], command)
	}

	LogI.Printf("%s expecting %x (%d)", ECUResponseTrace, r, size)
	return size
}

// package init function
func init() {
	// Response formats for commands that do not respond with the format [COMMAND][VALUE]
	// Generally these are either part of the initialisation sequence or are ECU data frames
	responseMap["0A"] = []byte{0x0A}
	responseMap["CA"] = []byte{0xCA}
	responseMap["75"] = []byte{0x75}

	// Format for DataFrames starts with [Command Echo][Data Size][Data Bytes (28 for 0x80 and 32 for 0x7D)]
	responseMap["80"] = []byte{0x80, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B}
	responseMap["7D"] = []byte{0x7d, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}
	responseMap["D0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}

	// heatbeat
	responseMap["F4"] = []byte{0xf4, 0x00}

	// adjustments
	responseMap["79"] = []byte{0x79, 0x8b} // increment STFT (default is 138)
	responseMap["7A"] = []byte{0x7a, 0x89} // decrement STFT (default is 138)
	responseMap["7B"] = []byte{0x7b, 0x1f} // increment LTFT (default is 30)
	responseMap["7C"] = []byte{0x7c, 0x1d} // decrement LTFT (default is 30)
	responseMap["89"] = []byte{0x89, 0x24} // increment Idle Decay (default is 35)
	responseMap["8A"] = []byte{0x8a, 0x22} // decrement Idle Decay (default is 35)
	responseMap["91"] = []byte{0x91, 0x81} // increment Idle Speed  (default is 128)
	responseMap["92"] = []byte{0x92, 0x7f} // decrement Idle Speed (default is 128)
	responseMap["93"] = []byte{0x93, 0x81} // increment Ignition Advance Offset (default is 128)
	responseMap["94"] = []byte{0x94, 0x7f} // decrement Ignition Advance Offset (default is 128)
	responseMap["FD"] = []byte{0xfd, 0x81} // increment IAC (default is 128)
	responseMap["FE"] = []byte{0xfe, 0x7f} // decrement IAC (default is 128)

	//resets
	responseMap["0F"] = []byte{0x0f, 0x00} // clear all adjustments
	responseMap["CC"] = []byte{0xcc, 0x00} // clear faults
	responseMap["FA"] = []byte{0xfa, 0x00} // clear all computed and learnt settings
	responseMap["FB"] = []byte{0xfb, 0x80} // Idle Air Control position

	// actuators
	responseMap["11"] = []byte{0x11, 0x00} // fuel pump on
	responseMap["01"] = []byte{0x01, 0x00} // fuel pump off
	responseMap["12"] = []byte{0x12, 0x00} // ptc relay on
	responseMap["02"] = []byte{0x02, 0x00} // ptc relay off
	responseMap["13"] = []byte{0x13, 0x00} // ac relay on
	responseMap["03"] = []byte{0x03, 0x00} // ac relay off
	responseMap["18"] = []byte{0x18, 0x00} // purge valve on
	responseMap["08"] = []byte{0x08, 0x00} // purge vavle off
	responseMap["19"] = []byte{0x19, 0x00} // O2 heater on
	responseMap["09"] = []byte{0x09, 0x00} // O2 heater off
	responseMap["1B"] = []byte{0x1b, 0x00} // boost valve on
	responseMap["0B"] = []byte{0x0b, 0x00} // boost valve off
	responseMap["1D"] = []byte{0x1d}       // fan 1 on
	responseMap["0D"] = []byte{0x0d, 0x00} // fan 1 off
	responseMap["1E"] = []byte{0x1e}       // fan 2 on
	responseMap["0E"] = []byte{0x0e, 0x00} // fan 2 off
	responseMap["EF"] = []byte{0xef, 0x03} // test mpi injectors
	responseMap["F7"] = []byte{0xf7, 0x03} // test injectors
	responseMap["F8"] = []byte{0xf8, 0x02} // fire coil

	// unknown command responses
	responseMap["65"] = []byte{0x65, 0x00}
	responseMap["6D"] = []byte{0x6d, 0x00}
	responseMap["7E"] = []byte{0x7e, 0x08}
	responseMap["7F"] = []byte{0x7f, 0x05}
	responseMap["82"] = []byte{0x82, 0x09, 0x9E, 0x1D, 0x00, 0x00, 0x60, 0x05, 0xFF, 0xFF}
	responseMap["CB"] = []byte{0xcb, 0x00}
	responseMap["CD"] = []byte{0xcd, 0x01}
	responseMap["D2"] = []byte{0xd2, 0x02, 0x01, 0x00, 0x01}
	responseMap["D3"] = []byte{0xd3, 0x02, 0x01, 0x00, 0x02}
	responseMap["E7"] = []byte{0xe7, 0x02}
	responseMap["E8"] = []byte{0xe8, 0x05, 0x26, 0x01, 0x00, 0x01}
	responseMap["ED"] = []byte{0xed, 0x00}
	responseMap["EE"] = []byte{0xee, 0x00}
	responseMap["F0"] = []byte{0xf0, 0x05}
	responseMap["F3"] = []byte{0xf3, 0x00}
	responseMap["F5"] = []byte{0xf5, 0x00}
	responseMap["F6"] = []byte{0xf6, 0x00}
	responseMap["FC"] = []byte{0xfc, 0x00}

	// generic response, expect command and single byte response
	responseMap["00"] = []byte{0x00, 0x00}
}