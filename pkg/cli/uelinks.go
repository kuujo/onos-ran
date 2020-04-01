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

const ueLinksCellsTemplate = "{{ printf \"IMSIs           UEs                      \" }}" +
	"{{range . }}" +
	"{{ printf \"%-10s\" . }}" +
	"{{end}}\n"

const ueLinksTemplate = "{{ printf \"%-16s\" .Imsi }}" +
	"{{printf \"%-25s\" .Ue}}" +
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

var tmplCellsList, _ = template.New("uelinkscells").Parse(ueLinksCellsTemplate)
var tmplUesList, _ = template.New("uelinksues").Parse(ueLinksTemplate)

func getGetUeLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uelinks",
		Short: "Get UE Links",
		RunE:  runUeLinksCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	cmd.Flags().String("ecgi", "", "optional station identifier")
	cmd.Flags().String("crnti", "", "optional UE identifier in the local station context")
	cmd.Flags().Bool("sortimsi", false, "Sort by UE Imsi instead of Cell")
	return cmd
}

func getWatchUeLinksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uelinks",
		Short: "Watch UE Links",
		RunE:  runUeLinksCommand,
	}
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
	sortImsi, _ := cmd.Flags().GetBool("sortimsi")

	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	request := nb.UELinkListRequest{Subscribe: subscribe}
	if subscribe {
		// TODO: indicate watch semantics in the request
		cli.Output("Watching list of UE links\n")
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
		return err
	}
	towerKeys := make([]string, 0)
	qualityMap := make(map[string][]*nb.ChannelQuality)
	imsiMap := make(map[string]string)
	imsiInvMap := make(map[string]string)
	imsiKeys := make([]string, 0)
	servingMap := make(map[string]string)
	qualKeys := make([]string, 0)
	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			cli.Output("recv error")
			return err
		}
		ecids := make([]string, 0)
		for _, cq := range response.ChannelQualities {
			ecids = append(ecids, cq.GetTargetEcgi().GetEcid())
		}
		ecids = append(ecids, response.GetEcgi().GetEcid())
		formerLen := len(towerKeys)
		towerKeys = addTowerKeys(towerKeys, ecids...)
		qualID := fmt.Sprintf("%s:%s:%s", response.Ecgi.Plmnid, response.Ecgi.Ecid, response.Crnti)
		if subscribe {
			if len(towerKeys) > formerLen {
				_ = tmplCellsList.Execute(cli.GetOutput(), towerKeys)
			}
			printQualMap(qualID, response.Imsi, response.Ecgi.Ecid, response.ChannelQualities, towerKeys, tmplUesList, true)
		} else {
			qualityMap[qualID] = response.ChannelQualities
			imsiMap[qualID] = response.Imsi
			imsiInvMap[response.Imsi] = qualID
			qualKeys = append(qualKeys, qualID)
			imsiKeys = append(imsiKeys, response.GetImsi())
			servingMap[qualID] = response.GetEcgi().GetEcid()
		}
	}
	if !subscribe {
		if !noHeaders {
			cli.Output("Cqi for UE-Cell links:    UEs: %4d      Cells %4d   (*=serving cell)\n", len(qualityMap), len(towerKeys))
		}
		_ = tmplCellsList.Execute(cli.GetOutput(), towerKeys)
		if !sortImsi {
			sort.Slice(qualKeys, func(i, j int) bool {
				return qualKeys[i] < qualKeys[j]
			})
		} else {
			sort.Slice(imsiKeys, func(i, j int) bool {
				return imsiKeys[i] < imsiKeys[j]
			})
			qualKeys = make([]string, 0)
			for _, imsi := range imsiKeys {
				if qualID, ok := imsiInvMap[imsi]; ok {
					qualKeys = append(qualKeys, qualID)
				}
			}
		}
		for _, qualID := range qualKeys {
			printQualMap(qualID, imsiMap[qualID], servingMap[qualID], qualityMap[qualID], towerKeys, tmplUesList, false)
		}
	}
	return nil
}

func printQualMap(qualid string, imsi string, serving string, qualities []*nb.ChannelQuality, towerKeys []string, tmplUesList *template.Template, printTime bool) {
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
	qualSet := QualSet{Ue: qualid, Imsi: imsi, UeQual: qualTable, PrintTime: printTime, Time: time.Now().Format("15:04:05.0")}
	_ = tmplUesList.Execute(cli.GetOutput(), qualSet)
}

func addTowerKeys(existing []string, towerIds ...string) []string {
	towers := make(map[string]interface{})
	for _, t := range existing {
		towers[t] = struct{}{}
	}
	for _, t := range towerIds { //Prevents duplicates
		towers[t] = struct{}{}
	}
	towerKeys := make([]string, len(towers))
	var i = 0
	for t := range towers { //Back to array
		towerKeys[i] = t
		i++
	}
	sort.Slice(towerKeys, func(i, j int) bool {
		return towerKeys[i] < towerKeys[j]
	})
	return towerKeys
}
