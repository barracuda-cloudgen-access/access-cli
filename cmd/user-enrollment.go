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

	if !multiOpCheckArgsPresent(cmd, args) {
		return fmt.Errorf("missing user ID argument")
	}

	return nil
}

// enrollmentGenerateCmd represents the generate command
var enrollmentGenerateCmd = &cobra.Command{
	Use:     "generate [user ID]...",
	Short:   "Generate user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		intArgs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"Slots",
			"Expiration",
			"URL",
		})
		tw.SetAllowedColumnLengths([]int{15, 10, 30, 140})

		createdList := []*apiusers.GenerateEnrollmentLinkCreatedBody{}

		for _, arg := range intArgs {
			params := apiusers.NewGenerateEnrollmentLinkParams()
			params.SetID(arg)

			resp, err := global.Client.Users.GenerateEnrollmentLink(params, global.AuthWriter)
			if err != nil {
				// best possible workaround for https://github.com/go-swagger/go-swagger/issues/1929
				// (without resorting to fixing the go-swagger code generator)
				if strings.Contains(err.Error(), "(*models.NotFoundResponse) is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface") {
					err = fmt.Errorf("user does not exist")
				}

				tw.AppendRow(table.Row{
					fmt.Sprintf("[ERR] %v", arg),
					"-",
					"-",
					processErrorResponse(err),
				})
				createdList = append(createdList, nil)

				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, createdList, tw, len(intArgs), err)
			}

			tw.AppendRow(table.Row{
				arg,
				resp.GetPayload().Count,
				resp.GetPayload().Expiration,
				resp.GetPayload().URL,
			})
			createdList = append(createdList, resp.Payload)
		}
		return printListOutputAndError(cmd, createdList, tw, len(intArgs), err)
	},
}

// enrollmentRevokeCmd represents the revoke command
var enrollmentRevokeCmd = &cobra.Command{
	Use:     "revoke [user ID]...",
	Short:   "Revoke user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		intArgs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		tw, j := multiOpBuildTableWriter()

		for _, arg := range intArgs {
			params := apiusers.NewRevokeEnrollmentLinkParams()
			params.SetID(arg)

			_, err = global.Client.Users.RevokeEnrollmentLink(params, global.AuthWriter)
			if err != nil {
				// best possible workaround for https://github.com/go-swagger/go-swagger/issues/1929
				// (without resorting to fixing the go-swagger code generator)
				if strings.Contains(err.Error(), "(*models.NotFoundResponse) is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface") {
					err = fmt.Errorf("user does not exist or does not have an enrollment link")
				}

				multiOpTableWriterAppend(tw, &j, arg, processErrorResponse(err))
				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, j, tw, len(args), err)
			}
			multiOpTableWriterAppend(tw, &j, arg, "success")
		}
		return printListOutputAndError(cmd, j, tw, len(args), err)
	},
}

// enrollmentGetCmd represents the get command
var enrollmentGetCmd = &cobra.Command{
	Use:     "get [user ID]...",
	Short:   "Get user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		intArgs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		cmd.SilenceUsage = true // errors beyond this point are no longer due to malformed input

		for _, arg := range intArgs {
			params := apiusers.NewGetUserParams()
			params.SetID(arg)

			resp, err := global.Client.Users.GetUser(params, global.AuthWriter)
			if err != nil {
				if loopControlContinueOnError(cmd) {
					cmd.PrintErrln(processErrorResponse(err))
					continue
				}
				return processErrorResponse(err)
			}

			if resp.Payload.Enrollment == nil ||
				resp.Payload.EnrollmentStatus == "revoked" ||
				resp.Payload.EnrollmentStatus == "expired" {
				cmd.Println("No shareable enrollment link available for this user")
			} else {
				cmd.Println(resp.Payload.Enrollment.URL)
			}
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

	initMultiOpArgFlags(enrollmentGenerateCmd, "user", "generate enrollments for", "id", "[]int64")
	initOutputFlags(enrollmentGenerateCmd)
	initLoopControlFlags(enrollmentGenerateCmd)

	initMultiOpArgFlags(enrollmentRevokeCmd, "user", "revoke enrollments for", "id", "[]int64")
	initOutputFlags(enrollmentRevokeCmd)
	initLoopControlFlags(enrollmentRevokeCmd)

	initMultiOpArgFlags(enrollmentGetCmd, "user", "get enrollments for", "id", "[]int64")
	initLoopControlFlags(enrollmentGetCmd)
}
