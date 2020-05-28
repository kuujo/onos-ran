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

package indications

import (
	"errors"
	"github.com/onosproject/onos-lib-go/pkg/atomix"
	"github.com/onosproject/onos-lib-go/pkg/cluster"
	"github.com/onosproject/onos-ric/api/store/indications"
	"google.golang.org/grpc"
	"sync"
)

func init() {
	atomix.RegisterService(func(id cluster.NodeID, server *grpc.Server) {
		indications.RegisterIndicationsServiceServer(server, newServer())
	})
}

var server *storeServer

func getServer() *storeServer {
	return server
}

func newServer() indications.IndicationsServiceServer {
	server = &storeServer{}
	return server
}

type storeHandler interface {
	subscribe(*indications.SubscribeRequest, indications.IndicationsService_SubscribeServer) error
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

func (s *storeServer) Subscribe(request *indications.SubscribeRequest, stream indications.IndicationsService_SubscribeServer) error {
	handler, err := s.getHandler()
	if err != nil {
		return err
	}
	return handler.subscribe(request, stream)
}