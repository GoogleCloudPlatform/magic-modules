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

type NightlyRun struct {
	MetadataByTest map[string][]TgcMetadataPayload
	Date           time.Time
}

// The metadata for each step in one test
type TgcMetadataPayload struct {
	TestName         string                       `json:"test_name"`
	StepNumber       int                          `json:"step_number"`
	RawConfig        string                       `json:"raw_config"`
	ResourceMetadata map[string]*ResourceMetadata `json:"resource_metadata"`
	PrimaryResource  string                       `json:"primary_resource"`
}

type ResourceTestData struct {
	ParsedRawConfig  map[string]struct{} `json:"parsed_raw_config"`
	ResourceMetadata `json:"resource_metadata"`
}

type StepTestData struct {
	StepNumber       int
	PrimaryResource  string
	ResourceTestData map[string]ResourceTestData // key is resource address
}

type Resource struct {
	Type       string              `json:"type"`
	Name       string              `json:"name"`
	Attributes map[string]struct{} `json:"attributes"`
}

const (
	ymdFormat   = "2006-01-02"
	maxAttempts = 3
)

var (
	TestsMetadata = make([]NightlyRun, maxAttempts)
	setupDone     = false
)

func ReadTestsDataFromGcs() ([]NightlyRun, error) {
	if !setupDone {
		// bucketName := "cai_assets_metadata"
		bucketName := "cai_assets" // Use the bucket in testing project for tansition
		currentDate := time.Now()
		ctx := context.Background()
		client, err := storage.NewClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("storage.NewClient: %v", err)
		}
		defer client.Close()

		bucket := client.Bucket(bucketName)

		var allErrs error
		retries := 0
		for i := 0; i < len(TestsMetadata); i++ {
			metadata, err := readTestsDataFromGCSForRun(ctx, currentDate, bucketName, bucket)
			if err != nil {
				if allErrs == nil {
					allErrs = fmt.Errorf("reading tests data from gcs: %v", err)
				} else {
					allErrs = fmt.Errorf("%v, %v", allErrs, err)
				}
			}
			if metadata == nil {
				// Keep looking until we find a date with metadata.
				i--
				retries++
				if retries > maxAttempts {
					// Stop looking when we find maxAttempts dates with no metadata.
					return nil, fmt.Errorf("too many retries, %v", allErrs)
				}
			} else {
				TestsMetadata[i] = NightlyRun{
					MetadataByTest: metadata,
					Date:           currentDate,
				}
			}
			currentDate = currentDate.AddDate(0, 0, -1)
		}

		if allErrs != nil {
			return nil, allErrs
		}

		if os.Getenv("WRITE_FILES") != "" {
			writeJSONFile("../../tests_metadata.json", TestsMetadata)
		}
		setupDone = true
	}
	return TestsMetadata, nil
}

func readTestsDataFromGCSForRun(ctx context.Context, currentDate time.Time, bucketName string, bucket *storage.BucketHandle) (map[string][]TgcMetadataPayload, error) {
	metadata := make(map[string][]TgcMetadataPayload)
	objectName := fmt.Sprintf("nightly_tests/%s/nightly_tests_meta.json", currentDate.Format(ymdFormat))
	log.Printf("Read object  %s from the bucket %s", objectName, bucketName)

	rc, err := bucket.Object(objectName).NewReader(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			log.Printf("Object '%s' in bucket '%s' does NOT exist.\n", objectName, bucketName)
			return nil, nil
		} else {
			return nil, fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
		}
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %v", err)
	}

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return metadata, nil
}

func getStepNumbers(testName string) ([]int, error) {
	var err error
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	TestsMetadata, err = ReadTestsDataFromGcs()
	if err != nil {
		return nil, err
	}

	stepNumbers := make([]int, 0)
	for _, run := range TestsMetadata {
		testMetadata, ok := run.MetadataByTest[testName]
		if ok && len(testMetadata) > 0 {
			for _, stepMetadata := range testMetadata {
				stepNumbers = append(stepNumbers, stepMetadata.StepNumber)
			}
			break
		}
	}
	return stepNumbers, nil
}

func prepareTestData(testName string, targetStep int, retries int) (*StepTestData, error) {
	var err error

	var testMetadata []TgcMetadataPayload

	run := TestsMetadata[retries]
	testMetadata, ok := run.MetadataByTest[testName]
	if !ok {
		log.Printf("Data of test is unavailable: %s", testName)
		return nil, nil
	}

	log.Printf("Found metadata for %s from run on %s", testName, run.Date.Format(ymdFormat))

	for _, stepMetadata := range testMetadata {
		stepN := stepMetadata.StepNumber
		if stepN == targetStep {
			resourceMetadata := stepMetadata.ResourceMetadata

			rawTfFile := fmt.Sprintf("%s_step%d.tf", testName, stepN)
			err = os.WriteFile(rawTfFile, []byte(stepMetadata.RawConfig), 0644)
			if err != nil {
				return nil, fmt.Errorf("error writing to file %s: %#v", rawTfFile, err)
			}
			if os.Getenv("WRITE_FILES") == "" {
				defer os.Remove(rawTfFile)
			}

			rawResourceConfigs, err := parseResourceConfigs(rawTfFile)
			if err != nil {
				return nil, fmt.Errorf("error parsing resource configs: %#v", err)
			}

			if len(rawResourceConfigs) == 0 {
				return nil, fmt.Errorf("Test %s fails: raw config is unavailable", testName)
			}

			rawConfigMap := convertToConfigMap(rawResourceConfigs)

			resourceTestData := make(map[string]ResourceTestData, 0)
			for address, metadata := range resourceMetadata {
				resourceTestData[address] = ResourceTestData{
					ParsedRawConfig:  rawConfigMap[address],
					ResourceMetadata: *metadata,
				}
			}
			return &StepTestData{
				StepNumber:       stepN,
				PrimaryResource:  stepMetadata.PrimaryResource,
				ResourceTestData: resourceTestData,
			}, nil
		}
	}

	return nil, nil
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
