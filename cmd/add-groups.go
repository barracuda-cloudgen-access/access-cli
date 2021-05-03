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
	"regexp"

	"github.com/spf13/cobra"

	apigroups "github.com/fyde/access-cli/client/groups"
)

// groupsAddCmd represents the add command
var groupsAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add groups",
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
		createdList := []*apigroups.CreateGroupCreatedBody{}
		total := 0
		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				group := &apigroups.CreateGroupParamsBodyGroup{}
				err := placeInputValues(cmd, values, group,
					func(s string) { group.Name = s },
					func(s string) { group.Description = s },
					func(s string) { group.Color = s })
				if err != nil {
					return nil, err
				}
				body := apigroups.CreateGroupBody{Group: group}
				params := apigroups.NewCreateGroupParams()
				params.SetGroup(body)

				resp, err := global.Client.Groups.CreateGroup(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				group := data.(*apigroups.CreateGroupCreatedBody)
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
	groupsCmd.AddCommand(groupsAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupsAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupsAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(groupsAddCmd)
	initLoopControlFlags(groupsAddCmd)

	initInputFlags(groupsAddCmd, "group",
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the name for the created group",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "name",
		},
		inputField{
			Name:            "Description",
			FlagName:        "description",
			FlagDescription: "specify the description for the created group",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Color",
			FlagName:        "color",
			FlagDescription: "specify the color for the created group (hexadecimal #RRGGBB format)",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			Validator:       validateHTMLHexColor,
		})
}

func validateHTMLHexColor(input interface{}) bool {
	c, ok := input.(string)
	if !ok {
		return false
	}
	if c == "" {
		return true
	}
	matched, err := regexp.MatchString(`#[0-9a-fA-F]{6}`, c)
	return err == nil && matched
}
