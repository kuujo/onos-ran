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
	"fmt"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/cluster"
	"github.com/spaolacci/murmur3"
	"io"
	"sync"
)

// Term is a monotonically increasing mastership term
type Term uint64

// Key is a mastership election key
type Key struct {
	PlmnID string
	Ecid   string
	Crnti  string
}

// PartitionID is the partition identifier
type PartitionID uint32

// String returns the mastership election key as a string
func (k Key) String() string {
	return fmt.Sprintf("%s:%s:%s", k.PlmnID, k.Ecid, k.Crnti)
}

// Hash returns the mastership election key as a hash
func (k Key) Hash() uint32 {
	return murmur3.Sum32([]byte(k.String()))
}

// Store is the mastership store
type Store interface {
	io.Closer

	// NodeID returns the local node identifier used in mastership elections
	NodeID() cluster.NodeID

	// IsMaster returns a boolean indicating whether the local node is the master for the given key
	IsMaster(key Key) (bool, error)

	// Watch watches the store for mastership changes
	Watch(Key, chan<- Mastership) error
}

// Mastership contains information about a mastership epoch
type Mastership struct {
	// PartitionID is the mastership partition identifier
	PartitionID PartitionID

	// Term is the mastership term
	Term Term

	// Master is the NodeID of the master for the key
	Master cluster.NodeID
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
		newElection: func(id PartitionID) (mastershipElection, error) {
			return newDistributedElection(id, database)
		},
		elections: make(map[PartitionID]mastershipElection),
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
		newElection: func(id PartitionID) (mastershipElection, error) {
			return newLocalElection(id, nodeID, address)
		},
		elections: make(map[PartitionID]mastershipElection),
	}, nil
}

// distributedStore is the default implementation of the NetworkConfig store
type distributedStore struct {
	nodeID      cluster.NodeID
	partitions  int
	newElection func(PartitionID) (mastershipElection, error)
	elections   map[PartitionID]mastershipElection
	mu          sync.RWMutex
}

// getElection gets the mastership election for the given partition
func (s *distributedStore) getElection(key Key) (mastershipElection, error) {
	partitionID := PartitionID(key.Hash() % uint32(s.partitions))
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

func (s *distributedStore) IsMaster(key Key) (bool, error) {
	election, err := s.getElection(key)
	if err != nil {
		return false, err
	}
	return election.isMaster()
}

func (s *distributedStore) Watch(key Key, ch chan<- Mastership) error {
	election, err := s.getElection(key)
	if err != nil {
		return err
	}
	return election.watch(ch)
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
