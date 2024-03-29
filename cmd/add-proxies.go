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

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apiproxies "github.com/barracuda-cloudgen-access/access-cli/client/access_proxies"
)

// proxiesAddCmd represents the add command
var proxiesAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add proxies",
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
		tw := proxyBuildTableWriterForCreation()
		createdList := []*apiproxies.CreateProxyCreatedBody{}
		total := 0
		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				proxy := &apiproxies.CreateProxyBody{}
				err := placeInputValues(cmd, values, proxy,
					func(s string) { proxy.Name = s },
					func(s string) { proxy.Location = s },
					func(s string) { proxy.Host = s },
					func(s int) { proxy.Port = int64(s) })
				if err != nil {
					return nil, err
				}
				params := apiproxies.NewCreateProxyParams()
				setTenant(cmd, params)
				params.SetProxy(*proxy)

				resp, err := global.Client.AccessProxies.CreateProxy(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				proxy := data.(*apiproxies.CreateProxyCreatedBody)
				createdList = append(createdList, proxy)
				proxyTableWriterAppendForCreation(tw, proxy)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				proxyTableWriterAppendErrorForCreation(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func proxyBuildTableWriterForCreation() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Location",
		"Proxy host:port",
		"Enrollment URL",
	})
	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 36},
		{Number: 2, WidthMax: 30},
		{Number: 3, WidthMax: 30},
		{Number: 4, WidthMax: 30},
		{Number: 5, WidthMax: 60},
	})
	return tw
}

func proxyTableWriterAppendForCreation(tw table.Writer, proxy *apiproxies.CreateProxyCreatedBody) {
	tw.AppendRow(table.Row{
		proxy.ID,
		proxy.Name,
		proxy.Location,
		fmt.Sprintf("%s:%d", proxy.Host, proxy.Port),
		proxy.EnrollmentURL,
	})
}

func proxyTableWriterAppendErrorForCreation(tw table.Writer, err error, id interface{}) {
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
	})
}

func init() {
	proxiesCmd.AddCommand(proxiesAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxiesAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxiesAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(proxiesAddCmd)
	initLoopControlFlags(proxiesAddCmd)
	initTenantFlags(proxiesAddCmd)
	initInputFlags(proxiesAddCmd, "proxy",
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the name for the created proxy",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "name",
		},
		inputField{
			Name:            "Location",
			FlagName:        "location",
			FlagDescription: "specify the location for the created proxy",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Host",
			FlagName:        "host",
			FlagDescription: "specify the host for the created proxy",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Port",
			FlagName:        "port",
			FlagDescription: "specify the port for the created proxy",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    0,
		})
}
