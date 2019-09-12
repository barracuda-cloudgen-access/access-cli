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
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/shibukawa/configdir"
	"github.com/spf13/viper"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/oNaiPs/fyde-cli/client"
)

var cfgFile string
var authFile string

var cfgViper *viper.Viper
var authViper *viper.Viper

type globalInfo struct {
	Transport  *httptransport.Runtime
	Client     *apiclient.FydeEnterpriseConsole
	AuthWriter runtime.ClientAuthInfoWriter
}

var global globalInfo

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fyde-cli",
	Short: "Command-line client for the Fyde Enterprise Console",
	Long:  `fyde-cli allows access to all Enterprise Console APIs from the command line`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initAuthConfig)
	cobra.OnInitialize(initClient)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/fyde/fyde-cli/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&authFile, "auth", "", "credentials file (default is $HOME/.config/fyde/fyde-cli/auth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	cfgViper = viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		cfgViper.SetConfigFile(cfgFile)
	} else {
		p := getUserConfigPath()

		// viper currently requires that config files exist in order to be able to write them
		// remove once https://github.com/spf13/viper/pull/723 is merged
		os.MkdirAll(p, os.ModePerm)
		fp := filepath.Join(p, "config.yaml")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			ioutil.WriteFile(fp, []byte{}, os.FileMode(0644))
		}
		// ---

		cfgViper.AddConfigPath(p)
		cfgViper.SetConfigName("config")
		cfgViper.SetConfigType("yaml")
	}

	cfgViper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := cfgViper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cfgViper.ConfigFileUsed())
	}
}

// initAuthConfig reads in credentials file and ENV variables if set.
func initAuthConfig() {
	authViper = viper.New()
	setAuthDefaults()
	if authFile != "" {
		// Use config file from the flag.
		authViper.SetConfigFile(authFile)
	} else {
		p := getUserConfigPath()

		// viper currently requires that config files exist in order to be able to write them
		// remove once https://github.com/spf13/viper/pull/723 is merged
		os.MkdirAll(p, os.ModePerm)
		fp := filepath.Join(p, "auth.yaml")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			ioutil.WriteFile(fp, []byte{}, os.FileMode(0644))
		}
		// ---

		authViper.AddConfigPath(p)
		authViper.SetConfigName("auth")
		authViper.SetConfigType("yaml")
	}

	authViper.AutomaticEnv() // read in environment variables that match

	// If a credentials file is found, read it in.
	if err := authViper.ReadInConfig(); err == nil {
		fmt.Println("Using credentials file:", authViper.ConfigFileUsed())
	}
}

func getUserConfigPath() string {
	configDirs := configdir.New("fyde", "fyde-cli")
	return configDirs.QueryFolders(configdir.Global)[0].Path
}

func initClient() {
	endpoint := authViper.GetString(ckeyAuthEndpoint)
	if endpoint == "" {
		return
	}
	global.Transport = httptransport.New(endpoint, "/api/v1", nil)
	global.Client = apiclient.New(global.Transport, strfmt.Default)

	switch authViper.GetString(ckeyAuthMethod) {
	case "bearerToken":
		accessToken := authViper.GetString(ckeyAuthAccessToken)
		client := authViper.GetString(ckeyAuthClient)
		uid := authViper.GetString(ckeyAuthUID)
		global.AuthWriter = FydeAPIKeyAuth(accessToken, client, uid)
	default:
	}

}

// FydeAPIKeyAuth provides an API key auth info writer
func FydeAPIKeyAuth(accessToken, client, uid string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		err := r.SetHeaderParam("access-token", accessToken)
		if err != nil {
			return err
		}

		err = r.SetHeaderParam("client", client)
		if err != nil {
			return err
		}

		return r.SetHeaderParam("uid", uid)
	})
}
