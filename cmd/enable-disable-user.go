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
	"github.com/fyde/fyde-cli/models"
)

// userEnableCmd represents the enable command
var userEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable user",
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
		enable := cmd.Use == "enable"
		tw := userBuildTableWriter()
		createdList := []*models.User{}
		for _, arg := range args {
			userID, err := strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return err
			}

			params := apiusers.NewEditUserParams()
			params.SetID(userID)
			params.SetUser(apiusers.EditUserBody{
				User: &apiusers.EditUserParamsBodyUser{
					Enabled: &enable,
				},
			})

			resp, err := global.Client.Users.EditUser(params, global.AuthWriter)
			if err != nil {
				if loopControlContinueOnError(cmd) {
					createdList = append(createdList, nil)
					userTableWriterAppendError(tw, err, userID)
					continue
				}
				return processErrorResponse(err)
			}
			createdList = append(createdList, &resp.Payload.User)
			userTableWriterAppend(tw, resp.Payload.User)
		}
		result, err := renderListOutput(cmd, createdList, tw, len(args))
		cmd.Println(result)
		return err
	},
}

// userDisableCmd represents the disable command
var userDisableCmd *cobra.Command

func init() {
	disableCmd := *userEnableCmd
	disableCmd.Use = "disable"
	disableCmd.Short = "disable user"
	userDisableCmd = &disableCmd
	usersCmd.AddCommand(userEnableCmd)
	usersCmd.AddCommand(userDisableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userEnableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userEnableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(userEnableCmd)
	initOutputFlags(userDisableCmd)

	initLoopControlFlags(userEnableCmd)
	initLoopControlFlags(userDisableCmd)
}
