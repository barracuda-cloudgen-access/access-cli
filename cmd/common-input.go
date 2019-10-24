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
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"

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
	IsIDOnError     bool
	SchemaName      string // only used if IsIDOnError is true, so that error handling functions can get an identifier for the failing record
	DefaultValue    interface{}
}

type inputEntry struct {
	Type    wholeObjectType
	CSVdata interface{}
	JSON    json.RawMessage
	Values  []interface{}
}

type wholeObjectType int

const (
	wholeCSVObject wholeObjectType = iota
	wholeJSONObject
	individualValues
)

func initInputFlags(cmd *cobra.Command, fields ...inputField) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitInput] = "yes"

	cmd.Flags().StringP("from-file", "f", "", "file from where to import users")
	cmd.Flags().StringP("file-format", "i", "json", "format for the file from where to import users (csv or json)")
	cmd.Flags().Bool("errors-only", false, "only include failed operations in output")

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

func forAllInput(cmd *cobra.Command,
	do func(entry *inputEntry) (interface{}, error),
	printSuccess func(interface{}),
	doOnError func(error, interface{})) error {
	if _, ok := cmd.Annotations[flagInitInput]; !ok {
		panic("forAllInput called for command where input flags were not initialized. This is a bug!")
	}
	data := global.InputData[cmd]

	if errorsOnly, err := cmd.Flags().GetBool("errors-only"); err == nil && errorsOnly {
		printSuccess = nil
	}

	fromFile, err := cmd.Flags().GetString("from-file")
	if err != nil {
		return err
	}
	if fromFile != "" {
		return forAllInputFromFile(cmd, do, printSuccess, doOnError)
	}

	entry := &inputEntry{
		Type:   individualValues,
		Values: make([]interface{}, len(data.fields)),
	}
	for i, field := range data.fields {
		var d interface{}

		if !cmd.Flags().Changed(field.FlagName) && field.Mandatory {
			// user did not supply the field value in a flag, must ask interactively
			d = interactivelyReadField(cmd, field)
		} else {
			d, err = getFlagValue(cmd, field.VarType, field.FlagName)
			if err != nil {
				continue
			}
		}
		entry.Values[i] = d
	}
	r, err := do(entry)
	if doOnError != nil && err != nil {
		doOnError(err, getIDinputValue(cmd, entry))
	} else if printSuccess != nil {
		printSuccess(r)
	}
	if !loopControlContinueOnError(cmd) {
		return err
	}
	return nil
}

func interactivelyReadField(cmd *cobra.Command, field inputField) interface{} {
	in := bufio.NewReader(os.Stdin)
	cmd.Printf("%s: ", field.Name)
	for {
		line, err := in.ReadString('\n')
		if err != nil {
			cmd.Println("error:", err)
			cmd.Printf("%s must be provided: ", field.Name)
			continue
		}
		line = strings.TrimRight(line, "\t\n\v\f\r")
		if len(line) == 0 {
			cmd.Printf("%s must be provided: ", field.Name)
			continue
		}
		if field.VarType == "string" || field.VarType == "[]string" {
			// we want fmt.Scanln's "magic" automatic types behavior, but we don't want the part where it stops at spaces
			line = strings.ReplaceAll(line, " ", "\uF8FF")
		}
		var d interface{}
		switch field.VarType {
		case "[]string", "[]int":
			// treat slices as string arrays, and convert back later
			s := ""
			spointer := &s
			d = spointer
		default:
			d = reflect.New(reflect.TypeOf(field.DefaultValue)).Elem().Addr().Interface()
		}
		i, err := fmt.Sscanln(line, d)
		if err != nil {
			cmd.Println("error:", err)
			cmd.Printf("%s must be provided: ", field.Name)
			continue
		}
		// the rest of our code doesn't expect a pointer
		d = reflect.ValueOf(d).Elem().Interface()
		if field.VarType == "string" || field.VarType == "[]string" {
			// undo "magic" transformation above
			d = strings.ReplaceAll(d.(string), "\uF8FF", " ")
		}
		// convert strings to slices, if applicable
		switch field.VarType {
		case "[]string":
			s := d.(string)
			d = strings.Split(s, ",")
		case "[]int":
			s := d.(string)
			slice := strings.Split(s, ",")
			intSlice := make([]int, len(slice))
			for i, e := range slice {
				intSlice[i], err = strconv.Atoi(e)
				if err != nil {
					break
				}
			}
			if err != nil {
				cmd.Println("error:", err)
				cmd.Printf("%s must be provided: ", field.Name)
				continue
			}
			d = intSlice
		}
		if i != 0 && field.Validator != nil {
			if !field.Validator(d) {
				cmd.Println("invalid value")
				continue
			}
		}
		return d
	}
}

func forAllInputFromFile(cmd *cobra.Command,
	do func(entry *inputEntry) (interface{}, error),
	printSuccess func(interface{}),
	doOnError func(error, interface{})) error {
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
		return forAllInputFromJSON(cmd, do, printSuccess, doOnError, reader)
	case "csv":
		return forAllInputFromCSV(cmd, do, printSuccess, doOnError, reader)
	}
	return nil
}

func forAllInputFromJSON(cmd *cobra.Command,
	do func(entry *inputEntry) (interface{}, error),
	printSuccess func(interface{}),
	doOnError func(error, interface{}),
	reader io.Reader) error {
	records := make([]json.RawMessage, 0)

	err := json.NewDecoder(reader).Decode(&records)
	if err != nil {
		return err
	}

	for _, record := range records {
		entry := &inputEntry{
			Type: wholeJSONObject,
			JSON: record,
		}
		r, err := do(entry)
		if err != nil {
			if !loopControlContinueOnError(cmd) {
				return err
			}
			if doOnError != nil {
				doOnError(err, getIDinputValue(cmd, entry))
			}
		} else if printSuccess != nil {
			printSuccess(r)
		}
	}
	return nil
}

func forAllInputFromCSV(cmd *cobra.Command,
	do func(entry *inputEntry) (interface{}, error),
	printSuccess func(interface{}),
	doOnError func(error, interface{}),
	reader io.Reader) error {
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

		entry := &inputEntry{
			Type:    wholeCSVObject,
			CSVdata: m,
		}
		res, err := do(entry)
		if err != nil {
			if doOnError != nil {
				doOnError(err, getIDinputValue(cmd, entry))
			}
			if !loopControlContinueOnError(cmd) {
				return err
			}
		} else if printSuccess != nil {
			printSuccess(res)
		}
	}
	return nil
}

func placeInputValues(cmd *cobra.Command,
	entry *inputEntry,
	object interface{},
	setterFuncs ...interface{}) error {
	if _, ok := cmd.Annotations[flagInitInput]; !ok {
		panic("forAllInput called for command where input flags were not initialized. This is a bug!")
	}
	data := global.InputData[cmd]

	switch entry.Type {
	case individualValues:
		for i, field := range data.fields {
			callApplyFunc(setterFuncs[i], entry.Values[i], field.VarType)
		}
		return nil
	case wholeJSONObject:
		return json.Unmarshal(entry.JSON, object)
	case wholeCSVObject:
		return mapstructure.WeakDecode(entry.CSVdata, object)
	}
	return fmt.Errorf("unknown input entry type")
}

// the sole purpose of this function is to attempt to recover the ID of an object
// (of unknown type) in order to be able to tell the user that an error ocurred
// for that ID
func getIDinputValue(cmd *cobra.Command, entry *inputEntry) interface{} {
	if _, ok := cmd.Annotations[flagInitInput]; !ok {
		panic("getIdInputValue called for command where input flags were not initialized. This is a bug!")
	}
	data := global.InputData[cmd]

	// identify the index of the field that was marked as IsIDOnError
	indexOfIDfield := func(fields []inputField) int {
		for i, field := range fields {
			if field.IsIDOnError {
				return i
			}
		}
		return -1
	}

	idIdx := indexOfIDfield(data.fields)
	if idIdx < 0 {
		// no field marked as ID on error, we can't help the user pinpoint the failure
		return nil
	}

	switch entry.Type {
	case individualValues:
		return entry.Values[idIdx]
	case wholeJSONObject:
		var object map[string]interface{}
		err := json.Unmarshal(entry.JSON, &object)
		if err != nil {
			return err
		}
		return object[data.fields[idIdx].SchemaName]
	case wholeCSVObject:
		m := entry.CSVdata.(map[string]interface{})
		return m[data.fields[idIdx].SchemaName]
	}
	return nil
}
