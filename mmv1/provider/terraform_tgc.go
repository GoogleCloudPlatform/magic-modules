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
		r.ImportPath = ImportPathFromVersion(versionName)
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
func (tgc *TerraformGoogleConversion) generateCaiIamResources(products []*api.Product) {
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
	log.Printf("Compiling common files for tgc.")

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

	resourceConverters := map[string]string{
		"converters/google/resources/services/compute/compute_instance_helpers.go": "third_party/terraform/services/compute/go/compute_instance_helpers.go.tmpl",
		"converters/google/resources/resource_converters.go":                       "third_party/tgc/resource_converters.go.tmpl",
		"converters/google/resources/services/kms/iam_kms_key_ring.go":             "third_party/terraform/services/kms/go/iam_kms_key_ring.go.tmpl",
		"converters/google/resources/services/kms/iam_kms_crypto_key.go":           "third_party/terraform/services/kms/go/iam_kms_crypto_key.go.tmpl",
		"converters/google/resources/services/compute/metadata.go":                 "third_party/terraform/services/compute/go/metadata.go.tmpl",
		"converters/google/resources/services/compute/compute_instance.go":         "third_party/tgc/services/compute/compute_instance.go.tmpl",
	}
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
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
		tgc.replaceImportPath(outputFolder, target)
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

// Gets all of files in the folder third_party/tgc/tests/data
func retrieveFullListOfTestTilesWithLocation() map[string]string {
	testFiles := make(map[string]string)
	files := retrieveFullListOfTestFiles()
	for _, file := range files {
		target := fmt.Sprintf("testdata/templates/%s", file)
		source := fmt.Sprintf("third_party/tgc/tests/data/%s", file)
		testFiles[target] = source
	}
	return testFiles
}

func retrieveTestSourceCodeWithLocation(suffix string) map[string]string {
	var fileNames []string
	var path string
	if suffix == ".tmpl" {
		path = "third_party/tgc/tests/source/go"
	} else {
		path = "third_party/tgc/tests/source"
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
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
	log.Printf("Copying common files.")

	if !generateCode {
		return
	}

	tgc.CopyFileList(outputFolder, retrieveFullListOfTestTilesWithLocation())
	tgc.CopyFileList(outputFolder, retrieveTestSourceCodeWithLocation(".go"))

	resourceConverters := map[string]string{
		"converters/google/resources/cai/constants.go":                             "third_party/tgc/cai/constants.go",
		"converters/google/resources/constants.go":                                 "third_party/tgc/constants.go",
		"converters/google/resources/cai.go":                                       "third_party/tgc/cai.go",
		"converters/google/resources/cai/cai.go":                                   "third_party/tgc/cai/cai.go",
		"converters/google/resources/cai/cai_test.go":                              "third_party/tgc/cai/cai_test.go",
		"converters/google/resources/org_policy_policy.go":                         "third_party/tgc/org_policy_policy.go",
		"converters/google/resources/getconfig.go":                                 "third_party/tgc/getconfig.go",
		"converters/google/resources/folder.go":                                    "third_party/tgc/folder.go",
		"converters/google/resources/getconfig_test.go":                            "third_party/tgc/getconfig_test.go",
		"converters/google/resources/cai/json_map.go":                              "third_party/tgc/cai/json_map.go",
		"converters/google/resources/project.go":                                   "third_party/tgc/project.go",
		"converters/google/resources/sql_database_instance.go":                     "third_party/tgc/sql_database_instance.go",
		"converters/google/resources/storage_bucket.go":                            "third_party/tgc/storage_bucket.go",
		"converters/google/resources/cloudfunctions_function.go":                   "third_party/tgc/cloudfunctions_function.go",
		"converters/google/resources/cloudfunctions_cloud_function.go":             "third_party/tgc/cloudfunctions_cloud_function.go",
		"converters/google/resources/bigquery_table.go":                            "third_party/tgc/bigquery_table.go",
		"converters/google/resources/bigtable_cluster.go":                          "third_party/tgc/bigtable_cluster.go",
		"converters/google/resources/bigtable_instance.go":                         "third_party/tgc/bigtable_instance.go",
		"converters/google/resources/cai/iam_helpers.go":                           "third_party/tgc/cai/iam_helpers.go",
		"converters/google/resources/cai/iam_helpers_test.go":                      "third_party/tgc/cai/iam_helpers_test.go",
		"converters/google/resources/services/resourcemanager/organization_iam.go": "third_party/tgc/organization_iam.go",
		"converters/google/resources/services/resourcemanager/project_iam.go":      "third_party/tgc/project_iam.go",
		"converters/google/resources/project_organization_policy.go":               "third_party/tgc/project_organization_policy.go",
		"converters/google/resources/folder_organization_policy.go":                "third_party/tgc/folder_organization_policy.go",
		"converters/google/resources/services/resourcemanager/folder_iam.go":       "third_party/tgc/folder_iam.go",
		"converters/google/resources/container.go":                                 "third_party/tgc/container.go",
		"converters/google/resources/project_service.go":                           "third_party/tgc/project_service.go",
		"converters/google/resources/services/monitoring/monitoring_slo_helper.go": "third_party/tgc/monitoring_slo_helper.go",
		"converters/google/resources/service_account.go":                           "third_party/tgc/service_account.go",
		"converters/google/resources/services/compute/image.go":                    "third_party/terraform/services/compute/image.go",
		"converters/google/resources/services/compute/disk_type.go":                "third_party/terraform/services/compute/disk_type.go",
		"converters/google/resources/services/kms/kms_utils.go":                    "third_party/terraform/services/kms/kms_utils.go",
		"converters/google/resources/services/sourcerepo/source_repo_utils.go":     "third_party/terraform/services/sourcerepo/source_repo_utils.go",
		"converters/google/resources/services/pubsub/pubsub_utils.go":              "third_party/terraform/services/pubsub/pubsub_utils.go",
		"converters/google/resources/services/resourcemanager/iam_organization.go": "third_party/terraform/services/resourcemanager/iam_organization.go",
		"converters/google/resources/services/resourcemanager/iam_folder.go":       "third_party/terraform/services/resourcemanager/iam_folder.go",
		"converters/google/resources/services/resourcemanager/iam_project.go":      "third_party/terraform/services/resourcemanager/iam_project.go",
		"converters/google/resources/services/privateca/privateca_utils.go":        "third_party/terraform/services/privateca/privateca_utils.go",
		"converters/google/resources/services/bigquery/iam_bigquery_dataset.go":    "third_party/terraform/services/bigquery/iam_bigquery_dataset.go",
		"converters/google/resources/services/bigquery/bigquery_dataset_iam.go":    "third_party/tgc/bigquery_dataset_iam.go",
		"converters/google/resources/compute_security_policy.go":                   "third_party/tgc/compute_security_policy.go",
		"converters/google/resources/kms_key_ring_iam.go":                          "third_party/tgc/kms_key_ring_iam.go",
		"converters/google/resources/kms_crypto_key_iam.go":                        "third_party/tgc/kms_crypto_key_iam.go",
		"converters/google/resources/project_iam_custom_role.go":                   "third_party/tgc/project_iam_custom_role.go",
		"converters/google/resources/organization_iam_custom_role.go":              "third_party/tgc/organization_iam_custom_role.go",
		"converters/google/resources/services/pubsub/iam_pubsub_subscription.go":   "third_party/terraform/services/pubsub/iam_pubsub_subscription.go",
		"converters/google/resources/services/pubsub/pubsub_subscription_iam.go":   "third_party/tgc/pubsub_subscription_iam.go",
		"converters/google/resources/services/spanner/iam_spanner_database.go":     "third_party/terraform/services/spanner/iam_spanner_database.go",
		"converters/google/resources/services/spanner/spanner_database_iam.go":     "third_party/tgc/spanner_database_iam.go",
		"converters/google/resources/services/spanner/iam_spanner_instance.go":     "third_party/terraform/services/spanner/iam_spanner_instance.go",
		"converters/google/resources/services/spanner/spanner_instance_iam.go":     "third_party/tgc/spanner_instance_iam.go",
		"converters/google/resources/storage_bucket_iam.go":                        "third_party/tgc/storage_bucket_iam.go",
		"converters/google/resources/organization_policy.go":                       "third_party/tgc/organization_policy.go",
		"converters/google/resources/iam_storage_bucket.go":                        "third_party/tgc/iam_storage_bucket.go",
		"ancestrymanager/ancestrymanager.go":                                       "third_party/tgc/ancestrymanager/ancestrymanager.go",
		"ancestrymanager/ancestrymanager_test.go":                                  "third_party/tgc/ancestrymanager/ancestrymanager_test.go",
		"ancestrymanager/ancestryutil.go":                                          "third_party/tgc/ancestrymanager/ancestryutil.go",
		"ancestrymanager/ancestryutil_test.go":                                     "third_party/tgc/ancestrymanager/ancestryutil_test.go",
		"converters/google/convert.go":                                             "third_party/tgc/convert.go",
		"converters/google/convert_test.go":                                        "third_party/tgc/convert_test.go",
		"converters/google/resources/compute_instance_group.go":                    "third_party/tgc/compute_instance_group.go",
		"converters/google/resources/job.go":                                       "third_party/tgc/job.go",
		"converters/google/resources/service_account_key.go":                       "third_party/tgc/service_account_key.go",
		"converters/google/resources/compute_target_pool.go":                       "third_party/tgc/compute_target_pool.go",
		"converters/google/resources/dataproc_cluster.go":                          "third_party/tgc/dataproc_cluster.go",
		"converters/google/resources/composer_environment.go":                      "third_party/tgc/composer_environment.go",
		"converters/google/resources/commitment.go":                                "third_party/tgc/commitment.go",
		"converters/google/resources/firebase_project.go":                          "third_party/tgc/firebase_project.go",
		"converters/google/resources/appengine_application.go":                     "third_party/tgc/appengine_application.go",
		"converters/google/resources/apikeys_key.go":                               "third_party/tgc/apikeys_key.go",
		"converters/google/resources/logging_folder_bucket_config.go":              "third_party/tgc/logging_folder_bucket_config.go",
		"converters/google/resources/logging_organization_bucket_config.go":        "third_party/tgc/logging_organization_bucket_config.go",
		"converters/google/resources/logging_project_bucket_config.go":             "third_party/tgc/logging_project_bucket_config.go",
		"converters/google/resources/logging_billing_account_bucket_config.go":     "third_party/tgc/logging_billing_account_bucket_config.go",
		"converters/google/resources/appengine_standard_version.go":                "third_party/tgc/appengine_standard_version.go",
	}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversion) CopyFileList(outputFolder string, files map[string]string) {
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

func (tgc TerraformGoogleConversion) replaceImportPath(outputFolder, target string) {
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
