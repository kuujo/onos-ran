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
	"context"
	"errors"
	"fmt"
	"github.com/atomix/go-client/pkg/client"
	"github.com/atomix/go-client/pkg/client/election"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util/net"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"io"
	"sync"
	"time"
)

// newDistributedElection returns a new distributed mastership election
func newDistributedElection(partitionID PartitionID, database *client.Database, cluster cluster.Cluster) (Election, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	election, err := database.GetElection(ctx, fmt.Sprintf("mastership-%d", partitionID), election.WithID(string(cluster.Node().ID)))
	cancel()
	if err != nil {
		return nil, err
	}
	return newMastershipElection(partitionID, election)
}

// newLocalElection returns a new local mastership election
func newLocalElection(partitionID PartitionID, nodeID cluster.NodeID, address net.Address) (Election, error) {
	name := primitive.Name{
		Namespace: "local",
		Name:      fmt.Sprintf("mastership-%d", partitionID),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	session, err := primitive.NewSession(ctx, primitive.Partition{ID: 1, Address: address})
	if err != nil {
		return nil, err
	}
	election, err := election.New(context.Background(), name, []*primitive.Session{session}, election.WithID(string(nodeID)))
	if err != nil {
		return nil, err
	}
	return newMastershipElection(partitionID, election)
}

// newDeviceMastershipElection creates and enters a new mastership election
func newMastershipElection(partitionID PartitionID, election election.Election) (Election, error) {
	mastershipElection := &distributedMastershipElection{
		partitionID: partitionID,
		election:    election,
		watchers:    make([]chan<- State, 0, 1),
	}
	if err := mastershipElection.enter(); err != nil {
		return nil, err
	}
	return mastershipElection, nil
}

// State contains information about a mastership term
type State struct {
	// PartitionID is the mastership partition identifier
	PartitionID PartitionID

	// Term is the mastership term
	Term Term

	// Master is the NodeID of the master for the key
	Master cluster.NodeID

	// Replicas is the replicas for the key
	Replicas []cluster.ReplicaID
}

// Election is an election for a single mastership
type Election interface {
	io.Closer

	// NodeID returns the local node identifier used in the election
	NodeID() cluster.NodeID

	// PartitionID returns the mastership election partition identifier
	PartitionID() PartitionID

	// GetState returns the mastership state
	GetState() (*State, error)

	// IsMaster returns a bool indicating whether the local node is the master for the device
	IsMaster() (bool, error)

	// Watch watches the election for changes
	Watch(ch chan<- State) error
}

// distributedMastershipElection is a persistent mastership election
type distributedMastershipElection struct {
	partitionID PartitionID
	election    election.Election
	mastership  *State
	watchers    []chan<- State
	mu          sync.RWMutex
}

func (e *distributedMastershipElection) NodeID() cluster.NodeID {
	return cluster.NodeID(e.election.ID())
}

func (e *distributedMastershipElection) PartitionID() PartitionID {
	return e.partitionID
}

// enter enters the election
func (e *distributedMastershipElection) enter() error {
	ch := make(chan *election.Event)
	if err := e.election.Watch(context.Background(), ch); err != nil {
		return err
	}

	// Enter the election to get the current leadership term
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	term, err := e.election.Enter(ctx)
	cancel()
	if err != nil {
		_ = e.election.Close(context.Background())
		return err
	}

	// Set the mastership term
	e.mu.Lock()
	replicas := make([]cluster.ReplicaID, 0, len(term.Candidates))
	for _, id := range term.Candidates {
		replicas = append(replicas, cluster.ReplicaID(id))
	}
	e.mastership = &State{
		PartitionID: e.partitionID,
		Master:      cluster.NodeID(term.Leader),
		Term:        Term(term.ID),
		Replicas:    replicas,
	}
	e.mu.Unlock()

	// Wait for the election event to be received before returning
	for event := range ch {
		if event.Term.ID == term.ID {
			go e.watchElection(ch)
			return nil
		}
	}

	_ = e.election.Close(context.Background())
	return errors.New("failed to enter election")
}

// watchElection watches the election events and updates mastership info
func (e *distributedMastershipElection) watchElection(ch <-chan *election.Event) {
	for event := range ch {
		var mastership *State
		e.mu.Lock()

		changed := false
		term := Term(event.Term.ID)
		if term > e.mastership.Term {
			changed = true
		} else if e.mastership.Term == Term(event.Term.ID) {
			if len(e.mastership.Replicas) != len(event.Term.Candidates) {
				changed = true
			} else {
				for _, candidateID := range event.Term.Candidates {
					for _, replicaID := range e.mastership.Replicas {
						if cluster.ReplicaID(candidateID) != replicaID {
							changed = true
							break
						}
					}
				}
			}
		}

		if changed {
			replicas := make([]cluster.ReplicaID, 0, len(event.Term.Candidates))
			for _, id := range event.Term.Candidates {
				replicas = append(replicas, cluster.ReplicaID(id))
			}
			mastership = &State{
				PartitionID: e.partitionID,
				Term:        term,
				Master:      cluster.NodeID(event.Term.Leader),
				Replicas:    replicas,
			}
			e.mastership = mastership
		}
		e.mu.Unlock()

		if mastership != nil {
			e.mu.RLock()
			for _, watcher := range e.watchers {
				watcher <- *mastership
			}
			e.mu.RUnlock()
		}
	}
}

func (e *distributedMastershipElection) GetState() (*State, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mastership, nil
}

func (e *distributedMastershipElection) IsMaster() (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.mastership == nil || string(e.mastership.Master) != e.election.ID() {
		return false, nil
	}
	return true, nil
}

func (e *distributedMastershipElection) Watch(ch chan<- State) error {
	e.mu.Lock()
	e.watchers = append(e.watchers, ch)
	e.mu.Unlock()
	return nil
}

func (e *distributedMastershipElection) Close() error {
	return e.election.Close(context.Background())
}

var _ Election = &distributedMastershipElection{}
