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
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

type inputData struct {
	fields []inputField
}

type inputField struct {
	Name            string
	FlagName        string
	FlagShorthand   string
	FlagDescription string
	Validator       func(interface{}) bool
	VarType         string
	Mandatory       bool
	DefaultValue    interface{}
}

func initInputFlags(cmd *cobra.Command, fields ...inputField) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitInput] = "yes"

	cmd.Flags().StringP("from-file", "f", "", "File from where to import users")
	cmd.Flags().StringP("file-format", "i", "json", "Format for the file from where to import users (csv or json)")

	for _, field := range fields {
		switch field.VarType {
		case "bool":
			cmd.Flags().BoolP(field.FlagName, field.FlagShorthand, field.DefaultValue.(bool), field.FlagDescription)
		case "int":
			cmd.Flags().IntP(field.FlagName, field.FlagShorthand, field.DefaultValue.(int), field.FlagDescription)
		case "string":
			cmd.Flags().StringP(field.FlagName, field.FlagShorthand, field.DefaultValue.(string), field.FlagDescription)
		case "[]int":
			cmd.Flags().IntSliceP(field.FlagName, field.FlagShorthand, field.DefaultValue.([]int), field.FlagDescription)
		case "[]string":
			cmd.Flags().StringSliceP(field.FlagName, field.FlagShorthand, field.DefaultValue.([]string), field.FlagDescription)
		default:
			panic("Unknown filter variable type " + field.VarType)
		}
	}

	if global.InputData == nil {
		global.InputData = make(map[*cobra.Command]*inputData)
	}
	global.InputData[cmd] = &inputData{
		fields: fields,
	}
}

func preRunFlagCheckInput(cmd *cobra.Command, args []string) error {
	data := global.InputData[cmd]

	input, err := cmd.Flags().GetString("file-format")
	if err != nil {
		return err
	}
	if !funk.Contains([]string{"json", "csv"}, input) {
		return fmt.Errorf("invalid input file format %s", input)
	}

	for _, field := range data.fields {
		if field.Validator == nil {
			continue
		}

		value, err := getFlagValue(cmd, field.VarType, field.FlagName)
		if err == nil {
			if !field.Validator(value) {
				return fmt.Errorf("invalid value for field %s", field.Name)
			}
		}
	}

	return nil
}

func getFlagValue(cmd *cobra.Command, varType, flagName string) (interface{}, error) {
	var value interface{}
	var err error
	switch varType {
	case "bool":
		value, err = cmd.Flags().GetBool(flagName)
	case "int":
		value, err = cmd.Flags().GetInt(flagName)
	case "string":
		value, err = cmd.Flags().GetString(flagName)
	case "[]int":
		value, err = cmd.Flags().GetIntSlice(flagName)
	case "[]string":
		value, err = cmd.Flags().GetStringSlice(flagName)
	default:
		panic("Unknown variable type " + varType)
	}
	return value, err
}

func forAllInput(cmd *cobra.Command, do func(values []interface{}) error) error {
	if _, ok := cmd.Annotations[flagInitInput]; !ok {
		panic("forAllInput called for command where input flags were not initialized. This is a bug!")
	}
	data := global.InputData[cmd]

	fromFile, err := cmd.Flags().GetString("from-file")
	if err != nil {
		return err
	}
	if fromFile != "" {
		return forAllInputFromFile(cmd, do)
	}

	values := make([]interface{}, len(data.fields))
	for i, field := range data.fields {
		var d interface{}

		if !cmd.Flags().Changed(field.FlagName) && field.Mandatory {
			// user did not supply the field value in a flag, must ask interactively
			cmd.Printf("%s: ", field.Name)
			for {
				d = reflect.New(reflect.TypeOf(field.DefaultValue)).Elem().Addr().Interface()
				i, err := fmt.Scanln(d)
				if err != nil {
					panic(err)
				}
				// the rest of our code doesn't expect a pointer
				d = reflect.ValueOf(d).Elem().Interface()
				if i != 0 && field.Validator != nil {
					if !field.Validator(d) {
						cmd.Println("invalid value")
						continue
					}
				}
				break
			}
		} else {
			d, err = getFlagValue(cmd, field.VarType, field.FlagName)
			if err != nil {
				continue
			}
		}
		values[i] = d
	}
	return do(values)
}

func forAllInputFromFile(cmd *cobra.Command, do func(values []interface{}) error) error {
	inputFormat, err := cmd.Flags().GetString("file-format")
	if err != nil {
		return err
	}

	inputFile, err := cmd.Flags().GetString("from-file")
	if err != nil {
		return err
	}

	reader, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	switch inputFormat {
	case "json":
		return forAllInputFromJSON(cmd, do, reader)
	case "csv":
		return forAllInputFromCSV(cmd, do, reader)
	}
	return nil
}

type wholeObjectFlagType struct{}

var wholeObjectFlag = wholeObjectFlagType{}

func forAllInputFromJSON(cmd *cobra.Command, do func(values []interface{}) error, reader io.Reader) error {
	records := make([]interface{}, 0)

	err := json.NewDecoder(reader).Decode(&records)
	if err != nil {
		return err
	}

	for _, record := range records {
		err = do([]interface{}{wholeObjectFlag, record})
		if err != nil {
			if loopControlContinueOnError(cmd) {
				cmd.PrintErrln(processErrorResponse(err))
				continue
			}
			return err
		}
	}
	return nil
}

func forAllInputFromCSV(cmd *cobra.Command, do func(values []interface{}) error, reader io.Reader) error {
	r := csv.NewReader(reader)

	header, err := r.Read()
	if err == io.EOF {
		return fmt.Errorf("CSV file is missing header")
	}
	if err != nil {
		return err
	}

	for lineNumber := 1; ; lineNumber++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != len(header) {
			return fmt.Errorf("record %d is malformed", lineNumber)
		}

		m := make(map[string]interface{})
		for i := range record {
			m[header[i]] = record[i]
		}

		err = do([]interface{}{wholeObjectFlag, m})
		if err != nil {
			if loopControlContinueOnError(cmd) {
				cmd.PrintErrln(processErrorResponse(err))
				continue
			}
			return err
		}
	}
	return nil
}

func placeInputValues(cmd *cobra.Command, values []interface{}, object interface{}, setterFuncs ...interface{}) error {
	if _, ok := cmd.Annotations[flagInitInput]; !ok {
		panic("forAllInput called for command where input flags were not initialized. This is a bug!")
	}
	data := global.InputData[cmd]

	if len(values) == 2 && values[0] == wholeObjectFlag {
		return mapstructure.WeakDecode(values[1], object)
	}

	for i, field := range data.fields {
		callApplyFunc(setterFuncs[i], values[i], field.VarType)
	}
	return nil
}
