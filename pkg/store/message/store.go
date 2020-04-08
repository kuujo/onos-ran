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
	timestore "github.com/onosproject/onos-ric/pkg/store/time"
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
func NewDistributedStore(name string, config config.Config, db atomix.DatabaseType, timeStore timestore.Store) (Store, error) {
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
		dist:    messages,
		entries: make(map[Key]entryStore),
		time:    timeStore,
	}
	if err := store.open(); err != nil {
		return nil, err
	}
	return store, nil
}

// NewLocalStore returns a new local message store
func NewLocalStore(name string, timeStore timestore.Store) (Store, error) {
	_, address := atomix.StartLocalNode()
	return newLocalStore(address, name, timeStore)
}

// newLocalStore creates a new local message store
func newLocalStore(address net.Address, name string, timeStore timestore.Store) (Store, error) {
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
		entries: make(map[Key]entryStore),
		time:    timeStore,
	}, nil
}

// Store is interface for message store
type Store interface {
	io.Closer

	// Gets a message based on a given ID
	Get(Key, ...GetOption) (*message.MessageEntry, error)

	// Puts a message to the store
	Put(Key, *message.MessageEntry) error

	// Removes a message from the store
	Delete(Key, ...DeleteOption) error

	// List all of the last up to date messages
	List(ch chan<- message.MessageEntry) error

	// Watch watches message updates
	Watch(ch chan<- message.MessageEntry, opts ...WatchOption) error

	// Clear deletes all entries from the store
	Clear() error
}

// GetOption is a message store get option
type GetOption interface {
	apply(options *getOptions)
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

func (o *getRevisionOption) apply(options *getOptions) {
	options.revision = o.revision
}

// DeleteOption is a message store delete option
type DeleteOption interface {
	apply(options *deleteOptions)
}

// deleteOptions is a struct of message delete options
type deleteOptions struct {
	revision Revision
}

// IfRevision returns a delete option that deletes the message only if its revision matches the given revision
func IfRevision(revision Revision) DeleteOption {
	return &deleteRevisionOption{
		revision: revision,
	}
}

type deleteRevisionOption struct {
	revision Revision
}

func (o *deleteRevisionOption) apply(options *deleteOptions) {
	options.revision = o.revision
}

// WatchOption is a message store watch option
type WatchOption interface {
	apply(options *watchOptions)
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

func (o *watchReplayOption) apply(options *watchOptions) {
	options.replay = o.replay
}

// distributedStore is an Atomix based store
type distributedStore struct {
	dist    _map.Map
	entries map[Key]entryStore
	time    timestore.Store
	mu      sync.RWMutex
}

// open opens the store
func (s *distributedStore) open() error {
	ch := make(chan *_map.Event)
	if err := s.dist.Watch(context.Background(), ch, _map.WithReplay()); err != nil {
		return err
	}
	go func() {
		for event := range ch {
			s.mu.RLock()
			store, ok := s.entries[Key(event.Entry.Key)]
			s.mu.RUnlock()
			if ok {
				if entry, err := decodeEntry(event.Entry); err == nil {
					tombstone := event.Type == _map.EventRemoved
					go store.update(entry, tombstone)
				}
			}
		}
	}()
	return nil
}

// getKeyStore uses a double checked lock to get or create a key store
func (s *distributedStore) getEntryStore(key Key) entryStore {
	s.mu.RLock()
	store, ok := s.entries[key]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		store, ok = s.entries[key]
		if !ok {
			store = newEntryStore(s.dist, key)
			s.entries[key] = store
		}
		s.mu.Unlock()
	}
	return store
}

func (s *distributedStore) Get(key Key, opts ...GetOption) (*message.MessageEntry, error) {
	options := &getOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	return s.getEntryStore(key).get(options.revision)
}

func (s *distributedStore) GetRevision(key Key, revision Revision) (*message.MessageEntry, error) {
	return s.getEntryStore(key).get(revision)
}

func (s *distributedStore) Put(key Key, entry *message.MessageEntry) error {
	clock, err := s.time.GetLogicalClock(timestore.Key(key))
	if err != nil {
		return err
	}
	timestamp, err := clock.GetTimestamp()
	if err != nil {
		return err
	}
	entry.Term = uint64(timestamp.Epoch)
	entry.Timestamp = uint64(timestamp.LogicalTime)
	return s.getEntryStore(key).put(entry)
}

func (s *distributedStore) Delete(key Key, opts ...DeleteOption) error {
	options := &deleteOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	return s.getEntryStore(key).delete(options.revision)
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

func (s *distributedStore) Watch(ch chan<- message.MessageEntry, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	watchOptions := []_map.WatchOption{}
	if options.replay {
		watchOptions = append(watchOptions, _map.WithReplay())
	}

	watchCh := make(chan *_map.Event)
	if err := s.dist.Watch(context.Background(), watchCh, watchOptions...); err != nil {
		return err
	}
	go func() {
		defer close(ch)
		for event := range watchCh {
			if event.Type != _map.EventRemoved {
				if entry, err := decodeEntry(event.Entry); err == nil {
					ch <- *entry
				}
			}
		}
	}()
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

// Key is a mesasge key
type Key string

// String returns the key as a string
func (k Key) String() string {
	return string(k)
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