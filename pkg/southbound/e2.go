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
	topodevice "github.com/onosproject/onos-topo/api/device"
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"

	"github.com/onosproject/onos-ric/api/sb"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("southbound")

// Session is responsible for managing connections to and interactions with the RAN southbound.
type Session struct {
	EndPoint sb.Endpoint
	Ecgi     sb.ECGI
	client   sb.InterfaceServiceClient

	controlResponses chan sb.ControlResponse
	controlUpdates   chan sb.ControlUpdate

	telemetryUpdates chan sb.TelemetryMessage
}

// HOEventMeasuredRIC is struct including UE ID and its eNB ID
type HOEventMeasuredRIC struct {
	Timestamp   time.Time
	Crnti       string
	DestPlmnID  string
	DestECID    string
	ElapsedTime int64
}

// ListHOEventMeasuredRIC is the list of HOEvent which is measured at ONOS RIC
var ListHOEventMeasuredRIC []HOEventMeasuredRIC
var mutexListHOEventMeasuredRIC sync.RWMutex

// NewSession creates a new southbound session controller.
func NewSession(ecgi sb.ECGI, endPoint sb.Endpoint) (*Session, error) {
	log.Infof("Creating Session for %v at %s", ecgi, endPoint)
	return &Session{Ecgi: ecgi, EndPoint: endPoint}, nil
}

// Run starts the southbound control loop.
func (s *Session) Run(tls topodevice.TlsConfig, creds topodevice.Credentials,
	controlUpdates chan sb.ControlUpdate, controlResponses chan sb.ControlResponse, telemetryUpdates chan sb.TelemetryMessage) {
	// Kick off a go routine that manages the connection to the simulator
	s.controlUpdates = controlUpdates
	s.controlResponses = controlResponses
	s.telemetryUpdates = telemetryUpdates
	go s.manageConnections(tls, creds)
}

// SendResponse sends the specified response on the control channel.
func (s *Session) SendResponse(response sb.ControlResponse) error {
	s.controlResponses <- response
	return nil
}

func (s *Session) manageConnections(tls topodevice.TlsConfig, creds topodevice.Credentials) {
	for {
		// Attempt to create connection to the simulator
		log.Infof("Connecting to simulator...%s with context", s.EndPoint)
		connection, err := southbound.Connect(context.Background(), string(s.EndPoint), tls.GetCert(), tls.GetKey())
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			s.manageConnection(connection)
		}
		time.Sleep(5 * time.Second)
	}
}

func (s *Session) manageConnection(connection *grpc.ClientConn) {
	// Offer the telemetry and control surfaces to the E2 devices
	s.client = sb.NewInterfaceServiceClient(connection)
	if s.client == nil {
		return
	}
	log.Infof("Connected to simulator on %s", s.EndPoint)
	// Setup coordination channel
	errors := make(chan error)

	go s.handleControl(errors)

	go s.handleTelemetry(errors)

	// Wait for the first error on the coordination channel.
	for i := 0; i < 2; i++ {
		<-errors
	}

	// Clean-up the connection, forcing other consumers to terminate and clean-up.
	_ = connection.Close()

	// FIXME: should be done separately
	//close(s.responses)
	//close(s.updates)

	close(errors)
	log.Info("Disconnected from simulator %s", s.EndPoint)
}

func (s *Session) handleControl(errors chan error) {

	cellConfigRequest := &sb.ControlResponse{
		MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		S: &sb.ControlResponse_CellConfigRequest{
			CellConfigRequest: &sb.CellConfigRequest{
				Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
			},
		},
	}

	stream, err := s.client.SendControl(context.Background())
	if err != nil {
		errors <- err
		return
	}

	waitc := make(chan error)
	go func() {
		for {
			update, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Infof("Got messageType %d from %s", update.MessageType, s.EndPoint)
			s.processControlUpdate(update)
		}
	}()

	if err := stream.Send(cellConfigRequest); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}

	for {
		response := <-s.controlResponses
		s.updateHOEventMeasuredRIC(response) // to monitor HO delay
		err := stream.Send(&response)
		if err != nil {
			waitc <- err
			return
		}
	}
}

func (s *Session) processControlUpdate(update *sb.ControlUpdate) {
	switch x := update.S.(type) {
	case *sb.ControlUpdate_CellConfigReport:
		log.Infof("%s CellConfigReport plmnid:%s, ecid:%s", s.EndPoint, x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
	case *sb.ControlUpdate_UEAdmissionRequest:
		log.Infof("%s UEAdmissionRequest plmnid:%s, ecid:%s, crnti:%s", s.EndPoint, x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti)
	case *sb.ControlUpdate_UEReleaseInd:
		log.Infof("%s UEReleaseInd plmnid:%s, ecid:%s, crnti:%s", s.EndPoint, x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)
	default:
		log.Fatalf("%s Control update has unexpected type %T", s.EndPoint, x)
	}
	s.controlUpdates <- *update
}

func (s *Session) handleTelemetry(errors chan error) {
	l2MeasConfig := &sb.L2MeasConfig{
		Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
	}

	stream, err := s.client.SendTelemetry(context.Background(), l2MeasConfig)
	if err != nil {
		errors <- err
		return
	}

	waitc := make(chan error)
	go func() {
		for {
			update, err := stream.Recv()
			if err != nil {
				close(waitc)
				return
			}
			log.Infof("%s telemetry messageType %d", s.EndPoint, update.MessageType)
			s.processTelemetryUpdate(update)
		}
	}()
}

func (s *Session) processTelemetryUpdate(update *sb.TelemetryMessage) {
	switch x := update.S.(type) {
	case *sb.TelemetryMessage_RadioMeasReportPerUE:
		s.addReceivedHOEventMeasuredRIC(update) // to monitor HO delay
	default:
		log.Fatalf("%s Telemetry update has unexpected type %T", s.EndPoint, x)
	}
	s.telemetryUpdates <- *update
}

func (s *Session) addReceivedHOEventMeasuredRIC(update *sb.TelemetryMessage) {
	servStationID := update.GetRadioMeasReportPerUE().GetEcgi()
	numNeighborCells := len(update.GetRadioMeasReportPerUE().GetRadioReportServCells())
	bestStationID := update.GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetEcgi()
	bestCQI := update.GetRadioMeasReportPerUE().GetRadioReportServCells()[0].GetCqiHist()[0]

	for i := 1; i < numNeighborCells; i++ {
		tmpCQI := update.GetRadioMeasReportPerUE().GetRadioReportServCells()[i].GetCqiHist()[0]
		if bestCQI < tmpCQI {
			bestStationID = update.GetRadioMeasReportPerUE().GetRadioReportServCells()[i].GetEcgi()
			bestCQI = tmpCQI
		}
	}

	if servStationID.GetEcid() != bestStationID.GetEcid() || servStationID.GetPlmnId() != bestStationID.GetPlmnId() {
		tmpHOEventMeasuredRIC := HOEventMeasuredRIC{
			Timestamp:  time.Now(),
			Crnti:      update.GetRadioMeasReportPerUE().GetCrnti(),
			DestPlmnID: bestStationID.GetPlmnId(),
			DestECID:   bestStationID.GetEcid(),
		}
		mutexListHOEventMeasuredRIC.Lock()
		ListHOEventMeasuredRIC = append(ListHOEventMeasuredRIC, tmpHOEventMeasuredRIC)
		mutexListHOEventMeasuredRIC.Unlock()
	}
}

func (s *Session) updateHOEventMeasuredRIC(ctrMsg sb.ControlResponse) {
	switch ctrMsg.S.(type) {
	case *sb.ControlResponse_HORequest:
		for i := 0; i < len(ListHOEventMeasuredRIC); i++ {
			mutexListHOEventMeasuredRIC.Lock()
			if ListHOEventMeasuredRIC[i].Crnti == ctrMsg.GetHORequest().GetCrnti() &&
				ListHOEventMeasuredRIC[i].DestECID == ctrMsg.GetHORequest().EcgiT.GetEcid() &&
				ListHOEventMeasuredRIC[i].DestPlmnID == ctrMsg.GetHORequest().EcgiT.GetPlmnId() &&
				ListHOEventMeasuredRIC[i].ElapsedTime == 0 {
				ListHOEventMeasuredRIC[i].ElapsedTime = time.Since(ListHOEventMeasuredRIC[i].Timestamp).Microseconds()
			}
			mutexListHOEventMeasuredRIC.Unlock()
		}
	default:
	}

}
