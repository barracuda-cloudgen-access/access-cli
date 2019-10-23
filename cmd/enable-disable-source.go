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

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apisources "github.com/fyde/fyde-cli/client/asset_sources"
)

// sourceEnableCmd represents the enable command
var sourceEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "enable source",
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
			return fmt.Errorf("missing source ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		enable := cmd.Use == "enable"
		for _, arg := range args {
			params := apisources.NewEditAssetSourceParams()
			params.SetID(strfmt.UUID(arg))
			params.SetAssetSource(apisources.EditAssetSourceBody{
				AssetSource: &apisources.EditAssetSourceParamsBodyAssetSource{
					Enabled: &enable,
				},
			})

			resp, err := global.Client.AssetSources.EditAssetSource(params, global.AuthWriter)
			if err != nil {
				if loopControlContinueOnError(cmd) {
					cmd.PrintErrln(processErrorResponse(err))
					continue
				}
				return processErrorResponse(err)
			}

			if resp.Payload.Enabled {
				cmd.Println("Source", resp.Payload.ID, "enabled")
			} else {
				cmd.Println("Source", resp.Payload.ID, "disabled")
			}
		}
		return nil
	},
}

// sourceDisableCmd represents the disable command
var sourceDisableCmd *cobra.Command

func init() {
	disableCmd := *sourceEnableCmd
	disableCmd.Use = "disable"
	disableCmd.Short = "disable source"
	sourceDisableCmd = &disableCmd
	sourcesCmd.AddCommand(sourceEnableCmd)
	sourcesCmd.AddCommand(sourceDisableCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sourceEnableCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sourceEnableCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initLoopControlFlags(sourceEnableCmd)
	initLoopControlFlags(sourceDisableCmd)
}
