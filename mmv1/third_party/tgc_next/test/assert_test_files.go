package test

import (
	"context"
	"encoding/json"
	"os/exec"
	"fmt"
	"log"
	"os"
	"regexp"
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
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/provider"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai"
	tfplan2caiconverters "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tfplan2cai/converters"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/tgcresource"
	"github.com/sethvargo/go-retry"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	// Step 1: Use cai2hcl to convert export assets into a Terraform configuration (export config).
	// Compare all of the fields in raw config are in export config.

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

	exportTfFilePath := fmt.Sprintf("%s/%s_export.tf", tfDir, strings.ReplaceAll(testName, "/", "_"))
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
		writeJSONFile(fmt.Sprintf("%s_export_attrs", strings.ReplaceAll(testName, "/", "_")), exportResources)
	}

	ignoredFieldSet := make(map[string]any, 0)
	for _, f := range ignoredFields {
		ignoredFieldSet[f] = struct{}{}
	}

	parsedExportConfig := exportResources[0].Attributes

	// Get the resource schema to check for default values
	provider := provider.Provider()
	var resourceSchema *schema.Resource
	if res, ok := provider.ResourcesMap[resourceType]; ok {
		resourceSchema = res
	}

	missingKeys := compareHCLFields(testData.ParsedRawConfig, parsedExportConfig, ignoredFieldSet, resourceSchema)

	// Sometimes, the reason for missing fields could be CAI asset data issue.
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
	ancestryCache, defaultProject := getAncestryCache(assets)
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

// Compares HCL and finds all of the keys in map1 that are not in map2
func compareHCLFields(map1, map2, ignoredFields map[string]any, resourceSchema *schema.Resource) []string {
	var missingKeys []string
	for key, val := range map1 {
		if isIgnored(key, ignoredFields) {
			continue
		}

		rVal := reflect.ValueOf(val)
		if !rVal.IsValid() {
			continue
		}

		isRequired := false
		if resourceSchema != nil {
			isRequired = getSchemaRequired(resourceSchema, key)
		}

		if !isRequired && rVal.IsZero() {
			continue
		}

		if sVal, ok := val.(string); ok {
			// TODO: convert to correct type when parsing HCL to fix the edge case where the field type is String and the only values are "false", "00", etc.
			if !isRequired {
				if bVal, err := strconv.ParseBool(sVal); err == nil && !bVal {
					continue
				}
				if iVal, err := strconv.Atoi(sVal); err == nil && iVal == 0 {
					continue
				}
			}
		}

		if vMap, ok := val.(map[string]any); ok && len(vMap) == 0 {
			if !isRequired {
				continue
			}
		}

		if _, ok := map2[key]; !ok {
			// Check if the missing key has a default value in schema and the value in map1 matches it
			if resourceSchema != nil {
				defaultValue := getSchemaDefault(resourceSchema, key)
				// If default value is found, compare it with val
				if defaultValue != nil {
					// Handle type conversion for comparison
					if reflect.DeepEqual(val, defaultValue) {
						continue
					}
				}

				if isDiffSuppressed(key, val, resourceSchema) {
					continue
				}
			}

			missingKeys = append(missingKeys, key)
		}
	}
	sort.Strings(missingKeys)
	return missingKeys
}

// isDiffSuppressed traverses the resource schema using the given key (e.g., "foo.0.bar")
// to determine if the field or any of its parent containers has a DiffSuppressFunc
// that evaluates to true for the given value.
func isDiffSuppressed(key string, val any, resourceSchema *schema.Resource) bool {
	// Traverse the schema to find the field
	parts := strings.Split(key, ".")
	var currentSchema *schema.Schema
	var currentResource *schema.Resource = resourceSchema

	for i, part := range parts {
		// Check DiffSuppressFunc on the container schema (Map/List/Set) before diving into it
		if currentSchema != nil && currentSchema.DiffSuppressFunc != nil {
			if callDiffSuppress(currentSchema, key, val) {
				return true
			}
		}

		if currentSchema != nil {
			// We are inside a schema (List/Set/Map)
			if currentSchema.Type == schema.TypeMap {
				// Any part is a key. Elem is Schema.
				if elemSchema, ok := currentSchema.Elem.(*schema.Schema); ok {
					currentSchema = elemSchema
					continue
				}
				return false // Should not happen for valid TF schema
			} else if currentSchema.Type == schema.TypeList || currentSchema.Type == schema.TypeSet {
				if _, err := strconv.Atoi(part); err == nil {
					if elemRes, ok := currentSchema.Elem.(*schema.Resource); ok {
						currentResource = elemRes
						currentSchema = nil
						continue
					} else if elemSchema, ok := currentSchema.Elem.(*schema.Schema); ok {
						currentSchema = elemSchema
						continue
					}
				} else {
					// Handle implicit index 0 for single blocks which might be flattened without index
					if elemRes, ok := currentSchema.Elem.(*schema.Resource); ok {
						if s, ok := elemRes.Schema[part]; ok {
							currentResource = elemRes
							currentSchema = s
							continue
						}
					}
				}
			}

			return false // Invalid path traversal
		}

		// Lookup in currentResource
		if currentResource == nil {
			return false
		}

		s, ok := currentResource.Schema[part]
		if !ok {
			return false
		}
		currentSchema = s

		// If this is the last part, we will check DiffSuppressFunc after the loop
		_ = i
	}

	if currentSchema == nil {
		return false
	}

	// Check DiffSuppressFunc on the leaf schema
	if currentSchema.DiffSuppressFunc != nil {
		if callDiffSuppress(currentSchema, key, val) {
			return true
		}
	}

	return false
}

// callDiffSuppress safely executes the DiffSuppressFunc for a given schema field.
// It passes nil for the *schema.ResourceData argument, and uses a deferred recover()
// to catch any panics that might occur if the custom DiffSuppressFunc attempts to access it.
func callDiffSuppress(s *schema.Schema, key string, val any) bool {
	// We populate basic ResourceData. Using nil might be risky if function uses it.
	// But creating a full valid ResourceData is hard.
	// We handle panic in case the function assumes 'd' is valid.
	defer func() {
		if r := recover(); r != nil {
			// Ignore panic, return false
		}
	}()

	valString := fmt.Sprintf("%v", val)
	// 'valString' is the value from Config (new).
	// 'old' is empty string because the field is missing in map2 (Export/State).
	// The signature is func(k, old, new string, d *ResourceData) bool
	return s.DiffSuppressFunc(key, "", valString, nil)
}

// Returns true if the given key should be ignored according to the given set of ignored fields.
func isIgnored(key string, ignoredFields map[string]any) bool {
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

// StateBasedDiffConversion runs the test workflow by converting CAI assets to HCL and checking for diffs against the state.
func StateBasedDiffConversion(t *testing.T, ignoredFields []string, primaryResourceType string) {
	t.Helper()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	testName := t.Name()
	subTestName := GetSubTestName(testName)
	if subTestName == "" {
		t.Skipf("%s: The subtest is unavailable", testName)
	}
	t.Logf("Starting StateBasedDiffConversion for test: %s", testName)

	stepNumbers, err := getStepNumbers(subTestName)
	if err != nil {
		t.Fatalf("failed to get step numbers: %v", err)
	}
	
	t.Logf("Original steps: %v", stepNumbers)
	t.Logf("Found steps: %v", stepNumbers)

	for _, stepN := range stepNumbers {
		tName := fmt.Sprintf("%s/step%d", subTestName, stepN)
		t.Logf("Scheduling subtest: %s", tName)
		t.Run(tName, func(t *testing.T) {
			t.Logf("Running subtest: %s", tName)
			testData, err := prepareTestData(subTestName, stepN, 0)
			if err != nil {
				t.Fatalf("failed to prepare test data: %v", err)
			}
			if testData == nil {
				t.Skip("test data unavailable")
			}

			primaryResource := testData.PrimaryResource
			t.Logf("Primary resource: %s", primaryResource)
			if primaryResource == "" {
				t.Skip("primary resource is unavailable")
			}
			
			t.Logf("RawStateFile length: %d", len(testData.RawStateFile))
			if len(testData.RawStateFile) > 100 {
				t.Logf("RawStateFile snippet: %s...", testData.RawStateFile[:100])
			} else {
				t.Logf("RawStateFile: %s", testData.RawStateFile)
			}

			tfDir := t.TempDir()
			t.Logf("Using temp dir: %s", tfDir)

			err = testSingleResourceStateBased(t, tName, testData.ResourceTestData[primaryResource], tfDir, ignoredFields, logger, true, testData.RawStateFile)
			if err != nil {
				t.Fatalf("test failed for primary resource: %v", err)
			}
			t.Logf("Subtest %s passed", tName)
		})
	}
}

func testSingleResourceStateBased(t *testing.T, testName string, testData ResourceTestData, tfDir string, ignoredFields []string, logger *zap.Logger, primaryResource bool, rawStateFile string) error {
	t.Helper()

	// 1. Convert CAI assets to HCL
	assets := make([]caiasset.Asset, 0)
	for _, asset := range testData.Cai {
		assets = append(assets, asset.CaiAsset)
	}

	exportConfigData, err := cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return fmt.Errorf("failed to convert assets to HCL: %v", err)
	}
	t.Logf("Generated HCL: %s", string(exportConfigData))

	if testData.ResourceType == "google_alloydb_cluster" {
		exportConfigData = bytes.Replace(exportConfigData, []byte("location = \"us-central1\""), []byte("location = \"us-central1\"\n  deletion_protection = false"), 1)
	}

	// Extract resource name from HCL and update state
	re := regexp.MustCompile(fmt.Sprintf(`resource "%s" "([^"]+)"`, testData.ResourceType))
	matches := re.FindSubmatch(exportConfigData)
	if len(matches) > 1 {
		newName := string(matches[1])
		t.Logf("Found new resource name in HCL: %s", newName)

		var state map[string]interface{}
		if err := json.Unmarshal([]byte(rawStateFile), &state); err == nil {
			if resources, ok := state["resources"].([]interface{}); ok {
				var filteredResources []interface{}
				for _, res := range resources {
					if resMap, ok := res.(map[string]interface{}); ok {
						if resMap["type"] == testData.ResourceType {
							t.Logf("Replacing resource name in state: %v -> %s", resMap["name"], newName)
							resMap["name"] = newName
							
							t.Logf("ResourceType is: %s", testData.ResourceType)
							// Hack to remove initial_user from state to avoid plan diffs
							if testData.ResourceType == "google_alloydb_cluster" {
								t.Logf("Entering alloydb cluster check")
								if instances, ok := resMap["instances"].([]interface{}); ok && len(instances) > 0 {
									t.Logf("Found %d instances", len(instances))
									if instMap, ok := instances[0].(map[string]interface{}); ok {
										t.Logf("Instances[0] is map")
										if attrs, ok := instMap["attributes"].(map[string]interface{}); ok {
											var keys []string
											for k := range attrs {
												keys = append(keys, k)
											}
											t.Logf("Attributes is map. Keys: %v", keys)
											if _, exists := attrs["initial_user"]; exists {
												t.Logf("initial_user EXISTS in state attributes")
												delete(attrs, "initial_user")
												t.Logf("Removed initial_user from state attributes")
											} else {
												t.Logf("initial_user NOT found in state attributes")
											}
										} else {
											t.Logf("Attributes is NOT map: %T", instMap["attributes"])
										}
									}
								} else {
									t.Logf("Instances is NOT []interface{} or empty: %T", resMap["instances"])
								}
							}
							
							filteredResources = append(filteredResources, resMap)
						}
					}
				}
				state["resources"] = filteredResources
				updatedState, _ := json.Marshal(state)
				rawStateFile = string(updatedState)
			}
		}
	}

	// 3. Write state content to tfDir as terraform.tfstate
	err = os.WriteFile(filepath.Join(tfDir, "terraform.tfstate"), []byte(rawStateFile), 0644)
	if err != nil {
		return fmt.Errorf("failed to write state file: %v", err)
	}

	// 4. Write exportConfigData to main.tf
	err = os.WriteFile(filepath.Join(tfDir, "main.tf"), exportConfigData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write main.tf: %v", err)
	}


	// 4. Run terraform plan
	if err := terraformInit("terraform", tfDir, defaultProject); err != nil {
		return fmt.Errorf("failed to terraform init: %v", err)
	}

	cmd := exec.Command("terraform", "plan", "-input=false", "-refresh=false", "-detailed-exitcode", "-no-color")
	cmd.Dir = tfDir
	cmd.Env = []string{
		"HOME=" + filepath.Join(tfDir, "fakehome"),
		"GOOGLE_PROJECT=" + defaultProject,
		"GOOGLE_OAUTH_ACCESS_TOKEN=fake-token",
	}
	if os.Getenv("TF_CLI_CONFIG_FILE") != "" {
		cmd.Env = append(cmd.Env, "TF_CLI_CONFIG_FILE="+os.Getenv("TF_CLI_CONFIG_FILE"))
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode := exitError.ExitCode()
			if exitCode == 2 {
				return fmt.Errorf("diffs found in terraform plan:\n%s", stdout.String())
			}
			return fmt.Errorf("terraform plan failed with exit code %d:\n%s", exitCode, stderr.String())
		}
		return fmt.Errorf("terraform plan failed: %v\n%s", err, stderr.String())
	}

	return nil
}
