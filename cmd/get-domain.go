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
	"strconv"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apiassets "github.com/fyde/fyde-cli/client/assets"
	"github.com/fyde/fyde-cli/models"
)

// domainGetCmd represents the get command
var domainGetCmd = &cobra.Command{
	Use:   "get [domain ID]",
	Short: "Get domain",
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
			return fmt.Errorf("missing domain ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var domainID int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			domainID = int64(d)
		} else {
			domainID, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		params := apiassets.NewGetAssetParams()
		params.SetID(domainID)

		resp, err := global.Client.Assets.GetAsset(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := domainBuildTableWriter()
		domainTableWriterAppend(tw, resp.Payload)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func domainBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"Category",
		"Asset source",
	})
	tw.SetAllowedColumnLengths([]int{15, 30, 30, 36})
	return tw
}

func domainTableWriterAppend(tw table.Writer, asset *models.Asset) {
	tw.AppendRow(table.Row{
		asset.ID,
		asset.Name,
		asset.Category,
		asset.AssetSourceID,
	})
}

func domainTableWriterAppendError(tw table.Writer, err error, id interface{}) {
	idStr := "[ERR]"
	if id != nil {
		idStr += fmt.Sprintf(" %v", id)
	}
	tw.AppendRow(table.Row{
		idStr,
		processErrorResponse(err),
		"-",
		"-",
	})
}

func init() {
	domainsCmd.AddCommand(domainGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(domainGetCmd)
	domainGetCmd.Flags().Int("id", 0, "id of domain to get")
}
