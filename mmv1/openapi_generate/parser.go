// Copyright 2024 Google Inc.
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

// Code generator for a library converting terraform state to gcp objects.

package openapi_generate

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"log"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	r "github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/getkin/kin-openapi/openapi3"
	"gopkg.in/yaml.v2"
)

type Parser struct {
	Folder string
	Output string
}

func NewOpenapiParser(folder, output string) Parser {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf(err.Error())
	}

	parser := Parser{
		Folder: path.Join(wd, folder),
		Output: path.Join(wd, output),
	}

	return parser
}

func (parser Parser) Run() {
	f, err := os.Open(parser.Folder)
	if err != nil {
		log.Fatalf(err.Error())
		return
	}
	defer f.Close()
	files, err := f.Readdirnames(0)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// check if folder is empty
	if len(files) == 0 {
		log.Fatalf("No OpenAPI files found in %s", parser.Folder)
	}

	for _, file := range files {
		parser.WriteYaml(path.Join(parser.Folder, file))
	}
}

func (parser Parser) WriteYaml(filePath string) {
	log.Printf("Reading from file path %s", filePath)

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, _ := loader.LoadFromFile(filePath)
	_ = doc.Validate(ctx)

	header, err := os.ReadFile("openapi_generate/header.txt")
	if err != nil {
		log.Fatalf("error reading header %v", err)
	}

	resourcePaths := findResources(doc)
	productPath := buildProduct(filePath, parser.Output, doc, header)

	// Disables line wrap for long strings
	yaml.FutureLineWrap()
	log.Printf("Generated product %+v/product.yaml", productPath)
	for _, pathArray := range resourcePaths {
		resource := buildResource(filePath, pathArray[0], pathArray[1], doc)

		// marshal method
		resourceOutPathMarshal := filepath.Join(productPath, fmt.Sprintf("%s.yaml", resource.Name))
		bytes, err := yaml.Marshal(resource)
		if err != nil {
			log.Fatalf("error marshalling yaml %v: %v", resourceOutPathMarshal, err)
		}

		f, err := os.Create(resourceOutPathMarshal)
		if err != nil {
			log.Fatalf("error creating resource file %v", err)
		}
		_, err = f.Write(header)
		if err != nil {
			log.Fatalf("error writing resource file header %v", err)
		}
		_, err = f.Write(bytes)
		if err != nil {
			log.Fatalf("error writing resource file %v", err)
		}
		err = f.Close()
		if err != nil {
			log.Fatalf("error closing resource file %v", err)
		}
		log.Printf("Generated resource %s", resourceOutPathMarshal)
	}
}

func findResources(doc *openapi3.T) [][]string {
	var resourcePaths [][]string

	pathMap := doc.Paths.Map()
	for key, pathValue := range pathMap {
		if pathValue.Post == nil {
			continue
		}

		// Not very clever way of identifying create resource methods
		if strings.HasPrefix(pathValue.Post.OperationID, "Create") {
			resourcePath := key
			resourceName := strings.Replace(pathValue.Post.OperationID, "Create", "", 1)
			resourcePaths = append(resourcePaths, []string{resourcePath, resourceName})
		}
	}

	return resourcePaths
}

func buildProduct(filePath, output string, root *openapi3.T, header []byte) string {

	version := root.Info.Version
	server := root.Servers[0].URL

	productName := strings.Split(filepath.Base(filePath), "_")[0]
	productPath := filepath.Join(output, productName)

	if err := os.MkdirAll(productPath, os.ModePerm); err != nil {
		log.Fatalf("error creating product output directory %v: %v", productPath, err)
	}

	apiProduct := &api.Product{}
	apiVersion := &product.Version{}

	apiVersion.BaseUrl = fmt.Sprintf("%s/%s/", server, version)
	// TODO(slevenick) figure out how to tell the API version
	apiVersion.Name = "ga"
	apiProduct.Versions = []*product.Version{apiVersion}

	// Standard titling is "Service Name API"
	displayName := strings.Replace(root.Info.Title, " API", "", 1)
	apiProduct.Name = strings.ReplaceAll(displayName, " ", "")
	apiProduct.DisplayName = displayName

	//Scopes should be added soon to OpenAPI, until then use global scope
	apiProduct.Scopes = []string{"https://www.googleapis.com/auth/cloud-platform"}

	productOutPathMarshal := filepath.Join(output, fmt.Sprintf("/%s/product.yaml", productName))

	// Default yaml marshaller
	bytes, err := yaml.Marshal(apiProduct)
	if err != nil {
		log.Fatalf("error marshalling yaml %v: %v", productOutPathMarshal, err)
	}

	f, err := os.Create(productOutPathMarshal)
	if err != nil {
		log.Fatalf("error creating product file %v", err)
	}
	_, err = f.Write(header)
	if err != nil {
		log.Fatalf("error writing product file header %v", err)
	}
	_, err = f.Write(bytes)
	if err != nil {
		log.Fatalf("error writing product file %v", err)
	}
	err = f.Close()
	if err != nil {
		log.Fatalf("error closing product file %v", err)
	}
	return productPath
}

func baseUrl(resourcePath string) string {
	base := strings.ReplaceAll(resourcePath, "{", "{{")
	base = strings.ReplaceAll(base, "}", "}}")
	// Some APIs use projectsId and locationsId, but we have standardized on these
	base = strings.ReplaceAll(base, "projectsId", "project")
	base = strings.ReplaceAll(base, "locationsId", "location")
	base = stripVersion(base)
	r := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := r.FindStringSubmatch(base)
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		base = strings.ReplaceAll(base, match, google.Underscore(match))
	}
	return base
}

// OpenAPI paths are prefixed with the version of the API, which already exists
// in the product. Strip it out here
func stripVersion(path string) string {
	pattern := `^(/.*v\d[^/]*/)`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(path, "")
}

func buildResource(filePath, resourcePath, resourceName string, root *openapi3.T) api.Resource {
	resource := api.Resource{}

	parsedObjects := parseOpenApi(resourcePath, resourceName, root)

	parameters := parsedObjects[0].([]*api.Type)
	properties := parsedObjects[1].([]*api.Type)
	queryParam := parsedObjects[2].(string)

	baseUrl := baseUrl(resourcePath)
	selfLink := fmt.Sprintf("%s/{{%s}}", baseUrl, google.Underscore(queryParam))

	resource.Name = resourceName
	resource.BaseUrl = baseUrl
	resource.Parameters = parameters
	resource.Properties = properties
	resource.SelfLink = selfLink
	resource.IdFormat = selfLink
	resource.ImportFormat = []string{selfLink}
	resource.CreateUrl = fmt.Sprintf("%s?%s={{%s}}", baseUrl, queryParam, google.Underscore(queryParam))
	resource.Description = "Description"

	resource.AutogenAsync = true
	async := api.NewAsync()
	async.Operation.BaseUrl = "{{op_id}}"
	async.Result.ResourceInsideResponse = true
	resource.Async = async

	example := r.Examples{}
	example.Name = "name_of_example_file"
	example.PrimaryResourceId = "example"
	example.Vars = map[string]string{"resource_name": "test-resource"}

	resource.Examples = []r.Examples{example}

	return resource
}

func parseOpenApi(resourcePath, resourceName string, root *openapi3.T) []any {
	returnArray := []any{}
	path := root.Paths.Find(resourcePath)

	parameters := []*api.Type{}
	var idParam string
	for _, param := range path.Post.Parameters {
		if strings.Contains(strings.ToLower(param.Value.Name), strings.ToLower(resourceName)) {
			idParam = param.Value.Name
		}
		paramObj := writeObject(param.Value.Name, param.Value.Schema, propType(param.Value.Schema), true)
		description := param.Value.Description
		if strings.TrimSpace(description) == "" {
			description = "No description"
		}
		paramObj.Description = trimSpacesFromDescription(description)

		if param.Value.Name == "requestId" || param.Value.Name == "validateOnly" || paramObj.Name == "" {
			continue
		}

		// All parameters are immutable
		paramObj.Immutable = true
		parameters = append(parameters, &paramObj)
	}

	properties := buildProperties(path.Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties, path.Post.RequestBody.Value.Content["application/json"].Schema.Value.Required)

	returnArray = append(returnArray, parameters)
	returnArray = append(returnArray, properties)
	returnArray = append(returnArray, idParam)

	return returnArray
}

func propType(prop *openapi3.SchemaRef) openapi3.Types {
	if len(prop.Value.AllOf) > 0 {
		return *prop.Value.AllOf[0].Value.Type
	} else {
		return *prop.Value.Type
	}
}

func writeObject(name string, obj *openapi3.SchemaRef, objType openapi3.Types, urlParam bool) api.Type {
	var field api.Type

	switch name {
	case "projectsId", "project":
		// projectsId and project are omitted in MMv1 as they are inferred from
		// the presence of {{project}} in the URL
		return field
	case "locationsId":
		name = "location"
	}
	additionalDescription := ""

	if len(obj.Value.AllOf) > 0 {
		obj = obj.Value.AllOf[0]
		objType = *obj.Value.Type
	}

	field.Name = name
	switch objType[0] {
	case "string":
		field.Type = "String"
		if len(obj.Value.Enum) > 0 {
			var enums []string
			for _, enum := range obj.Value.Enum {
				enums = append(enums, fmt.Sprintf("%v", enum))
			}
			additionalDescription = fmt.Sprintf("\n Possible values:\n %s", strings.Join(enums, "\n"))
		}
	case "integer":
		field.Type = "Integer"
	case "number":
		field.Type = "Double"
	case "boolean":
		field.Type = "Boolean"
	case "object":
		if field.Name == "labels" {
			// Standard labels implementation
			field.Type = "KeyValueLabels"
			break
		}

		if obj.Value.AdditionalProperties.Schema != nil && obj.Value.AdditionalProperties.Schema.Value.Type.Is("string") {
			// AdditionalProperties with type string is a string -> string map
			field.Type = "KeyValuePairs"
			break
		}

		field.Type = "NestedObject"

		field.Properties = buildProperties(obj.Value.Properties, obj.Value.Required)
	case "array":
		field.Type = "Array"
		var subField api.Type
		typ := *obj.Value.Items.Value.Type
		switch typ[0] {
		case "string":
			subField.Type = "String"
		case "integer":
			subField.Type = "Integer"
		case "number":
			subField.Type = "Double"
		case "boolean":
			subField.Type = "Boolean"
		case "object":
			subField.Type = "NestedObject"
			subField.Properties = buildProperties(obj.Value.Items.Value.Properties, obj.Value.Items.Value.Required)
		}
		field.ItemType = &subField
	default:
		panic(fmt.Sprintf("Failed to identify field type for %s %s", field.Name, objType[0]))
	}

	description := fmt.Sprintf("%s %s", obj.Value.Description, additionalDescription)
	if strings.TrimSpace(description) == "" {
		description = "No description"
	}

	field.Description = trimSpacesFromDescription(description)

	if urlParam {
		field.UrlParamOnly = true
		field.Required = true
	}

	// These methods are only available when the field is set
	if obj.Value.ReadOnly {
		field.Output = true
	}

	// x-google-identifier fields are described by AIP 203 and are represented
	// as output only in Terraform.
	xGoogleId, err := obj.JSONLookup("x-google-identifier")
	if err == nil && xGoogleId != nil {
		field.Output = true
	}

	xGoogleImmutable, err := obj.JSONLookup("x-google-immutable")
	if err == nil && xGoogleImmutable != nil {
		field.Immutable = true
	}

	return field
}

func buildProperties(props openapi3.Schemas, required []string) []*api.Type {
	properties := []*api.Type{}
	for k, prop := range props {
		propObj := writeObject(k, prop, propType(prop), false)
		if slices.Contains(required, k) {
			propObj.Required = true
		}
		properties = append(properties, &propObj)
	}
	return properties
}

// Trims whitespace from the ends of lines in a description to force multiline
// formatting for strings with newlines present
func trimSpacesFromDescription(description string) string {
	lines := strings.Split(description, "\n")
	var trimmedDescription []string
	for _, line := range lines {
		trimmedDescription = append(trimmedDescription, strings.Trim(line, " "))
	}
	return strings.Join(trimmedDescription, "\n")
}
