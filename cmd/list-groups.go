// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2023 Barracuda Networks, Inc.

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

	apigroups "github.com/barracuda-cloudgen-access/access-cli/client/groups"
)

// groupsListCmd represents the list command
var groupsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List groups",
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
		params := apigroups.NewListGroupsParams()
		setTenant(cmd, params)
		setSort(cmd, params)
		setSearchQuery(cmd, params)
		completePayload := []*apigroups.ListGroupsOKBodyItems0{}
		total := 0
		cutStart, cutEnd, err := forAllPages(cmd, params, func() (int, int64, error) {
			resp, err := global.Client.Groups.ListGroups(params, global.AuthWriter)
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

		tw := groupBuildTableWriter()

		for _, item := range completePayload {
			groupTableWriterAppendFromMultiple(tw, item)
		}

		return printListOutputAndError(cmd, completePayload, tw, total, err)
	},
}

func init() {
	groupsCmd.AddCommand(groupsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupsListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupsListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initPaginationFlags(groupsListCmd)
	initSortFlags(groupsListCmd)
	initSearchFlags(groupsListCmd)
	initOutputFlags(groupsListCmd)
	initTenantFlags(groupsListCmd)
	groupsListCmd.Flags().StringP("filter", "f", "", "filter groups")
}
