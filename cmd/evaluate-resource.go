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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apidevices "github.com/barracuda-cloudgen-access/access-cli/client/devices"
)

func getBoolNoError(cmd *cobra.Command, argName string) bool {
	var val bool
	val, _ = cmd.Flags().GetBool(argName)
	return val
}

// evaluateResourceCmd represents the get command
var evaluateResourceCmd = &cobra.Command{
	Use:   "evaluate [device ID] [resource domain/ip]",
	Short: "Evaluate device against a resource",
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
			return fmt.Errorf("missing device ID argument")
		}

		if len(args) < 2 && !cmd.Flags().Changed("resource") {
			return fmt.Errorf("missing resource argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var deviceID string
		var resource string
		var err error
		if cmd.Flags().Changed("id") {
			deviceID, err = cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
		} else {
			deviceID = args[0]
		}

		if cmd.Flags().Changed("resource") {
			resource, err = cmd.Flags().GetString("resource")
			if err != nil {
				return err
			}
		} else {
			resource = args[1]
		}

		params := apidevices.NewEvaluateResourceParams()
		setTenant(cmd, params)

		attributes := apidevices.EvaluateResourceParamsBodyAttributes{
			Antivirus:  &[]bool{getBoolNoError(cmd, "antivirus")}[0],
			Fde:        &[]bool{getBoolNoError(cmd, "fde")}[0],
			Firewall:   &[]bool{getBoolNoError(cmd, "firewall")}[0],
			Jailbroken: &[]bool{!(getBoolNoError(cmd, "not-jailbroken"))}[0],
			ScreenLock: &[]bool{getBoolNoError(cmd, "screenlock")}[0],
		}
		body := apidevices.EvaluateResourceBody{
			Attributes:     &attributes,
			ResourceDomain: &resource,
		}

		params.SetID(strfmt.UUID(deviceID))
		params.SetEvaluateResource(body)

		resp, err := global.Client.Devices.EvaluateResource(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"Result",
			"Remediations",
			"Errors",
		})
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 36},
			{Number: 2, WidthMax: 6},
			{Number: 3, WidthMax: 60},
			{Number: 4, WidthMax: 60},
		})

		remediations, _ := renderJSON(resp.Payload.Remediations)
		errors, _ := renderJSON(resp.Payload.Errors)

		tw.AppendRow(table.Row{
			strfmt.UUID(deviceID),
			resp.Payload.Result,
			remediations,
			errors,
		})

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func init() {
	devicesCmd.AddCommand(evaluateResourceCmd)

	evaluateResourceCmd.Flags().BoolP("antivirus", "a", false, "Device has antivirus")
	evaluateResourceCmd.Flags().BoolP("fde", "d", false, "Device has full-disk encryption")
	evaluateResourceCmd.Flags().BoolP("firewall", "f", false, "Device has firewall enabled")
	evaluateResourceCmd.Flags().BoolP("not-jailbroken", "j", false, "Device is not jailbreak")
	evaluateResourceCmd.Flags().BoolP("screenlock", "s", false, "Device has screen lock enabled")

	initOutputFlags(evaluateResourceCmd)
	evaluateResourceCmd.Flags().String("id", "", "id of device to evaluate")
	evaluateResourceCmd.Flags().String("resource", "", "Domain or IP of resource to evaluate")

	initTenantFlags(evaluateResourceCmd)
}
