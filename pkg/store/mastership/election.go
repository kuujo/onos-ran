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
	"github.com/onosproject/onos-ric/pkg/store/cluster"
	"io"
	"sync"
	"time"
)

// newDistributedElection returns a new distributed device mastership election
func newDistributedElection(partitionID PartitionID, database *client.Database) (mastershipElection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	election, err := database.GetElection(ctx, fmt.Sprintf("mastership-%d", partitionID), election.WithID(string(cluster.GetNodeID())))
	cancel()
	if err != nil {
		return nil, err
	}
	return newMastershipElection(partitionID, election)
}

// newLocalElection returns a new local device mastership election
func newLocalElection(partitionID PartitionID, nodeID cluster.NodeID, address net.Address) (mastershipElection, error) {
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

// newDeviceMastershipElection creates and enters a new device mastership election
func newMastershipElection(partitionID PartitionID, election election.Election) (mastershipElection, error) {
	mastershipElection := &distributedMastershipElection{
		partitionID: partitionID,
		election:    election,
		watchers:    make([]chan<- MastershipState, 0, 1),
	}
	if err := mastershipElection.enter(); err != nil {
		return nil, err
	}
	return mastershipElection, nil
}

// mastershipElection is an election for a single device mastership
type mastershipElection interface {
	io.Closer

	// NodeID returns the local node identifier used in the election
	NodeID() cluster.NodeID

	// PartitionID returns the mastership election partition identifier
	PartitionID() PartitionID

	// getState returns the mastership state
	getState() (*MastershipState, error)

	// isMaster returns a bool indicating whether the local node is the master for the device
	isMaster() (bool, error)

	// watch watches the election for changes
	watch(ch chan<- MastershipState) error
}

// distributedMastershipElection is a persistent device mastership election
type distributedMastershipElection struct {
	partitionID PartitionID
	election    election.Election
	mastership  *MastershipState
	watchers    []chan<- MastershipState
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
	e.mastership = &MastershipState{
		PartitionID: e.partitionID,
		Master:      cluster.NodeID(term.Leader),
		Term:        Term(term.ID),
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
		var mastership *MastershipState
		e.mu.Lock()
		if uint64(e.mastership.Term) != event.Term.ID {
			mastership = &MastershipState{
				PartitionID: e.partitionID,
				Term:        Term(event.Term.ID),
				Master:      cluster.NodeID(event.Term.Leader),
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

func (e *distributedMastershipElection) getState() (*MastershipState, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.mastership, nil
}

func (e *distributedMastershipElection) isMaster() (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.mastership == nil || string(e.mastership.Master) != e.election.ID() {
		return false, nil
	}
	return true, nil
}

func (e *distributedMastershipElection) watch(ch chan<- MastershipState) error {
	e.mu.Lock()
	e.watchers = append(e.watchers, ch)
	e.mu.Unlock()
	return nil
}

func (e *distributedMastershipElection) Close() error {
	return e.election.Close(context.Background())
}

var _ mastershipElection = &distributedMastershipElection{}
