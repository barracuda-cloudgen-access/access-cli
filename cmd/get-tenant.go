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
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apitenants "github.com/barracuda-cloudgen-access/access-cli/client/tenants"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// tenantGetCmd represents the get command
var tenantGetCmd = &cobra.Command{
	Use:   "get [tenant ID]",
	Short: "Get tenant",
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
			return fmt.Errorf("missing tenant ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var tenantID string
		var err error
		if cmd.Flags().Changed("id") {
			tenantID, err = cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
		} else {
			tenantID = args[0]
		}

		params := apitenants.NewGetTenantParams()
		params.SetID(strfmt.UUID(tenantID))

		resp, err := global.Client.Tenants.GetTenant(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := tenantBuildTableWriter()
		tenantTableWriterAppend(tw, *resp.Payload)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func tenantBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Created At",
		"Updated At",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 36, Align: text.AlignRight},
		{Number: 2, WidthMax: 30, Align: text.AlignRight},
		{Number: 3, WidthMax: 30, Align: text.AlignLeft},
		{Number: 4, WidthMax: 30, Align: text.AlignLeft},
	})
	return tw
}

func tenantTableWriterAppend(tw table.Writer, tenant models.Tenant) {
	tw.AppendRow(table.Row{
		tenant.ID,
		tenant.Name,
		tenant.CreatedAt,
		tenant.UpdatedAt,
	})
}

func init() {
	tenantsCmd.AddCommand(tenantGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tenantGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tenantGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(tenantGetCmd)
	tenantGetCmd.Flags().Int("id", 0, "id of tenant to get")
}
