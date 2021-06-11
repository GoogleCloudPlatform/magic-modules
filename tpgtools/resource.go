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
	"sort"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/nasa9084/go-openapi"
)

const GoPkgTerraformSdkValidation = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

// Resource is tpgtools' model of what a information is necessary to generate a
// resource in TPG.
type Resource struct {
	productMetadata *ProductMetadata

	// ID is the Terraform resource id format as a pattern string. Additionally,
	// import formats can be derived from it.
	ID string

	// ImportFormats are pattern format strings for importing the Terraform resource.
	// TODO: if none are set, the resource does not support import.
	ImportFormats []string

	// title is the name of the resource in snake_case. For example,
	// "instance", "backend_service".
	title string

	// dclname is the name of the DCL resource in snake_case. For example,
	// "instance", "backend_service".
	dclname string

	// Description of the Terraform resource
	Description string

	// Lock name for a mutex to prevent concurrent API calls for a given resource
	Mutex string

	// Properties are the fields of a resource. Properties may be nested.
	Properties []Property

	// InsertTimeoutMinutes is the timeout value in minutes for the resource
	// create operation
	InsertTimeoutMinutes int

	// UpdateTimeoutMinutes is the timeout value in minutes for the resource
	// update operation
	UpdateTimeoutMinutes int

	// DeleteTimeoutMinutes is the timeout value in minutes for the resource
	// delete operation
	DeleteTimeoutMinutes int

	// PreCreateFunction is the name of a function that's called before the
	// Creation call for a resource- specifically, before the id is recorded.
	PreCreateFunction *string

	// PostCreateFunction is the name of a function that's called immediately
	// after the Creation call for a resource.
	PostCreateFunction *string

	// PreDeleteFunction is the name of a function that's called immediately
	// prior to the Delete call for a resource.
	PreDeleteFunction *string

	// CustomImportFunction is the name of a function that's called in place
	// of the standard import template code
	CustomImportFunction *string

	// CustomizeDiff is a list of functions to set as the Terraform schema
	// CustomizeDiff field
	CustomizeDiff []string

	// List of other Golang packages to import in a resources' generated Go file
	additionalFileImportSet stringset.Set

	// ListFields is the list of fields required for a list call.
	ListFields []string

	// HasSweeper says if this resource has a generated sweeper.
	HasSweeper bool

	// These are all of the reused types.
	ReusedTypes []Property

	// If this resource requires a state hint to update correctly
	StateHint bool

	// CustomCreateDirectiveFunction is the name of a function that takes the
	// object to be created and returns a list of directive to use for the apply
	// call
	CustomCreateDirectiveFunction *string

	// SkipDeleteFunction is the name of a function that takes the
	// object and config and returns a boolean for if Terraform should make
	// the delete call for the resource
	SkipDeleteFunction *string
}

// Name is the shortname of a resource. For example, "instance".
func (r Resource) Name() string {
	return r.title
}

func (r Resource) DCLName() string {
	if r.dclname != "" {
		return r.dclname
	}
	return r.title
}

// Path is the provider name of a resource, product_name. For example,
// "cloud_run_service".
func (r Resource) Path() string {
	return r.Package() + "_" + r.Name()
}

// TerraformName is the Terraform resource type used in HCL configuration.
// For example, "google_compute_instance"
func (r Resource) TerraformName() string {
	return "google_" + r.Path()
}

// Type is the title-cased name of a resource, for printing information about
// the type". For example, "Instance".
func (r Resource) Type() string {
	return snakeToTitleCase(r.DCLName())
}

// PathType is the title-cased name of a resource preceded by it's package, for
// often used to namespace functions. For example, "RedisInstance".
func (r Resource) PathType() string {
	return snakeToTitleCase(r.Path())
}

// Package is the namespace of the package within the dcl
// the Package is normally a lowercase variant of ProductName
func (r Resource) Package() string {
	return r.productMetadata.PackageName
}

// ProductType is the title-cased product name of a resource. For example,
// "NetworkServices".
func (r Resource) ProductType() string {
	return r.productMetadata.ProductType()
}

func (r Resource) ProductMetadata() *ProductMetadata {
	copy := *r.productMetadata
	return &copy
}

// DCLPackage is the package name of the DCL client library to use for this
// resource. For example, the Package "access_context_manager" would have a
// DCLPackage of "accesscontextmanager"
func (r Resource) DCLPackage() string {
	return strings.Replace(r.productMetadata.PackagePath, "_", "", -1)
}

// SidebarCurrent is the website sidebar identifier, for example
// docs-google-compute-instance
// TODO: is this still needed?
func (r Resource) SidebarCurrent() string {
	return "docs-" + strings.Replace(r.TerraformName(), "_", "-", -1)
}

// Updatable returns true if the resource should have an update method.
// This will avoid the error message:
// "All fields are ForceNew or Computed w/out Optional, Update is superfluous"
func (r Resource) Updatable() bool {
	for _, p := range r.SchemaProperties() {
		if !p.ForceNew && !(!p.Optional && p.Computed) {
			return true
		}
	}
	return false
}

// Objects are properties with sub-properties
func (r Resource) Objects() (props []Property) {
	for _, v := range r.Properties {
		if len(v.Properties) != 0 {
			// If this property uses a reused type, add it one-time afterwards to avoid multiple creations.
			if v.ref != "" {
				continue
			}
			props = append(props, v)
			props = append(props, v.Objects()...)
		}
	}

	for _, v := range r.ReusedTypes {
		props = append(props, v)
		props = append(props, v.Objects()...)
	}

	return props
}

// SchemaProperties is the list of resource properties filtered to those included in schema.
func (r Resource) SchemaProperties() (props []Property) {
	return collapsedProperties(r.Properties)
}

// Enum arrays are not arrays of strings in the DCL and require special
// conversion
func (r Resource) EnumArrays() (props []Property) {
	// Top level
	for _, v := range r.Properties {
		if v.Type.typ.Items != nil && len(v.Type.typ.Items.Enum) > 0 {
			props = append(props, v)
		}
	}
	// Look for nested
	for _, n := range r.Objects() {
		for _, v := range n.Properties {
			if v.Type.typ.Items != nil && len(v.Type.typ.Items.Enum) > 0 {
				props = append(props, v)
			}
		}
	}

	return props
}

// AdditionalGoPackages returns a sorted list of additional Go packages to import.
func (r Resource) AdditionalFileImports() []string {
	sl := make([]string, 0, len(r.additionalFileImportSet))
	for k := range r.additionalFileImportSet {
		sl = append(sl, k)
	}
	sort.Strings(sl)
	return sl
}

// If this resource has a server generated field that is used to read the
// resource. This must be set during create
func (r Resource) HasServerGeneratedName() bool {
	identityFields := idParts(r.ID)
	for _, p := range r.Properties {
		if stringInSlice(p.Name(), identityFields) {
			if !p.Settable {
				return true
			}
		}
	}
	return false
}

// SweeperName returns the name of the Sweeper for this resource.
func (r Resource) SweeperName() string {
	return r.ProductType() + strings.Title(r.Name())
}

// SweeperFunctionArgs returns a list of arguments for calling a sweeper function.
func (r Resource) SweeperFunctionArgs() string {
	vals := make([]string, 0)
	for _, v := range r.ListFields {
		vals = append(vals, fmt.Sprintf("d[\"%s\"]", v))
	}

	if len(vals) > 0 {
		return strings.Join(vals, ",") + ","
	} else {
		return ""
	}
}

// Returns the name of the ID function for the Terraform resource.
func (r Resource) IdFunction() string {
	for _, p := range r.Properties {
		if p.forwardSlashAllowed {
			return "replaceVars"
		}
	}
	return "replaceVarsForId"
}

// ResourceInput is a Resource along with additional generation metadata.
type ResourceInput struct {
	Resource
}

// RegisterReusedType adds a new reused type if the type does not already exist.
func (r Resource) RegisterReusedType(p Property) []Property {
	found := false
	for _, v := range r.ReusedTypes {
		if v.ref == p.ref {
			found = true
		}
	}
	if !found {
		r.ReusedTypes = append(r.ReusedTypes, p)
	}
	return r.ReusedTypes
}

func createResource(schema *openapi.Schema, typeFetcher *TypeFetcher, overrides Overrides, product *ProductMetadata, location string) (*Resource, error) {
	resourceTitle := schema.Title

	// Attempt to construct the resource name using location. Other than
	// zonal resources (which never include "zone"), there is no consistency
	// for when to include the location in the resource name.
	// A resource name override will often need to be used for on of the localized
	// resource versions.
	if location != "zone" {
		resourceTitle = location + resourceTitle
	}
	res := Resource{
		title:                jsonToSnakeCase(resourceTitle),
		dclname:              jsonToSnakeCase(schema.Title),
		productMetadata:      product,
		Description:          schema.Description,
		InsertTimeoutMinutes: 10,
		UpdateTimeoutMinutes: 10,
		DeleteTimeoutMinutes: 10,
	}

	crname := CustomResourceNameDetails{}
	crnameOk, err := overrides.ResourceOverrideWithDetails(CustomResourceName, &crname, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom resource name details: %v", err)
	}

	if crnameOk {
		res.title = jsonToSnakeCase(crname.Title)
	}

	id, err := findResourceId(schema, overrides, location)
	if err != nil {
		return nil, err
	}
	res.ID = id

	// Resource Override: Custom Import Function
	cifd := CustomImportFunctionDetails{}
	cifdOk, err := overrides.ResourceOverrideWithDetails(CustomImport, &cifd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom import function details: %v", err)
	}
	if cifdOk {
		res.CustomImportFunction = &cifd.Function
	}

	// Resource Override: Import formats
	ifd := ImportFormatDetails{}
	ifdOk, err := overrides.ResourceOverrideWithDetails(ImportFormat, &ifd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode import format details: %v", err)
	}
	if ifdOk {
		res.ImportFormats = ifd.Formats
	} else {
		res.ImportFormats = defaultImportFormats(res.ID)
	}

	// Resource Override: Mutex
	md := MutexDetails{}
	mdOk, err := overrides.ResourceOverrideWithDetails(Mutex, &md, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode mutex details: %v", err)
	}
	if mdOk {
		res.Mutex = md.Mutex
	}

	props, err := createPropertiesFromSchema(schema, typeFetcher, overrides, &res, nil, location)
	if err != nil {
		return nil, err
	}

	res.Properties = props

	// Resource Override: Virtual field
	for _, vfd := range overrides.ResourceOverridesWithDetails(VirtualField, location) {
		vf := VirtualFieldDetails{}
		if err := convert(vfd, &vf); err != nil {
			return nil, fmt.Errorf("error converting type: %v", err)
		}

		res.Properties = append(res.Properties, readVirtualField(vf))
	}

	// Resource-level pre and post action functions
	preCreate := PreCreateFunctionDetails{}
	preCreOk, err := overrides.ResourceOverrideWithDetails(PreCreate, &preCreate, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pre create function details: %v", err)
	}
	if preCreOk {
		res.PreCreateFunction = &preCreate.Function
	}

	postCreate := PostCreateFunctionDetails{}
	postCreOk, err := overrides.ResourceOverrideWithDetails(PostCreate, &postCreate, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode post create function details: %v", err)
	}
	if postCreOk {
		res.PostCreateFunction = &postCreate.Function
	}

	pd := PreDeleteFunctionDetails{}
	pdOk, err := overrides.ResourceOverrideWithDetails(PreDelete, &pd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pre delete function details: %v", err)
	}
	if pdOk {
		res.PreDeleteFunction = &pd.Function
	}

	// Resource Override: Customize Diff
	cdiff := CustomizeDiffDetails{}
	cdOk, err := overrides.ResourceOverrideWithDetails(CustomizeDiff, &cdiff, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode customize diff details: %v", err)
	}

	if cdOk {
		res.CustomizeDiff = cdiff.Functions
	}

	// ListFields
	if parameters, ok := typeFetcher.doc.Paths["list"]; ok {
		for _, param := range parameters.Parameters {
			if param.Name != "" {
				res.ListFields = append(res.ListFields, param.Name)
			}
		}
	}

	// Resource Override: No Sweeper
	res.HasSweeper = true
	if overrides.ResourceOverride(NoSweeper, location) {
		res.HasSweeper = false
	}

	stateHint, ok := schema.Extension["x-dcl-uses-state-hint"].(bool)
	if ok {
		res.StateHint = stateHint
	}

	// Resource Override: CustomCreateDirectiveFunction
	createDirectiveFunc := CustomCreateDirectiveDetails{}
	createDirectiveFuncOk, err := overrides.ResourceOverrideWithDetails(CustomCreateDirective, &createDirectiveFunc, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom create directive function details: %v", err)
	}
	if createDirectiveFuncOk {
		res.CustomCreateDirectiveFunction = &createDirectiveFunc.Function
	}

	// Resource Override: SkipDeleteFunction
	skipDeleteFunc := SkipDeleteFunctionDetails{}
	skipDeleteFuncOk, err := overrides.ResourceOverrideWithDetails(SkipDelete, &skipDeleteFunc, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode skip delete function details: %v", err)
	}
	if skipDeleteFuncOk {
		res.SkipDeleteFunction = &skipDeleteFunc.Function
	}

	return &res, nil
}

func readVirtualField(vf VirtualFieldDetails) Property {
	prop := Property{
		title: vf.Name,
		Type:  Type{&openapi.Schema{Type: vf.Type}},
	}

	if vf.Type == "boolean" {
		def := "false"
		prop.Default = &def
	}

	if vf.Output {
		prop.Computed = true
	} else {
		prop.Optional = true
		prop.Settable = true
	}
	return prop
}
