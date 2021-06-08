// Package serial implements access-cli serializables
package serial

import (
	"encoding/json"
)

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

type NullableOptionalString struct {
	Value *string
}

// MarshalJSON returns the NullableOptionalString as JSON
func (n NullableOptionalString) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

// UnmarshalJSON sets the NullableOptionalString from JSON
func (n *NullableOptionalString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*n = NullableOptionalString{Value: &value}
	return nil
}

// converting the struct to String format.
func (n NullableOptionalString) String() string {
	if n.Value == nil {
		return "<null>"
	}
	return *n.Value
}
