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

	log.Printf("%+v", resourcePaths)
	log.Printf("%+v", productPath)
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

func WriteGoTemplate(templatePath, filePath string, input any) {
	contents := bytes.Buffer{}

	templateFileName := filepath.Base(templatePath)
	templates := []string{
		templatePath,
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
