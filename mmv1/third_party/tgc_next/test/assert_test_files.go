package test

import (
	"context"
	"fmt"
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
	hcl "github.com/joselitofilho/hcl-parser-go/pkg/parser/hcl"
)

func AssertTestFile(t *testing.T, fileName, resource, assetType string, excludedFields []string) {
	GlobalSetup()

	// Create a temporary directory for running terraform.
	tfDir, err := os.MkdirTemp(tmpDir, "terraform")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tfDir)

	testMetadata := TestConfig[fileName]
	exportAssets := testMetadata.Assets
	targetAsset, found := findTargetAsset(exportAssets, assetType)
	if !found {
		log.Fatalf("error finding the target asset with type %s", assetType)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	exportConfigData, err := cai2hcl.Convert(exportAssets, &cai2hcl.Options{
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

	exportConfigMap, err := getConfig(exportTfFilePath, resource)
	if err != nil {
		log.Fatal(err)
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

	rawConfigMap, err := getConfig(rawTfFile, resource)
	if err != nil {
		log.Fatal(err)
	}
	if len(rawConfigMap) == 0 {
		log.Fatalf("raw config for test %s is unavailable", fileName)
	}

	excludedFieldMap := make(map[string]bool, 0)
	for _, f := range excludedFields {
		excludedFieldMap[f] = true
	}

	for address, rawConfig := range rawConfigMap {
		if exportConfig, ok := exportConfigMap[address]; !ok {
			log.Fatalf("%s - Missing resource after cai2hcl conversion: %s.", t.Name(), address)
		} else {
			missingKeys := compareHCLFields(rawConfig.(map[string]interface{}), exportConfig.(map[string]interface{}), "", excludedFieldMap)
			if len(missingKeys) > 0 {
				log.Fatalf("%s - Missing fields in resource %s after cai2hcl conversion:\n%s", t.Name(), address, missingKeys)
			}
		}
	}

	// Get the ancestry cache for tfplan2cai conversion
	ancestors := targetAsset.Ancestors
	ancestryCache := make(map[string]string, 0)
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
		ancestryCache[ancestors[0]] = path

		project := utils.ParseFieldValue(targetAsset.Name, "projects")
		projectKey := fmt.Sprintf("projects/%s", project)
		if strings.HasPrefix(ancestors[0], "projects") && ancestors[0] != projectKey {
			ancestryCache[projectKey] = path
		}
	}

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

		log.Fatalf("cmp.Diff() got diff (-want +got): %s", diff)
	}
}

func findTargetAsset(assets []caiasset.Asset, assetType string) (*caiasset.Asset, bool) {
	for _, asset := range assets {
		if asset.Type == assetType {
			return &asset, true
		}
	}
	return nil, false
}

func getConfig(filePath, target string) (map[string]interface{}, error) {
	files := []string{filePath}

	// Parse Terraform configurations
	config, err := hcl.Parse([]string{}, files)
	if err != nil {
		return nil, err
	}

	configMap := make(map[string]interface{}, 0)
	for _, r := range config.Resources {
		if r.Type != target {
			continue
		}
		addr := fmt.Sprintf("%s.%s", r.Type, r.Name)
		configMap[addr] = r.Attributes
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
			// if !ok {
			// 	fmt.Printf("Type mismatch for key: %s\n", currentPath)
			// 	continue
			// }
			missingKeys = append(missingKeys, compareHCLFields(v1, v2, currentPath, excludedFields)...)
		case []interface{}:
			v2, _ := value2.([]interface{})
			// if !ok {
			// 	fmt.Printf("Type mismatch for key: %s\n", currentPath)
			// 	continue
			// }
			// if len(v1) != len(v2) {
			// 	fmt.Printf("List length mismatch for key: %s\n", currentPath)
			// 	continue
			// }

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
