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
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"

	apiusers "github.com/fyde/fyde-cli/client/users"
	"github.com/fyde/fyde-cli/models"
)

// userGetCmd represents the get command
var userGetCmd = &cobra.Command{
	Use:   "get",
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

		if len(args) == 0 {
			return fmt.Errorf("missing user ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		params := apiusers.NewGetUserParams()
		params.SetID(userID)

		resp, err := global.Client.Users.GetUser(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

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
		tw.SetAllowedColumnLengths([]int{15, 30, 30, 30, 10, 15, 16})

		groups := strings.Join(funk.Map(resp.Payload.Groups, func(g *models.UserGroupsItems0) string {
			return g.Name
		}).([]string), ",")

		tw.AppendRow(table.Row{
			resp.Payload.ID,
			resp.Payload.Name,
			resp.Payload.Email,
			groups,
			resp.Payload.Enabled,
			resp.Payload.Status,
			resp.Payload.EnrollmentStatus,
		})

		result, err := renderListOutput(cmd, resp.Payload, tw, 1)
		cmd.Println(result)
		return err
	},
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
}
