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
	"os"
	"strings"

	api "github.com/barracuda-cloudgen-access/access-cli/client/settings_agent_configuration"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/barracuda-cloudgen-access/access-cli/serial"
	"github.com/spf13/cobra"
)

// setAgentConfigCmd represents the get command
var setAgentConfigCmd = &cobra.Command{
	Use:   "set",
	Short: "Set agent configuration",
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
		tw := agentConfigBuildTableWriter()
		createdList := []*api.EditSettingsAgentConfigurationBody{}
		total := 0
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := api.NewEditSettingsAgentConfigurationParams()
				setTenant(cmd, params)
				config := &models.SettingsAgentConfiguration{}
				err := placeInputValues(cmd, values, config,
					func(s string) {
						config.CertificatePeriod = &serial.NullableOptionalInt{}
						config.CertificatePeriod.AssignFromString(s)
					},
					func(s string) {
						config.ContactAdminAction = &serial.NullableOptionalString{
							Value: &s,
						}
					},
					func(s string) {
						config.EnrollmentPollingTimeInSeconds = &serial.NullableOptionalInt{}
						config.EnrollmentPollingTimeInSeconds.AssignFromString(s)
					},
					func(s bool) {
						config.HistoryScreenDisabled = &serial.NullableOptionalBoolean{
							Value: &s,
						}
					},
					func(s string) {
						parts := strings.SplitN(s, ":", 2)
						if len(parts) != 2 {
							fmt.Fprint(os.Stderr, "Error: Invalid dns_servers format.\n\n")
							cmd.Usage()
							os.Exit(1)
						}

						config.DNSServers = &models.SettingsAgentConfigurationDNSServers{
							Protocol: parts[0],
							List:     parts[1],
						}
					},
				)
				if err != nil {
					return nil, err
				}

				body := api.EditSettingsAgentConfigurationBody{AppConfiguration: config}
				params.SetAppConfiguration(body)

				resp, err := global.Client.SettingsAgentConfiguration.EditSettingsAgentConfiguration(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				config := data.(*models.SettingsAgentConfiguration)
				agentConfigTableWriterAppend(tw, *config)
			}, func(err error, id interface{}) { // doOnError func
				agentConfigTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	settingsAgentConfigCmd.AddCommand(setAgentConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setAgentConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setAgentConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(setAgentConfigCmd)
	initLoopControlFlags(setAgentConfigCmd)
	initTenantFlags(setAgentConfigCmd)
	initInputFlags(setAgentConfigCmd, "agent configuration",
		inputField{
			Name:            "Certificate Period",
			FlagName:        "certificate_period",
			FlagDescription: "specify period (in days) for the device certificate",
			VarType:         "string", // use string to read "null" pseudo value
			DefaultValue:    "30",
		},
		inputField{
			Name:            "Contact Admin Action",
			FlagName:        "contact_admin_action",
			FlagDescription: "Configured url action to contact admin. Example: mailto:support@acme.corp",
			VarType:         "string",
			DefaultValue:    "",
		},
		inputField{
			Name:            "Enrollment Polling Time",
			FlagName:        "enrollment_polling_time",
			FlagDescription: "Configured polling time to update agent settings (in seconds)",
			VarType:         "string", // use string to read "null" pseudo value
			DefaultValue:    "600",
		},
		inputField{
			Name:            "History Screen Disabled",
			FlagName:        "history_screen_disabled",
			FlagDescription: "Turn on to disable the CloudGen Access App history screen.",
			VarType:         "bool",
			DefaultValue:    false,
		},
		inputField{
			Name:            "DNS Servers List",
			FlagName:        "dns_servers",
			FlagDescription: "Enforce agent to use a DNS server config. Format: \"protocol:IP1,IP2\". Protocols: [plain, dns_over_tls]",
			VarType:         "string",
			DefaultValue:    "",
		},
	)
}
