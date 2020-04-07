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
	"fmt"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sync"
)

// newLogicalClock creates a new logical clock for the given mastership Election
func newLogicalClock(election mastership.Election) LogicalClock {
	return &logicalClock{
		election: election,
	}
}

// LogicalTimestamp is a logical timestamp
type LogicalTimestamp struct {
	// Epoch is the epoch
	Epoch Epoch
	// LogicalTime is the logical time
	LogicalTime LogicalTime
}

// Epoch is a monotonically increasing, globally unique system timestamp backed by mastership election
type Epoch uint64

// LogicalTime is a monotonically increasing timestamp that resets for each unique Epoch
type LogicalTime uint64

// LogicalClock is a mastership based logical clock
type LogicalClock interface {
	// GetTimestamp gets the next logical timestamp
	GetTimestamp() (LogicalTimestamp, error)
}

// logicalClock is the default implementation of a LogicalClock
type logicalClock struct {
	election mastership.Election
	epoch    Epoch
	time     LogicalTime
	mu       sync.Mutex
}

func (c *logicalClock) GetTimestamp() (LogicalTimestamp, error) {
	mastership, err := c.election.GetState()
	if err != nil {
		return LogicalTimestamp{}, err
	}
	if mastership.Master != c.election.NodeID() {
		return LogicalTimestamp{}, fmt.Errorf("node %s is not the master", c.election.NodeID())
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	epoch := Epoch(mastership.Term)
	if epoch > c.epoch {
		c.epoch = epoch
		c.time = 0
	}
	c.time++
	return LogicalTimestamp{
		Epoch:       epoch,
		LogicalTime: c.time,
	}, nil
}

var _ LogicalClock = &logicalClock{}
