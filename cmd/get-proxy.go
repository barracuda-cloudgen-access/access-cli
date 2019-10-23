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
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apiproxies "github.com/fyde/fyde-cli/client/access_proxies"
)

// proxyGetCmd represents the get command
var proxyGetCmd = &cobra.Command{
	Use:   "get",
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

		if len(args) == 0 {
			return fmt.Errorf("missing proxy ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		params := apiproxies.NewGetProxyParams()
		params.SetID(strfmt.UUID(args[0]))

		resp, err := global.Client.AccessProxies.GetProxy(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

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
		tw.SetAllowedColumnLengths([]int{36, 30, 30, 30, 9, 12, 12, 30})

		lastAccess := fmt.Sprint(resp.Payload.LastAccessAt)
		if resp.Payload.LastAccessAt == nil {
			lastAccess = "never"
		}

		tw.AppendRow(table.Row{
			resp.Payload.ID,
			resp.Payload.Name,
			resp.Payload.Location,
			fmt.Sprintf("%s:%d", resp.Payload.Host, resp.Payload.Port),
			len(resp.Payload.AccessResources),
			resp.Payload.AccessCount.Granted,
			resp.Payload.AccessCount.Granted + resp.Payload.AccessCount.Denied,
			lastAccess,
		})

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
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
}
