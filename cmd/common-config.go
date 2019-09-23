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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/shibukawa/configdir"
	"github.com/spf13/viper"
)

func setAuthDefaults() {
	authViper.SetDefault(ckeyAuthEndpoint, DefaultEndpoint)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgViper != nil {
		// already init (e.g. in tests)
		return
	}
	global.WriteFiles = true
	cfgViper = viper.New()
	if cfgFile != "" {
		// Use config file from the flag.
		cfgViper.SetConfigFile(cfgFile)
	} else {
		p := getUserConfigPath()

		// viper currently requires that config files exist in order to be able to write them
		// remove once https://github.com/spf13/viper/pull/723 is merged
		os.MkdirAll(p, os.ModePerm)
		fp := filepath.Join(p, ConfigFileName)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			ioutil.WriteFile(fp, []byte{}, os.FileMode(0644))
		}
		// ---

		cfgViper.AddConfigPath(p)
		dotIdx := strings.LastIndex(ConfigFileName, ".")
		cfgViper.SetConfigName(ConfigFileName[0:dotIdx])
		cfgViper.SetConfigType(ConfigFileName[dotIdx+1 : len(ConfigFileName)])
	}

	cfgViper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := cfgViper.ReadInConfig(); err == nil && global.VerboseLevel > 0 {
		fmt.Println("Using config file:", cfgViper.ConfigFileUsed())
	}
}

// initAuthConfig reads in credentials file and ENV variables if set.
func initAuthConfig() {
	if authViper != nil {
		// already init (e.g. in tests)
		return
	}
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
		fp := filepath.Join(p, AuthFileName)
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			ioutil.WriteFile(fp, []byte{}, os.FileMode(0644))
		}
		// ---

		authViper.AddConfigPath(p)
		dotIdx := strings.LastIndex(AuthFileName, ".")
		authViper.SetConfigName(AuthFileName[0:dotIdx])
		authViper.SetConfigType(AuthFileName[dotIdx+1 : len(AuthFileName)])
	}

	authViper.AutomaticEnv() // read in environment variables that match

	// If a credentials file is found, read it in.
	if err := authViper.ReadInConfig(); err == nil && global.VerboseLevel > 0 {
		fmt.Println("Using credentials file:", authViper.ConfigFileUsed())
	}
}

func getUserConfigPath() string {
	configDirs := configdir.New(ConfigVendorName, ConfigApplicationName)
	return configDirs.QueryFolders(configdir.Global)[0].Path
}
