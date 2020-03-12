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

package cli

import (
	"context"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/spf13/cobra"
)

// TestPlmnID - the default PlmnID
const TestPlmnID = "001001"

func getHandOverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handover <crnti> <src-ecgi> <dst-ecgi>",
		Short: "Trigger UE handover between two base stations",
		Args:  cobra.ExactArgs(3),
		RunE:  runHandOverCommand,
	}
	return cmd
}

func runHandOverCommand(cmd *cobra.Command, args []string) error {
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := nb.HandOverRequest{Crnti: args[0],
		SrcStation: &nb.ECGI{Ecid: args[1], Plmnid: TestPlmnID},
		DstStation: &nb.ECGI{Ecid: args[2], Plmnid: TestPlmnID}}
	client := nb.NewC1InterfaceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.TriggerHandOver(ctx, &request)
	if err != nil {
		log.Errorf("handover error %s", err.Error())
		return err
	}
	if !response.Success {
		log.Error("Handover failed")
	}

	return nil
}
