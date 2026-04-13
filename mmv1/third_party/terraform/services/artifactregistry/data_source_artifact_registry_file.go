package artifactregistry

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
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

// buildFileResourceURL constructs the AR file resource URL with fileID properly URL-encoded.
// AR file IDs may contain slashes and colons (e.g. Maven artifact paths).
// url.PathEscape encodes slashes but leaves colons unescaped (valid per RFC 3986 path segments).
// AR API requires colons to be percent-encoded as well, so we encode them explicitly.
func buildFileResourceURL(base, project, location, repository, fileID string) string {
	encoded := strings.ReplaceAll(url.PathEscape(fileID), ":", "%3A")
	return fmt.Sprintf(
		"%sprojects/%s/locations/%s/repositories/%s/files/%s",
		base, project, location, repository, encoded,
	)
}

func sha256Hashes(b []byte) (hexStr, b64Str string) {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:]), base64.StdEncoding.EncodeToString(sum[:])
}

func init() {
	registry.Schema{
		Name:        "google_artifact_registry_file",
		ProductName: "artifactregistry",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceArtifactRegistryFile(),
	}.Register()
}
