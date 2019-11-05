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

	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"

	apigroups "github.com/fyde/fyde-cli/client/groups"
)

// groupDeleteCmd represents the delete command
var groupDeleteCmd = &cobra.Command{
	Use:     "delete",
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

		if len(args) == 0 {
			return fmt.Errorf("missing group ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		groupIDs := make([]int64, len(args))
		var err error
		for i, arg := range args {
			groupIDs[i], err = strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return err
			}
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

		if loopControlContinueOnError(cmd) {
			// then we must delete individually, because on a request for multiple deletions,
			// the server does nothing if one fails
			i := 0
			for _, groupID := range groupIDs {
				err = delete([]int64{groupID})
				if err != nil {
					cmd.PrintErrln(err)
				} else {
					// only keep successful deletions in list of groupIDs
					// this rewrites the array in place and lets us "delete" as we iterate
					// (junk is removed after the loop)
					groupIDs[i] = groupID
					i++
				}
			}
			// remove junk left at end of slice
			groupIDs = groupIDs[:i]
		} else {
			err = delete(groupIDs)
			if err != nil {
				return err
			}
		}

		cmd.Println("Groups",
			strings.Join(
				funk.Map(
					groupIDs,
					func(i int64) string {
						return strconv.Itoa(int(i))
					}).([]string),
				", "), "deleted")
		return nil
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

	initLoopControlFlags(groupDeleteCmd)
}
