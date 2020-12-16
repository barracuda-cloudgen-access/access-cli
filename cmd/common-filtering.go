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

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

type filterData struct {
	types []filterType
}

type filterType struct {
	name    string
	vartype string
}

func initFilterFlags(cmd *cobra.Command, filterTypes ...filterType) {
	cmd.Flags().SortFlags = false
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
		case "[]int":
			// see https://github.com/spf13/pflag/issues/222
			// we will accept a string slice instead, and convert to a int slice later
			cmd.Flags().StringSlice(name, []string{}, desc)
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
		d, err := getFlagValue(cmd, filterType.vartype, flagName)
		if err != nil {
			continue
		}

		callApplyFunc(filterApplyFuncs[i], d, filterType.vartype)
	}
}

func callApplyFunc(f, value interface{}, varType string) {
	switch f := f.(type) {
	// add more types, as needed. don't forget to add in initFilterFlags too
	case func(bool):
		f(value.(bool))
	case func(*bool):
		dd := value.(bool)
		f(&dd)
	case func(int):
		f(value.(int))
	case func(*int):
		dd := value.(int)
		f(&dd)
	case func(int64):
		f(int64(value.(int)))
	case func(*int64):
		dd := int64(value.(int))
		f(&dd)
	case func(string):
		f(value.(string))
	case func(*string):
		dd := value.(string)
		f(&dd)
	case func([]int):
		f(value.([]int))
	case func([]int64):
		dconv := funk.Map(value, func(x int) int64 {
			return int64(x)
		}).([]int64)
		f(dconv)
	case func([]string):
		f(value.([]string))
	case func([]strfmt.UUID):
		dconv := funk.Map(value, func(x string) strfmt.UUID {
			return strfmt.UUID(x)
		}).([]strfmt.UUID)
		f(dconv)
	default:
		panic(fmt.Errorf("callApplyFunc called with inadequate function in parameters (function is %T for vartype %s)", f, varType))
	}
}
