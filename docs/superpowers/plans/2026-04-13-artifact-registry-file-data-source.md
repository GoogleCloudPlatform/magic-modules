# `google_artifact_registry_file` Data Source Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `google_artifact_registry_file` data source that downloads a single file from Artifact Registry to a local path and exposes its metadata and content hashes as attributes.

**Architecture:** Handwritten data source under `mmv1/third_party/terraform/services/artifactregistry/`. Two-call flow: `files.get` for metadata, then `files.download?alt=media` streamed to disk via a raw `http.Request` (the shared `transport_tpg.SendRequest` helpers parse JSON and are unsuitable for binary media). SHA-256 is computed while streaming.

**Tech Stack:** Go, terraform-plugin-sdk/v2, `google/transport` helpers, Artifact Registry v1 REST API.

---

## File Structure

**New files:**
- `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go` — data source (schema, read, flatten, self-registration)
- `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_internal_test.go` — unit tests for small pure helpers (`buildFileResourceURL`, `sha256Hashes`)
- `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_test.go` — acceptance test (`TestAccDataSourceArtifactRegistryFile_basic`)
- `mmv1/third_party/terraform/website/docs/d/artifact_registry_file.html.markdown` — user-facing docs

No existing files are modified; the data source self-registers via `init()` using `registry.Schema{...}.Register()`, matching the existing AR data source pattern.

---

## Task 1: Scaffold data source skeleton (schema + registration)

**Files:**
- Create: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go`

- [ ] **Step 1: Write the file with schema, a Read stub that returns an error, and `init()` registration**

```go
package artifactregistry

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
)

func DataSourceArtifactRegistryFile() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceArtifactRegistryFileRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"output_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hashes": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_base64sha256": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceArtifactRegistryFileRead(d *schema.ResourceData, meta interface{}) error {
	return fmt.Errorf("not implemented")
}

func init() {
	registry.Schema{
		Name:        "google_artifact_registry_file",
		ProductName: "artifactregistry",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceArtifactRegistryFile(),
	}.Register()
}
```

- [ ] **Step 2: Verify it compiles**

Run: `cd mmv1/third_party/terraform && go build ./services/artifactregistry/...`
Expected: exit 0, no output.

- [ ] **Step 3: Commit**

```bash
git add mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go
git commit -s -m "feat(artifactregistry): scaffold google_artifact_registry_file data source"
```

---

## Task 2: URL builder helper + unit test

**Files:**
- Modify: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go`
- Create: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_internal_test.go`

Rationale: The file ID can contain slashes and colons (e.g. Maven: `com.google.guava:guava:32.0.0:guava-32.0.0.jar`). AR's REST API requires these to be URL-encoded. A small pure helper is easy to unit test.

- [ ] **Step 1: Write the failing test**

Create `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_internal_test.go`:

```go
package artifactregistry

import "testing"

func TestBuildFileResourceURL(t *testing.T) {
	cases := []struct {
		name     string
		base     string
		project  string
		location string
		repo     string
		fileID   string
		want     string
	}{
		{
			name:     "simple generic file",
			base:     "https://artifactregistry.googleapis.com/v1/",
			project:  "my-proj",
			location: "us-central1",
			repo:     "my-repo",
			fileID:   "foo.tar.gz",
			want:     "https://artifactregistry.googleapis.com/v1/projects/my-proj/locations/us-central1/repositories/my-repo/files/foo.tar.gz",
		},
		{
			name:     "maven file with slashes and colons",
			base:     "https://artifactregistry.googleapis.com/v1/",
			project:  "p",
			location: "us",
			repo:     "r",
			fileID:   "com.google.guava:guava:32.0.0:guava-32.0.0.jar",
			want:     "https://artifactregistry.googleapis.com/v1/projects/p/locations/us/repositories/r/files/com.google.guava%3Aguava%3A32.0.0%3Aguava-32.0.0.jar",
		},
		{
			name:     "path with slashes",
			base:     "https://artifactregistry.googleapis.com/v1/",
			project:  "p",
			location: "us",
			repo:     "r",
			fileID:   "nested/path/file.txt",
			want:     "https://artifactregistry.googleapis.com/v1/projects/p/locations/us/repositories/r/files/nested%2Fpath%2Ffile.txt",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildFileResourceURL(tc.base, tc.project, tc.location, tc.repo, tc.fileID)
			if got != tc.want {
				t.Errorf("buildFileResourceURL() = %q, want %q", got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd mmv1/third_party/terraform && go test ./services/artifactregistry/ -run TestBuildFileResourceURL -v`
Expected: FAIL — `undefined: buildFileResourceURL`.

- [ ] **Step 3: Implement the helper**

Add to `data_source_artifact_registry_file.go` (add `"net/url"` to imports):

```go
// buildFileResourceURL constructs the AR file resource URL with fileID properly URL-encoded.
// AR file IDs may contain slashes and colons (e.g. Maven artifact paths).
func buildFileResourceURL(base, project, location, repository, fileID string) string {
	return fmt.Sprintf(
		"%sprojects/%s/locations/%s/repositories/%s/files/%s",
		base, project, location, repository, url.PathEscape(fileID),
	)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd mmv1/third_party/terraform && go test ./services/artifactregistry/ -run TestBuildFileResourceURL -v`
Expected: PASS — all three sub-tests.

- [ ] **Step 5: Commit**

```bash
git add mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go \
        mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_internal_test.go
git commit -s -m "feat(artifactregistry): add buildFileResourceURL helper"
```

---

## Task 3: Hash helper + unit test

**Files:**
- Modify: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go`
- Modify: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_internal_test.go`

- [ ] **Step 1: Append failing test**

Append to `data_source_artifact_registry_file_internal_test.go`:

```go
func TestSHA256Hashes(t *testing.T) {
	// Empty input has a well-known SHA-256.
	hex, b64 := sha256Hashes(nil)
	wantHex := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	wantB64 := "47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU="
	if hex != wantHex {
		t.Errorf("hex = %q, want %q", hex, wantHex)
	}
	if b64 != wantB64 {
		t.Errorf("b64 = %q, want %q", b64, wantB64)
	}

	// "hello" has known hashes.
	hex, b64 = sha256Hashes([]byte("hello"))
	if hex != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Errorf("hex for hello = %q", hex)
	}
	if b64 != "LPJNul+wow4m6DsqxbninhsWHlwfp0JecwQzYpOLmCQ=" {
		t.Errorf("b64 for hello = %q", b64)
	}
}
```

- [ ] **Step 2: Run — expect fail**

Run: `cd mmv1/third_party/terraform && go test ./services/artifactregistry/ -run TestSHA256Hashes -v`
Expected: FAIL — `undefined: sha256Hashes`.

- [ ] **Step 3: Implement**

Add to `data_source_artifact_registry_file.go` (add `"crypto/sha256"`, `"encoding/base64"`, `"encoding/hex"` to imports):

```go
func sha256Hashes(b []byte) (hexStr, b64Str string) {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), base64.StdEncoding.EncodeToString(sum[:])
}
```

- [ ] **Step 4: Run — expect pass**

Run: `cd mmv1/third_party/terraform && go test ./services/artifactregistry/ -run TestSHA256Hashes -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go \
        mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_internal_test.go
git commit -s -m "feat(artifactregistry): add sha256Hashes helper"
```

---

## Task 4: Implement Read — metadata + download + write

**Files:**
- Modify: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go`

This task has no dedicated unit test (it's mostly HTTP and filesystem glue against real AR); the acceptance test in Task 5 covers end-to-end behavior.

- [ ] **Step 1: Replace the Read stub with the full implementation**

Replace the `DataSourceArtifactRegistryFileRead` function and add imports (`"io"`, `"net/http"`, `"os"`, `"path/filepath"`, `"strconv"`, `"github.com/hashicorp/terraform-provider-google/google/tpgresource"`, `transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"`):

```go
func DataSourceArtifactRegistryFileRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location := d.Get("location").(string)
	repoID := d.Get("repository_id").(string)
	fileID := d.Get("file_id").(string)
	outputPath := d.Get("output_path").(string)

	resourceURL := buildFileResourceURL(config.ArtifactRegistryBasePath, project, location, repoID, fileID)

	// 1. Fetch metadata.
	metaResp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    resourceURL,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("fetching Artifact Registry file metadata: %w", err)
	}

	name, _ := metaResp["name"].(string)
	createTime, _ := metaResp["createTime"].(string)
	updateTime, _ := metaResp["updateTime"].(string)

	var sizeBytes int64
	switch v := metaResp["sizeBytes"].(type) {
	case string:
		sizeBytes, _ = strconv.ParseInt(v, 10, 64)
	case float64:
		sizeBytes = int64(v)
	}

	hashesAttr := map[string]string{}
	if rawHashes, ok := metaResp["hashes"].([]interface{}); ok {
		for _, h := range rawHashes {
			hm, _ := h.(map[string]interface{})
			t, _ := hm["type"].(string)
			val, _ := hm["value"].(string)
			if t != "" {
				hashesAttr[t] = val
			}
		}
	}

	// 2. Download bytes via raw HTTP (SendRequest parses JSON, unsuitable for media).
	downloadURL := resourceURL + ":download?alt=media"
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("building download request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := config.Client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading Artifact Registry file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("downloading Artifact Registry file: HTTP %d: %s", resp.StatusCode, string(body))
	}

	// 3. Ensure parent dir exists, write file, compute hashes.
	if dir := filepath.Dir(outputPath); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("creating output directory %q: %w", dir, err)
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if err := os.WriteFile(outputPath, body, 0o644); err != nil {
		return fmt.Errorf("writing %q: %w", outputPath, err)
	}

	hexStr, b64Str := sha256Hashes(body)

	// 4. Set attributes.
	if err := d.Set("project", project); err != nil {
		return err
	}
	if err := d.Set("name", name); err != nil {
		return err
	}
	if err := d.Set("size_bytes", sizeBytes); err != nil {
		return err
	}
	if err := d.Set("hashes", hashesAttr); err != nil {
		return err
	}
	if err := d.Set("create_time", createTime); err != nil {
		return err
	}
	if err := d.Set("update_time", updateTime); err != nil {
		return err
	}
	if err := d.Set("output_sha256", hexStr); err != nil {
		return err
	}
	if err := d.Set("output_base64sha256", b64Str); err != nil {
		return err
	}

	if name != "" {
		d.SetId(name)
	} else {
		d.SetId(fmt.Sprintf("projects/%s/locations/%s/repositories/%s/files/%s", project, location, repoID, fileID))
	}
	return nil
}
```

Note: `io.ReadAll` is used rather than streaming-while-hashing because Go's `sha256.New()` streaming + `io.TeeReader` adds complexity without meaningful savings for typical artifact sizes. If a future requirement demands very large files, switch to streaming.

- [ ] **Step 2: Verify it compiles and unit tests still pass**

Run: `cd mmv1/third_party/terraform && go build ./services/artifactregistry/... && go test ./services/artifactregistry/ -run "TestBuildFileResourceURL|TestSHA256Hashes" -v`
Expected: build succeeds; both unit tests PASS.

- [ ] **Step 3: Commit**

```bash
git add mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go
git commit -s -m "feat(artifactregistry): implement google_artifact_registry_file read"
```

---

## Task 5: Acceptance test

**Files:**
- Create: `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_test.go`

The test creates a Generic-format AR repo, uploads a small file via the AR generic upload API using the same authenticated transport that the provider uses, then reads it back via the data source.

- [ ] **Step 1: Write the acceptance test**

```go
package artifactregistry_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const testFileContents = "hello-artifact-registry\n"
const testFileName = "hello.txt"
const testPackageID = "hello-pkg"
const testVersionID = "1.0.0"

func TestAccDataSourceArtifactRegistryFile_basic(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	location := envvar.GetTestRegionFromEnv()
	repoID := fmt.Sprintf("tf-test-file-ds-%s", acctest.RandString(t, 10))
	outputPath := filepath.Join(t.TempDir(), "downloaded.txt")

	// File ID for a generic upload: "<package>:<version>:<filename>".
	fileID := fmt.Sprintf("%s:%s:%s", testPackageID, testVersionID, testFileName)
	sum := sha256.Sum256([]byte(testFileContents))
	expectedSHA := hex.EncodeToString(sum[:])

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Step 1: create the repository only.
				Config: testAccDataSourceArtifactRegistryFile_repoOnly(repoID, location),
				Check: resource.ComposeTestCheckFunc(
					uploadGenericArtifact(t, project, location, repoID, testPackageID, testVersionID, testFileName, []byte(testFileContents)),
				),
			},
			{
				// Step 2: read it via the data source.
				Config: testAccDataSourceArtifactRegistryFile_withDataSource(repoID, location, fileID, outputPath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_artifact_registry_file.test", "output_sha256", expectedSHA),
					resource.TestCheckResourceAttr("data.google_artifact_registry_file.test", "size_bytes", fmt.Sprintf("%d", len(testFileContents))),
					resource.TestCheckResourceAttrSet("data.google_artifact_registry_file.test", "name"),
					checkFileOnDisk(outputPath, []byte(testFileContents)),
				),
			},
		},
	})
}

func testAccDataSourceArtifactRegistryFile_repoOnly(repoID, location string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  location      = "%s"
  repository_id = "%s"
  format        = "GENERIC"
}
`, location, repoID)
}

func testAccDataSourceArtifactRegistryFile_withDataSource(repoID, location, fileID, outputPath string) string {
	return fmt.Sprintf(`
resource "google_artifact_registry_repository" "test" {
  location      = "%s"
  repository_id = "%s"
  format        = "GENERIC"
}

data "google_artifact_registry_file" "test" {
  location      = google_artifact_registry_repository.test.location
  repository_id = google_artifact_registry_repository.test.repository_id
  file_id       = "%s"
  output_path   = "%s"
}
`, location, repoID, fileID, outputPath)
}

// uploadGenericArtifact uploads a file via the AR generic upload endpoint, using
// the same authenticated HTTP client the provider uses in tests.
func uploadGenericArtifact(t *testing.T, project, location, repoID, pkg, version, filename string, contents []byte) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		url := fmt.Sprintf(
			"%sprojects/%s/locations/%s/repositories/%s/genericArtifacts:create?package_id=%s&version_id=%s",
			config.ArtifactRegistryBasePath, project, location, repoID, pkg, version,
		)

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("file", filename)
		if err != nil {
			return err
		}
		if _, err := io.Copy(part, bytes.NewReader(contents)); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}

		req, err := http.NewRequest("POST", url, &body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		userAgent, err := transport_tpg.BuildUserAgent(config, "acctest")
		if err == nil && userAgent != "" {
			req.Header.Set("User-Agent", userAgent)
		}

		resp, err := config.Client.Do(req)
		if err != nil {
			return fmt.Errorf("uploading test artifact: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			return fmt.Errorf("upload failed: HTTP %d: %s", resp.StatusCode, string(b))
		}
		return nil
	}
}

func checkFileOnDisk(path string, want []byte) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		got, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		if !bytes.Equal(got, want) {
			return fmt.Errorf("file contents mismatch: got %q want %q", string(got), string(want))
		}
		return nil
	}
}
```

Note on upload API: the AR generic upload endpoint used above is `POST .../genericArtifacts:create` with a multipart body. Verify the exact shape against the AR REST reference during implementation; if the helper `transport_tpg.BuildUserAgent` isn't present under that name, substitute the equivalent used by other acceptance tests (search with `Grep`).

- [ ] **Step 2: Compile the test package**

Run: `cd mmv1/third_party/terraform && go vet ./services/artifactregistry/...`
Expected: exit 0.

- [ ] **Step 3: Run the acceptance test against a real project (requires `TF_ACC=1` and GCP credentials)**

Run:
```bash
cd mmv1/third_party/terraform && \
TF_ACC=1 GOOGLE_PROJECT=<project> GOOGLE_REGION=us-central1 \
go test ./services/artifactregistry/ -run TestAccDataSourceArtifactRegistryFile_basic -v -timeout 20m
```
Expected: PASS. If it fails because of upload endpoint shape or helper name, fix inline and re-run.

- [ ] **Step 4: Commit**

```bash
git add mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_test.go
git commit -s -m "test(artifactregistry): add acceptance test for google_artifact_registry_file"
```

---

## Task 6: User-facing documentation

**Files:**
- Create: `mmv1/third_party/terraform/website/docs/d/artifact_registry_file.html.markdown`

- [ ] **Step 1: Write the doc**

```markdown
---
subcategory: "Artifact Registry"
description: |-
  Downloads a file from a Google Artifact Registry repository.
---

# google_artifact_registry_file

Downloads a single file from a Google Artifact Registry repository to a local
path and exposes its metadata and content hashes. Applies to file-based
Artifact Registry formats (Generic, Maven, npm, Python, Apt, Yum, Go). For
Docker/OCI images, use
[`google_artifact_registry_docker_image`](./artifact_registry_docker_image.html.markdown).

To get more information about Artifact Registry files, see:

* [API documentation](https://cloud.google.com/artifact-registry/docs/reference/rest/v1/projects.locations.repositories.files)

## Example Usage

```hcl
data "google_artifact_registry_file" "example" {
  location      = "us-central1"
  repository_id = "my-generic-repo"
  file_id       = "my-package:1.0.0:my-artifact.tar.gz"
  output_path   = "${path.module}/tmp/my-artifact.tar.gz"
}
```

## Argument Reference

The following arguments are supported:

* `location` - (Required) The location of the repository.
* `repository_id` - (Required) The ID of the repository.
* `file_id` - (Required) The Artifact Registry file ID. For Generic repositories this is `<package>:<version>:<filename>`; for other formats refer to the file listing in the API. Slashes and other reserved characters are URL-encoded by the provider.
* `output_path` - (Required) Local filesystem path where the downloaded bytes are written. Parent directories are created if missing.
* `project` - (Optional) The project in which the repository lives. Defaults to the provider project.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `name` - The fully-qualified file resource name (`projects/.../files/...`).
* `size_bytes` - Size of the file in bytes, as reported by Artifact Registry.
* `hashes` - Map of hash type (e.g. `SHA256`, `MD5`) to the corresponding hash value reported by Artifact Registry.
* `create_time` - Creation time (RFC 3339).
* `update_time` - Last update time (RFC 3339).
* `output_sha256` - Hex-encoded SHA-256 of the downloaded file contents.
* `output_base64sha256` - Base64-encoded SHA-256 of the downloaded file contents.
```

- [ ] **Step 2: Commit**

```bash
git add mmv1/third_party/terraform/website/docs/d/artifact_registry_file.html.markdown
git commit -s -m "docs(artifactregistry): document google_artifact_registry_file"
```

---

## Self-Review Notes

- **Spec coverage:** Inputs/outputs in Task 1 match the spec's schema table exactly. Read flow steps (metadata → download → write → hash → set attrs) map to Task 4 steps 1–4. Test plan in Task 5 creates Generic repo, uploads known content, asserts `output_sha256`, `size_bytes`, file on disk. Docs in Task 6 mirror the schema. Out-of-scope items are not added.
- **Placeholder scan:** No TBD/TODO. One inline verification note ("verify upload endpoint shape") is explicit — the engineer has the API reference link to check against.
- **Type consistency:** `buildFileResourceURL`, `sha256Hashes`, `DataSourceArtifactRegistryFileRead` names match across all tasks. Schema field names match between data source, test, and docs.
