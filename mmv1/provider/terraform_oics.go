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
)

type TerraformOiCS struct {
	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time
}

func NewTerraformOiCS(product *api.Product, versionName string, startTime time.Time) TerraformOiCS {
	toics := TerraformOiCS{
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
	}

	toics.Product.SetPropertiesBasedOnVersion(&toics.Version)

	return toics
}

func (toics TerraformOiCS) Generate(outputFolder, productPath, resourceToGenerate string, generateCode, generateDocs bool) {
	toics.GenerateObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)
}

func (toics TerraformOiCS) GenerateObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	for _, object := range toics.Product.Objects {
		object.ExcludeIfNotInVersion(&toics.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			log.Printf("Excluding %s per user request", object.Name)
			continue
		}

		toics.GenerateObject(*object, outputFolder, toics.TargetVersionName, generateCode, generateDocs)
	}
}

func (toics TerraformOiCS) GenerateObject(object api.Resource, outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	templateData := NewTemplateData(outputFolder, toics.TargetVersionName)

	if !object.IsExcluded() {
		log.Printf("Generating %s resource", object.Name)
		toics.GenerateResource(object, *templateData, outputFolder, generateCode, generateDocs)
	}
}

func (toics TerraformOiCS) GenerateResource(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if !generateDocs {
		return
	}

	for _, example := range object.TestExamples() {
		if len(example.TestEnvVars) > 0 {
			continue
		}

		example.SetOiCSHCLText()

		targetFolder := path.Join(outputFolder, example.Name)

		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating oics example directory %v: %v", targetFolder, err))
		}

		oicsExampleTemplatePath := "templates/terraform/examples/base_configs/oics_example_file.tf.tmpl"
		oicsExampleTemplates := []string{
			oicsExampleTemplatePath,
		}
		templateData.GenerateFile(path.Join(targetFolder, "main.tf"), oicsExampleTemplatePath, example, false, oicsExampleTemplates...)

		tutorialTemplatePath := "templates/terraform/examples/base_configs/tutorial.md.tmpl"
		tutorialTemplates := []string{
			tutorialTemplatePath,
		}
		templateData.GenerateFile(path.Join(targetFolder, "tutorial.md"), tutorialTemplatePath, example, false, tutorialTemplates...)

		backingTemplatePath := "templates/terraform/examples/base_configs/example_backing_file.tf.tmpl"
		backingTemplates := []string{
			backingTemplatePath,
		}
		templateData.GenerateFile(path.Join(targetFolder, "backing_file.tf"), backingTemplatePath, example, false, backingTemplates...)

		motdTemplatePath := "templates/terraform/examples/static/motd.tmpl"
		motdTemplates := []string{
			motdTemplatePath,
		}
		templateData.GenerateFile(path.Join(targetFolder, "motd"), motdTemplatePath, example, false, motdTemplates...)
	}
}

func (toics TerraformOiCS) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {

}

func (toics TerraformOiCS) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {

}
