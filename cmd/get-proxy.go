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

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apiproxies "github.com/barracuda-cloudgen-access/access-cli/client/access_proxies"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// proxyGetCmd represents the get command
var proxyGetCmd = &cobra.Command{
	Use:   "get [proxy ID]",
	Short: "Get proxy",
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
			return fmt.Errorf("missing proxy ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var id string
		var err error
		if cmd.Flags().Changed("id") {
			id, err = cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
		} else {
			id = args[0]
		}

		params := apiproxies.NewGetProxyParams()
		setTenant(cmd, params)
		params.SetID(strfmt.UUID(id))

		resp, err := global.Client.AccessProxies.GetProxy(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := proxyBuildTableWriter()
		proxyTableWriterAppend(tw, resp.Payload.AccessProxy, len(resp.Payload.AccessResources))

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func proxyBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Location",
		"Proxy host:port",
		"Resources",
		"Granted req.",
		"Total req.",
		"Last access",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 36},
		{Number: 2, WidthMax: 30},
		{Number: 3, WidthMax: 30},
		{Number: 4, WidthMax: 30},
		{Number: 5, WidthMax: 9},
		{Number: 6, WidthMax: 12},
		{Number: 7, WidthMax: 12},
		{Number: 8, WidthMax: 30},
	})

	return tw
}

func proxyTableWriterAppend(tw table.Writer, proxy models.AccessProxy, accessResourceCount int) {
	lastAccess := fmt.Sprint(proxy.LastAccessAt)
	if proxy.LastAccessAt == nil {
		lastAccess = "never"
	}

	granted := "?"
	total := "?"
	if proxy.AccessCount != nil {
		granted = fmt.Sprint(proxy.AccessCount.Granted)
		total = fmt.Sprint(proxy.AccessCount.Granted + proxy.AccessCount.Denied)
	}

	tw.AppendRow(table.Row{
		proxy.ID,
		proxy.Name,
		proxy.Location,
		fmt.Sprintf("%s:%d", proxy.Host, proxy.Port),
		accessResourceCount,
		granted,
		total,
		lastAccess,
	})
}

func proxyTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	idStr := "[ERR]"
	if id != nil {
		idStr += fmt.Sprintf(" %v", id)
	}
	tw.AppendRow(table.Row{
		idStr,
		processErrorResponse(err),
		"-",
		"-",
		"-",
		"-",
		"-",
		"-",
	})
}

func init() {
	proxiesCmd.AddCommand(proxyGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(proxyGetCmd)
	initTenantFlags(proxyGetCmd)
	proxyGetCmd.Flags().String("id", "", "id of proxy to get")
}
