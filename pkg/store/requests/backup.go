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

package requests

import (
	"context"
	"errors"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"google.golang.org/grpc"
	"sync"
)

func newBackupStore(deviceKey device.Key, cluster cluster.Cluster, mastership mastership.State, log Log) (storeHandler, error) {
	handler := &backupStore{
		deviceKey:  deviceKey,
		mastership: mastership,
		cluster:    cluster,
		log:        log,
	}
	if err := handler.open(); err != nil {
		return nil, err
	}
	return handler, nil
}

type backupStore struct {
	deviceKey  device.Key
	cluster    cluster.Cluster
	mastership mastership.State
	log        Log
	conn       *grpc.ClientConn
	client     requests.RequestsServiceClient
	mu         sync.RWMutex
}

func (s *backupStore) open() error {
	master := s.cluster.Replica(cluster.ReplicaID(s.mastership.Master))
	if master == nil {
		return errors.New("unknown master node")
	}
	conn, err := master.Connect()
	if err != nil {
		return err
	}
	s.client = requests.NewRequestsServiceClient(conn)
	return nil
}

func (s *backupStore) Append(request *Request, opts ...AppendOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.Append(ctx, &requests.AppendRequest{
		DeviceID: string(s.deviceKey),
		Request:  request.RicControlRequest,
	})
	return err
}

func (s *backupStore) Ack(request *Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.Ack(ctx, &requests.AckRequest{
		DeviceID: string(s.deviceKey),
		Index:    uint64(request.Index),
	})
	return err
}

func (s *backupStore) Watch(deviceID device.ID, ch chan<- Event, opts ...WatchOption) error {
	options := &watchOptions{}
	for _, opt := range opts {
		opt.applyWatch(options)
	}

	var reader Reader
	nextIndex := s.log.Writer().Index() + 1
	if options.replay {
		reader = s.log.OpenReader(0)
	} else {
		reader = s.log.OpenReader(nextIndex)
	}

	go func() {
		for {
			batch := reader.Read()
			for _, entry := range batch.Entries {
				var eventType EventType
				if entry.Index < nextIndex {
					eventType = EventNone
				} else {
					eventType = EventAppend
				}
				request := New(entry.Value.(*e2ap.RicControlRequest))
				request.Index = entry.Index
				ch <- Event{
					Type:    eventType,
					Request: *request,
				}
			}
		}
	}()
	return nil
}

func (s *backupStore) append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	return nil, errors.New("not the master")
}

func (s *backupStore) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	return nil, errors.New("not the master")
}

func (s *backupStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	prevIndex := Index(request.PrevIndex)
	lastIndex := s.log.Writer().Index()
	if prevIndex != 0 && lastIndex != prevIndex {
		return &requests.BackupResponse{
			DeviceID: string(s.deviceKey),
			Index:    uint64(lastIndex),
			Term:     uint64(s.mastership.Term),
		}, nil
	}

	for _, entry := range request.Entries {
		s.log.Writer().Write(entry.Request)
	}
	return &requests.BackupResponse{
		DeviceID: string(s.deviceKey),
		Index:    uint64(s.log.Writer().Index()),
		Term:     uint64(s.mastership.Term),
	}, nil
}

func (s *backupStore) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

var _ storeHandler = &backupStore{}
