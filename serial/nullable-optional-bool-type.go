// Package serial implements access-cli serializables
package serial

import (
	"encoding/json"
	"log"
	"strconv"
)

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

type NullableOptionalBoolean struct {
	Value *bool
}

// MarshalJSON returns the NullableOptionalBoolean as JSON
func (n NullableOptionalBoolean) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

// UnmarshalJSON sets the NullableOptionalBoolean from JSON
func (n *NullableOptionalBoolean) UnmarshalJSON(data []byte) error {
	var value bool
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*n = NullableOptionalBoolean{Value: &value}
	return nil
}

// converting the struct to String format.
func (n NullableOptionalBoolean) String() string {
	if n.Value == nil {
		return "<null>"
	} else if *n.Value {
		return "true"
	} else {
		return "false"
	}
}

func (n *NullableOptionalBoolean) AssignFromString(s string) {
	if s == "null" {
		return
	}
	value, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal(err)
	}
	n.Value = &value
}
