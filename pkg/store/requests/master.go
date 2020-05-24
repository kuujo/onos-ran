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
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"io"
	"sync"
)

func newMasterStore(deviceKey device.Key, cluster cluster.Cluster, mastership mastership.State, log Log) (storeHandler, error) {
	handler := &masterStore{
		deviceKey:  deviceKey,
		cluster:    cluster,
		mastership: mastership,
		log:        log,
		backups:    make([]*masterBackup, 0),
	}
	if err := handler.open(); err != nil {
		return nil, err
	}
	return handler, nil
}

type masterStore struct {
	deviceKey  device.Key
	cluster    cluster.Cluster
	mastership mastership.State
	log        Log
	backups    []*masterBackup
	mu         sync.RWMutex
}

func (s *masterStore) open() error {
	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	for _, replicaID := range s.mastership.Replicas {
		wg.Add(1)
		go func(replicaID cluster.ReplicaID) {
			replica := s.cluster.Replica(replicaID)
			if replica == nil {
				errCh <- fmt.Errorf("unknown replica node %s", replicaID)
			} else {
				conn, err := replica.Connect()
				if err != nil {
					errCh <- err
				} else {
					client := requests.NewRequestsServiceClient(conn)
					backup, err := newMasterBackup(s.deviceKey, replicaID, client, s.log.OpenReader(0))
					if err != nil {
						errCh <- err
					} else {
						s.mu.Lock()
						s.backups = append(s.backups, backup)
						s.mu.Unlock()
					}
				}
			}
			wg.Done()
		}(replicaID)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}
	return nil
}

func (s *masterStore) Append(request *Request, opts ...AppendOption) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.append(ctx, &requests.AppendRequest{
		DeviceID: string(s.deviceKey),
		Request:  request.RicControlRequest,
	})
	return err
}

func (s *masterStore) Ack(request *Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()
	_, err := s.ack(ctx, &requests.AckRequest{
		DeviceID: string(s.deviceKey),
		Index:    uint64(request.Index),
	})
	return err
}

func (s *masterStore) Watch(deviceID device.ID, ch chan<- Event, opts ...WatchOption) error {
	panic("implement me")
}

func (s *masterStore) append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	panic("implement me")
}

func (s *masterStore) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	panic("implement me")
}

func (s *masterStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	panic("implement me")
}

func (s *masterStore) Close() error {
	s.mu.RLock()
	for _, backup := range s.backups {
		backup.close()
	}
	return nil
}

func newMasterBackup(deviceKey device.Key, replicaID cluster.ReplicaID, mastership mastership.State, client requests.RequestsServiceClient, reader Reader) (*masterBackup, error) {
	backup := &masterBackup{
		deviceKey:  deviceKey,
		replicaID:  replicaID,
		mastership: mastership,
		client:     client,
		reader:     reader,
		backupCh:   make(chan Index),
	}
	if err := backup.open(); err != nil {
		return nil, err
	}
	return backup, nil
}

type masterBackup struct {
	deviceKey  device.Key
	replicaID  cluster.ReplicaID
	mastership mastership.State
	client     requests.RequestsServiceClient
	reader     Reader
	cancel     context.CancelFunc
	backupCh   chan Index
}

func (b *masterBackup) open() error {
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := b.client.Backup(ctx)
	if err != nil {
		cancel()
		return err
	}
	b.cancel = cancel
	go b.receiveStream(stream)
	go b.sendStream(stream)
	return nil
}

func (b *masterBackup) receiveStream(stream requests.RequestsService_BackupClient) {
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			return
		} else if err != nil {
			logger.Errorf("Failed receiving stream from %s: %s", b.replicaID, err)
		}
	}
}

func (b *masterBackup) sendStream(stream requests.RequestsService_BackupClient) {
	for {
		select {
		case <-b.backupCh:
			entry := b.reader.NextEntry()
			if entry != nil {
				err := stream.Send(&requests.BackupRequest{
					DeviceID: string(b.deviceKey),
					Index:    uint64(entry.Index),
					Term:     uint64(b.mastership.Term),
					Request:  entry.Request.RicControlRequest,
				})
				if err != nil {
					logger.Errorf("Failed to backup request %d to %s: %s", entry.Index, b.replicaID, err)
					b.reader.Reset(entry.Index)
				}
			}
		}
	}
}

func (b *masterBackup) close() error {
	b.cancel()
	return nil
}
