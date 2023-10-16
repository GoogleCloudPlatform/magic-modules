package testing

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v2/caiasset"
	"go.uber.org/zap"

	"github.com/google/go-cmp/cmp"
)

type _TestCase struct {
	name         string
	sourceFolder string
}

func AssertTestFiles(t *testing.T, converterNames map[string]string, converterMap map[string]common.Converter, folder string, fileNames []string) {
	cases := []_TestCase{}

	for _, name := range fileNames {
		cases = append(cases, _TestCase{name: name, sourceFolder: folder})
	}

	for i := range cases {
		c := cases[i]

		t.Run(c.name, func(t *testing.T) {
			t.Parallel()

			err := assertTestData(c, converterNames, converterMap)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func assertTestData(testCase _TestCase, converterNames map[string]string, converterMap map[string]common.Converter) (err error) {
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

	got, err := common.Convert(assets, converterNames, converterMap)
	if err != nil {
		return err
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		logger.Debug(fmt.Sprintf("Expected %s to be:\n%s\nBut was:\n%s", expectedTfFilePath, string(want), string(got)))

		return fmt.Errorf("cmp.Diff() got diff (-want +got): %s", diff)
	}

	return nil
}
