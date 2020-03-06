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
	"sync"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-ric/api/sb"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("southbound")

// Sessions is responsible for managing connections to and interactions with the RAN southbound.
type Sessions struct {
	Simulator *string
	client    sb.InterfaceServiceClient

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

// NewSessions creates a new southbound sessions controller.
func NewSessions() (*Sessions, error) {
	log.Info("Creating Sessions")
	return &Sessions{}, nil
}

// Run starts the southbound control loop.
func (m *Sessions) Run(controlUpdates chan sb.ControlUpdate, controlResponses chan sb.ControlResponse, telemetryUpdates chan sb.TelemetryMessage) {
	// Kick off a go routine that manages the connection to the simulator
	m.controlUpdates = controlUpdates
	m.controlResponses = controlResponses
	m.telemetryUpdates = telemetryUpdates
	go m.manageConnections()
}

// SendResponse sends the specified response on the control channel.
func (m *Sessions) SendResponse(response sb.ControlResponse) error {
	m.controlResponses <- response
	return nil
}

func (m *Sessions) manageConnections() {
	for {
		// Attempt to create connection to the simulator

		log.Infof("Connecting to simulator...%s with context", *m.Simulator)
		connection, err := getConnection(*m.Simulator)
		if err == nil {
			// If successful, manage this connection and don't return until it is
			// no longer valid and all related resources have been properly cleaned-up.
			m.manageConnection(connection)
		}
		time.Sleep(5 * time.Second)
	}
}

func (m *Sessions) manageConnection(connection *grpc.ClientConn) {
	// Offer the telemetry and control surfaces to the E2 devices
	m.client = sb.NewInterfaceServiceClient(connection)
	if m.client == nil {
		return
	}

	log.Info("Connected to simulator")
	// Setup coordination channel
	errors := make(chan error)

	go m.handleControl(errors)

	go m.handleTelemetry(errors)

	// Wait for the first error on the coordination channel.
	for i := 0; i < 2; i++ {
		<-errors
	}

	// Clean-up the connection, forcing other consumers to terminate and clean-up.
	_ = connection.Close()

	// FIXME: should be done separately
	//close(m.responses)
	//close(m.updates)

	close(errors)
	log.Info("Disconnected from simulator")
}

func (m *Sessions) handleControl(errors chan error) {

	cellConfigRequest := &sb.ControlResponse{
		MessageType: sb.MessageType_CELL_CONFIG_REQUEST,
		S: &sb.ControlResponse_CellConfigRequest{
			CellConfigRequest: &sb.CellConfigRequest{
				Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
			},
		},
	}

	stream, err := m.client.SendControl(context.Background())
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
			log.Infof("Got messageType %d", update.MessageType)
			m.processControlUpdate(update)
		}
	}()

	if err := stream.Send(cellConfigRequest); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}

	for {
		response := <-m.controlResponses
		m.updateHOEventMeasuredRIC(response) // to monitor HO delay
		err := stream.Send(&response)
		if err != nil {
			waitc <- err
			return
		}
	}
}

func (m *Sessions) processControlUpdate(update *sb.ControlUpdate) {
	switch x := update.S.(type) {
	case *sb.ControlUpdate_CellConfigReport:
		log.Infof("plmnid:%s, ecid:%s", x.CellConfigReport.Ecgi.PlmnId, x.CellConfigReport.Ecgi.Ecid)
	case *sb.ControlUpdate_UEAdmissionRequest:
		log.Infof("plmnid:%s, ecid:%s, crnti:%s", x.UEAdmissionRequest.Ecgi.PlmnId, x.UEAdmissionRequest.Ecgi.Ecid, x.UEAdmissionRequest.Crnti)
	case *sb.ControlUpdate_UEReleaseInd:
		log.Infof("plmnid:%s, ecid:%s, crnti:%s", x.UEReleaseInd.Ecgi.PlmnId, x.UEReleaseInd.Ecgi.Ecid, x.UEReleaseInd.Crnti)
	default:
		log.Fatalf("Control update has unexpected type %T", x)
	}
	m.controlUpdates <- *update
}

func (m *Sessions) handleTelemetry(errors chan error) {
	l2MeasConfig := &sb.L2MeasConfig{
		Ecgi: &sb.ECGI{PlmnId: "test", Ecid: "test"},
	}

	stream, err := m.client.SendTelemetry(context.Background(), l2MeasConfig)
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
			log.Infof("Got telemetry messageType %d", update.MessageType)
			m.processTelemetryUpdate(update)
		}
	}()
}

func (m *Sessions) processTelemetryUpdate(update *sb.TelemetryMessage) {
	switch x := update.S.(type) {
	case *sb.TelemetryMessage_RadioMeasReportPerUE:
		m.addReceivedHOEventMeasuredRIC(update) // to monitor HO delay
	default:
		log.Fatalf("Telemetry update has unexpected type %T", x)
	}
	m.telemetryUpdates <- *update
}

func (m *Sessions) addReceivedHOEventMeasuredRIC(update *sb.TelemetryMessage) {
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

func (m *Sessions) updateHOEventMeasuredRIC(ctrMsg sb.ControlResponse) {
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
