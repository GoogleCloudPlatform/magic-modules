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
	"reflect"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/product"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
)

const TERRAFORM_PROVIDER_GA = "github.com/hashicorp/terraform-provider-google"
const TERRAFORM_PROVIDER_BETA = "github.com/hashicorp/terraform-provider-google-beta"
const TERRAFORM_PROVIDER_PRIVATE = "internal/terraform-next"
const RESOURCE_DIRECTORY_GA = "google"
const RESOURCE_DIRECTORY_BETA = "google-beta"
const RESOURCE_DIRECTORY_PRIVATE = "google-private"

type Terraform struct {
	ResourceCount int

	IAMResourceCount int

	ResourcesForVersion []map[string]string

	TargetVersionName string

	Version product.Version

	Product *api.Product

	StartTime time.Time
}

func NewTerraform(product *api.Product, versionName string, startTime time.Time) *Terraform {
	t := Terraform{
		ResourceCount:     0,
		IAMResourceCount:  0,
		Product:           product,
		TargetVersionName: versionName,
		Version:           *product.VersionObjOrClosest(versionName),
		StartTime:         startTime,
	}

	t.Product.SetPropertiesBasedOnVersion(&t.Version)
	for _, r := range t.Product.Objects {
		r.SetCompiler(t.providerName())
		r.ImportPath = t.ImportPathFromVersion(versionName)
	}

	return &t
}

func (t *Terraform) Generate(outputFolder, productPath string, generateCode, generateDocs bool) {
	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	t.GenerateObjects(outputFolder, generateCode, generateDocs)

	if generateCode {
		t.GenerateOperation(outputFolder)
	}
}

func (t *Terraform) GenerateObjects(outputFolder string, generateCode, generateDocs bool) {
	for _, object := range t.Product.Objects {
		object.ExcludeIfNotInVersion(&t.Version)

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
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s.go", t.FullResourceName(object)))
		templateData.GenerateResourceFile(targetFilePath, object)
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

func (t *Terraform) GenerateResourceTests(object api.Resource, templateData TemplateData, outputFolder string) {
	eligibleExample := false
	for _, example := range object.Examples {
		if !example.SkipTest {
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
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_generated_test.go", t.FullResourceName(object)))
	templateData.GenerateTestFile(targetFilePath, object)
}

func (t *Terraform) GenerateResourceSweeper(object api.Resource, templateData TemplateData, outputFolder string) {
	if object.SkipSweeper || object.CustomCode.CustomDelete != "" || object.CustomCode.PreDelete != "" || object.CustomCode.PostDelete != "" || object.SkipDelete {
		return
	}

	productName := t.Product.ApiName
	targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("resource_%s_sweeper.go", t.FullResourceName(object)))
	templateData.GenerateSweeperFile(targetFilePath, object)
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
// def generate_iam_policy(pwd, data, generate_code, generate_docs)
func (t *Terraform) GenerateIamPolicy(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	if generateCode && object.IamPolicy != nil && (object.IamPolicy.MinVersion == "" || object.IamPolicy.MinVersion >= t.TargetVersionName) {
		productName := t.Product.ApiName
		targetFolder := path.Join(outputFolder, t.FolderName(), "services", productName)
		if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
			log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
		}
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("iam_%s.go", t.FullResourceName(object)))
		templateData.GenerateIamPolicyFile(targetFilePath, object)

		// Only generate test if testable examples exist.
		examples := google.Reject(object.Examples, func(e resource.Examples) bool {
			return e.SkipTest
		})
		if len(examples) != 0 {
			targetFilePath := path.Join(targetFolder, fmt.Sprintf("iam_%s_generated_test.go", t.FullResourceName(object)))
			templateData.GenerateIamPolicyTestFile(targetFilePath, object)
		}
	}
	if generateDocs {
		t.GenerateIamDocumentation(object, templateData, outputFolder, generateCode, generateDocs)
	}
}

// def generate_iam_documentation(pwd, data)
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
	return "google-beta"
}

func (t *Terraform) FullResourceName(object api.Resource) string {
	if object.LegacyName != "" {
		return strings.Replace(object.LegacyName, "google_", "", 1)
	}

	var name string
	if object.FilenameOverride != "" {
		name = object.FilenameOverride
	} else {
		name = google.Underscore(object.Name)
	}

	var productName string
	if t.Product.LegacyName != "" {
		productName = t.Product.LegacyName
	} else {
		productName = google.Underscore(t.Product.Name)
	}

	return fmt.Sprintf("%s_%s", productName, name)
}

// def copy_common_files(output_folder, generate_code, generate_docs, provider_name = nil)
func (t Terraform) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	log.Printf("Copying common files for %s", t.providerName())

	files := t.getCommonCopyFiles(t.TargetVersionName, generateCode, generateDocs)
	t.CopyFileList(outputFolder, files)
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
		foldersCopiedToGoogleDir = []string{"third_party/terraform/services", "third_party/terraform/acctest", "third_party/terraform/sweeper", "third_party/terraform/provider", "third_party/terraform/tpgdclresource", "third_party/terraform/tpgiamresource", "third_party/terraform/tpgresource", "third_party/terraform/transport", "third_party/terraform/fwmodels", "third_party/terraform/fwprovider", "third_party/terraform/fwtransport", "third_party/terraform/fwresource", "third_party/terraform/verify", "third_party/terraform/envvar", "third_party/terraform/functions", "third_party/terraform/test-fixtures"}
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
		"go.mod":                           "third_party/terraform/go/go.mod",
		".go-version":                      "third_party/terraform/.go-version",
		"terraform-registry-manifest.json": "third_party/terraform/go/terraform-registry-manifest.json",
	}
	maps.Copy(commonCopyFiles, singleFiles)

	return commonCopyFiles
}

func (t Terraform) getCopyFilesInFolder(folderPath, targetDir string) map[string]string {
	m := make(map[string]string, 0)
	filepath.WalkDir(folderPath, func(path string, di fs.DirEntry, err error) error {
		if !di.IsDir() && !strings.HasSuffix(di.Name(), ".tmpl") && !strings.HasSuffix(di.Name(), ".erb") {
			fname := strings.TrimPrefix(strings.Replace(path, "/go/", "/", 1), "third_party/terraform/")
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

// def copy_file_list(output_folder, files)
func (t Terraform) CopyFileList(outputFolder string, files map[string]string) {
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

		err = os.WriteFile(targetFile, sourceByte, 0644)
		if err != nil {
			log.Fatalf("Cannot write target file %s while copying: %s", target, err)
		}

		// Replace import path based on version (beta/alpha)
		if filepath.Ext(target) == ".go" || filepath.Ext(target) == ".mod" {
			t.replaceImportPath(outputFolder, target)
		}

		if filepath.Ext(target) == ".go" {
			t.addHashicorpCopyRightHeader(outputFolder, target)
		}
	}
}

// Compiles files that are shared at the provider level
//
//	def compile_common_files(
//	  output_folder,
//	  products,
//	  common_compile_file,
//	  override_path = nil
//	)
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
		"main.go":                       "third_party/terraform/go/main.go.tmpl",
		".goreleaser.yml":               "third_party/terraform/go/.goreleaser.yml.tmpl",
		".release/release-metadata.hcl": "third_party/terraform/go/release-metadata.hcl.tmpl",
		".copywrite.hcl":                "third_party/terraform/go/.copywrite.hcl.tmpl",
	}
	maps.Copy(commonCompileFiles, singleFiles)

	return commonCompileFiles
}

func (t Terraform) getCompileFilesInFolder(folderPath, targetDir string) map[string]string {
	m := make(map[string]string, 0)
	filepath.WalkDir(folderPath, func(path string, di fs.DirEntry, err error) error {
		if !di.IsDir() && strings.HasSuffix(di.Name(), ".tmpl") {
			fname := strings.TrimPrefix(strings.Replace(path, "/go/", "/", 1), "third_party/terraform/")
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

// def compile_file_list(output_folder, files, file_template, pwd = Dir.pwd)
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
		t.replaceImportPath(outputFolder, target)
		t.addHashicorpCopyRightHeader(outputFolder, target)
	}
}

// def add_hashicorp_copyright_header(output_folder, target)
func (t Terraform) addHashicorpCopyRightHeader(outputFolder, target string) {
	if !expectedOutputFolder(outputFolder) {
		log.Printf("Unexpected output folder (%s) detected"+
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
	ignoredFiles := []string{"go.mod", ".goreleaser.yml", ".golangci.yml", "terraform-registry-manifest.json"}
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

// def expected_output_folder?(output_folder)
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

// def replace_import_path(output_folder, target)
func (t Terraform) replaceImportPath(outputFolder, target string) {
	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to replace import path: %s", targetFile, err)
	}

	data := string(sourceByte)

	gaImportPath := t.ImportPathFromVersion("ga")
	betaImportPath := t.ImportPathFromVersion("beta")

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
		tpg = TERRAFORM_PROVIDER_PRIVATE
		dir = RESOURCE_DIRECTORY_PRIVATE

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

func (t Terraform) ImportPathFromVersion(v string) string {
	var tpg, dir string
	switch v {
	case "ga":
		tpg = TERRAFORM_PROVIDER_GA
		dir = RESOURCE_DIRECTORY_GA
	case "beta":
		tpg = TERRAFORM_PROVIDER_BETA
		dir = RESOURCE_DIRECTORY_BETA
	default:
		tpg = TERRAFORM_PROVIDER_PRIVATE
		dir = RESOURCE_DIRECTORY_PRIVATE
	}
	return fmt.Sprintf("%s/%s", tpg, dir)
}

func (t Terraform) ProviderFromVersion() string {
	var dir string
	switch t.TargetVersionName {
	case "ga":
		dir = RESOURCE_DIRECTORY_GA
	case "beta":
		dir = RESOURCE_DIRECTORY_BETA
	default:
		dir = RESOURCE_DIRECTORY_PRIVATE
	}
	return dir
}

// Gets the list of services dependent on the version ga, beta, and private
// If there are some resources of a servcie is in GA,
// then this service is in GA. Otherwise, the service is in BETA
// def get_mmv1_services_in_version(products, version)
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

// def generate_newyaml(pwd, data)
//
//	# @api.api_name is the service folder name
//	product_name = @api.api_name
//	target_folder = File.join(folder_name(data.version), 'services', product_name)
//	FileUtils.mkpath target_folder
//	data.generate(pwd,
//	              '/templates/terraform/yaml_conversion.erb',
//	              "#{target_folder}/go_#{data.object.name}.yaml",
//	              self)
//	return if File.exist?("#{target_folder}/go_product.yaml")
//
//	data.generate(pwd,
//	              '/templates/terraform/product_yaml_conversion.erb',
//	              "#{target_folder}/go_product.yaml",
//	              self)
//
// end
//
// def build_env
//
//	{
//	  goformat_enabled: @go_format_enabled,
//	  start_time: @start_time
//	}
//
// end
//
// # used to determine and separate objects that have update methods
// # that target individual fields
// def field_specific_update_methods(properties)
//
//	properties_by_custom_update(properties).length.positive?
//
// end
//
// # Filter the properties to keep only the ones requiring custom update
// # method and group them by update url & verb.
// def properties_by_custom_update(properties)
//
//	update_props = properties.reject do |p|
//	  p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP ||
//	    p.is_a?(Api::Type::KeyValueTerraformLabels) ||
//	    p.is_a?(Api::Type::KeyValueLabels) # effective_labels is used for update
//	end
//
//	update_props.group_by do |p|
//	  {
//	    update_url: p.update_url,
//	    update_verb: p.update_verb,
//	    update_id: p.update_id,
//	    fingerprint_name: p.fingerprint_name
//	  }
//	end
//
// end
//
// # Filter the properties to keep only the ones don't have custom update
// # method and group them by update url & verb.
// def properties_without_custom_update(properties)
//
//	properties.select do |p|
//	  p.update_url.nil? || p.update_verb.nil? || p.update_verb == :NOOP
//	end
//
// end
//
// # Takes a update_url and returns the list of custom updatable properties
// # that can be updated at that URL. This allows flattened objects
// # to determine which parent property in the API should be updated with
// # the contents of the flattened object
// def custom_update_properties_by_key(properties, key)
//
//	properties_by_custom_update(properties).select do |k, _|
//	  k[:update_url] == key[:update_url] &&
//	    k[:update_id] == key[:update_id] &&
//	    k[:fingerprint_name] == key[:fingerprint_name]
//	end.first.last
//	# .first is to grab the element from the select which returns a list
//	# .last is because properties_by_custom_update returns a list of
//	# [{update_url}, [properties,...]] and we only need the 2nd part
//
// end
//
// def update_url(resource, url_part)
//
//	[resource.__product.base_url, update_uri(resource, url_part)].flatten.join
//
// end
//
// def generating_hashicorp_repo?
//
//	# The default Provider is used to generate TPG and TPGB in HashiCorp-owned repos.
//	# The compiler deviates from the default behaviour with a -f flag to produce
//	# non-HashiCorp downstreams.
//	true
//
// end
//
// # ProductFileTemplate with Terraform specific fields
// class TerraformProductFileTemplate < Provider::ProductFileTemplate
//
//	# The async object used for making operations.
//	# We assume that all resources share the same async properties.
//	attr_accessor :async
//
//	# When generating OiCS examples, we attach the example we're
//	# generating to the data object.
//	attr_accessor :example
//
//	attr_accessor :resource_name
//
// end
//
// # Sorts properties in the order they should appear in the TF schema:
// # Required, Optional, Computed
// def order_properties(properties)
//
//	properties.select(&:required).sort_by(&:name) +
//	  properties.reject(&:required).reject(&:output).sort_by(&:name) +
//	  properties.select(&:output).sort_by(&:name)
//
// end
//
// def tf_type(property)
//
//	tf_types[property.class]
//
// end
//
// # Converts between the Magic Modules type of an object and its type in the
// # TF schema
// def tf_types
//
//	{
//	  Api::Type::Boolean => 'schema.TypeBool',
//	  Api::Type::Double => 'schema.TypeFloat',
//	  Api::Type::Integer => 'schema.TypeInt',
//	  Api::Type::String => 'schema.TypeString',
//	  # Anonymous string property used in array of strings.
//	  'Api::Type::String' => 'schema.TypeString',
//	  Api::Type::Time => 'schema.TypeString',
//	  Api::Type::Enum => 'schema.TypeString',
//	  Api::Type::ResourceRef => 'schema.TypeString',
//	  Api::Type::NestedObject => 'schema.TypeList',
//	  Api::Type::Array => 'schema.TypeList',
//	  Api::Type::KeyValuePairs => 'schema.TypeMap',
//	  Api::Type::KeyValueLabels => 'schema.TypeMap',
//	  Api::Type::KeyValueTerraformLabels => 'schema.TypeMap',
//	  Api::Type::KeyValueEffectiveLabels => 'schema.TypeMap',
//	  Api::Type::KeyValueAnnotations => 'schema.TypeMap',
//	  Api::Type::Map => 'schema.TypeSet',
//	  Api::Type::Fingerprint => 'schema.TypeString'
//	}
//
// end
//
// def updatable?(resource, properties)
//
//	!resource.immutable || !properties.reject { |p| p.update_url.nil? }.empty?
//
// end
//
// # Returns tuples of (fieldName, list of update masks) for
// #  top-level updatable fields. Schema path refers to a given Terraform
// # field name (e.g. d.GetChange('fieldName)')
// def get_property_update_masks_groups(properties, mask_prefix: ”)
//
//	mask_groups = []
//	properties.each do |prop|
//	  if prop.flatten_object
//	    mask_groups += get_property_update_masks_groups(
//	      prop.properties, mask_prefix: "#{prop.api_name}."
//	    )
//	  elsif prop.update_mask_fields
//	    mask_groups << [prop.name.underscore, prop.update_mask_fields]
//	  else
//	    mask_groups << [prop.name.underscore, [mask_prefix + prop.api_name]]
//	  end
//	end
//	mask_groups
//
// end
//
// # Capitalize the first letter of a property name.
// # E.g. "creationTimestamp" becomes "CreationTimestamp".
// def titlelize_property(property)
//
//	property.name.camelize(:upper)
//
// end
//
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
// def generate_resources_for_version(products, version)
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

				if !(iamPolicy.MinVersion != "" && iamPolicy.MinVersion < t.TargetVersionName) {
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

	// @resources_for_version = @resources_for_version.compact
}

// # TODO(nelsonjr): Review all object interfaces and move to private methods
// # that should not be exposed outside the object hierarchy.
// def provider_name
func (t Terraform) providerName() string {
	return reflect.TypeOf(t).Name()
}

// # Adapted from the method used in templating
// # See: mmv1/compile/core.rb
// def comment_block(text, lang)
func commentBlock(text []string, lang string) string {
	var headers []string
	switch lang {
	case "ruby", "python", "yaml", "gemfile":
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

// def language_from_filename(filename)
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

//	  # Returns the id format of an object, or self_link_uri if none is explicitly defined
//	  # We prefer the long name of a resource as the id so that users can reference
//	  # resources in a standard way, and most APIs accept short name, long name or self_link
//	  def id_format(object)
//	    object.id_format || object.self_link_uri
//	  end

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
		supported = append(supported, v)
		if v == t.TargetVersionName {
			break
		}
	}
	return supported
}

type ProviderWithProducts struct {
	Terraform
	Products []*api.Product
}
