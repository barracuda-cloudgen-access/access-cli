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

	apipolicies "github.com/barracuda-cloudgen-access/access-cli/client/access_policies"
)

// policiesEditCmd represents the edit command
var policiesEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit policies",
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
		tw := policyBuildTableWriter()
		createdList := []*apipolicies.EditPolicyOKBody{}
		total := 0
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apipolicies.NewEditPolicyParams()
				// IDs are not part of the request body, so we use this workaround
				policy := &struct {
					apipolicies.EditPolicyParamsBodyAccessPolicy
					ID int64 `json:"id"`
				}{}
				err := placeInputValues(cmd, values, policy,
					func(s int) { policy.ID = int64(s) },
					func(s string) { policy.Name = s },
					func(s []strfmt.UUID) { policy.AccessResourceIds = s },
					func(rbac bool) {
						if policy.Conditions == nil {
							policy.Conditions = &apipolicies.EditPolicyParamsBodyAccessPolicyConditions{}
							policy.Conditions.Rbac = &apipolicies.EditPolicyParamsBodyAccessPolicyConditionsRbac{
								UserIds:  []int64{},
								GroupIds: []int64{},
							}
						}
						policy.Conditions.Rbac.Enabled = &rbac
					},
					func(groups []int64) {
						if policy.Conditions == nil {
							t := true
							policy.Conditions = &apipolicies.EditPolicyParamsBodyAccessPolicyConditions{}
							policy.Conditions.Rbac = &apipolicies.EditPolicyParamsBodyAccessPolicyConditionsRbac{
								Enabled: &t,
								UserIds: []int64{},
							}
						}
						policy.Conditions.Rbac.GroupIds = groups
					},
					func(users []int64) {
						if policy.Conditions == nil {
							t := true
							policy.Conditions = &apipolicies.EditPolicyParamsBodyAccessPolicyConditions{}
							policy.Conditions.Rbac = &apipolicies.EditPolicyParamsBodyAccessPolicyConditionsRbac{
								Enabled:  &t,
								GroupIds: []int64{},
							}
						}
						policy.Conditions.Rbac.UserIds = users
					})
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(policy.ID)
				body := apipolicies.EditPolicyBody{AccessPolicy: &policy.EditPolicyParamsBodyAccessPolicy}
				params.SetPolicy(body)

				resp, err := global.Client.AccessPolicies.EditPolicy(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				policy := data.(*apipolicies.EditPolicyOKBody)
				createdList = append(createdList, policy)
				policyTableWriterAppend(tw, policy.AccessPolicy, len(policy.AccessResources))
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				policyTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	policiesCmd.AddCommand(policiesEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policiesEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policiesEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(policiesEditCmd)
	initLoopControlFlags(policiesEditCmd)

	initInputFlags(policiesEditCmd, "policy",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the policy to edit",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    0,
			MainField:       true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the new name for the policy",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Resources",
			FlagName:        "resources",
			FlagDescription: "specify the new resources for the policy",
			VarType:         "[]string",
			Mandatory:       false,
			DefaultValue:    []string{},
		},
		inputField{
			Name:            "RBAC",
			FlagName:        "rbac",
			FlagDescription: "whether to enable role-based access control (RBAC) for the policy",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    false,
		},
		inputField{
			Name:            "Groups",
			FlagName:        "groups",
			FlagDescription: "specify the new RBAC groups for the policy",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
		},
		inputField{
			Name:            "Users",
			FlagName:        "users",
			FlagDescription: "specify the new RBAC users for the policy",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
		})
}
