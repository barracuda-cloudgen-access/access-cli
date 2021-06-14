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
	"fmt"

	api "github.com/barracuda-cloudgen-access/access-cli/client/identity_providers"

	"github.com/spf13/cobra"
)

// deleteIdpCmd represents the get command
var deleteIdpCmd = &cobra.Command{
	Use:     "delete [idp ID]...",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete idps",
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
			return fmt.Errorf("missing idp ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		adminIDs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		delete := func(ids []int64) error {
			params := api.NewDeleteIdentityProviderParams()
			setTenant(cmd, params)
			params.SetID(ids)

			_, err = global.Client.IdentityProviders.DeleteIdentityProvider(params, global.AuthWriter)
			if err != nil {
				return processErrorResponse(err)
			}
			return nil
		}

		tw, j := multiOpBuildTableWriter()

		if loopControlContinueOnError(cmd) {
			// then we must delete individually, because on a request for multiple deletions,
			// the server does nothing if one fails

			for _, id := range adminIDs {
				err = delete([]int64{id})
				var result interface{}
				result = "success"
				if err != nil {
					result = err
				}
				multiOpTableWriterAppend(tw, &j, id, result)
			}
			err = nil
		} else {
			err = delete(adminIDs)
			var result interface{}
			result = "success"
			if err != nil {
				result = err
			}
			multiOpTableWriterAppend(tw, &j, "*", result)
		}

		return printListOutputAndError(cmd, j, tw, len(adminIDs), err)
	},
}

func init() {
	settingsIdentityProviderConfigCmd.AddCommand(deleteIdpCmd)

	initMultiOpArgFlags(deleteIdpCmd, "idp", "delete", "id", "[]int64")
	initOutputFlags(deleteIdpCmd)
	initLoopControlFlags(deleteIdpCmd)
	initTenantFlags(deleteIdpCmd)
}
