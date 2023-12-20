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

	api "github.com/barracuda-cloudgen-access/access-cli/client/settings_agent_configuration"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// getAgentConfigCmd represents the get command
var getAgentConfigCmd = &cobra.Command{
	Use:   "get",
	Short: "Get agent configuration",
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

		params := api.NewSettingsAgentConfigurationParams()
		setTenant(cmd, params)

		resp, err := global.Client.SettingsAgentConfiguration.SettingsAgentConfiguration(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := agentConfigBuildTableWriter()
		agentConfigTableWriterAppend(tw, *resp.Payload)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func agentConfigBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"CertificatePeriod",
		"ContactAdminAction",
		"DnsServers",
		"EnrollmentPollingTime",
		"HistoryScreenDisabled",
	})

	return tw
}

func agentConfigTableWriterAppend(tw table.Writer, config models.SettingsAgentConfiguration) table.Writer {
	dnsConfig := "<null>"
	if config.DNSServers.List != "" {
		dnsConfig = config.DNSServers.Protocol.(string) + ": " + config.DNSServers.List
	}
	tw.AppendRow(table.Row{
		config.CertificatePeriod,
		config.ContactAdminAction,
		dnsConfig,
		config.EnrollmentPollingTimeInSeconds,
		config.HistoryScreenDisabled,
	})
	return tw
}

func agentConfigTableWriterAppendError(tw table.Writer, err error, id interface{}) {
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
	settingsAgentConfigCmd.AddCommand(getAgentConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getAgentConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getAgentConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(getAgentConfigCmd)
	initTenantFlags(getAgentConfigCmd)
}
