package storage_test

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/provider"
	_ "github.com/hashicorp/terraform-provider-google/google/services/storage"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type crc32cMismatchRoundTripper struct {
	http.RoundTripper
	bucketName string
	objectName string
}

func (t *crc32cMismatchRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := t.RoundTripper.RoundTrip(r)
	if err != nil || response.StatusCode != http.StatusOK {
		return response, err
	}

	// Intercept object metadata GET requests and mutate the returned CRC32c checksum.
	// We distinguish metadata GET requests from payload download requests by ensuring
	// the `alt=media` query parameter is NOT present.
	if r.Method == http.MethodGet &&
		strings.Contains(r.URL.Path, fmt.Sprintf("/b/%s/o/%s", t.bucketName, t.objectName)) &&
		r.URL.Query().Get("alt") != "media" {

		responseBytes, readErr := io.ReadAll(response.Body)
		response.Body.Close()
		if readErr != nil {
			return nil, readErr
		}

		var responseMap map[string]interface{}
		if jsonErr := json.Unmarshal(responseBytes, &responseMap); jsonErr != nil {
			return nil, jsonErr
		}

		// Inject a mismatched base64-encoded CRC32c checksum.
		// Standard valid GCS CRC32c hashes are 4 bytes, base64-encoded.
		responseMap["crc32c"] = "MismatchedCrc32cForTesting=="

		newBytes, marshalErr := json.Marshal(responseMap)
		if marshalErr != nil {
			return nil, marshalErr
		}

		response.Body = io.NopCloser(bytes.NewReader(newBytes))
		response.ContentLength = int64(len(newBytes))
		response.Header.Set("Content-Length", strconv.Itoa(len(newBytes)))
	}

	return response, err
}

func TestAccDataSourceStorageBucketObjectContent_Crc32cMismatch(t *testing.T) {
	acctest.SkipIfVcr(t)

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"
	objectName := "butterfly01"

	providerInstance := provider.Provider()
	oldConfigureFunc := providerInstance.ConfigureContextFunc
	providerInstance.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		c, diagnostics := oldConfigureFunc(ctx, d)
		if diagnostics.HasError() {
			return c, diagnostics
		}
		config := c.(*transport_tpg.Config)
		config.Client.Transport = &crc32cMismatchRoundTripper{
			RoundTripper: config.Client.Transport,
			bucketName:   bucket,
			objectName:   objectName,
		}
		return c, diagnostics
	}

	providers := map[string]*schema.Provider{
		"google": providerInstance,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:  func() { acctest.AccTestPreCheck(t) },
		Providers: providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceStorageBucketObjectContent_Basic(content, bucket),
				ExpectError: regexp.MustCompile("CRC32C checksum mismatch for storage bucket object"),
			},
		},
	})
}

func TestAccDataSourceStorageBucketObjectContent_Basic(t *testing.T) {

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageBucketObjectContent_Basic(content, bucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content", content),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content_base64"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content_base64", base64.StdEncoding.EncodeToString([]byte(content))),
				),
			},
		},
	})
}

func TestAccDataSourceStorageBucketObjectContent_FileContentBase64(t *testing.T) {
	acctest.SkipIfVcr(t)

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	folderName := "tf-folder-" + acctest.RandString(t, 10)

	if err := os.Mkdir(folderName, 0777); err != nil {
		t.Errorf("error creating directory: %v", err)
	}

	data := []byte("data data data")
	testFile := getTmpTestFile(t, folderName, "tf-test")
	if err := ioutil.WriteFile(testFile.Name(), data, 0644); err != nil {
		t.Errorf("error writing file: %v", err)
	}
	defer os.Remove(testFile.Name()) // clean up

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"local": resource.ExternalProvider{
				VersionConstraint: "> 2.5.0",
			},
			"archive": resource.ExternalProvider{
				VersionConstraint: "> 2.5.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStorageBucketObjectContent_FileContentBase64(bucket, folderName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.this", "content_base64"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.this", "content_hexsha512"),
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.this", "content_base64sha512"),
					verifyValidZip(),
				),
			},
		},
	})
}

func verifyValidZip() func(*terraform.State) error {
	return func(s *terraform.State) error {
		var outputFilePath string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "local_file" {
				outputFilePath = rs.Primary.Attributes["filename"]
				break
			}
		}
		archive, err := zip.OpenReader(outputFilePath)
		if err != nil {
			return err
		}
		defer archive.Close()
		return nil
	}
}

func testAccDataSourceStorageBucketObjectContent_Basic(content, bucket string) string {
	return fmt.Sprintf(`
data "google_storage_bucket_object_content" "default" {
	bucket = google_storage_bucket.contenttest.name
	name   = google_storage_bucket_object.object.name      
}

resource "google_storage_bucket_object" "object" {
	name    = "butterfly01"
	content = "%s"
	bucket  = google_storage_bucket.contenttest.name
}

resource "google_storage_bucket" "contenttest" {
	name          = "%s"
	location      = "US"
	force_destroy = true
}`, content, bucket)
}

func testAccDataSourceStorageBucketObjectContent_FileContentBase64(bucket, folderName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "this" {
  name                        = "%s"
  location                    = "us-east4"
  uniform_bucket_level_access = true
}

data "archive_file" "this" {
  type       = "zip"
  source_dir = "${path.cwd}/%s"
  output_path = "${path.cwd}/archive.zip"
}

resource "google_storage_bucket_object" "this" {
  name   = "archive.zip"
  bucket = google_storage_bucket.this.name
  source = data.archive_file.this.output_path
}

data "google_storage_bucket_object_content" "this" {
  name   = google_storage_bucket_object.this.name
  bucket = google_storage_bucket.this.name
}

resource "local_file" "this" {
  content_base64 = (data.google_storage_bucket_object_content.this.content_base64)
  filename = "${path.cwd}/content.zip"
}`, bucket, folderName)
}

func TestAccDataSourceStorageBucketObjectContent_Issue15717(t *testing.T) {

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	config := fmt.Sprintf(`
%s

output "output" {
	value = replace(data.google_storage_bucket_object_content.default.content, "q", "Q")
}`, testAccDataSourceStorageBucketObjectContent_Basic(content, bucket))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.default", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.default", "content", content),
				),
			},
		},
	})
}

func TestAccDataSourceStorageBucketObjectContent_Issue15717BackwardCompatibility(t *testing.T) {

	bucket := "tf-bucket-object-content-" + acctest.RandString(t, 10)
	content := "qwertyuioasdfghjk1234567!!@#$*"

	config := fmt.Sprintf(`
%s

data "google_storage_bucket_object_content" "new" {
	bucket  = google_storage_bucket.contenttest.name
	content = "%s"
	name    = google_storage_bucket_object.object.name
}`, testAccDataSourceStorageBucketObjectContent_Basic(content, bucket), content)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_storage_bucket_object_content.new", "content"),
					resource.TestCheckResourceAttr("data.google_storage_bucket_object_content.new", "content", content),
				),
			},
		},
	})
}

func getTmpTestFile(t *testing.T, folderName, prefix string) *os.File {
	testFile, err := ioutil.TempFile(folderName, prefix)
	if err != nil {
		t.Fatalf("Cannot create temp file: %s", err)
	}
	return testFile
}
