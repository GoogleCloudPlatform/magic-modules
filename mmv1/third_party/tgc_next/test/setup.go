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
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/caiasset"
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
	MetadataByTestAndStep map[string]map[int]TgcMetadataPayload
	Date                  time.Time
}

// The metadata for each step in one test
type TgcMetadataPayload struct {
	TestName         string                       `json:"test_name"`
	StepNumber       int                          `json:"step_number"`
	RawConfig        string                       `json:"raw_config"`
	ResourceMetadata map[string]*ResourceMetadata `json:"resource_metadata"`
	PrimaryResource  string                       `json:"primary_resource"`
	CaiReadTime      time.Time                    `json:"cai_read_time"`
}

type ResourceTestData struct {
	ParsedRawConfig  map[string]any `json:"parsed_raw_config"`
	ResourceMetadata `json:"resource_metadata"`
}

type StepTestData struct {
	StepNumber       int
	PrimaryResource  string
	ResourceTestData map[string]ResourceTestData // key is resource address
}

type Resource struct {
	Type       string         `json:"type"`
	Name       string         `json:"name"`
	Attributes map[string]any `json:"attributes"`
}

const (
	ymdFormat   = "2006-01-02"
	maxAttempts = 5
)

var (
	TestsMetadata = make([]NightlyRun, maxAttempts)
	setupDone     = false
)

func ReadTestsDataFromGcs() ([]NightlyRun, error) {
	if !setupDone {
		bucketName := "cai_assets_metadata"
		currentDate := time.Now()
		ctx := context.Background()

		var client *storage.Client
		var bucket *storage.BucketHandle
		var err error

		var allErrs error
		retries := 0
		for i := 0; i < len(TestsMetadata); i++ {
			var metadata map[string]map[int]TgcMetadataPayload
			if os.Getenv("WRITE_FILES") != "" {
				filename := fmt.Sprintf("../../tests_metadata_%s.json", currentDate.Format(ymdFormat))
				_, err := os.Stat(filename)
				if !os.IsNotExist(err) {
					metadata = readTestsDataFromLocalFile(filename)
				}
			}
			if metadata == nil {
				if client == nil {
					client, err = storage.NewClient(ctx)
					if err != nil {
						return nil, fmt.Errorf("storage.NewClient: %v", err)
					}
					defer client.Close()
					bucket = client.Bucket(bucketName)
				}
				metadata, err = readTestsDataFromGCSForRun(ctx, currentDate, bucketName, bucket)
				if os.Getenv("WRITE_FILES") != "" {
					writeJSONFile(fmt.Sprintf("../../tests_metadata_%s.json", currentDate.Format(ymdFormat)), metadata)
				}

				if err != nil {
					if allErrs == nil {
						allErrs = fmt.Errorf("reading tests data from gcs: %v", err)
					} else {
						allErrs = fmt.Errorf("%v, %v", allErrs, err)
					}
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
					MetadataByTestAndStep: metadata,
					Date:                  currentDate,
				}
			}
			currentDate = currentDate.AddDate(0, 0, -1)
		}

		if allErrs != nil {
			return nil, allErrs
		}
		setupDone = true
	}
	return TestsMetadata, nil
}

func readTestsDataFromGCSForRun(ctx context.Context, currentDate time.Time, bucketName string, bucket *storage.BucketHandle) (map[string]map[int]TgcMetadataPayload, error) {
	metadata := make(map[string]map[int]TgcMetadataPayload)
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

func readTestsDataFromLocalFile(filename string) map[string]map[int]TgcMetadataPayload {
	metadata := make(map[string]map[int]TgcMetadataPayload, 0)
	log.Printf("Read the the local file %s", filename)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	if err != nil {
		return nil
	}

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return nil
	}

	return metadata
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
		testMetadata, ok := run.MetadataByTestAndStep[testName]
		if ok && len(testMetadata) > 0 {
			for stepNumber := range testMetadata {
				stepNumbers = append(stepNumbers, stepNumber)
			}
			break
		}
	}
	return stepNumbers, nil
}

func prepareTestData(testName string, stepNumber int, retries int) (*StepTestData, error) {
	var err error

	var testMetadata map[int]TgcMetadataPayload

	run := TestsMetadata[retries]
	testMetadata, ok := run.MetadataByTestAndStep[testName]
	if !ok {
		log.Printf("Data of test is unavailable: %s", testName)
		return nil, nil
	}

	log.Printf("Found metadata for %s from run on %s", testName, run.Date.Format(ymdFormat))

	if stepMetadata, ok := testMetadata[stepNumber]; ok {
		resourceMetadata := stepMetadata.ResourceMetadata

		rawTfFile := fmt.Sprintf("%s_step%d.tf", testName, stepNumber)
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
			return nil, fmt.Errorf("test %s fails: raw config is unavailable", testName)
		}

		if os.Getenv("WRITE_FILES") != "" {
			writeJSONFile(fmt.Sprintf("%s_attrs", testName), rawResourceConfigs)
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
			StepNumber:       stepNumber,
			PrimaryResource:  stepMetadata.PrimaryResource,
			ResourceTestData: resourceTestData,
		}, nil
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
func convertToConfigMap(resources []Resource) map[string]map[string]any {
	configMap := make(map[string]map[string]any, 0)

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
