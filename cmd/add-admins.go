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

	apiadmins "github.com/fyde/fyde-cli/client/admins"
	"github.com/fyde/fyde-cli/models"
)

// adminsAddCmd represents the add command
var adminsAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add admins",
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
		tw := adminBuildTableWriter()
		createdList := []*models.Admin{}
		total := 0
		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				admin := &apiadmins.CreateAdminParamsBodyAdmin{}
				err := placeInputValues(cmd, values, admin,
					func(s string) { admin.Name = s },
					func(s string) { admin.Email = strfmt.Email(s) },
					func(s string) { admin.AuthenticationType = s },
					func(s string) { admin.Password = s },
					func(s string) { admin.TenantOwner = s })
				if err != nil {
					return nil, err
				}
				body := apiadmins.CreateAdminBody{Admin: admin}
				params := apiadmins.NewCreateAdminParams()
				params.SetAdmin(body)

				resp, err := global.Client.Admins.CreateAdmin(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				admin := data.(*models.Admin)
				createdList = append(createdList, admin)
				adminTableWriterAppend(tw, admin)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				adminTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	adminsCmd.AddCommand(adminsAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adminsAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adminsAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(adminsAddCmd)
	initLoopControlFlags(adminsAddCmd)

	initInputFlags(adminsAddCmd, "admin",
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the name for the admin",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "name",
		},
		inputField{
			Name:            "Email",
			FlagName:        "email",
			FlagDescription: "specify the email for the created admin",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
		},
		inputField{
			Name:            "Authentication Type",
			FlagName:        "authn-type",
			FlagDescription: "specify the new authentication type for the admin",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "authentication_type",
		},
		inputField{
			Name:            "Password",
			FlagName:        "password",
			FlagDescription: "specify the new password for the admin",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "password",
		},
		inputField{
			Name:            "Tenant Owner",
			FlagName:        "owner",
			FlagDescription: "Is this admin a Tenant Owner (able to manage admins)",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "tenant_owner",
		})
}
