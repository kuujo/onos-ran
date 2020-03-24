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
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/onosproject/onos-ric/api/nb"
	"github.com/spf13/cobra"
	"io"
	"sort"
	"text/template"
	"time"
)

// QualSet - used for the Go template/text output
type QualSet struct {
	Ue        string
	Imsi      string
	UeQual    []string
	PrintTime bool
	Time      string
}

const ueLinksCellsTemplate = "{{ printf \"IMSIs           UEs       \" }}" +
	"{{range . }}" +
	"{{ printf \"%-10s\" . }}" +
	"{{end}}\n"

const ueLinksTemplate = "{{ printf \"%-16s\" .Imsi }}" +
	"{{printf \"%-10s\" .Ue}}" +
	"{{ range $idx, $qual := .UeQual }}" +
	"{{ if $qual}}" +
	"{{printf \"%-10s\" $qual}}" +
	"{{ else}}" +
	"{{printf \"%-10s\" \"\"}}" +
	"{{ end}}" +
	"{{end}}" +
	"{{ if .PrintTime }}" +
	"{{printf \"%-20s\" .Time }}" +
	"{{end}}\n"

func getGetUeLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uelinks",
		Short: "Get UE Links",
		RunE:  runUeLinksCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("ecgi", "", "optional station identifier")
	cmd.Flags().String("crnti", "", "optional UE identifier in the local station context")
	return cmd
}

func getWatchUeLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uelinks",
		Short: "Watch UE Links",
		RunE:  runUeLinksCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("ecgi", "", "optional station identifier")
	cmd.Flags().String("crnti", "", "optional UE identifier in the local station context")
	return cmd
}

func getCRNTI(cmd *cobra.Command) string {
	crnti, _ := cmd.Flags().GetString("crnti")
	return crnti
}

func runUeLinksCommand(cmd *cobra.Command, args []string) error {
	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	subscribe, _ := cmd.Flags().GetBool(_subscribe)

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	tmplCellsList, _ := template.New("uelinkscells").Parse(ueLinksCellsTemplate)
	tmplUesList, _ := template.New("uelinksues").Parse(ueLinksTemplate)

	request := nb.UELinkListRequest{Subscribe: subscribe}
	if subscribe {
		// TODO: indicate watch semantics in the request
		Output("Watching list of UE links\n")
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

	stream, err := client.ListUELinks(context.Background(), &request)
	if err != nil {
		Output("list error %s", err.Error())
		return err
	}
	towers := make(map[string]interface{})
	var towerKeys []string
	qualityMap := make(map[string][]*nb.ChannelQuality)
	imsiMap := make(map[string]string)
	if !noHeaders {
		Output("                        Cell sites\n")
	}
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			Output("rcv error %s", err.Error())
			return err
		}
		if _, ok := towers[response.Ecgi.Ecid]; !ok {
			towers[response.Ecgi.Ecid] = struct{}{}
			towerKeys = make([]string, len(towers))
			var i = 0
			for k := range towers {
				towerKeys[i] = k
				i++
			}
			sort.Slice(towerKeys, func(i, j int) bool {
				return towerKeys[i] < towerKeys[j]
			})
			if subscribe {
				_ = tmplCellsList.Execute(GetOutput(), towerKeys)
			}
		}
		qualityMap[response.Crnti] = response.ChannelQualities
		imsiMap[response.Crnti] = response.Imsi
		if subscribe {
			printQualMap(response.Crnti, response.Imsi, response.Ecgi.Ecid, response.ChannelQualities, towerKeys, tmplUesList, true)
		}
	}
	if !subscribe {
		_ = tmplCellsList.Execute(GetOutput(), towerKeys)
		for k, qualities := range qualityMap {
			printQualMap(k, imsiMap[k], "", qualities, towerKeys, tmplUesList, false)
		}
	}
	return nil
}

func printQualMap(crnti string, imsi string, serving string, qualities []*nb.ChannelQuality, towerKeys []string, tmplUesList *template.Template, printTime bool) {
	qualTable := make([]string, len(towerKeys))
	for _, q := range qualities {
		for i, t := range towerKeys {
			if t == q.TargetEcgi.Ecid {
				qualTable[i] = fmt.Sprintf("%d", q.CqiHist)
				if t == serving {
					qualTable[i] = fmt.Sprintf("*%d", q.CqiHist)
				}
				break
			}
		}
	}
	qualSet := QualSet{Ue: crnti, Imsi: imsi, UeQual: qualTable, PrintTime: printTime, Time: time.Now().Format("15:04:05.0")}
	_ = tmplUesList.Execute(GetOutput(), qualSet)
}
