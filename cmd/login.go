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
	"strings"

	"github.com/howeyc/gopass"
	"github.com/oNaiPs/fyde-cli/models"
	"github.com/spf13/cobra"

	apiauth "github.com/oNaiPs/fyde-cli/client/auth"
)

var username string
var password string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Short:   "Sign in to the console and store access token",
	PreRunE: preRunCheckEndpoint,
	RunE: func(cmd *cobra.Command, args []string) error {
		passwordfd, err := cmd.Flags().GetInt("password-fd")
		if err == nil && passwordfd >= 0 {
			// TODO reading from FD is broken, figure out why later
			file := os.NewFile(uintptr(passwordfd), "pipe")
			if file == nil {
				return fmt.Errorf("invalid file descriptor %d", passwordfd)
			}
			defer file.Close()
			fmt.Println("Reading password from file descriptor", passwordfd)
			pwbytes, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			endIdx := strings.IndexAny(string(pwbytes), "\n\r")
			if endIdx >= 0 {
				password = string(pwbytes)[0:endIdx]
			} else {
				password = string(pwbytes)
			}
			fmt.Println(password)
		}

		// read username from terminal, if not obtained by other means
		if username == "" {
			fmt.Print("Username: ")
			i, err := fmt.Scanln(&username)
			if i == 0 || err != nil {
				return err
			}
		}

		// read password from terminal, if not obtained by other means
		if password == "" {
			fmt.Print("Password: ")
			passwordbytes, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			password = string(passwordbytes)
		}

		// send sign-in request
		params := apiauth.NewSignInParams()
		params.WithBody(&models.SignInRequest{
			Email:    username,
			Password: password,
		})
		signInResponse, err := global.Client.Auth.SignIn(params)
		if err != nil {
			return err
		}

		// store access tokens
		authViper.Set(ckeyAuthAccessToken, signInResponse.AccessToken)
		authViper.Set(ckeyAuthClient, signInResponse.Client)
		authViper.Set(ckeyAuthUID, signInResponse.UID)
		authViper.Set(ckeyAuthMethod, "bearerToken")

		err = authViper.WriteConfig()
		if err != nil {
			return err
		}
		fmt.Println("Logged in successfully, access token stored in", authViper.ConfigFileUsed())
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "username to use when logging in")
	loginCmd.Flags().IntP("password-fd", "d", -1, "read password from file descriptor, terminated by end of file, '\\r' or '\\n'.")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password to use when logging in. Note that the password can be viewed by other processes. Prefer --password-fd instead.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
