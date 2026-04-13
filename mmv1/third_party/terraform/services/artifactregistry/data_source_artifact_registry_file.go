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
