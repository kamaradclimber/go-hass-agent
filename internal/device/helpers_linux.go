// Copyright (c) 2023 Joshua Rich <joshua.rich@gmail.com>
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package device

import (
	"context"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/rs/zerolog/log"
)

const (
	sessionBus dbusType = iota
	systemBus
)

type dbusType int
type bus struct {
	conn    *dbus.Conn
	events  chan *dbus.Signal
	busType dbusType
}

// NewBus sets up DBus connections and channels for receiving signals. It creates both a system and session bus connection.
func NewBus(ctx context.Context, t dbusType) *bus {
	var conn *dbus.Conn
	var err error
	switch t {
	case sessionBus:
		conn, err = dbus.ConnectSessionBus(dbus.WithContext(ctx))
	case systemBus:
		conn, err = dbus.ConnectSystemBus(dbus.WithContext(ctx))
	}
	if err != nil {
		log.Error().Stack().Err(err).
			Msg("Could not connect to bus")
		return nil
	}
	events := make(chan *dbus.Signal)
	conn.Signal(events)
	return &bus{
		conn:    conn,
		events:  events,
		busType: t,
	}
}

// DBusWatchRequest contains all the information required to set-up a DBus match
// signal watcher.
type DBusWatchRequest struct {
	bus          dbusType
	path         dbus.ObjectPath
	match        []dbus.MatchOption
	event        string
	eventHandler func(*dbus.Signal)
}

func (d *deviceAPI) bus(t dbusType) *dbus.Conn {
	switch t {
	case sessionBus:
		return d.dBusSession.conn
	case systemBus:
		return d.dBusSystem.conn
	default:
		log.Warn().Msg("Could not discern DBus bus type.")
		return nil
	}
}

// monitorDBus listens for DBus watch requests and ensures the appropriate
// signal watches are created. It will also dispatch to a handler function when
// a signal is matched.
func (d *deviceAPI) monitorDBus(ctx context.Context) {
	events := make(map[dbusType]map[string]func(*dbus.Signal))
	events[sessionBus] = make(map[string]func(*dbus.Signal))
	events[systemBus] = make(map[string]func(*dbus.Signal))
	watches := make(map[dbusType]*DBusWatchRequest)
	defer close(d.WatchEvents)
	defer d.dBusSession.conn.RemoveSignal(d.dBusSession.events)
	defer d.dBusSystem.conn.RemoveSignal(d.dBusSystem.events)
	for {
		select {
		// When the context is finished/cancelled, try to clean up gracefully.
		case <-ctx.Done():
			log.Debug().Caller().Msg("Stopping DBus Monitor.")
			close(d.WatchEvents)
			d.dBusSession.conn.RemoveSignal(d.dBusSession.events)
			d.dBusSystem.conn.RemoveSignal(d.dBusSystem.events)
			return
		// When a new watch request is received, send it to the right DBus
		// connection and record it so it can be matched to a handler.
		case watch := <-d.WatchEvents:
			d.AddDBusWatch(watch.bus, watch.match)
			events[watch.bus][string(watch.path)] = watch.eventHandler
			watches[watch.bus] = watch
		// For each bus signal handler, try to match first on an exact path
		// match, then try a substr match. Whichever matches, run the
		// handler function associated with it.
		case systemSignal := <-d.dBusSystem.events:
			if handlerFunc, ok := events[systemBus][string(systemSignal.Path)]; ok {
				handlerFunc(systemSignal)
			} else {
				for matchPath, handlerFunc := range events[systemBus] {
					if strings.Contains(string(systemSignal.Path), matchPath) {
						handlerFunc(systemSignal)
					}
				}
			}
		case sessionSignal := <-d.dBusSession.events:
			if handlerFunc, ok := events[sessionBus][string(sessionSignal.Path)]; ok {
				handlerFunc(sessionSignal)
			} else {
				for matchPath, handlerFunc := range events[systemBus] {
					if strings.Contains(string(sessionSignal.Path), matchPath) {
						handlerFunc(sessionSignal)
					}
				}
			}
		}
	}
}

// AddDBusWatch will add a matcher to the specified bus monitoring for
// the specified path and interface.
func (d *deviceAPI) AddDBusWatch(t dbusType, matches []dbus.MatchOption) error {
	if err := d.bus(t).AddMatchSignal(matches...); err != nil {
		return err
	} else {
		return nil
	}
}

// RemoveDBusWatch will remove a matcher from the specified bus to stop it
// generating signals for that match.
func (d *deviceAPI) RemoveDBusWatch(t dbusType, path dbus.ObjectPath, intr string) error {
	if err := d.bus(t).RemoveMatchSignal(
		dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(intr),
	); err != nil {
		return err
	} else {
		return nil
	}
}

func (d *deviceAPI) GetDBusObject(t dbusType, dest string, path dbus.ObjectPath) dbus.BusObject {
	return d.bus(t).Object(dest, path)
}

// GetDBusProp will retrieve the specified property value from the given path
// and destination.
func (d *deviceAPI) GetDBusProp(t dbusType, dest string, path dbus.ObjectPath, prop string) dbus.Variant {
	obj := d.bus(t).Object(dest, path)
	res, err := obj.GetProperty(prop)
	if err != nil {
		log.Error().Err(err).
			Msgf("Unable to retrieve property %s (%s)", prop, dest)
		return dbus.MakeVariant("")
	}
	return res
}

func (d *deviceAPI) GetDBusDataAsMap(t dbusType, dest string, path dbus.ObjectPath, method string, args string) map[string]dbus.Variant {
	obj := d.bus(t).Object(dest, path)
	var data map[string]dbus.Variant
	var err error
	if args != "" {
		err = obj.Call(method, 0, args).Store(&data)
	} else {
		err = obj.Call(method, 0).Store(&data)
	}
	if err != nil {
		log.Error().Err(err).
			Msgf("Unable to execute %s on %s (args: %s)", method, dest, args)
		return nil
	}
	return data
}

func (d *deviceAPI) GetDBusDataAsList(t dbusType, dest string, path dbus.ObjectPath, method string, args string) []string {
	obj := d.bus(t).Object(dest, path)
	var data []string
	var err error
	if args != "" {
		err = obj.Call(method, 0, args).Store(&data)
	} else {
		err = obj.Call(method, 0).Store(&data)
	}
	if err != nil {
		log.Error().Err(err).
			Msgf("Unable to execute %s on %s (args: %s)", method, dest, args)
		return nil
	}
	return data
}

func (d *deviceAPI) GetDBusData(t dbusType, dest string, path dbus.ObjectPath, method string, args string) interface{} {
	obj := d.bus(t).Object(dest, path)
	var data interface{}
	var err error
	if args != "" {
		err = obj.Call(method, 0, args).Store(&data)
	} else {
		err = obj.Call(method, 0).Store(&data)
	}
	if err != nil {
		log.Error().Err(err).
			Msgf("Unable to execute %s on %s (args: %s)", method, dest, args)
		return nil
	}
	return data
}

// FindPortal is a helper function to work out which portal interface should be
// used for getting information on running apps.
func FindPortal() string {
	switch os.Getenv("XDG_CURRENT_DESKTOP") {
	case "KDE":
		return "org.freedesktop.impl.portal.desktop.kde"
	case "GNOME":
		return "org.freedesktop.impl.portal.desktop.kde"
	default:
		log.Warn().Msg("Unsupported desktop/window environment.")
		return ""
	}
}
