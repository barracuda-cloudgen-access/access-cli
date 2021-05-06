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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiusers "github.com/barracuda-cloudgen-access/access-cli/client/users"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// usersEditCmd represents the edit command
var usersEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit users",
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
		tw := userBuildTableWriter()
		createdList := []*models.User{}
		total := 0
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apiusers.NewEditUserParams()
				setTenant(cmd, params)
				// IDs are not part of the request body, so we use this workaround
				enabledDefault := true
				user := &struct {
					apiusers.EditUserParamsBodyUser
					ID int64 `json:"id"`
				}{
					EditUserParamsBodyUser: apiusers.EditUserParamsBodyUser{
						Enabled: &enabledDefault, // the UI on the web console enables by default
					},
				}
				err := placeInputValues(cmd, values, user,
					func(s int) { user.ID = int64(s) },
					func(s string) { user.Name = s },
					func(s string) { user.Email = strfmt.Email(s) },
					func(s string) { user.PhoneNumber = s },
					func(s []int64) { user.GroupIds = s },
					func(s bool) { user.Enabled = &s })
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(user.ID)
				body := apiusers.EditUserBody{User: &user.EditUserParamsBodyUser}
				params.SetUser(body)

				resp, err := global.Client.Users.EditUser(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload.User, nil
			}, func(data interface{}) { // printSuccess func
				user := data.(models.User)
				createdList = append(createdList, &user)
				userTableWriterAppend(tw, user)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				userTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	usersCmd.AddCommand(usersEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(usersEditCmd)
	initLoopControlFlags(usersEditCmd)
	initTenantFlags(usersEditCmd)

	initInputFlags(usersEditCmd, "user",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the user to edit",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    0,
			MainField:       true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "Username",
			FlagName:        "username",
			FlagDescription: "specify the new username for the user",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Email",
			FlagName:        "email",
			FlagDescription: "specify the new email for the user",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Phone",
			FlagName:        "phone",
			FlagDescription: "specify the new phone for the user",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Groups",
			FlagName:        "groups",
			FlagDescription: "specify the new group IDs for the user",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
		},
		inputField{
			Name:            "Enabled",
			FlagName:        "enabled",
			FlagDescription: "whether the user is enabled",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    true,
		})
}
