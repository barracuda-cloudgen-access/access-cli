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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	api "github.com/barracuda-cloudgen-access/access-cli/client/settings_analytics"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// getAnalyticsCmd represents the get command
var getAnalyticsCmd = &cobra.Command{
	Use:   "get",
	Short: "Get analytics settings",
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
		params := api.NewSettingsAnalyticsParams()
		setTenant(cmd, params)

		cmd.SilenceUsage = true // errors beyond this point are no longer due to malformed input

		resp, err := global.Client.SettingsAnalytics.SettingsAnalytics(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := analyticsBuildTableWriter()
		analyticsTableWriterAppend(tw, *resp.Payload)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func analyticsBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"URL",
		"DisableSSL",
		"InterceptAllDomains",
	})

	return tw
}

func analyticsTableWriterAppend(tw table.Writer, config models.SettingsAnalytics) table.Writer {
	if config.ExternalServer.URL == "" {
		tw.AppendRow(table.Row{
			"<disabled>",
		})
	} else {
		tw.AppendRow(table.Row{
			config.ExternalServer.URL,
			config.ExternalServer.DisableSsl,
			config.ExternalServer.InterceptAllDomains,
		})
	}
	return tw
}

func analyticsTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	tw.AppendRow(table.Row{
		"[ERR]",
		processErrorResponse(err),
		"-",
		"-",
		"-",
		"-",
	})
}

func init() {
	settingsAnalyticsCmd.AddCommand(getAnalyticsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getAnalyticsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getAnalyticsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(getAnalyticsCmd)
	initTenantFlags(getAnalyticsCmd)
}
