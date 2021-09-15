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

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

// When adding a new override, write a doc about it using go/tpgtools-new-feature and add it to go/tpgtools-overrides.

type OverrideType string // enum

// Product-level Overrides
const (
	ProductBasePath OverrideType = "PRODUCT_BASE_PATH"
	ProductTitle    OverrideType = "PRODUCT_TITLE"
)

// Resource-level Overrides
const (
	VirtualField          OverrideType = "VIRTUAL_FIELD"
	CustomID                           = "CUSTOM_ID"
	CustomizeDiff                      = "CUSTOMIZE_DIFF"
	ImportFormat                       = "IMPORT_FORMAT"
	AppendToBasePath                   = "APPEND_TO_BASE_PATH"
	Mutex                              = "MUTEX"
	PreCreate                          = "PRE_CREATE_FUNCTION"
	PostCreate                         = "POST_CREATE_FUNCTION"
	PreDelete                          = "PRE_DELETE_FUNCTION"
	CustomResourceName                 = "CUSTOM_RESOURCE_NAME"
	NoSweeper                          = "NO_SWEEPER"
	CustomImport                       = "CUSTOM_IMPORT_FUNCTION"
	CustomCreateDirective              = "CUSTOM_CREATE_DIRECTIVE_FUNCTION"
	Undeletable                        = "UNDELETABLE"
	SkipDeleteFunction                 = "SKIP_DELETE_FUNCTION"
	SerializationOnly                  = "SERIALIZATION_ONLY"
	CustomSerializer                   = "CUSTOM_SERIALIZER"
	TerraformProductName               = "CUSTOM_TERRAFORM_PRODUCT_NAME"
	UseDCLID                           = "USE_DCL_ID"
)

// Field-level Overrides
const (
	CustomConfigMode     OverrideType = "CUSTOM_CONFIG_MODE"
	CustomDescription                 = "CUSTOM_DESCRIPTION"
	NamePrefix                        = "NAME_PREFIX"
	CustomName                        = "CUSTOM_NAME"
	CustomStateGetter                 = "CUSTOM_STATE_GETTER"
	CustomStateSetter                 = "CUSTOM_STATE_SETTER"
	CustomValidation                  = "CUSTOM_VALIDATION"
	Deprecated                        = "DEPRECATED"
	DiffSuppressFunc                  = "DIFF_SUPPRESS_FUNC"
	EnumBool                          = "ENUM_BOOL"
	Exclude                           = "EXCLUDE"
	CustomIdentityGetter              = "CUSTOM_IDENTITY_GETTER"
	Removed                           = "REMOVED"
	SetHashFunc                       = "SET_HASH_FUNC"
	CollapsedObject                   = "COLLAPSED_OBJECT"
	IgnoreRead                        = "IGNORE_READ"
	GenerateIfNotSet                  = "GENERATE_IF_NOT_SET"
	CustomListSize                    = "CUSTOM_LIST_SIZE_CONSTRAINT"
	CustomDefault                     = "CUSTOM_DEFAULT"
	CustomRequired                    = "REQUIRED_OVERRIDE"
)

// Overrides represents the type a resource's override file can be marshalled
// into.
type Overrides []Override

// Overrides handle minor quirks in behaviour for a resource (or one of its
// fields) by injecting modifications into the generated code for the resource.
// Every Override will have a Type; Overrides with a Field defined apply to that
// field in the resource and overrives without apply to the resource; Overrides
// may contain Details with structured metadata.
type Override struct {
	Type     OverrideType
	Field    *string     // may be nil
	Details  interface{} // may be nil
	Location *string     // may be nil
}

// ResourceOverride returns whether a single override with a single OverrideType
// is present on the resource.
func (o Overrides) ResourceOverride(typ OverrideType, location string) bool {
	found := false
	for _, v := range o {
		if v.Field == nil && v.Type == typ && compareLocation(v.Location, location) {
			if found {
				glog.Fatalf("found duplicate override of type %v", typ)
			}

			found = true
		}
	}

	return found
}

// ResourceOverrideWithDetails returns whether a single OverrideType is present
// on the resource, and includes the override's Details in the i interface if so.
func (o Overrides) ResourceOverrideWithDetails(typ OverrideType, i interface{}, location string) (bool, error) {
	found := false
	for _, v := range o {
		if v.Field == nil && v.Type == typ && compareLocation(v.Location, location) {
			if found {
				return false, fmt.Errorf("found duplicate override of type %v", typ)
			}

			found = true
			if err := convert(v.Details, i); err != nil {
				return false, fmt.Errorf("error converting type: %v", err)
			}
		}
	}

	return found, nil
}

// ResourceOverridesWithDetails returns all OverrideTypes of a given type on the
// resource, returning their details as yaml in an array.
// TODO: make this generic when Go supports generics
// ResourceOverrideWithDetails (the singular variant) is preferred when multiple
// overrides of a given type are not expected.
func (o Overrides) ResourceOverridesWithDetails(typ OverrideType, location string) (overrides []interface{}) {
	for _, v := range o {
		if v.Field == nil && v.Type == typ && compareLocation(v.Location, location) {
			overrides = append(overrides, v.Details)
		}
	}

	return overrides
}

// PropertyOverride returns whether a single override with a single OverrideType
// is present on a given property.
func (o Overrides) PropertyOverride(typ OverrideType, p Property, location string) bool {
	found := false
	for _, v := range o {
		if v.Field != nil && *v.Field == p.overridePath() && v.Type == typ && compareLocation(v.Location, location) {
			if found {
				glog.Fatalf("found duplicate override of type %v", typ)
			}

			found = true
		}
	}

	return found
}

// PropertyOverrideWithDetails returns whether a single OverrideType is present
// on a property, and includes the override's Details in the i interface if so.
func (o Overrides) PropertyOverrideWithDetails(typ OverrideType, p Property, i interface{}, location string) (bool, error) {
	found := false
	for _, v := range o {
		if v.Field != nil && *v.Field == p.overridePath() && v.Type == typ && compareLocation(v.Location, location) {
			if found {
				return false, fmt.Errorf("found duplicate override of type %v", typ)
			}

			found = true
			if err := convert(v.Details, i); err != nil {
				return false, fmt.Errorf("error converting type: %v", err)
			}
		}
	}

	return found, nil
}

// ProductWithDetails returns whether a single OverrideType is present
// on the resource, and includes the override's Details in the i interface if so.
func (o Overrides) ProductOverrideWithDetails(typ OverrideType, i interface{}) (bool, error) {
	found := false
	for _, v := range o {
		if v.Field == nil && v.Type == typ {
			if found {
				return false, fmt.Errorf("found duplicate override of type %v", typ)
			}

			found = true
			if err := convert(v.Details, i); err != nil {
				return false, fmt.Errorf("error converting type: %v", err)
			}
		}
	}

	return found, nil
}

func compareLocation(overrideLocation *string, objectLocation string) bool {
	// Overrides without a location field are considered valid for every location
	if overrideLocation == nil {
		return true
	}
	return *overrideLocation == objectLocation
}

func convert(item, out interface{}) error {
	bytes, err := yaml.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal: %v", err)
	}

	err = yaml.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}
