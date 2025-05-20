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
	tgc.CompileTfToCaiCommonFiles(outputFolder, products)
	tgc.CompileCaiToHclCommonFiles(outputFolder, products)
}

func (tgc TerraformGoogleConversionNext) CompileTfToCaiCommonFiles(outputFolder string, products []*api.Product) {
	log.Printf("Compiling common files for tgc tfplan2cai.")

	resourceConverters := map[string]string{
		"pkg/tfplan2cai/converters/resource_converters.go":                       "templates/tgc_next/tfplan2cai/resource_converters.go.tmpl",
		"pkg/tfplan2cai/converters/services/compute/compute_instance_helpers.go": "third_party/terraform/services/compute/compute_instance_helpers.go.tmpl",
		"pkg/tfplan2cai/converters/services/compute/metadata.go":                 "third_party/terraform/services/compute/metadata.go.tmpl",
	}
	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
}

func (tgc TerraformGoogleConversionNext) CompileCaiToHclCommonFiles(outputFolder string, products []*api.Product) {
	log.Printf("Compiling common files for tgc tfplan2cai.")

	resourceConverters := map[string]string{
		"pkg/cai2hcl/converters/resource_converters.go": "templates/tgc_next/cai2hcl/resource_converters.go.tmpl",
	}
	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
}

func (tgc TerraformGoogleConversionNext) CompileFileList(outputFolder string, files map[string]string, fileTemplate TemplateData, products []*api.Product) {
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

	tgc.CopyTfToCaiCommonFiles(outputFolder)
	tgc.CopyCaiToHclCommonFiles(outputFolder)
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
	sourceByte = bytes.Replace(sourceByte, []byte(gaImportPath), []byte(TERRAFORM_PROVIDER_BETA+"/"+RESOURCE_DIRECTORY_BETA), -1)
	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to replace import path: %s", target, err)
	}
}
