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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apiassets "github.com/barracuda-cloudgen-access/access-cli/client/assets"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// domainsListCmd represents the list command
var domainsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List domains",
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
		params := apiassets.NewListAssetsParams()
		setTenant(cmd, params)
		setSort(cmd, params)
		setFilter(cmd, params.SetCategory)
		completePayload := []*models.Asset{}
		total := 0
		cutStart, cutEnd, err := forAllPages(cmd, params, func() (int, int64, error) {
			resp, err := global.Client.Assets.ListAssets(params, global.AuthWriter)
			if err != nil {
				return 0, 0, err
			}
			completePayload = append(completePayload, resp.Payload...)
			total = int(resp.Total)
			return len(resp.Payload), resp.Total, err
		})
		if err != nil {
			return processErrorResponse(err)
		}
		completePayload = completePayload[cutStart:cutEnd]

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"Name",
			"Category",
			"Asset source",
		})
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 15},
			{Number: 2, WidthMax: 30},
			{Number: 3, WidthMax: 30},
			{Number: 4, WidthMax: 36},
		})
		for _, item := range completePayload {
			tw.AppendRow(table.Row{
				item.ID,
				item.Name,
				item.Category,
				item.AssetSourceID,
			})
		}

		return printListOutputAndError(cmd, completePayload, tw, total, err)
	},
}

func init() {
	domainsCmd.AddCommand(domainsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainsListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainsListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initPaginationFlags(domainsListCmd)
	initSortFlags(domainsListCmd)
	initFilterFlags(domainsListCmd,
		filterType{"category", "string"})
	initOutputFlags(domainsListCmd)
	initTenantFlags(domainsListCmd)
}
