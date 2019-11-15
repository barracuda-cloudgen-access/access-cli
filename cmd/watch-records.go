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
	"sort"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apievents "github.com/fyde/fyde-cli/client/device_events"
	"github.com/fyde/fyde-cli/models"
)

// recordsWatchCmd represents the watch command
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

		refreshPeriod, err := cmd.Flags().GetInt("refresh-period")
		if err != nil {
			return err
		}

		if refreshPeriod < 1 {
			return fmt.Errorf("invalid refresh period of %d seconds", refreshPeriod)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		recordChan := make(chan *models.DeviceEventListItem)

		outputFormat, _ := cmd.Flags().GetString("output")
		detailedEvents, _ := cmd.Flags().GetBool("detailed-info")
		detailedEvents = detailedEvents && (outputFormat == "json" || outputFormat == "json-pretty")

		refreshPeriod, _ := cmd.Flags().GetInt("refresh-period")

		var innerError error

		// launch producer thread
		go func() {
			lastSeenID := ""
			var fetchStart time.Time

			// newRecords contains the records created since last update check
			// most recent are always first
			newRecords := []*models.DeviceEventListItem{}
			for page := int64(1); ; page++ {
				params := apievents.NewListDeviceEventsParams()
				params.SetPage(&page)
				setFilter(cmd, params.SetEventName, params.SetUserID)
				resp, err := global.Client.DeviceEvents.ListDeviceEvents(params, global.AuthWriter)
				if err != nil {
					innerError = err
					close(recordChan)
					return
				}

				// workaround for the fact that the server returns records
				// ordered by descending timestamp, but the sorting is not
				// stable, the order of events with the same timestamp changes.
				// this breaks the check using lastSeenID
				// sort the events ourselves
				// order by timestamp desc, id asc
				sort.Slice(resp.Payload, func(iIdx, jIdx int) bool {
					i := resp.Payload[iIdx]
					j := resp.Payload[jIdx]
					if time.Time(i.Date).Equal(time.Time(j.Date)) {
						return i.ID < j.ID
					}
					// use After for descending order
					return time.Time(i.Date).After(time.Time(j.Date))
				})
				// end of workaround

				// see if and where the last seen record is in this response
				// (if not, we'll need to fetch another page)
				lastSeenIdx := -1
				if lastSeenID == "" {
					lastSeenIdx = len(resp.Payload)
				} else {
					for i, record := range resp.Payload {
						if record.ID == lastSeenID {
							lastSeenIdx = i
						}
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
					waitFor := time.Duration(refreshPeriod)*time.Second - time.Since(fetchStart)
					if waitFor > 0 {
						time.Sleep(waitFor)
					}
					fetchStart = time.Now()
				} else {
					// collect all records from this page and move on to the next page
					// (we loop until we find the last record we saw)
					newRecords = append(resp.Payload, newRecords...)
				}
			}
		}()

		// we are the consumer thread
		isFirst := true
		for record := range recordChan {
			tw := table.NewWriter()
			tw.Style().Format.Header = text.FormatDefault

			if isFirst {
				tw.AppendHeader(table.Row{
					"ID",
					"Name",
					"User",
					"Date",
				})
			}

			// fix column width so it looks consistent across "tables"
			tw.SetColumnConfigs([]table.ColumnConfig{
				table.ColumnConfig{
					Number:   1,
					WidthMin: 38,
					WidthMax: 38,
				},
				table.ColumnConfig{
					Number:   2,
					WidthMin: 25,
					WidthMax: 25,
				},
				table.ColumnConfig{
					Number:   3,
					WidthMin: 20,
					WidthMax: 20,
				},
				table.ColumnConfig{
					Number:   4,
					WidthMin: 24,
					WidthMax: 24,
				},
			})

			user := "Unknown"
			var toRender interface{}
			if detailedEvents {
				params := apievents.NewGetDeviceEventParams()
				params.SetID(record.ID)
				params.SetDate(record.Date)

				resp, err := global.Client.DeviceEvents.GetDeviceEvent(params, global.AuthWriter)
				if err != nil {
					return processErrorResponse(err)
				}
				if resp.Payload.User != nil {
					user = resp.Payload.User.Name
				}
				tw.AppendRow(table.Row{
					resp.Payload.ID,
					resp.Payload.Name,
					user,
					resp.Payload.Date.Utc,
				})
				toRender = resp.Payload
			} else {
				if record.User != nil {
					user = record.User.Name
				}
				tw.AppendRow(table.Row{
					record.ID,
					record.Name,
					user,
					record.Date,
				})
				toRender = record
			}

			isTable, result, err := renderWatchOutput(cmd, toRender, tw)
			if err != nil {
				return err
			}

			if isTable {
				if !isFirst {
					//remove top border
					result = result[strings.Index(result, "\n")+1 : len(result)]
				}
				result = result[0:strings.LastIndex(result, "\n")]
			}

			cmd.Println(result)
			isFirst = false
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

	initFilterFlags(recordsWatchCmd,
		filterType{"event_name", "[]string"},
		filterType{"user-id", "int"})
	initOutputFlags(recordsWatchCmd)

	recordsWatchCmd.Flags().IntP("refresh-period", "r", 60, "period, in seconds, at which to check for new events")
	recordsWatchCmd.Flags().BoolP("detailed-info", "d", false, "show detailed info for each record (slower, only for JSON output)")
}
