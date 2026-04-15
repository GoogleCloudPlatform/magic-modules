package cai2hcl

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/cmd/tgc/common"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func testHCLBlocks() []byte {
	testBlock := `
resource "google_compute_forwarding_rule" "test-2" {
  backend_service       = "projects/myproj/regions/us-central1/backendServices/test-bs-1"
  ip_address            = "projects/myproj/regions/us-central1/addresses/test-ip-1"
  ip_protocol           = "TCP"
  ip_version            = "IPV6"
  load_balancing_scheme = "EXTERNAL"
  name                  = "test-2"
  ports                 = ["80", "81"]
  region                = "us-central1"
}
	`
	return []byte(testBlock)
}

func mockConvertHCL(ctx context.Context, path string, errorLogger *zap.Logger) ([]byte, error) {
	return testHCLBlocks(), nil
}

func TestConvertRun(t *testing.T) {
	convertFunc = mockConvertHCL
	defer func() {
		convertFunc = origConvertFunc
	}()
	a := assert.New(t)
	verbosity := "debug"
	useStructuredLogging := true
	errorLogger, errorBuf := common.NewTestErrorLogger(verbosity, useStructuredLogging)
	outputLogger, outputBuf := common.NewTestOutputLogger()
	ro := &common.RootOptions{
		Verbosity:            verbosity,
		UseStructuredLogging: useStructuredLogging,
		ErrorLogger:          errorLogger,
		OutputLogger:         outputLogger,
	}
	o := convertOptions{
		rootOptions: ro,
	}

	path := "/path/to/cai_assets"
	err := o.run(path)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.Bytes()

	a.Equal(errorJSON, "")

	var output map[string]interface{}
	json.Unmarshal(outputJSON, &output)

	// On a successful run, we should see a list of google assets in the resource_body field
	a.Contains(output, "resource_body")

	var expectedHCLBlocks string
	expectedHCLJSON, _ := json.Marshal(string(testHCLBlocks()))
	json.Unmarshal(expectedHCLJSON, &expectedHCLBlocks)
	a.Equal(expectedHCLBlocks, output["resource_body"])
}

func TestConvertRunLegacy(t *testing.T) {
	convertFunc = mockConvertHCL
	defer func() {
		convertFunc = origConvertFunc
	}()
	a := assert.New(t)
	verbosity := "debug"
	useStructuredLogging := false
	errorLogger, errorBuf := common.NewTestErrorLogger(verbosity, useStructuredLogging)
	outputLogger, outputBuf := common.NewTestOutputLogger()
	ro := &common.RootOptions{
		Verbosity:            verbosity,
		UseStructuredLogging: useStructuredLogging,
		ErrorLogger:          errorLogger,
		OutputLogger:         outputLogger,
	}
	o := convertOptions{
		rootOptions: ro,
	}

	err := o.run("/path/to/plan")
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	// On a successful legacy run, we don't output anything via loggers.
	a.Equal(errorJSON, "")
	a.Equal(outputJSON, "")
}

func TestConvertRunOutputFile(t *testing.T) {
	convertFunc = mockConvertHCL
	defer func() {
		convertFunc = origConvertFunc
	}()
	a := assert.New(t)
	verbosity := "debug"
	useStructuredLogging := false
	errorLogger, errorBuf := common.NewTestErrorLogger(verbosity, useStructuredLogging)
	outputLogger, outputBuf := common.NewTestOutputLogger()
	ro := &common.RootOptions{
		Verbosity:            verbosity,
		UseStructuredLogging: useStructuredLogging,
		ErrorLogger:          errorLogger,
		OutputLogger:         outputLogger,
	}
	outputPath := path.Join(t.TempDir(), "converted.json")
	o := convertOptions{
		rootOptions: ro,
		outputPath:  outputPath,
	}

	path := "/path/to/plan"
	err := o.run(path)
	a.Nil(err)

	errorJSON := errorBuf.String()
	outputJSON := outputBuf.String()

	a.Equal(errorJSON, "")
	a.Equal(outputJSON, "")

	b, err := os.ReadFile(outputPath)
	if err != nil {
		a.Failf("Unable to read file %s: %s", outputPath, err)
	}

	expectedHCLBlocks := testHCLBlocks()
	a.Equal(expectedHCLBlocks, b)
}
