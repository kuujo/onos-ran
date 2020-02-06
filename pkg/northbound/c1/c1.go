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
	"context"

	"github.com/onosproject/onos-ran/api/sb"

	"github.com/onosproject/onos-ran/api/nb"

	"github.com/onosproject/onos-ran/pkg/manager"
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

// ListStations returns a stream of base station records.
func (s Server) ListStations(req *nb.StationListRequest, stream nb.C1InterfaceService_ListStationsServer) error {
	if req.Ecgi == nil {
		controlUpdates, err := manager.GetManager().GetControlUpdates()
		if err != nil {
			return nil
		}
		for _, update := range controlUpdates {
			switch update.GetMessageType() {
			case sb.MessageType_CELL_CONFIG_REPORT:
				cellConfigReport := update.GetCellConfigReport()
				ecgi := nb.ECGI{
					Ecid:   cellConfigReport.GetEcgi().GetEcid(),
					Plmnid: cellConfigReport.GetEcgi().GetPlmnId(),
				}
				baseStationInfo := nb.StationInfo{
					Ecgi: &ecgi,
				}
				baseStationInfo.MaxNumConnectedUes = cellConfigReport.GetMaxNumConnectedUes()
				err = stream.Send(&baseStationInfo)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ListUEs returns a stream of UE records.
func (s Server) ListUEs(req *nb.UEListRequest, stream nb.C1InterfaceService_ListUEsServer) error {
	panic("implement me")
}

// ListStationLinks returns a stream of links between neighboring base stations.
func (s Server) ListStationLinks(req *nb.StationLinkListRequest, stream nb.C1InterfaceService_ListStationLinksServer) error {
	panic("implement me")
}

// ListUELinks returns a stream of UI and base station links; one-time or (later) continuous subscribe.
func (s Server) ListUELinks(*nb.UELinkListRequest, nb.C1InterfaceService_ListUELinksServer) error {
	panic("implement me")
}

// TriggerHandOver returns a hand-over response indicating success or failure.
func (s Server) TriggerHandOver(context.Context, *nb.HandOverRequest) (*nb.HandOverResponse, error) {
	panic("implement me")
}

// SetRadioPower returns a response indicating success or failure.
func (s Server) SetRadioPower(context.Context, *nb.RadioPowerRequest) (*nb.RadioPowerResponse, error) {
	panic("implement me")
}
