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
	"github.com/spf13/cobra"
)

// settingsCmd represents the settings command
var settingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Operations on settings",
}

var settingsAgentConfigCmd = &cobra.Command{
	Use:   "agent",
	Short: "Operations on agent configuration",
}

var settingsEnrollmentCmd = &cobra.Command{
	Use:   "enrollment",
	Short: "Operations on enrollment",
}

var settingsAnalyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Operations on analytics",
}

var settingsIdentityProviderConfigCmd = &cobra.Command{
	Use:   "idp",
	Short: "Configure Identity Provider settings",
}

func init() {
	rootCmd.AddCommand(settingsCmd)

	settingsCmd.AddCommand(settingsAgentConfigCmd)
	settingsCmd.AddCommand(settingsEnrollmentCmd)
	settingsCmd.AddCommand(settingsAnalyticsCmd)
	settingsCmd.AddCommand(settingsIdentityProviderConfigCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// settingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// settingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
