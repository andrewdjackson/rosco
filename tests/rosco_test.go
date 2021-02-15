// +build standard

package tests

import (
	"github.com/andrewdjackson/rosco"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func TestGetScenarios(t *testing.T) {
	scenarios, err := rosco.GetScenarios()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, len(scenarios), is.GreaterThan(0))
}

func TestConnectAndInitialise(t *testing.T) {
	port := getPort(false)
	connectAndInitialise(t, port)
}

func TestConnectAndInitialiseScenario(t *testing.T) {
	port := getPort(true)
	connectAndInitialise(t, port)
}

func connectAndInitialise(t *testing.T, port string) {
	//homeFolder, _ := homedir.Dir()
	//logfolder := fmt.Sprintf("%s/memsfcr/logs", homeFolder)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	then.AssertThat(t, mems.Status.Connected, is.True())
	then.AssertThat(t, mems.Status.Initialised, is.True())

	mems.Disconnect()
	then.AssertThat(t, mems.Status.Connected, is.False())
	then.AssertThat(t, mems.Status.Initialised, is.False())
}

func TestStatusWithoutConnection(t *testing.T) {
	mems := rosco.NewMemsConnection(".")
	then.AssertThat(t, mems.Status.Connected, is.False())
}

func TestScenarioGetDataframe(t *testing.T) {
	port := getPort(true)
	getDataframe(t, port)
}

func TestGetDataframe(t *testing.T) {
	port := getPort(false)
	getDataframe(t, port)
}

func getDataframe(t *testing.T, port string) {
	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	data := mems.GetDataframes()

	then.AssertThat(t, data.BatteryVoltage, is.GreaterThanOrEqualTo(11.0))
	then.AssertThat(t, data.IdleSpeedOffset, is.EqualTo(10))
	then.AssertThat(t, mems.CommandResponse.Command, is.EqualTo(rosco.MEMSDataFrame))
	then.AssertThat(t, mems.CommandResponse.Response, is.EqualTo(rosco.MEMSDataFrame))
	then.AssertThat(t, mems.CommandResponse.MemsDataFrame.BatteryVoltage, is.GreaterThanOrEqualTo(11.0))
}

func TestStats(t *testing.T) {
	port := getPort(false)
	mems := rosco.NewMemsConnection(port)
	mems.ConnectAndInitialiseECU(port)

	mems.GetDataframes()
	mems.GetDataframes()

	mems.Diagnostics.Analyse()
	stats := mems.Diagnostics
	then.AssertThat(t, stats.Stats["LambdaVoltage"].Mean, is.GreaterThanOrEqualTo(0.0))
}

func TestAdjustSTFT(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	trim := mems.AdjustShortTermFuelTrim(1)
	then.AssertThat(t, trim, is.EqualTo(rosco.MEMSFuelTrimDefault+1))

	trim = mems.AdjustShortTermFuelTrim(-1)
	then.AssertThat(t, trim, is.EqualTo(rosco.MEMSFuelTrimDefault-1))

	trim = mems.AdjustShortTermFuelTrim(0)
	then.AssertThat(t, trim, is.EqualTo(rosco.MEMSFuelTrimDefault))
}

func TestResetAdjustments(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	success := mems.ResetAdjustments()
	then.AssertThat(t, success, is.True())
}

func TestResetECU(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	success := mems.ResetECU()
	then.AssertThat(t, success, is.True())
}

func TestClearFaults(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	success := mems.ClearFaults()
	then.AssertThat(t, success, is.True())
}

func TestGetIACPosition(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	success := mems.GetIACPosition()
	then.AssertThat(t, success, is.EqualTo(rosco.MEMSIACPositionDefault))
}

func TestIdleDecay(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	value := mems.AdjustIdleDecay(1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIdleDecayDefault+1))

	value = mems.AdjustIdleDecay(-1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIdleDecayDefault-1))

	value = mems.AdjustIdleDecay(0)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIdleDecayDefault))
}

func TestIdleSpeed(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	value := mems.AdjustIdleSpeed(1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIdleSpeedDefault+1))

	value = mems.AdjustIdleSpeed(-1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIdleSpeedDefault-1))

	value = mems.AdjustIdleSpeed(0)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIdleSpeedDefault))
}

func TestIgnitionAdvanceOffset(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	value := mems.AdjustIgnitionAdvanceOffset(1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIgnitionAdvanceOffsetDefault+1))

	value = mems.AdjustIgnitionAdvanceOffset(-1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIgnitionAdvanceOffsetDefault-1))

	value = mems.AdjustIgnitionAdvanceOffset(0)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIgnitionAdvanceOffsetDefault))
}

func TestIACPosition(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	value := mems.AdjustIACPosition(1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIACPositionDefault+1))

	value = mems.AdjustIACPosition(-1)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIACPositionDefault-1))

	value = mems.AdjustIACPosition(0)
	then.AssertThat(t, value, is.EqualTo(rosco.MEMSIACPositionDefault))
}

func TestHeartbeat(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	value := mems.SendHeartbeat()
	then.AssertThat(t, value, is.True())
	then.AssertThat(t, mems.CommandResponse.Command, is.EqualTo(rosco.MEMSHeartbeat))
	then.AssertThat(t, mems.CommandResponse.Response[0], is.EqualTo(rosco.MEMSHeartbeat[0]))
}

func TestActuators(t *testing.T) {
	port := getPort(false)

	mems := rosco.NewMemsConnection(".")
	mems.ConnectAndInitialiseECU(port)

	value := mems.TestFuelPump(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestFuelPump(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestPTCRelay(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestPTCRelay(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestACRelay(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestACRelay(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestPurgeValve(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestPurgeValve(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestO2Heater(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestO2Heater(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestBoostValve(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestBoostValve(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestInjectors(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestInjectors(false)
	then.AssertThat(t, value, is.True())

	value = mems.TestCoil(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestCoil(false)
	then.AssertThat(t, value, is.True())

	//value = mems.TestFan1(true)
	//then.AssertThat(t, value, is.True())

	//value = mems.TestFan1(false)
	//then.AssertThat(t, value, is.True())

	value = mems.TestFan2(true)
	then.AssertThat(t, value, is.True())

	value = mems.TestFan2(false)
	then.AssertThat(t, value, is.True())
}
