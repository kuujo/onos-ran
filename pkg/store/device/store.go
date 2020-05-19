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

package device

import (
	"context"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/onosproject/onos-ric/api/sb"
	"io"
	"strings"
	"time"

	topodevice "github.com/onosproject/onos-topo/api/device"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("store", "device")

const e2NodeType = topodevice.Type("E2Node")

// ID is a topology service device identifier
type ID sb.ECGI

// newID creates a new device identifier
func newID(id topodevice.ID) (ID, error) {
	if !strings.Contains(string(id), "-") {
		return ID{}, fmt.Errorf("unexpected format for E2Node ID %s", id)
	}
	parts := strings.Split(string(id), "-")
	if len(parts) != 2 {
		return ID{}, fmt.Errorf("unexpected format for E2Node ID %s", id)
	}
	return ID(sb.ECGI{Ecid: parts[1], PlmnId: parts[0]}), nil
}

// getDeviceID returns the topo device identifier for the given device ID
func getDeviceID(id ID) topodevice.ID {
	return topodevice.ID(fmt.Sprintf("%s-%s", id.PlmnId, id.Ecid))
}

// newDevice creates a new device
func newDevice(device *topodevice.Device) (*Device, error) {
	id, err := newID(device.ID)
	if err != nil {
		return nil, err
	}
	return &Device{
		ID:     id,
		Device: device,
	}, nil
}

// Device is a topology service device
type Device struct {
	ID ID
	*topodevice.Device
}

// EventType is a device event type
type EventType string

const (
	// EventNone is an event indicating the replay of an existing device
	EventNone EventType = ""
	// EventAdded indicates a newly added device
	EventAdded EventType = "added"
	// EventUpdated indicates an updated device
	EventUpdated EventType = "updated"
	// EventRemoved indicates a removed device
	EventRemoved EventType = "removed"
)

// Event is a device event
type Event struct {
	// Type is the event type
	Type EventType

	// Device is the updated device
	Device Device
}

// Store is a device store
type Store interface {
	// Get gets a device by ID
	Get(ID) (*Device, error)

	// Update updates a given device
	Update(*Device) (*Device, error)

	// List lists the devices in the store
	List(chan<- Device) error

	// Watch watches the device store for changes
	Watch(chan<- Event) error
}

// NewTopoStore returns a new topo-based device store
func NewTopoStore(topoEndpoint string, opts ...grpc.DialOption) (Store, error) {
	if len(opts) == 0 {
		return nil, fmt.Errorf("no opts given when creating topo store")
	}
	opts = append(opts, grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100*time.Millisecond)))
	conn, err := getTopoConn(topoEndpoint, opts...)

	if err != nil {
		return nil, err
	}
	client := topodevice.NewDeviceServiceClient(conn)

	return &topoStore{
		client: client,
	}, nil
}

// NewStore returns a new device store for the given client
func NewStore(client topodevice.DeviceServiceClient) (Store, error) {
	return &topoStore{client: client}, nil
}

// A device Store that uses the topo service to propagate devices
type topoStore struct {
	client topodevice.DeviceServiceClient
}

func (s *topoStore) Get(id ID) (*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := s.client.Get(ctx, &topodevice.GetRequest{
		ID: getDeviceID(id),
	})
	if err != nil {
		return nil, err
	} else if !isE2Device(response.Device) {
		return nil, nil
	}
	return newDevice(response.Device)
}

func (s *topoStore) Update(device *Device) (*Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	request := &topodevice.UpdateRequest{
		Device: device.Device,
	}
	response, err := s.client.Update(ctx, request)
	if err != nil {
		return nil, err
	} else if !isE2Device(response.Device) {
		return nil, nil
	}
	return newDevice(response.Device)
}

func (s *topoStore) List(ch chan<- Device) error {
	list, err := s.client.List(context.Background(), &topodevice.ListRequest{})
	if err != nil {
		return err
	}

	go func() {
		for {
			response, err := list.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			if isE2Device(response.Device) {
				device, err := newDevice(response.Device)
				if err != nil {
					log.Error(err)
				} else {
					ch <- *device
				}
			}
		}
	}()
	return nil
}

func (s *topoStore) Watch(ch chan<- Event) error {
	list, err := s.client.List(context.Background(), &topodevice.ListRequest{
		Subscribe: true,
	})
	if err != nil {
		return err
	}
	go func() {
		for {
			response, err := list.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			if isE2Device(response.Device) {
				var eventType EventType
				switch response.Type {
				case topodevice.ListResponse_NONE:
					eventType = EventNone
				case topodevice.ListResponse_ADDED:
					eventType = EventAdded
				case topodevice.ListResponse_UPDATED:
					eventType = EventUpdated
				case topodevice.ListResponse_REMOVED:
					eventType = EventRemoved
				}
				device, err := newDevice(response.Device)
				if err != nil {
					log.Error(err)
				} else {
					ch <- Event{
						Type:   eventType,
						Device: *device,
					}
				}
			}
		}
	}()
	return nil
}

func isE2Device(device *topodevice.Device) bool {
	return device.GetType() == e2NodeType
}

// getTopoConn gets a gRPC connection to the topology service
func getTopoConn(topoEndpoint string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(topoEndpoint, opts...)
}
