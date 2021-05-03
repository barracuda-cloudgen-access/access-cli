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

	apiadmins "github.com/fyde/access-cli/client/admins"
	"github.com/fyde/access-cli/models"
)

// adminsEditCmd represents the edit command
var adminsEditCmd = &cobra.Command{
	Use:                "edit",
	Short:              "Edit admins",
	FParseErrWhitelist: cobra.FParseErrWhitelist{},
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
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apiadmins.NewEditAdminParams()
				// IDs are not part of the request body, so we use this workaround
				admin := &struct {
					apiadmins.EditAdminParamsBodyAdmin
					ID int64 `json:"id"`
				}{}
				err := placeInputValues(cmd, values, admin,
					func(s int) { admin.ID = int64(s) },
					func(s string) { admin.Name = s },
					func(s string) { admin.AuthenticationType = s },
					func(s string) { admin.AuthenticationEmail = strfmt.Email(s) },
					func(s string) { admin.Password = s },
					func(s []string) { admin.RoleNames = s })
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(admin.ID)
				body := apiadmins.EditAdminBody{Admin: &admin.EditAdminParamsBodyAdmin}
				params.SetAdmin(body)

				resp, err := global.Client.Admins.EditAdmin(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				resp := data.(*apiadmins.EditAdminOKBody)
				createdList = append(createdList, &resp.Admin)
				adminTableWriterAppend(tw, &resp.Admin)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				adminTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	adminsCmd.AddCommand(adminsEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adminsEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adminsEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(adminsEditCmd)
	initLoopControlFlags(adminsEditCmd)

	initInputFlags(adminsEditCmd, "admin",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the admin to edit",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    0,
			MainField:       true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the new name for the admin",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "name",
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
			Name:            "Authentication Email",
			FlagName:        "authn-email",
			FlagDescription: "specify the new email used for authentication for this admin",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "authentication_email",
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
			Name:            "Roles",
			FlagName:        "roles",
			FlagDescription: "List of roles for this admin",
			VarType:         "[]string",
			Mandatory:       false,
			DefaultValue:    []string{},
			SchemaName:      "role_names",
		})
}
