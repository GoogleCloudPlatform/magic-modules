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
	"unicode"

	"github.com/golang/glog"
	"github.com/nasa9084/go-openapi"
)

// Property is the representation of a TPG resource property in tpgtools.
type Property struct {
	// title is the name of a property.
	title string
	// PackageName is the title-cased shortname of a field as it appears as a
	// property in the DCL. For example, "MachineType".
	PackageName string
	// Type is the type of a property.
	Type

	// Settable indicates that a field is settable in desired specs provided to
	// Apply.
	// Settable fields may be Computed- fields are sometimes Optional + Computed
	// to indicate that they have complex default values in the API.
	Settable bool

	// Only applies to nested object properties.
	// Indicates this property should be excluded and its
	// subproperties should be brought up a level.
	Collapsed bool

	// Elem is the value to insert into the Elem field. If empty, none will be
	// inserted.
	// In most cases, this should be a function call to a schema function. For
	// TypeMaps, this will be a one-liner adding a primitive elem.
	Elem            *string
	ElemIsBasicType bool

	// ConflictsWith is the list of fields that this field conflicts with
	ConflictsWith ConflictsWith

	// ConflictsWith is the list of fields that this field conflicts with
	// in JSONCase. For example, "machineType"
	JSONCaseConflictsWith []string

	// Default is the default for the field.
	Default *string

	// raw schema values
	Required  bool
	Optional  bool
	Computed  bool
	Sensitive bool

	ForceNew    bool
	Description string

	DiffSuppressFunc *string
	ValidateFunc     *string
	SetHashFunc      *string
	MaxItems         *int64
	MinItems         *int64
	ConfigMode       *string

	Removed    *string
	Deprecated *string
	// end raw schema values

	// StateGetter is the line of code to retrieve a field from the `d`
	// ResourceData or (TODO:) from a map[string]interface{}
	StateGetter *string

	// StateSetter is the line of code to set a field in the `d` ResourceData
	// or (TODO:) a map[string]interface{}
	StateSetter *string

	// If this field is a three-state boolean in DCL which is represented as a
	// string in terraform. This is done so that the DCL can distinguish between
	// the field being unset and being set to false.
	EnumBool bool

	// An IdentityGetter is a function to retrieve the value of an "identity" field
	// from state. Identity fields will sometimes allow retrieval from multiple
	// fields or from the user's environment variables.
	// In the most common case, project/region/zone will use special resource-level
	// properties instead of IdentityGetters. However, if they have atypical
	// behaviour, such as sourcing a region from a zone, an IdentityGetter will be
	// used instead.
	IdentityGetter *string

	// Sub-properties of nested objects or arrays with nested objects
	Properties []Property

	// Reference to the parent resource.
	// note: "Properties" will not be available.
	resource *Resource

	// Reference to the parent property of a sub-property. If the property is
	// top-level, this will be unset.
	parent *Property

	// customName is the Terraform-specific name that overrides title.
	customName string

	// Ref is the name of the shared reference type.
	ref string

	// If this property allows forward slashes in its value (only important for
	// properties sent in the URL)
	forwardSlashAllowed bool
}

// An IdentityGetter is a function to retrieve the value of an "identity" field
// from state.
type IdentityGetter struct {
	// If HasError is set to true, the function called by FunctionCall will
	// return (v, error)
	HasError bool
	// Rendered function call to insert into the template. For example,
	// d.Get("name").(string)
	FunctionCall string
}

// Name is the shortname of a field. For example, "machine_type".
func (p Property) Name() string {
	if len(p.customName) > 0 {
		return p.customName
	}
	return p.title
}

// overridePath is the path of a property used in override names. For example,
// "node_config.machine_type".
func (p Property) overridePath() string {
	if p.parent != nil {
		return p.parent.overridePath() + "." + p.title
	}

	return p.title
}

// PackageJSONName is the camel-cased shortname of a field as it appears in the
// DCL's json serialization.  For example, "machineType".
func (p Property) PackageJSONName() string {
	s := p.PackageName
	a := []rune(s)

	if len(a) == 0 {
		return ""
	}
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

// PackagePath is the title-cased path of a type (relative to the resource) for
// use in naming functions. For example, "MachineType" or "NodeConfigPreemptible".
func (p Property) PackagePath() string {
	if p.ref != "" {
		return p.ref
	}
	if p.parent != nil {
		return p.parent.PackagePath() + p.PackageName
	}

	return p.PackageName
}

func (p Property) ObjectType() string {
	parent := p
	// Look up chain to see if we are within a reference
	// types within a reference should not use the parent resource's type
	for {
		if parent.ref != "" {
			return p.PackagePath()
		}
		if parent.parent == nil {
			break
		}
		parent = *parent.parent
	}
	return fmt.Sprintf("%s%s", p.resource.DCLTitle(), p.PackagePath())
}

func (p Property) IsArray() bool {
	return (p.Type.String() == SchemaTypeList || p.Type.String() == SchemaTypeSet) && !p.Type.IsObject()
}

func (t Type) IsSet() bool {
	return t.String() == SchemaTypeSet
}

// ShouldGenerateNestedSchema returns true if an object's nested schema function should be generated.
func (p Property) ShouldGenerateNestedSchema() bool {
	return len(p.Properties) > 0 && !p.Collapsed
}

func (p Property) IsServerGeneratedName() bool {
	return p.StateGetter != nil && !p.Settable
}

// DefaultStateGetter returns the line of code to retrieve a field from the `d`
// ResourceData or (TODO:) from a map[string]interface{}
func (p Property) DefaultStateGetter() string {
	rawGetter := fmt.Sprintf("d.Get(%q)", p.Name())
	return buildGetter(p, rawGetter)
}

func (p Property) ChangeStateGetter() string {
	return buildGetter(p, fmt.Sprintf("oldValue(d.GetChange(%q))", p.Name()))
}

// Builds a Getter for constructing a shallow
// version of the object for destory purposes
func (p Property) StateGetterForDestroyTest() string {
	pullValueFromState := fmt.Sprintf(`rs.Primary.Attributes["%s"]`, p.Name())

	switch p.Type.String() {
	case SchemaTypeBool:
		return fmt.Sprintf(`dcl.Bool(%s == "true")`, pullValueFromState)
	case SchemaTypeString:
		if p.Type.IsEnum() {
			return fmt.Sprintf("%s.%sEnumRef(%s)", p.resource.Package(), p.ObjectType(), pullValueFromState)
		}
		if p.Computed {
			return fmt.Sprintf("dcl.StringOrNil(%s)", pullValueFromState)
		}
		return fmt.Sprintf("dcl.String(%s)", pullValueFromState)
	}

	return ""
}

// Builds a Getter for a property with given raw value
func buildGetter(p Property, rawGetter string) string {
	switch p.Type.String() {
	case SchemaTypeBool:
		return fmt.Sprintf("dcl.Bool(%s.(bool))", rawGetter)
	case SchemaTypeString:
		if p.Type.IsEnum() {
			return fmt.Sprintf("%s.%sEnumRef(%s.(string))", p.resource.Package(), p.ObjectType(), rawGetter)
		}
		if p.EnumBool {
			return fmt.Sprintf("expandEnumBool(%s.(string))", rawGetter)
		}
		if p.Computed {
			return fmt.Sprintf("dcl.StringOrNil(%s.(string))", rawGetter)
		}
		return fmt.Sprintf("dcl.String(%s.(string))", rawGetter)
	case SchemaTypeFloat:
		if p.Computed {
			return fmt.Sprintf("dcl.Float64OrNil(%s.(float64))", rawGetter)
		}
		return fmt.Sprintf("dcl.Float64(%s.(float64))", rawGetter)
	case SchemaTypeInt:
		if p.Computed {
			return fmt.Sprintf("dcl.Int64OrNil(int64(%s.(int)))", rawGetter)
		}
		return fmt.Sprintf("dcl.Int64(int64(%s.(int)))", rawGetter)
	case SchemaTypeMap:
		return fmt.Sprintf("tpgresource.CheckStringMap(%s)", rawGetter)
	case SchemaTypeList, SchemaTypeSet:
		if p.Type.IsEnumArray() {
			return fmt.Sprintf("expand%s%sArray(%s)", p.resource.PathType(), p.PackagePath(), rawGetter)
		}
		if p.Type.typ.Items != nil && p.Type.typ.Items.Type == "string" {
			return fmt.Sprintf("expandStringArray(%s)", rawGetter)
		}

		if p.Type.typ.Items != nil && p.Type.typ.Items.Type == "integer" {
			return fmt.Sprintf("expandIntegerArray(%s)", rawGetter)
		}

		if p.Type.typ.Items != nil && len(p.Properties) > 0 {
			return fmt.Sprintf("expand%s%sArray(%s)", p.resource.PathType(), p.PackagePath(), rawGetter)
		}
	}

	if p.typ.Type == "object" {
		return fmt.Sprintf("expand%s%s(%s)", p.resource.PathType(), p.PackagePath(), rawGetter)
	}

	return "<unknown>"
}

// DefaultStateSetter returns the line of code to set a field in the `d`
// ResourceData or (TODO:) a map[string]interface{}
func (p Property) DefaultStateSetter() string {
	switch p.Type.String() {
	case SchemaTypeBool:
		fallthrough
	case SchemaTypeString:
		fallthrough
	case SchemaTypeInt:
		fallthrough
	case SchemaTypeFloat:
		fallthrough
	case SchemaTypeMap:
		return fmt.Sprintf("d.Set(%q, res.%s)", p.Name(), p.PackageName)
	case SchemaTypeList, SchemaTypeSet:
		if p.typ.Items != nil && ((p.typ.Items.Type == "string" && len(p.typ.Items.Enum) == 0) || p.typ.Items.Type == "integer") {
			return fmt.Sprintf("d.Set(%q, res.%s)", p.Name(), p.PackageName)
		}
		if p.typ.Items != nil && (len(p.Properties) > 0 || len(p.typ.Items.Enum) > 0) {
			return fmt.Sprintf("d.Set(%q, flatten%s%sArray(res.%s))", p.Name(), p.resource.PathType(), p.PackagePath(), p.PackageName)
		}
	}

	if p.typ.Type == "object" {
		return fmt.Sprintf("d.Set(%q, flatten%s%s(res.%s))", p.Name(), p.resource.PathType(), p.PackagePath(), p.PackageName)
	}

	return "<unknown>"
}

// ExpandGetter needs to return a snippet of code that produces the DCL-type
// for the field from a map[string]interface{} named obj that represents the
// parent object in Terraform.
func (p Property) ExpandGetter() string {
	rawGetter := fmt.Sprintf("obj[%q]", p.Name())
	return buildGetter(p, rawGetter)
}

// FlattenGetter needs to return a snippet of code that returns an interface{} which
// can be used in the d.Set() call, given a DCL-type for the parent object named `obj`.
func (p Property) FlattenGetter() string {
	return p.flattenGetterWithParent("obj")
}

func (p Property) flattenGetterWithParent(parent string) string {
	switch p.Type.String() {
	case SchemaTypeBool:
		fallthrough
	case SchemaTypeString:
		fallthrough
	case SchemaTypeInt:
		fallthrough
	case SchemaTypeFloat:
		fallthrough
	case SchemaTypeMap:
		if p.EnumBool {
			return fmt.Sprintf("flattenEnumBool(%s.%s)", parent, p.PackageName)
		}
		return fmt.Sprintf("%s.%s", parent, p.PackageName)
	case SchemaTypeList, SchemaTypeSet:
		if p.Type.IsEnumArray() {
			return fmt.Sprintf("flatten%s%sArray(obj.%s)", p.resource.PathType(), p.PackagePath(), p.PackageName)
		}
		if p.Type.typ.Items != nil && p.Type.typ.Items.Type == "integer" {
			return fmt.Sprintf("%s.%s", parent, p.PackageName)
		}
		if p.Type.typ.Items != nil && p.Type.typ.Items.Type == "string" {
			return fmt.Sprintf("%s.%s", parent, p.PackageName)
		}

		if p.Type.typ.Items != nil && len(p.Properties) > 0 {
			return fmt.Sprintf("flatten%s%sArray(%s.%s)", p.resource.PathType(), p.PackagePath(), parent, p.PackageName)
		}
	}

	if p.typ.Type == "object" {
		return fmt.Sprintf("flatten%s%s(%s.%s)", p.resource.PathType(), p.PackagePath(), parent, p.PackageName)
	}

	return "<unknown>"
}

func getSchemaExtensionMap(v interface{}) map[interface{}]interface{} {
	if v != nil {
		return nil
	}
	ls, ok := v.([]interface{})
	if ok && len(ls) > 0 {
		return ls[0].(map[interface{}]interface{})
	}
	return nil
}

func (p Property) DefaultDiffSuppress() *string {
	switch p.Type.String() {
	case SchemaTypeString:
		// Field is reference to another resource
		if _, ok := p.typ.Extension["x-dcl-references"]; ok {
			dsf := "compareSelfLinkOrResourceName"
			return &dsf
		}
	}
	return nil
}

func (p Property) GetRequiredFileImports() (imports []string) {
	if p.ValidateFunc != nil && strings.Contains(*p.ValidateFunc, "validation.") {
		imports = append(imports, GoPkgTerraformSdkValidation)
	}
	return imports
}

func (p Property) DefaultSetHashFunc() *string {
	switch p.Type.String() {
	case SchemaTypeSet:
		if p.ElemIsBasicType {
			shf := "schema.HashString"
			return &shf
		}
		shf := fmt.Sprintf("schema.HashResource(%s)", *p.Elem)
		return &shf
	}
	glog.Fatalf("Failed to find valid hash func")
	return nil
}

// Objects returns a flatmap of the sub-properties within a Property which are
// objects (eg: have sub-properties themselves).
func (p Property) Objects() (props []Property) {
	// if p.Properties is set, p is an object
	for _, v := range p.Properties {
		if len(v.Properties) != 0 {
			if v.ref == "" {
				props = append(props, v)
				props = append(props, v.Objects()...)
			}
		}
	}

	return props
}

// collapsedProperties returns the input list of properties with nested objects collapsed if needed.
func collapsedProperties(props []Property) (collapsed []Property) {
	for _, v := range props {
		if len(v.Properties) != 0 && v.Collapsed {
			collapsed = append(collapsed, collapsedProperties(v.Properties)...)
		} else {
			collapsed = append(collapsed, v)
		}
	}
	return collapsed
}

// Alias []string so that append etc. still work, but we can attach a rendering
// function to the object
type ConflictsWith []string

func (c ConflictsWith) String() string {
	var quoted []string
	for _, s := range c {
		quoted = append(quoted, fmt.Sprintf("%q", s))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(quoted, ","))
}

// Builds a list of properties from an OpenAPI schema.
func createPropertiesFromSchema(schema *openapi.Schema, typeFetcher *TypeFetcher, overrides Overrides, resource *Resource, parent *Property, location string) (props []Property, err error) {
	identityFields := []string{} // always empty if parent != nil
	if parent == nil {
		identityFields = idParts(resource.ID)
	}

	// Maps PackageJSONName back to property Name
	// for conflict fields
	conflictsMap := make(map[string]string)

	for k, v := range schema.Properties {
		ref := ""
		packageName := ""

		if pName, ok := v.Extension["x-dcl-go-name"].(string); ok {
			packageName = pName
		}

		if v.Ref != "" {
			ref = v.Ref
			v, err = typeFetcher.ResolveSchema(v.Ref)
			if err != nil {
				return nil, err
			}
			ref = typeFetcher.PackagePathForReference(ref, v.Extension["x-dcl-go-type"].(string))
		}

		// Sub-properties are referenced by name, and the explicit title value
		// won't be set initially.
		v.Title = k

		if parent == nil && v.Title == "id" {
			// If top-level field is named `id`, rename to avoid collision with Terraform id
			v.Title = fmt.Sprintf("%s%s", resource.Name(), "Id")
		}

		p := Property{
			title:       jsonToSnakeCase(v.Title).snakecase(),
			Type:        Type{typ: v},
			PackageName: packageName,
			Description: v.Description,
			resource:    resource,
			parent:      parent,
			ref:         ref,
		}

		if overrides.PropertyOverride(Exclude, p, location) {
			continue
		}

		do := CustomDefaultDetails{}
		doOk, err := overrides.PropertyOverrideWithDetails(CustomDefault, p, &do, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom list size details")
		}

		if v.Default != "" || doOk {
			def := v.Default
			if doOk {
				def = do.Default
			}
			d, err := renderDefault(p.Type, def)
			if err != nil {
				return nil, fmt.Errorf("failed to render default: %v", err)
			}
			p.Default = &d
		}

		cn := CustomNameDetails{}
		cnOk, err := overrides.PropertyOverrideWithDetails(CustomName, p, &cn, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom name details: %v", err)
		}
		if cnOk {
			p.customName = cn.Name
		}

		if p.Type.String() == SchemaTypeMap {
			e := "&schema.Schema{Type: schema.TypeString}"
			p.Elem = &e
			p.ElemIsBasicType = true
		}

		if sens, ok := v.Extension["x-dcl-sensitive"].(bool); ok {
			p.Sensitive = sens
		}

		if v, ok := v.Extension["x-dcl-conflicts"].([]interface{}); ok {
			// NOTE: DCL not label x-dcl-conflicts for reused types
			// TODO(shuya): handle nested field when b/213503595 got fixed

			if parent == nil {
				for _, ci := range v {
					p.JSONCaseConflictsWith = append(p.JSONCaseConflictsWith, ci.(string))
				}

				conflictsMap[p.PackageJSONName()] = p.Name()
			}
		}

		// Do this before handling properties so we can check if the parent is readOnly
		isSGP := false
		if sgp, ok := v.Extension["x-dcl-server-generated-parameter"].(bool); ok {
			isSGP = sgp
		}
		if v.ReadOnly || isSGP || (parent != nil && parent.Computed) {
			p.Computed = true

			if stringInSlice(p.Name(), identityFields) {
				sg := p.DefaultStateGetter()
				p.StateGetter = &sg
			}
		}

		// Handle object properties
		if len(v.Properties) > 0 {
			props, err := createPropertiesFromSchema(v, typeFetcher, overrides, resource, &p, location)
			if err != nil {
				return nil, err
			}

			p.Properties = props
			if !p.Computed {
				// Computed fields cannot specify MaxItems
				mi := int64(1)
				p.MaxItems = &mi
			}
			e := fmt.Sprintf("%s%sSchema()", resource.PathType(), p.PackagePath())
			p.Elem = &e
			p.ElemIsBasicType = false
		}

		// Handle array properties
		if v.Items != nil {
			ls := CustomListSizeConstraintDetails{}
			lsOk, err := overrides.PropertyOverrideWithDetails(CustomListSize, p, &ls, location)
			if err != nil {
				return nil, fmt.Errorf("failed to decode custom list size details")
			}
			if lsOk {
				if ls.Max > 0 {
					p.MaxItems = &ls.Max
				}
				if ls.Min > 0 {
					p.MinItems = &ls.Min
				}
			}

			// We end up handling arrays of objects very similarly to nested objects
			// themselves
			if len(v.Items.Properties) > 0 {
				props, err := createPropertiesFromSchema(v.Items, typeFetcher, overrides, resource, &p, location)
				if err != nil {
					return nil, err
				}

				p.Properties = props
				e := fmt.Sprintf("%s%sSchema()", resource.PathType(), p.PackagePath())
				p.Elem = &e
				p.ElemIsBasicType = false
			} else {
				i := Type{typ: v.Items}
				e := fmt.Sprintf("&schema.Schema{Type: schema.%s}", i.String())
				if _, ok := v.Extension["x-dcl-references"]; ok {
					e = fmt.Sprintf("&schema.Schema{Type: schema.%s, DiffSuppressFunc: compareSelfLinkOrResourceName, }", i.String())
				}
				p.Elem = &e
				p.ElemIsBasicType = true
			}
		}

		if !p.Computed {
			if stringInSlice(v.Title, schema.Required) {
				p.Required = true
			} else {
				p.Optional = true
			}
		}
		cr := CustomSchemaValuesDetails{}
		crOk, err := overrides.PropertyOverrideWithDetails(CustomSchemaValues, p, &cr, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom required details")
		}
		if crOk {
			p.Required = cr.Required
			p.Optional = cr.Optional
			p.Computed = cr.Computed
		}

		// Handle settable fields. If the field is computed it's not settable but
		// if it's also optional (O+C), it is.
		if !p.Computed || (p.Optional) {
			p.Settable = true

			// NOTE: x-kubernetes-immmutable implies that all children of a field
			// are actually immutable. However, in practice, DCL specs will label
			// every immutable subfield.
			if isImmutable, ok := v.Extension["x-kubernetes-immutable"].(bool); ok && isImmutable {
				p.ForceNew = true
			}

			serverDefault, _ := v.Extension["x-dcl-server-default"].(bool)
			extractIfEmpty, _ := v.Extension["x-dcl-extract-if-empty"].(bool)

			if serverDefault || extractIfEmpty {
				p.Computed = true
			}

			if forwardSlashAllowed, ok := v.Extension["x-dcl-forward-slash-allowed"].(bool); ok && forwardSlashAllowed {
				p.forwardSlashAllowed = true
			}

			// special handling for project/region/zone/other fields with
			// provider defaults
			if stringInSlice(p.title, []string{"project", "region", "zone"}) || stringInSlice(p.customName, []string{"region", "project", "zone"}) {
				p.Optional = true
				p.Required = false
				p.Computed = true

				sg := fmt.Sprintf("dcl.String(%v)", p.Name())
				p.StateGetter = &sg

				cig := &CustomIdentityGetterDetails{}
				cigOk, err := overrides.PropertyOverrideWithDetails(CustomIdentityGetter, p, cig, location)
				if err != nil {
					return nil, fmt.Errorf("failed to decode custom identity getter details")
				}

				propertyName := p.title
				if p.customName != "" {
					propertyName = p.customName
				}
				ig := fmt.Sprintf("get%s(d, config)", renderSnakeAsTitle(miscellaneousNameSnakeCase(propertyName)))
				if cigOk {
					ig = fmt.Sprintf("%s(d, config)", cig.Function)
				}

				p.IdentityGetter = &ig
			} else {
				sg := p.DefaultStateGetter()
				p.StateGetter = &sg
			}
		}

		ss := p.DefaultStateSetter()
		p.StateSetter = &ss

		if p.Sensitive && p.Settable {
			p.StateSetter = nil
		}

		css := CustomStateSetterDetails{}
		cssOk, err := overrides.PropertyOverrideWithDetails(CustomStateSetter, p, &css, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom stateSetter func: %v", err)
		}

		if cssOk {
			p.StateSetter = &css.Function
		}

		irOk := overrides.PropertyOverride(IgnoreRead, p, location)
		if irOk {
			p.StateSetter = nil
		}

		cd := CustomDescriptionDetails{}
		cdOk, err := overrides.PropertyOverrideWithDetails(CustomDescription, p, &cd, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom description details: %v", err)
		}

		if cdOk {
			p.Description = cd.Description
		}

		dsf := CustomDiffSuppressFuncDetails{}
		dsfOk, err := overrides.PropertyOverrideWithDetails(DiffSuppressFunc, p, &dsf, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom diff suppress func: %v", err)
		}

		if dsfOk {
			p.DiffSuppressFunc = &dsf.DiffSuppressFunc
		} else if !(p.Computed && !p.Optional) {
			p.DiffSuppressFunc = p.DefaultDiffSuppress()
		}

		vf := CustomValidationDetails{}
		vfOk, err := overrides.PropertyOverrideWithDetails(CustomValidation, p, &vf, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom validation func: %v", err)
		}

		if vfOk {
			p.ValidateFunc = &vf.Function
		}

		if p.Type.String() == SchemaTypeSet {
			shf := SetHashFuncDetails{}
			shfOk, err := overrides.PropertyOverrideWithDetails(SetHashFunc, p, &shf, location)
			if err != nil {
				return nil, fmt.Errorf("failed to decode set hash func: %v", err)
			}

			if shfOk {
				p.SetHashFunc = &shf.Function
			} else {
				p.SetHashFunc = p.DefaultSetHashFunc()
			}
		}

		cm := CustomConfigModeDetails{}
		cmOk, err := overrides.PropertyOverrideWithDetails(CustomConfigMode, p, &cm, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom config mode func: %v", err)
		}
		if cmOk {
			p.ConfigMode = &cm.Mode
		}

		rd := &RemovedDetails{}
		rdOk, err := overrides.PropertyOverrideWithDetails(Removed, p, rd, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode removed details")
		}
		if rdOk {
			p.Removed = &rd.Message
		}

		dd := &DeprecatedDetails{}
		ddOk, err := overrides.PropertyOverrideWithDetails(Deprecated, p, dd, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode deprecated details")
		}
		if ddOk {
			p.Deprecated = &dd.Message
		}

		if overrides.PropertyOverride(CollapsedObject, p, location) {
			p.Collapsed = true
			if p.parent == nil {
				collapseSS := fmt.Sprintf("setStateForCollapsedObject(d, %s)", p.flattenGetterWithParent("res"))
				p.StateSetter = &collapseSS
			}
			collapseSG := fmt.Sprintf("expand%s%sCollapsed(d)", p.resource.PathType(), p.PackagePath())
			p.StateGetter = &collapseSG
		}

		// Add any new imports as needed
		if ls := p.GetRequiredFileImports(); len(ls) > 0 {
			resource.additionalFileImportSet.Add(ls...)
		}

		csgd := CustomStateGetterDetails{}
		csgdOk, err := overrides.PropertyOverrideWithDetails(CustomStateGetter, p, &csgd, location)
		if err != nil {
			return nil, fmt.Errorf("failed to decode custom state getter details with err %v", err)
		}
		if csgdOk {
			p.StateGetter = &csgd.Function
		}

		if overrides.PropertyOverride(EnumBool, p, location) {
			p.EnumBool = true
			p.Type.typ.Type = "string"
			var parent string
			if p.parent == nil {
				parent = "res"
			} else {
				parent = "obj"
			}
			enumBoolSS := fmt.Sprintf("d.Set(%q, flattenEnumBool(%s.%s))", p.Name(), parent, p.PackageName)
			p.StateSetter = &enumBoolSS
			enumBoolSG := fmt.Sprintf("expandEnumBool(d.Get(%q))", p.Name())
			p.StateGetter = &enumBoolSG
		}

		if overrides.PropertyOverride(GenerateIfNotSet, p, location) {
			p.Computed = true
			p.Required = false
			p.Optional = true
			ig := fmt.Sprintf("generateIfNotSet(d, %q, %q)", p.Name(), "tfgen")
			p.IdentityGetter = &ig
			n := fmt.Sprintf("&%s", p.Name())
			p.StateGetter = &n
		}

		if overrides.PropertyOverride(NamePrefix, p, location) {
			p.Computed = true
			p.Required = false
			p.Optional = true
			ig := fmt.Sprintf("generateIfNotSet(d, %q, d.Get(%q).(string))", p.Name(), "name_prefix")
			p.IdentityGetter = &ig
			n := fmt.Sprintf("&%s", p.Name())
			p.StateGetter = &n

			// plus, add the "name_prefix" property.
			props = append(props, Property{
				title:    "name_prefix",
				Type:     p.Type,
				resource: resource,
				parent:   parent,
				Optional: true,
				Computed: true,
				ForceNew: true,
			})
		}

		if p.ref != "" {
			resource.ReusedTypes = resource.RegisterReusedType(p)
		}

		props = append(props, p)
	}

	// handle conflict fields
	for i, _ := range props {
		p := &props[i]
		if p.JSONCaseConflictsWith != nil {
			for _, cf := range p.JSONCaseConflictsWith {
				if val, ok := conflictsMap[cf]; ok {
					p.ConflictsWith = append(p.ConflictsWith, val)
				} else {
					return nil, fmt.Errorf("Error generating conflict fields. %s is not labeled as a conflict field in DCL", cf)
				}

			}
		}
	}

	// sort the properties so they're in a nice order
	sort.SliceStable(props, propComparator(props))

	return props, nil
}
