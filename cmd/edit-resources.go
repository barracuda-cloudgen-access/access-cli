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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiresources "github.com/fyde/fyde-cli/client/access_resources"
	"github.com/fyde/fyde-cli/models"
)

// resourcesEditCmd represents the get command
var resourcesEditCmd = &cobra.Command{
	Use:                "edit",
	Short:              "Edit resources",
	FParseErrWhitelist: cobra.FParseErrWhitelist{},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := preRunCheckAuth(cmd, args)
		if err != nil {
			return err
		}

		err = preRunFlagChecks(cmd, args)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		tw := resourceBuildTableWriter()
		createdList := []*models.AccessResource{}
		total := 0
		err := forAllInput(cmd, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apiresources.NewEditResourceParams()
				// IDs are not part of the request body, so we use this workaround
				resource := &struct {
					models.AccessResource
					ID              string      `json:"id"`
					AccessProxyID   strfmt.UUID `json:"access_proxy_id"`
					AccessPolicyIds []int64     `json:"access_policy_ids"`
				}{}
				err := placeInputValues(cmd, values, resource,
					func(s string) { resource.ID = s },
					func(s string) { resource.Name = s },
					func(s string) { resource.PublicHost = s },
					func(s string) { resource.InternalHost = s },
					func(s []string) { resource.Ports = s },
					func(s string) { resource.AccessProxyID = strfmt.UUID(s) },
					func(s int) {
						if s >= 0 {
							resource.AccessPolicyIds = []int64{int64(s)}
						}
					},
					func(s string) { resource.Notes = s })
				if err != nil {
					return nil, err
				}
				// here, map the ID from the "fake request body" to the correct place
				params.SetID(strfmt.UUID(resource.ID))
				body := apiresources.EditResourceBody{}
				body.AccessResource.AccessResource = resource.AccessResource
				body.AccessResource.AccessProxyID = resource.AccessProxyID
				body.AccessResource.AccessPolicyIds = resource.AccessPolicyIds
				body.AccessResource.Enabled = true
				params.SetResource(body)

				resp, err := global.Client.AccessResources.EditResource(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				resp := data.(*apiresources.EditResourceOKBody)
				createdList = append(createdList, &resp.AccessResource)
				resourceTableWriterAppend(tw, resp.AccessResource, resp.AccessProxyName)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				resourceTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	resourcesCmd.AddCommand(resourcesEditCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resourcesEditCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resourcesEditCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(resourcesEditCmd)
	initLoopControlFlags(resourcesEditCmd)

	initInputFlags(resourcesEditCmd,
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the resource to edit",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			IsIDOnError:     true,
			SchemaName:      "id",
		},
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the new name for the resource",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
			SchemaName:      "name",
		},
		inputField{
			Name:            "Public host",
			FlagName:        "public-host",
			FlagDescription: "specify the new public host for the resource",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Resource host",
			FlagName:        "resource-host",
			FlagDescription: "specify the new resource host for the resource",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Port mappings",
			FlagName:        "ports",
			FlagDescription: "specify the new port mappings (external:internal) for the resource",
			VarType:         "[]string",
			Mandatory:       true,
			DefaultValue:    []string{},
		},
		inputField{
			Name:            "Proxy",
			FlagName:        "proxy",
			FlagDescription: "specify the new proxy ID for the resource",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Policy",
			FlagName:        "policy",
			FlagDescription: "specify the new policy ID for the resource",
			VarType:         "int",
			Mandatory:       false,
			DefaultValue:    -1,
		},
		inputField{
			Name:            "Notes",
			FlagName:        "notes",
			FlagDescription: "specify notes for the resource",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		})
}
