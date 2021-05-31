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
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apisources "github.com/barracuda-cloudgen-access/access-cli/client/asset_sources"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// sourcesListCmd represents the list command
var sourcesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List sources",
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
		params := apisources.NewListAssetSourcesParams()
		setTenant(cmd, params)
		completePayload := []*models.AssetSource{}
		total := 0
		cutStart, cutEnd, err := forAllPages(cmd, params, func() (int, int64, error) {
			resp, err := global.Client.AssetSources.ListAssetSources(params, global.AuthWriter)
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
			"Domains",
			"Modified",
			"Enabled",
			"Read-only",
		})
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 36},
			{Number: 2, WidthMax: 30},
			{Number: 3, WidthMax: 10},
			{Number: 4, WidthMax: 30},
			{Number: 5, WidthMax: 10},
			{Number: 6, WidthMax: 10},
		})

		for _, item := range completePayload {
			assetsCount := item.AssetsCount
			if item.ReadOnlyCount != nil {
				assetsCount -= *item.ReadOnlyCount
			}
			var assetsCountStr string
			if item.ReadOnlyCount != nil && *item.ReadOnlyCount > 0 && assetsCount == 0 {
				assetsCountStr = "-"
			} else {
				assetsCountStr = strconv.Itoa(int(assetsCount))
			}
			tw.AppendRow(table.Row{
				item.ID,
				item.Name,
				assetsCountStr,
				item.ModifiedAt,
				item.Enabled,
				item.ReadOnly,
			})
		}

		return printListOutputAndError(cmd, completePayload, tw, total, err)
	},
}

func init() {
	sourcesCmd.AddCommand(sourcesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sourcesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sourcesListCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initPaginationFlags(sourcesListCmd)
	initOutputFlags(sourcesListCmd)
	initTenantFlags(sourcesListCmd)
}
