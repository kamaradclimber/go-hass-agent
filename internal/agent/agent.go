package agent

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jeandeaual/go-locale"
	"github.com/joshuar/go-hass-agent/internal/hass"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	Name    = "go-hass-agent"
	Version = "0.0.1"
)

type Agent struct {
	App    fyne.App
	Tray   fyne.Window
	config AppConfig
	// hassConfig    *hass.ConfigResponse
	Name, Version string
	configLoaded  chan bool
	requestsCh    chan interface{}
	responsesCh   chan interface{}
	MsgPrinter    *message.Printer
}

func NewAgent() *Agent {
	a := &Agent{
		App:          NewUI(),
		Name:         Name,
		Version:      Version,
		configLoaded: make(chan bool, 1),
		requestsCh:   make(chan interface{}),
		responsesCh:  make(chan interface{}),
	}

	userLocales, err := locale.GetLocales()
	if err != nil {
		log.Warn().Msg("Could not find a suitable locale. Defaulting to English.")
		a.MsgPrinter = message.NewPrinter(message.MatchLanguage(language.English.String()))
	} else {
		a.MsgPrinter = message.NewPrinter(message.MatchLanguage(userLocales...))
		log.Debug().Caller().Msgf("Setting language to %v.", userLocales)
	}

	return a
}

// func (a *Agent) GetHassConfig(configLoaded chan bool) {
// 	<-configLoaded
// 	a.hassConfig = hass.GetConfig(a.config.RestAPIURL)
// }

func (a *Agent) SetupSystemTray() {
	// a.hassConfig = hass.GetConfig(a.config.RestAPIURL)
	a.Tray = a.App.NewWindow("System Tray")
	a.Tray.SetMaster()
	if desk, ok := a.App.(desktop.App); ok {
		log.Debug().Caller().
			Msg("Config loaded successfully. Creating tray icon.")

		menuItemAbout := fyne.NewMenuItem("About", func() {
			w := a.App.NewWindow(a.MsgPrinter.Sprintf("About ", a.Name))
			w.SetContent(container.New(layout.NewVBoxLayout(),
				widget.NewLabel(a.MsgPrinter.Sprintf("App Version: %s", a.Version)),
				// widget.NewLabel("Home Assistant Version: "+a.hassConfig.Version),
				widget.NewButton("Ok", func() {
					w.Close()
				}),
			))
			w.Show()
		})
		menu := fyne.NewMenu(a.Name, menuItemAbout)
		desk.SetSystemTrayMenu(menu)
	}
	a.Tray.Hide()
}

func (a *Agent) RunWorkers() {
	<-a.configLoaded
	go hass.RequestDispatcher(a.config.RestAPIURL, a.requestsCh, a.responsesCh)
	go a.runLocationWorker()
	go a.runSensorWorker()
}
