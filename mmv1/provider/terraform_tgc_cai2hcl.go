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

package provider

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/otiai10/copy"
)

// Code generator for a library converting GCP CAI objects to Terraform state.
type CaiToTerraformConversion struct {
	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time
}

func NewCaiToTerraformConversion(product *api.Product, versionName string, startTime time.Time) CaiToTerraformConversion {
	t := CaiToTerraformConversion{
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

func (cai2hcl CaiToTerraformConversion) Generate(outputFolder, productPath, resourceToGenerate string, generateCode, generateDocs bool) {
}

func (cai2hcl CaiToTerraformConversion) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
}

func (cai2hcl CaiToTerraformConversion) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	if !generateCode {
		return
	}
	log.Printf("Coping cai2hcl common files")

	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	if err := copy.Copy("third_party/cai2hcl", outputFolder); err != nil {
		log.Println(fmt.Errorf("error copying directory %v: %v", outputFolder, err))
	}
}
