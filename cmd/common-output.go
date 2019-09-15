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

import (
	"encoding/json"
	"fmt"

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

func renderJSON(data interface{}) string {
	var r []byte
	var err error
	if global.Verbose {
		r, err = json.MarshalIndent(data, "", "  ")
	} else {
		r, err = json.Marshal(data)
	}
	if err != nil {
		return ""
	}
	return string(r)
}