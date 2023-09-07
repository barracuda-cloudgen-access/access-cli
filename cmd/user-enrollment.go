// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2023 Barracuda Networks, Inc.

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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/cobra"

	apiusers "github.com/barracuda-cloudgen-access/access-cli/client/users"
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
			"Classification",
			"Slots",
			"Expiration",
			"URL",
		})
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 15},
			{Number: 2, WidthMax: 10},
			{Number: 3, WidthMax: 10},
			{Number: 4, WidthMax: 30},
			{Number: 5, WidthMax: 140},
		})
		createdList := []*apiusers.GenerateEnrollmentLinkCreatedBody{}

		for _, arg := range intArgs {
			device_classification := cmd.Flag("classification").Value.String()
			count, err := cmd.Flags().GetInt("slots")
			if err != nil {
				return err
			}

			ref_count := int64(count)
			enrollment := apiusers.GenerateEnrollmentLinkBody{
				Enrollment: &apiusers.GenerateEnrollmentLinkParamsBodyEnrollment{
					DeviceClassification: &device_classification,
					Refcount:             &ref_count,
				},
			}
			params := apiusers.NewGenerateEnrollmentLinkParams()
			setTenant(cmd, params)
			params.SetID(arg)
			params.SetEnrollment(enrollment)
			if err != nil {
				return err
			}

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
				resp.GetPayload().DeviceClassification,
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
		device_classification := cmd.Flag("classification").Value.String()
		enrollment_id := strfmt.UUID(cmd.Flag("enrollment_id").Value.String())

		for _, arg := range intArgs {
			if enrollment_id == strfmt.UUID("") && device_classification != "" {
				enrollment_id, err = _enrollmentIdForClassification(cmd, arg, device_classification)
				if err != nil {
					multiOpTableWriterAppend(tw, &j, arg, processErrorResponse(err))
					if loopControlContinueOnError(cmd) {
						err = nil
						continue
					}
					return printListOutputAndError(cmd, j, tw, len(intArgs), err)
				}
			}
			if enrollment_id == strfmt.UUID("") {
				multiOpTableWriterAppend(tw, &j, arg, "no enrollment link found for the specified device classification")
				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, j, tw, len(intArgs), err)
			}

			params := apiusers.NewRevokeEnrollmentLinkParams()
			setTenant(cmd, params)
			params.SetID(arg)
			params.EnrollmentID = enrollment_id

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
				return printListOutputAndError(cmd, j, tw, len(intArgs), err)
			}
			multiOpTableWriterAppend(tw, &j, arg, "success")
		}
		return printListOutputAndError(cmd, j, tw, len(intArgs), err)
	},
}

func _enrollmentIdForClassification(cmd *cobra.Command, user_id int64, classification string) (strfmt.UUID, error) {

	enrollment_id := strfmt.UUID("")
	params := apiusers.NewGetUserParams()
	setTenant(cmd, params)
	params.SetID(user_id)

	resp, err := global.Client.Users.GetUser(params, global.AuthWriter)
	if err != nil {
		processErrorResponse(err)
		return enrollment_id, err
	}

	if len(resp.Payload.Enrollments) == 0 ||
		resp.Payload.EnrollmentStatus == "revoked" ||
		resp.Payload.EnrollmentStatus == "expired" {
		cmd.Println("No shareable enrollment link available for this user")
	} else {
		for _, enrollment := range resp.Payload.Enrollments {
			if enrollment.DeviceClassification == classification {
				enrollment_id = enrollment.ID
				break
			}
		}
	}

	if enrollment_id == strfmt.UUID("") {
		err = fmt.Errorf("no enrollment link found for the specified device classification")
	}
	return enrollment_id, err
}

// enrollmentChangeCmd represents the change command
var enrollmentChangeCmd = &cobra.Command{
	Use:   "change [user ID]...",
	Short: "Change user enrollment link slots",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := enrollmentPreRunE(cmd, args)
		if err != nil {
			return err
		}

		if !cmd.Flags().Changed("slots") {
			return fmt.Errorf("missing slots argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		intArgs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		slots, err := cmd.Flags().GetInt("slots")
		if err != nil {
			return err
		}

		classification, err := cmd.Flags().GetString("classification")
		if err != nil {
			return err
		}

		tw := table.NewWriter()
		tw.Style().Format.Header = text.FormatDefault
		tw.AppendHeader(table.Row{
			"ID",
			"Slots",
			"Classification",
			"Expiration",
			"URL",
		})
		tw.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, WidthMax: 15},
			{Number: 2, WidthMax: 5},
			{Number: 3, WidthMax: 14},
			{Number: 4, WidthMax: 30},
			{Number: 5, WidthMax: 140},
		})
		editedList := []*apiusers.ChangeEnrollmentLinkSlotsOKBody{}

		for _, arg := range intArgs {
			enrollment_id, err := _enrollmentIdForClassification(cmd, arg, classification)
			if err != nil {
				tw.AppendRow(table.Row{
					fmt.Sprintf("[ERR] %v", arg),
					"-",
					"-",
					"-",
					processErrorResponse(err),
				})
				editedList = append(editedList, nil)
				if loopControlContinueOnError(cmd) {
					cmd.PrintErrln(processErrorResponse(err))
					continue
				}
				return processErrorResponse(err)
			}

			change_params := apiusers.NewChangeEnrollmentLinkSlotsParams()
			setTenant(cmd, change_params)
			change_params.SetID(arg)
			change_params.SetEnrollmentID(enrollment_id)

			enrollment := &apiusers.ChangeEnrollmentLinkSlotsParamsBodyEnrollment{
				Refcount: int64(slots),
			}
			body := apiusers.ChangeEnrollmentLinkSlotsBody{
				UserID:     arg,
				Enrollment: enrollment,
			}
			change_params.SetRequest(body)

			change_resp, err := global.Client.Users.ChangeEnrollmentLinkSlots(change_params, global.AuthWriter)
			if err != nil {
				// best possible workaround for https://github.com/go-swagger/go-swagger/issues/1929
				// (without resorting to fixing the go-swagger code generator)
				if strings.Contains(err.Error(), "(*models.NotFoundResponse) is not supported by the TextConsumer, can be resolved by supporting TextUnmarshaler interface") {
					err = fmt.Errorf("user does not exist or does not have an enrollment link")
				}

				tw.AppendRow(table.Row{
					fmt.Sprintf("[ERR] %v", arg),
					"-",
					"-",
					processErrorResponse(err),
				})
				editedList = append(editedList, nil)

				if loopControlContinueOnError(cmd) {
					err = nil
					continue
				}
				return printListOutputAndError(cmd, editedList, tw, len(intArgs), err)
			}

			tw.AppendRow(table.Row{
				arg,
				change_resp.GetPayload().Count,
				change_resp.GetPayload().DeviceClassification,
				change_resp.GetPayload().Expiration,
				change_resp.GetPayload().URL,
			})
			editedList = append(editedList, change_resp.Payload)
		}
		return printListOutputAndError(cmd, editedList, tw, len(intArgs), err)
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

		jsonArray := make([]map[string]interface{}, 0)

		for _, arg := range intArgs {
			params := apiusers.NewGetUserParams()
			setTenant(cmd, params)
			params.SetID(arg)
			jsonObject := make(map[string]interface{})

			resp, err := global.Client.Users.GetUser(params, global.AuthWriter)
			if err != nil {
				if loopControlContinueOnError(cmd) {
					cmd.PrintErrln(processErrorResponse(err))
					continue
				}
				return processErrorResponse(err)
			}

			if len(resp.Payload.Enrollments) == 0 ||
				resp.Payload.EnrollmentStatus == "revoked" ||
				resp.Payload.EnrollmentStatus == "expired" {
				cmd.Println("No shareable enrollment link available for this user")
			} else {
				jsonObject["id"] = resp.Payload.ID
				for _, enrollment := range resp.Payload.Enrollments {
					key := enrollment.DeviceClassification
					value := enrollment.URL
					jsonObject[key] = value
				}
			}
			jsonArray = append(jsonArray, jsonObject)
		}

		jsonData, err := renderPrettyJSON(jsonArray)
		if err != nil {
			return processErrorResponse(err)
		}
		cmd.Println(string(jsonData))
		return nil
	},
}

// enrollmentEmailCmd represents the email command
var enrollmentEmailCmd = &cobra.Command{
	Use:     "email [user ID]...",
	Aliases: []string{"send"},
	Short:   "Send email with user enrollment link",
	PreRunE: enrollmentPreRunE,
	RunE: func(cmd *cobra.Command, args []string) error {
		intArgs, err := multiOpParseInt64Args(cmd, args, "id")
		if err != nil {
			return err
		}

		tw, j := multiOpBuildTableWriter()
		device_classification := cmd.Flag("classification").Value.String()
		slots, err := cmd.Flags().GetInt("slots")
		if err != nil {
			return err
		}

		device_classifications := &apiusers.SendEnrollmentEmailParamsBodyDeviceClassifications{}

		switch device_classification {
		case "supervised":
			device_classifications = &apiusers.SendEnrollmentEmailParamsBodyDeviceClassifications{
				Supervised: int64(slots),
			}
		case "managed":
			device_classifications = &apiusers.SendEnrollmentEmailParamsBodyDeviceClassifications{
				Managed: int64(slots),
			}
		case "personal":
			device_classifications = &apiusers.SendEnrollmentEmailParamsBodyDeviceClassifications{
				Personal: int64(slots),
			}
		default:
			return fmt.Errorf("invalid device classification")
		}

		for _, arg := range intArgs {
			params := apiusers.NewSendEnrollmentEmailParams()

			setTenant(cmd, params)
			params.SetID(arg)
			params.SetDeviceClassifications(apiusers.SendEnrollmentEmailBody{
				DeviceClassifications: device_classifications,
			})

			_, err = global.Client.Users.SendEnrollmentEmail(params, global.AuthWriter)
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
				return printListOutputAndError(cmd, j, tw, len(intArgs), err)
			}
			multiOpTableWriterAppend(tw, &j, arg, "success")
		}
		return printListOutputAndError(cmd, j, tw, len(intArgs), err)
	},
}

func init() {
	usersCmd.AddCommand(enrollmentCmd)
	enrollmentCmd.AddCommand(enrollmentGenerateCmd)
	enrollmentCmd.AddCommand(enrollmentRevokeCmd)
	enrollmentCmd.AddCommand(enrollmentChangeCmd)
	enrollmentCmd.AddCommand(enrollmentGetCmd)
	enrollmentCmd.AddCommand(enrollmentEmailCmd)

	initMultiOpArgFlags(enrollmentGenerateCmd, "user", "generate enrollments for", "id", "[]int64")
	initInputFlags(
		enrollmentGenerateCmd, "enrollment",
		inputField{
			Name:            "Device Classification",
			FlagName:        "classification",
			FlagDescription: "specify the device classification for the enrollment link",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "supervised",
		},
		inputField{
			Name:            "Slots",
			FlagName:        "slots",
			FlagDescription: "specify the number of slots for the enrollment link",
			VarType:         "int",
			Mandatory:       true,
			DefaultValue:    10,
			MainField:       true,
			SchemaName:      "slots",
		},
	)
	initOutputFlags(enrollmentGenerateCmd)
	initLoopControlFlags(enrollmentGenerateCmd)
	initTenantFlags(enrollmentGenerateCmd)

	initMultiOpArgFlags(enrollmentRevokeCmd, "user", "revoke enrollments for", "id", "[]int64")
	enrollmentRevokeCmd.Flags().String("classification", "supervised", "specify the classification")
	enrollmentRevokeCmd.Flags().String("enrollment_id", "", "the enrollment ID to revoke")
	initOutputFlags(enrollmentRevokeCmd)
	initLoopControlFlags(enrollmentRevokeCmd)
	initTenantFlags(enrollmentRevokeCmd)

	initMultiOpArgFlags(enrollmentChangeCmd, "user", "change enrollments for", "id", "[]int64")
	enrollmentChangeCmd.Flags().Int("slots", 10, "specify the new number of slots for the enrollment link")
	enrollmentChangeCmd.Flags().String("classification", "supervised", "specify the classification to change")

	initOutputFlags(enrollmentChangeCmd)
	initLoopControlFlags(enrollmentChangeCmd)
	initTenantFlags(enrollmentChangeCmd)

	initMultiOpArgFlags(enrollmentGetCmd, "user", "get enrollments for", "id", "[]int64")

	initLoopControlFlags(enrollmentGetCmd)
	initTenantFlags(enrollmentGetCmd)

	initMultiOpArgFlags(enrollmentEmailCmd, "user", "send enrollment emails. If provided classification for this user does not exist, create one", "id", "[]int64")
	enrollmentEmailCmd.Flags().Int("slots", 10, "specify the number of slots for the enrollment link")
	enrollmentEmailCmd.Flags().String("classification", "supervised", "specify the classification to send the enrollment email for")
	enrollmentEmailCmd.Flags().Bool("all", false, "send enrollment link to all classifications for the user")
	initOutputFlags(enrollmentEmailCmd)
	initLoopControlFlags(enrollmentEmailCmd)
	initTenantFlags(enrollmentEmailCmd)
}
