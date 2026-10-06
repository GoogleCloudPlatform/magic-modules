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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

type Terraform struct {
	ResourceCount int

	IAMResourceCount int

	ResourcesForVersion []map[string]string

	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time

	templateFS fs.FS
}

func NewTerraform(product *api.Product, versionName string, startTime time.Time, templateFS fs.FS) Terraform {
	t := Terraform{
		ResourceCount:     0,
		IAMResourceCount:  0,
		Product:           product,
		TargetVersionName: versionName,
		StartTime:         startTime,
		templateFS:        templateFS,
	}

	if product != nil {
		t.Version = *product.VersionObjOrClosest(versionName)
		t.Product.ImportPath = ImportPathFromVersion(versionName)
		for _, r := range t.Product.Objects {
			r.ImportPath = t.Product.ImportPath
		}
	}

	return t
}

func (t Terraform) Generate(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	t.GenerateObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)

	if generateCode {
		t.GenerateProduct(outputFolder)
		t.GenerateOperation(outputFolder)
	}
}

// resourceSem bounds how many resources are generated concurrently across all
// products. Products are already generated in parallel with one another; this
// additionally lets the resources *within* a product proceed in parallel, which
// matters because a single large product (compute) would otherwise be the
// critical path of the whole run.
var resourceSem = make(chan struct{}, runtime.GOMAXPROCS(0))

func (t *Terraform) GenerateObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	var objects []*api.Resource
	for _, object := range t.Product.Objects {
		object.ExcludeIfNotInVersion(t.Product.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			google.LogVerbose("Excluding %s per user request", object.Name)
			continue
		}
		objects = append(objects, object)
	}

	// Each resource writes only its own output files and only reads shared
	// product state, so resources are independent and can be generated concurrently.
	var wg sync.WaitGroup
	for _, object := range objects {
		wg.Add(1)
		go func(object api.Resource) {
			defer wg.Done()
			resourceSem <- struct{}{}
			defer func() { <-resourceSem }()
			t.GenerateObject(object, outputFolder, t.TargetVersionName, generateCode, generateDocs)
		}(*object)
	}
	wg.Wait()
}

func (t *Terraform) GenerateObject(object api.Resource, outputFolder, productPath string, generateCode, generateDocs bool) {
	templateData := NewTemplateData(outputFolder, t.TargetVersionName, t.templateFS)

	if !object.IsExcluded() {
		google.LogVerbose("Generating %s resource", object.Name)
		google.IncrementResourceGenerated()
		t.GenerateResource(object, *templateData, outputFolder, generateCode, generateDocs)
		t.GenerateSingularDataSource(object, *templateData, outputFolder, generateCode, generateDocs)

		if generateCode {
			// log.Printf("Generating %s tests", object.Name)
			t.GenerateResourceTests(object, *templateData, outputFolder)
			t.GenerateResourceSweeper(object, *templateData, outputFolder)
			t.GenerateSingularDataSourceTests(object, *templateData, outputFolder)
			// log.Printf("Generating %s metadata", object.Name)
			t.GenerateResourceMetadata(object, *templateData, outputFolder)
		}
	}

	// if iam_policy is not defined or excluded, don't generate it
	if object.IamPolicy == nil || object.IamPolicy.Exclude {
		return
	}

	t.GenerateIamPolicy(object, *templateData, outputFolder, generateCode, generateDocs)
}

func (t *Terraform) makeFolder(filePath ...string) string {
	targetFolder := path.Join(filePath...)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	return targetFolder
}

func (t *Terraform) GenerateResource(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if generateCode {
		targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
		if object.FrameworkResource {
			fmt.Printf("\n\x1b[1;33mWARNING:\x1b[0m\n")
			fmt.Printf("The plugin framework generation code is considered a WIP and experimental.\nAre you sure you want to use it for %s? (y/n) ", t.ResourceGoFilename(object))

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading input: %v", err)
			}

			if strings.ToLower(strings.TrimSpace(input)) == "y" {
				targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_fw_%s.go", t.ResourceGoFilename(object)))
				templateData.GenerateFWResourceFile(targetFilePath, object)
			} else {
				log.Fatalf("please remove \"plugin_framework_experimental: true\" from the YAML configuration.")
			}
		} else {
			targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s.go", t.ResourceGoFilename(object)))
			templateData.GenerateResourceFile(targetFilePath, object)
		}

		t.GenerateListResource(object, templateData, targetFolder)
	}

	if generateDocs {
		targetFolder := t.makeFolder(outputFolder, "website", "docs", "r")
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s.html.markdown", t.FullResourceName(object)))
		templateData.GenerateDocumentationFile(targetFilePath, object)

		if object.GenerateListResource {
			listDocFolder := t.makeFolder(outputFolder, "website", "docs", "list-resources")
			listDocFilePath := path.Join(listDocFolder, fmt.Sprintf("%s.html.markdown", object.TerraformName()))
			templateData.GenerateListResourceDocumentationFile(listDocFilePath, object)
		}

		if object.IamPolicy.AnyListResource() {
			listDocFolder := t.makeFolder(outputFolder, "website", "docs", "list-resources")
			listDocFilePath := path.Join(listDocFolder, fmt.Sprintf("%s_iam.html.markdown", object.TerraformName()))
			templateData.GenerateIamListResourceDocumentationFile(listDocFilePath, object)
		}
	}
}

func (t *Terraform) GenerateListResource(object api.Resource, templateData TemplateData, targetFolder string) {
	if object.GenerateListResource {
		if object.ExcludeIdentityGeneration {
			log.Fatalf("generate_list_resource requires identity support; remove exclude_identity_generation from resource %q or disable generate_list_resource", object.Name)
		}
		if object.ExcludeRead {
			log.Fatalf("generate_list_resource requires read support; remove exclude_read from resource %q or disable generate_list_resource", object.Name)
		}
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("list_%s.go", t.ResourceGoFilename(object)))
		templateData.GenerateFile(targetFilePath, "templates/terraform/list_resource.go.tmpl", object, true,
			"templates/terraform/list_resource.go.tmpl",
			"templates/terraform/list_resource_method.go.tmpl",
		)

		t.GenerateListResourceQueryTest(object, templateData, targetFolder)
	}
}

// GenerateIamListResource emits list_iam_<resource>.go containing every IAM list
// kind requested via iam_policy.generate_list_resource, mirroring iam_policy.go.tmpl.
func (t *Terraform) GenerateIamListResource(object api.Resource, templateData TemplateData, targetFolder string) {
	if !object.IamPolicy.AnyListResource() {
		return
	}

	targetFilePath := path.Join(targetFolder, fmt.Sprintf("list_iam_%s.go", t.ResourceGoFilename(object)))
	templatePath := "templates/terraform/iam_list_resource.go.tmpl"
	templateData.GenerateFile(targetFilePath, templatePath, object, true, templatePath)
	t.GenerateIamListResourceQueryTest(object, templateData, targetFolder)
}

// GenerateResourceFile is the Bazel counterpart to GenerateResource(), generating *only() the .go file and
// taking the full path to the output file to generate rather than implicitly generating the path.
func (t *Terraform) GenerateResourceFile(object api.Resource, targetFilePath string) {
	if object.FrameworkResource {
		log.Fatalf("Framework resources are currently unsupported")
	}
	targetFolder := path.Dir(targetFilePath)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	templateData := NewTemplateData("", t.TargetVersionName, t.templateFS)
	templateData.GenerateResourceFile(targetFilePath, object)
}

func (t *Terraform) GenerateResourceMetadata(object api.Resource, templateData TemplateData, outputFolder string) {
	targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	target := fmt.Sprintf("resource_%s_generated_meta.yaml", t.FullResourceName(object))
	targetFilePath := path.Join(targetFolder, target)
	templateData.GenerateMetadataFile(targetFilePath, object)
	t.addHashicorpCopyRightHeader(outputFolder, path.Join(t.FolderName(), "services", t.Product.ApiName, target))
}

// GenerateResourceMetadataFile is used by the Bazel version of the MM compiler to generate the specified
// resource's `generated_meta.yaml` file.
func (t *Terraform) GenerateResourceMetadataFile(object api.Resource, targetFilePath string) {
	targetFolder := path.Dir(targetFilePath)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	templateData := NewTemplateData("", t.TargetVersionName, t.templateFS)
	templateData.GenerateMetadataFile(targetFilePath, object)
}

func (t *Terraform) hasEligibleSample(object api.Resource) bool {
	for _, sample := range object.Samples {
		if sample.ExcludeTest {
			continue
		}
		if object.ProductMetadata.VersionObjOrClosest(t.Product.Version.Name).CompareTo(object.ProductMetadata.VersionObjOrClosest(sample.MinVersion)) >= 0 {
			return true
		}
	}
	return false
}

func (t *Terraform) GenerateResourceTests(object api.Resource, templateData TemplateData, outputFolder string) {
	if object.Examples != nil {
		log.Fatalf("Examples block exists in %v", object.Name)
	}

	if !t.hasEligibleSample(object) {
		return
	}

	targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_generated_test.go", t.ResourceGoFilename(object)))
	templateData.GenerateTestFile(targetFilePath, object)
}

func (t *Terraform) GenerateListResourceQueryTest(object api.Resource, templateData TemplateData, targetFolder string) {
	if object.Examples != nil {
		log.Fatalf("Examples block exists in %v", object.Name)
	}
	if object.Samples == nil || object.FirstRunnableTestConfig().Sample == nil {
		return
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("list_%s_generated_test.go", t.ResourceGoFilename(object)))
	templateData.GenerateQueryTestFile(targetFilePath, object)
}

// GenerateIamListResourceQueryTest emits list_iam_<resource>_generated_test.go.
func (t *Terraform) GenerateIamListResourceQueryTest(object api.Resource, templateData TemplateData, targetFolder string) {
	samples := google.Reject(object.Samples, func(s *resource.Sample) bool {
		return s.ExcludeTest
	})
	if len(samples) == 0 {
		log.Printf("[WARNING] No IAM list resource test generated for %s: no non-excluded samples available", object.Name)
		return
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("list_iam_%s_generated_test.go", t.ResourceGoFilename(object)))
	templateData.GenerateIamQueryTestFile(targetFilePath, object)
}

func (t *Terraform) GenerateResourceSweeper(object api.Resource, templateData TemplateData, outputFolder string) {
	if !object.ShouldGenerateSweepers() {
		return
	}

	targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_sweeper.go", t.ResourceGoFilename(object)))
	templateData.GenerateSweeperFile(targetFilePath, object)
}

// GenerateResourceMetadataFile is used by the Bazel version of the MM compiler to generate the sweeper for
// the specified resource. It panics if the resource does not use a sweeper.
func (t *Terraform) GenerateResourceSweeperFile(object api.Resource, targetFilePath string) {
	if !object.ShouldGenerateSweepers() {
		log.Fatalf("attempting to generate a sweeper for unswept resource %q", object.Name)
	}
	targetFolder := path.Dir(targetFilePath)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	templateData := NewTemplateData("", t.TargetVersionName, t.templateFS)
	templateData.GenerateSweeperFile(targetFilePath, object)
}

func (t *Terraform) GenerateSingularDataSource(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if !object.ShouldGenerateSingularDataSource() {
		return
	}

	if generateCode {
		targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("data_source_%s.go", t.ResourceGoFilename(object)))
		templateData.GenerateDataSourceFile(targetFilePath, object)
	}

	if generateDocs {
		targetFolder := t.makeFolder(outputFolder, "website", "docs", "d")
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s.html.markdown", t.FullResourceName(object)))
		templateData.GenerateDataSourceDocumentationFile(targetFilePath, object)
	}
}

func (t *Terraform) GenerateSingularDataSourceTests(object api.Resource, templateData TemplateData, outputFolder string) {
	if object.Examples != nil {
		log.Fatalf("Examples block exists in %v", object.Name)
	}

	if !object.ShouldGenerateSingularDataSourceTests() {
		return
	}

	targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("data_source_%s_test.go", t.ResourceGoFilename(object)))
	templateData.GenerateDataSourceTestFile(targetFilePath, object)

}

// GenerateProduct creates the product.go file for a given service directory.
// This will be used to seed the directory and add a package-level comment
// specific to the product.
func (t *Terraform) GenerateProduct(outputFolder string) {
	targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	targetFilePath := path.Join(targetFolder, "product.go")
	templateData := NewTemplateData(outputFolder, t.TargetVersionName, t.templateFS)
	templateData.GenerateProductFile(targetFilePath, *t.Product)
}

// GenerateProduct creates the product.go file for the bazel version of the MM compiler.
func (t *Terraform) GenerateProductFile(targetFilePath string) {
	targetFolder := path.Dir(targetFilePath)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}

	templateData := NewTemplateData("", t.TargetVersionName, t.templateFS)
	templateData.GenerateProductFile(targetFilePath, *t.Product)
}

func (t *Terraform) GenerateOperation(outputFolder string) {
	asyncObjects := google.Select(t.Product.Objects, func(o *api.Resource) bool {
		return o.AutogenAsync
	})

	if len(asyncObjects) == 0 {
		return
	}

	targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_operation.go", google.Underscore(t.Product.Name)))
	templateData := NewTemplateData(outputFolder, t.TargetVersionName, t.templateFS)
	templateData.GenerateOperationFile(targetFilePath, *asyncObjects[0])
}

// GenerateProduct creates the operation.go file for the bazel version of the MM compiler.
func (t *Terraform) GenerateOperationFile(object api.Resource, targetFilePath string) {
	targetFolder := path.Dir(targetFilePath)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	templateData := NewTemplateData("", t.TargetVersionName, t.templateFS)
	templateData.GenerateOperationFile(targetFilePath, object)
}

func (t *Terraform) GenerateIamPolicy(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if object.Examples != nil {
		log.Fatalf("Examples block exists in %v", object.Name)
	}

	if object.IamPolicy.SampleConfigBody == "" {
		object.IamPolicy.SampleConfigBody = "templates/terraform/iam/iam_attributes.go.tmpl"
	}

	if generateCode && object.IamPolicy != nil && (object.IamPolicy.MinVersion == "" || slices.Index(product.ORDER, object.IamPolicy.MinVersion) <= slices.Index(product.ORDER, t.TargetVersionName)) {
		targetFolder := t.makeFolder(outputFolder, t.FolderName(), "services", t.Product.ApiName)
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("iam_%s.go", t.ResourceGoFilename(object)))
		templateData.GenerateIamPolicyFile(targetFilePath, object)

		// Iam list resource (terraform query support)
		t.GenerateIamListResource(object, templateData, targetFolder)

		// Only generate test if testable example configs exist.
		samples := google.Reject(object.Samples, func(s *resource.Sample) bool {
			return s.ExcludeTest
		})
		if len(samples) != 0 {
			targetFilePath := path.Join(targetFolder, fmt.Sprintf("iam_%s_generated_test.go", t.ResourceGoFilename(object)))
			templateData.GenerateIamPolicyTestFile(targetFilePath, object)
		}
	}
	if generateDocs {
		t.GenerateIamDocumentation(object, templateData, outputFolder, generateCode, generateDocs)
	}
}

func (t *Terraform) GenerateIamDocumentation(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	resourceDocFolder := t.makeFolder(outputFolder, "website", "docs", "r")
	targetFilePath := path.Join(resourceDocFolder, fmt.Sprintf("%s_iam.html.markdown", t.FullResourceName(object)))
	templateData.GenerateIamResourceDocumentationFile(targetFilePath, object)

	datasourceDocFolder := t.makeFolder(outputFolder, "website", "docs", "d")
	targetFilePath = path.Join(datasourceDocFolder, fmt.Sprintf("%s_iam_policy.html.markdown", t.FullResourceName(object)))
	templateData.GenerateIamDatasourceDocumentationFile(targetFilePath, object)
}

// Finds the folder name for a given version of the terraform provider
func (t *Terraform) FolderName() string {
	if t.TargetVersionName == "ga" {
		return "google"
	}
	return "google-" + t.TargetVersionName
}

// Similar to FullResourceName, but override-aware to prevent things like ending in _test.
// Non-Go files should just use FullResourceName.
func (t *Terraform) ResourceGoFilename(object api.Resource) string {
	// early exit if no override is set
	if object.FilenameOverride == "" {
		return t.FullResourceName(object)
	}

	resName := object.FilenameOverride

	var productName string
	if t.Product.LegacyName != "" {
		productName = t.Product.LegacyName
	} else {
		productName = google.Underscore(t.Product.Name)
	}

	return fmt.Sprintf("%s_%s", productName, resName)
}

func (t *Terraform) FullResourceName(object api.Resource) string {
	// early exit- resource-level legacy names override the product too
	if object.LegacyName != "" {
		return strings.Replace(object.LegacyName, "google_", "", 1)
	}

	var productName string
	if t.Product.LegacyName != "" {
		productName = t.Product.LegacyName
	} else {
		productName = google.Underscore(t.Product.Name)
	}

	return fmt.Sprintf("%s_%s", productName, google.Underscore(object.Name))
}

func (t Terraform) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	google.LogVerbose("Copying common files for %s", ProviderName(t))

	files := t.getCommonCopyFiles(t.TargetVersionName, generateCode, generateDocs)
	t.CopyFileList(outputFolder, files, generateCode)
}

// To copy a new folder, add the folder to foldersCopiedToRootDir or foldersCopiedToGoogleDir.
// To copy a file, add the file to singleFiles
func (t Terraform) getCommonCopyFiles(versionName string, generateCode, generateDocs bool) map[string]string {
	// key is the target file and value is the source file
	commonCopyFiles := make(map[string]string, 0)

	// Case 0: If we're generating a specific product, only copy files for that product.
	if t.Product != nil {
		if !generateCode {
			return commonCopyFiles
		}
		googleDir := "google"
		if versionName != "ga" {
			googleDir = fmt.Sprintf("google-%s", versionName)
		}
		files := t.getCopyFilesInFolder("third_party/terraform/services/"+t.Product.ApiName, googleDir)
		maps.Copy(commonCopyFiles, files)
		return commonCopyFiles
	}

	// Case 1: When copy all of files except .tmpl in a folder to the root directory of downstream repository,
	// save the folder name to foldersCopiedToRootDir
	foldersCopiedToRootDir := []string{"third_party/terraform/META.d", "third_party/terraform/version"}
	// Copy TeamCity-related Kotlin & Markdown files to TPG only, not TPGB
	if versionName == "ga" {
		foldersCopiedToRootDir = append(foldersCopiedToRootDir, "third_party/terraform/.teamcity")
	}
	if generateCode {
		foldersCopiedToRootDir = append(foldersCopiedToRootDir, "third_party/terraform/scripts")
	}
	if generateDocs {
		foldersCopiedToRootDir = append(foldersCopiedToRootDir, "third_party/terraform/website")
	}
	for _, folder := range foldersCopiedToRootDir {
		files := t.getCopyFilesInFolder(folder, ".")
		maps.Copy(commonCopyFiles, files)
	}

	// Case 2: When copy all of files except .tmpl in a folder to the google directory of downstream repository,
	// save the folder name to foldersCopiedToGoogleDir
	var foldersCopiedToGoogleDir []string
	if generateCode {
		foldersCopiedToGoogleDir = []string{
			"third_party/terraform/acctest",
			"third_party/terraform/allservices",
			"third_party/terraform/envvar",
			"third_party/terraform/functions",
			"third_party/terraform/fwmodels",
			"third_party/terraform/fwprovider",
			"third_party/terraform/fwresource",
			"third_party/terraform/fwtransport",
			"third_party/terraform/fwutils",
			"third_party/terraform/fwvalidators",
			"third_party/terraform/provider",
			"third_party/terraform/registry",
			"third_party/terraform/sweeper",
			"third_party/terraform/test-fixtures",
			"third_party/terraform/tpgdclresource",
			"third_party/terraform/tpgiamresource",
			"third_party/terraform/tpgresource",
			"third_party/terraform/transport",
			"third_party/terraform/verify",
		}
	}
	googleDir := "google"
	if versionName != "ga" {
		googleDir = fmt.Sprintf("google-%s", versionName)
	}
	// Copy files to google(or google-beta or google-private) folder in downstream
	for _, folder := range foldersCopiedToGoogleDir {
		files := t.getCopyFilesInFolder(folder, googleDir)
		maps.Copy(commonCopyFiles, files)
	}

	// Case 3: When copy a single file, save the target as key and source as value to the map singleFiles
	singleFiles := map[string]string{
		"go.sum":                           "third_party/terraform/go.sum",
		"go.mod":                           "third_party/terraform/go.mod",
		".go-version":                      "third_party/terraform/.go-version",
		"terraform-registry-manifest.json": "third_party/terraform/terraform-registry-manifest.json",
	}
	maps.Copy(commonCopyFiles, singleFiles)

	return commonCopyFiles
}

func (t Terraform) getCopyFilesInFolder(folderPath, targetDir string) map[string]string {
	m := make(map[string]string, 0)
	fs.WalkDir(t.templateFS, folderPath, func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !di.IsDir() && !strings.HasSuffix(di.Name(), ".tmpl") && !strings.HasSuffix(di.Name(), ".erb") { // Exception files
			if di.Name() == "gha-branch-renaming.png" || di.Name() == "clock-timings-of-branch-making-and-usage.png" {
				return nil
			}

			fname := strings.TrimPrefix(path, "third_party/terraform/")
			target := fname
			if targetDir != "." {
				target = fmt.Sprintf("%s/%s", targetDir, fname)
			}
			m[target] = path
		}
		return nil
	})

	return m
}

func (t Terraform) CopyFileList(outputFolder string, files map[string]string, generateCode bool) {
	// Files are independent of one another, so process them concurrently.
	var wg sync.WaitGroup
	for target, source := range files {
		wg.Add(1)
		go func(target, source string) {
			defer wg.Done()
			resourceSem <- struct{}{}
			defer func() { <-resourceSem }()
			t.copyFile(outputFolder, target, source, generateCode)
		}(target, source)
	}
	wg.Wait()
}

// copyFile copies a single handwritten file from the template FS to the output
// folder, applying the version-specific import path rewrite and header
// insertions in memory so the file is formatted at most once and written once.
func (t Terraform) copyFile(outputFolder, target, source string, generateCode bool) {
	targetFile := filepath.Join(outputFolder, target)
	targetDir := filepath.Dir(targetFile)

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", targetDir, err))
	}
	// If we've modified a file since starting an MM run, it's a reasonable
	// assumption that it was this run that modified it.
	if info, err := os.Stat(targetFile); !errors.Is(err, os.ErrNotExist) && t.StartTime.Before(info.ModTime()) {
		log.Fatalf("%s was already modified during this run at %s", targetFile, info.ModTime().String())
	}

	sourceByte, err := fs.ReadFile(t.templateFS, source)
	if err != nil {
		log.Fatalf("Cannot read source file %s while copying: %s", source, err)
	}

	var permission fs.FileMode
	if strings.HasSuffix(targetDir, "scripts") {
		permission = 0755
	} else {
		permission = 0644
	}

	ext := filepath.Ext(target)
	needsFormat := false

	// Replace import path based on version (beta/alpha)
	if ext == ".go" || (ext == ".mod" && generateCode) {
		var replaced bool
		sourceByte, replaced = t.replaceImportPathBytes(target, sourceByte)
		if replaced && ext == ".go" {
			needsFormat = true
		}
	}
	if ext == ".go" || ext == ".markdown" {
		var added bool
		sourceByte, added = copyfileHeaderBytes(source, target, sourceByte)
		if added && ext == ".go" {
			needsFormat = true
		}
	}
	if needsFormat {
		sourceByte = formatGoSource(targetFile, sourceByte)
	}
	if ext == ".go" || ext == ".yaml" {
		sourceByte = t.hashicorpCopyRightHeaderBytes(outputFolder, target, sourceByte)
	}

	err = os.WriteFile(targetFile, sourceByte, permission)
	if err != nil {
		log.Fatalf("Cannot write target file %s while copying: %s", target, err)
	}
}

// Compiles files that are shared at the provider level
func (t Terraform) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	google.LogVerbose("Generating common files for %s", ProviderName(t))
	if t.Product == nil {
		t.generateResourcesForVersion(products)
	}
	files := t.getCommonCompileFiles(t.TargetVersionName)
	templateData := NewTemplateData(outputFolder, t.TargetVersionName, t.templateFS)
	t.CompileFileList(outputFolder, files, *templateData, products)
}

// To compile a new folder, add the folder to foldersCompiledToRootDir or foldersCompiledToGoogleDir.
// To compile a file, add the file to singleFiles
func (t Terraform) getCommonCompileFiles(versionName string) map[string]string {
	// key is the target file and the value is the source file
	commonCompileFiles := make(map[string]string, 0)

	if t.Product != nil {
		googleDir := "google"
		if versionName != "ga" {
			googleDir = fmt.Sprintf("google-%s", versionName)
		}
		return t.getCompileFilesInFolder("third_party/terraform/services/"+t.Product.ApiName, googleDir)
	}

	// Case 1: When compile all of files except .tmpl in a folder to the root directory of downstream repository,
	// save the folder name to foldersCopiedToRootDir
	foldersCompiledToRootDir := []string{"third_party/terraform/scripts"}
	for _, folder := range foldersCompiledToRootDir {
		files := t.getCompileFilesInFolder(folder, ".")
		maps.Copy(commonCompileFiles, files)
	}

	// Case 2: When compile all of files except .tmpl in a folder to the google directory of downstream repository,
	// save the folder name to foldersCopiedToGoogleDir
	foldersCompiledToGoogleDir := []string{
		"third_party/terraform/acctest",
		"third_party/terraform/allservices",
		"third_party/terraform/envvar",
		"third_party/terraform/functions",
		"third_party/terraform/fwmodels",
		"third_party/terraform/fwprovider",
		"third_party/terraform/fwresource",
		"third_party/terraform/fwtransport",
		"third_party/terraform/provider",
		"third_party/terraform/sweeper",
		"third_party/terraform/test-fixtures",
		"third_party/terraform/tpgdclresource",
		"third_party/terraform/tpgiamresource",
		"third_party/terraform/tpgresource",
		"third_party/terraform/transport",
		"third_party/terraform/verify",
	}
	googleDir := "google"
	if versionName != "ga" {
		googleDir = fmt.Sprintf("google-%s", versionName)
	}
	for _, folder := range foldersCompiledToGoogleDir {
		files := t.getCompileFilesInFolder(folder, googleDir)
		maps.Copy(commonCompileFiles, files)
	}

	// Case 3: When compile a single file, save the target as key and source as value to the map singleFiles
	singleFiles := map[string]string{
		"main.go":                       "third_party/terraform/main.go.tmpl",
		".goreleaser.yml":               "third_party/terraform/.goreleaser.yml.tmpl",
		".release/release-metadata.hcl": "third_party/terraform/release-metadata.hcl.tmpl",
		".copywrite.hcl":                "third_party/terraform/.copywrite.hcl.tmpl",
	}
	maps.Copy(commonCompileFiles, singleFiles)

	return commonCompileFiles
}

func (t Terraform) getCompileFilesInFolder(folderPath, targetDir string) map[string]string {
	m := make(map[string]string, 0)
	fs.WalkDir(t.templateFS, folderPath, func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !di.IsDir() && strings.HasSuffix(di.Name(), ".tmpl") {
			fname := strings.TrimPrefix(path, "third_party/terraform/")
			fname = strings.TrimSuffix(fname, ".tmpl")
			target := fname
			if targetDir != "." {
				target = fmt.Sprintf("%s/%s", targetDir, fname)
			}
			m[target] = path
		}
		return nil
	})

	return m
}

func (t Terraform) CompileFileList(outputFolder string, files map[string]string, fileTemplate TemplateData, products []*api.Product) {
	providerWithProducts := ProviderWithProducts{
		Terraform: t,
		Products:  products,
	}

	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	// Files are independent of one another, so process them concurrently.
	var wg sync.WaitGroup
	for target, source := range files {
		wg.Add(1)
		go func(target, source string) {
			defer wg.Done()
			resourceSem <- struct{}{}
			defer func() { <-resourceSem }()
			t.compileFile(outputFolder, target, source, fileTemplate, providerWithProducts)
		}(target, source)
	}
	wg.Wait()
}

// compileFile renders a single handwritten template to the output folder,
// applying the version-specific import path rewrite and header insertions in
// memory so the file is formatted at most once and written once.
func (t Terraform) compileFile(outputFolder, target, source string, fileTemplate TemplateData, input any) {
	targetFile := filepath.Join(outputFolder, target)
	targetDir := filepath.Dir(targetFile)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", targetDir, err))
	}

	sourceByte := fileTemplate.renderFile(targetFile, source, input, source)
	// skip this file if no file was generated
	if sourceByte == nil {
		return
	}

	ext := filepath.Ext(targetFile)
	if ext == ".go" {
		// Format before inserting the header: rendered templates frequently start
		// with blank lines, and whether gofmt treats the header as the package
		// doc comment (re-indenting it) depends on those being stripped first.
		sourceByte = formatGoSource(targetFile, sourceByte)
	}
	sourceByte, _ = t.replaceImportPathBytes(target, sourceByte)
	if ext == ".go" || ext == ".markdown" {
		sourceByte, _ = copyfileHeaderBytes(source, target, sourceByte)
	}
	if ext == ".go" {
		sourceByte = formatGoSource(targetFile, sourceByte)
	}
	sourceByte = t.hashicorpCopyRightHeaderBytes(outputFolder, target, sourceByte)

	if err := os.WriteFile(targetFile, sourceByte, 0644); err != nil {
		log.Fatalf("Cannot write file %s: %s", targetFile, err)
	}
}

// copyfileHeaderBytes prepends the "AUTO GENERATED CODE / Type: Handwritten"
// banner to content, unless the banner is already present. It returns the
// (possibly) updated content and whether the banner was added. Callers are
// responsible for running gofmt on Go sources afterwards.
func copyfileHeaderBytes(srcpath, target string, content []byte) ([]byte, bool) {
	githubPrefix := "https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/"
	if !strings.HasPrefix(srcpath, githubPrefix) {
		srcpath = githubPrefix + srcpath
	}

	srcStr := string(content)
	if strings.Contains(srcStr, "***     AUTO GENERATED CODE    ***    Type: Handwritten     ***") {
		return content, false
	}

	templateFormat := `// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: Handwritten     ***
//
// ----------------------------------------------------------------------------
//
//     This code is generated by Magic Modules using the following:
//
//     Source file: %s
//
//     DO NOT EDIT this file directly. Any changes made to this file will be
//     overwritten during the next generation cycle.
//
// ----------------------------------------------------------------------------
%s`
	body := srcStr
	if filepath.Ext(target) == ".markdown" {
		// insert the header after ---
		templateFormat = "---\n" + strings.Replace(templateFormat, "//", "#", -1)
		body = strings.TrimPrefix(srcStr, "---\n")
	}

	return []byte(fmt.Sprintf(templateFormat, srcpath, body)), true
}

// addHashicorpCopyRightHeader applies hashicorpCopyRightHeaderBytes to a file on disk.
func (t Terraform) addHashicorpCopyRightHeader(outputFolder, target string) {
	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to add Hashicorp copy right: %s", targetFile, err)
	}

	updated := t.hashicorpCopyRightHeaderBytes(outputFolder, target, sourceByte)
	if bytes.Equal(updated, sourceByte) {
		return
	}

	err = os.WriteFile(targetFile, updated, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to add Hashicorp copy right: %s", target, err)
	}
}

// hashicorpCopyRightHeaderBytes prepends the HashiCorp copyright header to
// content when target is a file that should carry one, and returns the result.
// Content that already carries the header is returned unchanged.
func (t Terraform) hashicorpCopyRightHeaderBytes(outputFolder, target string, content []byte) []byte {
	if !t.shouldAddHashicorpCopyRightHeader(outputFolder, target) {
		return content
	}

	lang := languageFromFilename(target)

	// File is not ignored and is appropriate file type to add header to
	copyrightHeader := []string{"Copyright IBM Corp. 2014, 2026", "SPDX-License-Identifier: MPL-2.0"}
	header := commentBlock(copyrightHeader, lang)

	if bytes.Contains(content, []byte("Copyright IBM Corp. 2014, 2026")) &&
		bytes.Contains(content, []byte("SPDX-License-Identifier: MPL-2.0")) {
		return content
	}

	return google.Concat([]byte(header), content)
}

// shouldAddHashicorpCopyRightHeader reports whether target (relative to
// outputFolder) is a file that should carry the HashiCorp copyright header.
func (t Terraform) shouldAddHashicorpCopyRightHeader(outputFolder, target string) bool {
	if !expectedOutputFolder(outputFolder) {
		log.Printf("Unexpected output folder (%s) detected "+
			"when deciding to add HashiCorp copyright headers.\n"+
			"Watch out for unexpected changes to copied files", outputFolder)
	}
	// only add copyright headers when generating TPG, TPGB, and TPGN
	if !(strings.HasSuffix(outputFolder, "terraform-provider-google") || strings.HasSuffix(outputFolder, "terraform-provider-google-beta") || strings.HasSuffix(outputFolder, "terraform-provider-google-nightly")) {
		return false
	}

	// Prevent adding copyright header to files with paths or names matching the strings below
	// NOTE: these entries need to match the content of the .copywrite.hcl file originally
	//       created in https://github.com/GoogleCloudPlatform/magic-modules/pull/7336
	//       The test-fixtures folder is not included here as it's copied as a whole,
	//       not file by file
	ignoredFolders := []string{".release/", ".changelog/", "examples/", "scripts/"}
	ignoredFiles := []string{"go.mod", ".goreleaser.yml", ".golangci.yml", "terraform-registry-manifest.json"}
	for _, folder := range ignoredFolders {
		// folder will be path leading to file
		if strings.HasPrefix(target, folder) {
			return false
		}
	}

	for _, file := range ignoredFiles {
		// file will be the filename and extension, with no preceding path
		if strings.HasSuffix(target, file) {
			return false
		}
	}

	// Some file types we don't want to add headers to
	// e.g. .sh where headers are functional
	// Also, this guards against new filetypes being added and triggering build errors
	return languageFromFilename(target) != "unsupported"
}

func expectedOutputFolder(outputFolder string) bool {
	expectedFolders := []string{"terraform-provider-google", "terraform-provider-google-beta", "terraform-provider-google-nightly", "terraform-next", "terraform-provider-google-internal", "terraform-google-conversion", "tfplan2cai"}
	folderName := filepath.Base(outputFolder) // Possible issue with Windows OS
	isExpected := false
	for _, folder := range expectedFolders {
		if folderName == folder {
			isExpected = true
			break
		}
	}

	return isExpected
}

// replaceImportPathBytes rewrites GA provider import paths in content to the
// paths for the target version. It returns the updated content and whether a
// rewrite was performed (always false for the GA provider). Callers are
// responsible for running gofmt on Go sources afterwards. It exits if content
// imports from the beta module directly, which is never allowed.
func (t Terraform) replaceImportPathBytes(target string, content []byte) ([]byte, bool) {
	gaImportPath := ImportPathFromVersion("ga")
	betaImportPath := ImportPathFromVersion("beta")

	if bytes.Contains(content, []byte(betaImportPath)) {
		log.Fatalf("Importing a package from module %s is not allowed in file %s. Please import a package from module %s.", betaImportPath, filepath.Base(target), gaImportPath)
	}

	if t.TargetVersionName == "ga" {
		return content, false
	}

	// Replace the import pathes in utility files
	var tpg, dir string
	switch t.TargetVersionName {
	case "beta":
		tpg = TERRAFORM_PROVIDER_BETA
		dir = RESOURCE_DIRECTORY_BETA
	default:
		tpg = "github.com/hashicorp/terraform-provider-google-" + t.TargetVersionName
		dir = "google-" + t.TargetVersionName
	}

	content = bytes.Replace(content, []byte(gaImportPath), []byte(tpg+"/"+dir), -1)
	content = bytes.Replace(content, []byte(TERRAFORM_PROVIDER_GA+"/version"), []byte(tpg+"/version"), -1)
	content = bytes.Replace(content, []byte("module "+TERRAFORM_PROVIDER_GA), []byte("module "+tpg), -1)
	return content, true
}

func (t Terraform) ProviderFromVersion() string {
	var dir string
	switch t.TargetVersionName {
	case "ga":
		dir = RESOURCE_DIRECTORY_GA
	case "beta":
		dir = RESOURCE_DIRECTORY_BETA
	default:
		dir = "google-" + t.TargetVersionName
	}
	return dir
}

// Gets the list of services dependent on the version ga, beta, and private
// If there are some resources of a servcie is in GA,
// then this service is in GA. Otherwise, the service is in BETA
func (t Terraform) GetMmv1ServicesInVersion(products []*api.Product) []string {
	var services []string
	for _, product := range products {
		if t.TargetVersionName == "ga" {
			someResourceInGA := false
			for _, object := range product.Objects {
				if someResourceInGA {
					break
				}

				if !object.Exclude && !object.NotInVersion(product.VersionObjOrClosest(t.TargetVersionName)) {
					someResourceInGA = true
				}
			}

			if someResourceInGA {
				services = append(services, strings.ToLower(product.Name))
			}
		} else {
			services = append(services, strings.ToLower(product.Name))
		}
	}
	return services
}

// # Generates the list of resources, and gets the count of resources and iam resources
// # dependent on the version ga, beta or private.
// # The resource object has the format
// # {
// #    terraform_name:
// #    resource_name:
// #    iam_class_name:
// # }
// # The variable resources_for_version is used to generate resources in file
// # mmv1/third_party/terraform/provider/provider_mmv1_resources.go.erb
func (t *Terraform) generateResourcesForVersion(products []*api.Product) {
	for _, productDefinition := range products {
		service := strings.ToLower(productDefinition.Name)
		for _, object := range productDefinition.Objects {
			if object.Exclude || object.NotInVersion(productDefinition.VersionObjOrClosest(t.TargetVersionName)) {
				continue
			}

			var resourceName string

			if !object.IsExcluded() {
				t.ResourceCount++
				resourceName = fmt.Sprintf("%s.Resource%s", service, object.ResourceName())
			}

			var iamClassName string
			iamPolicy := object.IamPolicy
			if iamPolicy != nil && !iamPolicy.Exclude {
				t.IAMResourceCount += 3

				if slices.Index(product.ORDER, iamPolicy.MinVersion) <= slices.Index(product.ORDER, t.TargetVersionName) {
					iamClassName = fmt.Sprintf("%s.%s", service, object.ResourceName())
				}
			}

			t.ResourcesForVersion = append(t.ResourcesForVersion, map[string]string{
				"TerraformName": object.TerraformName(),
				"ResourceName":  resourceName,
				"IamClassName":  iamClassName,
			})
		}
	}
}

// # Adapted from the method used in templating
// # See: mmv1/compile/core.rb
func commentBlock(text []string, lang string) string {
	var headers []string
	switch lang {
	case "python", "yaml":
		headers = commentText(text, "#")
	case "go":
		headers = commentText(text, "//")
	default:
		log.Fatalf("Unknown language for comment: %s", lang)
	}

	headerString := strings.Join(headers, "\n")
	return fmt.Sprintf("%s\n", headerString) // add trailing newline to returned value
}

func commentText(text []string, symbols string) []string {
	var header []string
	for _, t := range text {
		var comment string
		if t == "" {
			comment = symbols
		} else {
			comment = fmt.Sprintf("%s %s", symbols, t)
		}
		header = append(header, comment)
	}
	return header
}

func languageFromFilename(filename string) string {
	switch extension := filepath.Ext(filename); extension {
	case ".go":
		return "go"
	case ".rb":
		return "rb"
	case ".yaml", ".yml":
		return "yaml"
	default:
		return "unsupported"
	}
}

// Returns the extension for DCL packages for the given version. This is needed
// as the DCL uses "alpha" for preview resources, while we use "private"
func (t Terraform) DCLVersion() string {
	switch t.TargetVersionName {
	case "beta", "nightly":
		return "/beta"
	case "private", "internal":
		return "/alpha"
	default:
		return ""
	}
}

// Gets the provider versions supported by a version
func (t Terraform) SupportedProviderVersions() []string {
	var supported []string
	for i, v := range product.ORDER {
		if i == 0 {
			continue
		}
		if i > slices.Index(product.ORDER, t.TargetVersionName) {
			break
		}
		supported = append(supported, v)
	}
	return supported
}

type ProviderWithProducts struct {
	Terraform
	Compiler string
	Products []*api.Product
}
