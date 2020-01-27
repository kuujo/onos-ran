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

package c1

import (
	"github.com/onosproject/onos-ran/api/nb"
	"github.com/onosproject/onos-ran/pkg/service"
	"google.golang.org/grpc"
)

// NewService returns a new device Service
func NewService() (service.Service, error) {
	return &Service{}, nil
}

// Service is an implementation of C1 service.
type Service struct {
	service.Service
}

// Register registers the C1 Service with the gRPC server.
func (s Service) Register(r *grpc.Server) {
	server := &Server{}
	nb.RegisterC1InterfaceServiceServer(r, server)
}

// Server implements the C1 gRPC service for administrative facilities.
type Server struct {
}

// GetRNIB blah blah. WIP!
func (s Server) GetRNIB(*nb.C1RequestMessage, nb.C1InterfaceService_GetRNIBServer) error {
	panic("implement me")
}

// PostRNIB blah blah. WIP!
func (s Server) PostRNIB(*nb.C1RequestMessage, nb.C1InterfaceService_PostRNIBServer) error {
	panic("implement me")
}
