// Package cmd implements fyde-cli commands
package cmd

/*
Copyright © 2019 Fyde, Inc. <hello@fyde.com>

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

	"github.com/spf13/cobra"

	apidevices "github.com/fyde/fyde-cli/client/devices"
)

// deviceRevokeCmd represents the revoke command
var deviceRevokeCmd = &cobra.Command{
	Use:   "revoke",
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

		if len(args) == 0 {
			return fmt.Errorf("missing device ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		deviceID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		params := apidevices.NewRevokeDeviceParams()
		params.SetID(deviceID)

		resp, err := global.Client.Devices.RevokeDevice(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		cmd.Println("Authentication revoked for device", resp.Payload.ID)
		return nil
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

	initOutputFlags(deviceRevokeCmd)
}