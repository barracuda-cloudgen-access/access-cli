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
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	api "github.com/barracuda-cloudgen-access/access-cli/client/user_directories"
	"github.com/barracuda-cloudgen-access/access-cli/models"
)

// userDirectoryGetCmd represents the get command
var userDirectoryGetCmd = &cobra.Command{
	Use:   "get [userdirectory ID]",
	Short: "Get user directory",
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
			return fmt.Errorf("missing user ID argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var userID int64
		var err error
		if cmd.Flags().Changed("id") {
			var d int
			d, err = cmd.Flags().GetInt("id")
			userID = int64(d)
		} else {
			userID, err = strconv.ParseInt(args[0], 10, 64)
		}
		if err != nil {
			return err
		}

		params := api.NewGetUserDirectoryParams()
		setTenant(cmd, params)
		params.SetID(userID)

		resp, err := global.Client.UserDirectories.GetUserDirectory(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := userDirectoryBuildTableWriter()
		userDirectoryTableWriterAppend(tw, resp.Payload.UserDirectory)

		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

func userDirectoryBuildTableWriter() table.Writer {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Name",
		"ShortCode",
		"DirectoryType",
		"TotalGroups",
		"TotalUsers",
		"LastSuccessfulSyncAt",
		"LastJobState",
		"LastJobErrors",
	})

	return tw
}

func userDirectoryTableWriterAppend(tw table.Writer, config models.UserDirectory) table.Writer {
	errors := ""
	state := "never_started"
	if config.LastJob != nil {
		state = config.LastJob.State
		errors = strings.Join(config.LastJob.Errors, ", ")
	}

	tw.AppendRow(table.Row{
		config.ID,
		config.Name,
		config.ShortCode,
		config.DirectoryType,
		config.TotalGroups,
		config.TotalUsers,
		config.LastSuccessfulSyncAt,
		state,
		errors,
	})
	return tw
}

func userDirectoryTableWriterAppendError(tw table.Writer, err error, id interface{}) {
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
		"-",
		"-",
		"-",
		"-",
	})
}

func init() {
	settingsUserDirectoryCmd.AddCommand(userDirectoryGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userDirectoryGetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userDirectoryGetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(userDirectoryGetCmd)
	initTenantFlags(userDirectoryGetCmd)

	userDirectoryGetCmd.Flags().Int("id", 0, "id of user to get")
}
