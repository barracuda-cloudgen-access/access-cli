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

type searchable interface {
	SetQ(query *string)
}

func initSearchFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitSearch] = "yes"
	cmd.Flags().StringP("search", "q", "", "full-text search query for result filtering")
}

func preRunFlagCheckSearch(cmd *cobra.Command, args []string) error {
	// TODO perform some sort of parameter validation?
	return nil
}

func setSearchQuery(cmd *cobra.Command, s searchable) {
	if _, ok := cmd.Annotations[flagInitSearch]; !ok {
		panic("setSearchQuery called for command where search flag was not initialized. This is a bug!")
	}
	query, err := cmd.Flags().GetString("search")
	if err == nil && query != "" {
		s.SetQ(&query)
	}
}
