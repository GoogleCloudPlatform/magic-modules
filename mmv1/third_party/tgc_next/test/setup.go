package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/caiasset"
)

type TestMetadata struct {
	Test       string
	RawConfig  string
	Service    string
	Address    string
	AssetNames []string
	Assets     []caiasset.Asset
}

var (
	TestConfig = make(map[string]TestMetadata)
	setupDone  = false
	cacheMutex = sync.Mutex{}
)

func ReadTestsDataFromGcs() error {
	if !setupDone {
		cacheMutex.Lock()

		bucketName := "cai_assets_metadata"
		currentDate := time.Now()

		for len(TestConfig) == 0 {
			objectName := fmt.Sprintf("nightly_tests/%s/nightly_tests_meta.json", currentDate.Format("2006-01-02"))
			log.Printf("Read object  %s from the bucket %s", objectName, bucketName)

			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("storage.NewClient: %v", err)
			}
			defer client.Close()

			currentDate = currentDate.AddDate(0, 0, -1)

			rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
			if err != nil {
				if err == storage.ErrObjectNotExist {
					log.Printf("Object '%s' in bucket '%s' does NOT exist.\n", objectName, bucketName)
					continue
				} else {
					return fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
				}
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return fmt.Errorf("io.ReadAll: %v", err)
			}

			err = json.Unmarshal(data, &TestConfig)
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %v", err)
			}

			// generateTests(TestConfig, "google_compute_instance", "compute.googleapis.com/Instance")

		}
		setupDone = true
		cacheMutex.Unlock()
	}
	return nil
}
