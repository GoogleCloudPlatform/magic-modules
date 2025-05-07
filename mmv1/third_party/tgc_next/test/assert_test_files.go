package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tfplan2cai"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func AssertTestFile(t *testing.T, excludedFields []string) {
	err := ReadTestsDataFromGcs()
	if err != nil {
		log.Fatal(err)
	}

	// Create a temporary directory for running terraform.
	tfDir, err := os.MkdirTemp(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tfDir)

	fileName := t.Name()
	testMetadata := TestConfig[fileName]
	resource := testMetadata.Resource

	jsonData, err := json.MarshalIndent(testMetadata.Assets, "", "  ")
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %s", err)
	}
	assetFile := fmt.Sprintf("%s.json", fileName)
	err = ioutil.WriteFile(assetFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON data to file %s: %s", assetFile, err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	exportConfigData, err := cai2hcl.Convert(testMetadata.Assets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	exportFileName := fmt.Sprintf("%s_export", fileName)
	exportTfFile := fmt.Sprintf("%s.tf", exportFileName)
	exportTfFilePath := fmt.Sprintf("%s/%s", tfDir, exportTfFile)
	err = os.WriteFile(exportTfFilePath, exportConfigData, 0644)
	if err != nil {
		log.Fatalf("error writing to file %s: %#v", exportTfFilePath, err)
	}

	// TODO: remove it later
	err = os.WriteFile(exportTfFile, exportConfigData, 0644)
	if err != nil {
		log.Fatalf("error writing to file %s: %#v", exportTfFile, err)
	}

	exportConfigMap, err := parseConfig(exportTfFilePath, resource)
	if err != nil {
		log.Fatal(err)
	}

	if len(exportConfigMap) == 0 {
		log.Fatalf("%s - Missing resource after cai2hcl conversion: %s.", t.Name(), resource)
	}

	rawTfFile := fmt.Sprintf("%s.tf", fileName)
	err = os.WriteFile(rawTfFile, []byte(testMetadata.RawConfig), 0644)
	if err != nil {
		log.Fatalf("error writing to file %s: %#v", rawTfFile, err)
	}

	// defer func() {
	// 	err := os.Remove(rawTfFile)
	// 	if err != nil {
	// 		t.Errorf("Failed to remove file %s: %v", rawTfFile, err)
	// 	}
	// 	t.Logf("Removed temporary file: %s", rawTfFile)
	// }()

	rawConfigMap, err := parseConfig(rawTfFile, resource)
	if err != nil {
		log.Fatal(err)
	}
	if len(rawConfigMap) == 0 {
		log.Printf("Skip test %s: raw config is unavailable", fileName)
		return
	}

	excludedFieldMap := make(map[string]bool, 0)
	for _, f := range excludedFields {
		excludedFieldMap[f] = true
	}

	missingKeys := compareHCLFields(rawConfigMap, exportConfigMap, "", excludedFieldMap)
	if len(missingKeys) > 0 {
		log.Fatalf("%s - Missing fields in resource %s after cai2hcl conversion:\n%s", t.Name(), resource, missingKeys)
	}

	ancestryCache := getAncestryCache(testMetadata.Assets)

	// Convert the export config to roundtrip assets and then convert the roundtrip assets back to roundtrip config
	roundtripConfigData, err := getRoundtripConfig(t, exportFileName, tfDir, ancestryCache)
	if err != nil {
		log.Fatal(err)
	}

	roundtripFileName := fmt.Sprintf("%s_roundtrip", fileName)
	roundtripTfFilePath := fmt.Sprintf("%s.tf", roundtripFileName)
	err = os.WriteFile(roundtripTfFilePath, roundtripConfigData, 0644)
	if err != nil {
		log.Fatalf("error writing to file %s: %#v", roundtripTfFilePath, err)
	}

	if diff := cmp.Diff(string(roundtripConfigData), string(exportConfigData)); diff != "" {
		logger.Debug(fmt.Sprintf("Roundtrip config is different with the export config.\nroundtrip config:\n%s\nexport config:\n%s", string(roundtripConfigData), string(exportConfigData)))

		log.Fatalf("Test %s got diff (-want +got): %s", fileName, diff)
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

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
		{
			Type:       "provider",
			LabelNames: []string{"name"},
		},
	},
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
func parseConfig(filePath, target string) (map[string]interface{}, error) {
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
				log.Printf("Skipping resource block with unexpected number of labels: %v", block.Labels)
				continue
			}

			resType := block.Labels[0]
			resName := block.Labels[1]
			attrs, procDiags := parseHCLBody(block.Body, filePath)

			if procDiags.HasErrors() {
				log.Printf("Diagnostics while processing resource %s.%s body in %s:", resType, resName, filePath)
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

	configMap := make(map[string]interface{}, 0)
	for _, r := range allParsedResources {
		addr := fmt.Sprintf("%s.%s", r.Type, r.Name)
		if addr == target {
			configMap = r.Attributes
			break
		}
	}

	jsonData, err := json.MarshalIndent(allParsedResources, "", "  ") // "" is prefix, "  " is indent string (2 spaces)
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %s", err)
	}

	outputFileName := "output.json"

	err = ioutil.WriteFile(outputFileName, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON data to file %s: %s", outputFileName, err)
	}

	return configMap, nil
}

// Compares HCL and finds all of the keys in map1 are in map2
func compareHCLFields(map1, map2 map[string]interface{}, path string, excludedFields map[string]bool) []string {
	var missingKeys []string
	for key, value1 := range map1 {
		if value1 == nil {
			continue
		}

		currentPath := path + "." + key
		if path == "" {
			currentPath = key
		}

		if excludedFields[currentPath] {
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
			missingKeys = append(missingKeys, compareHCLFields(v1, v2, currentPath, excludedFields)...)
		case []interface{}:
			v2, _ := value2.([]interface{})

			for i := 0; i < len(v1); i++ {
				nestedMap1, ok1 := v1[i].(map[string]interface{})
				nestedMap2, ok2 := v2[i].(map[string]interface{})
				if ok1 && ok2 {
					keys := compareHCLFields(nestedMap1, nestedMap2, fmt.Sprintf("%s[%d]", currentPath, i), excludedFields)
					missingKeys = append(missingKeys, keys...)
				}
			}
		default:
		}
	}

	return missingKeys
}

// Converts a tfplan to CAI asset, and then converts the CAI asset into HCL
func getRoundtripConfig(t *testing.T, fileName, tfDir string, ancestryCache map[string]string) ([]byte, error) {
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

	jsonData, err := json.MarshalIndent(roundtripAssets, "", "  ") // "" is prefix, "  " is indent string (2 spaces)
	if err != nil {
		log.Fatalf("Error marshalling to JSON: %s", err)
	}

	roundtripAssetFile := fmt.Sprintf("%s_roundtrip.json", fileName)
	err = ioutil.WriteFile(roundtripAssetFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON data to file %s: %s", roundtripAssetFile, err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	data, err := cai2hcl.Convert(roundtripAssets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
