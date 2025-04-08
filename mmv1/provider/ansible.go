// Copyright 2025 Google Inc.
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

type Ansible struct {
	ResourceCount       int
	IAMResourceCount    int
	ResourcesForVersion []map[string]string
	TargetVersionName   string
	Version             product.Version
	Product             *api.Product
	StartTime           time.Time
	TemplateData        *TemplateData
}

func NewAnsible(product *api.Product, versionName string, startTime time.Time) Ansible {
	a := Ansible{
		ResourceCount:     0,
		IAMResourceCount:  0,
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
	}

	a.Product.SetPropertiesBasedOnVersion(&a.Version)
	for _, r := range a.Product.Objects {
		r.SetCompiler(ProviderName(a))
		r.ImportPath = ImportPathFromVersion(versionName)
	}

	return a
}

func (a Ansible) Generate(outputFolder, productName, resourceToGenerate string, generateCode, generateDocs, generateTests bool) {
	a.TemplateData = NewTemplateData(outputFolder, a.TargetVersionName, ANSIBLE_PROVIDER)
	if generateCode {
		a.generateAnsibleModules(outputFolder, resourceToGenerate)
	}
	if generateTests {
		a.generateAnsibleTests(outputFolder, resourceToGenerate)
	}
}

func (a *Ansible) generateAnsibleModules(outputFolder, resourceToGenerate string) {
	for _, object := range a.Product.Objects {
		object.ExcludeIfNotInVersion(&a.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			log.Printf("Excluding %s.%s per user request", object.ProductMetadata.Name, object.Name)
			continue
		}
		a.generateAnsibleModule(*object, outputFolder)
	}
}

func (a *Ansible) generateAnsibleModule(object api.Resource, outputFolder string) {
	targetFolder := path.Join(outputFolder, a.TemplateData.AnsiblePluginDirectory)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}
	module := object.AnsibleName()

	log.Printf("Generating module file %s.py for resource %s.%s", module, object.ProductMetadata.Name, object.Name)
	log.Printf("%v", object)
	a.TemplateData.GenerateAnsibleModuleFile(targetFolder, module, object)
}

func (a *Ansible) generateAnsibleTests(outputFolder, resourceToGenerate string) {
	for _, object := range a.Product.Objects {
		object.ExcludeIfNotInVersion(&a.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			log.Printf("Excluding tests for %s.%s per user request", object.ProductMetadata.Name, object.Name)
			continue
		}
		a.generateAnsibleIntegrationTest(*object, outputFolder)
	}
}

func (a *Ansible) generateAnsibleIntegrationTest(object api.Resource, outputFolder string) {
	log.Printf("Generating integration test for %s.%s resource", object.ProductMetadata.Name, object.Name)
	testType := "integration"
	testFolder := path.Join(outputFolder, a.TemplateData.AnsibleTestDirectories[testType])
	module := object.AnsibleName()
	targetFolder := path.Join(testFolder, module)
	subfolder := []string{"defaults", "tasks"}
	for _, sub := range subfolder {
		if err := os.MkdirAll(path.Join(targetFolder, sub), os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
		}
	}
	files := map[string]bool{
		"aliases":                         false,
		path.Join("defaults", "main.yml"): false,
		path.Join("tasks", "autogen.yml"): true,
		path.Join("tasks", "main.yml"):    false,
	}
	for filePath, overwrite := range files {
		a.TemplateData.GenerateAnsibleTestFile(targetFolder, testType, filePath, object, overwrite)
	}
}

// only needed to fully implement Provider interface
func (a Ansible) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
}

// only needed to fully implement Provider interface
func (a Ansible) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {}
