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
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"io"
	"sync"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	service "github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/pkg/manager"
	"github.com/onosproject/onos-ric/pkg/store/indications"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("northbound", "c1")

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
	log.Debugf("Received StationListRequest %+v", req)
	if req.Ecgi == nil {
		ch := make(chan e2ap.RicIndication)
		if req.Subscribe {
			watchCh := make(chan indications.Event)
			if err := manager.GetManager().SubscribeIndications(watchCh); err != nil {
				return err
			}
			go func() {
				defer close(ch)
				for event := range watchCh {
					ch <- *event.Indication.RicIndication
				}
			}()
		} else {
			if err := manager.GetManager().ListIndications(ch); err != nil {
				return err
			}
		}

		for update := range ch {
			switch update.GetHdr().GetMessageType() {
			case sb.MessageType_CELL_CONFIG_REPORT:
				cellConfigReport := update.GetMsg().GetCellConfigReport()
				ecgi := nb.ECGI{
					Ecid:   cellConfigReport.GetEcgi().GetEcid(),
					Plmnid: cellConfigReport.GetEcgi().GetPlmnId(),
				}
				baseStationInfo := nb.StationInfo{
					Ecgi: &ecgi,
				}
				baseStationInfo.MaxNumConnectedUes = cellConfigReport.GetMaxNumConnectedUes()
				log.Debugf("Sending StationInfo %v", baseStationInfo)
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
	log.Debugf("Received StationLinkListRequest %+v", req)
	if req.Ecgi == nil {
		ch := make(chan e2ap.RicIndication)
		if req.Subscribe {
			watchCh := make(chan indications.Event)
			if err := manager.GetManager().SubscribeIndications(watchCh); err != nil {
				return err
			}
			go func() {
				defer close(ch)
				for event := range watchCh {
					ch <- *event.Indication.RicIndication
				}
			}()
		} else {
			if err := manager.GetManager().ListIndications(ch); err != nil {
				return err
			}
		}

		for update := range ch {
			switch update.GetHdr().GetMessageType() {
			case sb.MessageType_CELL_CONFIG_REPORT:
				cellConfigReport := update.GetMsg().GetCellConfigReport()
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
				log.Debugf("Sending StationLinkInfo %v", stationLinkInfo)
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

// ListUEs returns a stream of UEs
func (s Server) ListUEs(req *nb.UEListRequest, stream nb.C1InterfaceService_ListUEsServer) error {
	log.Debugf("Received UEListRequest %+v", req)
	if req.Ecgi == nil {
		ch := make(chan e2ap.RicIndication)
		if req.Subscribe {
			watchCh := make(chan indications.Event)
			if err := manager.GetManager().SubscribeIndications(watchCh); err != nil {
				return err
			}
			go func() {
				defer close(ch)
				for event := range watchCh {
					ch <- *event.Indication.RicIndication
				}
			}()
		} else {
			if err := manager.GetManager().ListIndications(ch); err != nil {
				return err
			}
		}

		for update := range ch {
			switch update.GetHdr().GetMessageType() {
			case sb.MessageType_UE_ADMISSION_REQUEST:
				ueAdmReq := update.GetMsg().GetUEAdmissionRequest()
				crnti := ueAdmReq.GetCrnti()
				ecgi := nb.ECGI{
					Plmnid: ueAdmReq.GetEcgi().GetPlmnId(),
					Ecid:   ueAdmReq.GetEcgi().GetEcid(),
				}
				imsi := fmt.Sprintf("%d", ueAdmReq.GetImsi())
				ueInfo := nb.UEInfo{
					Crnti: crnti,
					Ecgi:  &ecgi,
					Imsi:  imsi,
				}
				log.Debugf("Sending UEInfo %v", ueInfo)
				if err := stream.Send(&ueInfo); err != nil {
					return err
				}
			}
		}

	} else {
		return fmt.Errorf("list UEs for specific ecgi and crnti not yet implemented")
	}

	return nil
}

// ListUELinks returns a stream of UI and base station links; one-time or (later) continuous subscribe.
func (s Server) ListUELinks(req *nb.UELinkListRequest, stream nb.C1InterfaceService_ListUELinksServer) error {
	log.Debugf("Received UELinkListRequest %+v", req)
	if req.Ecgi == nil {
		ch := make(chan e2ap.RicIndication)
		if req.Subscribe {
			watchCh := make(chan indications.Event)
			var err error
			if req.NoReplay {
				err = manager.GetManager().SubscribeIndications(watchCh)
			} else {
				err = manager.GetManager().SubscribeIndications(watchCh, indications.WithReplay())
			}
			if err != nil {
				return err
			}
			go func() {
				defer close(ch)
				for event := range watchCh {
					ch <- *event.Indication.RicIndication
				}
			}()
		} else {
			if err := manager.GetManager().ListIndications(ch); err != nil {
				return err
			}
		}

		for telemetry := range ch {
			switch telemetry.GetHdr().GetMessageType() {
			case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
				radioReportUe := telemetry.GetMsg().GetRadioMeasReportPerUE()
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
				log.Debugf("Sending UELinkInfo %v", ueLinkInfo)
				if err := stream.Send(&ueLinkInfo); err != nil {
					return err
				}
			default:
			}
		}
	} else {
		return fmt.Errorf("listuelinks for specific crnti and ecgi not yet implemented %v", req)
	}
	return nil
}

// TriggerHandOver returns a hand-over response indicating success or failure.
func (s Server) TriggerHandOver(ctx context.Context, req *nb.HandOverRequest) (*nb.HandOverResponse, error) {
	log.Debugf("Received HandOverRequest %+v", req)
	if req == nil || req.GetCrnti() == "" ||
		req.DstStation == nil || req.DstStation.Plmnid == "" || req.DstStation.Ecid == "" ||
		req.SrcStation == nil || req.SrcStation.Plmnid == "" || req.SrcStation.Ecid == "" {

		return nil, fmt.Errorf("HandOverRequest is missing values %v", req)
	}
	return sendHandoverTrigger(req)
}

// TriggerHandOverStream is a version that stays open all the time.
func (s Server) TriggerHandOverStream(stream nb.C1InterfaceService_TriggerHandOverStreamServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&nb.HandOverResponse{Success: true})
		}
		if err != nil {
			return err
		}
		log.Debugf("Received HandOverRequest %+v", req)
		if req == nil || req.GetCrnti() == "" ||
			req.DstStation == nil || req.DstStation.Plmnid == "" || req.DstStation.Ecid == "" ||
			req.SrcStation == nil || req.SrcStation.Plmnid == "" || req.SrcStation.Ecid == "" {

			log.Errorf("HandOverRequest is missing values %v", req)
			continue
		}

		if _, err = sendHandoverTrigger(req); err != nil {
			log.Warn("Error in sending HO trigger %v", err)
		}
	}
}

func sendHandoverTrigger(req *nb.HandOverRequest) (*nb.HandOverResponse, error) {
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

	srcHoReq := &e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &srcEcgi,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_HORequest{
				HORequest: &sb.HORequest{
					Crntis: []string{crnti},
					EcgiS:  &srcEcgi,
					EcgiT:  &dstEcgi,
				},
			},
		},
	}

	dstHoReq := &e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
			Ecgi:        &dstEcgi,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_HORequest{
				HORequest: &sb.HORequest{
					Crntis: []string{crnti},
					EcgiS:  &srcEcgi,
					EcgiT:  &dstEcgi,
				},
			},
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	errCh := make(chan error)
	go func() {
		err := manager.GetManager().StoreRicControlRequest(srcHoReq)
		if err != nil {
			errCh <- err
		}
		wg.Done()
	}()
	go func() {
		err := manager.GetManager().StoreRicControlRequest(dstHoReq)
		if err != nil {
			errCh <- err
		}
		wg.Done()
	}()
	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
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
	var pa sb.XICICPA
	switch offset {
	case nb.StationPowerOffset_PA_DB_0:
		pa = sb.XICICPA_XICIC_PA_DB_0
	case nb.StationPowerOffset_PA_DB_1:
		pa = sb.XICICPA_XICIC_PA_DB_1
	case nb.StationPowerOffset_PA_DB_2:
		pa = sb.XICICPA_XICIC_PA_DB_2
	case nb.StationPowerOffset_PA_DB_3:
		pa = sb.XICICPA_XICIC_PA_DB_3
	case nb.StationPowerOffset_PA_DB_MINUS3:
		pa = sb.XICICPA_XICIC_PA_DB_MINUS3
	case nb.StationPowerOffset_PA_DB_MINUS6:
		pa = sb.XICICPA_XICIC_PA_DB_MINUS6
	case nb.StationPowerOffset_PA_DB_MINUS1DOT77:
		pa = sb.XICICPA_XICIC_PA_DB_MINUS1DOT77
	case nb.StationPowerOffset_PA_DB_MINUX4DOT77:
		pa = sb.XICICPA_XICIC_PA_DB_MINUX4DOT77

	}

	ecgi := sb.ECGI{
		Ecid:   req.GetEcgi().GetEcid(),
		PlmnId: req.GetEcgi().GetPlmnid(),
	}

	var p []sb.XICICPA
	p = append(p, pa)
	rrmConfigReq := &e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_RRM_CONFIG,
			Ecgi:        &ecgi,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_RRMConfig{
				RRMConfig: &sb.RRMConfig{
					Ecgi: &ecgi,
					PA:   p,
				},
			},
		},
	}

	err := manager.GetManager().StoreRicControlRequest(rrmConfigReq)
	if err != nil {
		return nil, err
	}
	return &nb.RadioPowerResponse{Success: true}, nil
}
