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
	"io"
	"text/tabwriter"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/spf13/cobra"
)

func getGetUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Get UEs",
		RunE:  runUEsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("ecgi", "", "optional station identifier")
	cmd.Flags().String("crnti", "", "optional UE identifier in the local station context")
	return cmd
}

func getWatchUEsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ues",
		Short: "Get UEs",
		RunE:  runUEsCommand,
	}
	cmd.Flags().String("ecgi", "", "optional station identifier")
	cmd.Flags().String("crnti", "", "optional UE identifier in the local station context")
	return cmd
}

func runUEsCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	subscribe, _ := cmd.Flags().GetBool(_subscribe)

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()

	request := nb.UEListRequest{Subscribe: subscribe}
	if subscribe {
		// TODO: indicate watch semantics in the request
		cli.Output("Watching list of UEs\n")
	}

	// Populate optional ECGI qualifier
	ecgi := getECGI(cmd)
	if ecgi != "" {
		request.Ecgi = &nb.ECGI{Ecid: ecgi}
	}

	// Populate optional CRNTI qualifier
	crnti := getCRNTI(cmd)
	if ecgi != "" {
		request.Crnti = crnti
	}

	client := nb.NewC1InterfaceServiceClient(conn)

	stream, err := client.ListUEs(context.Background(), &request)
	if err != nil {
		return err
	}

	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, "PLMNID\tECID\tC-RNTI\tIMSI")
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("recv error")
			return err
		}

		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", response.Ecgi.GetPlmnid(), response.Ecgi.GetEcid(), response.GetCrnti(), response.GetImsi())
	}

	_ = writer.Flush()

	return nil
}
