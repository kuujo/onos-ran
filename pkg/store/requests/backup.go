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
	"fmt"
	"github.com/atomix/go-client/pkg/client/util"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"google.golang.org/grpc"
	"time"
)

func newBackupStore(deviceKey device.Key, cluster cluster.Cluster, state *deviceStoreState, log Log, config config.RequestsStoreConfig) (storeHandler, error) {
	handler := &backupStore{
		config:    config,
		deviceKey: deviceKey,
		state:     state,
		cluster:   cluster,
		log:       log,
	}
	if err := handler.open(); err != nil {
		return nil, err
	}
	return handler, nil
}

type backupStore struct {
	config    config.RequestsStoreConfig
	deviceKey device.Key
	cluster   cluster.Cluster
	state     *deviceStoreState
	log       Log
	conn      *grpc.ClientConn
	client    requests.RequestsServiceClient
}

func (s *backupStore) open() error {
	mastership := s.state.getMastership()
	master := s.cluster.Replica(cluster.ReplicaID(mastership.Master))
	if master == nil {
		return fmt.Errorf("unknown master %s", mastership.Master)
	}
	conn, err := master.Connect(grpc.WithInsecure(), grpc.WithStreamInterceptor(util.RetryingStreamClientInterceptor(time.Second)))
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
	panic("not implemented")
}

func (s *backupStore) append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	return nil, errors.New("not the master")
}

func (s *backupStore) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	return nil, errors.New("not the master")
}

func (s *backupStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	s.state.mu.Lock()
	defer s.state.mu.Unlock()

	logger.Debugf("Received BackupRequest %s", request)

	// TODO: Backup should compare entry terms when appending to avoid inconsistencies
	prevIndex := Index(request.PrevIndex)
	lastIndex := s.log.Writer().Index()
	if prevIndex != 0 && lastIndex != prevIndex {
		return &requests.BackupResponse{
			DeviceID: string(s.deviceKey),
			Index:    uint64(lastIndex),
			Term:     uint64(s.state.getMastership().Term),
		}, nil
	}

	commitIndex := Index(request.CommitIndex)
	ackIndex := Index(request.AckIndex)
	for _, entry := range request.Entries {
		entry := s.log.Writer().Write(entry.Request)
		s.state.setAppendIndex(entry.Index)
	}

	// Discard entries up to the ackIndex
	s.log.Writer().Discard(ackIndex)

	// Update the commitIndex and ackIndex to trigger watch events
	if commitIndex > s.state.getCommitIndex() {
		s.state.setCommitIndex(commitIndex)
		logger.Debugf("Committed entries up to %d for device %s", commitIndex, s.deviceKey)
	}

	if ackIndex > s.state.getAckIndex() {
		s.state.setAckIndex(ackIndex)
		logger.Debugf("Acknowledged entries up to %d for device %s", ackIndex, s.deviceKey)
	}

	response := &requests.BackupResponse{
		DeviceID: string(s.deviceKey),
		Index:    uint64(s.log.Writer().Index()),
		Term:     uint64(s.state.getMastership().Term),
	}
	logger.Debugf("Sending BackupResponse %s", response)
	return response, nil
}

func (s *backupStore) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

var _ storeHandler = &backupStore{}
