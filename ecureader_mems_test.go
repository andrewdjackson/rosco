package rosco

import (
	"encoding/json"
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//
// To Run these tests you need Memsulator running!
//

type TestFixtures struct {
	Initialised    bool
	UseVirtualPort bool   `json:"useVirtualPort"`
	Port           string `json:"port"`
}

var testFixtures TestFixtures

func getFixtures() {
	var jsonFile *os.File
	var err error

	if !testFixtures.Initialised {
		testFixtures.UseVirtualPort = true

		if jsonFile, err = os.Open("testdata/fixtures.json"); err == nil {
			// read our opened jsonFile as a byte array.
			byteValue, _ := ioutil.ReadAll(jsonFile)

			// we unmarshal our byteArray which contains our
			// jsonFile's content into 'users' which we defined above
			json.Unmarshal(byteValue, &testFixtures)
		}

		if testFixtures.UseVirtualPort {
			testFixtures.Port = getVirtualPort()
		}

		testFixtures.Initialised = true
		defer jsonFile.Close()
	}
}

func getVirtualPort() string {
	homefolder, _ := homedir.Dir()
	return filepath.ToSlash(homefolder + "/ttyecu")
}

func Test_mems_Connect(t *testing.T) {
	getFixtures()

	r := NewMEMSReader(testFixtures.Port)
	connected, err := r.Connect()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, r.connected, is.True())
	then.AssertThat(t, connected, is.True())

	_ = r.Disconnect()

	r = NewMEMSReader(invalidPort)
	connected, err = r.Connect()

	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, r.connected, is.False())
	then.AssertThat(t, connected, is.False())

	_ = r.Disconnect()

	r = NewMEMSReader(loopbackPort)
	connected, err = r.Connect()

	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, r.connected, is.False())
	then.AssertThat(t, connected, is.False())
}

func Test_mems_connectToSerialPort(t *testing.T) {
	getFixtures()

	r := NewMEMSReader(testFixtures.Port)
	err := r.connectToSerialPort(testFixtures.Port)

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, r.serialPort, is.Not(is.Nil()))

	r = NewMEMSReader(testFixtures.Port)
	err = r.connectToSerialPort(invalidPort)

	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, r.serialPort, is.Nil())
}

func Test_mems_Disconnect(t *testing.T) {
	getFixtures()

	r := NewMEMSReader(testFixtures.Port)
	err := r.Disconnect()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, r.connected, is.False())
}

func Test_mems_SendAndReceive(t *testing.T) {
	getFixtures()

	r := NewMEMSReader(testFixtures.Port)

	connected, err := r.Connect()
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, connected, is.True())

	// expect echo of command
	response, err := r.SendAndReceive([]byte{0x0A})
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, response, is.EqualTo([]byte{0x0A}))

	// expect id string response
	response, err = r.SendAndReceive([]byte{0xD0})
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, response, is.EqualTo([]byte{0xD0, 0x99, 0x00, 0x02, 0x03}))
}

func Test_mems_commandMatchesResponse(t *testing.T) {
	getFixtures()

	r := NewMEMSReader(testFixtures.Port)
	err := r.commandMatchesResponse([]byte{0xca}, []byte{0xca})
	then.AssertThat(t, err, is.Nil())

	err = r.commandMatchesResponse([]byte{0xca}, []byte{0xca, 0x00})
	then.AssertThat(t, err, is.Nil())

	err = r.commandMatchesResponse([]byte{0xca}, []byte{0xd1, 0x00})
	then.AssertThat(t, err, is.Not(is.Nil()))
}
