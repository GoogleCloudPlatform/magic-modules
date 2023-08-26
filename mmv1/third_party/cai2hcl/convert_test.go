package cai2hcl

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
)

type TestCase struct {
	name         string
	sourceFolder string
}

// Files from "testdata/" to test json -> tf conversion.
var testDataFileNames = []string{
	"compute_instance_iam",
	"full_compute_instance",
	"project_create",
	"project_iam",
	"full_compute_forwarding_rule",
	"full_compute_backend_service",
	"full_compute_health_check",
}

func TestCai2HclConvert(t *testing.T) {
	cases := []TestCase{}

	for _, name := range testDataFileNames {
		cases = append(cases, TestCase{name: name, sourceFolder: "./testdata"})
	}

	for i := range cases {
		c := cases[i]

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			err := assertTestData(c)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func assertTestData(testCase TestCase) (err error) {
	fileName := testCase.name
	folder := testCase.sourceFolder

	assetFilePath := fmt.Sprintf("%s/%s.json", folder, fileName)
	expectedTfFilePath := fmt.Sprintf("%s/%s.tf", folder, fileName)
	assetPayload, err := os.ReadFile(assetFilePath)
	if err != nil {
		return fmt.Errorf("cannot open %s, got: %s", assetFilePath, err)
	}
	want, err := os.ReadFile(expectedTfFilePath)
	if err != nil {
		return fmt.Errorf("cannot open %s, got: %s", expectedTfFilePath, err)
	}

	var assets []*caiasset.Asset
	if err := json.Unmarshal(assetPayload, &assets); err != nil {
		return fmt.Errorf("cannot unmarshal: %s", err)
	}

	logger, err := zap.NewDevelopment()

	got, err := Convert(assets, &ConvertOptions{
		ErrorLogger: logger,
	})
	if err != nil {
		return err
	}
	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		return fmt.Errorf("cmp.Diff() got diff (-want +got): %s", diff)
	}

	return nil
}
