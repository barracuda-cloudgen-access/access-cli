// Package cmd implements access-cli commands
package cmd

import (
	"fmt"

	apiwebpolicies "github.com/barracuda-cloudgen-access/access-cli/client/web_policies"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

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

var webPoliciesEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit web policies",
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
		tw := setWebPolicyBuildTableWriter()
		gparams := apiwebpolicies.NewListWebPoliciesParams()
		setTenant(cmd, gparams)
		resp, err := global.Client.WebPolicies.ListWebPolicies(gparams, global.AuthWriter)

		if err != nil {
			return processErrorResponse(err)
		}
		mainRulesetId := resp.Payload.ID
		policies := resp.Payload

		jumpRules := policies.Rules
		mapRules := make(map[strfmt.UUID]strfmt.UUID)
		for _, rule := range jumpRules {
			if rule.Type == "jump" {
				// make sure there is a rule in the jump ruleset
				if len(rule.RulesetJump.Rules) >= 1 {
					mapRules[rule.RulesetJump.Rules[0].ID] = rule.ID
				}
			}
		}

		modifiedList := []struct {
			models.WebPolicyData
			ID strfmt.UUID `json:"id"`
		}{}

		modifiedPolicy := &struct {
			models.WebPolicyData
			ID strfmt.UUID `json:"id"`
		}{}
		total := 0
		err = forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure

				params := apiwebpolicies.NewEditWebPolicyParams()
				setTenant(cmd, params)
				params.SetRulesetID(mainRulesetId)
				policy := &struct {
					apiwebpolicies.EditWebPolicyParamsBodyData
					ID strfmt.UUID `json:"id"`
				}{}

				err := placeInputValues(cmd, values, policy,
					func(s string) { policy.ID = strfmt.UUID(s) },
					func(s string) { policy.Label = s },
					func(s string) { policy.Action = s },
					func(s string) { policy.Type = s },
					func(s []string) { policy.Domains = s },
					func(s []int64) { policy.Categories = s },
					func(s []int64) { policy.UserIds = s },
					func(s []int64) { policy.GroupIds = s },
					func(s int) {
						idx := int64(s)
						policy.Index = idx
					},
					func(s bool) { policy.Log = s },
					func(s bool) { policy.Notify = s },
					func(s bool) { policy.Disabled = s },
					func(s bool) { policy.Alert = s },
				)
				if err != nil {
					fmt.Println("error: ", err)
					return nil, err
				}

				// here, map the ID from the "fake request body" to the correct place
				params = apiwebpolicies.NewEditWebPolicyParams()
				setTenant(cmd, params)
				params.SetRulesetID(mainRulesetId)
				params.SetRuleID(policy.ID)
				params.SetWebPolicy(apiwebpolicies.EditWebPolicyBody{Data: &policy.EditWebPolicyParamsBodyData})

				_, err = global.Client.WebPolicies.EditWebPolicy(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}

				if policy.GroupIds != nil || policy.UserIds != nil {
					//assign users / group ids to the jump rule
					if mapRules[policy.ID] == "" {
						return nil, fmt.Errorf("no jump rule found for policy %s", policy.ID)
					}
					jumpRuleId, ok := mapRules[policy.ID]
					if !ok {
						return nil, fmt.Errorf("no jump rule found for policy %s", policy.ID)
					}

					jumpRule := &struct {
						apiwebpolicies.EditWebPolicyParamsBodyData
						ID strfmt.UUID `json:"id"`
					}{}
					jumpRule.ID = jumpRuleId
					jumpRule.UserIds = policy.UserIds
					jumpRule.GroupIds = policy.GroupIds

					params.SetRuleID(jumpRuleId)
					params.SetWebPolicy(apiwebpolicies.EditWebPolicyBody{Data: &jumpRule.EditWebPolicyParamsBodyData})
					fmt.Print("jump rule params: ", params)

					//update the jump rule
					_, err = global.Client.WebPolicies.EditWebPolicy(params, global.AuthWriter)
					if err != nil {
						return nil, err
					}
				}

				modifiedPolicy.Action = policy.Action
				modifiedPolicy.Alert = policy.Alert
				modifiedPolicy.Categories = policy.Categories
				modifiedPolicy.Disabled = policy.Disabled
				modifiedPolicy.Domains = policy.Domains
				modifiedPolicy.GroupIds = policy.GroupIds
				modifiedPolicy.Index = policy.Index
				modifiedPolicy.Label = policy.Label
				modifiedPolicy.Log = policy.Log
				modifiedPolicy.Notify = policy.Notify
				modifiedPolicy.Type = policy.Type
				modifiedPolicy.UserIds = policy.UserIds
				modifiedPolicy.ID = policy.ID

				modifiedList = append(modifiedList, *modifiedPolicy)
				return modifiedPolicy, nil
			}, func(data interface{}) {

				//printSuccess func
				modifiedList = append(modifiedList, *modifiedPolicy)
				setWebPolicyTableWriterAppend(tw, modifiedPolicy.ID, modifiedPolicy.WebPolicyData)
			},
			func(err error, id interface{}) {
				// doOnError func
				modifiedList = append(modifiedList, struct {
					models.WebPolicyData
					ID strfmt.UUID `json:"id"`
				}{})
				webPolicyTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, modifiedList, tw, total, err)
	},
}

func init() {
	webPoliciesCmd.AddCommand(webPoliciesEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webPoliciesEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webPoliciesEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initOutputFlags(webPoliciesEditCmd)
	initLoopControlFlags(webPoliciesEditCmd)
	initTenantFlags(webPoliciesEditCmd)
	initInputFlags(webPoliciesEditCmd, "webpolicy",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the webpolicy to edit",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Label",
			FlagName:        "label",
			FlagDescription: "specify the new name for the policy",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Action",
			FlagName:        "action",
			FlagDescription: "specify the new action for the policy [block|allow]",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "PolicyType",
			FlagName:        "policyType",
			FlagDescription: "specify the type for the created policy [domain|category]",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "type",
		},
		inputField{
			Name:            "Domains",
			FlagName:        "domains",
			FlagDescription: "specify the domains for the created policy",
			VarType:         "[]string",
			Mandatory:       false,
			DefaultValue:    []string{},
		},
		inputField{
			Name:            "Categories",
			FlagName:        "categories",
			FlagDescription: "specify the category ids for the created policy",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
		},
		inputField{
			Name:            "Users",
			FlagName:        "user_ids",
			FlagDescription: "specify the user ids for the created policy",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
			MainField:       true,
			SchemaName:      "user_ids",
		},
		inputField{
			Name:            "Groups",
			FlagName:        "group_ids",
			FlagDescription: "specify the group ids for the created policy",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
			MainField:       true,
			SchemaName:      "group_ids",
		},
		inputField{
			Name:            "Index",
			FlagName:        "index",
			FlagDescription: "specify the index (position within the rule set)  for the created policy",
			VarType:         "int",
			Mandatory:       false,
			DefaultValue:    0,
			MainField:       true,
			SchemaName:      "index",
		},
		inputField{
			Name:            "Log",
			FlagName:        "log",
			FlagDescription: "whether to enable logging for this policy",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    false,
			MainField:       true,
			SchemaName:      "log",
		},
		inputField{
			Name:            "Notify",
			FlagName:        "notify",
			FlagDescription: "whether to notify admins when this policy is triggered",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    false,
			MainField:       true,
			SchemaName:      "notify",
		},
		inputField{
			Name:            "Disabled",
			FlagName:        "disabled",
			FlagDescription: "whether to set this policy as disabled",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    false,
			MainField:       true,
			SchemaName:      "disabled",
		},
		inputField{
			Name:            "Alert",
			FlagName:        "alert",
			FlagDescription: "alert admins when this policy is triggered",
			VarType:         "bool",
			Mandatory:       false,
			DefaultValue:    false,
			MainField:       true,
			SchemaName:      "alert",
		})
}
