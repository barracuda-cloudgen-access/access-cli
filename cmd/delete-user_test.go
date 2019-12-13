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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gock.v1"
)

func TestDeleteUsersOneRequest(t *testing.T) {
	defer gock.Off()

	gock.New(baseURIinTests()).
		MatchType("json").
		Delete("/users/345,9845,2202").
		Reply(204)

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"users",
		"delete",
		"-o=json",
		"345",
		"9845",
		"2202",
	})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	output, err := ioutil.ReadAll(buf)
	st.Expect(t, err, nil)

	r := []multiOpJSONResult{}
	err = json.Unmarshal(output, &r)
	st.Expect(t, err, nil)
	st.Expect(t, len(r), 1)
	for _, o := range r {
		st.Expect(t, o.OK, true)
	}
}

func TestDeleteUsersIndividualRequests(t *testing.T) {
	defer gock.Off()

	gock.New(baseURIinTests()).
		MatchType("json").
		Delete("/users/345").
		Reply(204)

	gock.New(baseURIinTests()).
		MatchType("json").
		Delete("/users/9845").
		Reply(204)

	gock.New(baseURIinTests()).
		MatchType("json").
		Delete("/users/2202").
		Reply(204)

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"users",
		"delete",
		"--continue-on-error",
		"-o=json",
		"345",
		"9845",
		"2202",
	})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	output, err := ioutil.ReadAll(buf)
	st.Expect(t, err, nil)

	r := []multiOpJSONResult{}
	err = json.Unmarshal(output, &r)
	st.Expect(t, err, nil)
	st.Expect(t, len(r), 3)
	for _, o := range r {
		st.Expect(t, o.OK, true)
	}
}

func TestDeleteUsersIndividualRequestsOneFail(t *testing.T) {
	defer gock.Off()

	gock.New(baseURIinTests()).
		MatchType("json").
		Delete("/users/345").
		Reply(204)

	gock.New(baseURIinTests()).
		MatchType("json").
		Delete("/users/2202").
		Reply(204)

	cmd := rootCmd

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{
		"users",
		"delete",
		"--continue-on-error",
		"-o=json",
		"345",
		"9845",
		"2202",
	})
	err := cmd.Execute()
	if err != nil {
		t.Error(err)
	}

	output, err := ioutil.ReadAll(buf)
	st.Expect(t, err, nil)

	r := []multiOpJSONResult{}
	err = json.Unmarshal(output, &r)
	st.Expect(t, err, nil)
	st.Expect(t, len(r), 3)
	st.Expect(t, r[0].OK, true)
	st.Expect(t, r[1].OK, false)
	st.Expect(t, r[2].OK, true)
}
