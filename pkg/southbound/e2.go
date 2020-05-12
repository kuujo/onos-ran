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

package southbound

import (
	"context"
	"fmt"
	"sync"
	"time"

	topodevice "github.com/onosproject/onos-topo/api/device"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"

	"github.com/onosproject/onos-ric/api/sb"
	e2ap "github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("southbound")

//mapHOEventMeasuredRIC ...
var mapHOEventMeasuredRIC map[string]HOEventMeasuredRIC

// mutexmapHOEventMeasuredRIC ...
var mutexmapHOEventMeasuredRIC sync.RWMutex

// Session is responsible for managing connections to and interactions with the RAN southbound.
type Session struct {
	EndPoint                      sb.Endpoint
	Ecgi                          sb.ECGI
	client                        e2ap.E2APClient
	ricControlRequestChan         chan e2ap.RicControlRequest
	ricIndicationChan             chan e2ap.RicIndication
	controlIndications            chan e2ap.RicIndication
	telemetryIndications          chan e2ap.RicIndication
	RicControlResponseHandlerFunc RicControlResponseHandler
	ControlUpdateHandlerFunc      ControlUpdateHandler
	TelemetryUpdateHandlerFunc    TelemetryUpdateHandler
	TopoUpdateHandler             TopoUpdateHandler
	EnableMetrics                 bool
	e2Chan                        chan E2
	connection                    *grpc.ClientConn
}

// NewSession creates a new southbound session controller.
func NewSession() (E2, error) {
	return &Session{
		ricControlRequestChan: make(chan e2ap.RicControlRequest),
		ricIndicationChan:     make(chan e2ap.RicIndication),
		controlIndications:    make(chan e2ap.RicIndication),
		telemetryIndications:  make(chan e2ap.RicIndication),
	}, nil
}

// Run starts the southbound control loop.
func (s *Session) Run(ecgi sb.ECGI, endPoint sb.Endpoint,
	tls topodevice.TlsConfig, creds topodevice.Credentials,
	storeRicControlResponse RicControlResponseHandler,
	storeControlUpdate ControlUpdateHandler,
	storeTelemetry TelemetryUpdateHandler,
	topoUpdateHandler TopoUpdateHandler,
	enableMetrics bool,
	e2Chan chan E2) {
	s.EndPoint = endPoint
	s.Ecgi = ecgi
	s.RicControlResponseHandlerFunc = storeRicControlResponse
	s.ControlUpdateHandlerFunc = storeControlUpdate
	s.TelemetryUpdateHandlerFunc = storeTelemetry
	s.TopoUpdateHandler = topoUpdateHandler
	s.EnableMetrics = enableMetrics
	s.e2Chan = e2Chan
	mapHOEventMeasuredRIC = make(map[string]HOEventMeasuredRIC)
	go s.manageConnections(tls, creds)
	go s.recvTelemetryUpdates()
	go s.recvUpdates()
	go s.recvControlResponses()
}

// Setup ...
func (s *Session) Setup() error {
	req := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_CellConfigRequest{
				CellConfigRequest: &sb.CellConfigRequest{},
			},
		},
	}
	err := s.sendRicControlRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// UeHandover ...
func (s *Session) UeHandover(crnti []string, srcEcgi sb.ECGI, dstEcgi sb.ECGI) error {

	hoReq := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_HO_REQUEST,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_HORequest{
				HORequest: &sb.HORequest{
					Crntis: crnti,
					EcgiS:  &srcEcgi,
					EcgiT:  &dstEcgi,
				},
			},
		},
	}
	err := s.sendRicControlRequest(hoReq)
	if err != nil {
		return err
	}
	return nil
}

// RemoteAddress ...
func (s *Session) RemoteAddress() sb.Endpoint {
	return s.EndPoint
}

// RRMConfig ...
func (s *Session) RRMConfig(pa sb.XICICPA) error {
	var p []sb.XICICPA
	p = append(p, pa)
	rrmConfigReq := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_RRM_CONFIG,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_RRMConfig{
				RRMConfig: &sb.RRMConfig{
					Ecgi: &s.Ecgi,
					PA:   p,
				},
			},
		},
	}
	err := s.sendRicControlRequest(rrmConfigReq)
	if err != nil {
		return err
	}
	return nil
}

// L2MeasConfig ...
func (s *Session) L2MeasConfig(l2MeasConfig *sb.L2MeasConfig) error {
	req := e2ap.RicControlRequest{
		Hdr: &e2sm.RicControlHeader{
			MessageType: sb.MessageType_L2_MEAS_CONFIG,
		},
		Msg: &e2sm.RicControlMessage{
			S: &e2sm.RicControlMessage_L2MeasConfig{
				L2MeasConfig: l2MeasConfig,
				/*
							L2MeasConfig: &l2ReportInterval,
							},
						},
					},
				*/
			},
		},
	}
	err := s.sendRicControlRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// Close - close the session if cell is removed from topo
func (s *Session) Close() {
	s.connection.Close()
	close(s.ricControlRequestChan)
	close(s.controlIndications)
	close(s.ricIndicationChan)
	close(s.telemetryIndications)
	// Don't close(s.e2Chan) - common to all sessions
}

func (s *Session) recvTelemetryUpdates() {
	for update := range s.telemetryIndications {
		s.TelemetryUpdateHandlerFunc(update)
	}
}

func (s *Session) recvUpdates() {
	for update := range s.controlIndications {
		s.ControlUpdateHandlerFunc(update)
	}
}

func (s *Session) recvControlResponses() {
	for update := range s.ricIndicationChan {
		s.RicControlResponseHandlerFunc(update)
	}
}

// sendRicControlRequest sends the specified RicControlRequest on the control channel
func (s *Session) sendRicControlRequest(req e2ap.RicControlRequest) error {

	req.GetHdr().Ecgi = &s.Ecgi

	switch req.GetHdr().GetMessageType() {
	case sb.MessageType_HO_REQUEST:
		if s.EnableMetrics {
			tC := time.Now()
			go s.updateHOEventMeasuredRIC(req, tC)
		}
	default:
	}

	s.ricControlRequestChan <- req
	return nil
}

func (s *Session) manageConnections(tls topodevice.TlsConfig, creds topodevice.Credentials) {
	for {
		// Attempt to create connection to the E2Node
		log.Infof("Connecting to E2Node...%s with context", s.EndPoint)
		opts := []grpc.DialOption{
			grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(100 * time.Millisecond)),
		}
		var err error
		s.connection, err = southbound.Connect(context.Background(), string(s.EndPoint), tls.GetCert(), tls.GetKey(), opts...)
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			s.manageConnection(s.connection)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Session) manageConnection(connection *grpc.ClientConn) {
	// Offer the telemetry and control surfaces to the E2 devices
	s.client = e2ap.NewE2APClient(connection)
	if s.client == nil {
		return
	}
	log.Infof("Connected to E2Node on %s", s.EndPoint)
	// Setup coordination channel
	errors := make(chan error)

	go s.ricChan(errors)

	s.TopoUpdateHandler(&s.Ecgi, &topodevice.ProtocolState{
		Protocol:          topodevice.Protocol_E2AP,
		ConnectivityState: topodevice.ConnectivityState_REACHABLE,
		ChannelState:      topodevice.ChannelState_CONNECTED,
		ServiceState:      topodevice.ServiceState_AVAILABLE,
	})

	// Wait for the first error on the coordination channel.
	for i := 0; i < 2; i++ {
		<-errors
	}

	// Clean-up the connection, forcing other consumers to terminate and clean-up.
	_ = connection.Close()

	s.TopoUpdateHandler(&s.Ecgi, &topodevice.ProtocolState{
		Protocol:          topodevice.Protocol_E2AP,
		ConnectivityState: topodevice.ConnectivityState_UNKNOWN_CONNECTIVITY_STATE,
		ChannelState:      topodevice.ChannelState_DISCONNECTED,
		ServiceState:      topodevice.ServiceState_UNAVAILABLE,
	})

	// FIXME: should be done separately
	//close(s.responses)
	//close(s.updates)

	close(errors)
	log.Info("Disconnected from E2Node %s", s.EndPoint)
}

func (s *Session) ricChan(errors chan error) {
	stream, err := s.client.RicChan(context.Background())
	if err != nil {
		errors <- err
		return
	}
	waitc := make(chan error)
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Infof("Got ricIndication messageType %d from %s", resp.Hdr.MessageType, s.EndPoint)
			s.handleRicIndication(resp)
		}
	}()

	s.e2Chan <- s

	for {
		req := <-s.ricControlRequestChan
		err = stream.Send(&req)
		if err != nil {
			errors <- err
			return
		}
	}
}

func (s *Session) handleRicIndication(update *e2ap.RicIndication) {
	msgType := update.GetHdr().GetMessageType()
	log.Infof("ricIndication %T %s", msgType, msgType)
	switch msgType {
	case sb.MessageType_CELL_CONFIG_REPORT:
		msg := update.GetMsg().GetCellConfigReport()
		log.Infof("%s CellConfigReport plmnid:%s, ecid:%s", s.EndPoint, msg.GetEcgi().GetPlmnId(), msg.GetEcgi().GetEcid())
		s.ricIndicationChan <- *update
	case sb.MessageType_RADIO_MEAS_REPORT_PER_UE:
		if s.EnableMetrics {
			tT := time.Now()
			go s.addReceivedHOEventMeasuredRIC(update, tT) // to monitor HO delay]
		}
		s.telemetryIndications <- *update
	case sb.MessageType_UE_ADMISSION_REQUEST, sb.MessageType_UE_RELEASE_IND:
		s.controlIndications <- *update
	default:
		log.Fatalf("%s ControlResponse has unexpected type %T", s.EndPoint, msgType)
	}
}

func (s *Session) addReceivedHOEventMeasuredRIC(update *e2ap.RicIndication, tT time.Time) {
	servStationID := update.GetMsg().GetRadioMeasReportPerUE().GetEcgi()
	numNeighborCells := len(update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells())
	bestStationID := update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetEcgi()
	bestCQI := update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetCqiHist()[0]

	for i := 1; i < numNeighborCells; i++ {
		tmpCQI := update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[i].GetCqiHist()[0]
		if bestCQI < tmpCQI {
			bestStationID = update.GetMsg().GetRadioMeasReportPerUE().GetRadioReportServCells()[i].GetEcgi()
			bestCQI = tmpCQI
		}
	}

	if servStationID.GetEcid() != bestStationID.GetEcid() || servStationID.GetPlmnId() != bestStationID.GetPlmnId() {
		tmpHOEventMeasuredRIC := HOEventMeasuredRIC{
			Timestamp:  tT,
			Crnti:      update.GetMsg().GetRadioMeasReportPerUE().GetCrnti(),
			DestPlmnID: bestStationID.GetPlmnId(),
			DestECID:   bestStationID.GetEcid(),
		}
		key := fmt.Sprintf("%s:%s:%s:%s:%s", update.GetMsg().GetRadioMeasReportPerUE().GetCrnti(), servStationID.GetPlmnId(), servStationID.GetEcid(), bestStationID.GetPlmnId(), bestStationID.GetEcid())
		mutexmapHOEventMeasuredRIC.Lock()
		_, f := mapHOEventMeasuredRIC[key]
		if !f {
			mapHOEventMeasuredRIC[key] = tmpHOEventMeasuredRIC
		}
		mutexmapHOEventMeasuredRIC.Unlock()
	}
}

func (s *Session) updateHOEventMeasuredRIC(req e2ap.RicControlRequest, tC time.Time) {
	key := fmt.Sprintf("%s:%s:%s:%s:%s", req.GetMsg().GetHORequest().GetCrnti(),
		req.GetMsg().GetHORequest().GetEcgiS().GetPlmnId(),
		req.GetMsg().GetHORequest().GetEcgiS().GetEcid(),
		req.GetMsg().GetHORequest().GetEcgiT().GetPlmnId(),
		req.GetMsg().GetHORequest().GetEcgiT().GetEcid())
	mutexmapHOEventMeasuredRIC.Lock()
	v, f := mapHOEventMeasuredRIC[key]
	if !f {
		//log.Warnf("Telemetry message is missing when calculating HO latency - %s", key)
	} else {
		v.ElapsedTime = tC.Sub(v.Timestamp).Microseconds()
		ChanHOEvent <- v
		delete(mapHOEventMeasuredRIC, key)
	}
	mutexmapHOEventMeasuredRIC.Unlock()
}
