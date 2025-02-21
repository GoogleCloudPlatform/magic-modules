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
type TerraformGoogleConversionV7 struct {
	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time
}

func NewTerraformGoogleConversionV7(product *api.Product, versionName string, startTime time.Time) TerraformGoogleConversionV7 {
	t := TerraformGoogleConversionV7{
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

func (tgc TerraformGoogleConversionV7) Generate(outputFolder, productPath, resourceToGenerate string, generateCode, generateDocs bool) {
	tgc.GenerateTfToCaiObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)
	tgc.GenerateCaiToHclObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)
}

func (tgc TerraformGoogleConversionV7) GenerateTfToCaiObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
}

func (tgc TerraformGoogleConversionV7) GenerateCaiToHclObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
}

func (tgc TerraformGoogleConversionV7) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	tgc.CompileTfToCaiCommonFiles(outputFolder, products)
	tgc.CompileCaiToHclCommonFiles(outputFolder, products)
}

func (tgc TerraformGoogleConversionV7) CompileTfToCaiCommonFiles(outputFolder string, products []*api.Product) {
	log.Printf("Compiling common files for tgc v7 tfplan2cai.")

	resourceConverters := map[string]string{
		"tfplan2cai/converters/resource_converters.go":                       "templates/tgc_v7/tfplan2cai/resource_converters.go.tmpl",
		"tfplan2cai/converters/services/compute/compute_instance_helpers.go": "third_party/terraform/services/compute/compute_instance_helpers.go.tmpl",
		"tfplan2cai/converters/services/compute/metadata.go":                 "third_party/terraform/services/compute/metadata.go.tmpl",
	}
	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
}

func (tgc TerraformGoogleConversionV7) CompileCaiToHclCommonFiles(outputFolder string, products []*api.Product) {
	log.Printf("Compiling common files for tgc v7 tfplan2cai.")

	resourceConverters := map[string]string{
		"cai2hcl/converters/resource_converters.go": "templates/tgc_v7/cai2hcl/resource_converters.go.tmpl",
	}
	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
}

func (tgc TerraformGoogleConversionV7) CompileFileList(outputFolder string, files map[string]string, fileTemplate TemplateData, products []*api.Product) {
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

		fileTemplate.GenerateFile(targetFile, source, tgc, formatFile, templates...)
		tgc.replaceImportPath(outputFolder, target)
	}
}

func (tgc TerraformGoogleConversionV7) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	if !generateCode {
		return
	}

	log.Printf("Copying common files for tgc v7.")

	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	if err := copy.Copy("third_party/tgc_v7", outputFolder); err != nil {
		log.Println(fmt.Errorf("error copying directory %v: %v", outputFolder, err))
	}

	tgc.CopyTfToCaiCommonFiles(outputFolder)
	tgc.CopyCaiToHclCommonFiles(outputFolder)
}

func (tgc TerraformGoogleConversionV7) CopyTfToCaiCommonFiles(outputFolder string) {
	resourceConverters := map[string]string{
		"tfplan2cai/converters/services/compute/image.go":     "third_party/terraform/services/compute/image.go",
		"tfplan2cai/converters/services/compute/disk_type.go": "third_party/terraform/services/compute/disk_type.go",
	}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionV7) CopyCaiToHclCommonFiles(outputFolder string) {
	resourceConverters := map[string]string{}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionV7) CopyFileList(outputFolder string, files map[string]string) {
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

func (tgc TerraformGoogleConversionV7) replaceImportPath(outputFolder, target string) {
	// Replace import paths to reference the resources dir instead of the google provider
	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to replace import path: %s", targetFile, err)
	}

	// replace google to google-beta
	gaImportPath := ImportPathFromVersion("ga")
	sourceByte = bytes.Replace(sourceByte, []byte(gaImportPath), []byte(TERRAFORM_PROVIDER_BETA+"/"+RESOURCE_DIRECTORY_BETA), -1)
	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to replace import path: %s", target, err)
	}
}
