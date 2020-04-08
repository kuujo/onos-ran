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
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"io"
)

// NewStore creates a new time store
func NewStore(mastershipStore mastership.Store) (Store, error) {
	return &store{
		mastership: mastershipStore,
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
}

func (s *store) GetLogicalClock(key Key) (LogicalClock, error) {
	election, err := s.mastership.GetElection(mastership.Key(key))
	if err != nil {
		return nil, err
	}
	return newLogicalClock(election), nil
}

func (s *store) Close() error {
	return nil
}

var _ Store = &store{}
