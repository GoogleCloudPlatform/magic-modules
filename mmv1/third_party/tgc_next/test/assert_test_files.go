package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl"
	cai2hclconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai"
	tfplan2caiconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai/converters"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	cacheMutex = sync.Mutex{}
	tmpDir     = os.TempDir()
)

func BidirectionalConversion(t *testing.T, ignoredFields []string) {
	resourceTestData, primaryResource, err := prepareTestData(t.Name())
	if err != nil {
		t.Fatal("Error preparing the input data:", err)
	}

	if resourceTestData == nil {
		t.Skipf("The test data is unavailable.")
	}

	// Create a temporary directory for running terraform.
	tfDir, err := os.MkdirTemp(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tfDir)

	logger := zaptest.NewLogger(t)

	// If the primary resource is available, only test the primary resource.
	// Otherwise, test all of the resources in the test.
	if primaryResource != "" {
		t.Logf("Test for the primary resource %s begins.", primaryResource)
		err = testSingleResource(t, t.Name(), resourceTestData[primaryResource], tfDir, ignoredFields, logger, true)
		if err != nil {
			t.Fatal("Test fails:", err)
		}
	} else {
		for _, testData := range resourceTestData {
			err = testSingleResource(t, t.Name(), testData, tfDir, ignoredFields, logger, false)
			if err != nil {
				t.Fatal("Test fails: ", err)
			}
		}
	}
}

// Tests a single resource
func testSingleResource(t *testing.T, testName string, testData ResourceTestData, tfDir string, ignoredFields []string, logger *zap.Logger, primaryResource bool) error {
	resourceType := testData.ResourceType
	var tfplan2caiSupported, cai2hclSupported bool
	if _, tfplan2caiSupported = tfplan2caiconverters.ConverterMap[resourceType]; !tfplan2caiSupported {
		log.Printf("%s is not supported in tfplan2cai conversion.", resourceType)
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
			log.Printf("Test for %s is skipped as conversion of the resource is not supported in tgc.", resourceType)
			return nil
		}
	}

	if !(tfplan2caiSupported && cai2hclSupported) {
		return fmt.Errorf("resource %s is supported in either tfplan2cai or cai2hcl within tgc, but not in both", resourceType)
	}

	// Uncomment these lines when debugging issues locally
	// assetFile := fmt.Sprintf("%s.json", t.Name())
	// writeJSONFile(assetFile, assets)

	// Step 1: Use cai2hcl to convert export assets into a Terraform configuration (export config).
	// Compare all of the fields in raw config are in export config.

	exportConfigData, err := cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return fmt.Errorf("error when converting the export assets into export config: %#v", err)
	}

	// Uncomment these lines when debugging issues locally
	// exportTfFile := fmt.Sprintf("%s_export.tf", t.Name())
	// err = os.WriteFile(exportTfFile, exportConfigData, 0644)
	// if err != nil {
	// 	return fmt.Errorf("error writing file %s", exportTfFile)
	// }
	// defer os.Remove(exportTfFile)

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

	ignoredFieldMap := make(map[string]bool, 0)
	for _, f := range ignoredFields {
		ignoredFieldMap[f] = true
	}

	parsedExportConfig := exportResources[0].Attributes
	missingKeys := compareHCLFields(testData.ParsedRawConfig, parsedExportConfig, "", ignoredFieldMap)
	if len(missingKeys) > 0 {
		return fmt.Errorf("missing fields in address %s after cai2hcl conversion:\n%s", testData.ResourceAddress, missingKeys)
	}
	log.Printf("Step 1 passes for resource %s. All of the fields in raw config are in export config", testData.ResourceAddress)

	// Step 2
	// Run a terraform plan using export_config.
	// Use tfplan2cai to convert the generated plan into CAI assets (roundtrip_assets).
	// Convert roundtrip_assets back into a Terraform configuration (roundtrip_config) using cai2hcl.
	// Compare roundtrip_config with export_config to ensure they are identical.

	// Convert the export config to roundtrip assets and then convert the roundtrip assets back to roundtrip config
	ancestryCache := getAncestryCache(assets)
	roundtripAssets, roundtripConfigData, err := getRoundtripConfig(t, testName, tfDir, ancestryCache, logger)
	if err != nil {
		return fmt.Errorf("error when converting the round-trip config: %#v", err)
	}

	roundtripTfFilePath := fmt.Sprintf("%s_roundtrip.tf", testName)
	err = os.WriteFile(roundtripTfFilePath, roundtripConfigData, 0644)
	if err != nil {
		return fmt.Errorf("error when writing the file %s", roundtripTfFilePath)
	}
	defer os.Remove(roundtripTfFilePath)

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
				cmpopts.IgnoreFields(caiasset.AssetResource{}, "Version", "Data"),
				// Consider DiscoveryDocumentURI equal if they have the same number of path segments when split by "/".
				cmp.FilterPath(func(p cmp.Path) bool {
					return p.Last().String() == ".DiscoveryDocumentURI"
				}, cmp.Comparer(func(x, y string) bool {
					parts1 := strings.Split(x, "/")
					parts2 := strings.Split(y, "/")
					return len(parts1) == len(parts2)
				})),
			); diff != "" {
				return fmt.Errorf("differences found between exported asset and roundtrip asset (-want +got):\n%s", diff)
			}
		}
	}
	log.Printf("Step 3 passes for resource %s. Exported asset and roundtrip asset are identical", testData.ResourceAddress)

	return nil
}

// Gets the ancestry cache for tfplan2cai conversion
func getAncestryCache(assets []caiasset.Asset) map[string]string {
	ancestryCache := make(map[string]string, 0)

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
			}

			project := utils.ParseFieldValue(asset.Name, "projects")
			projectKey := fmt.Sprintf("projects/%s", project)
			if strings.HasPrefix(ancestors[0], "projects") && ancestors[0] != projectKey {
				if _, ok := ancestryCache[projectKey]; !ok {
					ancestryCache[projectKey] = path
				}
			}
		}
	}
	return ancestryCache
}

// Compares HCL and finds all of the keys in map1 are in map2
func compareHCLFields(map1, map2 map[string]interface{}, path string, ignoredFields map[string]bool) []string {
	var missingKeys []string
	for key, value1 := range map1 {
		if value1 == nil {
			continue
		}

		currentPath := path + "." + key
		if path == "" {
			currentPath = key
		}

		if ignoredFields[currentPath] {
			continue
		}

		value2, ok := map2[key]
		if !ok || value2 == nil {
			missingKeys = append(missingKeys, currentPath)
			continue
		}

		switch v1 := value1.(type) {
		case map[string]interface{}:
			v2, _ := value2.(map[string]interface{})
			missingKeys = append(missingKeys, compareHCLFields(v1, v2, currentPath, ignoredFields)...)
		case []interface{}:
			v2, _ := value2.([]interface{})

			for i := 0; i < len(v1); i++ {
				nestedMap1, ok1 := v1[i].(map[string]interface{})
				nestedMap2, ok2 := v2[i].(map[string]interface{})
				if ok1 && ok2 {
					keys := compareHCLFields(nestedMap1, nestedMap2, fmt.Sprintf("%s[%d]", currentPath, i), ignoredFields)
					missingKeys = append(missingKeys, keys...)
				}
			}
		default:
		}
	}

	return missingKeys
}

// Converts a tfplan to CAI asset, and then converts the CAI asset into HCL
func getRoundtripConfig(t *testing.T, testName string, tfDir string, ancestryCache map[string]string, logger *zap.Logger) ([]caiasset.Asset, []byte, error) {
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
		DefaultProject: "ci-test-project-nightly-beta",
		DefaultRegion:  "",
		DefaultZone:    "",
		UserAgent:      "",
		AncestryCache:  ancestryCache,
	})

	if err != nil {
		return nil, nil, err
	}

	// Uncomment these lines when debugging issues locally
	// roundtripAssetFile := fmt.Sprintf("%s_roundtrip.json", t.Name())
	// writeJSONFile(roundtripAssetFile, roundtripAssets)

	roundtripConfig, err := cai2hcl.Convert(roundtripAssets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return nil, nil, err
	}

	return roundtripAssets, roundtripConfig, nil
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
