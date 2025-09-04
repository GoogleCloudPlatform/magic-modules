package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl"
	cai2hclconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai"
	tfplan2caiconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters"
	"github.com/sethvargo/go-retry"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	cacheMutex = sync.Mutex{}
	tmpDir     = os.TempDir()
)

func BidirectionalConversion(t *testing.T, ignoredFields []string, ignoredAssetFields []string) {
	retries := 0
	flakyAction := func(ctx context.Context) error {
		log.Printf("Starting the retry %d", retries)
		resourceTestData, primaryResource, err := prepareTestData(t.Name(), retries)
		retries++
		if err != nil {
			return fmt.Errorf("error preparing the input data: %v", err)
		}

		if resourceTestData == nil {
			return retry.RetryableError(fmt.Errorf("fail: test data is unavailable"))
		}

		// Create a temporary directory for running terraform.
		tfDir, err := os.MkdirTemp(tmpDir, "terraform")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tfDir)

		logger := zaptest.NewLogger(t)

		// If the primary resource is specified, only test the primary resource.
		// Otherwise, test all of the resources in the test.
		if primaryResource != "" {
			t.Logf("Test for the primary resource %s begins.", primaryResource)
			err = testSingleResource(t, t.Name(), resourceTestData[primaryResource], tfDir, ignoredFields, ignoredAssetFields, logger, true)
			if err != nil {
				return err
			}
		} else {
			for _, testData := range resourceTestData {
				err = testSingleResource(t, t.Name(), testData, tfDir, ignoredFields, ignoredAssetFields, logger, false)
				if err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Note maxAttempts-1 is retries, not attempts.
	backoffPolicy := retry.WithMaxRetries(maxAttempts-1, retry.NewConstant(50*time.Millisecond))

	t.Log("Starting test with retry logic.")

	if err := retry.Do(context.Background(), backoffPolicy, flakyAction); err != nil {
		if strings.Contains(err.Error(), "test data is unavailable") {
			t.Skipf("Test skipped because data was unavailable after all retries: %v", err)
		} else {
			t.Fatalf("Failed after all retries %d: %v", retries, err)
		}
	}
}

// Tests a single resource
func testSingleResource(t *testing.T, testName string, testData ResourceTestData, tfDir string, ignoredFields []string, ignoredAssetFields []string, logger *zap.Logger, primaryResource bool) error {
	resourceType := testData.ResourceType
	var tfplan2caiSupported, cai2hclSupported bool
	if _, tfplan2caiSupported = tfplan2caiconverters.ConverterMap[resourceType]; !tfplan2caiSupported {
		log.Printf("%s is not supported in tfplan2cai conversion.", resourceType)
	}

	if testData.Cai == nil {
		log.Printf("SKIP: cai asset is unavailable for resource %s", testData.ResourceAddress)
		return nil
	}

	assets := make([]caiasset.Asset, 0)
	for assetName, assetData := range testData.Cai {
		assets = append(assets, assetData.CaiAsset)
		assetType := assetData.CaiAsset.Type
		if assetType == "" {
			return fmt.Errorf("cai asset is unavailable for %s", assetName)
		}
		if _, cai2hclSupported = cai2hclconverters.ConverterMap[assetType]; !cai2hclSupported {
			log.Printf("%s is not supported in cai2hcl conversion.", assetType)
		}
	}

	if !tfplan2caiSupported && !cai2hclSupported {
		if primaryResource {
			return fmt.Errorf("conversion of the primary resource %s is not supported in tgc", testData.ResourceAddress)
		} else {
			log.Printf("SKIP: conversion of the resource %s is not supported in tgc.", resourceType)
			return nil
		}
	}

	if !(tfplan2caiSupported && cai2hclSupported) {
		return fmt.Errorf("resource %s is supported in either tfplan2cai or cai2hcl within tgc, but not in both", resourceType)
	}

	if os.Getenv("WRITE_FILES") != "" {
		assetFile := fmt.Sprintf("%s.json", t.Name())
		writeJSONFile(assetFile, assets)
	}

	// Step 1: Use cai2hcl to convert export assets into a Terraform configuration (export config).
	// Compare all of the fields in raw config are in export config.

	exportConfigData, err := cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return fmt.Errorf("error when converting the export assets into export config: %#v", err)
	}

	if os.Getenv("WRITE_FILES") != "" {
		exportTfFile := fmt.Sprintf("%s_export.tf", t.Name())
		err = os.WriteFile(exportTfFile, exportConfigData, 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s", exportTfFile)
		}
	}

	exportTfFilePath := fmt.Sprintf("%s/%s_export.tf", tfDir, t.Name())
	err = os.WriteFile(exportTfFilePath, exportConfigData, 0644)
	if err != nil {
		return fmt.Errorf("error when writing the file %s", exportTfFilePath)
	}

	exportResources, err := parseResourceConfigs(exportTfFilePath)
	if err != nil {
		return err
	}

	if len(exportResources) == 0 {
		return fmt.Errorf("missing hcl after cai2hcl conversion for resource %s", testData.ResourceType)
	}

	ignoredFieldSet := make(map[string]struct{}, 0)
	for _, f := range ignoredFields {
		ignoredFieldSet[f] = struct{}{}
	}

	parsedExportConfig := exportResources[0].Attributes
	missingKeys := compareHCLFields(testData.ParsedRawConfig, parsedExportConfig, ignoredFieldSet)

	// Sometimes, the reason for missing fields could be CAI asset data issue.
	if len(missingKeys) > 0 {
		log.Printf("missing fields in resource %s after cai2hcl conversion:\n%s", testData.ResourceAddress, missingKeys)
		return retry.RetryableError(fmt.Errorf("missing fields"))
	}
	log.Printf("Step 1 passes for resource %s. All of the fields in raw config are in export config", testData.ResourceAddress)

	// Step 2
	// Run a terraform plan using export_config.
	// Use tfplan2cai to convert the generated plan into CAI assets (roundtrip_assets).
	// Convert roundtrip_assets back into a Terraform configuration (roundtrip_config) using cai2hcl.
	// Compare roundtrip_config with export_config to ensure they are identical.

	// Convert the export config to roundtrip assets and then convert the roundtrip assets back to roundtrip config
	ancestryCache, defaultProject := getAncestryCache(assets)
	roundtripAssets, roundtripConfigData, err := getRoundtripConfig(t, testName, tfDir, ancestryCache, defaultProject, logger, ignoredAssetFields)
	if err != nil {
		return fmt.Errorf("error when converting the round-trip config: %#v", err)
	}

	roundtripTfFilePath := fmt.Sprintf("%s_roundtrip.tf", testName)
	err = os.WriteFile(roundtripTfFilePath, roundtripConfigData, 0644)
	if err != nil {
		return fmt.Errorf("error when writing the file %s", roundtripTfFilePath)
	}
	if os.Getenv("WRITE_FILES") == "" {
		defer os.Remove(roundtripTfFilePath)
	}

	if diff := cmp.Diff(string(roundtripConfigData), string(exportConfigData)); diff != "" {
		log.Printf("Roundtrip config is different from the export config.\nroundtrip config:\n%s\nexport config:\n%s", string(roundtripConfigData), string(exportConfigData))
		return fmt.Errorf("test %s got diff (-want +got): %s", testName, diff)
	}
	log.Printf("Step 2 passes for resource %s. Roundtrip config and export config are identical", testData.ResourceAddress)

	// Step 3
	// Compare most fields between the exported asset and roundtrip asset, except for "data" field for resource
	assetMap := convertToAssetMap(assets)
	roundtripAssetMap := convertToAssetMap(roundtripAssets)
	for assetType, asset := range assetMap {
		if roundtripAsset, ok := roundtripAssetMap[assetType]; !ok {
			return fmt.Errorf("roundtrip asset for type %s is missing", assetType)
		} else {
			if err := compareAssetName(asset.Name, roundtripAsset.Name); err != nil {
				return err
			}
			if diff := cmp.Diff(
				asset.Resource,
				roundtripAsset.Resource,
				cmpopts.IgnoreFields(caiasset.AssetResource{}, "Version", "Data", "Location", "DiscoveryDocumentURI"),
				// Consider DiscoveryDocumentURI equal if they have the same number of path segments when split by "/".
				cmp.FilterPath(func(p cmp.Path) bool {
					return p.Last().String() == ".DiscoveryDocumentURI"
				}, cmp.Comparer(func(x, y string) bool {
					parts1 := strings.Split(x, "/")
					parts2 := strings.Split(y, "/")
					return len(parts1) == len(parts2)
				})),
				cmp.FilterPath(func(p cmp.Path) bool {
					return p.Last().String() == ".DiscoveryName"
				}, cmp.Comparer(func(x, y string) bool {
					xParts := strings.Split(x, "/")
					yParts := strings.Split(y, "/")
					return xParts[len(xParts)-1] == yParts[len(yParts)-1]
				})),
			); diff != "" {
				return fmt.Errorf("differences found between exported asset and roundtrip asset (-want +got):\n%s", diff)
			}
		}
	}
	log.Printf("Step 3 passes for resource %s. Exported asset and roundtrip asset are identical", testData.ResourceAddress)

	return nil
}

// Gets the ancestry cache for tfplan2cai conversion and the default project
func getAncestryCache(assets []caiasset.Asset) (map[string]string, string) {
	ancestryCache := make(map[string]string, 0)
	defaultProject := ""

	for _, asset := range assets {
		ancestors := asset.Ancestors
		if len(ancestors) != 0 {
			var path string
			for i := len(ancestors) - 1; i >= 0; i-- {
				curr := ancestors[i]
				if path == "" {
					path = curr
				} else {
					path = fmt.Sprintf("%s/%s", path, curr)
				}
			}

			if _, ok := ancestryCache[ancestors[0]]; !ok {
				ancestryCache[ancestors[0]] = path
				if defaultProject == "" {
					if s, hasPrefix := strings.CutPrefix(ancestors[0], "projects/"); hasPrefix {
						defaultProject = s
					}
				}
			}

			project := utils.ParseFieldValue(asset.Name, "projects")
			if project != "" {
				projectKey := fmt.Sprintf("projects/%s", project)
				if strings.HasPrefix(ancestors[0], "projects") && ancestors[0] != projectKey {
					if _, ok := ancestryCache[projectKey]; !ok {
						ancestryCache[projectKey] = path
					}
				}

				if defaultProject == "" {
					defaultProject = project
				}
			}
		}
	}
	return ancestryCache, defaultProject
}

// Compares HCL and finds all of the keys in map1 that are not in map2
func compareHCLFields(map1, map2, ignoredFields map[string]struct{}) []string {
	var missingKeys []string
	for key := range map1 {
		if isIgnored(key, ignoredFields) {
			continue
		}

		if _, ok := map2[key]; !ok {
			missingKeys = append(missingKeys, key)
		}
	}
	sort.Strings(missingKeys)
	return missingKeys
}

// Returns true if the given key should be ignored according to the given set of ignored fields.
func isIgnored(key string, ignoredFields map[string]struct{}) bool {
	// Check for exact match first.
	if _, ignored := ignoredFields[key]; ignored {
		return true
	}

	// Check for partial matches.
	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return false
	}
	var nonIntegerParts []string
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			nonIntegerParts = append(nonIntegerParts, part)
		}
	}
	var partialKey string
	for _, part := range nonIntegerParts {
		if partialKey == "" {
			partialKey = part
		} else {
			partialKey += "." + part
		}
		if _, ignored := ignoredFields[partialKey]; ignored {
			return true
		}
	}
	return false
}

// Converts a tfplan to CAI asset, and then converts the CAI asset into HCL
func getRoundtripConfig(t *testing.T, testName string, tfDir string, ancestryCache map[string]string, defaultProject string, logger *zap.Logger, ignoredAssetFields []string) ([]caiasset.Asset, []byte, error) {
	fileName := fmt.Sprintf("%s_export", testName)

	// Run terraform init and terraform apply to generate tfplan.json files
	terraformWorkflow(t, tfDir, fileName)

	planFile := fmt.Sprintf("%s.tfplan.json", fileName)
	planfilePath := filepath.Join(tfDir, planFile)
	jsonPlan, err := os.ReadFile(planfilePath)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	roundtripAssets, err := tfplan2cai.Convert(ctx, jsonPlan, &tfplan2cai.Options{
		ErrorLogger:    logger,
		Offline:        true,
		DefaultProject: defaultProject,
		DefaultRegion:  "",
		DefaultZone:    "",
		UserAgent:      "",
		AncestryCache:  ancestryCache,
	})

	if err != nil {
		return nil, nil, err
	}

	deleteFieldsFromAssets(roundtripAssets, ignoredAssetFields)

	if os.Getenv("WRITE_FILES") != "" {
		roundtripAssetFile := fmt.Sprintf("%s_roundtrip.json", t.Name())
		writeJSONFile(roundtripAssetFile, roundtripAssets)
	}

	roundtripConfig, err := cai2hcl.Convert(roundtripAssets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return nil, nil, err
	}

	return roundtripAssets, roundtripConfig, nil
}

// Example:
//
//	data := map[string]interface{}{
//		"database": map[string]interface{}{
//			"host": "localhost",
//			"user": "admin",
//		},
//	}
//
// Path of "host" in "data" is ["database", "host"]
type Field struct {
	Path []string
}

// Deletes fields from the resource data of CAI assets
func deleteFieldsFromAssets(assets []caiasset.Asset, ignoredResourceDataFields []string) []caiasset.Asset {
	// The key is the content type, such as "resource"
	ignoredFieldsMap := make(map[string][]Field, 0)
	for _, ignoredField := range ignoredResourceDataFields {
		parts := strings.Split(ignoredField, ".")
		if len(parts) <= 1 {
			continue
		}
		if parts[0] == "RESOURCE" {
			if _, ok := ignoredFieldsMap["RESOURCE"]; !ok {
				ignoredFieldsMap["RESOURCE"] = make([]Field, 0)
			}
			f := Field{Path: parts[1:]}
			ignoredFieldsMap["RESOURCE"] = append(ignoredFieldsMap["RESOURCE"], f)
		}
	}

	for _, asset := range assets {
		if asset.Resource != nil && asset.Resource.Data != nil {
			data := asset.Resource.Data
			for _, ignoredField := range ignoredFieldsMap["RESOURCE"] {
				path := ignoredField.Path
				deleteMapFieldByPath(data, path)
			}
		}
	}
	return assets
}

// Deletes a field from a map by its path.
// Example:
//
//	data := map[string]interface{}{
//		"database": map[string]interface{}{
//			"host": "localhost",
//			"user": "admin",
//		},
//	}
//
// path := ["database", "host"]
func deleteMapFieldByPath(data map[string]interface{}, path []string) {
	i := 0
	for i < len(path)-1 {
		k := path[i]
		if v, ok := data[k]; ok {
			if data, ok = v.(map[string]interface{}); ok && data != nil {
				i++
			} else {
				break
			}
		} else {
			break
		}
	}
	if i == len(path)-1 {
		delete(data, path[i])
	}
}

// Compares the asset name in export asset and roundtrip asset and ignores "null" in the name
// Example: //cloudresourcemanager.googleapis.com/projects/123456
func compareAssetName(want, got string) error {
	parts1 := strings.Split(want, "/")
	parts2 := strings.Split(got, "/")
	if len(parts1) != len(parts2) {
		return fmt.Errorf("differences found between two asset names: want %s, got %s", want, got)
	}

	for i, part := range parts1 {
		if parts2[i] == "null" {
			continue
		}

		if part != parts2[i] {
			return fmt.Errorf("differences found between two asset names: want %s, got %s", want, got)
		}
	}
	return nil
}
