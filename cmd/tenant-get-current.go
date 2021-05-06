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

// tenantCurrentGetCmd represents the enable command
var tenantCurrentGetCmd = &cobra.Command{
	Use:   "get-current",
	Short: "Get currently configured console tenant",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && args[0] != "get-current" {
			return fmt.Errorf("use `access-cli tenant set-current` to set the console tenant")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tenant := authViper.GetString(ckeyAuthCurrentTenant)
		if tenant == "" {
			cmd.Println("Current tenant not currently set")
		}
		cmd.Println("Currently configured tenant:")
		cmd.Println(tenant)
	},
}

func init() {
	tenantsCmd.AddCommand(tenantCurrentGetCmd)
}
