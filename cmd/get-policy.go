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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apipolicies "github.com/barracuda-cloudgen-access/access-cli/client/access_policies"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// policyGetCmd represents the get command
var policyGetCmd = &cobra.Command{
	Use:   "get [policy ID]",
	Short: "Get policy",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		if len(args) == 0 && !cmd.Flags().Changed("id") {
			return fmt.Errorf("missing policy ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var policyID int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			policyID = int64(d)
		} else {
			policyID, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		params := apipolicies.NewGetPolicyParams()
		setTenant(cmd, params)
		params.SetID(policyID)

		resp, err := global.Client.AccessPolicies.GetPolicy(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := policyBuildTableWriter()
		policyTableWriterAppend(tw, resp.Payload.AccessPolicy, len(resp.Payload.AccessResources))

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func policyBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Resources",
		"Created",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 12, Align: text.AlignRight},
		{Number: 2, WidthMax: 30, Align: text.AlignLeft},
		{Number: 3, WidthMax: 12, Align: text.AlignLeft},
		{Number: 4, WidthMax: 30, Align: text.AlignLeft},
	})
	return tw
}

func policyTableWriterAppend(tw table.Writer, policy models.AccessPolicy, accessResourcesCount interface{}) {
	tw.AppendRow(table.Row{
		policy.ID,
		policy.Name,
		accessResourcesCount,
		policy.CreatedAt,
	})
}

func policyTableWriterAppendError(tw table.Writer, err error, id interface{}) {
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
	policiesCmd.AddCommand(policyGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// policyGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// policyGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(policyGetCmd)
	initTenantFlags(policyGetCmd)
	policyGetCmd.Flags().Int("id", 0, "id of policy to get")
}
