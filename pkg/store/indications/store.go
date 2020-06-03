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

package indications

import (
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/store/indications"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"io"
	"sync"
)

var log = logging.GetLogger("store", "indications")

func init() {
	log.SetLevel(logging.DebugLevel)
}

// New creates a new indication
func New(indication *e2ap.RicIndication) *Indication {
	return &Indication{
		RicIndication: indication,
	}
}

// Indication is an indication message
type Indication struct {
	*e2ap.RicIndication
}

// Event is a store event
type Event struct {
	// Type is the event type
	Type EventType

	// Indication is the indication
	Indication Indication
}

// EventType is a store event type
type EventType string

const (
	// EventNone indicates an event that was not triggered but replayed
	EventNone EventType = ""
	// EventReceived indicates the message was received in the store
	EventReceived EventType = "received"
)

// Store is interface for store
type Store interface {
	io.Closer

	// Record records an indication in the store
	Record(*Indication) error

	// List all of the last up to date messages
	List(ch chan<- Indication) error

	// Subscribe subscribes to indications
	Subscribe(ch chan<- Event, opts ...SubscribeOption) error
}

// NewDistributedStore creates a new distributed indications store
func NewDistributedStore(cluster cluster.Cluster, devices device.Store, masterships mastership.Store, config config.Config) (Store, error) {
	log.Info("Creating distributed indications store")
	store := &store{
		cluster:     cluster,
		devices:     devices,
		masterships: masterships,
		subscribers: make([]chan<- Event, 0),
		indications: make(map[device.Key]*deviceIndicationsStore),
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

var _ Store = &store{}

type store struct {
	cluster     cluster.Cluster
	devices     device.Store
	masterships mastership.Store
	subscribers []chan<- Event
	indications map[device.Key]*deviceIndicationsStore
	mu          sync.RWMutex
}

func (s *store) open() error {
	getServer().registerHandler(s)
	ch := make(chan device.Event)
	err := s.devices.Watch(ch)
	if err != nil {
		return err
	}
	go func() {
		for event := range ch {
			log.Debugf("Received device Event %v", event)
			if event.Type != device.EventUpdated {
				_, err := s.getDeviceStore(event.Device.ID.Key())
				if err != nil {
					log.Errorf("Failed to initialize indications store for device %s: %s", event.Device.ID.Key(), err)
				}
			}
		}
	}()
	return nil
}

func (s *store) getDeviceStore(deviceKey device.Key) (*deviceIndicationsStore, error) {
	s.mu.RLock()
	store, ok := s.indications[deviceKey]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		store, ok = s.indications[deviceKey]
		if !ok {
			election, err := s.masterships.GetElection(mastership.Key(deviceKey))
			if err != nil {
				return nil, err
			}
			store, err = newDeviceIndicationsStore(deviceKey, s.cluster, election)
			if err != nil {
				return nil, err
			}
			for _, subscriber := range s.subscribers {
				err := store.Subscribe(subscriber)
				if err != nil {
					return nil, err
				}
			}
			s.indications[deviceKey] = store
		}
	}
	return store, nil
}

func getDeviceID(indication *Indication) device.ID {
	switch indication.GetHdr().GetMessageType() {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_CELL:
		return device.ID(*indication.GetMsg().GetRadioMeasReportPerCell().Ecgi)
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		return device.ID(*indication.GetMsg().GetRadioMeasReportPerUE().Ecgi)
	case sb.MessageType_CELL_CONFIG_REPORT:
		return device.ID(*indication.GetMsg().GetCellConfigReport().Ecgi)
	case sb.MessageType_UE_ADMISSION_REQUEST:
		return device.ID(*indication.GetMsg().GetUEAdmissionRequest().Ecgi)
	case sb.MessageType_UE_RELEASE_IND:
		return device.ID(*indication.GetMsg().GetUEReleaseInd().Ecgi)
	}
	return device.ID{}
}

func (s *store) Record(indication *Indication) error {
	store, err := s.getDeviceStore(getDeviceID(indication).Key())
	if err != nil {
		return err
	}
	return store.Record(indication)
}

func (s *store) List(ch chan<- Indication) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	listWg := &sync.WaitGroup{}
	chanWg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, store := range s.indications {
		listWg.Add(1)
		chanWg.Add(1)
		deviceCh := make(chan Indication)
		go func(store Store) {
			err := store.List(deviceCh)
			if err != nil {
				errCh <- err
			}
			listWg.Done()
			for indication := range deviceCh {
				ch <- indication
			}
			chanWg.Done()
		}(store)
	}

	listWg.Wait()
	close(errCh)

	go func() {
		chanWg.Wait()
		close(ch)
	}()

	for err := range errCh {
		return err
	}
	return nil
}

func (s *store) Subscribe(ch chan<- Event, opts ...SubscribeOption) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, store := range s.indications {
		wg.Add(1)
		go func(store Store) {
			err := store.Subscribe(ch, opts...)
			if err != nil {
				errCh <- err
			}
			wg.Done()
		}(store)
	}

	wg.Wait()
	close(errCh)

	s.subscribers = append(s.subscribers, ch)

	for err := range errCh {
		return err
	}
	return nil
}

func (s *store) subscribe(request *indications.SubscribeRequest, stream indications.IndicationsService_SubscribeServer) error {
	store, err := s.getDeviceStore(device.Key(request.Device))
	if err != nil {
		return err
	}
	return store.subscribe(request, stream)
}

func (s *store) Close() error {
	return nil
}
