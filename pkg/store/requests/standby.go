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
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"google.golang.org/grpc"
)

func newStandbyStore(deviceKey device.Key, cluster cluster.Cluster, state *deviceStoreState, log Log, config config.RequestsStoreConfig) (storeHandler, error) {
	handler := &standbyStore{
		config:    config,
		deviceKey: deviceKey,
		state:     state,
		cluster:   cluster,
	}
	if err := handler.open(); err != nil {
		return nil, err
	}
	return handler, nil
}

type standbyStore struct {
	config    config.RequestsStoreConfig
	deviceKey device.Key
	cluster   cluster.Cluster
	state     *deviceStoreState
	conn      *grpc.ClientConn
	client    requests.RequestsServiceClient
}

func (s *standbyStore) open() error {
	master := s.cluster.Replica(cluster.ReplicaID(s.state.getMastership().Master))
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

func (s *standbyStore) Append(request *Request, opts ...AppendOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.Append(ctx, &requests.AppendRequest{
		DeviceID: string(s.deviceKey),
		Request:  request.RicControlRequest,
	})
	return err
}

func (s *standbyStore) Ack(request *Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.client.Ack(ctx, &requests.AckRequest{
		DeviceID: string(s.deviceKey),
		Index:    uint64(request.Index),
	})
	return err
}

func (s *standbyStore) Watch(deviceID device.ID, ch chan<- Event, opts ...WatchOption) error {
	panic("not implemented")
}

func (s *standbyStore) append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	return nil, errors.New("not the master")
}

func (s *standbyStore) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	return nil, errors.New("not the master")
}

func (s *standbyStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	return nil, errors.New("not a replica")
}

func (s *standbyStore) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

var _ storeHandler = &standbyStore{}
