// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package agent

import (
	"context"
	"os"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/joshuar/go-hass-agent/internal/config"
	"github.com/joshuar/go-hass-agent/internal/device"
	"github.com/joshuar/go-hass-agent/internal/hass"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/message"
)

const (
	Name      = "go-hass-agent"
	Version   = "0.0.1"
	fyneAppID = "com.github.joshuar.go-hass-agent"
)

type Agent struct {
	App           fyne.App
	Tray          fyne.Window
	Name, Version string
	MsgPrinter    *message.Printer
}

func NewAgent() *Agent {
	return &Agent{
		App:        newUI(),
		Name:       Name,
		Version:    Version,
		MsgPrinter: newMsgPrinter(),
	}
}

func Run() {
	ctx, cancelfunc := context.WithCancel(context.Background())
	ctx = device.SetupContext(ctx)
	log.Info().Msg("Starting agent.")
	agent := NewAgent()

	// If possible, create and log to a file as well as the console.
	logFile, err := storage.Child(agent.App.Storage().RootURI(), "go-hass-app.log")
	if err != nil {
		log.Error().Err(err).
			Msg("Unable to create a log file. Will only write logs to stdout.")
	} else {
		logWriter, err := storage.Writer(logFile)
		if err != nil {
			log.Error().Err(err).
				Msg("Unable to open log file for writing.")
		} else {
			consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
			multiWriter := zerolog.MultiLevelWriter(consoleWriter, logWriter)
			log.Logger = log.Output(multiWriter)
		}
	}

	var wg sync.WaitGroup

	// Try to load the app config. If it is not valid, start a new registration
	// process. Keep trying until we successfully register with HA or the user
	// quits.
	wg.Add(1)
	go func() {
		defer wg.Done()
		appConfig := agent.loadConfig()
		for appConfig.Validate() != nil {
			log.Warn().Msg("No suitable existing config found! Starting new registration process")
			err := agent.runRegistrationWorker(ctx)
			if err != nil {
				log.Error().Err(err).
					Msgf("Error trying to register: %v. Exiting.")
				agent.stop()
			}
			appConfig = agent.loadConfig()
		}
	}()

	agent.setupSystemTray()

	go agent.tracker(ctx, &wg)

	agent.App.Run()
	agent.stop()
	cancelfunc()
}

func (agent *Agent) stop() {
	log.Info().Msg("Shutting down agent.")
	agent.Tray.Close()
}

// tracker should be run in a goroutine and is responsible for creating,
// tracking and updating HA with all sensors provided from the platform/device.
func (agent *Agent) tracker(agentCtx context.Context, configWG *sync.WaitGroup) {
	configWG.Wait()

	appConfig := agent.loadConfig()
	ctx := config.NewContext(agentCtx, appConfig)
	sensorInfo := device.SetupSensors()
	sensors := make(map[string]*sensorState)
	updateCh := make(chan interface{})

	go agent.runNotificationsWorker(ctx)

	// goroutine to listen for sensor updates. Sensors are tracked in a map to
	// handle registration and disabling/enabling. Updates are sent to Home
	// Assistant.
	go func() {
		for {
			select {
			case data := <-updateCh:
				switch data := data.(type) {
				case hass.SensorUpdate:
					sensorID := data.ID()
					if _, ok := sensors[sensorID]; !ok {
						sensors[sensorID] = newSensor(data)
						log.Debug().Caller().Msgf("New sensor discovered: %s", sensors[sensorID].name)
					} else {
						sensors[sensorID].updateSensor(ctx, data)
					}
					go hass.APIRequest(ctx, sensors[sensorID])
				case hass.LocationUpdate:
					l := &location{
						data: data,
					}
					go hass.APIRequest(ctx, l)
				}
			case <-ctx.Done():
				log.Debug().Caller().
					Msg("Stopping sensor tracking.")
				return
			}
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		device.LocationUpdater(ctx, agent.App.UniqueID(), updateCh)
	}()

	// Run all the defined sensor update functions.
	for name, workerFunction := range sensorInfo.Get() {
		wg.Add(1)
		log.Debug().Caller().
			Msgf("Setting up sensors for %s.", name)
		go func(worker func(context.Context, chan interface{})) {
			defer wg.Done()
			worker(ctx, updateCh)
		}(workerFunction)
	}
	wg.Wait()
}
