package rosco

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"reflect"
	"runtime"
	"testing"
)

const (
	invalidPort  = ""
	loopbackPort = "loopback"
	scenarioPort = "testdata/nofaults.csv"
)

func Test_ecureader_NewECUReader(t *testing.T) {
	r := NewECUReader("loopback:")
	then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&LoopbackReader{})))

	r = NewECUReader("mems:/dev/tty.serial")
	then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&MEMSReader{})))

	// test the FCR and CSV file extensions
	r = NewECUReader("filename.csv")
	then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&ScenarioReader{})))

	r = NewECUReader("filename.fcr")
	then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&ScenarioReader{})))

	// ensure only the extension determines the reader is a file reader
	r = NewECUReader("filenamefcr.file")
	then.AssertThat(t, reflect.TypeOf(r), is.Not(is.EqualTo(reflect.TypeOf(&ScenarioReader{}))))

	r = NewECUReader("filenamecsv.file")
	then.AssertThat(t, reflect.TypeOf(r), is.Not(is.EqualTo(reflect.TypeOf(&ScenarioReader{}))))

	// MEMSReader for serial ports
	r = NewECUReader("COM5")
	if runtime.GOOS == "darwin" {
		then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&LoopbackReader{})))
	} else {
		then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&MEMSReader{})))
	}

	r = NewECUReader("/dev/tty.Serial")
	if runtime.GOOS == "darwin" {
		then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&LoopbackReader{})))
	} else {
		then.AssertThat(t, reflect.TypeOf(r), is.EqualTo(reflect.TypeOf(&MEMSReader{})))
	}
}

func Test_ecureader_MemsReader(t *testing.T) {
	var connected bool
	var err error

	getFixtures()

	r := NewECUReader(testFixtures.Port)
	connected, err = r.Connect()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, connected, is.True())

	_, _ = r.SendAndReceive(MEMSReqData7D)
	_, _ = r.SendAndReceive(MEMSReqData80)

	r.Disconnect()
}

func Test_ecureader_LoopbackReader(t *testing.T) {
	var connected bool
	var err error

	r := NewECUReader("/loopback")
	connected, err = r.Connect()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, connected, is.True())
	r.Disconnect()
}

func Test_ecureader_createResponseMap(t *testing.T) {
	m := createResponseMap()
	then.AssertThat(t, len(m), is.GreaterThan(0))
	then.AssertThat(t, m["79"], is.EqualTo([]byte{0x79, 0x8b}))
	// unmapped command
	then.AssertThat(t, m["01"], is.EqualTo([]byte{0x01, 0x00}))
}

func Test_ecureader_getResponseSize(t *testing.T) {
	s, err := getResponseSize([]byte{0x79})
	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, s, is.EqualTo(2))

	// unmapped command, expect a default response sze of 2 bytes
	s, err = getResponseSize([]byte{0x20})
	then.AssertThat(t, err, is.Not(is.Nil()))
	then.AssertThat(t, s, is.EqualTo(2))
}
