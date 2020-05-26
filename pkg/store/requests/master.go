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
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"github.com/onosproject/onos-ric/pkg/store/mastership"
	"sort"
	"sync"
)

func newMasterStore(deviceKey device.Key, c cluster.Cluster, mastership mastership.State, log Log) (storeHandler, error) {
	handler := &masterStore{
		deviceKey:      deviceKey,
		cluster:        c,
		mastership:     mastership,
		backups:        make([]*masterBackup, 0),
		log:            log,
		commitCh:       make(chan backupCommit),
		commitChannels: make(map[Index]chan<- Index),
		commitIndexes:  make(map[cluster.ReplicaID]Index),
	}
	if err := handler.open(); err != nil {
		return nil, err
	}
	return handler, nil
}

type masterStore struct {
	deviceKey      device.Key
	cluster        cluster.Cluster
	mastership     mastership.State
	backups        []*masterBackup
	log            Log
	commitCh       chan backupCommit
	commitChannels map[Index]chan<- Index
	commitIndexes  map[cluster.ReplicaID]Index
	commitIndex    Index
	ackIndex       Index
	mu             sync.RWMutex
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
					backup, err := newMasterBackup(s, replicaID, client, s.log.OpenReader(0))
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

	go s.processCommits()
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
	s.mu.Lock()
	entry := s.log.Writer().Write(request.Request)
	s.mu.Unlock()

	ch := make(chan Index)
	go s.backupEntry(entry, ch)
	select {
	case <-ch:
		return &requests.AppendResponse{
			Index: uint64(entry.Index),
		}, nil
	case <-ctx.Done():
		return nil, errors.New("commit timed out")
	}
}

func (s *masterStore) ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	s.mu.Lock()
	s.ackIndex = Index(request.Index)
	s.mu.Unlock()
	return &requests.AckResponse{}, nil
}

func (s *masterStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	return &requests.BackupResponse{
		DeviceID: string(s.deviceKey),
		Term:     uint64(s.mastership.Term),
	}, nil
}

func (s *masterStore) processCommits() {
	for commit := range s.commitCh {
		s.processCommit(commit.replica, commit.index)
	}
}

func (s *masterStore) processCommit(replicaID cluster.ReplicaID, index Index) {
	prevIndex := s.commitIndexes[replicaID]
	if index > prevIndex {
		s.commitIndexes[replicaID] = index

		indexes := make([]Index, len(s.backups))
		i := 0
		for _, index := range s.commitIndexes {
			indexes[i] = index
			i++
		}
		sort.Slice(indexes, func(i, j int) bool {
			return indexes[i] < indexes[j]
		})

		commitIndex := indexes[syncBackupCount-1]
		s.mu.Lock()
		defer s.mu.Unlock()
		if commitIndex > s.commitIndex {
			for i := s.commitIndex + 1; i <= commitIndex; i++ {
				s.commitIndex = index
				ch, ok := s.commitChannels[index]
				if ok {
					ch <- index
					delete(s.commitChannels, index)
				}
			}
			logger.Debugf("Committed entries up to %d", commitIndex)
		}
	}
}

func (s *masterStore) backupEntry(entry *Entry, ch chan<- Index) {
	if len(s.backups) == 0 {
		s.mu.Lock()
		if entry.Index > s.commitIndex {
			s.commitIndex = entry.Index
		}
		index := s.commitIndex
		s.mu.Unlock()
		ch <- index
		return
	}

	// Acquire a write lock and add the channel to commitChannels
	s.mu.Lock()
	s.commitChannels[entry.Index] = ch
	s.mu.Unlock()
}

func (s *masterStore) Close() error {
	s.mu.RLock()
	for _, backup := range s.backups {
		backup.close()
	}
	close(s.commitCh)
	return nil
}

type backupCommit struct {
	replica cluster.ReplicaID
	index   Index
}

var _ storeHandler = &masterStore{}

func newMasterBackup(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader) (*masterBackup, error) {
	backup := &masterBackup{
		store:     store,
		replicaID: replicaID,
		client:    client,
		reader:    reader,
	}
	if err := backup.open(); err != nil {
		return nil, err
	}
	return backup, nil
}

type masterBackup struct {
	store     *masterStore
	replicaID cluster.ReplicaID
	client    requests.RequestsServiceClient
	reader    Reader
	closed    bool
}

func (b *masterBackup) open() error {
	go b.backupEntries()
	return nil
}

func (b *masterBackup) backupEntries() {
	for !b.closed {
		batch := b.reader.Read()
		b.store.mu.RLock()
		entries := make([]*requests.BackupEntry, len(batch.Entries))
		for i := 0; i < len(batch.Entries); i++ {
			entry := batch.Entries[i]
			entries[i] = &requests.BackupEntry{
				Index:   uint64(entry.Index),
				Request: entry.Value.(*e2ap.RicControlRequest),
			}
		}
		request := &requests.BackupRequest{
			DeviceID:    string(b.store.deviceKey),
			Term:        uint64(b.store.mastership.Term),
			Entries:     entries,
			PrevIndex:   uint64(batch.PrevIndex),
			CommitIndex: uint64(b.store.commitIndex),
			AckIndex:    uint64(b.store.ackIndex),
		}
		b.store.mu.RUnlock()

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		response, err := b.client.Backup(ctx, request)
		cancel()
		if err != nil {
			b.reader.Seek(batch.PrevIndex + 1)
		} else if Index(response.Index) != batch.Entries[len(batch.Entries)-1].Index {
			b.reader.Seek(Index(response.Index) + 1)
		} else {
			b.store.commitCh <- backupCommit{
				replica: b.replicaID,
				index:   Index(response.Index),
			}
		}
	}
}

func (b *masterBackup) close() error {
	b.closed = true
	return nil
}