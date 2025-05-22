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

package provider

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/otiai10/copy"
)

// This proivder is for both tfplan2cai and cai2hcl conversions,
// and copying other files, such as transport.go
type TerraformGoogleConversionNext struct {
	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time
}

func NewTerraformGoogleConversionNext(product *api.Product, versionName string, startTime time.Time) TerraformGoogleConversionNext {
	t := TerraformGoogleConversionNext{
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
	}

	t.Product.SetPropertiesBasedOnVersion(&t.Version)
	for _, r := range t.Product.Objects {
		r.SetCompiler(ProviderName(t))
		r.ImportPath = ImportPathFromVersion(versionName)
	}

	return t
}

func (tgc TerraformGoogleConversionNext) Generate(outputFolder, productPath, resourceToGenerate string, generateCode, generateDocs bool) {
	tgc.GenerateTfToCaiObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)
	tgc.GenerateCaiToHclObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)
}

func (tgc TerraformGoogleConversionNext) GenerateTfToCaiObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
}

func (tgc TerraformGoogleConversionNext) GenerateCaiToHclObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
}

func (tgc TerraformGoogleConversionNext) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	resourceConverters := map[string]string{
		// common
		"pkg/transport/config.go":                        "third_party/terraform/transport/config.go.tmpl",
		"pkg/transport/provider_handwritten_endpoint.go": "third_party/terraform/transport/provider_handwritten_endpoint.go.tmpl",
		"pkg/tpgresource/common_diff_suppress.go":        "third_party/terraform/tpgresource/common_diff_suppress.go.tmpl",
		"pkg/provider/provider.go":                       "third_party/terraform/provider/provider.go.tmpl",
		"pkg/provider/provider_validators.go":            "third_party/terraform/provider/provider_validators.go.tmpl",

		// tfplan2cai
		"pkg/tfplan2cai/converters/resource_converters.go":                       "templates/tgc_next/tfplan2cai/resource_converters.go.tmpl",
		"pkg/tfplan2cai/converters/services/compute/compute_instance_helpers.go": "third_party/terraform/services/compute/compute_instance_helpers.go.tmpl",
		"pkg/tfplan2cai/converters/services/compute/metadata.go":                 "third_party/terraform/services/compute/metadata.go.tmpl",

		// cai2hcl
		"pkg/cai2hcl/converters/resource_converters.go": "templates/tgc_next/cai2hcl/resource_converters.go.tmpl",
	}

	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
}

func (tgc TerraformGoogleConversionNext) CompileFileList(outputFolder string, files map[string]string, fileTemplate TemplateData, products []*api.Product) {
	providerWithProducts := TgcWithProducts{
		TerraformGoogleConversionNext: tgc,
		Compiler:                      "terraformgoogleconversion-codegen",
		Products:                      products,
	}

	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	for target, source := range files {
		targetFile := filepath.Join(outputFolder, target)
		targetDir := filepath.Dir(targetFile)
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating output directory %v: %v", targetDir, err))
		}

		templates := []string{
			source,
		}

		formatFile := filepath.Ext(targetFile) == ".go"

		fileTemplate.GenerateFile(targetFile, source, providerWithProducts, formatFile, templates...)
		tgc.replaceImportPath(outputFolder, target)
	}
}

func (tgc TerraformGoogleConversionNext) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	if !generateCode {
		return
	}

	log.Printf("Copying common files for tgc.")

	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	if err := copy.Copy("third_party/tgc_next", outputFolder); err != nil {
		log.Println(fmt.Errorf("error copying directory %v: %v", outputFolder, err))
	}

	resourceConverters := map[string]string{
		// common
		"pkg/transport/batcher.go":                 "third_party/terraform/transport/batcher.go",
		"pkg/transport/retry_transport.go":         "third_party/terraform/transport/retry_transport.go",
		"pkg/transport/retry_utils.go":             "third_party/terraform/transport/retry_utils.go",
		"pkg/transport/header_transport.go":        "third_party/terraform/transport/header_transport.go",
		"pkg/transport/error_retry_predicates.go":  "third_party/terraform/transport/error_retry_predicates.go",
		"pkg/transport/bigtable_client_factory.go": "third_party/terraform/transport/bigtable_client_factory.go",
		"pkg/transport/transport.go":               "third_party/terraform/transport/transport.go",
		"pkg/tpgresource/utils.go":                 "third_party/terraform/tpgresource/utils.go",
		"pkg/tpgresource/self_link_helpers.go":     "third_party/terraform/tpgresource/self_link_helpers.go",
		"pkg/tpgresource/hashcode.go":              "third_party/terraform/tpgresource/hashcode.go",
		"pkg/tpgresource/regional_utils.go":        "third_party/terraform/tpgresource/regional_utils.go",
		"pkg/tpgresource/field_helpers.go":         "third_party/terraform/tpgresource/field_helpers.go",
		"pkg/tpgresource/service_scope.go":         "third_party/terraform/tpgresource/service_scope.go",
		"pkg/provider/mtls_util.go":                "third_party/terraform/provider/mtls_util.go",
		"pkg/verify/validation.go":                 "third_party/terraform/verify/validation.go",
		"pkg/verify/path_or_contents.go":           "third_party/terraform/verify/path_or_contents.go",
		"pkg/version/version.go":                   "third_party/terraform/version/version.go",

		// tfplan2cai
		"pkg/tfplan2cai/converters/services/compute/image.go":     "third_party/terraform/services/compute/image.go",
		"pkg/tfplan2cai/converters/services/compute/disk_type.go": "third_party/terraform/services/compute/disk_type.go",
	}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionNext) CopyTfToCaiCommonFiles(outputFolder string) {
	resourceConverters := map[string]string{
		"pkg/tfplan2cai/converters/services/compute/image.go":     "third_party/terraform/services/compute/image.go",
		"pkg/tfplan2cai/converters/services/compute/disk_type.go": "third_party/terraform/services/compute/disk_type.go",
	}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionNext) CopyCaiToHclCommonFiles(outputFolder string) {
	resourceConverters := map[string]string{}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionNext) CopyFileList(outputFolder string, files map[string]string) {
	for target, source := range files {
		targetFile := filepath.Join(outputFolder, target)
		targetDir := filepath.Dir(targetFile)

		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating output directory %v: %v", targetDir, err))
		}
		// If we've modified a file since starting an MM run, it's a reasonable
		// assumption that it was this run that modified it.
		if info, err := os.Stat(targetFile); !errors.Is(err, os.ErrNotExist) && tgc.StartTime.Before(info.ModTime()) {
			log.Fatalf("%s was already modified during this run at %s", targetFile, info.ModTime().String())
		}

		sourceByte, err := os.ReadFile(source)
		if err != nil {
			log.Fatalf("Cannot read source file %s while copying: %s", source, err)
		}

		err = os.WriteFile(targetFile, sourceByte, 0644)
		if err != nil {
			log.Fatalf("Cannot write target file %s while copying: %s", target, err)
		}

		// Replace import path based on version (beta/alpha)
		if filepath.Ext(target) == ".go" || filepath.Ext(target) == ".mod" {
			tgc.replaceImportPath(outputFolder, target)
		}
	}
}

func (tgc TerraformGoogleConversionNext) replaceImportPath(outputFolder, target string) {
	// Replace import paths to reference the resources dir instead of the google provider
	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to replace import path: %s", targetFile, err)
	}

	// replace google to google-beta
	gaImportPath := ImportPathFromVersion("ga")
	sourceByte = bytes.Replace(sourceByte, []byte(gaImportPath), []byte(TGC_PROVIDER+"/"+RESOURCE_DIRECTORY_TGC), -1)
	sourceByte = bytes.Replace(sourceByte, []byte(TERRAFORM_PROVIDER_GA+"/version"), []byte(TGC_PROVIDER+"/"+RESOURCE_DIRECTORY_TGC+"/version"), -1)

	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to replace import path: %s", target, err)
	}
}

type TgcWithProducts struct {
	TerraformGoogleConversionNext
	Compiler string
	Products []*api.Product
}
