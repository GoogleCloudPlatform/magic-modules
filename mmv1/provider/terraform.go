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
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/fs"
	"log"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
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
}

func NewTerraform(product *api.Product, versionName string, startTime time.Time) Terraform {
	t := Terraform{
		ResourceCount:     0,
		IAMResourceCount:  0,
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
	}

	t.Product.SetPropertiesBasedOnVersion(&t.Version)
	t.Product.SetCompiler(ProviderName(t))
	for _, r := range t.Product.Objects {
		r.SetCompiler(ProviderName(t))
		r.ImportPath = ImportPathFromVersion(versionName)
	}

	return t
}

func (t Terraform) Generate(outputFolder, productPath, resourceToGenerate string, generateCode, generateDocs bool) {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	t.GenerateObjects(outputFolder, resourceToGenerate, generateCode, generateDocs)

	if generateCode {
		t.GenerateProduct(outputFolder)
		t.GenerateOperation(outputFolder)
	}
}

func (t *Terraform) GenerateObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	for _, object := range t.Product.Objects {
		object.ExcludeIfNotInVersion(&t.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			log.Printf("Excluding %s per user request", object.Name)
			continue
		}

		t.GenerateObject(*object, outputFolder, t.TargetVersionName, generateCode, generateDocs)
	}
}

func (t *Terraform) GenerateObject(object api.Resource, outputFolder, productPath string, generateCode, generateDocs bool) {
	templateData := NewTemplateData(outputFolder, t.TargetVersionName)

	if !object.IsExcluded() {
		log.Printf("Generating %s resource", object.Name)
		t.GenerateResource(object, *templateData, outputFolder, generateCode, generateDocs)

		if generateCode {
			// log.Printf("Generating %s tests", object.Name)
			t.GenerateResourceTests(object, *templateData, outputFolder)
			t.GenerateResourceSweeper(object, *templateData, outputFolder)
			t.GenerateSingularDataSource(object, *templateData, outputFolder)
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

func (t *Terraform) GenerateResource(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if generateCode {
		productName := t.Product.ApiName
		targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
		}
		if object.FrameworkResource {
			targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_fw_%s.go", t.ResourceGoFilename(object)))
			templateData.GenerateFWResourceFile(targetFilePath, object)
		} else {
			targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s.go", t.ResourceGoFilename(object)))
			templateData.GenerateResourceFile(targetFilePath, object)
		}
	}

	if generateDocs {
		targetFolder := path.Join(outputFolder, "website", "docs", "r")
		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
		}
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s.html.markdown", t.FullResourceName(object)))
		templateData.GenerateDocumentationFile(targetFilePath, object)
	}
}

func (t *Terraform) GenerateResourceMetadata(object api.Resource, templateData TemplateData, outputFolder string) {
	productName := t.Product.ApiName
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_generated_meta.yaml", t.FullResourceName(object)))
	templateData.GenerateMetadataFile(targetFilePath, object)
}

func (t *Terraform) GenerateResourceTests(object api.Resource, templateData TemplateData, outputFolder string) {
	eligibleExample := false
	for _, example := range object.Examples {
		if !example.ExcludeTest {
			if object.ProductMetadata.VersionObjOrClosest(t.Version.Name).CompareTo(object.ProductMetadata.VersionObjOrClosest(example.MinVersion)) >= 0 {
				eligibleExample = true
				break
			}
		}
	}
	if !eligibleExample {
		return
	}

	productName := t.Product.ApiName
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_generated_test.go", t.ResourceGoFilename(object)))
	templateData.GenerateTestFile(targetFilePath, object)
}

func (t *Terraform) GenerateResourceSweeper(object api.Resource, templateData TemplateData, outputFolder string) {
	if !object.ShouldGenerateSweepers() {
		return
	}

	productName := t.Product.ApiName
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_sweeper.go", t.ResourceGoFilename(object)))
	templateData.GenerateSweeperFile(targetFilePath, object)
}

func (t *Terraform) GenerateSingularDataSource(object api.Resource, templateData TemplateData, outputFolder string) {
	if !object.ShouldGenerateSingularDataSource() {
		return
	}

	productName := t.Product.ApiName
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("data_source_%s.go", t.ResourceGoFilename(object)))
	templateData.GenerateDataSourceFile(targetFilePath, object)
}

func (t *Terraform) GenerateSingularDataSourceTests(object api.Resource, templateData TemplateData, outputFolder string) {
	if !object.ShouldGenerateSingularDataSourceTests() {
		return
	}

	productName := t.Product.ApiName
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("data_source_%s_test.go", t.ResourceGoFilename(object)))
	templateData.GenerateDataSourceTestFile(targetFilePath, object)

}

// GenerateProduct creates the product.go file for a given service directory.
// This will be used to seed the directory and add a package-level comment
// specific to the product.
func (t *Terraform) GenerateProduct(outputFolder string) {
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}

	targetFilePath := path.Join(targetFolder, "product.go")
	templateData := NewTemplateData(outputFolder, t.TargetVersionName)
	templateData.GenerateProductFile(targetFilePath, *t.Product)
}

func (t *Terraform) GenerateOperation(outputFolder string) {
	asyncObjects := google.Select(t.Product.Objects, func(o *api.Resource) bool {
		return o.AutogenAsync
	})

	if len(asyncObjects) == 0 {
		return
	}

	targetFolder := path.Join(outputFolder, t.FolderName(), "services", t.Product.ApiName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_operation.go", google.Underscore(t.Product.Name)))
	templateData := NewTemplateData(outputFolder, t.TargetVersionName)
	templateData.GenerateOperationFile(targetFilePath, *asyncObjects[0])
}

// Generate the IAM policy for this object. This is used to query and test
// IAM policies separately from the resource itself
func (t *Terraform) GenerateIamPolicy(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if generateCode && object.IamPolicy != nil && (object.IamPolicy.MinVersion == "" || slices.Index(product.ORDER, object.IamPolicy.MinVersion) <= slices.Index(product.ORDER, t.TargetVersionName)) {
		productName := t.Product.ApiName
		targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
		}
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("iam_%s.go", t.ResourceGoFilename(object)))
		templateData.GenerateIamPolicyFile(targetFilePath, object)

		// Only generate test if testable examples exist.
		examples := google.Reject(object.Examples, func(e resource.Examples) bool {
			return e.ExcludeTest
		})
		if len(examples) != 0 {
			targetFilePath := path.Join(targetFolder, fmt.Sprintf("iam_%s_generated_test.go", t.ResourceGoFilename(object)))
			templateData.GenerateIamPolicyTestFile(targetFilePath, object)
		}
	}
	if generateDocs {
		t.GenerateIamDocumentation(object, templateData, outputFolder, generateCode, generateDocs)
	}
}

func (t *Terraform) GenerateIamDocumentation(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	resourceDocFolder := path.Join(outputFolder, "website", "docs", "r")
	if err := os.MkdirAll(resourceDocFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", resourceDocFolder, err))
	}
	targetFilePath := path.Join(resourceDocFolder, fmt.Sprintf("%s_iam.html.markdown", t.FullResourceName(object)))
	templateData.GenerateIamResourceDocumentationFile(targetFilePath, object)

	datasourceDocFolder := path.Join(outputFolder, "website", "docs", "d")
	if err := os.MkdirAll(datasourceDocFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", datasourceDocFolder, err))
	}
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
	log.Printf("Copying common files for %s", ProviderName(t))

	files := t.getCommonCopyFiles(t.TargetVersionName, generateCode, generateDocs)
	t.CopyFileList(outputFolder, files, generateCode)
}

// To copy a new folder, add the folder to foldersCopiedToRootDir or foldersCopiedToGoogleDir.
// To copy a file, add the file to singleFiles
func (t Terraform) getCommonCopyFiles(versionName string, generateCode, generateDocs bool) map[string]string {
	// key is the target file and value is the source file
	commonCopyFiles := make(map[string]string, 0)

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
		foldersCopiedToGoogleDir = []string{"third_party/terraform/services", "third_party/terraform/acctest", "third_party/terraform/sweeper", "third_party/terraform/provider", "third_party/terraform/tpgdclresource", "third_party/terraform/tpgiamresource", "third_party/terraform/tpgresource", "third_party/terraform/transport", "third_party/terraform/fwmodels", "third_party/terraform/fwprovider", "third_party/terraform/fwtransport", "third_party/terraform/fwresource", "third_party/terraform/fwutils", "third_party/terraform/fwvalidators", "third_party/terraform/verify", "third_party/terraform/envvar", "third_party/terraform/functions", "third_party/terraform/test-fixtures"}
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
	filepath.WalkDir(folderPath, func(path string, di fs.DirEntry, err error) error {
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
	for target, source := range files {
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

		sourceByte, err := os.ReadFile(source)
		if err != nil {
			log.Fatalf("Cannot read source file %s while copying: %s", source, err)
		}

		var permission fs.FileMode
		if strings.HasSuffix(targetDir, "scripts") {
			permission = 0755
		} else {
			permission = 0644
		}

		err = os.WriteFile(targetFile, sourceByte, permission)
		if err != nil {
			log.Fatalf("Cannot write target file %s while copying: %s", target, err)
		}

		// Replace import path based on version (beta/alpha)
		if filepath.Ext(target) == ".go" || (filepath.Ext(target) == ".mod" && generateCode) {
			t.replaceImportPath(outputFolder, target)
		}
		if filepath.Ext(target) == ".go" || filepath.Ext(target) == ".markdown" {
			t.addCopyfileHeader(source, outputFolder, target)
		}
		if filepath.Ext(target) == ".go" {
			t.addHashicorpCopyRightHeader(outputFolder, target)
		}
	}
}

// Compiles files that are shared at the provider level
func (t Terraform) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	t.generateResourcesForVersion(products)
	files := t.getCommonCompileFiles(t.TargetVersionName)
	templateData := NewTemplateData(outputFolder, t.TargetVersionName)
	t.CompileFileList(outputFolder, files, *templateData, products)
}

// To compile a new folder, add the folder to foldersCompiledToRootDir or foldersCompiledToGoogleDir.
// To compile a file, add the file to singleFiles
func (t Terraform) getCommonCompileFiles(versionName string) map[string]string {
	// key is the target file and the value is the source file
	commonCompileFiles := make(map[string]string, 0)

	// Case 1: When compile all of files except .tmpl in a folder to the root directory of downstream repository,
	// save the folder name to foldersCopiedToRootDir
	foldersCompiledToRootDir := []string{"third_party/terraform/scripts"}
	for _, folder := range foldersCompiledToRootDir {
		files := t.getCompileFilesInFolder(folder, ".")
		maps.Copy(commonCompileFiles, files)
	}

	// Case 2: When compile all of files except .tmpl in a folder to the google directory of downstream repository,
	// save the folder name to foldersCopiedToGoogleDir
	foldersCompiledToGoogleDir := []string{"third_party/terraform/services", "third_party/terraform/acctest", "third_party/terraform/sweeper", "third_party/terraform/provider", "third_party/terraform/tpgdclresource", "third_party/terraform/tpgiamresource", "third_party/terraform/tpgresource", "third_party/terraform/transport", "third_party/terraform/fwmodels", "third_party/terraform/fwprovider", "third_party/terraform/fwtransport", "third_party/terraform/fwresource", "third_party/terraform/verify", "third_party/terraform/envvar", "third_party/terraform/functions", "third_party/terraform/test-fixtures"}
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
	filepath.WalkDir(folderPath, func(path string, di fs.DirEntry, err error) error {
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

		fileTemplate.GenerateFile(targetFile, source, providerWithProducts, formatFile, templates...)
		// continue to next file if no file was generated
		if _, err := os.Stat(targetFile); errors.Is(err, os.ErrNotExist) {
			continue
		}
		t.replaceImportPath(outputFolder, target)
		if filepath.Ext(targetFile) == ".go" || filepath.Ext(targetFile) == ".markdown" {
			t.addCopyfileHeader(source, outputFolder, target)
		}
		t.addHashicorpCopyRightHeader(outputFolder, target)
	}
}

func (t Terraform) addCopyfileHeader(srcpath, outputFolder, target string) {
	githubPrefix := "https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/"
	if !strings.HasPrefix(srcpath, githubPrefix) {
		srcpath = githubPrefix + srcpath
	}

	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to add copy file header: %s", targetFile, err)
	}

	srcStr := string(sourceByte)
	if strings.Contains(srcStr, "***     AUTO GENERATED CODE    ***    Type: Handwritten     ***") {
		return
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
	content := srcStr
	if filepath.Ext(target) == ".markdown" {
		// insert the header after ---
		templateFormat = "---\n" + strings.Replace(templateFormat, "//", "#", -1)
		content = strings.TrimPrefix(srcStr, "---\n")
	}

	fileStr := fmt.Sprintf(templateFormat, srcpath, content)

	sourceByte = []byte(fileStr)
	// format go file
	if filepath.Ext(targetFile) == ".go" {
		sourceByte, err = format.Source(sourceByte)
		if err != nil {
			log.Printf("error formatting %s: %s\n", targetFile, err)
			return
		}
	}

	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to add copy file header: %s", target, err)
	}
}

func (t Terraform) addHashicorpCopyRightHeader(outputFolder, target string) {
	if !expectedOutputFolder(outputFolder) {
		log.Printf("Unexpected output folder (%s) detected "+
			"when deciding to add HashiCorp copyright headers.\n"+
			"Watch out for unexpected changes to copied files", outputFolder)
	}
	// only add copyright headers when generating TPG and TPGB
	if !(strings.HasSuffix(outputFolder, "terraform-provider-google") || strings.HasSuffix(outputFolder, "terraform-provider-google-beta")) {
		return
	}

	// Prevent adding copyright header to files with paths or names matching the strings below
	// NOTE: these entries need to match the content of the .copywrite.hcl file originally
	//       created in https://github.com/GoogleCloudPlatform/magic-modules/pull/7336
	//       The test-fixtures folder is not included here as it's copied as a whole,
	//       not file by file
	ignoredFolders := []string{".release/", ".changelog/", "examples/", "scripts/", "META.d/"}
	ignoredFiles := []string{"go.mod", ".goreleaser.yml", ".golangci.yml", "terraform-registry-manifest.json", "_meta.yaml"}
	shouldAddHeader := true
	for _, folder := range ignoredFolders {
		// folder will be path leading to file
		if strings.HasPrefix(target, folder) {
			shouldAddHeader = false
			break
		}
	}
	if !shouldAddHeader {
		return
	}

	for _, file := range ignoredFiles {
		// file will be the filename and extension, with no preceding path
		if strings.HasSuffix(target, file) {
			shouldAddHeader = false
			break
		}
	}
	if !shouldAddHeader {
		return
	}

	lang := languageFromFilename(target)
	// Some file types we don't want to add headers to
	// e.g. .sh where headers are functional
	// Also, this guards against new filetypes being added and triggering build errors
	if lang == "unsupported" {
		return
	}

	// File is not ignored and is appropriate file type to add header to
	copyrightHeader := []string{"Copyright (c) HashiCorp, Inc.", "SPDX-License-Identifier: MPL-2.0"}
	header := commentBlock(copyrightHeader, lang)

	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to add Hashicorp copy right: %s", targetFile, err)
	}

	sourceByte = google.Concat([]byte(header), sourceByte)
	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to add Hashicorp copy right: %s", target, err)
	}
}

func expectedOutputFolder(outputFolder string) bool {
	expectedFolders := []string{"terraform-provider-google", "terraform-provider-google-beta", "terraform-next", "terraform-google-conversion", "tfplan2cai"}
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

func (t Terraform) replaceImportPath(outputFolder, target string) {
	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to replace import path: %s", targetFile, err)
	}

	data := string(sourceByte)

	gaImportPath := ImportPathFromVersion("ga")
	betaImportPath := ImportPathFromVersion("beta")

	if strings.Contains(data, betaImportPath) {
		log.Fatalf("Importing a package from module %s is not allowed in file %s. Please import a package from module %s.", betaImportPath, filepath.Base(target), gaImportPath)
	}

	if t.TargetVersionName == "ga" {
		return
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

	sourceByte = bytes.Replace(sourceByte, []byte(gaImportPath), []byte(tpg+"/"+dir), -1)
	sourceByte = bytes.Replace(sourceByte, []byte(TERRAFORM_PROVIDER_GA+"/version"), []byte(tpg+"/version"), -1)
	sourceByte = bytes.Replace(sourceByte, []byte("module "+TERRAFORM_PROVIDER_GA), []byte("module "+tpg), -1)

	if filepath.Ext(targetFile) == (".go") {
		formatByte, err := format.Source(sourceByte)
		if err != nil {
			log.Printf("error formatting %s: %s", targetFile, err)
		} else {
			sourceByte = formatByte
		}
	}

	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to replace import path: %s", target, err)
	}
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
	case "beta":
		return "/beta"
	case "private":
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
