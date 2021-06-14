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
	"github.com/spf13/cobra"

	api "github.com/barracuda-cloudgen-access/access-cli/client/identity_providers"
)

// listIdentityProviderCmd represents the get command
var listIdentityProviderCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List IdP configurations",
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
		params := api.NewListIdentityProviderParams()
		setTenant(cmd, params)

		resp, err := global.Client.IdentityProviders.ListIdentityProvider(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := identityProviderConfigBuildTableWriter()
		if resp.Payload.ID > 0 {
			identityProviderTableWriterAppend(tw, *resp.Payload)
		}

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func init() {
	settingsIdentityProviderConfigCmd.AddCommand(listIdentityProviderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listIdentityProviderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listIdentityProviderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(listIdentityProviderCmd)
	initTenantFlags(listIdentityProviderCmd)
}
