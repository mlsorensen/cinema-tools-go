package main

import (
	"cinema-tools/internal/cfg"
	"errors"
	"flag"
	"fmt"
	grafikeye "github.com/mlsorensen/grafikeye/pkg/serial"
	lumagen "github.com/mlsorensen/lumagen/pkg/serial"
	lumagenMessage "github.com/mlsorensen/lumagen/pkg/serial/message"
	lumagenParsers "github.com/mlsorensen/lumagen/pkg/serial/parsers"
	urtsi2 "github.com/mlsorensen/urtsi2/pkg/serial"
	"log"
	"strconv"
)

var (
	config           *cfg.Config
	lumagenSession   *lumagen.LumagenSession
	grafikEyeSession *grafikeye.QSESession
	urtsiSession     *urtsi2.RTSSession
)

func main() {
	cfgPath := flag.String("config", "./config.yaml", "Path of config file")
	flag.Parse()

	err := cfg.ValidateConfigPath(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	config, err = cfg.LoadConfig(*cfgPath)

	if !config.LumagenConfigured() {
		log.Fatal("Please provide a serial port via the -lumagen-port flag")
	}

	if !config.UrtsiConfigured() {
		log.Fatal("Please provide a serial port via the -urtsi-port flag")
	}

	if config.GrafikEyeConfigured() {
		log.Printf("Connecting grafikEye at %s\n", config.GrafikEyePort)
		grafikEyeSession = &grafikeye.QSESession{SerialPort: config.GrafikEyePort, BaudRate: 115200}
		err = grafikEyeSession.StartMonitor(grafikEyeMessageHandler)
		if err != nil {
			log.Fatal(err)
		}
	}

	urtsiSession = &urtsi2.RTSSession{SerialPort: config.UrtsiPort}
	log.Printf("Connecting to URTSI at %s\n", config.UrtsiPort)

	lumagenSession = &lumagen.LumagenSession{SerialPort: config.LumagenPort}
	parser := lumagenParsers.ZQI22Parser{Handler: lumagenMessageHandler}
	err = lumagenSession.StartMessageMonitor([]lumagenParsers.Parser{parser})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Monitoring Lumagen at %s\n", config.LumagenPort)

	select {}
}

func lumagenMessageHandler(msg lumagenMessage.ZQI22Message) {
	log.Printf("Handling I22 message: %v\n", msg)

	var matched = false
	for _, action := range config.LumagenActions {
		if action.FramerateIn(msg.SourceFrameRate) && action.AspectIn(msg.SourceAspectRatio) {
			matched = true
			err := sendUrtsiCommand(action.UrtsiCommand)
			if err != nil {
				log.Printf("error while sending URTSI command: %v\n", err)
			}

			err = pressGrafikEyeButtonId(action.GrafikEyePressButtonId)
			if err != nil {
				log.Printf("error while sending GrafikEye command: %v\n", err)
			}
		}
	}
	if !matched {
		log.Printf("No action found matching aspect %d and framerate %d", msg.SourceAspectRatio, msg.SourceFrameRate)
	}
}

func grafikEyeMessageHandler(msg grafikeye.QSCommand) {
	if msg.IntegrationId != grafikeye.GrafikEye || msg.Type != grafikeye.TypeDevice {
		return
	}

	log.Printf("Handling GrafikEye message %v\n", msg)

	if len(msg.CommandFields) < 2 {
		log.Println("incoming Grafik Eye event didn't contain command field length of at least 2")
		return
	}

	var matched = false
	for _, action := range config.GrafikEyeActions {
		buttonId := strconv.Itoa(int(action.ButtonId))
		eventId := strconv.Itoa(int(action.Event))
		if msg.CommandFields[0] == buttonId && msg.CommandFields[1] == eventId {
			matched = true
			err := sendUrtsiCommand(action.UrtsiCommand)
			if err != nil {
				log.Printf("error while sending URTSI comand: %v\n", err)
			}
		}
	}

	if !matched {
		log.Printf("No GrafikEye action found matching button id %s event %s", msg.CommandFields[0], msg.CommandFields[1])
	}
}

func sendUrtsiCommand(command string) error {
	if len(command) == 0 {
		return errors.New("no URTSI command provided, skipping")
	}

	if len(config.UrtsiPort) == 0 {
		return errors.New("no URTSI port is configured, not sending command")
	}

	if len(command) == 0 {
		return errors.New("command was empty, not sending to URTSI")
	}

	if urtsiSession == nil {
		return errors.New("URTSI session not yet established, skipping")
	}

	return urtsiSession.Send(fmt.Sprintf("%s\n", command))
}

func pressGrafikEyeButtonId(buttonId uint8) error {
	if grafikEyeSession == nil {
		return nil
	}

	if buttonId == 0 {
		return errors.New("no GrafikEye command provided, skipping")
	}

	if len(config.GrafikEyePort) == 0 {
		return errors.New("no GrafikEye port is configured, not sending command")
	}

	return grafikEyeSession.PressButton(strconv.Itoa(int(buttonId)))
}
