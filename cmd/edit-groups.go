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
	"github.com/spf13/cobra"

	apigroups "github.com/barracuda-cloudgen-access/access-cli/client/groups"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// groupsEditCmd represents the edit command
var groupsEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit groups",
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
		tw := groupBuildTableWriter()
		createdList := []*apigroups.EditGroupOKBody{}
		total := 0
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apigroups.NewEditGroupParams()
				setTenant(cmd, params)
				// IDs are not part of the request body, so we use this workaround
				group := &struct {
					models.Group
					ID int64 `json:"id"`
				}{}
				err := placeInputValues(cmd, values, group,
					func(s int) { group.ID = int64(s) },
					func(s string) { group.Name = s },
					func(s string) { group.Description = s },
					func(s string) { group.Color = s })
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(group.ID)
				body := apigroups.EditGroupBody{Group: &group.Group}
				params.SetGroup(body)

				resp, err := global.Client.Groups.EditGroup(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				group := data.(*apigroups.EditGroupOKBody)
				createdList = append(createdList, group)
				groupTableWriterAppendFromSingle(tw, group.Group, len(group.Users))
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				groupTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	groupsCmd.AddCommand(groupsEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupsEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupsEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(groupsEditCmd)
	initLoopControlFlags(groupsEditCmd)
	initTenantFlags(groupsEditCmd)
	initInputFlags(groupsEditCmd, "group",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the group to edit",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    0,
			MainField:       true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the new name for the group",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Description",
			FlagName:        "description",
			FlagDescription: "specify the new description for the group",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Color",
			FlagName:        "color",
			FlagDescription: "specify the new color for the group (hexadecimal #RRGGBB format)",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			Validator:       validateHTMLHexColor,
		})
}
