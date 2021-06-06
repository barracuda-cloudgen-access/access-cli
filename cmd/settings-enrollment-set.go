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
	api "github.com/barracuda-cloudgen-access/access-cli/client/settings_enrollment"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/spf13/cobra"
)

// setEnrollmentCmd represents the get command
var setEnrollmentCmd = &cobra.Command{
	Use:   "set",
	Short: "Set enrollment settings",
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
		tw := enrollmentBuildTableWriter()
		createdList := []*api.EditSettingsEnrollmentBody{}
		total := 0
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := api.NewEditSettingsEnrollmentParams()
				setTenant(cmd, params)
				config := &models.SettingsEnrollment{}
				err := placeInputValues(cmd, values, config,
					func(s int) {
						config.ExpirationDays = int64(s)
					},
					func(s int) {
						config.Refcount = int64(s)
					},
				)
				if err != nil {
					return nil, err
				}

				body := api.EditSettingsEnrollmentBody{EnrollmentSettings: config}
				params.SetEnrollmentSettings(body)

				resp, err := global.Client.SettingsEnrollment.EditSettingsEnrollment(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				config := data.(*models.SettingsEnrollment)
				enrollmentTableWriterAppend(tw, *config)
			}, func(err error, id interface{}) { // doOnError func
				enrollmentTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	settingsEnrollmentCmd.AddCommand(setEnrollmentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setEnrollmentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setEnrollmentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(setEnrollmentCmd)
	initLoopControlFlags(setEnrollmentCmd)
	initTenantFlags(setEnrollmentCmd)
	initInputFlags(setEnrollmentCmd, "enrollment",
		inputField{
			Name:            "Expiration",
			FlagName:        "expiration",
			FlagDescription: "specify period (in days) for the expiration of the enrollment link",
			VarType:         "int",
			DefaultValue:    0,
		},
		inputField{
			Name:            "Refcount",
			FlagName:        "refcount",
			FlagDescription: "Available slots for enrollment",
			VarType:         "int",
			DefaultValue:    0,
		},
	)
}
