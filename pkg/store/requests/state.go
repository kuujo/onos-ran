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
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sync"
)

type deviceStoreState struct {
	mastership  *mastership.State
	commitIndex Index
	ackIndex    Index
	mu          sync.RWMutex
}

func (s *deviceStoreState) setMastership(mastership *mastership.State) {
	s.mu.Lock()
	s.mastership = mastership
	s.mu.Unlock()
}

func (s *deviceStoreState) getMastership() *mastership.State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.mastership
}

func (s *deviceStoreState) setCommitIndex(index Index) {
	s.mu.Lock()
	if index > s.commitIndex {
		s.commitIndex = index
	}
	s.mu.Unlock()
}

func (s *deviceStoreState) getCommitIndex() Index {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.commitIndex
}

func (s *deviceStoreState) setAckIndex(index Index) {
	s.mu.Lock()
	if index > s.ackIndex {
		s.ackIndex = index
	}
	s.mu.Unlock()
}

func (s *deviceStoreState) getAckIndex() Index {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ackIndex
}
