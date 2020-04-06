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
	"container/list"
	"context"
	"errors"
	"fmt"
	"github.com/atomix/go-client/pkg/client/map"
	"github.com/gogo/protobuf/proto"
	"github.com/onosproject/onos-ric/api/store/telemetry"
	"sync"
	"time"
)

func newEntryStore(dist _map.Map) entryStore {
	return &distributedEntryStore{
		dist:    dist,
		waiters: list.New(),
	}
}

// entryStore is a store for a single telemetry entry
type entryStore interface {
	// update updates an entry in the store
	update(updatedEntry *telemetry.TelemetryEntry, tombstone bool)

	// get gets a telemetry message from the store
	get(revision Revision) (*telemetry.TelemetryEntry, error)

	// put puts a telemetry message in the store
	put(entry *telemetry.TelemetryEntry) error

	// delete deletes a telemetry message from the store
	delete(revision Revision) error
}

// distributedEntryStore is an implementation of the entryStore interface
type distributedEntryStore struct {
	dist    _map.Map
	cache   *telemetry.TelemetryEntry
	waiters *list.List
	mu      sync.RWMutex
}

func (s *distributedEntryStore) update(updatedEntry *telemetry.TelemetryEntry, tombstone bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Get the current entry from the cache
	currentEntry := s.cache

	// If the key is not present in the cache or the updated entry is newer, update the cache
	// Otherwise, insert the entry into the cache.
	if currentEntry == nil || updatedEntry.Term > currentEntry.Term ||
		(updatedEntry.Term == currentEntry.Term &&
			updatedEntry.Timestamp > currentEntry.Timestamp) {
		s.cache = updatedEntry

		// Trigger any channels waiting for the entry.
		waiter := s.waiters.Front()
		for waiter != nil {
			ctx := waiter.Value.(*waiterContext)
			if Term(updatedEntry.Term) > ctx.term ||
				(Term(updatedEntry.Term) == ctx.term &&
					Timestamp(updatedEntry.Timestamp) > ctx.timestamp) {
				ctx.ch <- updatedEntry
				close(ctx.ch)
				next := waiter.Next()
				s.waiters.Remove(waiter)
				waiter = next
			} else {
				break
			}
		}
	}
}

func (s *distributedEntryStore) get(revision Revision) (*telemetry.TelemetryEntry, error) {
	// Read the entry from the cache
	s.mu.RLock()
	telemetryEntry := s.cache
	s.mu.RUnlock()

	// If the entry term is greater than the requested term or the terms are equal and the entry timestamp
	// is greater than or equal to the requested timestamp, the entry is up to date.
	if telemetryEntry != nil &&
		(Term(telemetryEntry.Term) > revision.Term ||
			(Term(telemetryEntry.Term) == revision.Term && Timestamp(telemetryEntry.Timestamp) >= revision.Timestamp)) {
		return telemetryEntry, nil
	}

	// If the key is not present in the cache or is older than the requested key, read it from the distributed store.
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	entry, err := s.dist.Get(ctx, revision.ID.String())
	cancel()
	if err != nil {
		return nil, err
	} else if entry == nil {
		return nil, nil
	}

	// Decode the stored entry
	telemetryEntry, err = decodeEntry(entry)
	if err != nil {
		return nil, err
	}

	// Again, determine whether the telemetry entry meets the requested revision info.
	if Term(telemetryEntry.Term) > revision.Term ||
		(Term(telemetryEntry.Term) == revision.Term && Timestamp(telemetryEntry.Timestamp) >= revision.Timestamp) {
		// Cache the entry if it's up-to-date.
		s.mu.Lock()
		s.cache = telemetryEntry
		s.mu.Unlock()
		return telemetryEntry, nil
	}

	// If the entry is not up to date, enqueue a waiter to wait for the update to be propagated
	ch := make(chan *telemetry.TelemetryEntry)
	waiter := &waiterContext{
		term:      revision.Term,
		timestamp: revision.Timestamp,
		ch:        ch,
	}

	// Add the waiter to the appropriate position in the queue and wait for the event.
	// The waiter queue is a linked list sorted by logical time from lowest to highest.
	s.mu.Lock()
	pos := s.waiters.Front()
	for pos != nil {
		ctx := pos.Value.(*waiterContext)
		if revision.Term > ctx.term || (revision.Term == ctx.term && revision.Timestamp >= ctx.timestamp) {
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
		return nil, errors.New("telemetry get timed out")
	}
}

func (s *distributedEntryStore) put(update *telemetry.TelemetryEntry) error {
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
	entry := &telemetry.TelemetryEntry{
		Term:      uint64(revision.Term),
		Timestamp: uint64(revision.Timestamp),
	}

	// Store the update in the local cache
	s.cache = entry

	// Enqueue the update to be written to the distributed store
	s.enqueueDelete(revision)
	return nil
}

func (s *distributedEntryStore) enqueuePut(entry *telemetry.TelemetryEntry) {
	go s.writePut(entry)
}

func (s *distributedEntryStore) requeuePut(update *telemetry.TelemetryEntry) {
	time.AfterFunc(retryInterval, func() {
		s.writePut(update)
	})
}

func (s *distributedEntryStore) writePut(update *telemetry.TelemetryEntry) {
	id := newID(update.Message)

	// Get the current entry from the store
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	entry, err := s.dist.Get(ctx, id.String())
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

		// If the stored entry is newer, ignore the put
		if current.Term > update.Term || (current.Term == update.Term && current.Timestamp > update.Timestamp) {
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
		_, err = s.dist.Put(ctx, id.String(), bytes, _map.IfNotSet())
	} else {
		_, err = s.dist.Put(ctx, id.String(), bytes, _map.IfVersion(entry.Version))
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
	entry, err := s.dist.Get(ctx, revision.ID.String())
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
	if Term(current.Term) > revision.Term || (Term(current.Term) == revision.Term && Timestamp(current.Timestamp) > revision.Timestamp) {
		return
	}

	// Remove the stored entry using an optimistic lock
	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err = s.dist.Remove(ctx, revision.ID.String(), _map.IfVersion(entry.Version))

	// If the remove failed, requeue it
	if err != nil {
		s.requeueDelete(revision)
		return
	}
}

func decodeEntry(entry *_map.Entry) (*telemetry.TelemetryEntry, error) {
	telemetryEntry := &telemetry.TelemetryEntry{}
	if err := proto.Unmarshal(entry.Value, telemetryEntry); err != nil {
		return nil, err
	}
	return telemetryEntry, nil
}

func encodeEntry(entry *telemetry.TelemetryEntry) ([]byte, error) {
	return proto.Marshal(entry)
}
