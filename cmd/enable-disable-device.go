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
	"strings"

	"github.com/spf13/cobra"

	apidevices "github.com/barracuda-cloudgen-access/access-cli/client/devices"
)

// deviceEnableCmd represents the enable command
var deviceEnableCmd = &cobra.Command{
	Use:   "enable [source ID]...",
	Short: "Enable device",
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
			return fmt.Errorf("missing device ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		enable := strings.HasPrefix(cmd.Use, "enable")

		uuidArgs, err := multiOpParseUUIDArgs(cmd, args, "id")
		if err != nil {
			return err
		}

		tw, j := multiOpBuildTableWriter()

		for _, arg := range uuidArgs {
			params := apidevices.NewEditDeviceParams()
			setTenant(cmd, params)
			params.SetID(arg)
			params.SetDevice(apidevices.EditDeviceBody{
				Device: &apidevices.EditDeviceParamsBodyDevice{
					Enabled: &enable,
				},
			})

			_, err = global.Client.Devices.EditDevice(params, global.AuthWriter)
			if err != nil {
				multiOpTableWriterAppend(tw, &j, arg, processErrorResponse(err))
				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, j, tw, len(uuidArgs), err)
			}
			multiOpTableWriterAppend(tw, &j, arg, "success")
		}
		return printListOutputAndError(cmd, j, tw, len(uuidArgs), err)
	},
}

// deviceDisableCmd represents the disable command
var deviceDisableCmd *cobra.Command

func init() {
	disableCmd := *deviceEnableCmd
	disableCmd.Use = "disable [source ID]..."
	disableCmd.Short = "Disable device"
	deviceDisableCmd = &disableCmd
	devicesCmd.AddCommand(deviceEnableCmd)
	devicesCmd.AddCommand(deviceDisableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deviceEnableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceEnableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initMultiOpArgFlags(deviceEnableCmd, "device", "enable", "id", "[]strfmt.UUID")
	initMultiOpArgFlags(deviceDisableCmd, "device", "disable", "id", "[]strfmt.UUID")

	initOutputFlags(deviceEnableCmd)
	initOutputFlags(deviceDisableCmd)

	initLoopControlFlags(deviceEnableCmd)
	initLoopControlFlags(deviceDisableCmd)

	initTenantFlags(deviceEnableCmd)
	initTenantFlags(deviceDisableCmd)
}
