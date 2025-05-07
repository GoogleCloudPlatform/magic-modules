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
	Resource   string
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

		bucketName := "cai_assets"
		currentDate := time.Now()

		for len(TestConfig) == 0 {
			previousDate := currentDate.AddDate(0, 0, -1)
			objectName := fmt.Sprintf("nightly_tests/%s/nightly_tests_meta.json", previousDate.Format("2006-01-02"))
			log.Printf("Read object  %s from the bucket %s", objectName, bucketName)

			ctx := context.Background()
			client, err := storage.NewClient(ctx)
			if err != nil {
				return fmt.Errorf("storage.NewClient: %v", err)
			}
			defer client.Close()

			rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
			if err != nil {
				return fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
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

			currentDate = previousDate

			// generateTests(TestConfig, "google_compute_instance", "compute.googleapis.com/Instance")

		}
		setupDone = true
		cacheMutex.Unlock()
	}
	return nil
}
