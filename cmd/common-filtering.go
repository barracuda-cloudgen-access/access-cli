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

	"github.com/spf13/cobra"
)

type filterData struct {
	types []filterType
}

type filterType struct {
	name    string
	vartype string
}

func initFilterFlags(cmd *cobra.Command, filterTypes ...filterType) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitFilter] = "yes"
	for _, filterType := range filterTypes {
		desc := fmt.Sprintf("filter output %s", filterType.name)
		name := fmt.Sprintf("filter-%s", filterType.name)

		switch filterType.vartype {
		// add more types, as needed. don't forget to add in setFilter too
		case "int":
			cmd.Flags().Int(name, 0, desc)
		case "string":
			cmd.Flags().String(name, "", desc)
		case "[]string":
			cmd.Flags().StringSlice(name, []string{}, desc)
		default:
			panic("Unknown filter variable type " + filterType.vartype)
		}
	}
	if global.FilterData == nil {
		global.FilterData = make(map[*cobra.Command]*filterData)
	}
	global.FilterData[cmd] = &filterData{
		types: filterTypes,
	}
}

func preRunFlagCheckFilter(cmd *cobra.Command, args []string) error {
	// TODO perform some sort of parameter validation?
	return nil
}

func setFilter(cmd *cobra.Command, filterApplyFuncs ...interface{}) {
	if _, ok := cmd.Annotations[flagInitFilter]; !ok {
		panic("setFilter called for command where filtering flag was not initialized. This is a bug!")
	}
	data := global.FilterData[cmd]
	if len(filterApplyFuncs) != len(data.types) {
		panic("setFilter called with insufficient parameters. This is a bug!")
	}
	for i, filterType := range data.types {
		flagName := fmt.Sprintf("filter-%s", filterType.name)
		switch f := filterApplyFuncs[i].(type) {
		// add more types, as needed. don't forget to add in initFilterFlags too
		case func(int):
			d, err := cmd.Flags().GetInt(flagName)
			if err == nil {
				f(d)
			}
		case func(*int):
			d, err := cmd.Flags().GetInt(flagName)
			if err == nil {
				f(&d)
			}
		case func(string):
			d, err := cmd.Flags().GetString(flagName)
			if err == nil {
				f(d)
			}
		case func(*string):
			d, err := cmd.Flags().GetString(flagName)
			if err == nil {
				f(&d)
			}
		case func([]string):
			d, err := cmd.Flags().GetStringSlice(flagName)
			if err == nil {
				f(d)
			}
		default:
			panic("setFilter called with inadequate function in parameters")
		}
	}
}
