package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl"
	cai2hclconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/cai2hcl/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai"
	tfplan2caiconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tgcresource"
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

func BidirectionalConversion(t *testing.T, ignoredFields []string, primaryResourceType string) {
	testName := t.Name()
	subTestName := GetSubTestName(testName)
	if subTestName == "" {
		t.Skipf("%s: The subtest is unavailable", testName)
	}

	stepNumbers, err := getStepNumbers(subTestName)
	if err != nil {
		t.Fatalf("%s: error preparing the input data: %v", subTestName, err)
	}

	if len(stepNumbers) == 0 {
		t.Skipf("%s: test steps are unavailable", subTestName)
	}

	// Create a temporary directory for running terraform.
	tfDir, err := os.MkdirTemp(tmpDir, "terraform")
	if err != nil {
		t.Fatalf("%s: error creating a temporary directory for running terraform: %v", subTestName, err)
	}
	defer os.RemoveAll(tfDir)

	logger := zaptest.NewLogger(t)

	for _, stepN := range stepNumbers {
		stepName := fmt.Sprintf("step%d", stepN)
		t.Run(stepName, func(t *testing.T) {
			retries := 0
			tName := fmt.Sprintf("%s_%s", subTestName, stepName)
			flakyAction := func(ctx context.Context) error {
				testData, err := prepareTestData(subTestName, stepN, retries)
				retries++
				log.Printf("%s: Starting the attempt %d", tName, retries)
				if err != nil {
					return fmt.Errorf("%s: error preparing the input data: %v", tName, err)
				}

				if testData == nil {
					return retry.RetryableError(fmt.Errorf("fail: test data is unavailable"))
				}

				// If the primary resource is specified, only test the primary resource.
				// Otherwise, test all of the resources in the test.
				primaryResource := testData.PrimaryResource
				resourceTestData := testData.ResourceTestData
				if primaryResource != "" {
					t.Logf("%s: Test for the primary resource %s begins.", tName, primaryResource)
					err = testSingleResource(t, tName, resourceTestData[primaryResource], tfDir, ignoredFields, logger, true)
					if err != nil {
						return err
					}
				} else {
					for address, testData := range resourceTestData {
						if !strings.HasPrefix(address, primaryResourceType) {
							continue
						}
						err = testSingleResource(t, tName, testData, tfDir, ignoredFields, logger, false)
						if err != nil {
							return err
						}
					}
				}

				return nil
			}

			// Note maxAttempts-1 is retries, not attempts.
			backoffPolicy := retry.WithMaxRetries(maxAttempts-1, retry.NewConstant(50*time.Millisecond))

			t.Logf("%s: Starting test with retry logic.", tName)

			if err := retry.Do(context.Background(), backoffPolicy, flakyAction); err != nil {
				if strings.Contains(err.Error(), "test data is unavailable") {
					t.Skipf("%s: Test skipped because data was unavailable after all retries: %v", tName, err)
				} else {
					t.Fatalf("%s: Failed after all attempts %d: %v", tName, maxAttempts, err)
				}
			}
		})
	}
}

// Tests a single resource
func testSingleResource(t *testing.T, testName string, testData ResourceTestData, tfDir string, ignoredFields []string, logger *zap.Logger, primaryResource bool) error {
	resourceType := testData.ResourceType
	var tfplan2caiSupported, cai2hclSupported bool
	if _, tfplan2caiSupported = tfplan2caiconverters.ConverterMap[resourceType]; !tfplan2caiSupported {
		log.Printf("%s: %s is not supported in tfplan2cai conversion.", testName, resourceType)
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
			log.Printf("%s: cai asset is unavailable for %s", testName, assetName)
			return retry.RetryableError(fmt.Errorf("fail: test data is unavailable"))
		}
		if _, cai2hclSupported = cai2hclconverters.ConverterMap[assetType]; !cai2hclSupported {
			log.Printf("%s: %s is not supported in cai2hcl conversion.", testName, assetType)
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
		assetFile := fmt.Sprintf("%s.json", testName)
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
		exportTfFile := fmt.Sprintf("%s_export.tf", testName)
		err = os.WriteFile(exportTfFile, exportConfigData, 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s", exportTfFile)
		}
	}

	exportTfFilePath := fmt.Sprintf("%s/%s_export.tf", tfDir, testName)
	defer os.Remove(exportTfFilePath)
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

	if os.Getenv("WRITE_FILES") != "" {
		writeJSONFile(fmt.Sprintf("%s_export_attrs", testName), exportResources)
	}

	ignoredFieldSet := make(map[string]any, 0)
	for _, f := range ignoredFields {
		ignoredFieldSet[f] = struct{}{}
	}

	parsedExportConfig := exportResources[0].Attributes
	missingKeys := compareHCLFields(testData.ParsedRawConfig, parsedExportConfig, ignoredFieldSet)

	// Sometimes, the reason for missing fields could be CAI asset data issue.
	if len(missingKeys) > 0 {
		log.Printf("%s: missing fields in resource %s after cai2hcl conversion:\n%s", testName, testData.ResourceAddress, missingKeys)
		return retry.RetryableError(fmt.Errorf("missing fields"))
	}
	log.Printf("%s: Step 1 passes for resource %s. All of the fields in raw config are in export config", testName, testData.ResourceAddress)

	// Step 2
	// Run a terraform plan using export_config.
	// Use tfplan2cai to convert the generated plan into CAI assets (roundtrip_assets).
	// Convert roundtrip_assets back into a Terraform configuration (roundtrip_config) using cai2hcl.
	// Compare roundtrip_config with export_config to ensure they are identical.

	// Convert the export config to roundtrip assets and then convert the roundtrip assets back to roundtrip config
	ancestryCache, defaultProject := getAncestryCache(assets)
	roundtripAssets, roundtripConfigData, err := getRoundtripConfig(t, testName, tfDir, ancestryCache, defaultProject, logger)
	if err != nil {
		return fmt.Errorf("error when converting the round-trip config: %#v", err)
	}

	rtTfFile := fmt.Sprintf("%s_roundtrip.tf", testName)
	roundtripTfFilePath := filepath.Join(tfDir, rtTfFile)
	defer os.Remove(roundtripTfFilePath)
	err = os.WriteFile(roundtripTfFilePath, roundtripConfigData, 0644)
	if err != nil {
		return fmt.Errorf("error when writing the file %s", roundtripTfFilePath)
	}
	if os.Getenv("WRITE_FILES") != "" {
		err = os.WriteFile(rtTfFile, roundtripConfigData, 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s", rtTfFile)
		}
	}

	if diff := cmp.Diff(string(roundtripConfigData), string(exportConfigData)); diff != "" {
		tfFileName := fmt.Sprintf("%s_roundtrip", testName)
		jsonFileName := fmt.Sprintf("%s_reexport.json", testName)

		reexportAssets, err := tfplan2caiConvert(t, tfFileName, jsonFileName, tfDir, ancestryCache, defaultProject, logger)
		if err != nil {
			return fmt.Errorf("error when converting the third round-trip config: %#v", err)
		}

		if err = compareCaiAssets(reexportAssets, roundtripAssets, ignoredFieldSet); err != nil {
			log.Printf("%s: Roundtrip config is different from the export config.\nroundtrip config:\n%s\nexport config:\n%s", testName, string(roundtripConfigData), string(exportConfigData))
			return fmt.Errorf("test %s got diff (-want +got): %s", testName, diff)
		}
	}
	log.Printf("%s: Step 2 passes for resource %s. Roundtrip config and export config are identical", testName, testData.ResourceAddress)

	// Step 3
	// Compare most fields between the exported asset and roundtrip asset, except for "data" field for resource
	if err = compareCaiAssets(assets, roundtripAssets, ignoredFieldSet); err != nil {
		return err
	}
	log.Printf("%s: Step 3 passes for resource %s. Exported asset and roundtrip asset are identical", testName, testData.ResourceAddress)

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

			project := tgcresource.ParseFieldValue(asset.Name, "projects")
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
func compareHCLFields(map1, map2, ignoredFields map[string]any) []string {
	var missingKeys []string
	for key, val := range map1 {
		if isIgnored(key, ignoredFields) {
			continue
		}

		if rVal := reflect.ValueOf(val); !rVal.IsValid() || rVal.IsZero() {
			continue
		}

		if sVal, ok := val.(string); ok {
			// TODO: convert to correct type when parsing HCL to fix the edge case where the field type is String and the only values are "false", "00", etc.
			if bVal, err := strconv.ParseBool(sVal); err == nil && !bVal {
				continue
			}
			if iVal, err := strconv.Atoi(sVal); err == nil && iVal == 0 {
				continue
			}
		}

		if vMap, ok := val.(map[string]any); ok && len(vMap) == 0 {
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
func isIgnored(key string, ignoredFields map[string]any) bool {
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
func getRoundtripConfig(t *testing.T, testName string, tfDir string, ancestryCache map[string]string, defaultProject string, logger *zap.Logger) ([]caiasset.Asset, []byte, error) {
	fileName := fmt.Sprintf("%s_export", testName)
	roundtripAssetFile := fmt.Sprintf("%s_roundtrip.json", testName)

	roundtripAssets, err := tfplan2caiConvert(t, fileName, roundtripAssetFile, tfDir, ancestryCache, defaultProject, logger)
	if err != nil {
		return nil, nil, err
	}

	var roundtripAssetsCopy []caiasset.Asset
	// Perform the deep copy in case the assets are transformed in cai2hcl
	err = DeepCopyMap(roundtripAssets, &roundtripAssetsCopy)
	if err != nil {
		fmt.Println("Error during deep copy:", err)
		return nil, nil, err
	}

	roundtripConfig, err := cai2hcl.Convert(roundtripAssetsCopy, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return nil, nil, err
	}

	return roundtripAssets, roundtripConfig, nil
}

// Converts tf file to CAI assets
func tfplan2caiConvert(t *testing.T, tfFileName, jsonFileName string, tfDir string, ancestryCache map[string]string, defaultProject string, logger *zap.Logger) ([]caiasset.Asset, error) {
	// Run terraform init and terraform apply to generate tfplan.json files
	terraformWorkflow(t, tfDir, tfFileName, defaultProject)

	planFile := fmt.Sprintf("%s.tfplan.json", tfFileName)
	planfilePath := filepath.Join(tfDir, planFile)
	defer os.Remove(planfilePath)
	jsonPlan, err := os.ReadFile(planfilePath)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	assets, err := tfplan2cai.Convert(ctx, jsonPlan, &tfplan2cai.Options{
		ErrorLogger:    logger,
		Offline:        true,
		DefaultProject: defaultProject,
		DefaultRegion:  "",
		DefaultZone:    "",
		UserAgent:      "",
		AncestryCache:  ancestryCache,
	})

	if err != nil {
		return nil, err
	}

	if os.Getenv("WRITE_FILES") != "" {
		writeJSONFile(jsonFileName, assets)
	}

	return assets, nil
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

func compareCaiAssets(assets1, assets2 []caiasset.Asset, ignoredFieldSet map[string]any) error {
	assetMap := convertToAssetMap(assets1)
	roundtripAssetMap := convertToAssetMap(assets2)
	for assetType, asset := range assetMap {
		if roundtripAsset, ok := roundtripAssetMap[assetType]; !ok {
			return fmt.Errorf("roundtrip asset for type %s is missing", assetType)
		} else {
			if _, ok := ignoredFieldSet["ASSETNAME"]; !ok {
				if err := compareAssetName(asset.Name, roundtripAsset.Name); err != nil {
					return err
				}
			}
			if diff := cmp.Diff(
				asset.Resource,
				roundtripAsset.Resource,
				// secretmanager.googleapis.com/SecretVersion has secret as parent, not project
				cmpopts.IgnoreFields(caiasset.AssetResource{}, "Version", "Data", "Location", "Parent", "DiscoveryDocumentURI"),
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
				cmp.FilterPath(func(p cmp.Path) bool {
					// Filter if "parent" field in original asset is an empty string
					// and then ingore comparing
					if p.Last().String() == ".Parent" {
						v1, _ := p.Index(-1).Values()
						return v1.IsZero()
					}
					return false
				}, cmp.Ignore()),
			); diff != "" {
				return fmt.Errorf("differences found between exported asset and roundtrip asset (-want +got):\n%s", diff)
			}
		}
	}
	return nil
}
