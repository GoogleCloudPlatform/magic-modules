# `google_artifact_registry_file` Data Source

**Issue:** hashicorp/terraform-provider-google#21449
**Branch:** `data-source-gar`
**Date:** 2026-04-13

## Purpose

Add a Terraform data source that downloads a single file from a Google Artifact Registry repository to a local path, exposing its metadata and content hashes as computed attributes. Analogous to the jFrog `artifactory_file` data source.

Applies to file-based AR formats: Generic, Maven, npm, Python, Apt, Yum, Go. Not applicable to Docker/OCI images (those are fetched via OCI clients, not as single files, and already have `google_artifact_registry_docker_image`).

## Schema

### Inputs

| Field | Type | Required | Description |
|---|---|---|---|
| `project` | string | optional | GCP project ID; falls back to the provider's default project |
| `location` | string | required | Repository location (e.g. `us-central1`) |
| `repository_id` | string | required | AR repository ID |
| `file_id` | string | required | File ID within the repository. Slashes and other reserved characters are URL-encoded by the provider before the API call |
| `output_path` | string | required | Local filesystem path where downloaded bytes are written. Parent directories are created if missing |

### Computed outputs

| Field | Type | Source |
|---|---|---|
| `name` | string | Full resource name: `projects/{project}/locations/{location}/repositories/{repository_id}/files/{file_id}` |
| `size_bytes` | int | `files.get` response field `sizeBytes` |
| `hashes` | `map(string)` | `files.get` response field `hashes` — hash type (e.g. `SHA256`, `MD5`) to base64 value |
| `create_time` | string | `files.get` field `createTime` |
| `update_time` | string | `files.get` field `updateTime` |
| `output_sha256` | string | Hex SHA-256 of the downloaded bytes, computed locally |
| `output_base64sha256` | string | Base64 SHA-256 of the downloaded bytes, computed locally |

The data source ID is set to the full resource name.

## Read flow

1. Resolve `project` via `tpgresource.GetProject`.
2. URL-encode `file_id` (using `url.PathEscape`) and build the metadata URL:
   `{{ArtifactRegistryBasePath}}projects/{project}/locations/{location}/repositories/{repository_id}/files/{file_id_encoded}`
3. Call `GET` via `transport_tpg.SendRequest` to retrieve metadata; populate `name`, `size_bytes`, `hashes`, `create_time`, `update_time`.
4. Build the download URL by appending `:download?alt=media` to the file resource URL.
5. Issue the download via a raw `http.Request` using `config.Client` (authenticated) with the provider user-agent, streaming the response body.
6. Write the body to `output_path`, creating parent directories as needed (`os.MkdirAll` on `filepath.Dir(output_path)`).
7. While writing, tee through `sha256.New()` to compute `output_sha256` (hex) and `output_base64sha256` (base64 std encoding).
8. `d.SetId(name)`.

## File layout

- `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file.go` — implementation
- `mmv1/third_party/terraform/services/artifactregistry/data_source_artifact_registry_file_test.go` — acceptance test
- `mmv1/third_party/terraform/website/docs/d/artifact_registry_file.html.markdown` — documentation

The data source self-registers in `init()` via `registry.Schema{...}.Register()`, matching the pattern of the existing AR data sources.

## Raw-bytes HTTP helper

The existing `transport_tpg.SendRequest*` helpers parse the response as JSON, which is unsuitable for binary media downloads. The implementation will construct an `http.Request` directly:

- Client: `config.Client` (already carries OAuth credentials)
- Headers: `User-Agent` from `tpgresource.GenerateUserAgentString`
- Method: `GET`
- Non-2xx responses surface as errors using the same formatting as `transport_tpg.SendRequest` where practical.

If a suitable streaming helper already exists in `transport_tpg`, the implementation will use it instead; otherwise the direct `http.Request` approach above is used. This will be confirmed during implementation.

## Testing

Acceptance test (`TestAccDataSourceArtifactRegistryFile_basic`):

1. Create a `google_artifact_registry_repository` with `format = "GENERIC"` in a random location.
2. Upload a small known file to the repository in a test pre-step (via the AR generic upload API, using the same auth as the test framework, or via `gcloud` invoked in the test).
3. Use `google_artifact_registry_file` with `output_path = <t.TempDir()>/<file>`.
4. Assert:
   - `output_sha256` equals the known hash of the uploaded content.
   - `size_bytes` equals the expected size.
   - The file exists on disk at `output_path` with the expected content.
   - `hashes["SHA256"]` is non-empty.

The test uses the project/credentials supplied by `TF_ACC` environment variables, same as other AR acceptance tests.

## Out of scope

- Docker/OCI image download (already covered by `google_artifact_registry_docker_image`).
- Bulk download or listing multiple files (could be a follow-up `google_artifact_registry_files`).
- Signed URL generation without downloading.
- Caching or conditional download based on hash.

## Open risks

- Large files: bytes are streamed through a hash, not held fully in memory, but Terraform state will include the file path and hashes. Users downloading multi-GB artifacts should be aware of plan/apply time implications.
- Authentication edge cases: any OAuth scopes required for `files.download` beyond what the provider already uses will be flagged and added if needed.
