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

package telemetry

import (
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/api/store/message"
	"github.com/onosproject/onos-ric/pkg/config"
	messagestore "github.com/onosproject/onos-ric/pkg/store/message"
	timestore "github.com/onosproject/onos-ric/pkg/store/time"
	"io"

	"github.com/onosproject/onos-ric/api/sb"
)

var log = logging.GetLogger("store", "telemetry")

const primitiveName = "telemetry"

const keySep = ":"

// ID is a telemetry store identifier
type ID struct {
	MessageType sb.MessageType
	PlmnID      string
	Ecid        string
	Crnti       string
}

func (i ID) String() string {
	return fmt.Sprintf("%d%s%s%s%s%s%s", i.MessageType, keySep, i.PlmnID, keySep, i.Ecid, keySep, i.Crnti)
}

// Revision is a telemetry message revision
type Revision struct {
	// Term is the term in which the revision was created
	Term Term
	// Timestamp is the timestamp at which the revision was created
	Timestamp Timestamp
}

// Term is a telemetry store term
type Term uint64

// Timestamp is a telemetry store timestamp
type Timestamp uint64

// NewDistributedStore creates a new distributed telemetry store
func NewDistributedStore(config config.Config, timeStore timestore.Store) (Store, error) {
	log.Info("Creating distributed telemetry store")
	messageStore, err := messagestore.NewDistributedStore(primitiveName, config, timeStore)
	if err != nil {
		return nil, err
	}
	return &distributedStore{
		messageStore: messageStore,
	}, nil
}

// NewLocalStore returns a new local telemetry store
func NewLocalStore(timeStore timestore.Store) (Store, error) {
	messageStore, err := messagestore.NewLocalStore(primitiveName, timeStore)
	if err != nil {
		return nil, err
	}
	return &distributedStore{
		messageStore: messageStore,
	}, nil
}

// Store is interface for telemetry store
type Store interface {
	io.Closer

	// Gets a telemetry message based on a given ID
	Get(ID, ...GetOption) (*sb.TelemetryMessage, error)

	// Puts a telemetry message to the store
	Put(*sb.TelemetryMessage) error

	// Removes a telemetry message from the store
	Delete(ID, ...DeleteOption) error

	// List all of the last up to date telemetry messages
	List(ch chan<- sb.TelemetryMessage) error

	// Watch watches telemetry updates
	Watch(ch chan<- sb.TelemetryMessage, opts ...WatchOption) error

	// Clear deletes all entries from the store
	Clear() error
}

// GetOption is a message store get option
type GetOption messagestore.GetOption

// getOptions is a struct of message get options
type getOptions struct {
	revision Revision
}

// WithRevision returns a GetOption that ensures the retrieved entry is newer than the given revision
func WithRevision(revision Revision) GetOption {
	return messagestore.WithRevision(toMessageRevision(revision))
}

// DeleteOption is a message store delete option
type DeleteOption messagestore.DeleteOption

// IfRevision returns a delete option that deletes the message only if its revision matches the given revision
func IfRevision(revision Revision) DeleteOption {
	return messagestore.IfRevision(toMessageRevision(revision))
}

// WatchOption is a telemetry store watch option
type WatchOption messagestore.WatchOption

// WithReplay returns a watch option that replays existing telemetry
func WithReplay() WatchOption {
	return messagestore.WithReplay()
}

// distributedStore is an Atomix based store
type distributedStore struct {
	messageStore messagestore.Store
}

func (s *distributedStore) Get(id ID, opts ...GetOption) (*sb.TelemetryMessage, error) {
	messageOpts := make([]messagestore.GetOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt
	}
	entry, err := s.messageStore.Get(messagestore.Key(id.String()), messageOpts...)
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}
	return entry.GetTelemetry(), nil
}

func (s *distributedStore) Put(telemetry *sb.TelemetryMessage) error {
	entry := &message.MessageEntry{
		Message: &message.MessageEntry_Telemetry{
			Telemetry: telemetry,
		},
	}
	return s.messageStore.Put(getKey(telemetry), entry)
}

func (s *distributedStore) Delete(id ID, opts ...DeleteOption) error {
	messageOpts := make([]messagestore.DeleteOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt
	}
	return s.messageStore.Delete(messagestore.Key(id.String()), messageOpts...)
}

func (s *distributedStore) List(ch chan<- sb.TelemetryMessage) error {
	entryCh := make(chan message.MessageEntry)
	if err := s.messageStore.List(entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			ch <- *entry.GetTelemetry()
		}
	}()
	return nil
}

func (s *distributedStore) Watch(ch chan<- sb.TelemetryMessage, opts ...WatchOption) error {
	messageOpts := make([]messagestore.WatchOption, len(opts))
	for i, opt := range opts {
		messageOpts[i] = opt
	}

	watchCh := make(chan message.MessageEntry)
	if err := s.messageStore.Watch(watchCh, messageOpts...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range watchCh {
			ch <- *entry.GetTelemetry()
		}
	}()
	return nil
}

func (s *distributedStore) Clear() error {
	return s.messageStore.Clear()
}

func (s *distributedStore) Close() error {
	return s.messageStore.Close()
}

func toMessageRevision(revision Revision) messagestore.Revision {
	return messagestore.Revision{
		Term:      messagestore.Term(revision.Term),
		Timestamp: messagestore.Timestamp(revision.Timestamp),
	}
}

func getID(telemetry *sb.TelemetryMessage) ID {
	var ecgi sb.ECGI
	var crnti string
	switch telemetry.MessageType {
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		ecgi = *telemetry.GetRadioMeasReportPerUE().GetEcgi()
		crnti = telemetry.GetRadioMeasReportPerUE().GetCrnti()
	case sb.MessageType_RADIO_MEAS_REPORT_PER_CELL:
		ecgi = *telemetry.GetRadioMeasReportPerCell().GetEcgi()
	case sb.MessageType_SCHED_MEAS_REPORT_PER_UE:
	case sb.MessageType_SCHED_MEAS_REPORT_PER_CELL:
		ecgi = *telemetry.GetSchedMeasReportPerCell().GetEcgi()
	}
	return ID{
		PlmnID:      ecgi.PlmnId,
		Ecid:        ecgi.Ecid,
		Crnti:       crnti,
		MessageType: telemetry.GetMessageType(),
	}
}

func getKey(telemetry *sb.TelemetryMessage) messagestore.Key {
	return messagestore.Key(getID(telemetry).String())
}
