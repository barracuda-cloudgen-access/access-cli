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
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/gbl08ma/mapstructure"
	"github.com/go-openapi/strfmt"
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
	MainField       bool
	SchemaName      string // only used if MainField is true, so that error handling functions can get an identifier for the failing record
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

func initInputFlags(cmd *cobra.Command, typeName string, fields ...inputField) {
	cmd.Flags().SortFlags = false

	// all output goes to stderr
	cmd.Flags().SetOutput(nil)

	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitInput] = "yes"
	typeName = pluralize(typeName)

	cmd.Flags().StringP("from-file", "f", "", "file from where to import "+typeName)
	cmd.Flags().StringP("file-format", "i", "json", "format for the file from where to import "+typeName+" (csv or json)")
	cmd.Flags().Bool("errors-only", false, "only include failed operations in output")

	for _, field := range fields {
		switch field.VarType {
		case "bool":
			cmd.Flags().BoolP(field.FlagName, field.FlagShorthand, field.DefaultValue.(bool), field.FlagDescription)
		case "int":
			cmd.Flags().IntP(field.FlagName, field.FlagShorthand, field.DefaultValue.(int), field.FlagDescription)
		case "string":
			cmd.Flags().StringP(field.FlagName, field.FlagShorthand, field.DefaultValue.(string), field.FlagDescription)
		case "string+file":
			cmd.Flags().StringP(field.FlagName, field.FlagShorthand, field.DefaultValue.(string), field.FlagDescription)
			// special case to read input from file, appends -file to argument
			cmd.Flags().StringP(wrapFlagForFile(field.FlagName), field.FlagShorthand, field.DefaultValue.(string), field.FlagDescription+" (from file path)")
		case "[]int":
			// see https://github.com/spf13/pflag/issues/222
			// cmd.Flags().IntSliceP(field.FlagName, field.FlagShorthand, field.DefaultValue.([]int), field.FlagDescription)
			// we will accept a string slice instead, and convert to a int slice later
			// convert default value from int slice to string slice:
			dconv := funk.Map(field.DefaultValue.([]int), func(x int) string {
				return strconv.Itoa(x)
			}).([]string)
			cmd.Flags().StringSliceP(field.FlagName, field.FlagShorthand, dconv, field.FlagDescription)
		case "[]string":
			cmd.Flags().StringSliceP(field.FlagName, field.FlagShorthand, field.DefaultValue.([]string), field.FlagDescription)
		case "[]string.skipcomma":
			cmd.Flags().StringArrayP(field.FlagName, field.FlagShorthand, field.DefaultValue.([]string), field.FlagDescription)
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

func wrapFlagForFile(flagName string) string {
	return flagName + "-file"
}

func getFlagValueChanged(cmd *cobra.Command, varType, flagName string) bool {
	changed := cmd.Flags().Changed(flagName)
	if !changed && varType == "string+file" {
		changed = cmd.Flags().Changed(wrapFlagForFile(flagName))
	}
	return changed
}

func getFlagValue(cmd *cobra.Command, varType, flagName string) (interface{}, error) {
	var err error
	var value interface{}
	switch varType {
	case "bool":
		value, err = cmd.Flags().GetBool(flagName)
	case "int":
		value, err = cmd.Flags().GetInt(flagName)
	case "string":
		value, err = cmd.Flags().GetString(flagName)
	case "string+file":
		// If string value was passed, return it, otherwise, attempt to read from file
		value, err = cmd.Flags().GetString(flagName)
		if err != nil {
			return nil, err
		}
		if value != "" {
			return value, nil
		}
		filename, err := cmd.Flags().GetString(wrapFlagForFile(flagName))
		if err != nil {
			return nil, err
		}
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		value = string(contents)
	case "[]int":
		// see https://github.com/spf13/pflag/issues/222
		// we accepted a string slice instead, and will now convert to a int slice
		strSlice, err := cmd.Flags().GetStringSlice(flagName)
		if err != nil {
			return nil, err
		}
		if len(strSlice) == 1 && strings.Trim(strSlice[0], "[] ") == "" {
			return []int{}, nil
		}
		out := make([]int, len(strSlice))
		for i, str := range strSlice {
			// support [1,2,3] syntax
			str = strings.Trim(str, "[] ")

			var err error
			out[i], err = strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
		}
		return out, nil
	case "[]string":
		value, err = cmd.Flags().GetStringSlice(flagName)
	case "[]string.skipcomma":
		value, err = cmd.Flags().GetStringArray(flagName)
	default:
		panic("Unknown variable type " + varType)
	}
	return value, err
}

func forAllInput(cmd *cobra.Command,
	args []string,
	writeDefaultValues bool,
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

		flagChanged := getFlagValueChanged(cmd, field.VarType, field.FlagName)
		if !flagChanged && field.Mandatory {
			// user did not supply the field value in a flag
			fieldInArgs := false
			if field.MainField {
				// attempt to get it from the arg
				d, fieldInArgs = readFieldFromArgs(args, field)
			}
			if !fieldInArgs {
				// must ask interactively
				d = interactivelyReadField(cmd, field)
			}
		} else if !flagChanged && !writeDefaultValues {
			// in placeInputValues, we only call the respective function for this field
			// if the value for the entry is not nil. so, we just leave it nil, to indicate this value is not provided
			continue
		} else {
			d, err = getFlagValue(cmd, field.VarType, field.FlagName)
			if err != nil {
				if field.Mandatory {
					return err
				} else if field.DefaultValue != nil && writeDefaultValues {
					d = field.DefaultValue
				}
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

func readFieldFromArgs(args []string, field inputField) (interface{}, bool) {
	if len(args) == 0 {
		return nil, false
	}

	switch field.VarType {
	case "int":
		i, err := strconv.Atoi(args[0])
		return i, err == nil
	case "string":
		return strings.Join(args, " "), true
	default:
		// we don't need to support other vartypes in this function
		return nil, false
	}
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
			// special behavior for access resource port mappings, to keep compatibility with previous versions
			if strings.ToLower(header[i]) == "ports" {
				header[i] = "PortMappings"
			}
			if strings.ToLower(header[i]) == "port_mappings" ||
				strings.ToLower(header[i]) == "portmappings" {
				m[header[i]] = []*models.AccessResourcePortMapping{
					colonMappingToPortMapping(strings.TrimRight(strings.TrimLeft(record[i], "["), "]")),
				}
			} else {
				m[header[i]] = record[i]
			}
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
		if len(setterFuncs) != len(data.fields) {
			panic("Setter funcs do not match number of input fields!")
		}
		for i, field := range data.fields {
			if entry.Values[i] != nil {
				callApplyFunc(setterFuncs[i], entry.Values[i], field.VarType)
			}
		}
		return nil
	case wholeJSONObject:
		return json.Unmarshal(entry.JSON, object)
	case wholeCSVObject:
		config := &mapstructure.DecoderConfig{
			DecodeHook:       csvMapstructureDecodeHook,
			WeaklyTypedInput: true,
			Squash:           true,
			Result:           object,
		}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return err
		}
		return decoder.Decode(entry.CSVdata)
	}
	return fmt.Errorf("unknown input entry type")
}

func csvMapstructureDecodeHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() != reflect.String || to.Kind() != reflect.Slice {
		return data, nil
	}
	s := data.(string)
	return commaSeparatedListToStringSlice(s), nil
}

func commaSeparatedListToStringSlice(s string) []string {
	s = strings.TrimRight(strings.TrimLeft(s, "["), "]")
	split := strings.Split(s, ",")
	for i := range split {
		split[i] = strings.TrimSpace(split[i])
	}
	return split
}

// the sole purpose of this function is to attempt to recover the ID of an object
// (of unknown type) in order to be able to tell the user that an error ocurred
// for that ID
func getIDinputValue(cmd *cobra.Command, entry *inputEntry) interface{} {
	if _, ok := cmd.Annotations[flagInitInput]; !ok {
		panic("getIdInputValue called for command where input flags were not initialized. This is a bug!")
	}
	data := global.InputData[cmd]

	// identify the index of the field that was marked as MainField
	indexOfMainfield := func(fields []inputField) int {
		for i, field := range fields {
			if field.MainField {
				return i
			}
		}
		return -1
	}

	idIdx := indexOfMainfield(data.fields)
	if idIdx < 0 {
		// no field marked as main, we can't help the user pinpoint the failure
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

type multiOpData struct {
	fieldName string
	fieldType string
}

func initMultiOpArgFlags(cmd *cobra.Command, typeName, operationVerb, fieldName, fieldType string) {
	cmd.Flags().SortFlags = false
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitMultiOpArg] = "yes"
	typeName = pluralize(typeName)

	description := fmt.Sprintf("%s of %s to %s", pluralize(fieldName), typeName, operationVerb)

	switch fieldType {
	case "[]int64":
		// see https://github.com/spf13/pflag/issues/222
		// we will accept a string slice instead, and convert to a int slice later
		cmd.Flags().StringSlice(fieldName, []string{}, description)
	case "[]string":
		fallthrough
	case "[]strfmt.UUID":
		cmd.Flags().StringSlice(fieldName, []string{}, description)
	default:
		panic("Unsupported field type " + fieldType)
	}

	if global.MultiOpData == nil {
		global.MultiOpData = make(map[*cobra.Command]*multiOpData)
	}
	global.MultiOpData[cmd] = &multiOpData{
		fieldName: fieldName,
		fieldType: fieldType,
	}
}

func multiOpCheckArgsPresent(cmd *cobra.Command, args []string) bool {
	if _, ok := cmd.Annotations[flagInitMultiOpArg]; !ok {
		panic("multiOpCheckArgsPresent called for command where multi-op arg flags were not initialized. This is a bug!")
	}

	if len(args) > 0 {
		return true
	}

	if cmd.Flags().Changed(global.MultiOpData[cmd].fieldName) {
		return true
	}

	// check if there's piped input
	stdinInfo, err := os.Stdin.Stat()
	return err == nil && (stdinInfo.Mode()&os.ModeCharDevice == 0)
}

func multiOpParseInt64Args(cmd *cobra.Command, args []string, idFieldName string) ([]int64, error) {
	if _, ok := cmd.Annotations[flagInitMultiOpArg]; !ok {
		panic("multiOpParseInt64Args called for command where multi-op arg flags were not initialized. This is a bug!")
	}

	if cmd.Flags().Changed(global.MultiOpData[cmd].fieldName) {
		switch global.MultiOpData[cmd].fieldType {
		case "[]int64":
			// see https://github.com/spf13/pflag/issues/222
			// we accepted a string slice instead, and will now convert to a int slice
			strSlice, err := cmd.Flags().GetStringSlice(global.MultiOpData[cmd].fieldName)
			if err != nil {
				return nil, err
			}
			if len(strSlice) == 1 && strings.Trim(strSlice[0], "[] ") == "" {
				return []int64{}, nil
			}
			out := make([]int64, len(strSlice))
			for i, str := range strSlice {
				// support [1,2,3] syntax
				str = strings.Trim(str, "[] ")

				var err error
				out[i], err = strconv.ParseInt(str, 10, 64)
				if err != nil {
					return nil, err
				}
			}
			return out, nil
		default:
			panic("multiOpParseInt64Args called where it shouldn't have been. This is a bug!")
		}
	}

	if len(args) > 0 {
		intArgs := make([]int64, len(args))
		for i, arg := range args {
			var err error
			intArgs[i], err = strconv.ParseInt(arg, 10, 64)
			if err != nil {
				return intArgs, err
			}
		}
		return intArgs, nil
	}

	stdinInfo, err := os.Stdin.Stat()
	hasPiped := err == nil && (stdinInfo.Mode()&os.ModeCharDevice == 0)
	if !hasPiped {
		return []int64{}, fmt.Errorf("missing arguments")
	}

	d := json.NewDecoder(os.Stdin)
	d.UseNumber() // important, otherwise all numbers are decoded as float64

	var jsonArray []map[string]interface{}
	if err := d.Decode(&jsonArray); err != nil {
		return []int64{}, fmt.Errorf("decoding JSON from pipe: %w", err)
	}

	intArgs := make([]int64, len(jsonArray))
	for i, itemMap := range jsonArray {
		numIface, present := itemMap[idFieldName]
		if !present {
			return intArgs, fmt.Errorf("key %s not present in piped JSON", idFieldName)
		}

		num, ok := numIface.(json.Number)
		if !ok {
			return intArgs, fmt.Errorf("value of %s is not a number in piped JSON", idFieldName)
		}

		intArgs[i], err = num.Int64()
		if err != nil {
			return intArgs, fmt.Errorf("error converting value of %s to integer: %w", idFieldName, err)
		}
	}

	return intArgs, nil
}

func multiOpParseUUIDArgs(cmd *cobra.Command, args []string, idFieldName string) ([]strfmt.UUID, error) {
	if _, ok := cmd.Annotations[flagInitMultiOpArg]; !ok {
		panic("multiOpParseUUIDArgs called for command where multi-op arg flags were not initialized. This is a bug!")
	}

	if cmd.Flags().Changed(global.MultiOpData[cmd].fieldName) {
		switch global.MultiOpData[cmd].fieldType {
		case "[]strfmt.UUID":
			strSlice, err := cmd.Flags().GetStringSlice(global.MultiOpData[cmd].fieldName)
			if err != nil {
				return nil, err
			}
			if len(strSlice) == 1 && strings.Trim(strSlice[0], "[] ") == "" {
				return []strfmt.UUID{}, nil
			}
			out := make([]strfmt.UUID, len(strSlice))
			for i, str := range strSlice {
				// support [1,2,3] syntax
				str = strings.Trim(str, "[] ")

				out[i] = strfmt.UUID(str)
			}
			return out, nil
		default:
			panic("multiOpParseUUIDArgs called where it shouldn't have been. This is a bug!")
		}
	}

	if len(args) > 0 {
		uuidArgs := make([]strfmt.UUID, len(args))
		for i, arg := range args {
			uuidArgs[i] = strfmt.UUID(arg)
		}
		return uuidArgs, nil
	}

	stdinInfo, err := os.Stdin.Stat()
	hasPiped := err == nil && (stdinInfo.Mode()&os.ModeCharDevice == 0)
	if !hasPiped {
		return []strfmt.UUID{}, fmt.Errorf("missing arguments")
	}

	d := json.NewDecoder(os.Stdin)

	var jsonArray []map[string]interface{}
	if err := d.Decode(&jsonArray); err != nil {
		return []strfmt.UUID{}, fmt.Errorf("decoding JSON from pipe: %w", err)
	}

	uuidArgs := make([]strfmt.UUID, len(jsonArray))
	for i, itemMap := range jsonArray {
		numIface, present := itemMap[idFieldName]
		if !present {
			return uuidArgs, fmt.Errorf("key %s not present in piped JSON", idFieldName)
		}

		s, ok := numIface.(string)
		if !ok {
			return uuidArgs, fmt.Errorf("value of %s is not a string in piped JSON", idFieldName)
		}

		uuidArgs[i] = strfmt.UUID(s)
	}

	return uuidArgs, nil
}
