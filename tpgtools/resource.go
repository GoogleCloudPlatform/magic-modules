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
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/golang/glog"
	"github.com/nasa9084/go-openapi"
	"gopkg.in/yaml.v2"
)

const GoPkgTerraformSdkValidation = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

// Resource is tpgtools' model of what a information is necessary to generate a
// resource in TPG.
type Resource struct {
	productMetadata *ProductMetadata

	// ID is the Terraform resource id format as a pattern string. Additionally,
	// import formats can be derived from it.
	ID string

	// If the Terraform ID format should be used instead of the DCL ID function.
	// For example, resources with a regional/global cannot use the DCL ID formatter.
	UseTerraformID bool

	// ImportFormats are pattern format strings for importing the Terraform resource.
	// TODO: if none are set, the resource does not support import.
	ImportFormats []string

	// AppendToBasePath is a string that will be appended to the end of the API base path.
	// rarely needed in cases where the shared mm basepath does not include the version
	// as in Montioring https://git.io/Jz4Wn
	AppendToBasePath string

	// ReplaceInBasePath contains a string replacement for the config base path,
	// replacing one substring with another.
	ReplaceInBasePath BasePathReplacement

	// SkipInProvider is true when the resource shouldn't be included in the dclResources
	// map for the provider. This is usually because it was already included through mmv1.
	SkipInProvider bool

	// title is the name of the resource in snake_case. For example,
	// "instance", "backend_service".
	title SnakeCaseTerraformResourceName

	// dclTitle is the name of the resource in TitleCase, as parsed from the relevant
	// DCL YAML file. For example, "Instance", "BackendService".
	// This is particularly useful for acronymizations that exist in
	// resource names, like OSPolicy.  We store it separately from title
	// because it can differ, especially in the case of split resources:
	// "region_backend_service" vs "BackendService".  We use this to
	// create the names of DCL functions - "Apply{{dclTitle}}".
	dclTitle TitleCaseResourceName
	// dclStructName is the name of the resource struct in the DCL.  In almost all cases
	// this is == to dclTitle, but sometimes, due to (for instance) reserved words in the
	// DCL conflicting with resource names, this differs.  We use this to create DCL
	// structs: `foo := compute.{{dclStructName}}{field1: "bar"}`.
	dclStructName TitleCaseResourceName

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

	// location is one of "zone", "region", or "global".
	location string

	// HasProject tells us if the resource has a project field.
	HasProject bool

	// HasCreate tells us if the resource has a create endpoint.
	HasCreate bool

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

	// SerializationOnly defines if this resource should not generate provider files
	// and only be used for serialization
	SerializationOnly bool

	// CustomSerializer defines the function this resource should use to serialize itself.
	CustomSerializer *string

	// TerraformProductName is the Product name overriden from the DCL
	TerraformProductName *SnakeCaseProductName

	// The array of Samples associated with the resource
	Samples []Sample

	// Versions specific information about this resource
	versionMetadata Version

	// Reference points to the rest API
	Reference *Link
	// Guides point to non-rest useful context for the resource.
	Guides []Link
}

type Link struct {
	text string
	url  string
}

type BasePathReplacement struct {
	Present bool
	Old     string
	New     string
}

func (l Link) Markdown() string {
	return fmt.Sprintf("[%s](%s)", l.text, l.url)
}

func (r *Resource) fillLinksFromExtensionsMap(m map[string]interface{}) {
	ref, ok := m["x-dcl-ref"].(map[string]interface{})
	if ok {
		r.Reference = &Link{url: ref["url"].(string), text: ref["text"].(string)}
	}
	gs, ok := m["x-dcl-guides"].([]interface{})
	if ok {
		for _, g := range gs {
			guide := g.(map[interface{}]interface{})
			r.Guides = append(r.Guides, Link{url: guide["url"].(string), text: guide["text"].(string)})
		}
	}
}

// Name is the shortname of a resource. For example, "instance".
func (r Resource) Name() SnakeCaseTerraformResourceName {
	return r.title
}
func (r Resource) TitleCaseFullName() TitleCaseFullName {
	return TitleCaseFullName(snakeToTitleCase(concatenateSnakeCase(r.ProductName(), r.Name())).titlecase())
}

func (r Resource) DCLTitle() TitleCaseResourceName {
	return r.dclTitle
}

func (r Resource) DCLStructName() TitleCaseResourceName {
	if r.dclStructName != "" {
		return r.dclStructName
	}
	return r.dclTitle
}

// TerraformName is the Terraform resource type used in HCL configuration.
// For example, "google_compute_instance"
func (r Resource) TerraformName() SnakeCaseFullName {
	googlePrefix := miscellaneousNameSnakeCase("google")
	if r.TerraformProductName != nil {
		if r.TerraformProductName.snakecase() == "" {
			return SnakeCaseFullName(concatenateSnakeCase(googlePrefix, r.Name()))
		}
		return SnakeCaseFullName(concatenateSnakeCase(googlePrefix, *r.TerraformProductName, r.Name()))
	}
	return SnakeCaseFullName(concatenateSnakeCase(googlePrefix, r.ProductName(), r.Name()))
}

// PathType is the title-cased name of a resource preceded by its package,
// often used to namespace functions. For example, "RedisInstance".
func (r Resource) PathType() TitleCaseFullName {
	return r.TitleCaseFullName()
}

// Package is the namespace of the package within the dcl
// the Package is normally a lowercase variant of ProductName
func (r Resource) Package() DCLPackageName {
	return r.productMetadata.PackageName
}

func (r Resource) TitleCasePackageName() RenderedString {
	return RenderedString(snakeToTitleCase(r.ProductName()).titlecase())
}

func (r Resource) DocsSection() miscellaneousNameTitleCase {
	return r.productMetadata.DocsSection()
}

// ProductName is the snakecase product name of a resource. For example,
// "network_services".
func (r Resource) ProductName() SnakeCaseProductName {
	return r.productMetadata.ProductName
}

func (r Resource) ProductMetadata() *ProductMetadata {
	copy := *r.productMetadata
	return &copy
}

// DCLPackage is the package name of the DCL client library to use for this
// resource. For example, the Package "access_context_manager" at version GA would have a
// DCLPackage of "accesscontextmanager", and at beta would be "accesscontextmanager/beta".
func (r Resource) DCLPackage() DCLPackageNameWithVersion {
	return r.productMetadata.PackageNameWithVersion()
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
func (r Resource) IDFunction() string {
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

func createResource(schema *openapi.Schema, info *openapi.Info, typeFetcher *TypeFetcher, overrides Overrides, product *ProductMetadata, version Version, location string) (*Resource, error) {
	resourceTitle := strings.Split(info.Title, "/")[1]

	res := Resource{
		title:                SnakeCaseTerraformResourceName(jsonToSnakeCase(resourceTitle).snakecase()),
		dclStructName:        TitleCaseResourceName(schema.Title),
		dclTitle:             TitleCaseResourceName(resourceTitle),
		productMetadata:      product,
		versionMetadata:      version,
		Description:          info.Description,
		location:             location,
		InsertTimeoutMinutes: 20,
		UpdateTimeoutMinutes: 20,
		DeleteTimeoutMinutes: 20,
	}

	// Since the resource's "info" extension field can't be accessed, the relevant
	// extensions have been copied into the schema objects.
	res.fillLinksFromExtensionsMap(schema.Extension)

	// Resource Override: Custom Timeout
	ctd := CustomTimeoutDetails{}
	ctdOk, err := overrides.ResourceOverrideWithDetails(CustomTimeout, &ctd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom timeout details: %v", err)
	}
	if ctdOk {
		res.InsertTimeoutMinutes = ctd.TimeoutMinutes
		res.UpdateTimeoutMinutes = ctd.TimeoutMinutes
		res.DeleteTimeoutMinutes = ctd.TimeoutMinutes
	}

	if overrides.ResourceOverride(SkipInProvider, location) {
		res.SkipInProvider = true
	}

	crname := CustomResourceNameDetails{}
	crnameOk, err := overrides.ResourceOverrideWithDetails(CustomResourceName, &crname, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom resource name details: %v", err)
	}

	if crnameOk {
		res.title = SnakeCaseTerraformResourceName(crname.Title)
	}

	id, customID, err := findResourceID(schema, overrides, location)
	if err != nil {
		return nil, err
	}
	res.ID = id
	res.UseTerraformID = customID

	// Resource Override: Custom Import Function
	cifd := CustomImportFunctionDetails{}
	cifdOk, err := overrides.ResourceOverrideWithDetails(CustomImport, &cifd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom import function details: %v", err)
	}
	if cifdOk {
		res.CustomImportFunction = &cifd.Function
	}

	// Resource Override: Append to Base Path
	atbpd := AppendToBasePathDetails{}
	atbpOk, err := overrides.ResourceOverrideWithDetails(AppendToBasePath, &atbpd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode append to base path details: %v", err)
	}
	if atbpOk {
		res.AppendToBasePath = atbpd.String
	}

	// Resource Override: Replace in Base Path
	ribpd := ReplaceInBasePathDetails{}
	ribpOk, err := overrides.ResourceOverrideWithDetails(ReplaceInBasePath, &ribpd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode replace in base path details: %v", err)
	}
	if ribpOk {
		res.ReplaceInBasePath.Present = true
		res.ReplaceInBasePath.Old = ribpd.Old
		res.ReplaceInBasePath.New = ribpd.New
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

	onlyLongFormFormat := shouldAllowForwardSlashInFormat(res.ID, res.Properties)
	// Resource Override: Import formats
	ifd := ImportFormatDetails{}
	ifdOk, err := overrides.ResourceOverrideWithDetails(ImportFormat, &ifd, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode import format details: %v", err)
	}
	if ifdOk {
		res.ImportFormats = ifd.Formats
	} else {
		res.ImportFormats = defaultImportFormats(res.ID, onlyLongFormFormat)
	}

	_, res.HasProject = schema.Properties["project"]

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

	// Determine if a resource has a create method.
	res.HasCreate, _ = schema.Extension["x-dcl-has-create"].(bool)

	// Determine if a resource can use a generated sweeper or not
	// We only supply a certain set of parent values to sweepers, so only generate
	// one if it will actually work- resources with resource parents are not
	// sweepable, in particular, such as nested resources or fine-grained
	// resources. Additional special cases can be handled with overrides.
	res.HasSweeper = true
	validSweeperParameters := []string{"project", "region", "location", "zone", "billingAccount"}
	if deleteAllInfo, ok := typeFetcher.doc.Paths["deleteAll"]; ok {
		for _, p := range deleteAllInfo.Parameters {
			// if any field isn't a standard sweeper parameter, don't make a sweeper
			if !stringInSlice(p.Name, validSweeperParameters) {
				res.HasSweeper = false
			}
		}
	} else {
		// if deleteAll wasn't found, the DCL hasn't published a sweeper
		res.HasSweeper = false
	}

	if overrides.ResourceOverride(NoSweeper, location) {
		if res.HasSweeper == false {
			return nil, fmt.Errorf("superfluous NO_SWEEPER specified for %q", res.TerraformName())
		}

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
	skipDeleteFuncOk, err := overrides.ResourceOverrideWithDetails(SkipDeleteFunction, &skipDeleteFunc, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode skip delete details: %v", err)
	}
	if skipDeleteFuncOk {
		res.SkipDeleteFunction = &skipDeleteFunc.Function
	}

	// Resource Override: SerializationOnly
	res.SerializationOnly = overrides.ResourceOverride(SerializationOnly, location)

	// Resource Override: CustomSerializer
	customSerializerFunc := CustomSerializerDetails{}
	customSerializerFuncOk, err := overrides.ResourceOverrideWithDetails(CustomSerializer, &customSerializerFunc, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode custom serializer function details: %v", err)
	}
	if customSerializerFuncOk {
		res.CustomSerializer = &customSerializerFunc.Function
	}

	// Resource Override: TerraformProductName
	terraformProductName := TerraformProductNameDetails{}
	terraformProductNameOk, err := overrides.ResourceOverrideWithDetails(TerraformProductName, &terraformProductName, location)
	if err != nil {
		return nil, fmt.Errorf("failed to decode terraform product name function details: %v", err)
	}
	if terraformProductNameOk {
		scpn := SnakeCaseProductName(terraformProductName.Product)
		res.TerraformProductName = &scpn
	}

	res.Samples = res.loadSamples()

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

func (r Resource) TestSamples() []Sample {
	return r.getSamples(false)
}

func (r Resource) DocSamples() []Sample {
	return r.getSamples(true)
}

func (r Resource) getSamples(docs bool) []Sample {
	out := []Sample{}
	if len(r.Samples) < 1 {
		return out
	}
	var hideList []string
	if docs {
		hideList = r.Samples[0].DocHide
		if len(r.Samples[0].DocHideConditional) > 0 {
			for _, dochidec := range r.Samples[0].DocHideConditional {
				if r.location == dochidec.Location {
					hideList = append (hideList, dochidec.Name)
				}
			}
		}
	} else {
		hideList = r.Samples[0].Testhide
                if len(r.Samples[0].TestHideConditional) > 0 {
                        for _, testhidec := range r.Samples[0].TestHideConditional {
                                if r.location == testhidec.Location {
                                        hideList = append (hideList, testhidec.Name)
                                }
                        }
                }
	}
	for _, sample := range r.Samples {
		shouldhide := false
		for _, hideName := range hideList {
			if sample.FileName == hideName {
				shouldhide = true
			}
		}
		if !shouldhide {
			out = append(out, sample)
		}
	}

	return out
}

func (r *Resource) getSampleAccessoryFolder() Filepath {
	resourceType := strings.ToLower(r.DCLTitle().titlecase())
	packageName := r.productMetadata.PackageName.lowercase()
	sampleAccessoryFolder := path.Join(*tPath, packageName, "samples", resourceType)
	return Filepath(sampleAccessoryFolder)
}

func (r *Resource) loadSamples() []Sample {
	samples := []Sample{}
	handWritten := r.loadHandWrittenSamples()
	dclSamples := r.loadDCLSamples()
	samples = append(samples, dclSamples...)
	samples = append(samples, handWritten...)
	return samples
}

func (r *Resource) loadHandWrittenSamples() []Sample {
	sampleAccessoryFolder := r.getSampleAccessoryFolder()
	sampleFriendlyMetaPath := path.Join(string(sampleAccessoryFolder), "meta.yaml")
	samples := []Sample{}

	if !pathExists(sampleFriendlyMetaPath) {
		return samples
	}

	files, err := ioutil.ReadDir(string(sampleAccessoryFolder))
	if err != nil {
		glog.Exit(err)
	}

	for _, file := range files {
		if fileName := strings.ToLower(file.Name()); !strings.HasSuffix(fileName, ".tf.tmpl") ||
			strings.Contains(fileName, "_update") {
			continue
		}
		sample := Sample{}
		sampleName := strings.Split(file.Name(), ".")[0]
		sampleDefinitionFile := path.Join(string(sampleAccessoryFolder), sampleName+".yaml")
		var tc []byte
		if pathExists(sampleDefinitionFile) {
			tc, err = mergeYaml(sampleDefinitionFile, sampleFriendlyMetaPath)
		} else {
			tc, err = ioutil.ReadFile(sampleFriendlyMetaPath)
		}
		if err != nil {
			glog.Exit(err)
		}
		err = yaml.UnmarshalStrict(tc, &sample)
		if err != nil {
			glog.Exit(err)
		}

		versionMatch := false

		// if no versions provided assume all versions
		if len(sample.Versions) <= 0 {
			sample.HasGAEquivalent = true
			versionMatch = true
		} else {
			for _, v := range sample.Versions {
				if v == r.versionMetadata.V {
					versionMatch = true
				}
				if v == "ga" {
					sample.HasGAEquivalent = true
				}
			}
		}

		if !versionMatch {
			glog.Errorf("skipping %q due to no version match", file.Name())
			continue
		}

		sample.SamplesPath = sampleAccessoryFolder
		sample.resourceReference = r
		sample.FileName = file.Name()
		sample.PrimaryResource = &(sample.FileName)
		if sample.Name == nil || *sample.Name == "" {
			sampleName = strings.Split(sample.FileName, ".")[0]
			sample.Name = &sampleName
		}
		sample.TestSlug = RenderedString(snakeToTitleCase(miscellaneousNameSnakeCase(sampleName)).titlecase() + "HandWritten")
		samples = append(samples, sample)
	}

	return samples
}

func (r *Resource) loadDCLSamples() []Sample {
	sampleAccessoryFolder := r.getSampleAccessoryFolder()
	packagePath := r.productMetadata.PackagePath
	version := r.versionMetadata.V
        resourceType := r.DCLTitle()
	sampleFriendlyMetaPath := path.Join(string(sampleAccessoryFolder), "meta.yaml")
	samples := []Sample{}

	if mode != nil && *mode == "serialization" {
		return samples
	}

	// Samples appear in the root product folder
	packagePath = Filepath(strings.Split(string(packagePath), "/")[0])
	samplesPath := Filepath(path.Join(*fPath, string(packagePath), "samples"))
	files, err := ioutil.ReadDir(string(samplesPath))
	if err != nil {
		// ignore the error if the file just doesn't exist
		if !os.IsNotExist(err) {
			glog.Exit(err)
		}
	}
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		sample := Sample{}
		sampleOGFilePath := path.Join(string(samplesPath), file.Name())
		var tc []byte
		if pathExists(sampleFriendlyMetaPath) {
			tc, err = mergeYaml(sampleOGFilePath, sampleFriendlyMetaPath)
		} else {
			glog.Errorf("warning : sample meta does not exist for %v at %q", r.TerraformName(), sampleFriendlyMetaPath)
			tc, err = ioutil.ReadFile(path.Join(string(samplesPath), file.Name()))
		}
		if err != nil {
			glog.Exit(err)
		}

		err = yaml.UnmarshalStrict(tc, &sample)
		if err != nil {
			glog.Exit(err)
		}

		versionMatch := false
		for _, v := range sample.Versions {
			if v == version {
				versionMatch = true
			}
			if v == "ga" {
				sample.HasGAEquivalent = true
				versionMatch = true
			}
		}

		primaryResource := *sample.PrimaryResource
		var parts []miscellaneousNameSnakeCase
		parts = assertSnakeArray(strings.Split(primaryResource, "."))
		primaryResourceName := snakeToTitleCase(parts[len(parts)-2])

		if !versionMatch {
			continue
		} else if !strings.EqualFold(primaryResourceName.titlecase(), resourceType.titlecase()) {
			// This scenario will occur for product folders with multiple resources
			continue
		}

		sample.SamplesPath = samplesPath
		sample.resourceReference = r
		sample.FileName = file.Name()

		var dependencies []Dependency
		mainResource := sample.generateSampleDependencyWithName(primaryResource, "primary")
		dependencies = append(dependencies, mainResource)
		for _, dFileName := range sample.DependencyFileNames {
			dependency := sample.generateSampleDependency(dFileName)
			dependencies = append(dependencies, dependency)
		}
		sample.DependencyList = dependencies
		sample.TestSlug = RenderedString(sampleNameToTitleCase(*sample.Name).titlecase())
		samples = append(samples, sample)
	}

	return samples
}
