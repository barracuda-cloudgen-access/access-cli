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

	"github.com/spf13/cobra"

	apigroups "github.com/barracuda-cloudgen-access/access-cli/client/groups"
)

// groupDeleteCmd represents the delete command
var groupDeleteCmd = &cobra.Command{
	Use:     "delete [group ID]...",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete groups",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		if !multiOpCheckArgsPresent(cmd, args) {
			return fmt.Errorf("missing group ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		groupIDs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		delete := func(ids []int64) error {
			params := apigroups.NewDeleteGroupParams()
			params.SetID(ids)

			_, err = global.Client.Groups.DeleteGroup(params, global.AuthWriter)
			if err != nil {
				return processErrorResponse(err)
			}
			return nil
		}

		tw, j := multiOpBuildTableWriter()

		if loopControlContinueOnError(cmd) {
			// then we must delete individually, because on a request for multiple deletions,
			// the server does nothing if one fails

			for _, id := range groupIDs {
				err = delete([]int64{id})
				var result interface{}
				result = "success"
				if err != nil {
					result = err
				}
				multiOpTableWriterAppend(tw, &j, id, result)
			}
			err = nil
		} else {
			err = delete(groupIDs)
			var result interface{}
			result = "success"
			if err != nil {
				result = err
			}
			multiOpTableWriterAppend(tw, &j, "*", result)
		}

		return printListOutputAndError(cmd, j, tw, len(groupIDs), err)
	},
}

func init() {
	groupsCmd.AddCommand(groupDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initMultiOpArgFlags(groupDeleteCmd, "group", "delete", "id", "[]int64")
	initOutputFlags(groupDeleteCmd)
	initLoopControlFlags(groupDeleteCmd)
}
