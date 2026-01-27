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
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/google"
	"github.com/otiai10/copy"
)

var testRegex = regexp.MustCompile("func (TestAcc[^(]+)")

// TerraformGoogleConversionNext is for both tfplan2cai and cai2hcl conversions
// and copying other files, such as transport.go
type TerraformGoogleConversionNext struct {
	ResourceCount int

	ResourcesForVersion []ResourceIdentifier

	// Multiple Terraform resources can share the same CAI resource type.
	// For example, "google_compute_region_autoscaler" and "google_region_autoscaler"
	ResourcesByCaiResourceType map[string][]ResourceIdentifier

	TargetVersionName string

	Product *api.Product

	StartTime time.Time

	templateFS fs.FS
}

type ResourceIdentifier struct {
	ServiceName        string
	TerraformName      string
	ResourceName       string
	AliasName          string // It can be "Default" or the same with ResourceName
	CaiAssetNameFormat string
	IdentityParam      string
}

func NewTerraformGoogleConversionNext(product *api.Product, versionName string, startTime time.Time, templateFS fs.FS) TerraformGoogleConversionNext {
	t := TerraformGoogleConversionNext{
		Product:                    product,
		TargetVersionName:          versionName,
		StartTime:                  startTime,
		ResourcesByCaiResourceType: make(map[string][]ResourceIdentifier),
		templateFS:                 templateFS,
	}

	t.Product.SetCompiler(ProviderName(t))
	for _, r := range t.Product.Objects {
		r.SetCompiler(ProviderName(t))
		r.ImportPath = ImportPathFromVersion(versionName)
	}

	return t
}

func (tgc TerraformGoogleConversionNext) Generate(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	for _, object := range tgc.Product.Objects {
		object.ExcludeIfNotInVersion(tgc.Product.Version)

		if resourceToGenerate != "" && object.Name != resourceToGenerate {
			log.Printf("Excluding %s per user request", object.Name)
			continue
		}

		tgc.GenerateObject(*object, outputFolder, tgc.TargetVersionName, generateCode, generateDocs)
	}
}

func (tgc TerraformGoogleConversionNext) GenerateObject(object api.Resource, outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
	if !object.IncludeInTGCNext {
		return
	}

	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName, tgc.templateFS)

	if !object.IsExcluded() {
		tgc.GenerateResource(object, *templateData, outputFolder, generateCode, generateDocs)
		tgc.addTestsFromSamples(&object)
		if err := tgc.addTestsFromHandwrittenTests(&object); err != nil {
			log.Printf("Error adding examples from handwritten tests: %v", err)
		}
		tgc.GenerateResourceTests(object, *templateData, outputFolder)
	}
}

func (tgc TerraformGoogleConversionNext) GenerateResource(object api.Resource, templateData TemplateData, outputFolder string, generateCode, generateDocs bool) {
	productName := tgc.Product.ApiName
	targetFolder := path.Join(outputFolder, "pkg/services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}

	converters := []string{"tfplan2cai", "cai2hcl"}
	for _, converter := range converters {
		templatePath := fmt.Sprintf("templates/tgc_next/%s/resource_converter.go.tmpl", converter)
		targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_%s_%s.go", productName, google.Underscore(object.Name), converter))
		templateData.GenerateTGCResourceFile(templatePath, targetFilePath, object)
	}

	templatePath := "templates/tgc_next/services/resource.go.tmpl"
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_%s.go", productName, google.Underscore(object.Name)))
	templateData.GenerateTGCResourceFile(templatePath, targetFilePath, object)
}

func (tgc TerraformGoogleConversionNext) GenerateCaiToHclObjects(outputFolder, resourceToGenerate string, generateCode, generateDocs bool) {
}

func (tgc *TerraformGoogleConversionNext) GenerateResourceTests(object api.Resource, templateData TemplateData, outputFolder string) {
	if len(object.TGCTests) == 0 {
		return
	}

	productName := tgc.Product.ApiName
	targetFolder := path.Join(outputFolder, "test", "services", productName)
	if err := os.MkdirAll(targetFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating parent directory %v: %v", targetFolder, err))
	}
	targetFilePath := path.Join(targetFolder, fmt.Sprintf("%s_%s_generated_test.go", productName, google.Underscore(object.Name)))
	templateData.GenerateTGCNextTestFile(targetFilePath, object)
}

func (tgc TerraformGoogleConversionNext) CompileCommonFiles(outputFolder string, products []*api.Product, overridePath string) {
	tgc.generateResourcesForVersion(products)

	resourceConverters := map[string]string{
		// common
		"pkg/transport/config.go":                        "third_party/terraform/transport/config.go.tmpl",
		"pkg/transport/provider_handwritten_endpoint.go": "third_party/terraform/transport/provider_handwritten_endpoint.go.tmpl",
		"pkg/tpgresource/common_diff_suppress.go":        "third_party/terraform/tpgresource/common_diff_suppress.go",
		"pkg/provider/provider.go":                       "third_party/terraform/provider/provider.go.tmpl",
		"pkg/provider/provider_validators.go":            "third_party/terraform/provider/provider_validators.go",
		"pkg/provider/provider_mmv1_resources.go":        "templates/tgc_next/provider/provider_mmv1_resources.go.tmpl",

		// services
		"pkg/services/compute/compute_instance_helpers.go": "third_party/terraform/services/compute/compute_instance_helpers.go.tmpl",
		"pkg/services/compute/metadata.go":                 "third_party/terraform/services/compute/metadata.go.tmpl",

		// tfplan2cai
		"pkg/tfplan2cai/converters/resource_converters.go": "templates/tgc_next/tfplan2cai/resource_converters.go.tmpl",

		// cai2hcl
		"pkg/cai2hcl/converters/resource_converters.go": "templates/tgc_next/cai2hcl/resource_converters.go.tmpl",
		"pkg/cai2hcl/converters/convert_resource.go":    "templates/tgc_next/cai2hcl/convert_resource.go.tmpl",
	}

	templateData := NewTemplateData(outputFolder, tgc.TargetVersionName, tgc.templateFS)
	tgc.CompileFileList(outputFolder, resourceConverters, *templateData, products)
}

func (tgc TerraformGoogleConversionNext) CompileFileList(outputFolder string, files map[string]string, fileTemplate TemplateData, products []*api.Product) {
	providerWithProducts := TgcWithProducts{
		TerraformGoogleConversionNext: tgc,
		Compiler:                      "terraformgoogleconversion-codegen",
		Products:                      products,
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
		tgc.replaceImportPath(outputFolder, target)
	}
}

func (tgc TerraformGoogleConversionNext) CopyCommonFiles(outputFolder string, generateCode, generateDocs bool) {
	if !generateCode {
		return
	}

	log.Printf("Copying common files for tgc.")

	if err := os.MkdirAll(outputFolder, os.ModePerm); err != nil {
		log.Println(fmt.Errorf("error creating output directory %v: %v", outputFolder, err))
	}

	if err := copy.Copy("third_party/tgc_next", outputFolder); err != nil {
		log.Println(fmt.Errorf("error copying directory %v: %v", outputFolder, err))
	}

	resourceConverters := map[string]string{
		// common
		"pkg/transport/batcher.go":                 "third_party/terraform/transport/batcher.go",
		"pkg/transport/retry_transport.go":         "third_party/terraform/transport/retry_transport.go",
		"pkg/transport/retry_utils.go":             "third_party/terraform/transport/retry_utils.go",
		"pkg/transport/header_transport.go":        "third_party/terraform/transport/header_transport.go",
		"pkg/transport/error_retry_predicates.go":  "third_party/terraform/transport/error_retry_predicates.go",
		"pkg/transport/bigtable_client_factory.go": "third_party/terraform/transport/bigtable_client_factory.go",
		"pkg/transport/transport.go":               "third_party/terraform/transport/transport.go",
		"pkg/tpgresource/utils.go":                 "third_party/terraform/tpgresource/utils.go",
		"pkg/tpgresource/self_link_helpers.go":     "third_party/terraform/tpgresource/self_link_helpers.go",
		"pkg/tpgresource/hashcode.go":              "third_party/terraform/tpgresource/hashcode.go",
		"pkg/tpgresource/regional_utils.go":        "third_party/terraform/tpgresource/regional_utils.go",
		"pkg/tpgresource/field_helpers.go":         "third_party/terraform/tpgresource/field_helpers.go",
		"pkg/tpgresource/service_scope.go":         "third_party/terraform/tpgresource/service_scope.go",
		"pkg/provider/mtls_util.go":                "third_party/terraform/provider/mtls_util.go",
		"pkg/verify/validation.go":                 "third_party/terraform/verify/validation.go",
		"pkg/verify/path_or_contents.go":           "third_party/terraform/verify/path_or_contents.go",
		"pkg/version/version.go":                   "third_party/terraform/version/version.go",

		// services
		"pkg/services/compute/image.go":             "third_party/terraform/services/compute/image.go",
		"pkg/services/compute/disk_type.go":         "third_party/terraform/services/compute/disk_type.go",
		"pkg/services/kms/kms_utils.go":             "third_party/terraform/services/kms/kms_utils.go",
		"pkg/services/privateca/privateca_utils.go": "third_party/terraform/services/privateca/privateca_utils.go",
	}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionNext) CopyTfToCaiCommonFiles(outputFolder string) {
	resourceConverters := map[string]string{
		"pkg/tfplan2cai/converters/services/compute/image.go":     "third_party/terraform/services/compute/image.go",
		"pkg/tfplan2cai/converters/services/compute/disk_type.go": "third_party/terraform/services/compute/disk_type.go",
	}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionNext) CopyCaiToHclCommonFiles(outputFolder string) {
	resourceConverters := map[string]string{}
	tgc.CopyFileList(outputFolder, resourceConverters)
}

func (tgc TerraformGoogleConversionNext) CopyFileList(outputFolder string, files map[string]string) {
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

		sourceByte, err := fs.ReadFile(tgc.templateFS, source)
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

func (tgc TerraformGoogleConversionNext) replaceImportPath(outputFolder, target string) {
	// Replace import paths to reference the resources dir instead of the google provider
	targetFile := filepath.Join(outputFolder, target)
	sourceByte, err := os.ReadFile(targetFile)
	if err != nil {
		log.Fatalf("Cannot read file %s to replace import path: %s", targetFile, err)
	}

	// replace google to google-beta
	gaImportPath := ImportPathFromVersion("ga")
	sourceByte = bytes.Replace(sourceByte, []byte(gaImportPath), []byte(TGC_PROVIDER+"/"+RESOURCE_DIRECTORY_TGC), -1)
	sourceByte = bytes.Replace(sourceByte, []byte(TERRAFORM_PROVIDER_GA+"/version"), []byte(TGC_PROVIDER+"/"+RESOURCE_DIRECTORY_TGC+"/version"), -1)

	err = os.WriteFile(targetFile, sourceByte, 0644)
	if err != nil {
		log.Fatalf("Cannot write file %s to replace import path: %s", target, err)
	}
}

func (tgc TerraformGoogleConversionNext) addTestsFromExamples(object *api.Resource) {
	for _, example := range object.Examples {
		if example.ExcludeTest {
			continue
		}
		if object.ProductMetadata.VersionObjOrClosest(tgc.Product.Version.Name).CompareTo(object.ProductMetadata.VersionObjOrClosest(example.MinVersion)) < 0 {
			continue
		}
		object.TGCTests = append(object.TGCTests, resource.TGCTest{
			Name: "TestAcc" + example.TestSlug(object.ProductMetadata.Name, object.Name),
			Skip: example.TGCSkipTest,
		})
	}
}

func (tgc TerraformGoogleConversionNext) addTestsFromSamples(object *api.Resource) {
	if object.Examples != nil {
		tgc.addTestsFromExamples(object)
		return
	}
	for _, sample := range object.Samples {
		if sample.ExcludeTest {
			continue
		}
		if object.ProductMetadata.VersionObjOrClosest(tgc.Product.Version.Name).CompareTo(object.ProductMetadata.VersionObjOrClosest(sample.MinVersion)) < 0 {
			continue
		}
		object.TGCTests = append(object.TGCTests, resource.TGCTest{
			Name: "TestAcc" + sample.TestSampleSlug(object.ProductMetadata.Name, object.Name),
			Skip: sample.TGCSkipTest,
		})
	}
}
func (tgc TerraformGoogleConversionNext) addTestsFromHandwrittenTests(object *api.Resource) error {
	if object.ProductMetadata == nil {
		return nil
	}
	productName := strings.ToLower(tgc.Product.Name)
	resourceFullName := tgc.ResourceGoFilename(*object)
	handwrittenTestFilePath := fmt.Sprintf("third_party/terraform/services/%s/resource_%s_test.go", productName, resourceFullName)
	data, err := fs.ReadFile(tgc.templateFS, handwrittenTestFilePath)
	for err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if strings.HasSuffix(handwrittenTestFilePath, ".tmpl") {
				log.Printf("no handwritten test file found for %s", resourceFullName)
				return nil
			}
			handwrittenTestFilePath += ".tmpl"
			data, err = fs.ReadFile(tgc.templateFS, handwrittenTestFilePath)
		} else {
			return fmt.Errorf("error reading handwritten test file %s: %v", handwrittenTestFilePath, err)
		}
	}

	// Skip adding handwritten tests that are already defined in yaml (because they have custom overrides etc.)
	testNamesInYAML := make(map[string]struct{})
	for _, test := range object.TGCTests {
		if test.Name != "" {
			testNamesInYAML[test.Name] = struct{}{}
		}
	}

	matches := testRegex.FindAllSubmatch(data, -1)
	tests := make([]resource.TGCTest, len(matches))
	for i, match := range matches {
		if len(match) == 2 {
			if _, ok := testNamesInYAML[string(match[1])]; ok {
				continue
			}
			tests[i] = resource.TGCTest{
				Name: string(match[1]),
			}
		}
	}

	object.TGCTests = append(object.TGCTests, tests...)

	return nil
}

// Similar to FullResourceName, but override-aware to prevent things like ending in _test.
// Non-Go files should just use FullResourceName.
func (tgc *TerraformGoogleConversionNext) ResourceGoFilename(object api.Resource) string {
	// early exit if no override is set
	if object.FilenameOverride == "" {
		return tgc.FullResourceName(object)
	}

	resName := object.FilenameOverride

	var productName string
	if tgc.Product.LegacyName != "" {
		productName = tgc.Product.LegacyName
	} else {
		productName = google.Underscore(tgc.Product.Name)
	}

	return fmt.Sprintf("%s_%s", productName, resName)
}

func (tgc *TerraformGoogleConversionNext) FullResourceName(object api.Resource) string {
	// early exit- resource-level legacy names override the product too
	if object.LegacyName != "" {
		return strings.Replace(object.LegacyName, "google_", "", 1)
	}

	var productName string
	if tgc.Product.LegacyName != "" {
		productName = tgc.Product.LegacyName
	} else {
		productName = google.Underscore(tgc.Product.Name)
	}

	return fmt.Sprintf("%s_%s", productName, google.Underscore(object.Name))
}

// Generates the list of resources, and gets the count of resources.
// The resource object has the format
//
//	{
//	   terraform_name:
//	   resource_name:
//	}
//
// The variable resources_for_version is used to generate resources in file
// mmv1/templates/tgc_next/provider/provider_mmv1_resources.go.tmpl
func (tgc *TerraformGoogleConversionNext) generateResourcesForVersion(products []*api.Product) {
	resourcesByCaiResourceType := make(map[string][]ResourceIdentifier)

	for _, productDefinition := range products {
		service := strings.ToLower(productDefinition.Name)
		for _, object := range productDefinition.Objects {
			if object.Exclude || object.NotInVersion(productDefinition.VersionObjOrClosest(tgc.TargetVersionName)) {
				continue
			}

			if !object.IncludeInTGCNext {
				continue
			}

			tgc.ResourceCount++

			resourceIdentifier := ResourceIdentifier{
				ServiceName:        service,
				TerraformName:      object.TerraformName(),
				ResourceName:       object.ResourceName(),
				AliasName:          object.ResourceName(),
				CaiAssetNameFormat: object.GetCaiAssetNameTemplate(),
			}
			tgc.ResourcesForVersion = append(tgc.ResourcesForVersion, resourceIdentifier)

			caiResourceType := object.CaiAssetType()
			if _, ok := resourcesByCaiResourceType[caiResourceType]; !ok {
				resourcesByCaiResourceType[caiResourceType] = make([]ResourceIdentifier, 0)
			}
			resourcesByCaiResourceType[caiResourceType] = append(resourcesByCaiResourceType[caiResourceType], resourceIdentifier)
		}
	}

	for caiResourceType, resources := range resourcesByCaiResourceType {
		// If no other Terraform resources share the API resource type, override the alias name as "Default"
		if len(resources) == 1 {
			for _, resourceIdentifier := range resources {
				resourceIdentifier.AliasName = "Default"
				tgc.ResourcesByCaiResourceType[caiResourceType] = []ResourceIdentifier{resourceIdentifier}
			}
		} else {
			tgc.ResourcesByCaiResourceType[caiResourceType] = FindIdentityParams(resources)
		}
	}
}

// Analyzes a list of CAI asset names and finds the single path segment
// (by index) that contains different values across all names.
// Example:
// "folders/{{folder}}/feeds/{{feed_id}}" -> folders
// "organizations/{{org_id}}/feeds/{{feed_id}} -> organizations
// "projects/{{project}}/feeds/{{feed_id}}" -> projects
func FindIdentityParams(rids []ResourceIdentifier) []ResourceIdentifier {
	segmentsList := make([][]string, len(rids))
	for i, rid := range rids {
		urlPath := rid.CaiAssetNameFormat
		urlPath = strings.Trim(urlPath, "/")

		processedURL := regexp.MustCompile(`\{\{%?(\w+)\}\}`).ReplaceAllString(urlPath, "")
		segments := strings.Split(processedURL, "/")
		var cleanSegments []string
		for _, seg := range segments {
			if seg != "" {
				cleanSegments = append(cleanSegments, seg)
			}
		}

		segmentsList[i] = cleanSegments
	}

	segmentsList = removeSharedElements(segmentsList)

	for i, segments := range segmentsList {
		if len(segments) == 0 {
			rids[i].IdentityParam = ""
		} else {
			rids[i].IdentityParam = segments[0]
		}
	}

	// Move the id with empty IdentityParam to the end of the list
	for i, ids := range rids {
		if ids.IdentityParam == "" {
			temp := ids
			lastIndex := len(rids) - 1
			if i != lastIndex {
				rids[i] = rids[lastIndex]
				rids[lastIndex] = temp
			}
			break
		}
	}

	return rids
}

// Finds elements common to ALL lists in a list of lists
// and returns a new list of lists with those common elements removed.
func removeSharedElements(list_of_lists [][]string) [][]string {
	if len(list_of_lists) <= 1 {
		return list_of_lists
	}

	sharedSet := make(map[string]bool)
	for _, element := range list_of_lists[0] {
		sharedSet[element] = true
	}

	for i := 1; i < len(list_of_lists); i++ {
		currentListSet := make(map[string]bool)
		for _, element := range list_of_lists[i] {
			currentListSet[element] = true
		}

		newSharedSet := make(map[string]bool)

		for element := range sharedSet {
			if currentListSet[element] {
				newSharedSet[element] = true
			}
		}

		sharedSet = newSharedSet

		if len(sharedSet) == 0 {
			break
		}
	}

	var new_list_of_lists [][]string

	for _, sublist := range list_of_lists {
		var newSublist []string
		for _, element := range sublist {
			if !sharedSet[element] {
				newSublist = append(newSublist, element)
			}
		}
		new_list_of_lists = append(new_list_of_lists, newSublist)
	}

	return new_list_of_lists
}

type TgcWithProducts struct {
	TerraformGoogleConversionNext
	Compiler string
	Products []*api.Product
}
