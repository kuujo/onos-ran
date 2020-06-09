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

package indications

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/store/indications"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"sync"
	"time"
)

// newDeviceIndicationsStore creates a new indications store for a single device
func newDeviceIndicationsStore(deviceKey device.Key, cluster cluster.Cluster, election mastership.Election) (*deviceIndicationsStore, error) {
	store := &deviceIndicationsStore{
		deviceKey:   deviceKey,
		cluster:     cluster,
		election:    election,
		cache:       make(map[sb.MessageType]*Indication),
		subscribers: make([]func(Event), 0),
		streams:     make(map[uuid.UUID]indications.IndicationsService_SubscribeServer),
		recordCh:    make(chan Indication, 1000),
	}
	err := backoff.Retry(store.open, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}
	return store, nil
}

// deviceIndicationsStore is a store of indications for a single device
type deviceIndicationsStore struct {
	deviceKey   device.Key
	cluster     cluster.Cluster
	election    mastership.Election
	cache       map[sb.MessageType]*Indication
	subscribers []func(Event)
	state       *mastership.State
	streams     map[uuid.UUID]indications.IndicationsService_SubscribeServer
	recordCh    chan Indication
	stateMu     sync.RWMutex
	cacheMu     sync.RWMutex
}

func (s *deviceIndicationsStore) open() error {
	ch := make(chan mastership.State)
	if err := s.election.Watch(ch); err != nil {
		return err
	}
	state, err := s.election.GetState()
	if err != nil {
		return err
	}
	s.state = state
	log.Debugf("Initializing store with mastership state %v for device %s", state, s.deviceKey)
	if state.Master != s.cluster.Node().ID {
		log.Debugf("Subscribing to events from %s for device %s", state.Master, s.deviceKey)
		ctx, cancel := context.WithCancel(context.Background())
		err := s.subscribeMaster(ctx, *state)
		if err != nil {
			cancel()
			return err
		}
		go s.watchMastership(ch, cancel)
	} else {
		go s.watchMastership(ch, nil)
	}
	go s.processRecords()
	return nil
}

func (s *deviceIndicationsStore) processRecords() {
	for indication := range s.recordCh {
		// Write the indication event to susbcriber channels
		s.stateMu.RLock()
		subscribers := s.subscribers
		state := s.state
		s.stateMu.RUnlock()
		event := Event{
			Type:       EventReceived,
			Indication: indication,
		}
		log.Debugf("Received event %v for device %s", event, s.deviceKey)
		for _, subscriber := range subscribers {
			subscriber(event)
		}

		// If this node is the master, send the indication to replicas
		if state != nil && state.Master == s.cluster.Node().ID {
			response := &indications.SubscribeResponse{
				Device:     string(s.deviceKey),
				Term:       uint64(state.Term),
				Indication: indication.RicIndication,
			}
			s.stateMu.RLock()
			for _, stream := range s.streams {
				log.Debugf("Sending SubscribeResponse %v for device %s", indication, s.deviceKey)
				err := stream.Send(response)
				if err != nil {
					log.Errorf("Failed to send indication: %s", err)
				}
			}
			s.stateMu.RUnlock()
		}

		// Update the local cache
		s.cacheMu.Lock()
		ref := &indication
		s.cache[indication.Hdr.MessageType] = ref
		s.cacheMu.Unlock()
	}
}

func (s *deviceIndicationsStore) watchMastership(ch <-chan mastership.State, cancel context.CancelFunc) {
	for state := range ch {
		log.Infof("Received mastership change %v for device %s", state, s.deviceKey)
		s.stateMu.Lock()
		ref := &state
		s.state = ref

		if state.Master != s.cluster.Node().ID && cancel == nil {
			log.Debugf("Subscribing to events from %s for device %s", state.Master, s.deviceKey)
			ctx, cancelFunc := context.WithCancel(context.Background())
			err := s.subscribeMaster(ctx, state)
			if err != nil {
				log.Errorf("Failed to subscribe to master node %s: %s", state.Master, err)
			} else {
				cancel = cancelFunc
			}
		} else if state.Master == s.cluster.Node().ID && cancel != nil {
			log.Debugf("Cancelling subscription for device %s", s.deviceKey)
			cancel()
			cancel = nil
		}
		s.stateMu.Unlock()
	}
}

func (s *deviceIndicationsStore) subscribeMaster(ctx context.Context, mastership mastership.State) error {
	master := s.cluster.Replica(cluster.ReplicaID(mastership.Master))
	if master == nil {
		return fmt.Errorf("cannot find master node %s", mastership.Master)
	}
	conn, err := master.Connect(grpc.WithInsecure(), grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(time.Second)))
	if err != nil {
		return err
	}
	client := indications.NewIndicationsServiceClient(conn)
	request := &indications.SubscribeRequest{
		Device: string(s.deviceKey),
		Term:   uint64(mastership.Term),
	}
	log.Debugf("Sending SubscribeRequest %v for device %s", request, s.deviceKey)
	stream, err := client.Subscribe(ctx, request)
	if err != nil {
		return err
	}

	// Start a goroutine to receive the stream
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Errorf("Received subscribe error from master node %s: %s", mastership.Master, err)
				break
			} else {
				log.Debugf("Received SubscribeResponse %v for device %s", response, s.deviceKey)
				s.recordCh <- *New(response.Indication)
			}
		}
	}()
	return nil
}

func (s *deviceIndicationsStore) Record(indication *Indication) error {
	log.Debugf("Recording indication %v for device %s", indication, s.deviceKey)
	s.stateMu.RLock()
	state := s.state
	s.stateMu.RUnlock()
	if state == nil || state.Master != s.cluster.Node().ID {
		return errors.New("not the master")
	}
	s.recordCh <- *indication
	return nil
}

func (s *deviceIndicationsStore) List(ch chan<- Indication) error {
	go func() {
		defer close(ch)
		s.cacheMu.RLock()
		defer s.cacheMu.RUnlock()
		for _, indication := range s.cache {
			ch <- *indication
		}
	}()
	return nil
}

func (s *deviceIndicationsStore) Subscribe(ch chan<- Event, opts ...SubscribeOption) error {
	options := &subscribeOptions{}
	for _, opt := range opts {
		opt.applySubscribe(options)
	}

	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	var f func(Event)
	if options.filter != nil {
		filter := options.filter
		f = func(event Event) {
			if filter(event) {
				ch <- event
			}
		}
	} else {
		f = func(event Event) {
			ch <- event
		}
	}
	s.subscribers = append(s.subscribers, f)

	if options.replay {
		go func() {
			s.cacheMu.RLock()
			defer s.cacheMu.RUnlock()
			for _, indication := range s.cache {
				event := Event{
					Type:       EventNone,
					Indication: *indication,
				}
				log.Debugf("Received event %v for device %s", event, s.deviceKey)
				ch <- event
			}
		}()
	}
	return nil
}

func (s *deviceIndicationsStore) subscribe(request *indications.SubscribeRequest, stream indications.IndicationsService_SubscribeServer) error {
	log.Debugf("Received SubscribeRequest %v for device %s", request, s.deviceKey)

	s.stateMu.RLock()
	state := s.state
	s.stateMu.RUnlock()

	if state == nil {
		return status.Error(codes.NotFound, "not the master")
	} else if state.Master != s.cluster.Node().ID || state.Term != mastership.Term(request.Term) {
		return status.Error(codes.NotFound, "not the master for the given term")
	}

	id := uuid.New()

	s.stateMu.Lock()
	s.streams[id] = stream
	s.stateMu.Unlock()

	<-stream.Context().Done()

	s.stateMu.Lock()
	delete(s.streams, id)
	s.stateMu.Unlock()

	return nil
}

func (s *deviceIndicationsStore) Close() error {
	return nil
}
