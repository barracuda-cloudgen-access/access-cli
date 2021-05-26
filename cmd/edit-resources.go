// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2020 Barracuda Networks, Inc.

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
	"log"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiresources "github.com/barracuda-cloudgen-access/access-cli/client/access_resources"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/barracuda-cloudgen-access/access-cli/serial"
)

// resourcesEditCmd represents the edit command
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
		err := forAllInput(cmd, args, false,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				params := apiresources.NewEditResourceParams()
				setTenant(cmd, params)
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
					func(s []string) {
						resource.PortMappings = []*models.AccessResourcePortMapping{}
						for _, mapping := range s {
							resource.PortMappings = append(resource.PortMappings, colonMappingToPortMapping(mapping))
						}
					},
					func(s string) { resource.AccessProxyID = strfmt.UUID(s) },
					func(s []int64) {
						if len(s) > 0 {
							resource.AccessPolicyIds = s
						}
					},
					func(s []string) { resource.WildcardExceptions = s },
					func(s string) { resource.Notes = &s },
					func(s string) {
						resource.FixedLastOctet = &serial.NullableOptionalInt{}
						if s == "null" {
							return
						}
						i, err := strconv.ParseUint(s, 10, 8)
						if err != nil {
							log.Fatal(err)
						}
						value := int64(i)
						resource.FixedLastOctet.Value = &value
					})
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
				resourceTableWriterAppend(tw, resp.AccessResource)
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
	initTenantFlags(resourcesEditCmd)
	initInputFlags(resourcesEditCmd, "resource",
		inputField{
			Name:            "ID",
			FlagName:        "id",
			FlagDescription: "specify the ID of the resource to edit",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
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
			FlagDescription: "specify the port mappings (external:internal:protocol) for the created resource. Also accepts (external:internal), considers TCP by default.",
			VarType:         "[]string.skipcomma",
			Mandatory:       false,
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
			Name:            "Policies",
			FlagName:        "policies",
			FlagDescription: "specify a list of comma-separated policy IDs for the created resource",
			VarType:         "[]int",
			Mandatory:       false,
			DefaultValue:    []int{},
		},
		inputField{
			Name:            "Wildcard Exceptions",
			FlagName:        "exceptions",
			FlagDescription: "specify a list of of sub-domain wildcard exceptions that wont be proxied over (comma separated)",
			VarType:         "[]string",
			Mandatory:       false,
			DefaultValue:    []string{},
		},
		inputField{
			Name:            "Notes",
			FlagName:        "notes",
			FlagDescription: "specify notes for the resource",
			VarType:         "string",
			Mandatory:       false,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Fixed Last Octet",
			FlagName:        "fixed-last-octet",
			FlagDescription: "forces the agent to bind the resource to a local IP in the format 192.0.2.X (null to disable)",
			VarType:         "string", // use string to read "null" pseudo value
			Mandatory:       false,
			DefaultValue:    "",
		})
}
