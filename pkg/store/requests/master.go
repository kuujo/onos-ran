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
	"sort"
	"sync"
)

func newMasterStore(deviceKey device.Key, c cluster.Cluster, state *deviceStoreState, log Log, config config.RequestsStoreConfig) (storeHandler, error) {
	handler := &masterStore{
		config:         config,
		deviceKey:      deviceKey,
		cluster:        c,
		backups:        make([]*masterBackup, 0),
		state:          state,
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
	config         config.RequestsStoreConfig
	deviceKey      device.Key
	cluster        cluster.Cluster
	backups        []*masterBackup
	state          *deviceStoreState
	log            Log
	commitCh       chan backupCommit
	commitChannels map[Index]chan<- Index
	commitIndexes  map[cluster.ReplicaID]Index
	mu             sync.RWMutex
}

func (s *masterStore) open() error {
	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	mastership := s.state.getMastership()
	for index, replicaID := range mastership.Replicas {
		if index > 0 && index < s.config.GetBackups()+s.config.GetAsyncBackups() {
			wg.Add(1)
			go func(replicaID cluster.ReplicaID) {
				err := s.addBackup(replicaID)
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

func (s *masterStore) addBackup(replicaID cluster.ReplicaID) error {
	replica := s.cluster.Replica(replicaID)
	if replica == nil {
		return fmt.Errorf("unknown replica node %s", replicaID)
	}
	conn, err := replica.Connect()
	if err != nil {
		return err
	}
	client := requests.NewRequestsServiceClient(conn)
	backup, err := newMasterBackup(s, replicaID, client, s.log.OpenReader(0))
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.backups = append(s.backups, backup)
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

		indexes := make([]Index, len(s.backups))
		i := 0
		for _, index := range s.commitIndexes {
			indexes[i] = index
			i++
		}
		sort.Slice(indexes, func(i, j int) bool {
			return indexes[i] < indexes[j]
		})

		index := indexes[s.config.GetBackups()-1]
		s.mu.Lock()
		defer s.mu.Unlock()
		commitIndex := s.state.getCommitIndex()
		if index > commitIndex {
			for i := commitIndex + 1; i <= index; i++ {
				s.state.setCommitIndex(i)
				ch, ok := s.commitChannels[i]
				if ok {
					ch <- i
					delete(s.commitChannels, i)
				}
			}
			logger.Debugf("Committed entries up to %d", index)
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
			Term:        uint64(b.store.state.getMastership().Term),
			Entries:     entries,
			PrevIndex:   uint64(batch.PrevIndex),
			CommitIndex: uint64(b.store.state.getCommitIndex()),
			AckIndex:    uint64(b.store.state.getAckIndex()),
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