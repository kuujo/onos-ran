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
	"fmt"

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
	if req.Subscribe {
		return fmt.Errorf("subscribe not yet implemented")
	}

	if req.Ecgi == nil {
		controlUpdates, err := manager.GetManager().GetControlUpdates()
		if err != nil {
			return err
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
	} else {
		return fmt.Errorf("list stations for specific ecgi not yet implemented")
	}
	return nil
}

// ListStationLinks returns a stream of links between neighboring base stations.
func (s Server) ListStationLinks(req *nb.StationLinkListRequest, stream nb.C1InterfaceService_ListStationLinksServer) error {
	if req.Ecgi == nil {
		controlUpdates, err := manager.GetManager().GetControlUpdates()
		if err != nil {
			return err
		}
		for _, update := range controlUpdates {
			switch update.GetMessageType() {
			case sb.MessageType_CELL_CONFIG_REPORT:
				cellConfigReport := update.GetCellConfigReport()
				ecgi := nb.ECGI{
					Ecid:   cellConfigReport.GetEcgi().GetEcid(),
					Plmnid: cellConfigReport.GetEcgi().GetPlmnId(),
				}
				stationLinkInfo := nb.StationLinkInfo{
					Ecgi: &ecgi,
				}
				candScells := cellConfigReport.GetCandScells()
				for _, candScell := range candScells {
					candCellEcgi := candScell.GetEcgi()
					nbEcgi := nb.ECGI{
						Ecid:   candCellEcgi.GetEcid(),
						Plmnid: candCellEcgi.GetPlmnId(),
					}
					stationLinkInfo.NeighborECGI = append(stationLinkInfo.NeighborECGI, &nbEcgi)
				}
				err = stream.Send(&stationLinkInfo)
				if err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("req ecgi is not nil")
	}
	return nil
}

// ListUELinks returns a stream of UI and base station links; one-time or (later) continuous subscribe.
func (s Server) ListUELinks(req *nb.UELinkListRequest, stream nb.C1InterfaceService_ListUELinksServer) error {
	if req.Ecgi == nil {
		telemetry, err := manager.GetManager().GetTelemetry()
		if err != nil {
			return err
		}
		for _, msg := range telemetry {
			switch msg.GetMessageType() {
			case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
				radioReportUe := msg.GetRadioMeasReportPerUE()
				ecgi := nb.ECGI{
					Ecid:   radioReportUe.GetEcgi().GetEcid(),
					Plmnid: radioReportUe.GetEcgi().GetPlmnId(),
				}
				radioReportServCells := radioReportUe.GetRadioReportServCells()
				var cqis []*nb.ChannelQuality
				for _, radioReportServCell := range radioReportServCells {
					servCellEcgi := radioReportServCell.GetEcgi()
					ecgi := nb.ECGI{
						Ecid:   servCellEcgi.GetEcid(),
						Plmnid: servCellEcgi.GetPlmnId(),
					}
					cqiHist := radioReportServCell.GetCqiHist()
					for _, cqi := range cqiHist {
						nbCqi := nb.ChannelQuality{
							TargetEcgi: &ecgi,
							CqiHist:    cqi,
						}
						cqis = append(cqis, &nbCqi)
					}

				}
				ueLinkInfo := nb.UELinkInfo{
					Ecgi:             &ecgi,
					Crnti:            radioReportUe.GetCrnti(),
					ChannelQualities: cqis,
				}

				err = stream.Send(&ueLinkInfo)
				if err != nil {
					return err
				}

			}
		}

	} else {

		return fmt.Errorf("UELinkListRequest is not empty")
	}
	return nil
}

// TriggerHandOver returns a hand-over response indicating success or failure.
func (s Server) TriggerHandOver(ctx context.Context, req *nb.HandOverRequest) (*nb.HandOverResponse, error) {
	if req != nil {
		src := req.GetSrcStation()
		dst := req.GetDstStation()
		crnti := req.GetCrnti()

		srcEcgi := sb.ECGI{
			Ecid:   src.GetEcid(),
			PlmnId: src.GetPlmnid(),
		}

		dstEcgi := sb.ECGI{
			Ecid:   dst.GetEcid(),
			PlmnId: dst.GetPlmnid(),
		}

		ctrlResponse := sb.ControlResponse{
			MessageType: sb.MessageType_HO_REQUEST,
			S: &sb.ControlResponse_HORequest{
				HORequest: &sb.HORequest{
					Crnti: crnti,
					EcgiS: &srcEcgi,
					EcgiT: &dstEcgi,
				},
			},
		}

		err := manager.GetManager().SB.SendResponse(ctrlResponse)
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("HandOverRequest is nil")
}

// SetRadioPower returns a response indicating success or failure.
func (s Server) SetRadioPower(ctx context.Context, req *nb.RadioPowerRequest) (*nb.RadioPowerResponse, error) {
	if req != nil {
		offset := req.GetOffset()
		pa := make([]sb.XICICPA, 1)
		switch offset {
		case nb.StationPowerOffset_PA_DB_0:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_0)
		case nb.StationPowerOffset_PA_DB_1:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_1)
		case nb.StationPowerOffset_PA_DB_2:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_2)
		case nb.StationPowerOffset_PA_DB_3:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_3)
		case nb.StationPowerOffset_PA_DB_MINUS3:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_MINUS3)
		case nb.StationPowerOffset_PA_DB_MINUS6:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_MINUS6)
		case nb.StationPowerOffset_PA_DB_MINUS1DOT77:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_MINUS1DOT77)
		case nb.StationPowerOffset_PA_DB_MINUX4DOT77:
			pa = append(pa, sb.XICICPA_XICIC_PA_DB_MINUX4DOT77)

		}

		ecgi := sb.ECGI{
			Ecid:   req.GetEcgi().GetEcid(),
			PlmnId: req.GetEcgi().GetPlmnid(),
		}

		ctrlResponse := sb.ControlResponse{
			MessageType: sb.MessageType_RRM_CONFIG,
			S: &sb.ControlResponse_RRMConfig{
				RRMConfig: &sb.RRMConfig{
					Ecgi: &ecgi,
					PA:   pa,
				},
			},
		}
		err := manager.GetManager().SB.SendResponse(ctrlResponse)
		if err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("SetRadioPower request cannot be nil")
}
