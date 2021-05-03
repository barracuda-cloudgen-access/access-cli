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
	"fmt"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apigroups "github.com/barracuda-cloudgen-access/access-cli/client/groups"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// groupGetCmd represents the get command
var groupGetCmd = &cobra.Command{
	Use:   "get [group ID]",
	Short: "Get group",
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
			return fmt.Errorf("missing group ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var groupID int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			groupID = int64(d)
		} else {
			groupID, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		params := apigroups.NewGetGroupParams()
		params.SetID(groupID)

		resp, err := global.Client.Groups.GetGroup(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := groupBuildTableWriter()
		groupTableWriterAppendFromSingle(tw, resp.Payload.Group, len(resp.Payload.Users))

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func groupBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Description",
		"Enrolled users",
		"Total users",
	})
	tw.SetAlign([]text.Align{
		text.AlignRight,
		text.AlignLeft,
		text.AlignLeft,
		text.AlignLeft,
		text.AlignLeft,
	})
	tw.SetAllowedColumnLengths([]int{15, 30, 30, 15, 15})
	return tw
}

func groupTableWriterAppendFromSingle(tw table.Writer, group models.Group, users int) {
	tw.AppendRow(table.Row{
		group.ID,
		group.DisplayName,
		group.Description,
		"?",
		users,
	})
}

func groupTableWriterAppendFromMultiple(tw table.Writer, item *apigroups.ListGroupsOKBodyItems0) {
	tw.AppendRow(table.Row{
		item.ID,
		item.DisplayName,
		item.Description,
		item.TotalUsers.Enrolled,
		item.TotalUsers.Enrolled + item.TotalUsers.Pending + item.TotalUsers.Unenrolled,
	})
}

func groupTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	idStr := "[ERR]"
	if id != nil {
		idStr += fmt.Sprintf(" %v", id)
	}
	tw.AppendRow(table.Row{
		idStr,
		processErrorResponse(err),
		"-",
		"-",
		"-",
	})
}

func init() {
	groupsCmd.AddCommand(groupGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(groupGetCmd)
	groupGetCmd.Flags().Int("id", 0, "id of group to get")
}
