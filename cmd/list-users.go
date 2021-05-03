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
	"github.com/spf13/cobra"

	apiusers "github.com/fyde/access-cli/client/users"
)

// usersListCmd represents the list command
var usersListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List users",
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
		params := apiusers.NewListUsersParams()
		setSort(cmd, params)
		setFilter(cmd, params.SetGroupID, params.SetGroupName, params.SetStatus, params.SetEnrollmentStatus)
		setSearchQuery(cmd, params)
		completePayload := []*apiusers.ListUsersOKBodyItems0{}
		total := 0
		cutStart, cutEnd, err := forAllPages(cmd, params, func() (int, int64, error) {
			resp, err := global.Client.Users.ListUsers(params, global.AuthWriter)
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

		tw := userBuildTableWriter()

		for _, item := range completePayload {
			userTableWriterAppend(tw, item.User)
		}

		return printListOutputAndError(cmd, completePayload, tw, total, err)
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

	initPaginationFlags(usersListCmd)
	initSortFlags(usersListCmd)
	initFilterFlags(usersListCmd,
		filterType{"group", "[]int"},
		filterType{"group-name", "[]string"},
		filterType{"enabled-status", "string"},
		filterType{"status", "string"})
	initSearchFlags(usersListCmd)
	initOutputFlags(usersListCmd)
}
