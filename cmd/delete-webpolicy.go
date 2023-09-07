// Package cmd implements access-cli commands
package cmd

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiwebpolicies "github.com/barracuda-cloudgen-access/access-cli/client/web_policies"
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

// policyDeleteCmd represents the delete command
var webPolicyDeleteCmd = &cobra.Command{
	Use:     "delete [policy ID]...",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete web policies",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		if !multiOpCheckArgsPresent(cmd, args) {
			return fmt.Errorf("missing policy ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		gparams := apiwebpolicies.NewListWebPoliciesParams()
		setTenant(cmd, gparams)
		resp, err := global.Client.WebPolicies.ListWebPolicies(gparams, global.AuthWriter)

		if err != nil {
			return processErrorResponse(err)
		}
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
		policyIDs, err := multiOpParseUUIDArgs(cmd, args, "id")
		if err != nil {
			return err
		}
		params := apiwebpolicies.NewDeleteWebPolicyParams()
		setTenant(cmd, params)

		delete := func(id strfmt.UUID) error {
			if id == "" {
				fmt.Println("Skip empty id")
				return nil
			}
			if mapRules[id] != "" {
				id = mapRules[id]
				fmt.Println("Delete jumprule ", id)
			}

			params.SetID(id)
			_, err = global.Client.WebPolicies.DeleteWebPolicy(params, global.AuthWriter)
			if err != nil {
				return processErrorResponse(err)
			}
			return nil
		}

		tw, j := multiOpBuildTableWriter()

		for _, id := range policyIDs {
			if loopControlContinueOnError(cmd) {
				// Note errors and continue on to the next item
				for _, id := range policyIDs {
					err = delete(id)
					var result interface{}
					result = "success"
					if err != nil {
						result = err
					}
					multiOpTableWriterAppend(tw, &j, id, result)
				}
				err = nil
			} else {
				//Note error and fail
				err = delete(id)
				var result interface{}
				result = "success"
				if err != nil {
					result = err
				}
				multiOpTableWriterAppend(tw, &j, id, result)
				if err != nil {
					return err
				}
			}
		}

		return printListOutputAndError(cmd, j, tw, len(policyIDs), err)
	},
}

func init() {
	webPoliciesCmd.AddCommand(webPolicyDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policyDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//policyDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initMultiOpArgFlags(webPolicyDeleteCmd, "webpolicy", "delete", "id", "[]strfmt.UUID")

	initOutputFlags(webPolicyDeleteCmd)
	initLoopControlFlags(webPolicyDeleteCmd)
	initTenantFlags(webPolicyDeleteCmd)
}
