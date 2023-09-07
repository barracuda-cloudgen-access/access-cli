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

import "github.com/spf13/cobra"

func initLoopControlFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitLoopControl] = "yes"
	cmd.Flags().Bool("continue-on-error", false, "whether to continue if the operation fails for one of the arguments")
}

func preRunFlagCheckLoopControl(cmd *cobra.Command, args []string) error {
	return nil
}

func loopControlContinueOnError(cmd *cobra.Command) bool {
	if _, ok := cmd.Annotations[flagInitLoopControl]; !ok {
		panic("loopControlContinueOnError called for command where looping flag was not initialized. This is a bug!")
	}
	c, err := cmd.Flags().GetBool("continue-on-error")
	if err == nil {
		return c
	}
	return false
}
