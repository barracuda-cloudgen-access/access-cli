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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiusers "github.com/barracuda-cloudgen-access/access-cli/client/users"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// usersAddCmd represents the add command
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

		// Assign deprecated username if name was not supplied
		name, _ := cmd.Flags().GetString("name")
		username, _ := cmd.Flags().GetString("username")
		if name == "" && username != "" {
			cmd.Flags().Set("name", username)
		}

		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				user := &struct {
					apiusers.CreateUserParamsBodyUser
					Groups []struct {
						ID int64 `json:"id"`
					}
				}{}

				err := placeInputValues(cmd, values, user,
					func(s string) { /* deprecated username already handled */ },
					func(s string) { user.Name = s },
					func(s string) { user.Email = strfmt.Email(s) },
					func(s string) { user.PhoneNumber = s },
					func(s []int64) { user.GroupIds = s },
					func(s bool) { user.Enabled = s },
					func(s bool) { user.SendEmailInvitation = s })
				if err != nil {
					return nil, err
				}
				body := apiusers.CreateUserBody{User: &user.CreateUserParamsBodyUser}

				// map group ids since GET and POST are not exactly the same (when adding from file)
				for _, group := range user.Groups {
					body.User.GroupIds = append(body.User.GroupIds, group.ID)
				}

				params := apiusers.NewCreateUserParams()

				setTenant(cmd, params)
				params.SetUser(body)

				resp, err := global.Client.Users.CreateUser(params, global.AuthWriter)
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
	initTenantFlags(usersAddCmd)

	initInputFlags(usersAddCmd, "user",
		//deprecated
		inputField{
			Name:            "Username",
			FlagName:        "username",
			FlagDescription: "specify the username for the created user",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the name for the created user",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "name",
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
	usersAddCmd.Flags().MarkDeprecated("username", "use name instead")

}
