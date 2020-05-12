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

package dispatcher

import (
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger("dispatcher")

// Request is a dispatcher request
type Request struct {
	// ID is the request identifier
	ID string
}

// Watcher is implemented by controllers to implement watching for specific events
// Type identifiers that are written to the watcher channel will eventually be processed by
// the dispatcher.
type Watcher interface {
	// Start starts watching for events
	Start(ch chan<- Request) error

	// Stop stops watching for events
	Stop()
}

// Handler reconciles objects
// The reconciler will be called for each type ID received from the Watcher. The reconciler may indicate
// whether to retry requests by returning either false or a non-nil error. Reconcilers should make no
// assumptions regarding the ordering of requests and should use the provided type IDs to resolve types
// against the current state of the cluster.
type Handler interface {
	// Reconcile is called to reconcile the state of an object
	Reconcile(Request) (Result, error)
}

// Result is a reconciler result
type Result struct {
	// Requeue is the identifier of an event to requeue
	Requeue bool
}

// NewController creates a new dispatcher
func NewController(name string) *Dispatcher {
	return &Dispatcher{
		name:        name,
		activator:   &UnconditionalActivator{},
		partitioner: &UnaryPartitioner{},
		watchers:    make([]Watcher, 0),
		partitions:  make(map[PartitionKey]chan Request),
	}
}

// Dispatcher is a control loop
// The Dispatcher is responsible for processing events provided by a Watcher. Events are processed by
// a configurable Handler. The dispatcher processes events in a loop, retrying requests until the
// Handler can successfully process them.
// The Dispatcher can be activated or deactivated by a configurable Activator. When inactive, the dispatcher
// will ignore requests, and when active it processes all requests.
// For per-request filtering, a Filter can be provided which provides a simple bool to indicate whether a
// request should be passed to the Handler.
// Once the Handler receives a request, it should process the request using the current state of the cluster
// Reconcilers should not cache state themselves and should instead rely on stores for consistency.
// If a Handler returns false, the request will be requeued to be retried after all pending requests.
// If a Handler returns an error, the request will be retried after a backoff period.
// Once a Handler successfully processes a request by returning true, the request will be discarded.
// Requests can be partitioned among concurrent goroutines by configuring a WorkPartitioner. The dispatcher
// will create a goroutine per PartitionKey provided by the WorkPartitioner, and requests to different
// partitions may be handled concurrently.
type Dispatcher struct {
	name        string
	mu          sync.RWMutex
	activator   Activator
	partitioner WorkPartitioner
	filter      Filter
	watchers    []Watcher
	reconciler  Handler
	partitions  map[PartitionKey]chan Request
}

// Activate sets an activator for the dispatcher
func (c *Dispatcher) Activate(activator Activator) *Dispatcher {
	c.mu.Lock()
	c.activator = activator
	c.mu.Unlock()
	return c
}

// Partition partitions work among multiple goroutines for the dispatcher
func (c *Dispatcher) Partition(partitioner WorkPartitioner) *Dispatcher {
	c.mu.Lock()
	c.partitioner = partitioner
	c.mu.Unlock()
	return c
}

// Filter sets a filter for the dispatcher
func (c *Dispatcher) Filter(filter Filter) *Dispatcher {
	c.mu.Lock()
	c.filter = filter
	c.mu.Unlock()
	return c
}

// Watch adds a watcher to the dispatcher
func (c *Dispatcher) Watch(watcher Watcher) *Dispatcher {
	c.mu.Lock()
	c.watchers = append(c.watchers, watcher)
	c.mu.Unlock()
	return c
}

// Reconcile sets the reconciler for the dispatcher
func (c *Dispatcher) Reconcile(reconciler Handler) *Dispatcher {
	c.mu.Lock()
	c.reconciler = reconciler
	c.mu.Unlock()
	return c
}

// Start starts the request dispatcher
func (c *Dispatcher) Start() error {
	ch := make(chan bool)
	if err := c.activator.Start(ch); err != nil {
		return err
	}
	go func() {
		active := false
		for activate := range ch {
			if activate {
				if !active {
					log.Infof("Activating dispatcher %s", c.name)
					c.activate()
					active = true
				}
			} else {
				if active {
					log.Infof("Deactivating dispatcher %s", c.name)
					c.deactivate()
					active = false
				}
			}
		}
	}()
	return nil
}

// Stop stops the dispatcher
func (c *Dispatcher) Stop() {
	c.activator.Stop()
}

// activate activates the dispatcher
func (c *Dispatcher) activate() {
	ch := make(chan Request)
	wg := &sync.WaitGroup{}
	for _, watcher := range c.watchers {
		if err := c.startWatcher(ch, wg, watcher); err == nil {
			wg.Add(1)
		}
	}
	go func() {
		wg.Wait()
		close(ch)
	}()
	go c.processEvents(ch)
}

// startWatcher starts a single watcher
func (c *Dispatcher) startWatcher(ch chan Request, wg *sync.WaitGroup, watcher Watcher) error {
	watcherCh := make(chan Request)
	if err := watcher.Start(watcherCh); err != nil {
		return err
	}

	go func() {
		for request := range watcherCh {
			ch <- request
		}
		wg.Done()
	}()
	return nil
}

// deactivate deactivates the dispatcher
func (c *Dispatcher) deactivate() {
	for _, watcher := range c.watchers {
		watcher.Stop()
	}
}

// processEvents processes the events from the given channel
func (c *Dispatcher) processEvents(ch chan Request) {
	c.mu.RLock()
	filter := c.filter
	partitioner := c.partitioner
	c.mu.RUnlock()
	for request := range ch {
		// Ensure the event is applicable to this dispatcher
		if filter == nil || filter.Accept(request) {
			c.partition(request, partitioner)
		}
	}
}

// partition writes the given request to a partition
func (c *Dispatcher) partition(request Request, partitioner WorkPartitioner) {
	iteration := 1
	for {
		// Get the partition key for the object ID
		key, err := partitioner.Partition(request)
		if err != nil {
			time.Sleep(time.Duration(iteration*2) * time.Millisecond)
		} else {
			// Get or create a partition channel for the partition key
			c.mu.RLock()
			partition, ok := c.partitions[key]
			c.mu.RUnlock()
			if !ok {
				c.mu.Lock()
				partition, ok = c.partitions[key]
				if !ok {
					partition = make(chan Request)
					c.partitions[key] = partition
					go c.processRequests(partition)
				}
				c.mu.Unlock()
			}
			partition <- request
			return
		}
		iteration++
	}
}

// processRequests processes requests from the given channel
func (c *Dispatcher) processRequests(ch chan Request) {
	c.mu.RLock()
	reconciler := c.reconciler
	c.mu.RUnlock()

	for request := range ch {
		// Reconcile the request. If the reconciliation is not successful, requeue the request to be processed
		// after the remaining enqueued events.
		result := c.reconcile(request, reconciler)
		if result.Requeue {
			go c.requeueRequest(ch, request)
		}
	}
}

// requeueRequest requeues the given request
func (c *Dispatcher) requeueRequest(ch chan Request, id Request) {
	ch <- id
}

// reconcile reconciles the given request ID until complete
func (c *Dispatcher) reconcile(request Request, reconciler Handler) Result {
	iteration := 1
	for {
		// Reconcile the request. If an error occurs, use exponential backoff to retry in order.
		// Otherwise, return the result.
		result, err := reconciler.Reconcile(request)
		if err != nil {
			log.Errorf("An error occurred during reconciliation of %s: %v", request, err)
			time.Sleep(time.Duration(iteration*2) * time.Millisecond)
		} else {
			return result
		}
		iteration++
	}
}
