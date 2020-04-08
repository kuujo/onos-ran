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
	"container/list"
	"context"
	"errors"
	"fmt"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-ric/api/store/message"
	"sync"
	"time"
)

func newEntryStore(dist _map.Map, key Key) entryStore {
	return &distributedEntryStore{
		dist:     dist,
		key:      key,
		watchers: list.New(),
		waiters:  list.New(),
	}
}

// entryWatcher is a watcher for a single entry
type entryWatcher struct {
	ch       chan<- message.MessageEntry
	revision Revision
}

// entryStore is a store for a single message entry
type entryStore interface {
	// update updates an entry in the store
	update(updatedEntry *message.MessageEntry, tombstone bool)

	// get gets a message from the store
	get(revision Revision) (*message.MessageEntry, error)

	// put puts a message in the store
	put(entry *message.MessageEntry) error

	// delete deletes a message from the store
	delete(revision Revision) error

	// watch watches the store
	watch(ch chan<- message.MessageEntry)
}

// distributedEntryStore is an implementation of the entryStore interface
type distributedEntryStore struct {
	dist     _map.Map
	key      Key
	cache    *message.MessageEntry
	watchers *list.List
	waiters  *list.List
	mu       sync.RWMutex
}

func (s *distributedEntryStore) update(updatedEntry *message.MessageEntry, tombstone bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the current entry from the cache
	currentEntry := s.cache

	// Compute the current and updated revision
	var currentRevision Revision
	if currentEntry != nil {
		currentRevision = newRevision(currentEntry.Term, currentEntry.Timestamp)
	}
	updateRevision := newRevision(updatedEntry.Term, updatedEntry.Timestamp)

	// If the updated revision is newer than the current revision, update the cache
	if updateRevision.isNewerThan(currentRevision) {
		s.cache = updatedEntry

		// Trigger all watchers
		element := s.watchers.Front()
		for element != nil {
			watcher := element.Value.(*entryWatcher)
			watcher.ch <- *updatedEntry
			watcher.revision = updateRevision
			element = element.Next()
		}

		// Trigger any channels waiting for the entry.
		waiter := s.waiters.Front()
		for waiter != nil {
			ctx := waiter.Value.(*waiterContext)
			if updateRevision.isEqualToOrNewerThan(ctx.revision) {
				ctx.ch <- updatedEntry
				close(ctx.ch)
				next := waiter.Next()
				s.waiters.Remove(waiter)
				waiter = next
			} else {
				break
			}
		}
	} else if currentRevision.isEqualTo(updateRevision) {
		// Ensure all watchers have seen the revision
		element := s.watchers.Front()
		for element != nil {
			watcher := element.Value.(*entryWatcher)
			if updateRevision.isNewerThan(watcher.revision) {
				watcher.ch <- *updatedEntry
				watcher.revision = updateRevision
			}
			element = element.Next()
		}
	}
}

func (s *distributedEntryStore) get(revision Revision) (*message.MessageEntry, error) {
	// Read the entry from the cache
	s.mu.RLock()
	messageEntry := s.cache
	s.mu.RUnlock()

	// If the entry term is greater than the requested term or the terms are equal and the entry timestamp
	// is greater than or equal to the requested timestamp, the entry is up to date.
	if messageEntry != nil &&
		(Term(messageEntry.Term) > revision.Term ||
			(Term(messageEntry.Term) == revision.Term && Timestamp(messageEntry.Timestamp) >= revision.Timestamp)) {
		return messageEntry, nil
	}

	// If the key is not present in the cache or is older than the requested key, read it from the distributed store.
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	entry, err := s.dist.Get(ctx, s.key.String())
	cancel()
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}

	// Decode the stored entry
	messageEntry, err = decodeEntry(entry)
	if err != nil {
		return nil, err
	}

	// Again, determine whether the message entry meets the requested revision info.
	messageRevision := newRevision(messageEntry.Term, messageEntry.Timestamp)
	if messageRevision.isEqualToOrNewerThan(revision) {
		// Cache the entry if it's up-to-date.
		s.mu.Lock()
		s.cache = messageEntry
		s.mu.Unlock()
		return messageEntry, nil
	}

	// If the entry is not up to date, enqueue a waiter to wait for the update to be propagated
	ch := make(chan *message.MessageEntry)
	waiter := &waiterContext{
		ch:       ch,
		revision: revision,
	}

	// Add the waiter to the appropriate position in the queue and wait for the event.
	// The waiter queue is a linked list sorted by logical time from lowest to highest.
	s.mu.Lock()
	pos := s.waiters.Front()
	for pos != nil {
		ctx := pos.Value.(*waiterContext)
		if ctx.revision.isNewerThan(revision) {
			s.waiters.InsertBefore(waiter, pos)
		} else if pos.Next() == nil {
			s.waiters.PushBack(waiter)
		} else {
			pos = pos.Next()
		}
	}
	s.mu.Unlock()

	select {
	case e := <-ch:
		return e, nil
	case <-time.After(requestTimeout):
		return nil, errors.New("message get timed out")
	}
}

func (s *distributedEntryStore) put(update *message.MessageEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Store the update in the local cache
	s.cache = update

	// Propagate the change to the store asynchronously
	s.enqueuePut(update)
	return nil
}

func (s *distributedEntryStore) delete(revision Revision) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create a timestamped entry for the store
	entry := &message.MessageEntry{
		Term:      uint64(revision.Term),
		Timestamp: uint64(revision.Timestamp),
	}

	// Store the update in the local cache
	s.cache = entry

	// Enqueue the update to be written to the distributed store
	s.enqueueDelete(revision)
	return nil
}

func (s *distributedEntryStore) enqueuePut(entry *message.MessageEntry) {
	go s.writePut(entry)
}

func (s *distributedEntryStore) requeuePut(update *message.MessageEntry) {
	time.AfterFunc(retryInterval, func() {
		s.writePut(update)
	})
}

func (s *distributedEntryStore) writePut(update *message.MessageEntry) {
	// Get the current entry from the store
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	entry, err := s.dist.Get(ctx, s.key.String())
	cancel()
	if err != nil {
		fmt.Println(err)
		s.requeuePut(update)
		return
	}

	// If the entry is already stored, verify the updated is newer than the stored entry
	if entry != nil {
		// Decode the entry and fail if decoding fails
		current, err := decodeEntry(entry)
		if err != nil {
			fmt.Println(err)
			return
		}

		currentRevision := newRevision(current.Term, current.Timestamp)
		updateRevision := newRevision(update.Term, update.Timestamp)

		// If the stored entry is newer, ignore the put
		if currentRevision.isNewerThan(updateRevision) {
			return
		}
	}

	// Encode the updated entry and fail if encoding fails
	bytes, err := encodeEntry(update)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Update the stored entry using an optimistic lock
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	if entry == nil {
		_, err = s.dist.Put(ctx, s.key.String(), bytes, _map.IfNotSet())
	} else {
		_, err = s.dist.Put(ctx, s.key.String(), bytes, _map.IfVersion(entry.Version))
	}

	// If the update failed, requeue it
	if err != nil {
		s.requeuePut(update)
		return
	}
}

func (s *distributedEntryStore) enqueueDelete(revision Revision) {
	go s.writeDelete(revision)
}

func (s *distributedEntryStore) requeueDelete(revision Revision) {
	time.AfterFunc(retryInterval, func() {
		s.writeDelete(revision)
	})
}

func (s *distributedEntryStore) writeDelete(revision Revision) {
	// Get the current entry from the store
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	entry, err := s.dist.Get(ctx, s.key.String())
	cancel()
	if err != nil {
		fmt.Println(err)
		s.requeueDelete(revision)
		return
	}

	// If the entry is already deleted, ignore the delete
	if entry == nil {
		return
	}

	// Decode the current entry and fail if decoding fails
	current, err := decodeEntry(entry)
	if err != nil {
		fmt.Println(err)
		return
	}

	// If the current entry is newer than the update entry, ignore the update
	currentRevision := newRevision(current.Term, current.Timestamp)
	if currentRevision.isNewerThan(revision) {
		return
	}

	// Remove the stored entry using an optimistic lock
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err = s.dist.Remove(ctx, s.key.String(), _map.IfVersion(entry.Version))

	// If the remove failed, requeue it
	if err != nil {
		s.requeueDelete(revision)
		return
	}
}

func (s *distributedEntryStore) watch(ch chan<- message.MessageEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.watchers.PushBack(&entryWatcher{
		ch: ch,
	})
}

func decodeEntry(entry *_map.Entry) (*message.MessageEntry, error) {
	messageEntry := &message.MessageEntry{}
	if err := proto.Unmarshal(entry.Value, messageEntry); err != nil {
		return nil, err
	}
	return messageEntry, nil
}

func encodeEntry(entry *message.MessageEntry) ([]byte, error) {
	return proto.Marshal(entry)
}

type waiterContext struct {
	ch       chan *message.MessageEntry
	revision Revision
}
