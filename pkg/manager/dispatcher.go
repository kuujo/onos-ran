// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"fmt"
	"sync"
)

// Event is a stream event
type Event struct {
	// Type is the stream event type
	Type string

	// Object is the event object
	Object interface{}
}

// Dispatcher :
type Dispatcher struct {
	telemetryListenersLock sync.RWMutex
	telemetryListeners     map[string]chan Event
}

// NewDispatcher creates and initializes a new event dispatcher
func newDispatcher() *Dispatcher {
	return &Dispatcher{
		telemetryListeners: make(map[string]chan Event),
	}
}

// ListenTelemetryEvents :
func (d *Dispatcher) listenTelemetryEvents(telemetryEventChannel <-chan Event) {
	log.Info("Telemetry Event listener initialized")

	for tEvent := range telemetryEventChannel {
		d.telemetryListenersLock.RLock()
		for _, nbiChan := range d.telemetryListeners {
			nbiChan <- tEvent
		}
		d.telemetryListenersLock.RUnlock()
	}
}

// RegisterUeListener :
func (d *Dispatcher) registerTelemetryListener(subscriber string) (chan Event, error) {
	d.telemetryListenersLock.Lock()
	defer d.telemetryListenersLock.Unlock()
	if _, ok := d.telemetryListeners[subscriber]; ok {
		return nil, fmt.Errorf("telemetry listener %s is already registered", subscriber)
	}
	channel := make(chan Event)
	d.telemetryListeners[subscriber] = channel
	log.Infof("Registered %s in Telemetry listeners", subscriber)
	return channel, nil
}

func (d *Dispatcher) unregisterTelemetryListener(subscriber string) {
	d.telemetryListenersLock.Lock()
	defer d.telemetryListenersLock.Unlock()
	channel, ok := d.telemetryListeners[subscriber]
	if !ok {
		log.Infof("Subscriber %s had not been registered", subscriber)
		return
	}
	delete(d.telemetryListeners, subscriber)
	close(channel)
	log.Infof("Unregistered %s from Telemetry listeners", subscriber)
}

// GetListeners returns a list of registered listeners names
func (d *Dispatcher) GetListeners() []string {
	listenerKeys := make([]string, 0)
	d.telemetryListenersLock.RLock()
	defer d.telemetryListenersLock.RUnlock()
	for k := range d.telemetryListeners {
		listenerKeys = append(listenerKeys, k)
	}
	return listenerKeys
}
