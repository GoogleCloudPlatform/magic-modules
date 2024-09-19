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
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/slices"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

type TerraformGoogleConversion struct {
	IamResources []map[string]string

	NonDefinedTests []string

	Tests []string

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

	tgc.GenerateIamPolicy(object, *templateData, outputFolder, generateCode, generateDocs)
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

// Generate the IAM policy for this object. This is used to query and test
// IAM policies separately from the resource itself
// Docs are generated for the terraform provider, not here.
func (tgc TerraformGoogleConversion) GenerateIamPolicy(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if !generateCode || object.IamPolicy.ExcludeTgc {
		return
	}

	productName := tgc.Product.ApiName
	targetFolder := path.Join(outputFolder, "converters/google/resources/services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}

	name := object.FilenameOverride
	if name == "" {
		name = google.Underscore(object.Name)
	}

	targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_%s_iam.go", productName, name))
	templateData.GenerateTGCIamResourceFile(targetFilePath, object)

	targetFilePath = path.Join(targetFolder, fmt.Sprintf("iam_%s_%s.go", productName, name))
	templateData.GenerateIamPolicyFile(targetFilePath, object)

	// Don't generate tests - we can rely on the terraform provider
	//  to test these.
}

// Generates the list of resources
func (tgc TerraformGoogleConversion) generateCaiIamResources(products []*api.Product) {
	for _, productDefinition := range products {
		service := strings.ToLower(productDefinition.Name)
		for _, object := range productDefinition.Objects {
			if object.MinVersionObj().Name != "ga" || object.Exclude || object.ExcludeTgc {
				continue
			}

			var iamClassName string
			iamPolicy := object.IamPolicy
			if iamPolicy != nil && !iamPolicy.Exclude && !iamPolicy.ExcludeTgc {

				iamClassName = fmt.Sprintf("%s.ResourceConverter%s", service, object.ResourceName())

				tgc.IamResources = append(tgc.IamResources, map[string]string{
					"TerraformName": object.TerraformName(),
					"IamClassName":  iamClassName,
				})
			}
		}
	}
}

func (tgc TerraformGoogleConversion) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	log.Printf("Compiling common files.")

	tgc.generateCaiIamResources(products)
	tgc.NonDefinedTests = retrieveFullManifestOfNonDefinedTests()

	files := retrieveFullListOfTestFiles()
	for _, file := range files {
		tgc.Tests = append(tgc.Tests, strings.Split(file, ".")[0])
	}
	tgc.Tests = slices.Compact(tgc.Tests)

	testSource := make(map[string]string)
	for target, source := range retrieveTestSourceCodeWithLocation(".tmpl") {
		target := strings.Replace(target, "go.tmpl", "go", 1)
		testSource[target] = source
	}

	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName)
	tgc.CompileFileList(outputFolder, testSource, *templateData, products)

	// compile_file_list(
	//   output_folder,
	//   [
	// 	['converters/google/resources/services/compute/compute_instance_helpers.go',
	// 	 'third_party/terraform/services/compute/compute_instance_helpers.go.erb'],
	// 	['converters/google/resources/resource_converters.go',
	// 	 'templates/tgc/resource_converters.go.erb'],
	// 	['converters/google/resources/services/kms/iam_kms_key_ring.go',
	// 	 'third_party/terraform/services/kms/iam_kms_key_ring.go.erb'],
	// 	['converters/google/resources/services/kms/iam_kms_crypto_key.go',
	// 	 'third_party/terraform/services/kms/iam_kms_crypto_key.go.erb'],
	// 	['converters/google/resources/services/compute/metadata.go',
	// 	 'third_party/terraform/services/compute/metadata.go.erb'],
	// 	['converters/google/resources/services/compute/compute_instance.go',
	// 	 'third_party/tgc/compute_instance.go.erb']
	//   ],
	//   file_template
	// )
}

func (tgc TerraformGoogleConversion) CompileFileList(outputFolder string, files map[string]string, fileTemplate TemplateData, products []*api.Product) {
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
		// tgc.replaceImportPath(outputFolder, target)
	}
}

func retrieveFullManifestOfNonDefinedTests() []string {
	var tests []string
	fileMap := make(map[string]bool)

	files := retrieveFullListOfTestFiles()
	for _, file := range files {
		tests = append(tests, strings.Split(file, ".")[0])
		fileMap[file] = true
	}
	tests = slices.Compact(tests)

	nonDefinedTests := google.Diff(tests, retrieveListOfManuallyDefinedTests())
	nonDefinedTests = google.Reject(nonDefinedTests, func(file string) bool {
		return strings.HasSuffix(file, "_without_default_project")
	})

	for _, test := range nonDefinedTests {
		_, ok := fileMap[fmt.Sprintf("%s.json", test)]
		if !ok {
			log.Fatalf("test file named %s.json expected but found none", test)
		}

		_, ok = fileMap[fmt.Sprintf("%s.tf", test)]
		if !ok {
			log.Fatalf("test file named %s.tf expected but found none", test)
		}
	}

	return nonDefinedTests
}

// Gets all of the test files in the folder third_party/tgc/tests/data
func retrieveFullListOfTestFiles() []string {
	var testFiles []string

	files, err := ioutil.ReadDir("third_party/tgc/tests/data")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		testFiles = append(testFiles, file.Name())
	}
	slices.Sort(testFiles)

	return testFiles
}

func retrieveTestSourceCodeWithLocation(suffix string) map[string]string {
	var fileNames []string
	path := "third_party/tgc/tests/source/go"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Printf("ext %s", filepath.Ext(file.Name()))
		if filepath.Ext(file.Name()) == suffix {
			fileNames = append(fileNames, file.Name())
		}
	}

	slices.Sort(fileNames)

	testSource := make(map[string]string)
	for _, file := range fileNames {
		target := fmt.Sprintf("test/%s", file)
		source := fmt.Sprintf("%s/%s", path, file)
		testSource[target] = source
	}
	return testSource
}

func retrieveListOfManuallyDefinedTests() []string {
	m1 := retrieveListOfManuallyDefinedTestsFromFile("third_party/tgc/tests/source/go/cli_test.go.tmpl")
	m2 := retrieveListOfManuallyDefinedTestsFromFile("third_party/tgc/tests/source/go/read_test.go.tmpl")
	return google.Concat(m1, m2)
}

// Reads the content of the file and then finds all of the tests in the contents
func retrieveListOfManuallyDefinedTestsFromFile(file string) []string {
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Cannot open the file: %v", file)
	}

	var tests []string
	testsReg := regexp.MustCompile(`\s*name\s*:\s*"([^,]+)"`)
	matches := testsReg.FindAllStringSubmatch(string(data), -1)
	for _, testWithName := range matches {
		tests = append(tests, testWithName[1])
	}
	return tests
}

func (tgc TerraformGoogleConversion) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {

}
