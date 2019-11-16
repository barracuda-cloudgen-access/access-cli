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
	"strings"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/cobra"

	apiusers "github.com/fyde/fyde-cli/client/users"
)

// enrollmentCmd represents the enrollment command
var enrollmentCmd = &cobra.Command{
	Use:   "enrollment",
	Short: "Operations on user enrollment",
}

var enrollmentPreRunE = func(cmd *cobra.Command, args []string) error {
	err := preRunCheckAuth(cmd, args)
	if err != nil {
		return err
	}

	err = preRunFlagChecks(cmd, args)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("missing user ID argument")
	}

	return nil
}

// enrollmentGenerateCmd represents the generate command
var enrollmentGenerateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generate user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		params := apiusers.NewGenerateEnrollmentLinkParams()
		params.SetID(userID)

		resp, err := global.Client.Users.GenerateEnrollmentLink(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"Slots",
			"Expiration",
			"URL",
		})
		tw.SetAllowedColumnLengths([]int{10, 30, 140})

		tw.AppendRow(table.Row{
			resp.GetPayload().Count,
			resp.GetPayload().Expiration,
			resp.GetPayload().URL,
		})
		return printListOutputAndError(cmd, resp.Payload, tw, 1, err)
	},
}

// enrollmentRevokeCmd represents the revoke command
var enrollmentRevokeCmd = &cobra.Command{
	Use:     "revoke",
	Short:   "Revoke user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		params := apiusers.NewRevokeEnrollmentLinkParams()
		params.SetID(userID)

		cmd.SilenceUsage = true
		_, err = global.Client.Users.RevokeEnrollmentLink(params, global.AuthWriter)
		if err != nil {
			// best possible workaround for https://github.com/go-swagger/go-swagger/issues/1929
			// (without resorting to fixing the go-swagger code generator)
			if strings.Contains(err.Error(), "(*models.NotFoundResponse) is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface") {
				cmd.Println("User", userID, "does not exist or does not have an enrollment link")
				return nil
			}
			return processErrorResponse(err)
		}

		cmd.Println("Enrollment link for user", userID, "revoked")
		return nil
	},
}

// enrollmentGetCmd represents the get command
var enrollmentGetCmd = &cobra.Command{
	Use:     "get",
	Short:   "Get user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		params := apiusers.NewGetUserParams()
		params.SetID(userID)

		cmd.SilenceUsage = true
		resp, err := global.Client.Users.GetUser(params, global.AuthWriter)
		if err != nil {
			return processErrorResponse(err)
		}

		if resp.Payload.Enrollment == nil ||
			resp.Payload.EnrollmentStatus == "revoked" ||
			resp.Payload.EnrollmentStatus == "expired" {
			cmd.Println("No shareable enrollment link available for this user")
		} else {
			cmd.Println(resp.Payload.Enrollment.URL)
		}
		return nil
	},
}

func init() {
	usersCmd.AddCommand(enrollmentCmd)
	enrollmentCmd.AddCommand(enrollmentGenerateCmd)
	enrollmentCmd.AddCommand(enrollmentRevokeCmd)
	enrollmentCmd.AddCommand(enrollmentGetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// enrollmentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// enrollmentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(enrollmentGenerateCmd)
}
