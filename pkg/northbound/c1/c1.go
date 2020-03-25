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

	"github.com/onosproject/onos-ric/api/nb"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"

	"github.com/onosproject/onos-ric/pkg/manager"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("c1")

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
		ch := make(chan sb.ControlUpdate)
		if req.Subscribe {
			if err := manager.GetManager().SubscribeControlUpdates(ch); err != nil {
				return err
			}
		} else {
			if err := manager.GetManager().ListControlUpdates(ch); err != nil {
				return err
			}
		}

		for update := range ch {
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
				if err := stream.Send(&baseStationInfo); err != nil {
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
	if req.Subscribe {
		return fmt.Errorf("subscribe not yet implemented")
	}
	if req.Ecgi == nil {
		ch := make(chan sb.ControlUpdate)
		if req.Subscribe {
			if err := manager.GetManager().SubscribeControlUpdates(ch); err != nil {
				return err
			}
		} else {
			if err := manager.GetManager().ListControlUpdates(ch); err != nil {
				return err
			}
		}

		for update := range ch {
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
				if err := stream.Send(&stationLinkInfo); err != nil {
					return err
				}
			}
		}
	} else {
		return fmt.Errorf("list station links for specific ecgi not yet implemented")
	}
	return nil
}

// ListUELinks returns a stream of UI and base station links; one-time or (later) continuous subscribe.
func (s Server) ListUELinks(req *nb.UELinkListRequest, stream nb.C1InterfaceService_ListUELinksServer) error {
	if req.Ecgi == nil {
		ch := make(chan sb.TelemetryMessage)
		if req.Subscribe {
			if err := manager.GetManager().SubscribeTelemetry(ch, !req.NoReplay); err != nil {
				return err
			}
		} else {
			if err := manager.GetManager().ListTelemetry(ch); err != nil {
				return err
			}
		}

		for telemetry := range ch {
			switch telemetry.GetMessageType() {
			case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
				radioReportUe := telemetry.GetRadioMeasReportPerUE()
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
				update, err := manager.GetManager().GetUEAdmissionByID(radioReportUe.GetEcgi(), radioReportUe.Crnti)
				if err == nil {
					ueLinkInfo.Imsi = fmt.Sprintf("%d", update.GetUEAdmissionRequest().GetImsi())
				} else {
					log.Infof("Cannot find Imsi for %v %s", radioReportUe.Ecgi, radioReportUe.Crnti)
				}

				if err := stream.Send(&ueLinkInfo); err != nil {
					return err
				}
			default:
				log.Errorf("Unhandled case %s", telemetry.GetMessageType())
			}
		}
	} else {
		return fmt.Errorf("listuelinks for specific crnti and ecgi not yet implemented %v", req)
	}
	return nil
}

// TriggerHandOver returns a hand-over response indicating success or failure.
func (s Server) TriggerHandOver(ctx context.Context, req *nb.HandOverRequest) (*nb.HandOverResponse, error) {
	if req == nil || req.Crnti == "" ||
		req.DstStation == nil || req.DstStation.Plmnid == "" || req.DstStation.Ecid == "" ||
		req.SrcStation == nil || req.SrcStation.Plmnid == "" || req.SrcStation.Ecid == "" {

		return nil, fmt.Errorf("HandOverRequest is missing values %v", req)
	}
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

	hoReq := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_HORequest{
				HORequest: &sb.HORequest{
					Crnti: crnti,
					EcgiS: &srcEcgi,
					EcgiT: &dstEcgi,
				},
			},
		},
	}

	srcSession, ok := manager.GetManager().SbSessions[srcEcgi]
	if !ok {
		return nil, fmt.Errorf("session not found for HO source %v", srcEcgi)
	}
	log.Infof("Sending HO for %v to %s %v", srcEcgi, srcSession.EndPoint, srcSession.Ecgi)
	err := srcSession.SendRicControlRequest(hoReq)
	if err != nil {
		return nil, err
	}

	dstSession, ok := manager.GetManager().SbSessions[dstEcgi]
	if !ok {
		return nil, fmt.Errorf("session not found for HO dest %v", dstEcgi)
	}
	log.Infof("Sending HO for %v to %s %v", srcEcgi, dstSession.EndPoint, dstSession.Ecgi)
	err = dstSession.SendRicControlRequest(hoReq)
	if err != nil {
		return nil, err
	}

	err = manager.GetManager().DeleteTelemetry(src.GetPlmnid(), src.GetEcid(), crnti)
	if err != nil {
		return nil, err
	}
	err = manager.GetManager().DeleteUEAdmissionRequest(src.GetPlmnid(), src.GetEcid(), crnti)
	if err != nil {
		return nil, err
	}

	return &nb.HandOverResponse{Success: true}, nil
}

// SetRadioPower returns a response indicating success or failure.
func (s Server) SetRadioPower(ctx context.Context, req *nb.RadioPowerRequest) (*nb.RadioPowerResponse, error) {
	if req == nil || req.Ecgi == nil || req.Ecgi.Plmnid == "" || req.Ecgi.Ecid == "" || req.Offset < 0 {
		return nil, fmt.Errorf("SetRadioPower is missing values %v", req)
	}
	offset := req.GetOffset()
	var pa []sb.XICICPA
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

	rrmConfigReq := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_RRM_CONFIG,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_RRMConfig{
				RRMConfig: &sb.RRMConfig{
					Ecgi: &ecgi,
					PA:   pa,
				},
			},
		},
	}
	session, ok := manager.GetManager().SbSessions[ecgi]
	if !ok {
		return nil, fmt.Errorf("session not found for Power Request %v", ecgi)
	}
	err := session.SendRicControlRequest(rrmConfigReq)
	if err != nil {
		return nil, err
	}

	return &nb.RadioPowerResponse{Success: true}, nil
}
