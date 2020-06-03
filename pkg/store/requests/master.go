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
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/store/requests"
	"github.com/onosproject/onos-ric/pkg/config"
	"github.com/onosproject/onos-ric/pkg/store/device"
	"google.golang.org/grpc"
	"sync"
	"time"
)

func newMasterStore(deviceKey device.Key, c cluster.Cluster, state *deviceStoreState, log Log, config config.RequestsStoreConfig) (storeHandler, error) {
	handler := &masterStore{
		config:         config,
		deviceKey:      deviceKey,
		cluster:        c,
		backups:        make([]*backupSynchronizer, 0),
		state:          state,
		log:            log,
		commitCh:       make(chan commit, 1000),
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
}

func (s *masterStore) open() error {
	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	mastership := s.state.getMastership()
	for index, replicaID := range mastership.Replicas {
		if index > 0 {
			sync := index <= s.config.GetSyncBackups()
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
		return fmt.Errorf("unknown replica %s", replicaID)
	}
	conn, err := replica.Connect(grpc.WithInsecure(), grpc.WithStreamInterceptor(util.RetryingStreamClientInterceptor(time.Second)))
	if err != nil {
		return err
	}
	client := requests.NewRequestsServiceClient(conn)
	synchronizer, err := factory(s, replicaID, client, s.log.OpenReader(0))
	if err != nil {
		return err
	}
	s.state.mu.Lock()
	s.backups = append(s.backups, synchronizer)
	s.state.mu.Unlock()
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
	s.state.mu.Lock()
	entry := s.log.Writer().Write(request.Request)
	s.state.mu.Unlock()

	ch := make(chan Index, 1)
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
	index := Index(request.Index)
	s.state.mu.Lock()
	if index > s.state.getAckIndex() {
		s.log.Writer().Discard(index)
		s.state.setAckIndex(index)
		logger.Debugf("Acknowledged entries up to %d for device %s", index, s.deviceKey)
	}
	s.state.mu.Unlock()
	return &requests.AckResponse{}, nil
}

func (s *masterStore) backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	s.state.mu.RLock()
	defer s.state.mu.RUnlock()
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

		s.state.mu.Lock()
		defer s.state.mu.Unlock()
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
			logger.Debugf("Committed entries up to %d for device %s", minIndex, s.deviceKey)
		}
	}
}

func (s *masterStore) backupEntry(entry *Entry, ch chan<- Index) {
	if len(s.backups) == 0 || s.config.GetSyncBackups() == 0 {
		s.state.mu.Lock()
		commitIndex := s.state.getCommitIndex()
		if entry.Index > commitIndex {
			s.state.setCommitIndex(entry.Index)
			logger.Debugf("Committed entries up to %d for device %s", entry.Index, s.deviceKey)
		}
		index := s.state.getCommitIndex()
		s.state.mu.Unlock()
		ch <- index
		return
	}

	// Acquire a write lock and add the channel to commitChannels
	s.state.mu.Lock()
	if s.state.getCommitIndex() >= entry.Index {
		ch <- entry.Index
	} else {
		s.commitChannels[entry.Index] = ch
	}
	s.state.mu.Unlock()
}

func (s *masterStore) Close() error {
	s.state.mu.RLock()
	for _, backup := range s.backups {
		backup.close()
	}
	s.state.mu.RUnlock()
	close(s.commitCh)
	return nil
}

type commit struct {
	replica cluster.ReplicaID
	index   Index
}

var _ storeHandler = &masterStore{}

func newBackupSynchronizer(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader, sync bool) (*backupSynchronizer, error) {
	synchronizer := &backupSynchronizer{
		store:     store,
		replicaID: replicaID,
		client:    client,
		reader:    reader,
		sync:      sync,
		closeCh:   make(chan struct{}),
	}
	if err := synchronizer.open(); err != nil {
		return nil, err
	}
	return synchronizer, nil
}

type backupSynchronizer struct {
	store     *masterStore
	replicaID cluster.ReplicaID
	client    requests.RequestsServiceClient
	reader    Reader
	closeCh   chan struct{}
	sync      bool
}

func (b *backupSynchronizer) open() error {
	commitCh := make(chan Index)
	ctx, cancel := context.WithCancel(context.Background())
	b.store.state.mu.Lock()
	b.store.state.watchCommitIndex(ctx, commitCh)
	ackCh := make(chan Index)
	b.store.state.watchAckIndex(ctx, ackCh)
	b.store.state.mu.Unlock()
	go b.backupEntries(commitCh, ackCh)
	go func() {
		<-b.closeCh
		cancel()
	}()
	return nil
}

func (b *backupSynchronizer) backupEntries(commitCh <-chan Index, ackCh <-chan Index) {
	var backedUpIndex Index
	var committedIndex Index
	var ackedIndex Index
	for {
		if b.sync {
			select {
			case batch := <-b.reader.AwaitBatch():
				b.store.state.mu.RLock()
				commitIndex := b.store.state.getCommitIndex()
				ackIndex := b.store.state.getAckIndex()
				b.store.state.mu.RUnlock()
				err := b.backupBatch(batch, commitIndex, ackIndex)
				if err == nil {
					backedUpIndex = batch.Entries[len(batch.Entries)-1].Index
					committedIndex = commitIndex
					ackedIndex = ackIndex
				}
			case commitIndex := <-commitCh:
				if commitIndex > committedIndex {
					b.store.state.mu.RLock()
					ackIndex := b.store.state.getAckIndex()
					b.store.state.mu.RUnlock()
					err := b.backupBatch(&Batch{
						PrevIndex: backedUpIndex,
						Entries:   []*Entry{},
					}, commitIndex, ackIndex)
					if err == nil {
						committedIndex = commitIndex
						ackedIndex = ackIndex
					}
				}
			case ackIndex := <-ackCh:
				if ackIndex > ackedIndex {
					b.store.state.mu.RLock()
					commitIndex := b.store.state.getCommitIndex()
					b.store.state.mu.RUnlock()
					err := b.backupBatch(&Batch{
						PrevIndex: backedUpIndex,
						Entries:   []*Entry{},
					}, commitIndex, ackIndex)
					if err == nil {
						committedIndex = commitIndex
						ackedIndex = ackIndex
					}
				}
			case <-b.closeCh:
				return
			}
		} else {
			select {
			case commitIndex := <-commitCh:
				if commitIndex > committedIndex {
					b.store.state.mu.RLock()
					ackIndex := b.store.state.getAckIndex()
					b.store.state.mu.RUnlock()
					batch := b.reader.ReadUntil(commitIndex)
					err := b.backupBatch(batch, commitIndex, ackIndex)
					if err == nil {
						backedUpIndex = batch.Entries[len(batch.Entries)-1].Index
						committedIndex = commitIndex
						ackedIndex = ackIndex
					}
				}
			case ackIndex := <-ackCh:
				if ackIndex > ackedIndex {
					b.store.state.mu.RLock()
					commitIndex := b.store.state.getCommitIndex()
					b.store.state.mu.RUnlock()
					err := b.backupBatch(&Batch{
						PrevIndex: backedUpIndex,
						Entries:   []*Entry{},
					}, commitIndex, ackIndex)
					if err == nil {
						committedIndex = commitIndex
						ackedIndex = ackIndex
					}
				}
			case <-b.closeCh:
				return
			}
		}
	}
}

func (b *backupSynchronizer) backupBatch(batch *Batch, commitIndex Index, ackIndex Index) error {
	b.store.state.mu.RLock()
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
		CommitIndex: uint64(commitIndex),
		AckIndex:    uint64(ackIndex),
	}
	b.store.state.mu.RUnlock()

	logger.Debugf("Sending BackupRequest %s", request)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	response, err := b.client.Backup(ctx, request)
	cancel()
	logger.Debugf("Received BackupResponse %s", response)
	if err != nil {
		b.reader.Seek(batch.PrevIndex + 1)
		return err
	} else if Index(response.Index) != batch.Entries[len(batch.Entries)-1].Index {
		b.reader.Seek(Index(response.Index) + 1)
		return errors.New("uncommitted backup")
	} else if b.sync {
		b.store.commitCh <- commit{
			replica: b.replicaID,
			index:   Index(response.Index),
		}
	}
	return nil
}

func (b *backupSynchronizer) close() error {
	close(b.closeCh)
	return nil
}

func newSyncBackup(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader) (*backupSynchronizer, error) {
	return newBackupSynchronizer(store, replicaID, client, reader, true)
}

func newAsyncBackup(store *masterStore, replicaID cluster.ReplicaID, client requests.RequestsServiceClient, reader Reader) (*backupSynchronizer, error) {
	return newBackupSynchronizer(store, replicaID, client, reader, false)
}
