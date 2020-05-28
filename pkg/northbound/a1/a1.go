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

package a1

import (
	"context"

	"github.com/onosproject/onos-lib-go/pkg/logging/service"
	"github.com/onosproject/onos-ric/api/nb/a1"
	"google.golang.org/grpc"
)

//var log = logging.GetLogger("a1")

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
	a1.RegisterA1Server(r, server)
}

// Server implements the C1 gRPC service for administrative facilities.
type Server struct {
}

// CreateOrUpdate creates or updates an A1 policy
func (s Server) CreateOrUpdate(ctx context.Context, req *a1.CreateOrUpdateRequest) (*a1.CreateOrUpdateResponse, error) {
	// TODO implement me
	return &a1.CreateOrUpdateResponse{}, nil
}

// Query queries about one or more than A1 policies
func (s Server) Query(req *a1.QueryRequest, stream a1.A1_QueryServer) error {
	// TODO implement me
	return nil
}

// Delete deletes an A1 policy
func (s Server) Delete(ctx context.Context, req *a1.DeleteRequest) (*a1.DeleteResponse, error) {
	// TODO implement me
	resp := &a1.DeleteResponse{}
	return resp, nil
}

// Notify notify about an enforcement status change of a policy between 'enforced' and 'not enforced'.
func (s Server) Notify(stream a1.A1_NotifyServer) error {
	// TODO implement me
	return nil
}
