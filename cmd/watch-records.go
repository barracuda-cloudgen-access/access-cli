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
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apievents "github.com/oNaiPs/fyde-cli/client/device_events"
	"github.com/oNaiPs/fyde-cli/models"
)

// recordsWatchCmd represents the list command
var recordsWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch records as they are created",
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
		recordChan := make(chan *models.DeviceEventListItem)

		var innerError error

		// launch producer thread
		go func() {
			lastSeenID := ""

			params := apievents.NewListDeviceEventsParams()
			setFilter(cmd, params.SetUser, params.SetEventName)
			resp, err := global.Client.DeviceEvents.ListDeviceEvents(params, global.AuthWriter)
			if err != nil {
				innerError = err
				close(recordChan)
				return
			}

			if len(resp.Payload) == 0 {
				innerError = fmt.Errorf("No records")
				close(recordChan)
				return
			}

			// most recent always comes first
			lastSeenID = resp.Payload[0].ID
			for i := len(resp.Payload) - 1; i >= 0; i-- {
				recordChan <- resp.Payload[i]
			}

			// newRecords contains the records created since last update check
			// most recent are always first
			newRecords := []*models.DeviceEventListItem{}
			for page := int64(1); ; page++ {
				params := apievents.NewListDeviceEventsParams()
				params.SetPage(&page)
				setFilter(cmd, params.SetUser, params.SetEventName)
				resp, err := global.Client.DeviceEvents.ListDeviceEvents(params, global.AuthWriter)
				if err != nil {
					innerError = err
					close(recordChan)
					return
				}

				// see if and where the last seen record is in this response
				// (if not, we'll need to fetch another page)
				lastSeenIdx := -1
				for i, record := range resp.Payload {
					if record.ID == lastSeenID {
						lastSeenIdx = i
					}
				}
				if lastSeenIdx != -1 {
					newRecords = append(resp.Payload[0:lastSeenIdx], newRecords...)
					// most recent always comes first in responses,
					// but we want to produce records from old to new
					for i := len(newRecords) - 1; i >= 0; i-- {
						recordChan <- newRecords[i]
						if i == 0 {
							lastSeenID = newRecords[0].ID
						}
					}

					// reset state and wait:
					page = 0 // will be set to 1 once we loop
					newRecords = []*models.DeviceEventListItem{}
					time.Sleep(5 * time.Second)
				} else {
					// collect all records from this page and move on to the next page
					// (we loop until we find the last record we saw)
					newRecords = append(resp.Payload, newRecords...)
				}
			}
		}()

		// we are the consumer thread
		for record := range recordChan {
			tw := table.NewWriter()
			tw.Style().Format.Header = text.FormatDefault
			tw.Style().Options.DrawBorder = false
			/*tw.AppendHeader(table.Row{
				"ID",
				"Name",
				"User",
				"Date",
			})*/
			user := "?"
			if record.User != nil {
				user = record.User.Name
			}
			tw.AppendRow(table.Row{
				record.ID,
				record.Name,
				user,
				record.Date,
			})
			tw.SetAllowedColumnLengths([]int{38, 30, 30, 30})
			cmd.Println(tw.Render())
		}
		if innerError != nil {
			return processErrorResponse(innerError)
		}
		return nil
	},
}

func init() {
	recordsCmd.AddCommand(recordsWatchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recordsWatchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recordsWatchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initPaginationFlags(recordsWatchCmd)
	initSortFlags(recordsWatchCmd)
	initFilterFlags(recordsWatchCmd,
		filterType{"event_name", "[]string"},
		filterType{"user", "string"})
	initOutputFlags(recordsWatchCmd)
}
