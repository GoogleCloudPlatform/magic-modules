// Copyright 2025 Red Hat Inc.
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
	"os"
	"path"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/golang/glog"
)

type Ansible struct {
	ResourceCount       int
	IAMResourceCount    int
	ResourcesForVersion []map[string]string
	TargetVersionName   string
	Version             product.Version
	Product             *api.Product
	StartTime           time.Time
	TemplateData        *AnsibleTemplateData
	Overwrite           bool
}

func NewAnsible(product *api.Product, versionName string, startTime time.Time, overwrite bool) Ansible {
	a := Ansible{
		ResourceCount:     0,
		IAMResourceCount:  0,
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
		Overwrite:         overwrite,
	}

	a.Product.SetPropertiesBasedOnVersion(&a.Version)
	for _, r := range a.Product.Objects {
		r.SetCompiler(ProviderName(a))
		r.ImportPath = ImportPathFromVersion(versionName)
	}

	return a
}

func (a Ansible) Generate(outputFolder, productName, resourceToGenerate string, generateCode, generateDocs bool) {
	a.TemplateData = NewAnsibleTemplateData(outputFolder, a.TargetVersionName, a.Overwrite)
	if generateCode {
		a.generateAnsibleModules(outputFolder, resourceToGenerate)
	}
	a.generateAnsibleTests(outputFolder, resourceToGenerate)
}

func (a *Ansible) generateAnsibleModules(outputFolder, resourceToGenerate string) {
	for _, r := range a.Product.Objects {
		r.ExcludeIfNotInVersion(&a.Version)

		if resourceToGenerate != "" && r.Name != resourceToGenerate {
			// glog.Infof("excluding %s.%s per user request", r.ProductMetadata.Name, r.Name)
			continue
		}

		targetFolder := path.Join(outputFolder, a.TemplateData.ModuleDirectory)
		targetFile := path.Join(targetFolder, fmt.Sprintf("%s.py", r.AnsibleName()))
		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			glog.Fatal(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
		}

		glog.Infof("generating module file %s for resource %s.%s", targetFile, r.ProductMetadata.Name, r.Name)

		if err := a.TemplateData.GenerateModuleFile(targetFile, r); err != nil {
			glog.Fatal(fmt.Errorf("error creating module file for %v: %v", r.AnsibleName(), err))
		}
	}
}

func (a *Ansible) generateAnsibleTests(outputFolder, resourceToGenerate string) {
	for _, object := range a.Product.Objects {
		object.ExcludeIfNotInVersion(&a.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			// glog.Infof("excluding tests for %s.%s", object.ProductMetadata.Name, object.Name)
			continue
		}
		a.generateAnsibleIntegrationTest(object, outputFolder)
	}
}

func (a *Ansible) generateAnsibleIntegrationTest(r *api.Resource, outputFolder string) {
	glog.Infof("generating integration test for %s.%s", r.ProductMetadata.Name, r.Name)
	testFolder := path.Join(outputFolder, a.TemplateData.TestDirectories["integration"], r.AnsibleName())
	subfolder := []string{"defaults", "tasks"}
	for _, sub := range subfolder {
		// glog.Infof("mkdir %s/%s", targetFolder, sub)
		targetFolder := path.Join(testFolder, sub)
		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			glog.Fatal(fmt.Errorf("error creating directory %v: %v", targetFolder, err))
		}
	}
	files := []string{
		"aliases",
		path.Join("defaults", "main.yml"),
		path.Join("tasks", "autogen.yml"),
		path.Join("tasks", "main.yml"),
	}
	for _, targetFile := range files {
		if err := a.TemplateData.GenerateTestFile(outputFolder, targetFile, "integration", r); err != nil {
			glog.Fatal(err)
		}
	}
}

// only needed to fully implement Provider interface
func (a Ansible) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	glog.Info("compile common files")
}

// only needed to fully implement Provider interface
func (a Ansible) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	glog.Info("copy common files")
}
