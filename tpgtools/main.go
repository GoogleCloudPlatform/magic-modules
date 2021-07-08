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
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/golang/glog"

	"github.com/nasa9084/go-openapi"
	"gopkg.in/yaml.v2"
)

var fPath = flag.String("path", "", "path to the root service directory holding openapi schemas")
var tPath = flag.String("overrides", "", "path to the root directory holding overrides files")
var cPath = flag.String("handwritten", "handwritten", "path to the root directory holding handwritten files to copy")
var oPath = flag.String("output", "", "path to output generated files to")

var sFilter = flag.String("service", "", "optional service name. If specified, only this service is generated")
var rFilter = flag.String("resource", "", "optional resource name (from filename). If specified, only resources with this name are generated")
var vFilter = flag.String("version", "", "optional version name. If specified, this version is preferred for resource generation when applicable")

var mode = flag.String("mode", "", "mode for the generator. If unset, creates the provider. Options: 'serialization'")

var terraformResourceDirectory = "google-beta"

func main() {
	resources, products := loadAndModelResources()

	if mode != nil && *mode == "serialization" {
		if vFilter != nil {
			glog.Warning("[WARNING] serialization mode uses all resource versions. version flag is ignored")
		}
		generateSerializationLogic(resources)
		return
	}

	var resourcesForVersion []*Resource
	var productsForVersion []*ProductMetadata
	var version *Version
	if vFilter != nil && *vFilter != "" {
		version = fromString(*vFilter)
		if version == nil {
			glog.Exitf("Failed finding version for input: %s", *vFilter)
		}
		resourcesForVersion = resources[*version]
		productsForVersion = products[*version]
	} else {
		resourcesForVersion = resources[allVersions()[0]]
		productsForVersion = products[allVersions()[0]]
	}
	if *version == GA_VERSION {
		terraformResourceDirectory = "google"
	}

	for _, resource := range resourcesForVersion {
		if skipResource(resource) {
			continue
		}
		resJSON, err := json.MarshalIndent(resource, "", "  ")
		if err != nil {
			glog.Errorf("Failed to marshal resource struct")
		} else {
			glog.Infof("Generating from resource %s", string(resJSON))
		}

		generateResourceFile(resource)
		generateSweeperFile(resource)
		generateResourceWebsiteFile(resource, resources, version)
		generateResourceTestFile(resource, resources, version)
	}

	// product specific generation
	generateProductsFile("provider_dcl_endpoints", productsForVersion)
	generateProductsFile("provider_dcl_client_creation", productsForVersion)

	if oPath == nil || *oPath == "" {
		glog.Info("Skipping copying handwritten files, no output specified")
		return
	}

	if cPath == nil || *cPath == "" {
		glog.Info("No handwritten path specified")
		return
	}

	copyHandwrittenFiles(*cPath, *oPath)
}

func skipResource(r *Resource) bool {
	// if a filter is specified, skip filtered services
	if sFilter != nil && *sFilter != "" && *sFilter != r.ProductMetadata().PackageName {
		return true
	}
	// skip filtered resources
	if rFilter != nil && *rFilter != "" && strings.ToLower(*rFilter) != snakeToLowercase(r.DCLName()) {
		return true
	}

	// skip if set to SerializationOnly
	return r.SerializationOnly
}

// TODO(rileykarson): Change interface to an error, handle exceptional stuff in
// main func.
func loadAndModelResources() (map[Version][]*Resource, map[Version][]*ProductMetadata) {
	flag.Parse()
	if fPath == nil || *fPath == "" {
		glog.Exit("No path specified")
	}

	dirs, err := ioutil.ReadDir(*fPath)
	if err != nil {
		glog.Fatal(err)
	}
	resources := make(map[Version][]*Resource)
	products := make(map[Version][]*ProductMetadata)

	for _, version := range allVersions() {
		resources[version] = make([]*Resource, 0)
		for _, v := range dirs {
			// skip flat files- we're looking for service directory
			if !v.IsDir() {
				continue
			}

			var specs []os.FileInfo
			var packagePath string
			if version == GA_VERSION {
				// GA has no separate directory
				packagePath = v.Name()
			} else {
				packagePath = path.Join(v.Name(), version.V)
			}

			specs, err = ioutil.ReadDir(path.Join(*fPath, packagePath))
			newProduct := getProductInformation(packagePath, specs)
			if newProduct == nil {
				// No resource at this version for this product
				continue
			}
			products[version] = append(products[version], newProduct)

			newResources := getResources(packagePath, specs, version)
			resources[version] = append(resources[version], newResources...)

		}
	}

	return resources, products
}

func getProductInformation(packagePath string, specs []os.FileInfo) *ProductMetadata {
	for _, f := range specs {
		if f.IsDir() {
			continue
		}

		document := loadDocument(packagePath, &f)
		productMetadata := GetProductMetadataFromDocument(document, packagePath)

		return productMetadata
	}
	return nil
}

func getResources(packagePath string, specs []os.FileInfo, version Version) []*Resource {
	var resources []*Resource
	for _, f := range specs {
		if f.IsDir() {
			continue
		}

		document := loadDocument(packagePath, &f)
		productMetadata := GetProductMetadataFromDocument(document, packagePath)
		titleParts := strings.Split(document.Info.Title, "/")

		var schema *openapi.Schema
		for k, v := range document.Components.Schemas {
			if k == titleParts[len(titleParts)-1] {
				schema = v
				schema.Title = k
			}
		}

		overrides := loadOverrides(packagePath, f.Name())

		samples := loadSamples(packagePath, f.Name(), version)

		if schema == nil {
			glog.Exit(fmt.Sprintf("Could not find document schema for %s", document.Info.Title))
		}

		if err := schema.Validate(); err != nil {
			glog.Exit(err)
		}

		lRaw := schema.Extension["x-dcl-locations"]
		var schemaLocations []interface{}
		if lRaw == nil {
			schemaLocations = make([]interface{}, 0)
		} else {
			schemaLocations = lRaw.([]interface{})
		}

		typeFetcher := NewTypeFetcher(document)
		var locations []string
		// If the schema cannot be split into two or more locations, we specify this
		// by passing a single empty location string.
		if len(schemaLocations) < 2 {
			locations = make([]string, 1)
		} else {
			locations = make([]string, 0, len(schemaLocations))
			for _, l := range schemaLocations {
				locations = append(locations, l.(string))
			}
		}

		for _, l := range locations {
			res, err := createResource(schema, typeFetcher, overrides, productMetadata, samples, l)
			if err != nil {
				glog.Exit(err)
			}

			resources = append(resources, res)
		}
	}

	return resources
}

func loadDocument(packagePath string, f *os.FileInfo) *openapi.Document {
	// TODO: use yaml.UnmarshalStrict once apply / list paths are changed to
	// specification extensions and we're using a datatype that supports them.
	document := &openapi.Document{}
	p := path.Join(*fPath, packagePath, (*f).Name())
	b, err := ioutil.ReadFile(p)
	if err != nil {
		glog.Exit(err)
	}
	err = yaml.Unmarshal(b, document)
	if err != nil {
		glog.Exit(err)
	}
	return document
}

// SerializationInput contains an array of resources along with additional generation metadata.
type SerializationInput struct {
	Resources map[Version][]*Resource
	Packages  map[string]string
}

func generateSerializationLogic(specs map[Version][]*Resource) {
	buf := bytes.Buffer{}
	tmpl, err := template.New("serialization.go.tmpl").Funcs(TemplateFunctions).ParseFiles(
		"templates/serialization.go.tmpl",
	)
	if err != nil {
		glog.Exit(err)
	}

	packageMap := make(map[string]string)
	for v, resList := range specs {
		for _, res := range resList {
			var pkgName, pkgPath string
			pkgName = res.Package() + v.SerializationSuffix
			if v == BETA_VERSION {
				pkgPath = path.Join(res.Package(), v.V)
			} else {
				pkgPath = res.Package()
			}

			if _, ok := packageMap[pkgPath]; !ok {
				packageMap[pkgName] = pkgPath
			}
		}
	}

	tmplInput := SerializationInput{
		Resources: specs,
		Packages:  packageMap,
	}

	if err = tmpl.ExecuteTemplate(&buf, "serialization.go.tmpl", tmplInput); err != nil {
		glog.Exit(err)
	}

	formatted, err := formatSource(&buf)
	if err != nil {
		glog.Error(fmt.Errorf("error formatting serialization logic: %v", err))
	}

	if oPath == nil || *oPath == "" {
		fmt.Printf("%v", string(formatted))
	} else {
		err := ioutil.WriteFile(path.Join(*oPath, "serialization.go"), formatted, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}

func loadOverrides(packagePath, fileName string) Overrides {
	overrides := Overrides{}
	if !(tPath == nil) && !(*tPath == "") {
		b, err := ioutil.ReadFile(path.Join(*tPath, packagePath, fileName))
		if err != nil {
			// ignore the error if the file just doesn't exist
			if !os.IsNotExist(err) {
				glog.Exit(err)
			}
		} else {
			err = yaml.UnmarshalStrict(b, &overrides)
			if err != nil {
				glog.Exit(err)
			}
		}
	}
	return overrides
}

func loadSamples(packagePath, fileName string, version Version) Samples {
	samples := Samples{}

	if mode != nil && *mode == "serialization" {
		return samples
	}

	// Samples appear in the root product folder
	packagePath = strings.Split(packagePath, "beta")[0]
	samplesPath := path.Join(*fPath, packagePath, "samples")
	files, err := ioutil.ReadDir(samplesPath)
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
		tc, err := ioutil.ReadFile(path.Join(samplesPath, file.Name()))
		if err != nil {
			glog.Exit(err)
		}

		err = yaml.UnmarshalStrict(tc, &sample)
		if err != nil {
			glog.Exit(err)
		}

		versionMatch := false
		hasGA := false
		for _, v := range sample.Versions {
			if v == version.V {
				versionMatch = true
			}
			if v == "ga" {
				hasGA = true
			}
		}
		if !versionMatch {
			continue
		}

		var dependencies []Dependency
		mainResource := loadSampleDependency(samplesPath, *sample.PrimaryResource, version)
		dependencies = append(dependencies, mainResource)
		for _, dFileName := range sample.DependencyFileNames {
			dependency := loadSampleDependency(samplesPath, dFileName, version)
			dependencies = append(dependencies, dependency)
		}
		sample.DependencyList = dependencies
		sample.TestSlug = sampleNameToTitleCase(*sample.Name)
		sample.HasGAEquivalent = hasGA
		samples = append(samples, sample)
	}

	return samples
}

func loadSampleDependency(samplesPath, fileName string, version Version) Dependency {
	dFileNameParts := strings.Split(fileName, "samples/")
	fileName = dFileNameParts[len(dFileNameParts)-1]
	dependencyBytes, err := ioutil.ReadFile(path.Join(samplesPath, fileName))
	d, err := BuildDependency(fileName, version, dependencyBytes)
	if err != nil {
		glog.Exit(err)
	}
	return *d
}

func generateResourceFile(res *Resource) {
	// Generate resource file
	tmplInput := ResourceInput{
		Resource: *res,
	}

	tmpl, err := template.New("resource.go.tmpl").Funcs(TemplateFunctions).ParseFiles(
		"templates/resource.go.tmpl",
	)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, "resource.go.tmpl", tmplInput); err != nil {
		glog.Exit(err)
	}

	if err != nil {
		glog.Exit(err)
	}

	formatted, err := formatSource(&contents)
	if err != nil {
		glog.Error(fmt.Errorf("error formatting %v: %v - resource \n ", res.ProductName()+res.Name(), err))
	}

	if oPath == nil || *oPath == "" {
		fmt.Printf("%v", string(formatted))
	} else {
		outname := fmt.Sprintf("resource_%s_%s.go", res.ProductName(), res.Name())
		err := ioutil.WriteFile(path.Join(*oPath, terraformResourceDirectory, outname), formatted, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}

func generateSweeperFile(res *Resource) {
	if !res.HasSweeper {
		return
	}

	// Generate resource file
	tmplInput := ResourceInput{
		Resource: *res,
	}

	tmpl, err := template.New("sweeper.go.tmpl").Funcs(TemplateFunctions).ParseFiles(
		"templates/sweeper.go.tmpl",
	)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, "sweeper.go.tmpl", tmplInput); err != nil {
		glog.Exit(err)
	}

	if err != nil {
		glog.Exit(err)
	}

	formatted, err := formatSource(&contents)
	if err != nil {
		glog.Error(fmt.Errorf("error formatting %v: %v - sweeper", res.ProductName()+res.Name(), err))
	}

	if oPath == nil || *oPath == "" {
		fmt.Printf("%v", string(formatted))
	} else {
		outname := fmt.Sprintf("resource_%s_%s_sweeper_test.go", res.ProductName(), res.Name())
		err := ioutil.WriteFile(path.Join(*oPath, terraformResourceDirectory, outname), formatted, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}

func generateProductsFile(fileName string, products []*ProductMetadata) {
	if len(products) <= 0 {
		return
	}
	templateFileName := fileName + ".go.tmpl"
	// Generate endpoints file
	tmpl, err := template.New(templateFileName).Funcs(TemplateFunctions).ParseFiles(
		"templates/" + templateFileName,
	)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, products); err != nil {
		glog.Exit(err)
	}

	if err != nil {
		glog.Exit(err)
	}

	formatted, err := formatSource(&contents)
	if err != nil {
		glog.Error(fmt.Errorf("error formatting package %s file: \n%w", fileName, err))
	}

	if oPath == nil || *oPath == "" {
		fmt.Printf("%v", string(formatted))
	} else {
		outname := fmt.Sprintf(fileName + ".go")
		err := ioutil.WriteFile(path.Join(*oPath, terraformResourceDirectory, outname), formatted, 0644)
		if err != nil {
			glog.Exit(err)
		}
	}
}

var TemplateFunctions = template.FuncMap{
	"title":          strings.Title,
	"patternToRegex": PatternToRegex,
	"replace":        strings.Replace,
}

// TypeFetcher fetches reused types, as marked by the $ref field being marked on an OpenAPI schema.
type TypeFetcher struct {
	doc *openapi.Document

	// Tracks if a property has already been generated.
	generates map[string]string
}

// NewTypeFetcher returns a TypeFetcher for a OpenAPI document.
func NewTypeFetcher(doc *openapi.Document) *TypeFetcher {
	return &TypeFetcher{
		doc:       doc,
		generates: make(map[string]string),
	}
}

// ResolveSchema resolves a #/components/schemas reference from a reused type.
func (r *TypeFetcher) ResolveSchema(ref string) (*openapi.Schema, error) {
	return openapi.ResolveSchema(r.doc, ref)
}

// PackagePathForReference returns either the packageName or the shared package name.
func (r *TypeFetcher) PackagePathForReference(ref, packageName string) string {
	if v, ok := r.generates[ref]; ok {
		return v
	} else {
		r.generates[ref] = packageName
		return packageName
	}
}
