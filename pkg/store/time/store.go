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

package time

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"io"
	"sync"
)

var log = logging.GetLogger("store", "time")

// NewStore creates a new time store
func NewStore(mastershipStore mastership.Store) (Store, error) {
	return &store{
		mastership: mastershipStore,
		clocks:     make(map[Key]LogicalClock),
	}, nil
}

// Key is the name of a logical clock
type Key string

// Store is a store for managing logical clocks
type Store interface {
	io.Closer

	// GetLogicalClock gets a mastership based logical clock
	GetLogicalClock(Key) (LogicalClock, error)
}

// store is the default implementation of the time store
type store struct {
	mastership mastership.Store
	clocks     map[Key]LogicalClock
	mu         sync.RWMutex
}

func (s *store) GetLogicalClock(key Key) (LogicalClock, error) {
	s.mu.RLock()
	clock, ok := s.clocks[key]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		clock, ok = s.clocks[key]
		if !ok {
			election, err := s.mastership.GetElection(mastership.Key(key))
			if err != nil {
				return nil, err
			}
			clock = newLogicalClock(election)
			s.clocks[key] = clock
		}
	}
	return clock, nil
}

func (s *store) Close() error {
	return nil
}

var _ Store = &store{}
