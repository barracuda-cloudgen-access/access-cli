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

	"github.com/spf13/cobra"

	apiusers "github.com/oNaiPs/fyde-cli/client/users"
)

func preRunCheckEndpoint(cmd *cobra.Command, args []string) error {
	if authViper.GetString(ckeyAuthEndpoint) == "" || global.Client == nil {
		return fmt.Errorf("endpoint not set! Run `fyde-cli endpoint [hostname]` first")
	}

	return nil
}

func preRunCheckAuth(cmd *cobra.Command, args []string) error {
	err := preRunCheckEndpoint(cmd, args)
	if err != nil {
		return err
	}

	switch authViper.GetString(ckeyAuthMethod) {
	case "bearerToken":
		if authViper.GetString(ckeyAuthAccessToken) == "" ||
			authViper.GetString(ckeyAuthClient) == "" ||
			authViper.GetString(ckeyAuthUID) == "" {
			return fmt.Errorf("not logged in! Run `fyde-cli login` first")
		}
	case "":
	default:
		return fmt.Errorf("not logged in! Run `fyde-cli login` first")
	}

	return nil
}

func processErrorResponse(err error) error {
	// TODO prepare for other error response types
	// (maybe use reflection if we can always get the Payload from within the error type)
	switch r := err.(type) {
	case *apiusers.ListUsersUnauthorized:
		return fmt.Errorf(strings.Join(r.Payload.Errors, "\n"))
	default:
		return err
	}
}

func preRunFlagChecks(cmd *cobra.Command, args []string) error {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	if _, ok := cmd.Annotations["pagination_flags_init"]; ok {
		err := preRunFlagCheckPagination(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations["sort_flags_init"]; ok {
		err := preRunFlagCheckSort(cmd, args)
		if err != nil {
			return err
		}
	}

	if _, ok := cmd.Annotations["output_flags_init"]; ok {
		err := preRunFlagCheckOutput(cmd, args)
		if err != nil {
			return err
		}
	}

	return nil
}

func int64min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
