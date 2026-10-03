package test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/tfplan"
	tfjson "github.com/hashicorp/terraform-json"
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
			var attemptErrors []error
			flakyAction := func(ctx context.Context) error {
				testData, err := prepareTestData(subTestName, stepN, retries)
				retries++
				log.Printf("%s: Starting the attempt %d", tName, retries)
				if err != nil {
					err = fmt.Errorf("%s: error preparing the input data: %v", tName, err)
					attemptErrors = append(attemptErrors, err)
					return err
				}

				if testData == nil {
					err = retry.RetryableError(fmt.Errorf("fail: test data is unavailable"))
					attemptErrors = append(attemptErrors, err)
					return err
				}

				// If the primary resource is specified, only test the primary resource.
				// Otherwise, test all of the resources in the test.
				primaryResource := testData.PrimaryResource
				resourceTestData := testData.ResourceTestData
				if primaryResource != "" {
					t.Logf("%s: Test for the primary resource %s begins.", tName, primaryResource)
					err = testSingleResource(t, tName, resourceTestData[primaryResource], tfDir, ignoredFields, logger, true)
					if err != nil {
						attemptErrors = append(attemptErrors, err)
						return err
					}
				} else {
					for address, testData := range resourceTestData {
						if !strings.HasPrefix(address, primaryResourceType) {
							continue
						}
						err = testSingleResource(t, tName, testData, tfDir, ignoredFields, logger, false)
						if err != nil {
							attemptErrors = append(attemptErrors, err)
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
				allUnavailable := len(attemptErrors) > 0
				for _, e := range attemptErrors {
					if !strings.Contains(e.Error(), "test data is unavailable") {
						allUnavailable = false
						break
					}
				}

				if allUnavailable {
					t.Skipf("%s: Test skipped because data was unavailable after all %d attempts: %v", tName, len(attemptErrors), err)
				} else {
					var firstRealError error
					for _, e := range attemptErrors {
						if !strings.Contains(e.Error(), "test data is unavailable") {
							firstRealError = e
							break
						}
					}
					t.Fatalf("%s: Failed after %d attempts. First real error: %v", tName, len(attemptErrors), firstRealError)
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
		return retry.RetryableError(fmt.Errorf("fail: test data is unavailable"))
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
		return fmt.Errorf("conversion of the resource %s is not supported in tgc", testData.ResourceAddress)
	}

	if !(tfplan2caiSupported && cai2hclSupported) {
		return fmt.Errorf("resource %s is supported in either tfplan2cai or cai2hcl within tgc, but not in both", resourceType)
	}

	if os.Getenv("WRITE_FILES") != "" {
		assetFile := fmt.Sprintf("%s.json", strings.ReplaceAll(testName, "/", "_"))
		writeJSONFile(assetFile, assets)
	}

	// Step 1: Generate plan JSON for both raw and export config, and compare resource changes.
	rawTfDir := filepath.Join(tfDir, "raw")
	if err := os.MkdirAll(rawTfDir, 0755); err != nil {
		return fmt.Errorf("error creating raw tf directory: %w", err)
	}
	rawTfFile := filepath.Join(rawTfDir, "main.tf")
	if err := os.WriteFile(rawTfFile, []byte(testData.RawConfig), 0644); err != nil {
		return fmt.Errorf("error writing raw tf file: %w", err)
	}
	if err := appendProviderOverride(rawTfFile); err != nil {
		return err
	}

	ancestryCache, defaultProject := getAncestryCache(assets)

	if err := terraformWorkflow(rawTfDir, "raw", defaultProject); err != nil {
		return fmt.Errorf("error running terraform on raw config: %w", err)
	}
	rawPlanPath := filepath.Join(rawTfDir, "raw.tfplan.json")
	rawPlanData, err := os.ReadFile(rawPlanPath)
	if err != nil {
		return fmt.Errorf("error reading raw plan JSON: %w", err)
	}
	rawChanges, err := tfplan.ReadResourceChanges(rawPlanData)
	if err != nil {
		return fmt.Errorf("error reading raw resource changes: %w", err)
	}

	exportConfigData, err := cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return fmt.Errorf("error when converting the export assets into export config: %v", err)
	}

	if os.Getenv("WRITE_FILES") != "" {
		exportTfFile := fmt.Sprintf("%s_export.tf", strings.ReplaceAll(testName, "/", "_"))
		err = os.WriteFile(exportTfFile, exportConfigData, 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s", exportTfFile)
		}
	}

	exportTfFilePath := filepath.Join(tfDir, "main.tf")
	if err := os.WriteFile(exportTfFilePath, exportConfigData, 0644); err != nil {
		return fmt.Errorf("error writing export tf file: %w", err)
	}
	if err := appendProviderOverride(exportTfFilePath); err != nil {
		return err
	}

	if err := terraformWorkflow(tfDir, "export", defaultProject); err != nil {
		return fmt.Errorf("error running terraform on export config: %w", err)
	}
	exportPlanPath := filepath.Join(tfDir, "export.tfplan.json")
	exportPlanData, err := os.ReadFile(exportPlanPath)
	if err != nil {
		return fmt.Errorf("error reading export plan JSON: %w", err)
	}
	exportChanges, err := tfplan.ReadResourceChanges(exportPlanData)
	if err != nil {
		return fmt.Errorf("error reading export resource changes: %w", err)
	}

	var rawConfig any
	for _, rc := range rawChanges {
		if rc.Address == testData.ResourceAddress && (tfplan.IsCreate(rc) || tfplan.IsUpdate(rc) || tfplan.IsDeleteCreate(rc)) {
			rawConfig = rc.Change.After
			break
		}
	}

	if rawConfig == nil {
		return fmt.Errorf("target resource address %s is missing in raw config plan changes", testData.ResourceAddress)
	}

	exportConfig, err := findExportConfig(rawConfig.(map[string]any), resourceType, exportChanges)
	if err != nil {
		return fmt.Errorf("target resource type %s: error finding export config: %w", resourceType, err)
	}

	ignoredFieldSet := make(map[string]any, 0)
	for _, f := range ignoredFields {
		ignoredFieldSet[f] = struct{}{}
	}

	missingKeys := findMissingKeys(rawConfig.(map[string]any), exportConfig.(map[string]any), "", ignoredFieldSet)

	if len(missingKeys) > 0 {
		log.Printf("%s: missing fields in resource %s after cai2hcl conversion:\n%s", testName, testData.ResourceAddress, missingKeys)
		return retry.RetryableError(fmt.Errorf("missing fields: %v", missingKeys))
	}
	log.Printf("%s: Step 1 passes for resource %s. All of the fields in raw config are in export config", testName, testData.ResourceAddress)

	// Step 2
	// Run a terraform plan using export_config.
	// Use tfplan2cai to convert the generated plan into CAI assets (roundtrip_assets).
	// Convert roundtrip_assets back into a Terraform configuration (roundtrip_config) using cai2hcl.
	// Compare roundtrip_config with export_config to ensure they are identical.

	// Convert the export config to roundtrip assets and then convert the roundtrip assets back to roundtrip config
	roundtripAssets, roundtripConfigData, err := getRoundtripConfig(t, testName, tfDir, ancestryCache, defaultProject, logger)
	if err != nil {
		return retry.RetryableError(fmt.Errorf("error when converting the round-trip config: %v", err))
	}

	rtTfFile := fmt.Sprintf("%s_roundtrip.tf", strings.ReplaceAll(testName, "/", "_"))
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
		tfFileName := fmt.Sprintf("%s_roundtrip", strings.ReplaceAll(testName, "/", "_"))
		jsonFileName := fmt.Sprintf("%s_reexport.json", strings.ReplaceAll(testName, "/", "_"))

		reexportAssets, err := tfplan2caiConvert(t, tfFileName, jsonFileName, tfDir, ancestryCache, defaultProject, logger)
		if err != nil {
			return fmt.Errorf("error when converting the third round-trip config: %v", err)
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
	resolvedProject := ""

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
				if resolvedProject == "" {
					if s, hasPrefix := strings.CutPrefix(ancestors[0], "projects/"); hasPrefix {
						resolvedProject = s
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

				if resolvedProject == "" {
					resolvedProject = project
				}
			}
		}
	}
	if resolvedProject == "" {
		resolvedProject = defaultProject
	}
	ancestryCache["projects/"+defaultProject] = "organizations/" + defaultOrganization

	return ancestryCache, resolvedProject
}

// findMissingKeys recursively finds all of the keys in map1 that are not in map2
func findMissingKeys(map1, map2 map[string]any, path string, ignoredFields map[string]any) []string {
	var missingKeys []string
	for key, value1 := range map1 {
		if value1 == nil {
			continue
		}
		currentPath := key
		if path != "" {
			currentPath = path + "." + key
		}
		if isIgnoredPath(currentPath, ignoredFields) {
			continue
		}
		value2, ok := map2[key]
		if !ok || value2 == nil {
			missingKeys = append(missingKeys, currentPath)
			continue
		}
		switch v1 := value1.(type) {
		case map[string]any:
			v2, ok := value2.(map[string]any)
			if ok {
				missingKeys = append(missingKeys, findMissingKeys(v1, v2, currentPath, ignoredFields)...)
			}
		case []any:
			v2, ok := value2.([]any)
			if ok {
				for i := 0; i < len(v1); i++ {
					if i < len(v2) {
						nMap1, ok1 := v1[i].(map[string]any)
						nMap2, ok2 := v2[i].(map[string]any)
						if ok1 && ok2 {
							missingKeys = append(missingKeys, findMissingKeys(nMap1, nMap2, fmt.Sprintf("%s.%d", currentPath, i), ignoredFields)...)
						}
					}
				}
			}
		}
	}
	return missingKeys
}

// isIgnoredPath returns true if the given key path is ignored.
func isIgnoredPath(key string, ignoredFields map[string]any) bool {
	// Global ignores for write-only fields
	if strings.HasSuffix(key, "_wo") || strings.HasSuffix(key, "_wo_version") {
		return true
	}

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
	fileName := fmt.Sprintf("%s_export", strings.ReplaceAll(testName, "/", "_"))
	roundtripAssetFile := fmt.Sprintf("%s_roundtrip.json", strings.ReplaceAll(testName, "/", "_"))

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
	if err := terraformWorkflow(tfDir, tfFileName, defaultProject); err != nil {
		return nil, err
	}

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

	for i := 0; i < len(parts2); i++ {
		if parts2[i] == "null" || parts2[i] == "unknown" {
			return nil
		}
		if i >= len(parts1) {
			return fmt.Errorf("differences found between two asset names: want %s, got %s", want, got)
		}
		if parts1[i] != parts2[i] {
			return fmt.Errorf("differences found between two asset names: want %s, got %s", want, got)
		}
	}

	if len(parts1) != len(parts2) {
		return fmt.Errorf("differences found between two asset names: want %s, got %s", want, got)
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
				cmpopts.IgnoreFields(caiasset.AssetResource{}, "Version", "Data", "Location", "Parent", "DiscoveryName", "DiscoveryDocumentURI"),
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

func findExportConfig(rawConfig map[string]any, resourceType string, exportChanges []*tfjson.ResourceChange) (any, error) {
	var candidateChanges []*tfjson.ResourceChange
	for _, ec := range exportChanges {
		if ec.Type == resourceType && (tfplan.IsCreate(ec) || tfplan.IsUpdate(ec) || tfplan.IsDeleteCreate(ec)) {
			candidateChanges = append(candidateChanges, ec)
		}
	}

	if len(candidateChanges) == 0 {
		return nil, fmt.Errorf("no export plan changes found for resource type %s", resourceType)
	}

	if len(candidateChanges) == 1 {
		return candidateChanges[0].Change.After, nil
	}

	for _, ec := range candidateChanges {
		parts := strings.Split(ec.Address, ".")
		if len(parts) < 2 {
			continue
		}
		resName := parts[len(parts)-1]

		if isValueInMap(rawConfig, resName) {
			return ec.Change.After, nil
		}
	}

	return candidateChanges[0].Change.After, nil
}

func isValueInMap(m map[string]any, target string) bool {
	for _, val := range m {
		if s, ok := val.(string); ok && s == target {
			return true
		}
		if subMap, ok := val.(map[string]any); ok {
			if isValueInMap(subMap, target) {
				return true
			}
		}
		if list, ok := val.([]any); ok {
			for _, item := range list {
				if subMap, ok := item.(map[string]any); ok {
					if isValueInMap(subMap, target) {
						return true
					}
				}
				if s, ok := item.(string); ok && s == target {
					return true
				}
			}
		}
	}
	return false
}

func appendProviderOverride(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s for provider override: %w", filePath, err)
	}
	override := `
terraform {
  required_providers {
    google = {
      source = "hashicorp/google-beta"
    }
  }
}
`
	if err := os.WriteFile(filePath, append(content, []byte(override)...), 0644); err != nil {
		return fmt.Errorf("error writing provider override to %s: %w", filePath, err)
	}
	return nil
}


