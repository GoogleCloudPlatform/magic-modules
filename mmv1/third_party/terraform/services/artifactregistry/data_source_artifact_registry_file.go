package artifactregistry

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

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
