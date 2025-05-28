package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

type ResourceMetadata struct {
	CaiAssetName    string         `json:"cai_asset_name"`
	CaiAssetData    caiasset.Asset `json:"cai_asset_data"`
	ResourceType    string         `json:"resource_type"`
	ResourceAddress string         `json:"resource_address"`
	ImportMetadata  ImportMetadata `json:"import_metadata,omitempty"`
	Service         string         `json:"service"`
}

type ImportMetadata struct {
	Id            string   `json:"id,omitempty"`
	IgnoredFields []string `json:"ignored_fields,omitempty"`
}

type TgcMetadataPayload struct {
	TestName         string                       `json:"test_name"`
	RawConfig        string                       `json:"raw_config"`
	ResourceMetadata map[string]*ResourceMetadata `json:"resource_metadata"`
	PrimaryResource  string                       `json:"primary_resource"`
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

		// Uncomment this line to debug issues locally
		// writeJSONFile("../../tests_metadata.json", TestsMetadata)
		setupDone = true
	}
	return TestsMetadata, nil
}
