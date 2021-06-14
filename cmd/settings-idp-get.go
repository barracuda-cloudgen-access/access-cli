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
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	api "github.com/barracuda-cloudgen-access/access-cli/client/identity_providers"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// getIdentityProviderCmd represents the get command
var getIdentityProviderCmd = &cobra.Command{
	Use:   "get [idp ID]",
	Short: "Get IdP configuration",
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
			return fmt.Errorf("missing user ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			id = int64(d)
		} else {
			id, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true // errors beyond this point are no longer due to malformed input

		params := api.NewGetIdentityProviderParams()
		setTenant(cmd, params)
		params.SetID(id)

		resp, err := global.Client.IdentityProviders.GetIdentityProvider(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := identityProviderConfigBuildTableWriter()
		if resp.Payload.IdentityProvider.ID > 0 {
			identityProviderTableWriterAppend(tw, resp.Payload.IdentityProvider)
		}

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func identityProviderConfigBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Type",
		"Name",
		"CreatedAt",
		"UpdatedAt",
	})

	return tw
}

func identityProviderTableWriterAppend(tw table.Writer, idp models.IdentityProvider) table.Writer {
	tw.AppendRow(table.Row{
		idp.ID,
		idp.IdpType,
		idp.Name,
		idp.CreatedAt,
		idp.UpdatedAt,
	})
	return tw
}

func identityProviderTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	tw.AppendRow(table.Row{
		"[ERR]",
		processErrorResponse(err),
		"-",
		"-",
		"-",
	})
}

func init() {
	settingsIdentityProviderConfigCmd.AddCommand(getIdentityProviderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getIdentityProviderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getIdentityProviderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(getIdentityProviderCmd)
	initTenantFlags(getIdentityProviderCmd)

	getIdentityProviderCmd.Flags().Int("id", 0, "id of user to get")
}
