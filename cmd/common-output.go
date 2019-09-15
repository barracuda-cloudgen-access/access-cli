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

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

func initOutputFlags(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations["output_flags_init"] = "yes"
	cmd.Flags().StringP("output", "o", "table", "output format (table, json or csv)")
}

func preRunFlagCheckOutput(cmd *cobra.Command, args []string) error {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		return err
	}
	if !funk.Contains([]string{"table", "json", "csv"}, output) {
		return fmt.Errorf("invalid output format %s", output)
	}
	return nil
}

func renderListOutput(cmd *cobra.Command, data interface{}, tableWriter table.Writer) (string, error) {
	if _, ok := cmd.Annotations["output_flags_init"]; !ok {
		panic("renderListOutput called for command where output flags were not initialized. This is a bug!")
	}

	outputFormat, err := cmd.Flags().GetString("output")
	if err != nil {
		return "", err
	}
	switch outputFormat {
	case "table":
		return tableWriter.Render(), nil
	case "csv":
		return tableWriter.RenderCSV(), nil
	case "json":
		return renderJSON(data)
	default:
		return "", fmt.Errorf("unsupported output format %s", outputFormat)
	}
}

func renderJSON(data interface{}) (string, error) {
	var r []byte
	var err error
	if global.Verbose {
		r, err = json.MarshalIndent(data, "", "  ")
	} else {
		r, err = json.Marshal(data)
	}
	if err != nil {
		return "", err
	}
	return string(r), nil
}
