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

package e2

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/indications"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"github.com/onosproject/onos-ric/pkg/store/requests"
	"sync"
)

var log = logging.GetLogger("southbound", "e2")

func init() {
	log.SetLevel(logging.DebugLevel)
}

// NewSessionManager creates a new session manager
func NewSessionManager(devices device.Store, masterships mastership.Store, requests requests.Store, indications indications.Store) (*SessionManager, error) {
	mgr := &SessionManager{
		sessions:    make(map[device.ID]*Session),
		devices:     devices,
		masterships: masterships,
		requests:    requests,
		indications: indications,
	}
	err := mgr.start()
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

// SessionManager is an E2 session manager
type SessionManager struct {
	sessions    map[device.ID]*Session
	devices     device.Store
	masterships mastership.Store
	requests    requests.Store
	indications indications.Store
	mu          sync.RWMutex
}

// start starts the manager
func (m *SessionManager) start() error {
	deviceCh := make(chan device.Event)
	err := m.devices.Watch(deviceCh)
	if err != nil {
		return err
	}
	go m.processDeviceEvents(deviceCh)
	return nil
}

// processDeviceEvents processes events from the device store
func (m *SessionManager) processDeviceEvents(ch <-chan device.Event) {
	for event := range ch {
		log.Infof("Received device update %v", event.Device)
		err := m.processDeviceEvent(event)
		if err != nil {
			log.Errorf("Error updating session %s", event.Device.ID, err)
		}
	}
}

// processDeviceEvent processes an event from the device store
func (m *SessionManager) processDeviceEvent(event device.Event) error {
	if event.Type == device.EventNone || event.Type == device.EventAdded {
		return m.addDevice(event.Device)
	} else if event.Type == device.EventRemoved {
		return m.removeDevice(event.Device)
	}
	return nil
}

// addDevice adds a device session
func (m *SessionManager) addDevice(device device.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	election, err := m.masterships.GetElection(mastership.NewKey(device.ID.PlmnId, device.ID.Ecid))
	if err != nil {
		return err
	}

	session, err := newSession(device, election)
	if err != nil {
		return err
	}

	requestsCh := make(chan requests.Event)
	err = m.requests.Watch(device.ID, requestsCh, requests.WithReplay())
	if err != nil {
		return err
	}

	responsesCh := make(chan e2ap.RicIndication)
	err = session.subscribe(responsesCh)
	if err != nil {
		return err
	}

	go m.processRequests(session, requestsCh)
	go m.processResponses(responsesCh)

	oldSession, ok := m.sessions[device.ID]
	if ok {
		oldSession.Close()
	}
	m.sessions[device.ID] = session
	return nil
}

// removeDevice removes a device session
func (m *SessionManager) removeDevice(device device.Device) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, ok := m.sessions[device.ID]
	if ok {
		session.Close()
		delete(m.sessions, device.ID)
	}
	return nil
}

// processRequests processes events from the requests store
func (m *SessionManager) processRequests(session *Session, ch <-chan requests.Event) {
	for event := range ch {
		if event.Type == requests.EventNone || event.Type == requests.EventAppend {
			log.Debugf("Processing request %+v", event.Request.RicControlRequest)
			err := m.processRequest(session, event.Request)
			if err != nil {
				log.Errorf("Error processing request %+v: %v", event.Request.RicControlRequest, err)
			}
		}
	}
}

// processRequest processes an event from the requests store
func (m *SessionManager) processRequest(session *Session, request requests.Request) error {
	err := session.sendRequest(request.RicControlRequest)
	if err != nil {
		return err
	}
	return m.requests.Ack(&request)
}

// processResponses processes responses from the session
func (m *SessionManager) processResponses(ch <-chan e2ap.RicIndication) {
	for response := range ch {
		log.Debugf("Processing response %v", response)
		err := m.processResponse(response)
		if err != nil {
			log.Errorf("Error processing response %v: %v", response, err)
		}
	}
}

// processResponse processes a response from the session
func (m *SessionManager) processResponse(response e2ap.RicIndication) error {
	return m.indications.Record(indications.New(&response))
}

// Close closes the session manager
func (m *SessionManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, session := range m.sessions {
		session.Close()
	}
	return nil
}
