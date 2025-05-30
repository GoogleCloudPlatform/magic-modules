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

	"go.uber.org/zap/zaptest"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

var (
	cacheMutex = sync.Mutex{}
	tmpDir     = os.TempDir()
)

func AssertTestFile(t *testing.T, ignoredFields []string) {
	var err error
	cacheMutex.Lock()
	TestsMetadata, err = ReadTestsDataFromGcs()
	if err != nil {
		log.Fatal(err)
	}
	cacheMutex.Unlock()

	// Create a temporary directory for running terraform.
	tfDir, err := os.MkdirTemp(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tfDir)

	testMetadata := TestsMetadata[t.Name()]
	resourceMetadata := testMetadata.ResourceMetadata
	if len(resourceMetadata) == 0 {
		log.Printf("Data of test is unavailable: %s", t.Name())
		return
	}

	rawTfFile := fmt.Sprintf("%s.tf", t.Name())
	err = os.WriteFile(rawTfFile, []byte(testMetadata.RawConfig), 0644)
	if err != nil {
		log.Fatalf("error writing to file %s: %#v", rawTfFile, err)
	}
	defer os.Remove(rawTfFile)

	rawResourceConfigs, err := parseResourceConfigs(rawTfFile)
	if err != nil {
		log.Fatalf("error parsing resource configs: %#v", err)
	}

	if len(rawResourceConfigs) == 0 {
		log.Fatalf("Test %s fails: raw config is unavailable", t.Name())
	}

	rawConfigMap := convertToConfigMap(rawResourceConfigs)

	// If the primary resource is available, only test the primary resource.
	// Otherwise, test all of the resources in the test.
	if testMetadata.PrimaryResource != "" {
		primaryMetadata := resourceMetadata[testMetadata.PrimaryResource]
		address := primaryMetadata.ResourceAddress
		testSingleResource(t, primaryMetadata, rawConfigMap[address], tfDir, ignoredFields)
	} else {
		for address, metadata := range resourceMetadata {
			testSingleResource(t, metadata, rawConfigMap[address], tfDir, ignoredFields)
		}
	}
}

// Tests a single resource
func testSingleResource(t *testing.T, resourceMetadata *ResourceMetadata, rawConfig map[string]interface{}, tfDir string, ignoredFields []string) {
	resourceType := resourceMetadata.ResourceType
	if _, ok := tfplan2caiconverters.ConverterMap[resourceType]; !ok {
		log.Printf("Test %s fails as tfplan2cai conversion for %s is not supported.", t.Name(), resourceType)
		return
	}

	assetType := resourceMetadata.CaiAssetData.Type
	if assetType == "" {
		log.Fatalf("Test %s fails: export cai asset is unavailable for %s", t.Name(), resourceMetadata.CaiAssetName)
	}
	if _, ok := cai2hclconverters.ConverterMap[assetType]; !ok {
		log.Printf("Test %s fails as cai2hcl conversion for %s is not supported.", t.Name(), assetType)
		return
	}

	assets := []caiasset.Asset{resourceMetadata.CaiAssetData}

	// Uncomment these lines when debugging issues locally
	// assetFile := fmt.Sprintf("%s.json", t.Name())
	// writeJSONFile(assetFile, assets)

	// Step 1: Use cai2hcl to convert export assets into a Terraform configuration (export config).
	// Compare all of the fields in raw config are in export config.

	exportConfigData, err := cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: zaptest.NewLogger(t),
	})
	if err != nil {
		log.Fatalf("error when converting the export assets into export config: %#v", err)
	}

	// Uncomment these lines when debugging issues locally
	// exportTfFile := fmt.Sprintf("%s_export.tf", t.Name())
	// err = os.WriteFile(exportTfFile, exportConfigData, 0644)
	// if err != nil {
	// 	log.Fatalf("error writing file %s", exportTfFile)
	// }
	// defer os.Remove(exportTfFile)

	exportTfFilePath := fmt.Sprintf("%s/%s_export.tf", tfDir, t.Name())
	err = os.WriteFile(exportTfFilePath, exportConfigData, 0644)
	if err != nil {
		log.Fatalf("error when writing the file %s", exportTfFilePath)
	}

	exportResources, err := parseResourceConfigs(exportTfFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if len(exportResources) == 0 {
		log.Fatalf("%s - Missing hcl after cai2hcl conversion for CAI asset %s.", t.Name(), resourceMetadata.CaiAssetName)
	}

	ignoredFieldMap := make(map[string]bool, 0)
	for _, f := range ignoredFields {
		ignoredFieldMap[f] = true
	}

	exportConfig := exportResources[0].Attributes
	missingKeys := compareHCLFields(rawConfig, exportConfig, "", ignoredFieldMap)
	if len(missingKeys) > 0 {
		log.Fatalf("%s - Missing fields in address %s after cai2hcl conversion:\n%s", t.Name(), resourceMetadata.ResourceAddress, missingKeys)
	}

	// Step 2
	// Run a terraform plan using export_config.
	// Use tfplan2cai to convert the generated plan into CAI assets (roundtrip_assets).
	// Convert roundtrip_assets back into a Terraform configuration (roundtrip_config) using cai2hcl.
	// Compare roundtrip_config with export_config to ensure they are identical.

	// Convert the export config to roundtrip assets and then convert the roundtrip assets back to roundtrip config
	ancestryCache := getAncestryCache(assets)
	roundtripConfigData, err := getRoundtripConfig(t, tfDir, ancestryCache)
	if err != nil {
		log.Fatalf("error when converting the round-trip config: %#v", err)
	}

	roundtripTfFilePath := fmt.Sprintf("%s_roundtrip.tf", t.Name())
	err = os.WriteFile(roundtripTfFilePath, roundtripConfigData, 0644)
	if err != nil {
		log.Fatalf("error when writing the file %s", roundtripTfFilePath)
	}
	defer os.Remove(roundtripTfFilePath)

	if diff := cmp.Diff(string(roundtripConfigData), string(exportConfigData)); diff != "" {
		log.Printf("Roundtrip config is different from the export config.\nroundtrip config:\n%s\nexport config:\n%s", string(roundtripConfigData), string(exportConfigData))
		log.Fatalf("Test %s got diff (-want +got): %s", t.Name(), diff)
	}
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

type Resource struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Attributes map[string]interface{} `json:"attributes"`
}

// parseHCLBody recursively parses attributes and nested blocks from an HCL body.
func parseHCLBody(body hcl.Body, filePath string) (
	attributes map[string]interface{},
	diags hcl.Diagnostics,
) {
	attributes = make(map[string]interface{})
	var allDiags hcl.Diagnostics

	if syntaxBody, ok := body.(*hclsyntax.Body); ok {
		for _, attr := range syntaxBody.Attributes {
			attributes[attr.Name] = true
		}

		for _, block := range syntaxBody.Blocks {
			nestedAttr, diags := parseHCLBody(block.Body, filePath)
			if diags.HasErrors() {
				allDiags = append(allDiags, diags...)
			}

			attributes[block.Type] = nestedAttr
		}
	} else {
		allDiags = append(allDiags, &hcl.Diagnostic{
			Severity: hcl.DiagWarning,
			Summary:  "Body type assertion to *hclsyntax.Body failed",
			Detail:   fmt.Sprintf("Cannot directly parse attributes for body of type %T. Attribute parsing may be incomplete.", body),
		})
	}

	return attributes, allDiags
}

// Parses a Terraform configuation file written with HCL
func parseResourceConfigs(filePath string) ([]Resource, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", filePath, err)
	}

	parser := hclparse.NewParser()
	hclFile, diags := parser.ParseHCL(src, filePath)
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse HCL: %w", diags)
	}

	if hclFile == nil {
		return nil, fmt.Errorf("parsed HCL file %s is nil cannot proceed", filePath)
	}

	var allParsedResources []Resource

	for _, block := range hclFile.Body.(*hclsyntax.Body).Blocks {
		if block.Type == "resource" {
			if len(block.Labels) != 2 {
				log.Printf("Skipping address block with unexpected number of labels: %v", block.Labels)
				continue
			}

			resType := block.Labels[0]
			resName := block.Labels[1]
			attrs, procDiags := parseHCLBody(block.Body, filePath)

			if procDiags.HasErrors() {
				log.Printf("Diagnostics while processing address %s.%s body in %s:", resType, resName, filePath)
				for _, diag := range procDiags {
					log.Printf("  - %s (Severity)", diag.Error())
				}
			}

			gr := Resource{
				Type:       resType,
				Name:       resName,
				Attributes: attrs,
			}
			allParsedResources = append(allParsedResources, gr)
		}
	}

	return allParsedResources, nil
}

// Converts the slice to map with resource address as the key
func convertToConfigMap(resources []Resource) map[string]map[string]interface{} {
	configMap := make(map[string]map[string]interface{}, 0)

	for _, r := range resources {
		addr := fmt.Sprintf("%s.%s", r.Type, r.Name)
		configMap[addr] = r.Attributes
	}

	return configMap
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
func getRoundtripConfig(t *testing.T, tfDir string, ancestryCache map[string]string) ([]byte, error) {
	fileName := fmt.Sprintf("%s_export", t.Name())

	// Run terraform init and terraform apply to generate tfplan.json files
	terraformWorkflow(t, tfDir, fileName)

	planFile := fmt.Sprintf("%s.tfplan.json", fileName)
	planfilePath := filepath.Join(tfDir, planFile)
	jsonPlan, err := os.ReadFile(planfilePath)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	roundtripAssets, err := tfplan2cai.Convert(ctx, jsonPlan, &tfplan2cai.Options{
		ErrorLogger:    zaptest.NewLogger(t),
		Offline:        true,
		DefaultProject: "ci-test-project-nightly-beta",
		DefaultRegion:  "",
		DefaultZone:    "",
		UserAgent:      "",
		AncestryCache:  ancestryCache,
	})

	if err != nil {
		return nil, err
	}

	// Uncomment these lines when debugging issues locally
	// roundtripAssetFile := fmt.Sprintf("%s_roundtrip.json", t.Name())
	// writeJSONFile(roundtripAssetFile, roundtripAssets)

	data, err := cai2hcl.Convert(roundtripAssets, &cai2hcl.Options{
		ErrorLogger: zaptest.NewLogger(t),
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
