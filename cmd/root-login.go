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
	"strings"

	"github.com/fyde/fyde-cli/models"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"

	apiauth "github.com/fyde/fyde-cli/client/auth"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"log-in", "signin", "sign-in", "authenticate"},
	Short:   "Sign in to the console and store access token",
	PreRunE: preRunCheckEndpoint,
	RunE: func(cmd *cobra.Command, args []string) error {
		// ignoring errors, as we'll just ask for these if they are blank
		email, _ := cmd.Flags().GetString("email")
		password, _ := cmd.Flags().GetString("password")

		passwordfd, err := cmd.Flags().GetInt("password-fd")
		if err == nil && passwordfd >= 0 {
			file := os.NewFile(uintptr(passwordfd), "pipe")
			if file == nil {
				return fmt.Errorf("invalid file descriptor %d", passwordfd)
			}
			defer file.Close()
			cmd.Println("Reading password from file descriptor", passwordfd)
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
		}

		// read email from terminal, if not obtained by other means
		if email == "" {
			cmd.Print("Email address: ")
			i, err := fmt.Scanln(&email)
			if i == 0 || err != nil {
				return err
			}
		}

		// read password from terminal, if not obtained by other means
		if password == "" {
			cmd.Print("Password: ")
			passwordbytes, err := gopass.GetPasswd()
			if err != nil {
				return err
			}
			password = string(passwordbytes)
		}

		// send sign-in request
		params := apiauth.NewSignInParams()
		params.WithBody(&models.SignInRequest{
			Email:    email,
			Password: password,
		})
		signInResponse, err := global.Client.Auth.SignIn(params)
		if err != nil {
			return processErrorResponse(err)
		}

		// store access tokens
		authViper.Set(ckeyAuthAccessToken, signInResponse.AccessToken)
		authViper.Set(ckeyAuthClient, signInResponse.Client)
		authViper.Set(ckeyAuthUID, signInResponse.UID)
		authViper.Set(ckeyAuthMethod, authMethodBearerToken)

		if global.WriteFiles {
			err = authViper.WriteConfig()
			if err != nil {
				return err
			}
			cmd.Println("Logged in successfully, access token stored in", authViper.ConfigFileUsed())
		} else {
			cmd.Println("Logged in successfully")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringP("email", "e", "", "email address to use when logging in")
	loginCmd.Flags().IntP("password-fd", "d", -1, "read password from file descriptor, terminated by end of file, '\\r' or '\\n'.")
	loginCmd.Flags().StringP("password", "p", "", "password to use when logging in. Note that the password can be viewed by other processes. Prefer --password-fd instead.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
