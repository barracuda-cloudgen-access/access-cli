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
package main

import "github.com/fyde/fyde-cli/cmd"

var (
	// GitCommit is provided by govvv/goreleaser at compile-time
	GitCommit = "???"
	// BuildDate is provided by govvv/goreleaser at compile-time
	BuildDate = "???"
	// GitState is provided by govvv/goreleaser at compile-time
	GitState = "???"
	// Version is provided by govvv/goreleaser at compile-time
	Version = "???"
)

func main() {
	cmd.Execute(&cmd.VersionInformation{
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GitState:  GitState,
		Version:   Version,
	})
}
