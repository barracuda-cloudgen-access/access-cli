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

	api "github.com/barracuda-cloudgen-access/access-cli/client/user_directories"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// userDirectoriesEditCmd represents the edit command
var userDirectoriesEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit user directories",
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
		tw := userDirectoryBuildTableWriter()
		createdList := []*models.UserDirectory{}
		total := 0

		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := api.NewEditUserDirectoryParams()
				setTenant(cmd, params)

				userdirectory := &models.UserDirectory{}

				err := placeInputValues(cmd, values, userdirectory,
					func(s int) { userdirectory.ID = int64(s) },
					func(s string) { userdirectory.ShortCode = s },
					func(s string) { userdirectory.Name = s },
					func(s string) { userdirectory.DirectoryType = s },
					func(s string) { userdirectory.Notes = s })
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(userdirectory.ID)
				body := api.EditUserDirectoryBody{UserDirectory: userdirectory}

				params.SetUserDirectory(body)

				resp, err := global.Client.UserDirectories.EditUserDirectory(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload.UserDirectory, nil
			}, func(data interface{}) { // printSuccess func
				userdirectory := data.(models.UserDirectory)
				createdList = append(createdList, &userdirectory)
				userDirectoryTableWriterAppend(tw, userdirectory)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				userDirectoryTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	settingsUserDirectoryCmd.AddCommand(userDirectoriesEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userDirectoriesEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userDirectoriesEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(userDirectoriesEditCmd)
	initLoopControlFlags(userDirectoriesEditCmd)
	initTenantFlags(userDirectoriesEditCmd)

	initInputFlags(userDirectoriesEditCmd, "userdirectory",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the user to edit",
			VarType:         "int",
			DefaultValue:    0,
			Mandatory:       true,
			MainField:       true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "ShortCode",
			FlagName:        "shortcode",
			FlagDescription: "specify a shortcode for the directory",
			VarType:         "string",
			DefaultValue:    "",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the name for the directory",
			VarType:         "string",
			DefaultValue:    "",
			MainField:       true,
		},
		inputField{
			Name:            "Type",
			FlagName:        "type",
			FlagDescription: "specify source for the user directory",
			VarType:         "string",
			DefaultValue:    "",
		},
		inputField{
			Name:            "Notes",
			FlagName:        "notes",
			FlagDescription: "additional notes",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		})
}
