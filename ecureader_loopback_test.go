package rosco

import (
	"github.com/corbym/gocrest/is"
	"github.com/corbym/gocrest/then"
	"testing"
)

func Test_loopback_Connect(t *testing.T) {
	r := NewLoopbackReader()
	connected, err := r.Connect()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, r.connected, is.True())
	then.AssertThat(t, connected, is.True())

	_ = r.Disconnect()
}

func Test_loopback_Disconnect(t *testing.T) {
	r := NewLoopbackReader()
	err := r.Disconnect()

	then.AssertThat(t, err, is.Nil())
	then.AssertThat(t, r.connected, is.False())
}

func Test_loopback_SendAndReceive(t *testing.T) {
	r := NewLoopbackReader()

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
