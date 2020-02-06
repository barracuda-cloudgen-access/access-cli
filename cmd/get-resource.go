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
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"

	apiresources "github.com/fyde/fyde-cli/client/access_resources"
	"github.com/fyde/fyde-cli/models"
)

// resourceGetCmd represents the get command
var resourceGetCmd = &cobra.Command{
	Use:   "get [resource ID]",
	Short: "Get resource",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		if len(args) == 0 && !cmd.Flags().Changed("id") {
			return fmt.Errorf("missing resource ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var id string
		var err error
		if cmd.Flags().Changed("id") {
			id, err = cmd.Flags().GetString("id")
			if err != nil {
				return err
			}
		} else {
			id = args[0]
		}

		params := apiresources.NewGetResourceParams()
		params.SetID(strfmt.UUID(id))

		resp, err := global.Client.AccessResources.GetResource(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := resourceBuildTableWriter()
		resourceTableWriterAppend(tw, resp.Payload.AccessResource)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func resourceBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Public host",
		"Access policy",
		"Access proxy",
	})
	tw.SetAllowedColumnLengths([]int{36, 30, 30, 30, 30, 36})
	return tw
}

func resourceTableWriterAppend(tw table.Writer, resource models.AccessResource) {
	accessPolicies := strings.Join(funk.Map(resource.AccessPolicies, func(g *models.AccessResourceAccessPoliciesItems0) string {
		return g.Name
	}).([]string), ",")

	tw.AppendRow(table.Row{
		resource.ID,
		resource.Name,
		resource.PublicHost,
		accessPolicies,
		resource.AccessProxy.ID,
	})
}

func resourceTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	idStr := "[ERR]"
	if id != nil {
		idStr += fmt.Sprintf(" %v", id)
	}
	tw.AppendRow(table.Row{
		idStr,
		processErrorResponse(err),
		"-",
		"-",
		"-",
	})
}

func init() {
	resourcesCmd.AddCommand(resourceGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resourceGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resourceGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(resourceGetCmd)
	resourceGetCmd.Flags().String("id", "", "id of resource to get")
}
