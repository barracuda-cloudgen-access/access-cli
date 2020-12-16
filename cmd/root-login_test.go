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
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestLogin(t *testing.T) {
	defer gock.Off()

	gock.New(baseURIinTests()).
		MatchType("json").
		Post("/auth/sign_in").
		JSON(map[string]string{
			"email":    "test@example.com",
			"password": "testpw",
		}).
		Reply(200).
		SetHeader("access-token", "testAccessToken").
		SetHeader("client", "testClient").
		SetHeader("uid", "testUID")

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"login",
		"--email=test@example.com",
		"--password=testpw",
	})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	output, err := ioutil.ReadAll(buf)
	st.Expect(t, err, nil)
	if !strings.Contains(string(output), "Logged in successfully") {
		t.Fatal("Unexpected output")
	}

	st.Expect(t, authViper.GetString(ckeyAuthMethod), authMethodBearerToken)
	st.Expect(t, authViper.GetString(ckeyAuthAccessToken), "testAccessToken")
	st.Expect(t, authViper.GetString(ckeyAuthClient), "testClient")
	st.Expect(t, authViper.GetString(ckeyAuthUID), "testUID")
}
