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

	"github.com/spf13/cobra"
)

// tenantCurrentSetCmd represents the enable command
var tenantCurrentSetCmd = &cobra.Command{
	Use:   "set-current [tenant ID]...",
	Short: "Set current tenant",
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
			return fmt.Errorf("missing tenant ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		authViper.Set(ckeyAuthCurrentTenant, args[0])

		if global.WriteFiles {
			err := authViper.WriteConfig()
			if err != nil {
				return err
			}
		}
		cmd.Printf("Tenant set to %s.\n\n", args[0])
		return nil
	},
}

func init() {
	tenantsCmd.AddCommand(tenantCurrentSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tenantCurrentSetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tenantCurrentSetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	tenantCurrentSetCmd.Flags().String("id", "", "id of tenant to set")

	initOutputFlags(tenantCurrentSetCmd)

	initLoopControlFlags(tenantCurrentSetCmd)
}
