package main

import (
	"fmt"

	"github.com/Grayda/go-phone"
)

func main() {
	// Connect to our COM port
	phone.Start("/dev/ttyACM0")

	for { // Loop forever
		select { // This lets us do non-blocking channel reads. If we have a message, process it. If not, loop
		case msg := <-phone.Events: // If there is an event waiting
			switch msg.Name { // What event is it?
			case "READY": // Connected to COM port, ready to start sending / receiving data!
				fmt.Println("Connected to serial. Waiting for data!")
				phone.Read()
			case "RING":
				fmt.Println("Phone is ringing!")
				phone.Read()
			case "OTHER":
				phone.Read()
			case "NMBR":
				if msg.Message == "P" {
					fmt.Println("Private number detected!")
				} else {
					fmt.Println("Number detected:", msg.Message)
				}
				fmt.Println("Number detected:", msg.Message)
				phone.Read()
			default: // If there are no messages to parse, look for more bytes from our port, then try again
				phone.Read()
			}
		}
	}
}
