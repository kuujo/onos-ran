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
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/store/requests"
	"google.golang.org/grpc"
	"sync"
)

func init() {
	atomix.RegisterService(func(id cluster.NodeID, server *grpc.Server) {
		requests.RegisterRequestsServiceServer(server, newServer())
	})
}

var server *storeServer

func getServer() *storeServer {
	return server
}

func newServer() requests.RequestsServiceServer {
	server = &storeServer{}
	return server
}

type storeHandler interface {
	Store
	append(context.Context, *requests.AppendRequest) (*requests.AppendResponse, error)
	ack(context.Context, *requests.AckRequest) (*requests.AckResponse, error)
	backup(context.Context, *requests.BackupRequest) (*requests.BackupResponse, error)
}

// storeServer is the server side for the indications store
type storeServer struct {
	handler storeHandler
	mu      sync.RWMutex
}

func (s *storeServer) registerHandler(handler storeHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handler = handler
}

func (s *storeServer) getHandler() (storeHandler, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.handler == nil {
		return nil, errors.New("no store handler registered")
	}
	return s.handler, nil
}

func (s *storeServer) Append(ctx context.Context, request *requests.AppendRequest) (*requests.AppendResponse, error) {
	handler, err := s.getHandler()
	if err != nil {
		return nil, err
	}
	return handler.append(ctx, request)
}

func (s *storeServer) Ack(ctx context.Context, request *requests.AckRequest) (*requests.AckResponse, error) {
	handler, err := s.getHandler()
	if err != nil {
		return nil, err
	}
	return handler.ack(ctx, request)
}

func (s *storeServer) Backup(ctx context.Context, request *requests.BackupRequest) (*requests.BackupResponse, error) {
	handler, err := s.getHandler()
	if err != nil {
		return nil, err
	}
	return handler.backup(ctx, request)
}
