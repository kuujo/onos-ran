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

package requests

import (
	"context"
	"github.com/google/uuid"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sync"
)

type deviceStoreState struct {
	mastership     *mastership.State
	commitIndex    Index
	commitWatchers map[string]chan<- Index
	ackWatchers    map[string]chan<- Index
	ackIndex       Index
	mu             sync.RWMutex
}

func (s *deviceStoreState) setMastership(mastership *mastership.State) {
	s.mastership = mastership
}

func (s *deviceStoreState) getMastership() *mastership.State {
	return s.mastership
}

func (s *deviceStoreState) setCommitIndex(index Index) {
	if index > s.commitIndex {
		s.commitIndex = index
		for _, watcher := range s.commitWatchers {
			go func(ch chan<- Index) {
				ch <- index
			}(watcher)
		}
	}
}

func (s *deviceStoreState) getCommitIndex() Index {
	return s.commitIndex
}

func (s *deviceStoreState) watchCommitIndex(ctx context.Context, ch chan<- Index) {
	id := uuid.New().String()
	commitCh := make(chan Index)
	s.commitWatchers[id] = commitCh
	go func() {
		var commitIndex Index
		for commit := range commitCh {
			if commit > commitIndex {
				ch <- commit
				commitIndex = commit
			}
		}
	}()
	go func() {
		<-ctx.Done()
		delete(s.commitWatchers, id)
	}()
}

func (s *deviceStoreState) setAckIndex(index Index) {
	if index > s.ackIndex {
		s.ackIndex = index
		for _, watcher := range s.ackWatchers {
			go func(ch chan<- Index) {
				ch <- index
			}(watcher)
		}
	}
}

func (s *deviceStoreState) getAckIndex() Index {
	return s.ackIndex
}

func (s *deviceStoreState) watchAckIndex(ctx context.Context, ch chan<- Index) {
	id := uuid.New().String()
	ackCh := make(chan Index)
	s.ackWatchers[id] = ackCh
	go func() {
		var ackIndex Index
		for ack := range ackCh {
			if ack > ackIndex {
				ch <- ack
				ackIndex = ack
			}
		}
	}()
	go func() {
		<-ctx.Done()
		delete(s.ackWatchers, id)
	}()
}
