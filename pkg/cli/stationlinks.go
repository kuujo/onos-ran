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
	"time"

	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/spf13/cobra"
)

func getGetStationLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stationlinks",
		Short: "Get Station Links",
		RunE:  runStationLinksCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("ecgi", "", "optional station identifier")
	return cmd
}

func getWatchStationLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stationlinks",
		Short: "Watch Station Links",
		RunE:  runStationLinksCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runStationLinksCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	subscribe, _ := cmd.Flags().GetBool(_subscribe)

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()

	request := nb.StationLinkListRequest{Subscribe: subscribe}

	// Populate optional ECGI qualifier
	ecgi := getECGI(cmd)
	if ecgi != "" {
		request.Ecgi = &nb.ECGI{Ecid: ecgi}
	}

	client := nb.NewC1InterfaceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	stream, err := client.ListStationLinks(ctx, &request)
	if err != nil {
		return err
	}
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, "ECID\tNEIGHBOR ECIDs")
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("rcv error")
			return err
		}

		var neighbours string
		for i, s := range response.NeighborECGI {
			if i > 0 {
				neighbours = neighbours + ", " + s.Ecid
			} else {
				neighbours = s.Ecid
			}
		}
		_, _ = fmt.Fprintf(writer, "%s\t%s\n", response.Ecgi.Ecid, neighbours)
	}
	_ = writer.Flush()

	return nil
}
