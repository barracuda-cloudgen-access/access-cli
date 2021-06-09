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
	api "github.com/barracuda-cloudgen-access/access-cli/client/settings_analytics"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/spf13/cobra"
)

// setAnalyticsCmd represents the get command
var setAnalyticsCmd = &cobra.Command{
	Use:   "set",
	Short: "Set analytics settings",
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
		tw := analyticsBuildTableWriter()
		createdList := []*api.EditSettingsAnalyticsBody{}
		total := 0
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := api.NewEditSettingsAnalyticsParams()
				setTenant(cmd, params)
				config := &models.SettingsAnalyticsExternalServer{}
				err := placeInputValues(cmd, values, config,
					func(s string) { config.URL = s },
					func(s bool) { config.DisableSsl = s },
					func(s bool) { config.InterceptAllDomains = s },
				)
				if err != nil {
					return nil, err
				}

				body := api.EditSettingsAnalyticsBody{
					AnalyticsSettings: &models.SettingsAnalytics{
						ExternalServer: config,
					},
				}
				params.SetAnalyticsSettings(body)

				resp, err := global.Client.SettingsAnalytics.EditSettingsAnalytics(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				config := data.(*models.SettingsAnalytics)
				analyticsTableWriterAppend(tw, *config)
			}, func(err error, id interface{}) { // doOnError func
				analyticsTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	settingsAnalyticsCmd.AddCommand(setAnalyticsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setAnalyticsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setAnalyticsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(setAnalyticsCmd)
	initLoopControlFlags(setAnalyticsCmd)
	initTenantFlags(setAnalyticsCmd)
	initInputFlags(setAnalyticsCmd, "agent configuration",
		inputField{
			Name:            "URL",
			FlagName:        "url",
			FlagDescription: "Analytics server host url. Use empty string to disable.",
			VarType:         "string",
			DefaultValue:    "",
		},
		inputField{
			Name:            "Disable SSL",
			FlagName:        "disable_ssl",
			FlagDescription: "Skips SSL server certificate checking for HTTPS events. WARNING: only use for development purposes.",
			VarType:         "bool",
			DefaultValue:    false,
		},
		inputField{
			Name:            "Intercept All Domains",
			FlagName:        "intercept_all_domains",
			FlagDescription: "Logs all domain resolutions. Note: think before enabling, since this is an intrusive setting.",
			VarType:         "bool",
			DefaultValue:    false,
		},
	)
}
