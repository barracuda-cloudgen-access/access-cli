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

import "github.com/spf13/cobra"

type sortable interface {
	SetSort(sort *string)
}

func initSortFlags(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations["sort_flags_init"] = "yes"
	cmd.Flags().String("sort", "id_asc", "sort output. Possible options include: id_{asc|desc}, name_{asc|desc}, created_{asc|desc}, updated_{asc|desc}")
}

func preRunFlagCheckSort(cmd *cobra.Command, args []string) error {
	// TODO perform some sort of parameter validation?
	return nil
}

func setSort(cmd *cobra.Command, s sortable) {
	if _, ok := cmd.Annotations["sort_flags_init"]; !ok {
		panic("setSort called for command where sorting flag was not initialized. This is a bug!")
	}
	sort, err := cmd.Flags().GetString("sort")
	if err == nil {
		s.SetSort(&sort)
	}
}
