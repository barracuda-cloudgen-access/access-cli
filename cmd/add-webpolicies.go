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
	apiwebpolicies "github.com/barracuda-cloudgen-access/access-cli/client/web_policies"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

// webPoliciesAddCmd represents the add command
var webPoliciesAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add policies",
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

		total := 0
		createdList := []*models.WebPolicyData{}
		createdPolicy := &models.WebPolicyData{}
		policy := &apiwebpolicies.AddWebPolicyParamsBodyData{}

		err = forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				policy = &apiwebpolicies.AddWebPolicyParamsBodyData{}
				err := placeInputValues(cmd, values, policy,
					func(s string) { policy.Label = &s },
					func(s string) { policy.Action = &s },
					func(s string) { policy.Type = &s },
					func(s []string) { policy.Domains = s },
					func(s []int64) { policy.Categories = s },
					func(s []int64) { policy.UserIds = s },
					func(s []int64) { policy.GroupIds = s },
					func(s int) {
						idx := int64(s)
						policy.Index = &idx
					},
					func(s bool) { policy.Log = &s },
					func(s bool) { policy.Notify = &s },
					func(s bool) { policy.Disabled = &s },
					func(s bool) { policy.Alert = &s },
				)
				if err != nil {
					createdPolicy = &models.WebPolicyData{}
					return nil, err
				}

				body := apiwebpolicies.AddWebPolicyBody{Data: policy}

				params := apiwebpolicies.NewAddWebPolicyParams()
				params.SetRulesetID(mainRulesetId)
				params.SetWebpolicy(body)

				setTenant(cmd, params)

				resp, err := global.Client.WebPolicies.AddWebPolicy(params, global.AuthWriter)

				if err != nil {
					return nil, err
				}

				createdPolicy.Action = *policy.Action
				createdPolicy.Categories = policy.Categories
				createdPolicy.Domains = policy.Domains
				createdPolicy.Index = *policy.Index
				createdPolicy.Label = *policy.Label
				createdPolicy.Log = *policy.Log
				createdPolicy.Notify = *policy.Notify
				createdPolicy.Type = *policy.Type
				createdPolicy.Disabled = *policy.Disabled
				createdPolicy.Alert = *policy.Alert
				createdPolicy.UserIds = policy.UserIds
				createdPolicy.GroupIds = policy.GroupIds

				return resp.Payload.RuleID, nil
			}, func(data interface{}) { // printSuccess func
				ruleId := data.(strfmt.UUID)
				createdList = append(createdList, createdPolicy)
				setWebPolicyTableWriterAppend(tw, ruleId, *createdPolicy)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, createdPolicy)
				webPolicyTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	webPoliciesCmd.AddCommand(webPoliciesAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policiesAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policiesAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(webPoliciesAddCmd)
	initLoopControlFlags(webPoliciesAddCmd)
	initTenantFlags(webPoliciesAddCmd)
	initInputFlags(webPoliciesAddCmd, "webpolicy",
		inputField{
			Name:            "Label",
			FlagName:        "label",
			FlagDescription: "specify a name for the created policy",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "label",
		},
		inputField{
			Name:            "Action",
			FlagName:        "action",
			FlagDescription: "specify the action for the created policy [block|allow]",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "block",
			MainField:       true,
			SchemaName:      "action",
		},
		inputField{
			Name:            "PolicyType",
			FlagName:        "policyType",
			FlagDescription: "specify the type for the created policy [domain|category]",
			VarType:         "string",
			Mandatory:       true,
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
			MainField:       true,
			SchemaName:      "domains",
		},
		inputField{
			Name:            "Categories",
			FlagName:        "categories",
			FlagDescription: "specify the category ids for the created policy",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
			MainField:       true,
			SchemaName:      "categories",
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
			Mandatory:       true,
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
		},
	)
}
