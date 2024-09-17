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
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

type TerraformGoogleConversion struct {
	ResourcesForVersion []map[string]string

	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time
}

func NewTerraformGoogleConversion(product *api.Product, versionName string, startTime time.Time) TerraformGoogleConversion {
	t := TerraformGoogleConversion{
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
	}

	t.Product.SetPropertiesBasedOnVersion(&t.Version)
	for _, r := range t.Product.Objects {
		r.SetCompiler(ProviderName(t))
		r.ImportPath = ImportPathFromVersion(t, versionName)
	}

	return t
}

func (tgc TerraformGoogleConversion) generatingHashicorpRepo() bool {
	// This code is not used when generating TPG/TPGB
	return false
}

func (tgc TerraformGoogleConversion) Generate(outputFolder, productPath, resourceToGenerate string, generateCode, generateDocs bool) {
	// Temporary shim to generate the missing resources directory. Can be removed
	// once the folder exists downstream.
	resourcesFolder := path.Join(outputFolder, "converters/google/resources")
	if err := os.MkdirAll(resourcesFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", resourcesFolder, err))
	}
	tgc.GenerateObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)
}

func (tgc TerraformGoogleConversion) GenerateObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	for _, object := range tgc.Product.Objects {
		object.ExcludeIfNotInVersion(&tgc.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			log.Printf("Excluding %s per user request", object.Name)
			continue
		}

		tgc.GenerateObject(*object, outputFolder, tgc.TargetVersionName, generateCode, generateDocs)
	}
}

func (tgc TerraformGoogleConversion) GenerateObject(object api.Resource, outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	if object.ExcludeTgc {
		log.Printf("Skipping fine-grained resource %s", object.Name)
		return
	}

	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)

	if !object.IsExcluded() {
		tgc.GenerateResource(object, *templateData, outputFolder, generateCode, generateDocs)

		if generateCode {
			// tgc.GenerateResourceTests(object, *templateData, outputFolder)
			// tgc.GenerateResourceSweeper(object, *templateData, outputFolder)
		}
	}

	// if iam_policy is not defined or excluded, don't generate it
	if object.IamPolicy == nil || object.IamPolicy.Exclude {
		return
	}

	// tgc.GenerateIamPolicy(object, *templateData, outputFolder, generateCode, generateDocs)
}

func (tgc TerraformGoogleConversion) GenerateResource(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	productName := tgc.Product.ApiName
	targetFolder := path.Join(outputFolder, "converters/google/resources/services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}

	targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_%s.go", productName, google.Underscore(object.Name)))
	templateData.GenerateTGCResourceFile(targetFilePath, object)
}

func (tgc TerraformGoogleConversion) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {

}

func (tgc TerraformGoogleConversion) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {

}
