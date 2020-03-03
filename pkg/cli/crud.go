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

import "github.com/spf13/cobra"

func getGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get {stations|stationlinks|uelinks} [args]",
		Short: "Get RIC resources",
	}
	cmd.AddCommand(getGetStationsCommand())
	cmd.AddCommand(getGetStationLinksCommand())
	cmd.AddCommand(getGetUeLinksCommand())
	return cmd
}

func getWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch {stations|stationlinks|uelinks} [args]",
		Short: "Watch for changes to a RIC resource type",
	}
	cmd.AddCommand(getWatchStationsCommand())
	cmd.AddCommand(getWatchStationLinksCommand())
	cmd.AddCommand(getWatchUeLinksCommand())
	return cmd
}

func getSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {power} [args]",
		Short: "Set RIC resource parameters",
	}
	cmd.AddCommand(getSetPowerCommand())
	return cmd
}
