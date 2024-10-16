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
	"bytes"
	"context"
	"fmt"
	"os"
	"regexp"
	"path"
	"path/filepath"
	"strings"

	"log"

	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang/glog"
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

	resourcePaths := findResources(doc)
	productPath := buildProduct(filePath, parser.Output, doc)

	log.Printf("Generated product %+v/product.yaml", productPath)
	for _, pathArray := range resourcePaths {
		resource := buildResource(filePath, pathArray[0], pathArray[1], doc)

		// template method
		resourceOutPathTemplate := filepath.Join(productPath, fmt.Sprintf("%s_template.yaml", resource.Name))
		templatePath := "openapi_generate/resource_yaml.tmpl"
		WriteGoTemplate(templatePath, resourceOutPathTemplate, resource)
		log.Printf("Generated resource %s", resourceOutPathTemplate)

		// marshal method
		resourceOutPathMarshal := filepath.Join(productPath, fmt.Sprintf("%s_marshal.yaml", resource.Name))
		bytes, err := yaml.Marshal(resource)
		if err != nil {
			log.Fatalf("error marshalling yaml %v: %v", resourceOutPathMarshal, err)
		}
		err = os.WriteFile(resourceOutPathMarshal, bytes, 0644)
		if err != nil {
			log.Fatalf("error writing product to path %v: %v", resourceOutPathMarshal, err)
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

func buildProduct(filePath, output string, root *openapi3.T) string {

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

	// productOutPath := filepath.Join(output, fmt.Sprintf("/%s/product.yaml", productName))
	templatePath := "openapi_generate/product_yaml.tmpl"

	productOutPathTemplate := filepath.Join(output, fmt.Sprintf("/%s/product_template.yaml", productName))
	WriteGoTemplate(templatePath, productOutPathTemplate, apiProduct)

	productOutPathMarshal := filepath.Join(output, fmt.Sprintf("/%s/product_marshal.yaml", productName))

	// Default yaml marshaller
	bytes, err := yaml.Marshal(apiProduct)
	if err != nil {
		log.Fatalf("error marshalling yaml %v: %v", productOutPathMarshal, err)
	}

	err = os.WriteFile(productOutPathMarshal, bytes, 0644)
	if err != nil {
		log.Fatalf("error writing product to path %v: %v", productOutPathMarshal, err)
	}

	return productPath
}

func baseUrl(resourcePath string) string {
	base := strings.ReplaceAll(resourcePath, "{", "{{")
	base = strings.ReplaceAll(base, "}", "}}")
	base = strings.ReplaceAll(base, "projectsId", "project")
	base = strings.ReplaceAll(base, "locationsId", "location")
	base = strings.ReplaceAll(base, "/v1/", "")
	base = strings.ReplaceAll(base, "/v1alpha/", "project")
	r := regexp.MustCompile(`\{\{(\w+)\}\}`)
	matches := r.FindStringSubmatch(base)
	for i := 0; i < len(matches); i++ {
		match := matches[i]
		base = strings.ReplaceAll(base, match, google.Underscore(match))
	}
	return base
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
	resource.CreateUrl = fmt.Sprintf("%s?%s={{%s}}", baseUrl, queryParam, google.Underscore(queryParam))
	resource.Description = "Description"

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
		paramObj := writeObject(param.Value.Name, param.Value.Schema, *param.Value.Schema.Value.Type, true)
		paramObj.Description = param.Value.Description

		if param.Value.Name == "requestId" || param.Value.Name == "validateOnly" || paramObj.Name == "" {
			continue
		}

		// All parameters are immutable
		paramObj.Immutable = true
		parameters = append(parameters, &paramObj)
	}

	log.Print("properties")
	properties := []*api.Type{}
	log.Print(path.Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties)
	for k, prop := range path.Post.RequestBody.Value.Content["application/json"].Schema.Value.Properties {
		log.Print(prop.Value.Type)
		log.Print(k)
		// TODO handle nested object
		if prop.Value.Type != nil {
			propObj := writeObject(k, prop, *prop.Value.Type, false)
			properties = append(properties, &propObj)
		}
		
	}

	returnArray = append(returnArray, parameters)
	returnArray = append(returnArray, properties)
	returnArray = append(returnArray, idParam)

	return returnArray
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

	// log.Printf("%s %+v", name, obj.Value.AllOf)

	if len(obj.Value.AllOf) > 0 {
		obj = obj.Value.AllOf[0]
		objType = *obj.Value.Type
	}

	field.Name = name
	switch objType[0] {
	case "string":
		field.Type = "string"
		if len(obj.Value.Enum) > 0 {
			var enums []string
			for _, enum := range obj.Value.Enum {
				enums = append(enums, fmt.Sprintf("%v", enum))
			}
			additionalDescription = fmt.Sprintf("\n Possible values:\n %s", strings.Join(enums, "\n"))
		}
	case "object":
		if field.Name == "labels" {
			field.Type = "KeyValueLabels"
			break
		}
		if obj.Value.AdditionalProperties.Schema.Value.Type.Is("string") {
			// AdditionalProperties with type string is a string -> string map
			field.Type = "KeyValuePairs"
			log.Print(obj.Value.AdditionalProperties.Schema.Value.Type.Is("string"))
			break
		}

		field.Type = "NestedObject"
	default:
	}

	description := fmt.Sprintf("%s %s", obj.Value.Description, additionalDescription)
	if strings.TrimSpace(description) == "" {
		description = "No description"
	}
	field.Description = description

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
	if obj.Value.ReadOnly || (err == nil && xGoogleImmutable != nil) {
		field.Immutable = true
	}

	return field
}

func WriteGoTemplate(templatePath, filePath string, input any) {
	contents := bytes.Buffer{}

	templateFileName := filepath.Base(templatePath)
	templates := []string{
		templatePath,
		"openapi_generate/property_yaml.tmpl",
		"openapi_generate/description_yaml.tmpl",
	}

	tmpl, err := template.New(templateFileName).Funcs(google.TemplateFunctions).ParseFiles(templates...)
	if err != nil {
		glog.Exit(fmt.Sprintf("error parsing %s for filepath %s ", templatePath, filePath), err)
	}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, input); err != nil {
		glog.Exit(fmt.Sprintf("error executing %s for filepath %s ", templatePath, filePath), err)
	}

	bytes := contents.Bytes()

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		log.Fatalf("error writing product to path %v: %v", filePath, err)
	}

}
