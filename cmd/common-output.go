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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"golang.org/x/crypto/ssh/terminal"
)

func initOutputFlags(cmd *cobra.Command) {
	cmd.Flags().SortFlags = false
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitOutput] = "yes"
	d := "json"
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		d = "table"
	}
	cmd.Flags().StringP("output", "o", d, "output format (table, json, json-pretty or csv) (default \"json\" if pipe)")
	cmd.Flags().SetNormalizeFunc(aliasNormalizeFunc)
}

func preRunFlagCheckOutput(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	if !funk.Contains([]string{"table", "json", "json-pretty", "csv"}, output) {
		return fmt.Errorf("invalid output format %s", output)
	}
	return nil
}

func renderListOutput(cmd *cobra.Command, data interface{}, tableWriter table.Writer, total int) (string, error) {
	if _, ok := cmd.Annotations[flagInitOutput]; !ok {
		panic("renderListOutput called for command where output flags were not initialized. This is a bug!")
	}

	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", err
	}
	switch outputFormat {
	case "table":
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
			if err == nil {
				tableWriter.SetAllowedRowLength(width)
			}
		}
		totalsMessage := ""
		if tableWriter.Length() != total {
			totalsMessage = fmt.Sprintf("\n(%d records out of %d)",
				tableWriter.Length(), total)
		} else {
			totalsMessage = fmt.Sprintf("\n(%d records)", total)
		}
		return tableWriter.Render() + totalsMessage, nil
	case "csv":
		return tableWriter.RenderCSV(), nil
	case "json":
		return renderJSON(data)
	case "json-pretty":
		return renderPrettyJSON(data)
	default:
		return "", fmt.Errorf("unsupported output format %s", outputFormat)
	}
}

func printListOutputAndError(cmd *cobra.Command, data interface{}, tableWriter table.Writer, total int, loopErr error) error {
	cmd.SilenceUsage = true
	result, err2 := renderListOutput(cmd, data, tableWriter, total)
	cmd.Println(result)
	if loopErr != nil {
		return processErrorResponse(loopErr)
	}
	return err2
}

func printMultiOpOutput(cmd *cobra.Command, itemName string, itemIDs interface{}, operationVerb string) {
	output := itemName
	length := 2
	niceIDs := fmt.Sprint(itemIDs)
	switch ids := itemIDs.(type) {
	case []string:
		length = len(ids)
		niceIDs = strings.Join(ids, ", ")
	case []strfmt.UUID:
		length = len(ids)
		niceIDs = strings.Join(
			funk.Map(
				ids,
				func(uuid strfmt.UUID) string {
					return string(uuid)
				}).([]string),
			", ")
	case []int:
		length = len(ids)
		niceIDs = strings.Join(
			funk.Map(
				ids,
				func(i int) string {
					return strconv.Itoa(i)
				}).([]string),
			", ")
	case []int64:
		length = len(ids)
		niceIDs = strings.Join(
			funk.Map(
				ids,
				func(i int64) string {
					return strconv.Itoa(int(i))
				}).([]string),
			", ")
	}
	if length != 1 {
		output = pluralize(output)
	}
	if length == 0 {
		output = "No " + strings.ToLower(output)
	} else {
		niceIDs += " "
	}
	output += " " + niceIDs + operationVerb
	cmd.Println(output)
}

func renderWatchOutput(cmd *cobra.Command, data interface{}, tableWriter table.Writer) (bool, string, error) {
	cmd.SilenceUsage = true
	if _, ok := cmd.Annotations[flagInitOutput]; !ok {
		panic("renderWatchOutput called for command where output flags were not initialized. This is a bug!")
	}

	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return false, "", err
	}
	switch outputFormat {
	case "table":
		if terminal.IsTerminal(int(os.Stdout.Fd())) {
			width, _, err := terminal.GetSize(int(os.Stdout.Fd()))
			if err == nil {
				tableWriter.SetAllowedRowLength(width)
			}
		}
		return true, tableWriter.Render(), nil
	case "csv":
		return false, tableWriter.RenderCSV(), nil
	case "json":
		o, err := renderJSON(data)
		return false, o, err
	case "json-pretty":
		o, err := renderPrettyJSON(data)
		return false, o, err
	default:
		return false, "", fmt.Errorf("unsupported output format %s", outputFormat)
	}
}

func renderJSON(data interface{}) (string, error) {
	var r []byte
	var err error
	r, err = json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func renderPrettyJSON(data interface{}) (string, error) {
	var r []byte
	var err error
	r, err = json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(r), nil
}

func pluralize(noun string) string {
	if strings.HasSuffix(noun, "s") {
		return noun
	}
	if strings.HasSuffix(noun, "y") {
		return noun[0:len(noun)-1] + "ies"
	}
	return noun + "s"
}
