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
	"io"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/onos-ric/pkg/manager"
	e2 "github.com/onosproject/onos-ric/test/mocks/southbound"
	store "github.com/onosproject/onos-ric/test/mocks/store/indications"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func createServerConnection(t *testing.T) *grpc.ClientConn {
	lis = bufconn.Listen(1024 * 1024)
	s, err := NewService()
	assert.NoError(t, err)
	assert.NotNil(t, s)
	server := grpc.NewServer()
	s.Register(server)

	go func() {
		if err := server.Serve(lis); err != nil {
			assert.NoError(t, err, "Server exited with error: %v", err)
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn
}

func IDToEcid(ID uint32) string {
	return fmt.Sprintf("ECID-%d", ID)
}

func IDToPlmnid(ID uint32) string {
	return fmt.Sprintf("PLMNID-%d", ID)
}

func IDToCrnti(ID uint32) string {
	return fmt.Sprintf("CRNTI-%d", ID)
}

func checkEcgi(t *testing.T, ecgi *nb.ECGI, ID uint32) {
	assert.Equal(t, IDToEcid(ID), ecgi.Ecid)
	assert.Equal(t, IDToPlmnid(ID), ecgi.Plmnid)
}

func generateRicIndicationCellConfigReport(ID uint32) e2ap.RicIndication {
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_CELL_CONFIG_REPORT},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_CellConfigReport{CellConfigReport: &sb.CellConfigReport{
			Ecgi: &sb.ECGI{
				PlmnId: IDToPlmnid(ID),
				Ecid:   IDToEcid(ID),
			},
			Pci:                    (ID * 10) + 1,
			CandScells:             nil,
			EarfcnDl:               "",
			EarfcnUl:               "",
			RbsPerTtiDl:            (ID * 10) + 2,
			RbsPerTtiUl:            (ID * 10) + 3,
			NumTxAntenna:           (ID * 10) + 4,
			DuplexMode:             "",
			MaxNumConnectedUes:     (ID * 10) + 5,
			MaxNumConnectedBearers: (ID * 10) + 6,
			MaxNumUesSchedPerTtiDl: (ID * 10) + 7,
			MaxNumUesSchedPerTtiUl: (ID * 10) + 8,
			DlfsSchedEnable:        "",
		}}},
	}
}

func generateRicIndicationCellUEAdmissionRequest(ID uint32) e2ap.RicIndication {
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_UE_ADMISSION_REQUEST},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_UEAdmissionRequest{UEAdmissionRequest: &sb.UEAdmissionRequest{
			Ecgi: &sb.ECGI{
				PlmnId: IDToPlmnid(ID),
				Ecid:   IDToEcid(ID),
			},
		}}},
	}
}

func generateRicIndicationRadioMeasReportPerUE(ID uint32) e2ap.RicIndication {
	radioReportServCells := []*sb.RadioRepPerServCell{
		{
			Ecgi: &sb.ECGI{
				PlmnId: IDToPlmnid(ID + 1000),
				Ecid:   IDToEcid(ID + 1000),
			},
			CqiHist:       []uint32{1, 2, 3},
			RiHist:        nil,
			PuschSinrHist: nil,
			PucchSinrHist: nil,
		},
	}
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{MessageType: sb.MessageType_RADIO_MEAS_REPORT_PER_UE},
		Msg: &e2sm.RicIndicationMessage{S: &e2sm.RicIndicationMessage_RadioMeasReportPerUE{RadioMeasReportPerUE: &sb.RadioMeasReportPerUE{
			Ecgi: &sb.ECGI{
				PlmnId: IDToPlmnid(ID),
				Ecid:   IDToEcid(ID),
			},
			Crnti:                IDToCrnti(ID),
			RadioReportServCells: radioReportServCells,
		}}},
	}
}

const indicationsStoreItemsCount = 7

func createMockIndicationsStore(t *testing.T) {
	indicationsStore := store.NewMockStore(gomock.NewController(t))

	indicationsStore.EXPECT().List(gomock.Any()).DoAndReturn(
		func(ch chan<- e2ap.RicIndication) error {
			go func() {
				for i := 1; i <= indicationsStoreItemsCount; i++ {
					ch <- generateRicIndicationCellConfigReport(uint32(i))
					ch <- generateRicIndicationCellUEAdmissionRequest(uint32(i))
					ch <- generateRicIndicationRadioMeasReportPerUE(uint32(i))
				}
				close(ch)
			}()

			return nil
		})
	indicationsStore.EXPECT().Get(gomock.Any()).Return(nil, nil).AnyTimes()
	manager.InitializeManager(indicationsStore, nil, false)
}

func checkListStationsResponse(t *testing.T, info *nb.StationInfo, ID uint32) {
	checkEcgi(t, info.Ecgi, ID)
	assert.Equal(t, (ID*10)+5, info.MaxNumConnectedUes)
}

func Test_ListStations(t *testing.T) {
	createMockIndicationsStore(t)

	conn := createServerConnection(t)
	client := nb.NewC1InterfaceServiceClient(conn)

	request := nb.StationListRequest{Subscribe: false}
	stream, err := client.ListStations(context.Background(), &request)
	assert.NoError(t, err)
	assert.NotNil(t, stream)

	count := uint32(0)
	for {
		info, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			assert.NoError(t, err, "recv error")
		}
		count++
		checkListStationsResponse(t, info, count)
	}
	assert.Equal(t, uint32(indicationsStoreItemsCount), count, "Wrong number of stations received")
	assert.NoError(t, conn.Close())
}

func checkListStationLinksResponse(t *testing.T, info *nb.StationLinkInfo, ID uint32) {
	checkEcgi(t, info.Ecgi, ID)
}

func checkListUEsResponse(t *testing.T, info *nb.UEInfo, ID uint32) {
	checkEcgi(t, info.Ecgi, ID)
}

func checkListUELinksResponse(t *testing.T, info *nb.UELinkInfo, ID uint32) {
	checkEcgi(t, info.Ecgi, ID)
	assert.Equal(t, IDToCrnti(ID), info.Crnti)
	for _, qual := range info.ChannelQualities {
		assert.Equal(t, IDToEcid(ID+1000), qual.TargetEcgi.Ecid)
	}
}

func Test_ListStationLinks(t *testing.T) {
	createMockIndicationsStore(t)

	conn := createServerConnection(t)
	client := nb.NewC1InterfaceServiceClient(conn)

	request := nb.StationLinkListRequest{
		Subscribe: false,
	}
	stream, err := client.ListStationLinks(context.Background(), &request)
	assert.NoError(t, err)
	assert.NotNil(t, stream)

	count := uint32(0)
	for {
		info, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			assert.NoError(t, err, "recv error")
		}
		count++
		checkListStationLinksResponse(t, info, count)
	}
	assert.Equal(t, uint32(indicationsStoreItemsCount), count, "Wrong number of links received")
	assert.NoError(t, conn.Close())
}

func Test_ListUEs(t *testing.T) {
	createMockIndicationsStore(t)

	conn := createServerConnection(t)
	client := nb.NewC1InterfaceServiceClient(conn)

	request := nb.UEListRequest{
		Subscribe: false,
	}
	stream, err := client.ListUEs(context.Background(), &request)
	assert.NoError(t, err)
	assert.NotNil(t, stream)

	count := uint32(0)
	for {
		info, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			assert.NoError(t, err, "recv error")
		}
		count++
		checkListUEsResponse(t, info, count)
	}
	assert.Equal(t, uint32(indicationsStoreItemsCount), count, "Wrong number of UEs received")
	assert.NoError(t, conn.Close())
}

func Test_ListUELinks(t *testing.T) {
	createMockIndicationsStore(t)

	conn := createServerConnection(t)
	client := nb.NewC1InterfaceServiceClient(conn)

	request := nb.UELinkListRequest{
		Subscribe: false,
	}
	stream, err := client.ListUELinks(context.Background(), &request)
	assert.NoError(t, err)
	assert.NotNil(t, stream)

	count := uint32(0)
	for {
		info, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			assert.NoError(t, err, "recv error")
		}
		count++
		checkListUELinksResponse(t, info, count)
	}
	assert.Equal(t, uint32(indicationsStoreItemsCount), count, "Wrong number of UE Links received")
	assert.NoError(t, conn.Close())
}

func Test_RadioPower(t *testing.T) {
	testCases := []struct {
		description  string
		powerSetting nb.StationPowerOffset
		expectError  bool
	}{
		{description: "Power out of range", powerSetting: -1, expectError: true},
		{description: "StationPowerOffset_PA_DB_MINUS6", powerSetting: nb.StationPowerOffset_PA_DB_MINUS6, expectError: false},
		{description: "StationPowerOffset_PA_DB_MINUX4DOT77", powerSetting: nb.StationPowerOffset_PA_DB_MINUX4DOT77, expectError: false},
		{description: "StationPowerOffset_PA_DB_MINUS3", powerSetting: nb.StationPowerOffset_PA_DB_MINUS3, expectError: false},
		{description: "StationPowerOffset_PA_DB_MINUS1DOT77", powerSetting: nb.StationPowerOffset_PA_DB_MINUS1DOT77, expectError: false},
		{description: "StationPowerOffset_PA_DB_0", powerSetting: nb.StationPowerOffset_PA_DB_0, expectError: false},
		{description: "StationPowerOffset_PA_DB_1", powerSetting: nb.StationPowerOffset_PA_DB_1, expectError: false},
		{description: "StationPowerOffset_PA_DB_2", powerSetting: nb.StationPowerOffset_PA_DB_2, expectError: false},
		{description: "StationPowerOffset_PA_DB_3", powerSetting: nb.StationPowerOffset_PA_DB_3, expectError: false},
	}

	mockE2 := e2.NewMockE2(gomock.NewController(t))
	mockE2.EXPECT().RRMConfig(gomock.Any()).Return(nil).AnyTimes()

	plmnid := IDToPlmnid(1)
	ecid := IDToEcid(1)

	sbEcgi := sb.ECGI{PlmnId: plmnid, Ecid: ecid}

	manager.InitializeManager(nil, nil, false)
	manager.GetManager().SbSessions[sbEcgi] = mockE2

	conn := createServerConnection(t)
	client := nb.NewC1InterfaceServiceClient(conn)

	for _, test := range testCases {

		request := nb.RadioPowerRequest{
			Ecgi: &nb.ECGI{
				Plmnid: plmnid,
				Ecid:   ecid,
			},
			Offset: test.powerSetting,
		}

		response, err := client.SetRadioPower(context.Background(), &request)
		if test.expectError {
			assert.Error(t, err, "%s failed to return an error", test.description)
			assert.Nil(t, response, "%s failed to return a nil result", test.description)
		} else {
			assert.NoError(t, err, "%s returned an error", test.description)
			assert.NotNil(t, response, "%s returned a nil result", test.description)
		}
	}

	assert.NoError(t, conn.Close())
}
