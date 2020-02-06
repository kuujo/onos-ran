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
	"fmt"
	"github.com/onosproject/onos-ran/api/nb"
	"github.com/spf13/cobra"
	"io"
	log "k8s.io/klog"
	"text/tabwriter"
	"time"
)

func getGetUesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Get UEs",
		RunE:  runUesCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func getWatchUesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Watch UEs",
		RunE:  runUesCommand,
	}
	cmd.SetArgs([]string{_subscribe})
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runUesCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	var subscribe bool
	if len(args) == 1 && args[0] == _subscribe {
		subscribe = true
	}

	conn, err := getConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := GetOutput()

	request := nb.UEListRequest{}
	if subscribe {
		// TODO: indicate watch semantics in the request
		Output("Watching list of UEs\n")
	}

	client := nb.NewC1InterfaceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stream, err := client.ListUEs(ctx, &request)
	if err != nil {
		log.Error("list error ", err)
		return err
	}
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		fmt.Fprintln(writer, "IMSI")
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error("rcv error ", err)
			return err
		}

		fmt.Fprintln(writer, response.Imsi)
	}
	writer.Flush()

	return nil
}
