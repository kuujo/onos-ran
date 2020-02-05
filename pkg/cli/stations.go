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
	"github.com/spf13/cobra"
)

const _subscribe = "subscribe"

func getGetStationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stations",
		Short: "Get Stations",
		RunE:  runStationsCommand,
	}
	return cmd
}

func getWatchStationsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stations",
		Short: "Watch Stations",
		RunE:  runStationsCommand,
	}
	cmd.SetArgs([]string{_subscribe})
	return cmd
}

func runStationsCommand(cmd *cobra.Command, args []string) error {
	var subscribe bool
	if len(args) == 1 && args[0] == _subscribe {
		subscribe = true
	}

	if !subscribe {
		Output("Getting list of Stations\n")
	} else {
		Output("Watching list of Stations\n")
	}

	conn, err := getConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()

	return nil
}
