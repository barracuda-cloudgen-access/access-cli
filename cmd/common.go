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
	"math"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"

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

type pageable interface {
	SetPerPage(perPage *int64)
	SetPage(page *int64)
}

// forAllPages is a pagination helper
// all int64 usage is because go-swagger really likes int64
func forAllPages(params pageable, do func() (int64, error)) error {
	// func do must return the total number of items
	perPage := int64(50)

	total := int64(math.MaxInt64)
	var err error
	for curPage := int64(0); perPage*curPage < total; curPage++ {
		p := curPage + 1
		params.SetPage(&p)
		params.SetPerPage(&perPage)
		total, err = do()
		if err != nil {
			return err
		}
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
