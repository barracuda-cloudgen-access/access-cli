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

	apievents "github.com/fyde/fyde-cli/client/device_events"
)

// recordGetCmd represents the get command
var recordGetCmd = &cobra.Command{
	Use:   "get [record ID] [record date]",
	Short: "Get record",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		if len(args) < 2 {
			return fmt.Errorf("missing arguments. Usage: records get [ID] [date]")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		params := apievents.NewGetDeviceEventParams()
		params.SetID(args[0])

		date, err := strfmt.ParseDateTime(args[1])
		if err != nil {
			return fmt.Errorf("invalid date argument")
		}
		params.SetDate(strfmt.DateTime(date))

		resp, err := global.Client.DeviceEvents.GetDeviceEvent(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"Name",
			"User",
			"Date",
		})
		tw.SetAllowedColumnLengths([]int{38, 30, 30, 30})

		user := "?"
		if resp.Payload.User != nil {
			user = resp.Payload.User.Name
		}

		tw.AppendRow(table.Row{
			resp.Payload.ID,
			resp.Payload.Name,
			user,
			resp.Payload.Date.Utc,
		})

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func init() {
	recordsCmd.AddCommand(recordGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(recordGetCmd)
}
