// Package cmd implements fyde-cli commands
package cmd

/*
Copyright © 2019 Fyde, Inc. <hello@fyde.com>

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
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apidevices "github.com/fyde/fyde-cli/client/devices"
)

// deviceGetCmd represents the get command
var deviceGetCmd = &cobra.Command{
	Use:   "get [device ID]",
	Short: "Get device",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		if len(args) == 0 && !cmd.Flags().Changed("id") {
			return fmt.Errorf("missing device ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var deviceID string
		var err error
		if cmd.Flags().Changed("id") {
			deviceID, err = cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
		} else {
			deviceID = args[0]
		}

		params := apidevices.NewGetDeviceParams()
		params.SetID(strfmt.UUID(deviceID))

		resp, err := global.Client.Devices.GetDevice(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"User",
			"User Name",
			"OS",
			"Brand",
			"Model",
			"Enabled",
			"Status",
			"Failed security checks",
			"Total security checks",
		})

		failedChecks := 0
		for _, check := range resp.Payload.SecurityChecks {
			if check.Status != "passed" {
				failedChecks++
			}
		}

		tw.AppendRow(table.Row{
			resp.Payload.ID,
			resp.Payload.User.ID,
			resp.Payload.User.Name,
			resp.Payload.Os,
			resp.Payload.Brand,
			resp.Payload.HardwareModel,
			resp.Payload.Enabled,
			resp.Payload.Status,
			failedChecks,
			len(resp.Payload.SecurityChecks),
		})

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func init() {
	devicesCmd.AddCommand(deviceGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deviceGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(deviceGetCmd)
	deviceGetCmd.Flags().String("id", "", "id of device to get")
}
