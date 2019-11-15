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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiproxies "github.com/fyde/fyde-cli/client/access_proxies"
	"github.com/fyde/fyde-cli/models"
)

// proxiesEditCmd represents the edit command
var proxiesEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit proxies",
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
		tw := proxyBuildTableWriter()
		createdList := []*apiproxies.EditProxyOKBody{}
		total := 0
		err := forAllInput(cmd, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apiproxies.NewEditProxyParams()
				// IDs are not part of the request body, so we use this workaround
				proxy := &struct {
					models.AccessProxy
					ID string `json:"id"`
				}{}
				err := placeInputValues(cmd, values, proxy,
					func(s string) { proxy.ID = s },
					func(s string) { proxy.AccessProxy.Name = s },
					func(s string) { proxy.AccessProxy.Location = s },
					func(s string) { proxy.AccessProxy.Host = s },
					func(s int) { proxy.AccessProxy.Port = int64(s) })
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(strfmt.UUID(proxy.ID))
				body := apiproxies.EditProxyBody{}
				body.AccessProxy.AccessProxy = proxy.AccessProxy
				params.SetProxy(body)

				resp, err := global.Client.AccessProxies.EditProxy(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				proxy := data.(*apiproxies.EditProxyOKBody)
				createdList = append(createdList, proxy)
				proxyTableWriterAppend(tw, proxy.AccessProxy, len(proxy.AccessResources))
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				proxyTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	proxiesCmd.AddCommand(proxiesEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxiesEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxiesEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(proxiesEditCmd)
	initLoopControlFlags(proxiesEditCmd)

	initInputFlags(proxiesEditCmd, "proxy",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the proxy to edit",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			IsIDOnError:     true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the new name for the proxy",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			IsIDOnError:     true,
			SchemaName:      "name",
		},
		inputField{
			Name:            "Location",
			FlagName:        "location",
			FlagDescription: "specify the new location for the proxy",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			IsIDOnError:     true,
			SchemaName:      "name",
		},
		inputField{
			Name:            "Host",
			FlagName:        "host",
			FlagDescription: "specify the new host for the proxy",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Port",
			FlagName:        "port",
			FlagDescription: "specify the new port for the proxy",
			VarType:         "int",
			Mandatory:       false,
			DefaultValue:    0,
		})
}
