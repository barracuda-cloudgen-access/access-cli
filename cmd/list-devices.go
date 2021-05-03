// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2020 Barracuda Networks, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	"github.com/barracuda-cloudgen-access/access-cli/client/devices"
)

// devicesListCmd represents the list command
var devicesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List devices",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		params := devices.NewListDevicesParams()
		//setSort(cmd, params) // TODO re-enable when/if devices supports sort
		completePayload := []*devices.ListDevicesOKBodyItems0{}
		total := 0
		cutStart, cutEnd, err := forAllPages(cmd, params, func() (int, int64, error) {
			resp, err := global.Client.Devices.ListDevices(params, global.AuthWriter)
			if err != nil {
				return 0, 0, err
			}
			completePayload = append(completePayload, resp.Payload...)
			total = int(resp.Total)
			return len(resp.Payload), resp.Total, err
		})
		if err != nil {
			return processErrorResponse(err)
		}
		completePayload = completePayload[cutStart:cutEnd]

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"User",
			"User Name",
			"OS",
			"Brand",
			"Model",
			"Status",
			"Failed security checks",
			"Total security checks",
		})

		for _, item := range completePayload {
			failedChecks := 0
			for _, check := range item.SecurityChecks {
				if check.Status != "passed" {
					failedChecks++
				}
			}
			tw.AppendRow(table.Row{
				item.ID,
				item.User.ID,
				item.User.Name,
				item.Os,
				item.Brand,
				item.HardwareModel,
				item.Status,
				failedChecks,
				len(item.SecurityChecks),
			})
		}

		return printListOutputAndError(cmd, completePayload, tw, total, err)
	},
}

func init() {
	devicesCmd.AddCommand(devicesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// devicesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// devicesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initPaginationFlags(devicesListCmd)
	//initSortFlags(devicesListCmd) // TODO re-enable when/if devices supports sort
	initOutputFlags(devicesListCmd)
	devicesListCmd.Flags().StringP("filter", "f", "", "filter devices")
}
