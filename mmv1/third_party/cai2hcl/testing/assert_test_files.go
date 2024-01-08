package testing

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"go.uber.org/zap"

	"github.com/google/go-cmp/cmp"
)

type _TestCase struct {
	name         string
	sourceFolder string
}

func AssertTestFiles(t *testing.T, folder string, fileNames []string) {
	cases := []_TestCase{}

	for _, name := range fileNames {
		cases = append(cases, _TestCase{name: name, sourceFolder: folder})
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

func assertTestData(testCase _TestCase) (err error) {
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
	if err != nil {
		return err
	}

	got, err := cai2hcl.Convert(assets, &cai2hcl.Options{
		ErrorLogger: logger,
	})
	if err != nil {
		return err
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		logger.Debug(fmt.Sprintf("Expected %s to be:\n%s\nBut was:\n%s", expectedTfFilePath, string(want), string(got)))

		return fmt.Errorf("cmp.Diff() got diff (-want +got): %s", diff)
	}

	return nil
}
