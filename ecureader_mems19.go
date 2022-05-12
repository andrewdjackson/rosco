package rosco

import (
	"github.com/distributed/sers"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	ecu19SpecificInitCommand  = []byte{0x7C}
	ecu19WokeResponse         = []byte{0x55, 0x76, 0x83}
	ecu19SpecificInitResponse = []byte{ecu19SpecificInitCommand[0], 0xE9}
)

type MEMS19Reader struct {
	serialPort sers.SerialPort
	ecuReader  *MEMSReader
}

func NewMEMS19Reader(connection string) *MEMS19Reader {
	log.Infof("created mems ecu reader")

	r := &MEMS19Reader{}
	r.ecuReader = &MEMSReader{}
	r.ecuReader.port = connection

	return r
}

func (r *MEMS19Reader) Connect() (bool, error) {
	var connected bool
	var err error

	if err = r.wakeUpConnect(); err == nil {
		if err = r.wakeUp(); err == nil {
			if err = r.serialPort.Close(); err == nil {
				connected, err = r.ecuReader.Connect()
			}
		}
	}

	return connected, err
}

func (r *MEMS19Reader) SendAndReceive(command []byte) ([]byte, error) {
	return r.ecuReader.SendAndReceive(command)
}

func (r *MEMS19Reader) Disconnect() error {
	return r.ecuReader.Disconnect()
}

func (r *MEMS19Reader) wakeUp() error {
	var err error

	// clear the line
	if err = r.serialPort.SetBreak(false); err == nil {
		time.Sleep(time.Millisecond * 2000)

		start := time.Now()

		// start bit
		_ = r.serialPort.SetBreak(true)
		r.sleepUntil(start, 200)

		// send the byte
		ecuAddress := 0x16
		for i := 0; i < 8; i++ {
			bit := (ecuAddress >> i) & 1

			if bit > 0 {
				_ = r.serialPort.SetBreak(false)
			} else {
				_ = r.serialPort.SetBreak(true)
			}

			r.sleepUntil(start, 200+((i+1)*200))

		}

		// stop bit
		_ = r.serialPort.SetBreak(false)
		r.sleepUntil(start, 2000)
	}

	return err
}

func (r *MEMS19Reader) wakeUpConnect() error {
	var err error

	if r.serialPort, err = sers.Open(r.ecuReader.port); err == nil {
		defer func(serialPort sers.SerialPort) {
			err := serialPort.Close()
			if err != nil {
				log.Errorf("unable to close serial port for mems 1.9 ecu (%s)", err)
			}
		}(r.serialPort)

		if err = r.serialPort.SetMode(9600, 8, sers.N, 1, sers.NO_HANDSHAKE); err == nil {
			// minread = 0: minimal buffering on read, return characters as early as possible
			// timeout = 1.0: time out if after 1.0 seconds nothing is received
			err = r.serialPort.SetReadParams(0, 0.001)
		}
	}

	return err
}

func (r *MEMS19Reader) sleepUntil(start time.Time, plus int) {
	target := start.Add(time.Duration(plus) * time.Millisecond)
	sleepMs := target.Sub(time.Now()).Milliseconds()
	if sleepMs < 0 {
		return
	}
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)
}
