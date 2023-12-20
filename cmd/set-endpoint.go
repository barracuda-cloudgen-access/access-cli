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
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

// endpointSetCmd represents the endpoint set command
var endpointSetCmd = &cobra.Command{
	Use:   "set [endpoint]",
	Short: "Set console endpoint to use",
	Long: `Set console endpoint to use.
Valid values for [endpoint]:
  - eu
  - us
  - HTTP/HTTPS URL
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("missing endpoint argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var endpointUrl string
		// if someone passes in a URL, ensure we only extract user:pass@host:port without protocol, slashes, etc.
		re := regexp.MustCompile(`^(?:https?:(?:\/\/)?)?([^\/?\n]+)`)
		switch strings.ToLower(args[0]) {
		case "eu":
			endpointUrl = "api.eu.access.barracuda.com"
		case "us":
			endpointUrl = "api.us.access.barracuda.com"
		default:
			endpointUrl = re.FindStringSubmatch(args[0])[1]
		}

		authViper.Set(ckeyAuthAccessToken, "")
		authViper.Set(ckeyAuthClient, "")
		authViper.Set(ckeyAuthUID, "")
		authViper.Set(ckeyAuthMethod, "")
		authViper.Set(ckeyAuthCurrentTenant, "")
		authViper.Set(ckeyAuthEndpoint, endpointUrl)

		insecureSkipVerify, _ := cmd.Flags().GetBool("insecure-skip-verify")
		authViper.Set(ckeyAuthSkipTLSVerify, insecureSkipVerify)

		insecureUseHTTP, _ := cmd.Flags().GetBool("insecure-use-http")
		authViper.Set(ckeyAuthUseInsecureHTTP, insecureUseHTTP)

		useCache, _ := cmd.Flags().GetBool("experimental-use-cache")
		authViper.Set(ckeyAuthUseCache, useCache)

		path := cfgViper.GetString(ckeyCachePath)
		path = filepath.Join(path, "httpcache")
		os.RemoveAll(path)

		if global.WriteFiles {
			err := authViper.WriteConfig()
			if err != nil {
				return err
			}
		}
		cmd.Printf("Endpoint set to %s.\nCredentials cleared, please login again using `%s login`\n", endpointUrl, ApplicationName)
		if insecureUseHTTP {
			cmd.Println("WARNING: HTTP, instead of HTTPS, is being used for API communication. THIS IS INSECURE.")
		} else if insecureSkipVerify {
			cmd.Println("WARNING: TLS certificate verification is being skipped for the endpoint. THIS IS INSECURE.")
		}
		return nil
	},
}

func init() {
	clusterCmd.AddCommand(endpointSetCmd)

	endpointSetCmd.Flags().Bool("insecure-skip-verify", false, "Skip TLS certificate verification for the endpoint. INSECURE, use only if you know what you are doing")
	endpointSetCmd.Flags().Bool("insecure-use-http", false, "Communicate with the management console over HTTP instead of HTTPS. INSECURE, use only if you know what you are doing")
	endpointSetCmd.Flags().Bool("experimental-use-cache", false, "Enable HTTP response caching (experimental)")
}
