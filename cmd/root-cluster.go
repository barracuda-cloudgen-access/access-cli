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
)

// endpointCmd represents the endpoint command
var clusterCmd = &cobra.Command{
	Use:     "cluster",
	Aliases: []string{"endpoint"},
	Short:   "Get currently configured console cluster",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && args[0] != "get" {
			return fmt.Errorf("use `access-cli cluster set` to set the console cluster")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		cluster := authViper.GetString(ckeyAuthEndpoint)
		if cluster == "" {
			cmd.Println("Cluster not currently set")
		}
		cmd.Println("Currently configured cluster:")
		cmd.Println(cluster)
	},
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
