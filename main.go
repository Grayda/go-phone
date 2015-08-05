package phone

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// EventStruct is used to pass messages back to our calling code. Think of it like Node's "EventEmitter"
type EventStruct struct {
	Name    string
	Message string
}

// This lets our calling code view events. We only want 1 event at a time (else our code hangs)
var Events = make(chan EventStruct, 1) // Events is our events channel which will notify calling code that we have an event happening

var Ringing = false   // Ringing returns true if the phone is ringing. Times out after 2 seconds
var LastNumber string // LastNumber is the last phone number that called
var Connected = false

// New Port for working with
var serialport *serial.Port

// Stop closes our serial port
func Stop() {
	serialport.Close()
}

// Start connects to our COM port for reading and writing
func Start(COMPort string) error {
	var err error

	// Configure our serial port
	var com = &serial.Config{Name: COMPort, Baud: 9600}

	// c := &serial.Config{Name: COMPort, Baud: 9600}
	// Connect to our COM port
	serialport, err = serial.OpenPort(com)

	if err != nil {
		return err
	}

	Connected = true

	// Reset the modem. We have to add our own newlines
	_, err = serialport.Write([]byte("ATZ\r\n"))
	if err != nil {
		return err
	}

	// This is the most common method of turning on Caller ID. But sometimes it doesn't work. The modem will return "ERROR", but that's okay, because it's not a show-stopper
	_, err = serialport.Write([]byte("AT#CID=1\r\n"))
	if err != nil {
		return err
	}

	// The second most common method of turning on Caller ID. Again, will ERROR if it doesn't succeed, but we simply carry on
	_, err = serialport.Write([]byte("AT+VCID=1\r\n"))
	if err != nil {
		return err
	}

	passMessage("READY", "")
	return nil

}

// Read does what it says. Polls our serial port until something comes through. It's usually a blocking call, but we run it in a goroutine in our calling code
func Read() {
	// 128 bytes should be plenty
	buf := make([]byte, 128)
	n, err := serialport.Read(buf)

	if err != nil {
		log.Fatal(err)
	}

	// Find out what we're dealing with
	switch {
	// If our string contains "RING", it's the phone ringing. This is sent by the modem each time the phone rings (as in, an audible noise is made by the phone)
	case strings.Contains(string(buf[:n]), "RING") == true:
		time.AfterFunc(time.Second*2, func() {
			Ringing = false
		})
		passMessage("RING", "")

		// If our string contains "NMBR", it's Caller ID coming through
	case strings.Contains(string(buf[:n]), "NMBR") == true:
		// Simple regex that grabs everything after NMBR (which is the last item in our data anyway)
		r := regexp.MustCompile(`NMBR = (.+)?`)
		res := r.FindStringSubmatch(string(buf[:n]))

		LastNumber = res[1]
		passMessage("NMBR", res[1])
		// Something else?
	default:
		passMessage("OTHER", string(buf[:n]))
	}

}

// Passes a message back to our calling code. It does this by way of Channels
func passMessage(msgType string, msg string) bool {

	select {
	case Events <- EventStruct{msgType, msg}:

	default:
	}

	return true
}
