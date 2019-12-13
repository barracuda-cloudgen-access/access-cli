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
	"github.com/spf13/cobra"

	apiproxies "github.com/fyde/fyde-cli/client/access_proxies"
)

// proxyDeleteCmd represents the delete command
var proxyDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete proxies",
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
		proxyIDs := make([]strfmt.UUID, len(args))
		for i, arg := range args {
			proxyIDs[i] = strfmt.UUID(arg)
		}

		delete := func(ids []strfmt.UUID) error {
			params := apiproxies.NewDeleteProxyParams()
			params.SetID(ids)

			_, err := global.Client.AccessProxies.DeleteProxy(params, global.AuthWriter)
			if err != nil {
				return processErrorResponse(err)
			}
			return nil
		}

		tw, j := multiOpBuildTableWriter()

		var err error
		if loopControlContinueOnError(cmd) {
			// then we must delete individually, because on a request for multiple deletions,
			// the server does nothing if one fails

			for _, id := range proxyIDs {
				err = delete([]strfmt.UUID{id})
				var result interface{}
				result = "success"
				if err != nil {
					result = err
				}
				multiOpTableWriterAppend(tw, &j, id, result)
			}
			err = nil
		} else {
			err = delete(proxyIDs)
			var result interface{}
			result = "success"
			if err != nil {
				result = err
			}
			multiOpTableWriterAppend(tw, &j, "*", result)
		}

		return printListOutputAndError(cmd, j, tw, len(proxyIDs), err)
	},
}

func init() {
	proxiesCmd.AddCommand(proxyDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// proxyDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// proxyDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(proxyDeleteCmd)
	initLoopControlFlags(proxyDeleteCmd)
}
