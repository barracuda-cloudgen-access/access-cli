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
	"fmt"
	"strconv"
	"strings"

	apiwebpolicies "github.com/barracuda-cloudgen-access/access-cli/client/web_policies"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

var webPoliciesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List web policies",
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
		params := apiwebpolicies.NewListWebPoliciesParams()
		setTenant(cmd, params)
		resp, err := global.Client.WebPolicies.ListWebPolicies(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}
		tw := webPolicyBuildTableWriter()
		for _, rule := range resp.Payload.Rules {
			for _, policy := range rule.RulesetJump.Rules {
				webPolicyTableWriterAppend(tw, rule.ID, *policy)
			}
		}
		return printListOutputAndError(cmd, resp.Payload, tw, 0, err)
	},
}

func webPolicyBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Label",
		"Action",
		"Type",
		"Categories",
		"Domains",
		"Disabled",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 36, Align: text.AlignLeft},
		{Number: 2, WidthMax: 16, Align: text.AlignLeft},
		{Number: 3, WidthMax: 6, Align: text.AlignLeft},
		{Number: 4, WidthMax: 8, Align: text.AlignLeft},
		{Number: 5, WidthMax: 20, Align: text.AlignLeft},
		{Number: 6, WidthMax: 20, Align: text.AlignLeft},
		{Number: 7, WidthMax: 8, Align: text.AlignLeft},
	})
	return tw
}

func webPolicyTableWriterAppend(tw table.Writer, rulesetID strfmt.UUID, policy models.GetWebPolicyRule) {

	var cats []string
	for _, x := range policy.Categories {
		newId := int(x)
		if int64(newId) != x {
			panic("overflows!")
		}
		cats = append(cats, strconv.Itoa(newId)) // note the = instead of :=
	}
	categories := strings.Join(cats, ",")
	domains := strings.Join(policy.Domains, ", ")
	tw.AppendRow(table.Row{
		policy.ID,
		policy.Label,
		policy.Action,
		policy.Type,
		categories,
		domains,
		policy.Disabled,
	})
}
func setWebPolicyBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Label",
		"Action",
		"Type",
		"Categories",
		"Domains",
		"Disabled",
		"UserIds",
		"GroupIds",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 36, Align: text.AlignRight},
		{Number: 2, WidthMax: 16, Align: text.AlignLeft},
		{Number: 3, WidthMax: 6, Align: text.AlignLeft},
		{Number: 4, WidthMax: 8, Align: text.AlignLeft},
		{Number: 5, WidthMax: 20, Align: text.AlignLeft},
		{Number: 6, WidthMax: 20, Align: text.AlignLeft},
		{Number: 7, WidthMax: 8, Align: text.AlignLeft},
		{Number: 8, WidthMax: 8, Align: text.AlignLeft},
		{Number: 9, WidthMax: 8, Align: text.AlignLeft},
	})
	return tw
}
func setWebPolicyTableWriterAppend(tw table.Writer, ruleID strfmt.UUID, policy models.WebPolicyData) {

	var cats []string
	for _, x := range policy.Categories {
		newId := int(x)
		if int64(newId) != x {
			panic("overflows!")
		}
		cats = append(cats, strconv.Itoa(newId)) // note the = instead of :=
	}
	categories := strings.Join(cats, ",")
	domains := strings.Join(policy.Domains, ", ")
	var user_ids []string
	for _, x := range policy.UserIds {
		newId := int(x)
		if int64(newId) != x {
			panic("overflows!")
		}
		user_ids = append(user_ids, strconv.Itoa(newId)) // note the = instead of :=
	}
	var group_ids []string
	for _, x := range policy.GroupIds {
		newId := int(x)
		if int64(newId) != x {
			panic("overflows!")
		}
		group_ids = append(group_ids, strconv.Itoa(newId)) // note the = instead of :=
	}
	tw.AppendRow(table.Row{
		ruleID,
		policy.Label,
		policy.Action,
		policy.Type,
		categories,
		domains,
		policy.Disabled,
		strings.Join(user_ids, ", "),
		strings.Join(group_ids, ", "),
	})
}

func webPolicyTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	idStr := "[ERR]"
	if id != nil {
		idStr += fmt.Sprintf(" %v", id)
	}
	tw.AppendRow(table.Row{
		idStr,
		processErrorResponse(err),
		"-",
		"-",
	})
}
func init() {
	webPoliciesCmd.AddCommand(webPoliciesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policiesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policiesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(webPoliciesListCmd)
	initTenantFlags(webPoliciesListCmd)
}
