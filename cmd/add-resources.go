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
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiresources "github.com/barracuda-cloudgen-access/access-cli/client/access_resources"
	"github.com/barracuda-cloudgen-access/access-cli/models"
	"github.com/barracuda-cloudgen-access/access-cli/serial"
)

// resourcesAddCmd represents the add command
var resourcesAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add resources",
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
		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				resource := &apiresources.CreateResourceParamsBodyAccessResource{}
				resource.Enabled = true
				err := placeInputValues(cmd, values, resource,
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
					func(s string) { resource.Notes = s },
					func(s string) {
						resource.FixedLastOctet = &serial.NullableOptionalInt{}
						resource.FixedLastOctet.AssignFromString(s)
					})
				if err != nil {
					return nil, err
				}
				body := apiresources.CreateResourceBody{AccessResource: resource}
				params := apiresources.NewCreateResourceParams()
				setTenant(cmd, params)
				params.SetResource(body)

				resp, err := global.Client.AccessResources.CreateResource(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				resp := data.(*models.AccessResource)
				createdList = append(createdList, resp)
				resourceTableWriterAppend(tw, *resp)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				resourceTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func colonMappingToPortMapping(mapping string) *models.AccessResourcePortMapping {
	var protocol string
	internalPorts := []string{}
	parts := strings.SplitN(mapping, ":", 3)
	publicPorts := strings.Split(strings.TrimSpace(parts[0]), ",")

	if len(parts) > 1 {
		internalPorts = strings.Split(strings.TrimSpace(parts[1]), ",")
	}

	if len(parts) > 2 {
		protocol = strings.TrimSpace(parts[2])
	} else {
		protocol = "tcp"
	}

	return &models.AccessResourcePortMapping{
		InternalPorts: internalPorts,
		PublicPorts:   publicPorts,
		Protocol:      protocol,
	}
}

func init() {
	resourcesCmd.AddCommand(resourcesAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resourcesAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resourcesAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(resourcesAddCmd)
	initLoopControlFlags(resourcesAddCmd)
	initTenantFlags(resourcesAddCmd)
	initInputFlags(resourcesAddCmd, "resource",
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the name for the created resource",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "name",
		},
		inputField{
			Name:            "Public host",
			FlagName:        "public-host",
			FlagDescription: "specify the public host for the created resource",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Resource host",
			FlagName:        "resource-host",
			FlagDescription: "specify the resource host for the created resource",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		},
		inputField{
			Name:            "Port mappings",
			FlagName:        "ports",
			FlagDescription: "specify the port mappings (external:internal:protocol) for the created resource. Also accepts (external:internal), considers TCP by default.",
			VarType:         "[]string.skipcomma",
			Mandatory:       true,
			DefaultValue:    []string{},
		},
		inputField{
			Name:            "Proxy",
			FlagName:        "proxy",
			FlagDescription: "specify the proxy ID for the created resource",
			VarType:         "string",
			Mandatory:       true,
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
