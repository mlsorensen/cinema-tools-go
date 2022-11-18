package main

import (
	"flag"
	"fmt"
	"github.com/mlsorensen/lumagen/pkg/serial"
	"github.com/mlsorensen/lumagen/pkg/serial/message"
	"github.com/mlsorensen/lumagen/pkg/serial/parsers"
	"log"
)

func main() {
	port := flag.String("port", "", "Serial port to listen on")
	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("Please provide a serial port via the -port flag")
	}

	mon := serial.LumagenSession{SerialPort: *port}
	parser := parsers.ZQI22Parser{Handler: handleZQI22Message}
	err := mon.StartMessageMonitor([]parsers.Parser{parser})
	if err != nil {
		panic(err)
	}

	select {}
}

func handleZQI22Message(msg message.ZQI22Message) {
	fmt.Printf("got I22 message: %v\n", msg)
}
