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
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"sync"
)

func newMasterStore(deviceKey device.Key, c cluster.Cluster, state *deviceStoreState, log Log, config config.RequestsStoreConfig) (storeHandler, error) {
	handler := &masterStore{
		config:         config,
		deviceKey:      deviceKey,
		cluster:        c,
		backups:        make([]*backupSynchronizer, 0),
		state:          state,
		log:            log,
		commitCh:       make(chan commit),
		commitChannels: make(map[Index]chan<- Index),
		commitIndexes:  make(map[cluster.ReplicaID]Index),
	}
	if err := handler.open(); err != nil {
		return nil, err
	}
	return handler, nil
}

type masterStore struct {
	config         config.RequestsStoreConfig
	deviceKey      device.Key
	cluster        cluster.Cluster
	backups        []*backupSynchronizer
	state          *deviceStoreState
	log            Log
	commitCh       chan commit
	commitChannels map[Index]chan<- Index
	commitIndexes  map[cluster.ReplicaID]Index
	mu             sync.RWMutex
}

func (s *masterStore) open() error {
	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	mastership := s.state.getMastership()
	for index, replicaID := range mastership.Replicas {
		if index > 0 {
			sync := index < s.config.GetSyncBackups()
			wg.Add(1)
			go func(replicaID cluster.ReplicaID) {
				var err error
				if sync {
					err = s.addBackup(replicaID, newSyncBackup)
				} else {
					err = s.addBackup(replicaID, newAsyncBackup)
				}
				if err != nil {
					errCh <- err
				}
				wg.Done()
			}(replicaID)
		}
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return err
	}

	go s.processCommits()
	return nil
}

func (s *masterStore) addBackup(replicaID cluster.ReplicaID, factory func(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader) (*backupSynchronizer, error)) error {
	replica := s.cluster.Replica(replicaID)
	if replica == nil {
		return fmt.Errorf("unknown replica node %s", replicaID)
	}
	conn, err := replica.Connect()
	if err != nil {
		return err
	}
	client := requests.NewRequestsServiceClient(conn)
	synchronizer, err := factory(s, replicaID, client, s.log.OpenReader(0))
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.backups = append(s.backups, synchronizer)
	s.mu.Unlock()
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
	panic("not implemented")
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
	index := Index(request.Index)
	s.log.Writer().Discard(index)
	s.state.setAckIndex(index)
	s.mu.Unlock()
	return &requests.AckResponse{}, nil
}

func (s *masterStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	return &requests.BackupResponse{
		DeviceID: string(s.deviceKey),
		Term:     uint64(s.state.getMastership().Term),
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

		var minIndex Index
		for _, index := range s.commitIndexes {
			if minIndex == 0 || index < minIndex {
				minIndex = index
			}
		}

		s.mu.Lock()
		defer s.mu.Unlock()
		commitIndex := s.state.getCommitIndex()
		if minIndex > commitIndex {
			for index := commitIndex + 1; index <= minIndex; index++ {
				s.state.setCommitIndex(index)
				ch, ok := s.commitChannels[index]
				if ok {
					ch <- index
					delete(s.commitChannels, index)
				}
			}
			logger.Debugf("Committed entries up to %d", minIndex)
		}
	}
}

func (s *masterStore) backupEntry(entry *Entry, ch chan<- Index) {
	if len(s.backups) == 0 {
		s.mu.Lock()
		commitIndex := s.state.getCommitIndex()
		if entry.Index > commitIndex {
			s.state.setCommitIndex(entry.Index)
		}
		index := s.state.getCommitIndex()
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

type commit struct {
	replica cluster.ReplicaID
	index   Index
}

type indexFunc func() Index

type commitFunc func(Index)

var _ storeHandler = &masterStore{}

func newBackupSynchronizer(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader, limitFunc indexFunc, committer commitFunc) (*backupSynchronizer, error) {
	synchronizer := &backupSynchronizer{
		store:      store,
		replicaID:  replicaID,
		client:     client,
		reader:     reader,
		limitFunc:  limitFunc,
		commitFunc: committer,
	}
	if err := synchronizer.open(); err != nil {
		return nil, err
	}
	return synchronizer, nil
}

type backupSynchronizer struct {
	store      *masterStore
	replicaID  cluster.ReplicaID
	client     requests.RequestsServiceClient
	reader     Reader
	closed     bool
	limitFunc  indexFunc
	commitFunc commitFunc
	mu         sync.RWMutex
}

func (b *backupSynchronizer) open() error {
	go b.backupEntries()
	return nil
}

func (b *backupSynchronizer) backupEntries() {
	b.mu.RLock()
	closed := b.closed
	b.mu.RUnlock()
	for !closed {
		var batch *Batch
		if b.limitFunc != nil {
			batch = b.reader.ReadUntil(b.limitFunc())
		} else {
			batch = b.reader.Read()
		}

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
			Term:        uint64(b.store.state.getMastership().Term),
			Entries:     entries,
			PrevIndex:   uint64(batch.PrevIndex),
			CommitIndex: uint64(b.store.state.getCommitIndex()),
			AckIndex:    uint64(b.store.state.getAckIndex()),
		}
		b.store.mu.RUnlock()

		logger.Debugf("Sending BackupRequest %s", request)
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		response, err := b.client.Backup(ctx, request)
		cancel()
		logger.Debugf("Received BackupResponse %s", response)
		if err != nil {
			b.reader.Seek(batch.PrevIndex + 1)
		} else if Index(response.Index) != batch.Entries[len(batch.Entries)-1].Index {
			b.reader.Seek(Index(response.Index) + 1)
		} else if b.commitFunc != nil {
			b.commitFunc(Index(response.Index))
		}
		b.mu.RLock()
		closed = b.closed
		b.mu.RUnlock()
	}
}

func (b *backupSynchronizer) close() error {
	b.mu.Lock()
	b.closed = true
	b.mu.Unlock()
	return nil
}

func newSyncBackup(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader) (*backupSynchronizer, error) {
	committer := func(index Index) {
		store.commitCh <- commit{
			replica: replicaID,
			index:   index,
		}
	}
	return newBackupSynchronizer(store, replicaID, client, reader, nil, committer)
}

func newAsyncBackup(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader) (*backupSynchronizer, error) {
	limiter := func() Index {
		return store.state.getCommitIndex()
	}
	return newBackupSynchronizer(store, replicaID, client, reader, limiter, nil)
}
