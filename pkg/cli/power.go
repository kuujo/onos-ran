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
	"github.com/onosproject/onos-ran/api/nb"
	"github.com/spf13/cobra"
	log "k8s.io/klog"
	"time"
)

func getSetPowerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power <ecgi> <power-offset>",
		Short: "Set transmit power of a base station",
		Args:  cobra.ExactArgs(2),
		RunE:  runSetPowerCommand,
	}
	return cmd
}

func runSetPowerCommand(cmd *cobra.Command, args []string) error {
	conn, err := getConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	offset, ok := nb.StationPowerOffset_value[args[1]]
	if !ok {
		log.Error("Unsupported power offset ", args[1])
		return err
	}

	request := nb.RadioPowerRequest{Ecgi: &nb.ECGI{Ecid: args[0]}, Offset: nb.StationPowerOffset(offset)}
	client := nb.NewC1InterfaceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	response, err := client.SetRadioPower(ctx, &request)
	if err != nil {
		log.Error("set power error ", err)
		return err
	}
	if !response.Success {
		log.Error("Failed to set power")
	}

	return nil
}
