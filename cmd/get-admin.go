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
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apiadmins "github.com/fyde/access-cli/client/admins"
	"github.com/fyde/access-cli/models"
)

// adminGetCmd represents the get command
var adminGetCmd = &cobra.Command{
	Use:   "get [admin ID]",
	Short: "Get admin",
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
			return fmt.Errorf("missing admin ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			id = int64(d)
		} else {
			id, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		params := apiadmins.NewGetAdminParams()
		params.SetID(id)

		resp, err := global.Client.Admins.GetAdmin(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := adminBuildTableWriter()
		adminTableWriterAppend(tw, &resp.Payload.Admin)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func adminBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Email",
		"Authentication Type",
		"Authentication Email",
		"Roles",
		"Last Sign In",
	})
	tw.SetAllowedColumnLengths([]int{36, 30, 30, 30, 30, 30, 36})
	return tw
}

func adminTableWriterAppend(tw table.Writer, admin *models.Admin) {
	tw.AppendRow(table.Row{
		admin.ID,
		admin.Name,
		admin.Email,
		admin.AuthenticationType,
		admin.AuthenticationEmail,
		strings.Join(admin.RoleNames, ","),
		admin.LastSignInAt,
	})
}

func adminTableWriterAppendError(tw table.Writer, err error, id interface{}) {
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
		"-",
	})
}

func init() {
	adminsCmd.AddCommand(adminGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adminGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adminGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(adminGetCmd)
	adminGetCmd.Flags().String("id", "", "id of admin to get")
}
