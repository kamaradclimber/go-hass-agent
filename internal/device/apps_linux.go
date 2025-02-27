// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package device

import (
	"context"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/iancoleman/strcase"
	"github.com/joshuar/go-hass-agent/internal/hass"
	"github.com/rs/zerolog/log"
)

//go:generate stringer -type=appSensorType -output appSensor_types_linux.go
const (
	appStateDBusMethod    = "org.freedesktop.impl.portal.Background.GetAppState"
	appStateDBusPath      = "/org/freedesktop/portal/desktop"
	appStateDBusInterface = "org.freedesktop.impl.portal.Background"
	appStateDBusEvent     = "org.freedesktop.impl.portal.Background.RunningApplicationsChanged"

	ActiveApp appSensorType = iota
	RunningApps
)

type appSensorType int

type appSensor struct {
	sensorType  appSensorType
	sensorValue map[string]dbus.Variant
}

// appSensor implements hass.SensorUpdate

func (s *appSensor) Name() string {
	return strcase.ToDelimited(s.sensorType.String(), ' ')
}

func (s *appSensor) ID() string {
	return strings.ToLower(strcase.ToSnake(s.sensorType.String()))
}

func (s *appSensor) Icon() string {
	return "mdi:application"
}

func (s *appSensor) SensorType() hass.SensorType {
	return hass.TypeSensor
}

func (s *appSensor) DeviceClass() hass.SensorDeviceClass {
	return 0
}

func (s *appSensor) StateClass() hass.SensorStateClass {
	switch s.sensorType {
	case RunningApps:
		return hass.Measurement
	default:
		return 0
	}
}

func (s *appSensor) State() interface{} {
	switch s.sensorType {
	case ActiveApp:
		for appName, state := range s.sensorValue {
			if state.Value().(uint32) == 2 {
				return appName
			}
		}
	case RunningApps:
		var count int
		for _, state := range s.sensorValue {
			if state.Value().(uint32) > 0 {
				count++
			}
		}
		return count
	}
	return ""
}

func (s *appSensor) Units() string {
	return ""
}

func (s *appSensor) Category() string {
	return ""
}

func (s *appSensor) Attributes() interface{} {
	var runningApps []string
	for appName, state := range s.sensorValue {
		if state.Value().(uint32) > 0 {
			runningApps = append(runningApps, appName)
		}
	}
	return struct {
		RunningApps []string `json:"Running Apps"`
	}{
		RunningApps: runningApps,
	}
}

func marshallAppStateUpdate(t appSensorType, v map[string]dbus.Variant) *appSensor {
	return &appSensor{
		sensorValue: v,
		sensorType:  t,
	}
}

func AppUpdater(ctx context.Context, update chan interface{}) {
	deviceAPI, deviceAPIExists := FromContext(ctx)
	if !deviceAPIExists {
		log.Debug().Caller().
			Msg("Could not connect to DBus to monitor app state.")
		return
	}

	portalDest := FindPortal()
	if portalDest == "" {
		log.Debug().Caller().
			Msgf("Unsupported or unknown portal")
		return
	}

	appChangeSignal := &DBusWatchRequest{
		bus:  sessionBus,
		path: appStateDBusPath,
		match: []dbus.MatchOption{
			dbus.WithMatchObjectPath(appStateDBusPath),
			dbus.WithMatchInterface(appStateDBusInterface),
		},
		event: appStateDBusEvent,
		eventHandler: func(s *dbus.Signal) {
			activeAppList := deviceAPI.GetDBusDataAsMap(sessionBus,
				portalDest,
				appStateDBusPath,
				appStateDBusMethod, "")
			if activeAppList == nil {
				log.Debug().Caller().
					Msg("No active apps found.")
			} else {
				update <- marshallAppStateUpdate(RunningApps, activeAppList)
				update <- marshallAppStateUpdate(ActiveApp, activeAppList)
			}
		},
	}
	deviceAPI.WatchEvents <- appChangeSignal
}
