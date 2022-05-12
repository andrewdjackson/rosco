package rosco

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func Test_mems19_connectToSerialPort(t *testing.T) {
	virtualPort := getVirtualPort()
	r := NewMEMS19Reader(virtualPort)
	connected, err := r.Connect()

	then.AssertThat(t, connected, is.True())
	then.AssertThat(t, err, is.Nil())
}
