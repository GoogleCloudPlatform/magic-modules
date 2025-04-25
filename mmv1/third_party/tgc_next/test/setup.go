package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

var TestConfig = make(map[string]TestMetadata)
var setupDone = false

func GlobalSetup() error {
	if !setupDone {
		bucketName := "cai_assets"
		objectName := fmt.Sprintf("nightly_tests/%s/nightly_tests_meta.json", "2025-04-21")

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

		setupDone = true
	}
	return nil
}

func yesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}
