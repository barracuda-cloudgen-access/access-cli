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
	"github.com/thoas/go-funk"

	apiusers "github.com/barracuda-cloudgen-access/access-cli/client/users"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// userGetCmd represents the get command
var userGetCmd = &cobra.Command{
	Use:   "get [user ID]",
	Short: "Get user",
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
			return fmt.Errorf("missing user ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var userID int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			userID = int64(d)
		} else {
			userID, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		params := apiusers.NewGetUserParams()
		setTenant(cmd, params)
		params.SetID(userID)

		resp, err := global.Client.Users.GetUser(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := userBuildTableWriter()
		userTableWriterAppend(tw, resp.Payload.User)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func userBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Email",
		"Groups",
		"Enabled",
		"Status",
		"EnrollmentStatus",
	})
	tw.SetAlign([]text.Align{
		text.AlignRight,
		text.AlignLeft,
		text.AlignLeft,
		text.AlignLeft,
		text.AlignLeft,
		text.AlignLeft,
		text.AlignLeft})
	tw.SetAllowedColumnLengths([]int{15, 30, 30, 30, 10, 15, 16})
	return tw
}

func userTableWriterAppend(tw table.Writer, user models.User) {
	groups := strings.Join(funk.Map(user.Groups, func(g *models.UserGroupsItems0) string {
		return g.Name
	}).([]string), ",")

	tw.AppendRow(table.Row{
		user.ID,
		user.Name,
		user.Email,
		groups,
		user.Enabled,
		user.Status,
		user.EnrollmentStatus,
	})
}

func userTableWriterAppendError(tw table.Writer, err error, id interface{}) {
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
		"-",
	})
}

func init() {
	usersCmd.AddCommand(userGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(userGetCmd)
	initTenantFlags(userGetCmd)

	userGetCmd.Flags().Int("id", 0, "id of user to get")
}
