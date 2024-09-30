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
	"os/exec"
	"path/filepath"
	"strings"

	"text/template"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/golang/glog"
)

type TemplateData struct {
	OutputFolder string
	VersionName  string

	TerraformResourceDirectory string
	TerraformProviderModule    string

	// TODO rewrite: is this needed?
	//     # Information about the local environment
	//     # (which formatters are enabled, start-time)
	//     attr_accessor :env
}

var GA_VERSION = "ga"
var BETA_VERSION = "beta"
var ALPHA_VERSION = "alpha"
var PRIVATE_VERSION = "private"

func NewTemplateData(outputFolder string, versionName string) *TemplateData {
	td := TemplateData{OutputFolder: outputFolder, VersionName: versionName}

	if versionName == GA_VERSION {
		td.TerraformResourceDirectory = "google"
		td.TerraformProviderModule = "github.com/hashicorp/terraform-provider-google"
	} else if versionName == ALPHA_VERSION || versionName == PRIVATE_VERSION {
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
		"templates/terraform/custom_flatten/bigquery_table_ref.go.tmpl",
		"templates/terraform/flatten_property_method.go.tmpl",
		"templates/terraform/expand_property_method.go.tmpl",
		"templates/terraform/update_mask.go.tmpl",
		"templates/terraform/nested_query.go.tmpl",
		"templates/terraform/unordered_list_customize_diff.go.tmpl",
	}
	td.GenerateFile(filePath, templatePath, resource, true, templates...)
}

func (td *TemplateData) GenerateOperationFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/operation.go.tmpl"
	templates := []string{
		templatePath,
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
		"templates/terraform/env_var_context.go.tmpl",
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

func (td *TemplateData) GenerateIamPolicyFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/iam_policy.go.tmpl"
	templates := []string{
		templatePath,
	}
	td.GenerateFile(filePath, templatePath, resource, true, templates...)
}

func (td *TemplateData) GenerateIamResourceDocumentationFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/resource_iam.html.markdown.tmpl"
	templates := []string{
		templatePath,
	}
	td.GenerateFile(filePath, templatePath, resource, false, templates...)
}

func (td *TemplateData) GenerateIamDatasourceDocumentationFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/datasource_iam.html.markdown.tmpl"
	templates := []string{
		templatePath,
	}
	td.GenerateFile(filePath, templatePath, resource, false, templates...)
}

func (td *TemplateData) GenerateIamPolicyTestFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/examples/base_configs/iam_test_file.go.tmpl"
	templates := []string{
		templatePath,
		"templates/terraform/env_var_context.go.tmpl",
		"templates/terraform/iam/iam_context.go.tmpl",
	}
	td.GenerateFile(filePath, templatePath, resource, true, templates...)
}

func (td *TemplateData) GenerateSweeperFile(filePath string, resource api.Resource) {
	templatePath := "templates/terraform/sweeper_file.go.tmpl"
	templates := []string{
		templatePath,
	}
	td.GenerateFile(filePath, templatePath, resource, false, templates...)
}

func (td *TemplateData) GenerateTGCResourceFile(filePath string, resource api.Resource) {
	templatePath := "templates/tgc/resource_converter.go.tmpl"
	templates := []string{
		templatePath,
		"templates/terraform/expand_property_method.go.tmpl",
	}
	td.GenerateFile(filePath, templatePath, resource, true, templates...)
}

func (td *TemplateData) GenerateTGCIamResourceFile(filePath string, resource api.Resource) {
	templatePath := "templates/tgc/resource_converter_iam.go.tmpl"
	templates := []string{
		templatePath,
	}
	td.GenerateFile(filePath, templatePath, resource, true, templates...)
}

func (td *TemplateData) GenerateFile(filePath, templatePath string, input any, goFormat bool, templates ...string) {
	templateFileName := filepath.Base(templatePath)

	tmpl, err := template.New(templateFileName).Funcs(google.TemplateFunctions).ParseFiles(templates...)
	if err != nil {
		glog.Exit(fmt.Sprintf("error parsing %s for filepath %s ", templateFileName, filePath), err)
	}

	contents := bytes.Buffer{}
	if err = tmpl.ExecuteTemplate(&contents, templateFileName, input); err != nil {
		glog.Exit(fmt.Sprintf("error executing %s for filepath %s ", templateFileName, filePath), err)
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

	if goFormat && !strings.Contains(templatePath, "third_party/terraform") {
		cmd := exec.Command("goimports", "-w", filepath.Base(filePath))
		cmd.Dir = filepath.Dir(filePath)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

func (td *TemplateData) ImportPath() string {
	if td.VersionName == GA_VERSION {
		return "github.com/hashicorp/terraform-provider-google/google"
	} else if td.VersionName == ALPHA_VERSION || td.VersionName == PRIVATE_VERSION {
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
