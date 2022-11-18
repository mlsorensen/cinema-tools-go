package main

import (
	"flag"
	. "github.com/mlsorensen/urtsi2/pkg/serial"
	"log"
)

func main() {
	port := flag.String("port", "", "Serial port to listen on")
	closeShade := flag.Bool("close", false, "Close shade")
	openShade := flag.Bool("open", false, "Open shade")
	flag.Parse()

	if len(*port) == 0 {
		log.Fatal("Please provide a serial port via the -port flag")
	}

	if *closeShade && *openShade {
		log.Fatal("Provide only -close or -open flag, not both")
	}

	if !*closeShade && !*openShade {
		log.Fatal("One of -close or -open is required")
	}

	cmd := CommandOpen
	if *closeShade {
		cmd = CommandClose
	}

	session := RTSSession{SerialPort: *port}
	err := session.Send(cmd)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Command sent... did it work?")
}
