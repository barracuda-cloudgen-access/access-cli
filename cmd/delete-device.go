// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2023 Barracuda Networks, Inc.

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

	apidevices "github.com/barracuda-cloudgen-access/access-cli/client/devices"
)

// deviceDeleteCmd represents the delete command
var deviceDeleteCmd = &cobra.Command{
	Use:     "delete [device ID]...",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete devices",
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
		deviceIDs, err := multiOpParseUUIDArgs(cmd, args, "id")
		if err != nil {
			return err
		}

		delete := func(id strfmt.UUID) error {
			gparams := apidevices.NewGetDeviceParams()
			setTenant(cmd, gparams)
			gparams.SetID(id)
			resp, err := global.Client.Devices.GetDevice(gparams, global.AuthWriter)
			if err != nil {
				return processErrorResponse(err)
			}

			params := apidevices.NewDeleteDeviceParams()
			setTenant(cmd, params)
			params.SetUserID(resp.Payload.User.ID)
			params.SetDeviceID(resp.Payload.ID)

			_, err = global.Client.Devices.DeleteDevice(params, global.AuthWriter)
			if err != nil {
				return processErrorResponse(err)
			}
			return nil
		}

		tw, j := multiOpBuildTableWriter()

		for _, arg := range deviceIDs {
			err = delete(arg)
			if err != nil {
				multiOpTableWriterAppend(tw, &j, arg, processErrorResponse(err))
				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, j, tw, len(deviceIDs), err)
			}
			multiOpTableWriterAppend(tw, &j, arg, "success")
		}
		return printListOutputAndError(cmd, j, tw, len(deviceIDs), err)
	},
}

func init() {
	devicesCmd.AddCommand(deviceDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deviceDeleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceDeleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initMultiOpArgFlags(deviceDeleteCmd, "device", "delete", "id", "[]strfmt.UUID")
	initOutputFlags(deviceDeleteCmd)
	initLoopControlFlags(deviceDeleteCmd)
	initTenantFlags(deviceDeleteCmd)
}
