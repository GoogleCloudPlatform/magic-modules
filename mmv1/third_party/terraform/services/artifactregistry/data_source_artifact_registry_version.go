package artifactregistry

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type Version struct {
	name        string
	description string
	tags        string
	createTime  string
	updateTime  string
	annotations map[string]string
}

func DataSourceArtifactRegistryVersion() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceArtifactRegistryVersionRead,

		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"view": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "BASIC",
				ValidateFunc: validateViewArtifactRegistryVersion,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceArtifactRegistryVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return fmt.Errorf("Error setting Artifact Registry user agent: %s", err)
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error setting Artifact Registry project: %s", err)
	}

	basePath, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}")
	if err != nil {
		return fmt.Errorf("Error setting Artifact Registry base path: %s", err)
	}

	resourcePath, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}/packages/{{package_name}}/versions/{{version_name}}"))
	if err != nil {
		return fmt.Errorf("Error setting resource path: %s", err)
	}

	urlRequest := basePath + resourcePath
	headers := make(http.Header)

	u, err := url.Parse(urlRequest)
	if err != nil {
		return fmt.Errorf("Error parsing URL: %s", err)
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    u.String(),
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return fmt.Errorf("Error getting Artifact Registry version: %s", err)
	}

	annotations := make(map[string]string)
	if anno, ok := res["annotations"].(map[string]interface{}); ok {
		for k, v := range anno {
			if val, ok := v.(string); ok {
				annotations[k] = val
			}
		}
	}

	getString := func(m map[string]interface{}, key string) string {
		if v, ok := m[key].(string); ok {
			return v
		}
		return ""
	}

	name := getString(res, "name")

	if err := d.Set("project", project); err != nil {
		return err
	}
	if err := d.Set("name", name); err != nil {
		return err
	}
	if err := d.Set("version", res["version"].(string)); err != nil {
		return err
	}

	d.SetId(name)

	return nil
}

func validateViewArtifactRegistryVersion(val interface{}, key string) ([]string, []error) {
	v := val.(string)
	var errs []error

	if v != "BASIC" && v != "FULL" {
		errs = append(errs, fmt.Errorf("%q must be either 'BASIC' or 'FULL', got %q", key, v))
	}

	return nil, errs
}
