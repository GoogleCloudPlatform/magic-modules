package artifactregistry

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceArtifactRegistryPythonPackages() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArtifactRegistryPythonPackagesRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"python_packages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"package_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArtifactRegistryPythonPackagesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	basePath, err := tpgresource.ReplaceVars(d, config, "{{ArtifactRegistryBasePath}}")
	if err != nil {
		return fmt.Errorf("Error setting Artifact Registry base path: %s", err)
	}

	resourcePath, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("projects/{{project}}/locations/{{location}}/repositories/{{repository_id}}/pythonPackages"))
	if err != nil {
		return fmt.Errorf("Error setting resource path: %s", err)
	}

	urlRequest := basePath + resourcePath

	headers := make(http.Header)
	pythonPackages := make([]map[string]interface{}, 0)
	pageToken := ""

	for {
		u, err := url.Parse(urlRequest)
		if err != nil {
			return fmt.Errorf("Error parsing URL: %s", err)
		}

		q := u.Query()
		if pageToken != "" {
			q.Set("pageToken", pageToken)
		}
		u.RawQuery = q.Encode()

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    u.String(),
			UserAgent: userAgent,
			Headers:   headers,
		})

		if err != nil {
			return fmt.Errorf("Error listing Artifact Registry Python packages: %s", err)
		}

		if items, ok := res["pythonPackages"].([]interface{}); ok {
			for _, item := range items {
				pkg := item.(map[string]interface{})

				name, ok := pkg["name"].(string)
				if !ok {
					return fmt.Errorf("Error getting Artifact Registry Python package name: %s", err)
				}

				lastComponent := name[strings.LastIndex(name, "/")+1:]
				packageName := strings.SplitN(lastComponent, ":", 2)[0]

				getString := func(m map[string]interface{}, key string) string {
					if v, ok := m[key].(string); ok {
						return v
					}
					return ""
				}

				pythonPackages = append(pythonPackages, map[string]interface{}{
					"package_name": packageName,
					"name":         name,
					"version":      getString(pkg, "version"),
					"create_time":  getString(pkg, "createTime"),
					"update_time":  getString(pkg, "updateTime"),
				})
			}
		}

		if nextToken, ok := res["nextPageToken"].(string); ok && nextToken != "" {
			pageToken = nextToken
		} else {
			break
		}
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if err := d.Set("python_packages", pythonPackages); err != nil {
		return fmt.Errorf("Error setting Artifact Registry Python packages: %s", err)
	}

	d.SetId(resourcePath)

	return nil
}
