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

	"github.com/spf13/cobra"

	apiusers "github.com/fyde/fyde-cli/client/users"
)

// userDeleteCmd represents the delete command
var userDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete user",
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
		for _, arg := range args {
			params := apiusers.NewDeleteUserParams()

			userID, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return err
			}
			params.SetID(userID)

			_, err = global.Client.Users.DeleteUser(params, global.AuthWriter)
			if err != nil {
				if loopControlContinueOnError(cmd) {
					cmd.PrintErrln(processErrorResponse(err))
					continue
				}
				return processErrorResponse(err)
			}

			cmd.Println("User", userID, "deleted")
		}
		return nil
	},
}

func init() {
	usersCmd.AddCommand(userDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initLoopControlFlags(userDeleteCmd)
}
