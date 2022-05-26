package rosco

import (
	"fmt"
	"github.com/distributed/sers"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type MEMSReader struct {
	connected  bool
	port       string
	serialPort sers.SerialPort
}

func NewMEMSReader(connection string) *MEMSReader {
	log.Infof("created mems ecu reader")

	r := &MEMSReader{}
	r.port = r.fixPort(connection)
	return r
}

func (r *MEMSReader) Connect() (bool, error) {
	r.connected = false

	if err := r.connectToSerialPort(r.port); err != nil {
		log.Errorf("error opening serial (%s) status : (%+v)", r.port, err)
		// connect failure if we cannot open the serialPort
		return false, err
	}

	if err := r.wakeUp(); err != nil {
		// wakeup failure if we cannot open the serialPort
		return false, err
	}

	r.flushSerialPort()

	if err := r.initialiseMemsECU(); err != nil {
		log.Errorf("error initialising ecu (%s) status : (%+v)", r.port, err)
		// connect failure if we cannot initialise successfully
		// disconnect from the ecu
		_ = r.Disconnect()
		return false, err
	}

	// connected, no errors
	r.connected = true

	return r.connected, nil
}

func (r *MEMSReader) SendAndReceive(command []byte) ([]byte, error) {
	var response []byte
	var err error

	if r.serialPort != nil {
		if r.connected {
			r.writeSerial(command)
			response, err = r.readSerial(command)
		} else {
			err = fmt.Errorf("ecu is not connected, unable to send %X", command)
			log.Errorf("%s", err)
		}
	} else {
		err = fmt.Errorf("serial connnection not initialised, unable to send %X", command)
		log.Errorf("%s", err)
	}

	return response, err
}

func (r *MEMSReader) Disconnect() error {
	var err error

	// don't try and close an uninitialised serial serialPort
	// the serial library throws and ugly fatal if that happens
	if r.serialPort != nil {
		if r.connected {
			if err = r.serialPort.Close(); err != nil {
				log.Warnf("error closing serial serialPort (%+v)", err)
			} else {
				log.Infof("serial port closed successully")
			}
		}
	}

	r.connected = false
	return err
}

func (r *MEMSReader) connectToSerialPort(port string) error {
	var err error

	log.Infof("attempting to open serial serialPort %s", port)

	// connect to the ecu
	if r.serialPort, err = sers.Open(port); err != nil {
		log.Errorf("error opening serial port (%s)", err)
	} else {
		if err = r.serialPort.SetMode(9600, 8, sers.N, 1, sers.NO_HANDSHAKE); err != nil {
			log.Errorf("error configuring serial port (%s)", err)
		} else {
			if err = r.serialPort.SetReadParams(0, 0.1); err != nil {
				log.Errorf("error setting serial port timeouts (%s)", err)
			}
		}
	}

	return err
}

// initialises the connection to the ECU
// The initialisation sequence is as follows:
//
// 1. Send command CA (MEMS_InitCommandA)
// 2. Receive response CA
// 3. Send command 75 (MEMS_InitCommandB)
// 4. Receive response 75
// 5. Send request ECU ID command D0 (MEMS_InitECUID)
// 6. Receive response D0 XX XX XX XX
//
func (r *MEMSReader) initialiseMemsECU() error {
	log.Infof("initialising ecu")

	r.writeSerial(MEMSInitCommandA)
	if response, err := r.readSerial(MEMSInitCommandA); err != nil {
		// abandon initialisation if error occurred
		log.Errorf("mems initialisation failed command %X (%s)", MEMSInitCommandA, err)
		return err
	} else {
		// if we get the command echoed back we can assume good connection and proceed.
		// This is to work around the issue  in Windows where the serialPort always connects even if it's not available.
		if len(response) == 0 {
			err = fmt.Errorf("0 bytes received, serial serialPort read error, timeout? (%s)", err)
			log.Errorf("%s", err)
			return err
		}

		if response[0] == MEMSInitCommandA[0] {
			r.writeSerial(MEMSInitCommandB)

			if response, err = r.readSerial(MEMSInitCommandB); err != nil {
				// abandon initialisation if error occurred
				log.Errorf("mems initialisation failed command %X (%s)", MEMSInitCommandB, err)
				return err
			}

			r.writeSerial(MEMSHeartbeat)
			if response, err = r.readSerial(MEMSHeartbeat); err != nil {
				// abandon initialisation if error occurred
				log.Errorf("mems initialisation failed command %X (%s)", MEMSHeartbeat, err)
				return err
			}

			r.writeSerial(MEMSInitECUID)
			if response, err = r.readSerial(MEMSInitECUID); err != nil {
				// abandon initialisation if error occurred
				log.Errorf("mems initialisation failed command %X (%s)", MEMSInitECUID, err)
				return err
			} else {
				log.Infof("mems initialisation ECU ID %X successful", response)
			}
		}
	}

	log.Infof("mems initialised")
	return nil
}

// readSerial read from MEMS
// read 1 byte at a time until we have all the expected bytes
func (r *MEMSReader) readSerial(command []byte) ([]byte, error) {
	var bytesRead int
	var err error
	var retry int

	size, err := getResponseSize(command)

	// serial read buffer
	b := make([]byte, size)

	//  receivedBytes frame buffer
	receivedBytes := make([]byte, 0)

	if r.serialPort != nil {
		// read all the expected bytes before returning the receivedBytes
		retry = 0

		for count := 0; count < size; {
			// wait for a response from MEMS
			bytesRead, err = r.serialPort.Read(b)

			if bytesRead == 0 {
				// wait for data and retry
				time.Sleep(50 * time.Millisecond)
				retry++

				if retry >= 5 {
					// drop out of loop, send back a 0x00 byte array response
					// this prevents the loop getting blocked on a read error
					err = fmt.Errorf("0 bytes received, serial port read error after %d retries (%s)", retry, err)
					log.Errorf("(%s)", err)
					return receivedBytes, err
				}
			} else {
				// append the read bytes to the receivedBytes frame
				receivedBytes = append(receivedBytes, b[:bytesRead]...)
			}

			// increment by the number of bytes read
			count = count + bytesRead
			if count > size {
				err = fmt.Errorf("received dataframe size mismatch, received %d, expected %d", count, size)
				log.Errorf("%s", err)
			}
		}
	}

	log.Infof("received %X from ecu, %d bytes", receivedBytes, bytesRead)

	if err == nil {
		err = r.commandMatchesResponse(command, receivedBytes)
	}

	return receivedBytes, err
}

func (r *MEMSReader) commandMatchesResponse(command []byte, receivedBytes []byte) error {
	var err error

	if command[0] != receivedBytes[0] {
		err = fmt.Errorf("expecting command echo of %X, received %X", command[0], receivedBytes[0])
		log.Errorf("%s", err)
	}

	return err
}

// writeSerial write to MEMS
func (r *MEMSReader) writeSerial(data []byte) {
	if r.serialPort != nil {
		bytesWritten, err := r.serialPort.Write(data)

		if err != nil {
			log.Errorf("error sending %X to the ecu (%s)", data, err)
		}

		if bytesWritten > 0 {
			log.Infof("sent %X to ecu", data)
		}
	}
}

func (r *MEMSReader) wakeUp() error {
	var err error

	log.Infof("performing serial port 5 baud wake-up")
	// clear the line
	if err = r.serialPort.SetBreak(false); err == nil {
		log.Debugf("setting up ecu wake-up")
		time.Sleep(time.Millisecond * 2000)
		start := time.Now()

		log.Debugf("sending ecu wake-up data")
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
		log.Debugf("clearing down ecu wake-up")
		_ = r.serialPort.SetBreak(false)
		r.sleepUntil(start, 2000)
	}

	log.Infof("completed serial port 5 baud wake-up")
	return err
}

func (r *MEMSReader) sleepUntil(start time.Time, plus int) {
	target := start.Add(time.Duration(plus) * time.Millisecond)
	sleepMs := target.Sub(time.Now()).Milliseconds()
	if sleepMs < 0 {
		return
	}
	time.Sleep(time.Duration(sleepMs) * time.Millisecond)
}

func (r *MEMSReader) flushSerialPort() {
	size := 200
	b := make([]byte, size)

	for {
		if bytesRead, err := r.serialPort.Read(b); err == nil {
			log.Infof("flushing serial port, read %d bytes", bytesRead)
		} else {
			log.Infof("serial port flushed")
			break
		}
	}
}

func (r *MEMSReader) fixPort(port string) string {
	if strings.Contains(port, "/dev/tty.") {
		port = strings.Replace(port, "/tty.", "/cu.", 1)
		log.Infof("fixed tty port to %s", port)
	}

	return port
}
