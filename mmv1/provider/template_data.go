// Copyright 2024 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
)

type TemplateData struct {
	//     include Compile::Core

	OutputFolder string
	Version      product.Version

	TerraformResourceDirectory string
	TerraformProviderModule    string

	// TODO Q2: is this needed?
	//     # Information about the local environment
	//     # (which formatters are enabled, start-time)
	//     attr_accessor :env
}

// Build a map(map[string]interface{}) from a list of paramerter
// The format of passed in parmeters are key1, value1, key2, value2 ...
func wrapMultipleParams(params ...interface{}) (map[string]interface{}, error) {
	if len(params)%2 != 0 {
		return nil, errors.New("invalid number of arguments")
	}
	m := make(map[string]interface{}, len(params)/2)
	for i := 0; i < len(params); i += 2 {
		key, ok := params[i].(string)
		if !ok {
			return nil, errors.New("keys must be strings")
		}
		m[key] = params[i+1]
	}
	return m, nil
}

var TemplateFunctions = template.FuncMap{
	"title":           google.SpaceSeparatedTitle,
	"replace":         strings.Replace,
	"camelize":        google.Camelize,
	"underscore":      google.Underscore,
	"plural":          google.Plural,
	"contains":        strings.Contains,
	"join":            strings.Join,
	"lower":           strings.ToLower,
	"upper":           strings.ToUpper,
	"dict":            wrapMultipleParams,
	"format2regex":    google.Format2Regex,
	"orderProperties": api.OrderProperties,
	"hasPrefix":       strings.HasPrefix,
}

var GA_VERSION = "ga"
var BETA_VERSION = "beta"
var ALPHA_VERSION = "alpha"

func NewTemplateData(outputFolder string, version product.Version) *TemplateData {
	td := TemplateData{OutputFolder: outputFolder, Version: version}

	if version.Name == GA_VERSION {
		td.TerraformResourceDirectory = "google"
		td.TerraformProviderModule = "github.com/hashicorp/terraform-provider-google"
	} else if version.Name == ALPHA_VERSION {
		td.TerraformResourceDirectory = "google-private"
		td.TerraformProviderModule = "internal/terraform-next"
	} else {
		td.TerraformResourceDirectory = "google-beta"
		td.TerraformProviderModule = "github.com/hashicorp/terraform-provider-google-beta"
	}

	return &td
}

func (td *TemplateData) GenerateResourceFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/resource.go.tmpl"
	templates := []string{
		templatePath,
		"templates/terraform/schema_property.go.tmpl",
		"templates/terraform/schema_subresource.go.tmpl",
		"templates/terraform/expand_resource_ref.tmpl",
		"templates/terraform/custom_flatten/go/bigquery_table_ref.go.tmpl",
		"templates/terraform/flatten_property_method.go.tmpl",
		"templates/terraform/expand_property_method.go.tmpl",
	}
	td.GenerateFile(filePath, templatePath, resource, true, templates...)
}

func (td *TemplateData) GenerateDocumentationFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/resource.html.markdown.tmpl"
	templates := []string{
		templatePath,
		"templates/terraform/property_documentation.html.markdown.tmpl",
		"templates/terraform/nested_property_documentation.html.markdown.tmpl",
	}
	td.GenerateFile(filePath, templatePath, resource, false, templates...)
}

func (td *TemplateData) GenerateTestFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/examples/base_configs/test_file.go.tmpl"
	templates := []string{
		// "templates/terraform//env_var_context.go.tmpl",
		templatePath,
	}
	tmplInput := TestInput{
		Res:                 resource,
		ImportPath:          td.ImportPath(),
		PROJECT_NAME:        "my-project-name",
		CREDENTIALS:         "my/credentials/filename.json",
		REGION:              "us-west1",
		ORG_ID:              "123456789",
		ORG_DOMAIN:          "example.com",
		ORG_TARGET:          "123456789",
		PROJECT_NUMBER:      "1111111111111",
		BILLING_ACCT:        "000000-0000000-0000000-000000",
		MASTER_BILLING_ACCT: "000000-0000000-0000000-000000",
		SERVICE_ACCT:        "my@service-account.com",
		CUST_ID:             "A01b123xz",
		IDENTITY_USER:       "cloud_identity_user",
		PAP_DESCRIPTION:     "description",
	}

	td.GenerateFile(filePath, templatePath, tmplInput, true, templates...)
}

func (td *TemplateData) GenerateFile(filePath, templatePath string, input any, goFormat bool, templates ...string) {
	log.Printf("Generating %s", filePath)

	templateFileName := filepath.Base(templatePath)

	tmpl, err := template.New(templateFileName).Funcs(TemplateFunctions).ParseFiles(templates...)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, input); err != nil {
		glog.Exit(err)
	}

	sourceByte := contents.Bytes()

	if goFormat {
		formattedByte, err := format.Source(sourceByte)
		if err != nil {
			glog.Error(fmt.Errorf("error formatting %s: %s", filePath, err))
		} else {
			sourceByte = formattedByte
		}
	}

	err = os.WriteFile(filePath, sourceByte, 0644)
	if err != nil {
		glog.Exit(err)
	}

	if goFormat {
		cmd := exec.Command("goimports", "-w", filePath)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

//     # path is the output name of the file
//     # template is used to determine metadata about the file based on how it is
//     # generated
//     def format_output_file(path, template)
//       return unless path.end_with?('.go') && @env[:goformat_enabled]

//       run_formatter("gofmt -w -s #{path}")
//       run_formatter("goimports -w #{path}") unless template.include?('third_party/terraform')
//     end

//     def run_formatter(command)
//       output = %x(#{command} 2>&1)
//       Google::LOGGER.error output unless $CHILD_STATUS.to_i.zero?
//     end

//     def relative_path(target, base)
//       Pathname.new(target).relative_path_from(Pathname.new(base))
//     end
//   end

//   # Responsible for compiling provider-level files, rather than product-specific ones
//   class ProviderFileTemplate < Provider::FileTemplate
//     # All the products that are being compiled with the provider on this run
//     attr_accessor :products

//     # Optional path to the directory where overrides reside. Used to locate files
//     # outside of the MM root directory
//     attr_accessor :override_path

//     def initialize(output_folder, version, env, products, override_path = nil)
//       super()

//       @output_folder = output_folder
//       @version = version
//       @env = env
//       @products = products
//       @override_path = override_path
//     end
//   end

//   # Responsible for generating a file in the context of a product
//   # with a given set of parameters.
//   class ProductFileTemplate < Provider::FileTemplate
//     # The name of the resource
//     attr_accessor :name
//     # The resource itself.
//     attr_accessor :object
//     # The entire API object.
//     attr_accessor :product

//     class << self
//       # Construct a new ProductFileTemplate based on a resource object
//       def file_for_resource(output_folder, object, version, env)
//         file_template = new(output_folder, object.name, object.__product, version, env)
//         file_template.object = object
//         file_template
//       end
//     end

//     def initialize(output_folder, name, product, version, env)
//       super()

//       @name = name
//       @product = product
//       @output_folder = output_folder
//       @version = version
//       @env = env
//     end
//   end
// end

//    def import_path
//      case @target_version_name
//      when 'ga'
//        "#{TERRAFORM_PROVIDER_GA}/#{RESOURCE_DIRECTORY_GA}"
//      when 'beta'
//        "#{TERRAFORM_PROVIDER_BETA}/#{RESOURCE_DIRECTORY_BETA}"
//      else
//        "#{TERRAFORM_PROVIDER_PRIVATE}/#{RESOURCE_DIRECTORY_PRIVATE}"
//      end
//    end

func (td *TemplateData) ImportPath() string {
	if td.Version.Name == GA_VERSION {
		return "github.com/hashicorp/terraform-provider-google/google"
	} else if td.Version.Name == ALPHA_VERSION {
		return "internal/terraform-next/google-private"
	}
	return "github.com/hashicorp/terraform-provider-google-beta/google-beta"
}

type TestInput struct {
	Res                 api.Resource
	ImportPath          string
	PROJECT_NAME        string
	CREDENTIALS         string
	REGION              string
	ORG_ID              string
	ORG_DOMAIN          string
	ORG_TARGET          string
	PROJECT_NUMBER      string
	BILLING_ACCT        string
	MASTER_BILLING_ACCT string
	SERVICE_ACCT        string
	CUST_ID             string
	IDENTITY_USER       string
	PAP_DESCRIPTION     string
}
