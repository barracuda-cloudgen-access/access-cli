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

	"github.com/spf13/cobra"

	apidevices "github.com/barracuda-cloudgen-access/access-cli/client/devices"
)

// deviceRevokeCmd represents the revoke command
var deviceRevokeCmd = &cobra.Command{
	Use:   "revoke [device ID]...",
	Short: "Revoke device authentication",
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
		uuidArgs, err := multiOpParseUUIDArgs(cmd, args, "id")
		if err != nil {
			return err
		}

		tw, j := multiOpBuildTableWriter()

		for _, id := range uuidArgs {
			params := apidevices.NewRevokeDeviceParams()
			setTenant(cmd, params)
			params.SetID(id)

			_, err = global.Client.Devices.RevokeDevice(params, global.AuthWriter)
			if err != nil {
				multiOpTableWriterAppend(tw, &j, id, processErrorResponse(err))
				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, j, tw, len(uuidArgs), err)
			}
			multiOpTableWriterAppend(tw, &j, id, "success")
		}
		return printListOutputAndError(cmd, j, tw, len(uuidArgs), err)
	},
}

func init() {
	devicesCmd.AddCommand(deviceRevokeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deviceRevokeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deviceRevokeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initMultiOpArgFlags(deviceRevokeCmd, "device", "revoke", "id", "[]strfmt.UUID")
	initOutputFlags(deviceRevokeCmd)
	initLoopControlFlags(deviceRevokeCmd)
	initTenantFlags(deviceRevokeCmd)
}
