package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

type ResourceMetadata struct {
	CaiAssetNames   []string            `json:"cai_asset_names"`
	ResourceType    string              `json:"resource_type"`
	ResourceAddress string              `json:"resource_address"`
	Service         string              `json:"service"`
	Cai             map[string]*CaiData `json:"cai_data,omitempty"` // Holds the fetched CAI assets data
}

// CaiData holds the fetched CAI asset and related error information.
type CaiData struct {
	CaiAsset caiasset.Asset `json:"cai_asset,omitempty"`
}

type TgcMetadataPayload struct {
	TestName         string                       `json:"test_name"`
	RawConfig        string                       `json:"raw_config"`
	ResourceMetadata map[string]*ResourceMetadata `json:"resource_metadata"`
	PrimaryResource  string                       `json:"primary_resource"`
}

type ResourceTestData struct {
	ParsedRawConfig  map[string]struct{} `json:"parsed_raw_config"`
	ResourceMetadata `json:"resource_metadata"`
}

type Resource struct {
	Type       string              `json:"type"`
	Name       string              `json:"name"`
	Attributes map[string]struct{} `json:"attributes"`
}

var (
	TestsMetadata = make(map[string]TgcMetadataPayload)
	setupDone     = false
)

func ReadTestsDataFromGcs() (map[string]TgcMetadataPayload, error) {
	if !setupDone {
		bucketName := "cai_assets_metadata"
		currentDate := time.Now()

		for len(TestsMetadata) == 0 {
			objectName := fmt.Sprintf("nightly_tests/%s/nightly_tests_meta.json", currentDate.Format("2006-01-02"))
			log.Printf("Read object  %s from the bucket %s", objectName, bucketName)

			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("storage.NewClient: %v", err)
			}
			defer client.Close()

			currentDate = currentDate.AddDate(0, 0, -1)

			rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
			if err != nil {
				if err == storage.ErrObjectNotExist {
					log.Printf("Object '%s' in bucket '%s' does NOT exist.\n", objectName, bucketName)
					continue
				} else {
					return nil, fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
				}
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, fmt.Errorf("io.ReadAll: %v", err)
			}

			err = json.Unmarshal(data, &TestsMetadata)
			if err != nil {
				return nil, fmt.Errorf("json.Unmarshal: %v", err)
			}
		}

		if os.Getenv("WRITE_FILES") != "" {
			writeJSONFile("../../tests_metadata.json", TestsMetadata)
		}
		setupDone = true
	}
	return TestsMetadata, nil
}

func prepareTestData(testName string) (map[string]ResourceTestData, string, error) {
	var err error
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	TestsMetadata, err = ReadTestsDataFromGcs()
	if err != nil {
		return nil, "", err
	}

	testMetadata := TestsMetadata[testName]
	resourceMetadata := testMetadata.ResourceMetadata
	if len(resourceMetadata) == 0 {
		log.Printf("Data of test is unavailable: %s", testName)
		return nil, "", nil
	}

	rawTfFile := fmt.Sprintf("%s.tf", testName)
	err = os.WriteFile(rawTfFile, []byte(testMetadata.RawConfig), 0644)
	if err != nil {
		return nil, "", fmt.Errorf("error writing to file %s: %#v", rawTfFile, err)
	}
	if os.Getenv("WRITE_FILES") == "" {
		defer os.Remove(rawTfFile)
	}

	rawResourceConfigs, err := parseResourceConfigs(rawTfFile)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing resource configs: %#v", err)
	}

	if len(rawResourceConfigs) == 0 {
		return nil, "", fmt.Errorf("Test %s fails: raw config is unavailable", testName)
	}

	rawConfigMap := convertToConfigMap(rawResourceConfigs)

	resourceTestData := make(map[string]ResourceTestData, 0)
	for address, metadata := range resourceMetadata {
		resourceTestData[address] = ResourceTestData{
			ParsedRawConfig:  rawConfigMap[address],
			ResourceMetadata: *metadata,
		}
	}

	return resourceTestData, testMetadata.PrimaryResource, nil
}

// Parses a Terraform configuation file written with HCL
func parseResourceConfigs(filePath string) ([]Resource, error) {
	src, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %s", filePath, err)
	}

	topLevel, err := parseHCLBytes(src, filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hcl bytes: %s", err)
	}

	var allParsedResources []Resource
	for addr, attrs := range topLevel {
		addrParts := strings.Split(addr, ".")
		if len(addrParts) != 2 {
			return nil, fmt.Errorf("invalid resource address %s", addr)
		}
		allParsedResources = append(allParsedResources, Resource{
			Type:       addrParts[0],
			Name:       addrParts[1],
			Attributes: attrs,
		})
	}
	return allParsedResources, nil
}

// Converts the slice to map with resource address as the key
func convertToConfigMap(resources []Resource) map[string]map[string]struct{} {
	configMap := make(map[string]map[string]struct{}, 0)

	for _, r := range resources {
		addr := fmt.Sprintf("%s.%s", r.Type, r.Name)
		configMap[addr] = r.Attributes
	}

	return configMap
}

// Converts the slice of assets to map with the asset name as the key
func convertToAssetMap(assets []caiasset.Asset) map[string]caiasset.Asset {
	assetMap := make(map[string]caiasset.Asset)

	for _, asset := range assets {
		assetMap[asset.Type] = asset
	}
	return assetMap
}
