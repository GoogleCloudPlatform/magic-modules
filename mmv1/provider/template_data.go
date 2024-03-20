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
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"

	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
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

var TemplateFunctions = template.FuncMap{
	"title": strings.Title,
	// "patternToRegex":                  PatternToRegex,
	"replace": strings.Replace,
	// "isLastIndex":                     isLastIndex,
	// "escapeDescription":               escapeDescription,
	// "shouldAllowForwardSlashInFormat": shouldAllowForwardSlashInFormat,
}

var GA_VERSION = "ga"
var BETA_VERSION = "beta"
var ALPHA_VERSION = "alpha"

func NewTemplateData(outputFolder string, version product.Version) *TemplateData {
	td := TemplateData{OutputFolder: outputFolder, Version: version}

	if version.Name == GA_VERSION {
		td.TerraformResourceDirectory = "google"
		td.TerraformProviderModule = "github.com/hashicorp/terraform-provider-google/google"
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

	log.Printf("Generating %s", filePath)

	tmpl, err := template.New("resource.go.tmpl").Funcs(TemplateFunctions).ParseFiles(
		"templates/terraform/resource.go.tmpl",
	)
	if err != nil {
		glog.Exit(err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, "resource.go.tmpl", resource); err != nil {
		glog.Exit(err)
	}

	if err != nil {
		glog.Exit(err)
	}

	formatted, err := td.FormatSource(&contents)
	if err != nil {
		glog.Error(fmt.Errorf("error formatting %s", filePath))
	}

	err = os.WriteFile(filePath, formatted, 0644)
	if err != nil {
		glog.Exit(err)
	}
}

func (td *TemplateData) GenerateDocumentationFile(filePath string, resource api.Resource) {

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

func (td *TemplateData) FormatSource(source *bytes.Buffer) ([]byte, error) {
	sourceByte := source.Bytes()
	// Replace import path based on version (beta/alpha)
	if td.TerraformResourceDirectory != "google" {
		sourceByte = bytes.Replace(sourceByte, []byte("github.com/hashicorp/terraform-provider-google/google"), []byte(td.TerraformProviderModule+"/"+td.TerraformResourceDirectory), -1)
	}

	output, err := format.Source(sourceByte)
	if err != nil {
		return []byte(source.String()), err
	}

	return output, nil
}
