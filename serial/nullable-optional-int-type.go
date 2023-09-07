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

type NullableOptionalInt struct {
	Value *int64
}

// MarshalJSON returns the NullableOptionalInt as JSON
func (n NullableOptionalInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

// UnmarshalJSON sets the NullableOptionalInt from JSON
func (n *NullableOptionalInt) UnmarshalJSON(data []byte) error {
	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*n = NullableOptionalInt{Value: &value}
	return nil
}

// converting the struct to String format.
func (n NullableOptionalInt) String() string {
	if n.Value == nil {
		return "<null>"
	}
	return strconv.FormatInt(*n.Value, 10)
}

func (n *NullableOptionalInt) AssignFromString(s string) {
	if s == "" || s == "null" {
		return
	}
	i, err := strconv.ParseUint(s, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	value := int64(i)
	n.Value = &value
}
