// Package cmd implements access-cli commands
package cmd

import (
	"strconv"

	apiwebcategories "github.com/barracuda-cloudgen-access/access-cli/client/web_categories"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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
// policiesListCmd represents the list command
var webCategoriesListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List web categories",
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
		params := apiwebcategories.NewListWebCategoriesParams()
		params.SetVersion("3.2")
		setTenant(cmd, params)
		resp, err := global.Client.WebCategories.ListWebCategories(params, global.AuthWriter)

		if err != nil {
			return processErrorResponse(err)
		}
		tw := webCategoryBuildTableWriter()

		categories := resp.Payload.Categories
		supercategories := resp.Payload.Supercategories

		for categoryId, categoryObj := range categories {
			displayName := categoryObj.Display
			parent := strconv.FormatInt(int64(categoryObj.Parent), 10)
			parent = supercategories[parent]
			webCategoryTableWriterAppend(tw, categoryId, displayName, parent)
		}

		return printListOutputAndError(cmd, resp.Payload, tw, 0, err)
	},
}

func webCategoryBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Category",
		"Super Category",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 8, Align: text.AlignRight},
		{Number: 2, WidthMax: 36, Align: text.AlignLeft},
		{Number: 3, WidthMax: 36, Align: text.AlignLeft},
	})
	return tw
}

func webCategoryTableWriterAppend(tw table.Writer, categoryId string, displayName string, parent string) {

	tw.AppendRow(table.Row{
		categoryId,
		displayName,
		parent,
	})
}

func init() {
	webCategoriesCmd.AddCommand(webCategoriesListCmd)

	initOutputFlags(webCategoriesListCmd)
	initTenantFlags(webCategoriesListCmd)
}
