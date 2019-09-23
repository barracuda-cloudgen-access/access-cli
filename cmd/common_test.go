// Package cmd implements fyde-cli commands
package cmd

import (
	"github.com/spf13/viper"
)

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

func init() {
	authViper = viper.New()
	cfgViper = viper.New()
	authViper.Set(ckeyAuthEndpoint, "mocked")
	authViper.Set(ckeyAuthMethod, authMethodBearerToken)
	authViper.Set(ckeyAuthAccessToken, "testAccessToken")
	authViper.Set(ckeyAuthClient, "testClient")
	authViper.Set(ckeyAuthUID, "test@example.com")
}

func baseURIinTests() string {
	return "https://mocked/api/v1"
}
