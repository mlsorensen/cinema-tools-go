package serial

import (
	"github.com/tarm/serial"
)

type Aspect string

const (
	DefaultBaud = 9600

	CommandOpen  = "0101U\n"
	CommandClose = "0101D\n"
	CommandMy    = "0101S\n"
	CommandStop  = CommandMy

	Aspect_24   Aspect = CommandOpen                 // 2.4:1
	Aspect_4x3  Aspect = CommandClose                // 4:3
	Aspect_16x9 Aspect = CommandMy                   // 16:9
	Aspect_185  Aspect = "0101S;W9;0101U;W1;0101S\n" // 1.85:1
	Aspect_20   Aspect = "0101S;W9;0101U;W2;0101S\n" // 2.0:1
	Aspect_22   Aspect = "0101U;W9;0101D;W2;0101S\n" // 2.2:1
	Aspect_143  Aspect = "0101D;W9;0101U;W1;0101S\n" // 1.43:1
)

type RTSSession struct {
	SerialPort string
	openPort   *serial.Port
}

// TODO: handle closing session

// NewSession establishes a new serial connection to a URTSI II
// device.
func (r *RTSSession) NewSession() error {
	c := &serial.Config{Name: r.SerialPort, Baud: DefaultBaud}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	r.openPort = s
	return nil
}

// Send sends an arbitrary command string over serial
func (r *RTSSession) Send(command string) error {
	if r.openPort == nil {
		err := r.NewSession()
		if err != nil {
			return err
		}
	}

	_, err := r.openPort.Write([]byte(command))
	if err != nil {
		return err
	}

	return nil
}

// SetAspect sends a defined Aspect command string over serial
func (r *RTSSession) SetAspect(aspect Aspect) error {
	return r.Send(string(aspect))
}
