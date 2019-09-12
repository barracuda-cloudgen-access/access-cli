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
package cmd

import (
	"fmt"
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"

	apiusers "github.com/oNaiPs/fyde-cli/client/users"
	"github.com/oNaiPs/fyde-cli/models"
)

// usersListCmd represents the list command
var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagCheckOutput(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		uparams := apiusers.NewListUsersParams()
		resp, err := global.Client.Users.ListUsers(uparams, global.AuthWriter)
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

		for _, item := range resp.Payload {
			groups := strings.Join(funk.Map(item.Groups, func(g *models.UserGroupsItems0) string {
				return g.Name
			}).([]string), ",")

			tw.AppendRow(table.Row{
				item.ID,
				item.Name,
				item.Email,
				groups,
				item.Enabled,
				item.Status,
				item.EnrollmentStatus,
			})
		}

		var result string
		outputFormat, err := cmd.Flags().GetString("output")
		if err != nil {
			return err
		}
		switch outputFormat {
		case "table":
			result = tw.Render()
		case "csv":
			result = tw.RenderCSV()
		}
		fmt.Println(result)
		return nil
	},
}

func init() {
	usersCmd.AddCommand(usersListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	usersListCmd.Flags().StringP("filter", "f", "", "filter users")
	usersListCmd.Flags().StringP("output", "o", "table", "output format (table, json or csv)")
}
