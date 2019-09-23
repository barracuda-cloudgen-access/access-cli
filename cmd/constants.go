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

const (
	// ApplicationName is the user-facing name of the application
	ApplicationName = "fyde-cli"

	// ConfigFileName is the file name (not path, that's platform-dependent)
	// used for the config file
	ConfigFileName = "config.yaml"

	// AuthFileName is the file name (not path, that's platform-dependent)
	// used for the auth file
	AuthFileName = "auth.yaml"

	// ConfigVendorName is the vendor name used to select the default path for
	// configuration storage
	ConfigVendorName = "fyde"

	// ConfigApplicationName is the application name used to select the default
	// path for configuration storage
	ConfigApplicationName = "fyde-cli"

	// DefaultEndpoint is the default endpoint used by the client
	DefaultEndpoint = "enterprise.fyde.com"

	flagInitFilter     = "filter_flags_init"
	flagInitPagination = "pagination_flags_init"
	flagInitSort       = "sort_flags_init"
	flagInitOutput     = "output_flags_init"

	authMethodBearerToken = "bearerToken"
)
