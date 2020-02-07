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

	"github.com/fyde/fyde-cli/models"
)

func preRunCheckEndpoint(cmd *cobra.Command, args []string) error {
	if authViper.GetString(ckeyAuthEndpoint) == "" || global.Client == nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("endpoint not set! Run `%s endpoint [hostname]` first", ApplicationName)
	}

	return nil
}

func preRunCheckAuth(cmd *cobra.Command, args []string) error {
	err := preRunCheckEndpoint(cmd, args)
	if err != nil {
		return err
	}

	switch authViper.GetString(ckeyAuthMethod) {
	case authMethodBearerToken:
		if authViper.GetString(ckeyAuthAccessToken) == "" ||
			authViper.GetString(ckeyAuthClient) == "" ||
			authViper.GetString(ckeyAuthUID) == "" {
			cmd.SilenceUsage = true
			return fmt.Errorf("not logged in! Run `%s login` first", ApplicationName)
		}
	case "":
		fallthrough
	default:
		cmd.SilenceUsage = true
		return fmt.Errorf("not logged in! Run `%s login` first", ApplicationName)
	}

	return nil
}

type unprocessableEntityResponse interface {
	GetPayload() *models.UnprocessableEntityResponse
}

func processErrorResponse(err error) error {
	type unauthorizedResponse interface {
		GetPayload() *models.UnauthorizedResponse
	}

	type forbiddenResponse interface {
		GetPayload() models.ForbiddenResponse
	}

	type notFoundResponse interface {
		GetPayload() models.NotFoundResponse
	}

	switch r := err.(type) {
	case unauthorizedResponse:
		return fmt.Errorf(strings.Join(r.GetPayload().Errors, "\n"))
	case forbiddenResponse:
		return fmt.Errorf("forbidden")
	case notFoundResponse:
		return fmt.Errorf("not found")
	case unprocessableEntityResponse:
		return parseUnprocessableEntityResponse(r)
	default:
		return err
	}
}

func parseUnprocessableEntityResponse(r unprocessableEntityResponse) error {
	msgs := []string{}
	if r.GetPayload().Error != "" {
		msgs = []string{r.GetPayload().Error}
	}
	for k, v := range r.GetPayload().UnprocessableEntityResponse {
		if ifaceArray, ok := v.([]interface{}); ok {
			for _, arrElem := range ifaceArray {
				switch conv := arrElem.(type) {
				case string:
					msgs = append(msgs, fmt.Sprintf("%s: %s", k, conv))
				case map[string]interface{}:
					msgs = append(msgs, parseLimitsError(conv))
				}
			}
		}
	}
	return fmt.Errorf(strings.Join(msgs, "\n"))
}

func parseLimitsError(data map[string]interface{}) string {
	_, present := data["limit"]
	if _, present2 := data["kind"]; !present || !present2 {
		return ""
	}
	limitFloat, ok := data["limit"].(float64)
	if !ok {
		return ""
	}
	return fmt.Sprintf("limit of %d %ss reached", int(limitFloat), data["kind"])
}

func preRunFlagChecks(cmd *cobra.Command, args []string) error {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	if _, ok := cmd.Annotations[flagInitPagination]; ok {
		err := preRunFlagCheckPagination(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations[flagInitSort]; ok {
		err := preRunFlagCheckSort(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations[flagInitFilter]; ok {
		err := preRunFlagCheckFilter(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations[flagInitSearch]; ok {
		err := preRunFlagCheckSearch(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations[flagInitOutput]; ok {
		err := preRunFlagCheckOutput(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations[flagInitInput]; ok {
		err := preRunFlagCheckInput(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations[flagInitLoopControl]; ok {
		err := preRunFlagCheckLoopControl(cmd, args)
		if err != nil {
			return err
		}
	}

	return nil
}

type multiOpJSONResult struct {
	ID     interface{} `json:"id"`
	OK     bool        `json:"ok"`
	Result string      `json:"result"`
}

func multiOpBuildTableWriter() (table.Writer, []multiOpJSONResult) {
	tw := table.NewWriter()
	tw.Style().Format.Header = text.FormatDefault
	tw.AppendHeader(table.Row{
		"ID",
		"Result",
	})
	tw.SetAlign([]text.Align{
		text.AlignRight,
		text.AlignLeft})
	tw.SetAllowedColumnLengths([]int{36, 60})
	return tw, make([]multiOpJSONResult, 0)
}

func multiOpTableWriterAppend(tw table.Writer, j *[]multiOpJSONResult, id interface{}, result interface{}) {
	tw.AppendRow(table.Row{
		id,
		result,
	})

	_, isError := result.(error)
	r := multiOpJSONResult{
		OK:     !isError,
		Result: fmt.Sprint(result),
	}
	switch v := id.(type) {
	case int:
		r.ID = v
	case int64:
		r.ID = v
	default:
		r.ID = fmt.Sprint(id)
	}

	*j = append(*j, r)
}
