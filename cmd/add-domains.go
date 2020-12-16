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
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"

	apiassets "github.com/fyde/access-cli/client/assets"
	"github.com/fyde/access-cli/models"
)

// domainsAddCmd represents the add command
var domainsAddCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"create", "new"},
	Short:   "Add domains",
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
		tw := domainBuildTableWriter()
		createdList := []*models.Asset{}
		total := 0
		err := forAllInput(cmd, args, true,
			func(values *inputEntry) (interface{}, error) { // do func
				total++ // this is the total of successful+failures, must increment before failure
				asset := &apiassets.CreateAssetBody{
					Category: "domain",
				}
				err := placeInputValues(cmd, values, asset,
					func(s string) { asset.Name = s },
					func(s string) { asset.AssetSourceID = strfmt.UUID(s) })
				if err != nil {
					return nil, err
				}
				params := apiassets.NewCreateAssetParams()
				params.SetAsset(*asset)

				resp, err := global.Client.Assets.CreateAsset(params, global.AuthWriter)
				if err != nil {
					return nil, err
				}
				return resp.Payload, nil
			}, func(data interface{}) { // printSuccess func
				asset := data.(*models.Asset)
				createdList = append(createdList, asset)
				domainTableWriterAppend(tw, asset)
			}, func(err error, id interface{}) { // doOnError func
				createdList = append(createdList, nil)
				domainTableWriterAppendError(tw, err, id)
			})
		return printListOutputAndError(cmd, createdList, tw, total, err)
	},
}

func init() {
	domainsCmd.AddCommand(domainsAddCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// domainsAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// domainsAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initOutputFlags(domainsAddCmd)
	initLoopControlFlags(domainsAddCmd)

	initInputFlags(domainsAddCmd, "domain",
		inputField{
			Name:            "Name",
			FlagName:        "name",
			FlagDescription: "specify the domain name for the created domain",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
			MainField:       true,
			SchemaName:      "name",
		},
		inputField{
			Name:            "Source",
			FlagName:        "source",
			FlagDescription: "specify the source ID for the created domain",
			VarType:         "string",
			Mandatory:       true,
			DefaultValue:    "",
		})
}
