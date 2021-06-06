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

	"github.com/spf13/cobra"

	api "github.com/barracuda-cloudgen-access/access-cli/client/user_directories"
)

// userDirectoryDeleteCmd represents the delete command
var userDirectoryDeleteCmd = &cobra.Command{
	Use:     "delete [userdirectory ID]",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete user directories",
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
			return fmt.Errorf("missing proxy ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int64
		var err error
		if cmd.Flags().Changed("id") {
			id, err = cmd.Flags().GetInt64("id")
			if err != nil {
				return err
			}
		} else {
			id, err = strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
		}

		params := api.NewDeleteUserDirectoryParams()
		setTenant(cmd, params)
		params.SetID(id)

		tw, j := multiOpBuildTableWriter()
		_, err = global.Client.UserDirectories.DeleteUserDirectory(params, global.AuthWriter)

		var result interface{}
		result = "success"
		if err != nil {
			result = err
		}
		multiOpTableWriterAppend(tw, &j, "*", result)

		return printListOutputAndError(cmd, j, tw, 1, err)
	},
}

func init() {
	settingsUserDirectoryCmd.AddCommand(userDirectoryDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userDirectoryDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userDirectoryDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initInputFlags(userDirectoryDeleteCmd, "userdirectory",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the user to edit",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    0,
			MainField:       true,
		})
	initOutputFlags(userDirectoryDeleteCmd)
	initLoopControlFlags(userDirectoryDeleteCmd)
	initTenantFlags(userDirectoryDeleteCmd)
}
