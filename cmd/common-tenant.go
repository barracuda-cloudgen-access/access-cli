// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2023 Barracuda Networks, Inc. <hello@barracuda.com>

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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

type tenantable interface {
	SetTenantID(tenant strfmt.UUID)
}

func initTenantFlags(cmd *cobra.Command) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}
	cmd.Annotations[flagInitTenant] = "yes"
	cmd.Flags().StringP("tenant", "t", "", "tenant ID to perform operation on")
}

func setTenant(cmd *cobra.Command, t tenantable) {
	if _, ok := cmd.Annotations[flagInitTenant]; !ok {
		panic("setTenant called for command where tenant flag was not initialized. This is a bug!")
	}
	tenant, err := cmd.Flags().GetString("tenant")
	if err != nil || tenant == "" {
		tenant = global.CurrentTenant
	}
	t.SetTenantID(strfmt.UUID(tenant))
}
