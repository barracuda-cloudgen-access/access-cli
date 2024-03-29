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

const ckeyAuthEndpoint = "endpoint"
const ckeyAuthSkipTLSVerify = "skipTLSVerify"
const ckeyAuthUseInsecureHTTP = "useInsecureHTTP"
const ckeyAuthUseCache = "useCache"
const ckeyAuthMethod = "method"
const ckeyAuthAccessToken = "accessToken"
const ckeyAuthClient = "client"
const ckeyAuthUID = "uid"
const ckeyAuthCurrentTenant = "currentTenant"

const ckeyOutputFormat = "outputFormat"
const ckeyPipeOutputFormat = "pipeOutputFormat"
const ckeyRecordsPerGetRequest = "recordsPerGetRequest"
const ckeyDefaultRangeSize = "defaultRangeSize"
const ckeyCachePath = "cachePath"
