// Copyright 2019-present Open Networking Foundation.
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

package mastership

import (
	"github.com/DataDog/mmh3"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/cluster"
	"io"
	"sync"
)

// Term is a monotonically increasing mastership term
type Term uint64

// Key is a mastership election key
type Key string

// Hash returns the mastership election key as a hash
func (k Key) Hash() uint32 {
	return mmh3.Hash32([]byte(k))
}

// PartitionID is the partition identifier
type PartitionID uint32

// Store is the mastership store
type Store interface {
	io.Closer

	// NodeID returns the local node identifier used in mastership elections
	NodeID() cluster.NodeID

	// GetElection gets the mastership election for the given key
	GetElection(key Key) (Election, error)
}

// NewDistributedStore returns a new distributed Store
func NewDistributedStore(config config.Config) (Store, error) {
	database, err := atomix.GetDatabase(config.Atomix, config.Atomix.GetDatabase(atomix.DatabaseTypeConsensus))
	if err != nil {
		return nil, err
	}
	return &distributedStore{
		nodeID:     cluster.GetNodeID(),
		partitions: config.Mastership.GetPartitions(),
		newElection: func(id PartitionID) (Election, error) {
			return newDistributedElection(id, database)
		},
		elections: make(map[PartitionID]Election),
	}, nil
}

var localAddresses = make(map[string]net.Address)

// NewLocalStore returns a new local election store
func NewLocalStore(clusterID string, nodeID cluster.NodeID) (Store, error) {
	address, ok := localAddresses[clusterID]
	if !ok {
		_, address = atomix.StartLocalNode()
		localAddresses[clusterID] = address
	}
	return newLocalStore(nodeID, address)
}

// newLocalStore returns a new local mastership store
func newLocalStore(nodeID cluster.NodeID, address net.Address) (Store, error) {
	return &distributedStore{
		nodeID:     nodeID,
		partitions: 16,
		newElection: func(id PartitionID) (Election, error) {
			return newLocalElection(id, nodeID, address)
		},
		elections: make(map[PartitionID]Election),
	}, nil
}

// distributedStore is the default implementation of the NetworkConfig store
type distributedStore struct {
	nodeID      cluster.NodeID
	partitions  int
	newElection func(PartitionID) (Election, error)
	elections   map[PartitionID]Election
	mu          sync.RWMutex
}

// getPartitionFor gets the partition for the given key
func getPartitionFor(key Key, partitions int) PartitionID {
	return PartitionID(key.Hash() % uint32(partitions))
}

// GetElection gets the mastership election for the given key
func (s *distributedStore) GetElection(key Key) (Election, error) {
	partitionID := getPartitionFor(key, s.partitions)
	s.mu.RLock()
	election, ok := s.elections[partitionID]
	s.mu.RUnlock()
	if !ok {
		s.mu.Lock()
		election, ok = s.elections[partitionID]
		if !ok {
			e, err := s.newElection(partitionID)
			if err != nil {
				s.mu.Unlock()
				return nil, err
			}
			election = e
			s.elections[partitionID] = election
		}
		s.mu.Unlock()
	}
	return election, nil
}

func (s *distributedStore) NodeID() cluster.NodeID {
	return s.nodeID
}

func (s *distributedStore) Close() error {
	var returnErr error
	for _, election := range s.elections {
		if err := election.Close(); err != nil && returnErr == nil {
			returnErr = err
		}
	}
	return returnErr
}

var _ Store = &distributedStore{}