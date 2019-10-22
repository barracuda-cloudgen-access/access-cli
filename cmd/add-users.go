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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiusers "github.com/fyde/fyde-cli/client/users"
	"github.com/fyde/fyde-cli/models"
)

// usersAddCmd represents the get command
var usersAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add users",
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
		err := forAllInput(cmd, func(values []interface{}) error {
			total++ // this is the total of successful+failures, must increment before failure
			user := &apiusers.CreateUserParamsBodyUser{}
			err := placeInputValues(cmd, values, user,
				func(s string) { user.Name = s },
				func(s string) { user.Email = strfmt.Email(s) },
				func(s string) { user.PhoneNumber = s },
				func(s []int64) { user.GroupIds = s },
				func(s bool) { user.Enabled = s },
				func(s bool) { user.SendEmailInvitation = s })
			if err != nil {
				return err
			}
			body := apiusers.CreateUserBody{User: user}
			params := apiusers.NewCreateUserParams()
			params.SetUser(body)

			resp, err := global.Client.Users.CreateUser(params, global.AuthWriter)
			if err != nil {
				return err
			}
			createdList = append(createdList, &resp.Payload.User)
			userTableWriterAppend(tw, resp.Payload.User)
			return nil
		}, func(err error) {
			createdList = append(createdList, nil)
			userTableWriterAppendError(tw, err)
		})
		if err != nil {
			return processErrorResponse(err)
		}
		result, err := renderListOutput(cmd, createdList, tw, total)
		cmd.Println(result)
		return err
	},
}

func init() {
	usersCmd.AddCommand(usersAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// usersAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// usersAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(usersAddCmd)
	initLoopControlFlags(usersAddCmd)

	initInputFlags(usersAddCmd,
		inputField{
			Name:            "Username",
			FlagName:        "username",
			FlagDescription: "specify the username for the created user",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Email",
			FlagName:        "email",
			FlagDescription: "specify the email for the created user",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Phone",
			FlagName:        "phone",
			FlagDescription: "specify the phone for the created user",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Groups",
			FlagName:        "groups",
			FlagDescription: "specify the group IDs for the created user",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
		},
		inputField{
			Name:            "Enabled",
			FlagName:        "enabled",
			FlagDescription: "whether the created user will be enabled",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    true,
		},
		inputField{
			Name:            "Send email invitation",
			FlagName:        "invitation",
			FlagDescription: "whether to send an email invitation",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    false,
		})
}
