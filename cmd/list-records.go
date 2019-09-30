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
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apievents "github.com/oNaiPs/fyde-cli/client/device_events"
	"github.com/oNaiPs/fyde-cli/models"
)

// recordsListCmd represents the list command
var recordsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List records",
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
		params := apievents.NewListDeviceEventsParams()
		setFilter(cmd, params.SetUser, params.SetEventName)
		completePayload := []*models.DeviceEventListItem{}
		total := 0
		cutStart, cutEnd, err := forAllPages(cmd, params, func() (int, int64, error) {
			resp, err := global.Client.DeviceEvents.ListDeviceEvents(params, global.AuthWriter)
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
			"Name",
			"User",
			"Date",
		})
		tw.SetAllowedColumnLengths([]int{38, 30, 30, 30})

		for _, item := range completePayload {
			user := "?"
			if item.User != nil {
				user = item.User.Name
			}
			tw.AppendRow(table.Row{
				item.ID,
				item.Name,
				user,
				item.Date,
			})
		}

		result, err := renderListOutput(cmd, completePayload, tw, total)
		cmd.Println(result)
		return err
	},
}

func init() {
	recordsCmd.AddCommand(recordsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordsListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordsListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initPaginationFlags(recordsListCmd)
	initSortFlags(recordsListCmd)
	initFilterFlags(recordsListCmd,
		filterType{"event_name", "[]string"},
		filterType{"user", "string"})
	initOutputFlags(recordsListCmd)
}
