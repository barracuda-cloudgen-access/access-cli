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
	"time"

	apiwebcategories "github.com/barracuda-cloudgen-access/access-cli/client/web_categories"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"
)

// userGetCmd represents the get command
var webCategoryGetCmd = &cobra.Command{
	Use:     "get [domain]...",
	Aliases: []string{"domain"},
	Short:   "Get web categories for domains",
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
			return fmt.Errorf("missing domain argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		domains, err := multiOpParseStringArgs(cmd, args, "domain")
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		params := apiwebcategories.QueryWebCategoriesParams{}
		params.WithTimeout(30 * time.Second)
		setTenant(cmd, &params)
		params.SetDomains(domains)
		resp, err := global.Client.WebCategories.QueryWebCategories(&params, global.AuthWriter)

		if err != nil {
			return processErrorResponse(err)
		}

		tw := domainLookupBuildTableWriter()
		for _, item := range resp.Payload.Domains {
			domainLookupTableWriterAppend(tw, *item)
		}
		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func domainLookupBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"Domain",
		"Categories",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 30, Align: text.AlignLeft},
		{Number: 2, WidthMax: 30, Align: text.AlignLeft},
	})
	return tw
}

func domainLookupTableWriterAppend(tw table.Writer, item models.DomainLookupResultItem) {
	tw.AppendRow(table.Row{
		item.Domain,
		item.Categories,
	})
}

func init() {
	webCategoriesCmd.AddCommand(webCategoryGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initMultiOpArgFlags(webCategoryGetCmd, "webcategory", "get", "domain", "[]string")

	initOutputFlags(webCategoryGetCmd)
	initTenantFlags(webCategoryGetCmd)
}
