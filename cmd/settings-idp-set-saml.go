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
	api "github.com/barracuda-cloudgen-access/access-cli/client/identity_providers"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/spf13/cobra"
)

// setSamlIdpCmd represents the get command
var setSamlIdpCmd = &cobra.Command{
	Use:   "saml",
	Short: "Set SAML idp configuration",
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
		tw := identityProviderConfigBuildTableWriter()
		createdList := []*models.IdentityProvider{}
		total := 0

		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := api.NewCreateIdentityProviderParams()

				idp := &api.CreateIdentityProviderParamsBodyIdentityProvider{
					IdpType: "saml",
					Details: map[string]interface{}{},
				}

				err := placeInputValues(cmd, values, idp,
					func(s string) { idp.Details["entity_id"] = s },
					func(s string) { idp.Details["sso_url"] = s },
					func(s string) { idp.Details["certificate"] = s })
				if err != nil {
					return nil, err
				}

				setTenant(cmd, params)

				body := api.CreateIdentityProviderBody{
					IdentityProvider: idp,
				}
				params.SetIdentityProvider(body)

				resp, err := global.Client.IdentityProviders.CreateIdentityProvider(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}

				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				idp := data.(*models.IdentityProvider)
				if idp.ID > 0 {
					identityProviderTableWriterAppend(tw, *idp)
				}
				createdList = append(createdList, idp)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				identityProviderTableWriterAppendError(tw, err, id)
			})

		return printListOutputAndError(cmd, createdList, tw, 1, err)
	},
}

func init() {
	setIdpCmd.AddCommand(setSamlIdpCmd)

	initOutputFlags(setSamlIdpCmd)
	initLoopControlFlags(setSamlIdpCmd)
	initTenantFlags(setSamlIdpCmd)

	initInputFlags(setSamlIdpCmd, "idp",
		inputField{
			Name:            "Entity ID",
			FlagName:        "entity_id",
			FlagDescription: "SAML SSO Entity ID",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "SSO URL",
			FlagName:        "sso_url",
			FlagDescription: "SAML SSO Redirect URL",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Certificate",
			FlagName:        "certificate",
			FlagDescription: "Certificate of the SAML provider",
			VarType:         "string+file",
			Mandatory:       true,
			DefaultValue:    "",
		})
}
