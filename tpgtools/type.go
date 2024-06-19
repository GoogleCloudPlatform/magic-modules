// Copyright 2021 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/nasa9084/go-openapi"
)

const (
	SchemaTypeBool   = "TypeBool"
	SchemaTypeInt    = "TypeInt"
	SchemaTypeList   = "TypeList"
	SchemaTypeMap    = "TypeMap"
	SchemaTypeString = "TypeString"
	SchemaTypeFloat  = "TypeFloat"
	SchemaTypeSet    = "TypeSet"
)

// Type is the type of a TPG property.
type Type struct {
	// The raw string value of the type.
	typ *openapi.Schema
}

// IsObject returns whether a Type represents an object type. This is useful to
// disambiguate arrays from nested objects, both of which map to SchemaTypeList
// in terms of Terraform typing.
func (t Type) IsObject() bool {
	return t.String() == SchemaTypeList && t.typ.Type == "object"
}

func (t Type) IsEnum() bool {
	return len(t.typ.Enum) > 0
}

// Enum arrays are different from string arrays, and must be parsed differently
func (t Type) IsEnumArray() bool {
	return t.typ.Items != nil && len(t.typ.Items.Enum) > 0
}

// String returns the Terraform type of the Type.
func (t Type) String() string {
	switch t.typ.Type {
	case "boolean":
		return SchemaTypeBool
	case "string":
		return SchemaTypeString
	case "integer":
		return SchemaTypeInt
	case "number":
		if t.typ.Format == "double" {
			return SchemaTypeFloat
		}
		return "unknown number type"
	case "object":
		if t.typ.AdditionalProperties != nil {
			if v := t.typ.AdditionalProperties.Type; v == "string" {
				return SchemaTypeMap
			} else {
				// Complex maps are handled as sets with an extra value for the
				// name of the object
				return SchemaTypeSet
			}
		}
		return SchemaTypeList
	case "array":
		if t.typ.Extension["x-dcl-list-type"] == "set" {
			return SchemaTypeSet
		}
		return SchemaTypeList
	case "":
		return "<nil>"
	default:
		return fmt.Sprintf("undefined type: %s", t.typ)
	}
}

// IsDateTime returns true if the field in question is an openapi "date-time".
func (t Type) IsDateTime() bool {
	return t.typ.Format == "date-time"
}
