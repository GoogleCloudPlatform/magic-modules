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
		// assume if this is set, it's a string -> string map for now.
		// https://swagger.io/docs/specification/data-models/dictionaries/
		// describes the behaviour of AdditionalProperties for type: object
		if t.typ.AdditionalProperties != nil {
			if v := t.typ.AdditionalProperties.Type; v == "string" {
				return SchemaTypeMap
			} else {
				return fmt.Sprintf("unknown AdditionalProperties: %q", v)
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
