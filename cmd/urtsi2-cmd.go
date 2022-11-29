package main

import (
	"flag"
	. "github.com/mlsorensen/urtsi2/pkg/serial"
	"log"
	"os"
	"time"
)

func main() {
	port := flag.String("port", "", "Serial port to listen on")
	closeShade := flag.Bool("close", false, "Close shade")
	openShade := flag.Bool("open", false, "Open shade")
	flag.Parse()

	if len(*port) == 0 {
		log.Println("Please provide a serial port via the -port flag")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}

	if *closeShade && *openShade {
		log.Println("Provide only -close or -open flag, not both")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}

	if !*closeShade && !*openShade {
		log.Println("Please provide an action, one of flags -close or -open is required")
		time.Sleep(time.Second * 5)
		os.Exit(1)
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
