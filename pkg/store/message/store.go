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

package message

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-ric/api/store/message"
	"github.com/onosproject/onos-ric/pkg/config"
	clocks "github.com/onosproject/onos-ric/pkg/store/time"
	"io"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/atomix/go-client/pkg/client/map"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
)

var log = logging.GetLogger("store", "message")

const requestTimeout = 15 * time.Second
const retryInterval = time.Second

// NewDistributedStore creates a new distributed message store
func NewDistributedStore(name string, config config.Config, db atomix.DatabaseType, timeStore clocks.Store) (Store, error) {
	log.Info("Creating distributed message store")
	database, err := atomix.GetDatabase(config.Atomix, config.Atomix.GetDatabase(db))

	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	messages, err := database.GetMap(ctx, name)
	if err != nil {
		return nil, err
	}
	store := &distributedStore{
		dist:     messages,
		entries:  make(map[ID]entryStore),
		watchers: make([]chan<- Event, 0),
		clocks:   timeStore,
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

// NewLocalStore returns a new local message store
func NewLocalStore(name string, timeStore clocks.Store) (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address, name, timeStore)
}

// newLocalStore creates a new local message store
func newLocalStore(address net.Address, name string, timeStore clocks.Store) (Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session, err := primitive.NewSession(ctx, primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}
	primitiveName := primitive.Name{
		Namespace: "local",
		Name:      name,
	}
	messages, err := _map.New(context.Background(), primitiveName, []*primitive.Session{session})
	if err != nil {
		return nil, err
	}

	return &distributedStore{
		dist:    messages,
		entries: make(map[ID]entryStore),
		clocks:  timeStore,
	}, nil
}

// Event is a message store event
type Event struct {
	// Type is the event type
	Type EventType
	// ID is the event message ID
	ID ID
	// Message is the event message
	Message message.MessageEntry
}

// EventType is a message store event type
type EventType string

const (
	// EventNone indicates an event that was not triggered but replayed
	EventNone EventType = ""
	// EventInsert indicates the message was inserted into the store
	EventInsert EventType = "insert"
	// EventUpdate indicates the message was updated in the store
	EventUpdate EventType = "update"
	// EventDelete indicates the message was deleted from the store
	EventDelete EventType = "remove"
)

// Store is interface for message store
type Store interface {
	io.Closer

	// Gets a message based on a given ID
	Get(PartitionKey, ID, ...GetOption) (*message.MessageEntry, error)

	// Puts a message to the store
	Put(PartitionKey, ID, *message.MessageEntry, ...PutOption) error

	// Removes a message from the store
	Delete(PartitionKey, ID, ...DeleteOption) error

	// List all of the last up to date messages
	List(ch chan<- message.MessageEntry) error

	// Watch watches message updates
	Watch(ch chan<- Event, opts ...WatchOption) error

	// Clear deletes all entries from the store
	Clear() error
}

// GetOption is a message store get option
type GetOption interface {
	configureGet(options *getOptions)
}

// getOptions is a struct of message get options
type getOptions struct {
	revision Revision
}

// WithRevision returns a GetOption that ensures the retrieved entry is newer than the given revision
func WithRevision(revision Revision) GetOption {
	return &getRevisionOption{
		revision: revision,
	}
}

type getRevisionOption struct {
	revision Revision
}

func (o *getRevisionOption) configureGet(options *getOptions) {
	options.revision = o.revision
}

// PutOption is a message store put option
type PutOption interface {
	configurePut(options *putOptions)
}

// putOptions is a struct of message put options
type putOptions struct{}

// DeleteOption is a message store delete option
type DeleteOption interface {
	configureDelete(options *deleteOptions)
}

// deleteOptions is a struct of message delete options
type deleteOptions struct {
	revision Revision
}

// IfRevision returns a delete option that deletes the message only if its revision matches the given revision
func IfRevision(revision Revision) DeleteOption {
	return &updateRevisionOption{
		revision: revision,
	}
}

type updateRevisionOption struct {
	revision Revision
}

func (o *updateRevisionOption) configureDelete(options *deleteOptions) {
	options.revision = o.revision
}

// WatchOption is a message store watch option
type WatchOption interface {
	configureWatch(options *watchOptions)
}

// watchOptions is a struct of message store watch options
type watchOptions struct {
	replay bool
}

// WithReplay returns a watch option that replays existing messages
func WithReplay() WatchOption {
	return &watchReplayOption{
		replay: true,
	}
}

// watchReplayOption is an option for configuring whether to replay message on watch calls
type watchReplayOption struct {
	replay bool
}

func (o *watchReplayOption) configureWatch(options *watchOptions) {
	options.replay = o.replay
}

// distributedStore is an Atomix based store
type distributedStore struct {
	dist     _map.Map
	entries  map[ID]entryStore
	watchers []chan<- Event
	clocks   clocks.Store
	mu       sync.RWMutex
}

// open opens the store
func (s *distributedStore) open() error {
	ch := make(chan *_map.Event)
	if err := s.dist.Watch(context.Background(), ch, _map.WithReplay()); err != nil {
		return err
	}
	go func() {
		for event := range ch {
			if entry, err := decodeEntry(event.Entry); err == nil {
				if store, err := s.getEntryStore(PartitionKey(entry.PartitionKey), ID(entry.Id)); err == nil {
					go func() {
						_ = store.update(Revision{}, entry, event.Type == _map.EventRemoved)
					}()
				}
			}
		}
	}()
	return nil
}

// getKeyStore uses a double checked lock to get or create a key store
func (s *distributedStore) getEntryStore(key PartitionKey, id ID) (entryStore, error) {
	s.mu.RLock()
	store, ok := s.entries[id]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		store, ok = s.entries[id]
		if !ok {
			clock, err := s.clocks.GetLogicalClock(clocks.Key(key))
			if err != nil {
				return nil, err
			}
			store = newEntryStore(s.dist, id, clock)
			for _, watcher := range s.watchers {
				store.watch(watcher)
			}
			s.entries[id] = store
		}
		s.mu.Unlock()
	}
	return store, nil
}

func (s *distributedStore) Get(key PartitionKey, id ID, opts ...GetOption) (*message.MessageEntry, error) {
	options := &getOptions{}
	for _, opt := range opts {
		opt.configureGet(options)
	}
	store, err := s.getEntryStore(key, id)
	if err != nil {
		return nil, err
	}
	return store.get(options.revision)
}

func (s *distributedStore) Put(key PartitionKey, id ID, entry *message.MessageEntry, opts ...PutOption) error {
	store, err := s.getEntryStore(key, id)
	if err != nil {
		return err
	}
	revision := newRevision(entry.Term, entry.Timestamp)
	return store.put(revision, entry)
}

func (s *distributedStore) Delete(key PartitionKey, id ID, opts ...DeleteOption) error {
	options := &deleteOptions{}
	for _, opt := range opts {
		opt.configureDelete(options)
	}
	store, err := s.getEntryStore(key, id)
	if err != nil {
		return err
	}
	return store.delete(options.revision)
}

func (s *distributedStore) List(ch chan<- message.MessageEntry) error {
	entryCh := make(chan *_map.Entry)
	if err := s.dist.Entries(context.Background(), entryCh); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for entry := range entryCh {
			message, err := decodeEntry(entry)
			if err == nil {
				ch <- *message
			}
		}
	}()
	return nil
}

func (s *distributedStore) Watch(ch chan<- Event, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.configureWatch(options)
	}

	// Add the watcher to all stores
	s.mu.Lock()
	s.watchers = append(s.watchers, ch)
	for _, store := range s.entries {
		store.watch(ch)
	}
	s.mu.Unlock()

	// If replay is enabled, list the entries in the store to replay messages
	if options.replay {
		entryCh := make(chan *_map.Entry)
		if err := s.dist.Entries(context.Background(), entryCh); err != nil {
			return err
		}
		go func() {
			for entry := range entryCh {
				if message, err := decodeEntry(entry); err == nil {
					if store, err := s.getEntryStore(PartitionKey(message.PartitionKey), ID(message.Id)); err == nil {
						_ = store.update(Revision{}, message, false)
					}
				}
			}
		}()
	}
	return nil
}

func (s *distributedStore) Clear() error {
	return s.dist.Clear(context.Background())
}

func (s *distributedStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	return s.dist.Close(ctx)
}

// PartitionKey is a message partition key
type PartitionKey string

// String returns the partition key as a string
func (k PartitionKey) String() string {
	return string(k)
}

// ID is a message ID
type ID string

// String returns the ID as a string
func (i ID) String() string {
	return string(i)
}

// Term is a message term
type Term uint64

// Timestamp is a message timestamp
type Timestamp uint64

// Revision is a message revision
type Revision struct {
	// Term is the term in which the revision was created
	Term Term
	// Timestamp is the timestamp at which the revision was created
	Timestamp Timestamp
}

// isEqualTo returns a bool indicating whether the revision is equal to the given revision
func (r Revision) isEqualTo(revision Revision) bool {
	return r.Term == revision.Term && r.Timestamp == revision.Timestamp
}

// isNewerThan returns a bool indicating whether the revision is newer than the given revision
func (r Revision) isNewerThan(revision Revision) bool {
	return r.Term > revision.Term || (r.Term == revision.Term && r.Timestamp > revision.Timestamp)
}

// isEqualToOrNewerThan returns a bool indicating whether the revision is equal to or newer than the given revision
func (r Revision) isEqualToOrNewerThan(revision Revision) bool {
	return r.Term > revision.Term || (r.Term == revision.Term && r.Timestamp >= revision.Timestamp)
}

// isZero returns a bool indicating whether the revision is not set
func (r Revision) isZero() bool {
	return r.Term == 0 || r.Timestamp == 0
}

func newRevision(term, timestamp uint64) Revision {
	return Revision{
		Term:      Term(term),
		Timestamp: Timestamp(timestamp),
	}
}
