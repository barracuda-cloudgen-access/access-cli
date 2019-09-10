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
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// endpointCmd represents the endpoint command
var endpointCmd = &cobra.Command{
	Use:   "endpoint",
	Short: "Get currently configured console endpoint",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && args[0] != "get" {
			return fmt.Errorf("use `fyde-cli endpoint set` to set the console endpoint")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		endpoint := authViper.GetString("endpoint")
		if endpoint == "" {
			fmt.Println("Endpoint not currently set")
		}
		fmt.Println("Currently configured endpoint:")
		fmt.Println(endpoint)
	},
}

func init() {
	rootCmd.AddCommand(endpointCmd)
}
