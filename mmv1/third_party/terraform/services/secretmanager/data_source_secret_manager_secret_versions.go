// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0
package secretmanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceSecretManagerSecretVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretManagerSecretVersionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"destroy_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSecretManagerSecretVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for SecretVersion: %s", err)
	}

	secret := d.Get("secret").(string)
	if !strings.HasPrefix(secret, "projects/") {
		secret = fmt.Sprintf("projects/%s/secrets/%s", project, secret)
	}

	url := fmt.Sprintf("%s%s/versions", transport_tpg.BaseUrl(Product, config), secret)

	filter, hasFilter := d.GetOk("filter")
	if hasFilter {
		url, err = transport_tpg.AddQueryParams(url, map[string]string{"filter": filter.(string)})
		if err != nil {
			return err
		}
	}

	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	allVersions := make([]interface{}, 0)
	token := ""
	for paginate := true; paginate; {
		if token != "" {
			url, err = transport_tpg.AddQueryParams(url, map[string]string{"pageToken": token})
			if err != nil {
				return err
			}
		}
		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SecretManagerSecretVersions %q", d.Id()))
		}
		versionsInterface := resp["versions"]
		if versionsInterface == nil {
			break
		}
		allVersions = append(allVersions, versionsInterface.([]interface{})...)
		tokenInterface := resp["nextPageToken"]
		if tokenInterface == nil {
			paginate = false
		} else {
			paginate = true
			token = tokenInterface.(string)
		}
	}

	if err := d.Set("versions", flattenSecretManagerSecretVersionsList(allVersions)); err != nil {
		return fmt.Errorf("error setting versions: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	id := fmt.Sprintf("%s/versions", secret)
	if hasFilter {
		id += "/filter=" + filter.(string)
	}
	d.SetId(id)

	return nil
}

func flattenSecretManagerSecretVersionsList(v interface{}) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			continue
		}
		name := ""
		if n, ok := original["name"].(string); ok {
			name = n
		}
		version := ""
		if name != "" {
			parts := strings.Split(name, "/")
			version = parts[len(parts)-1]
		}
		enabled := false
		if state, ok := original["state"].(string); ok {
			enabled = state == "ENABLED"
		}
		transformed = append(transformed, map[string]interface{}{
			"name":         name,
			"version":      version,
			"enabled":      enabled,
			"create_time":  original["createTime"],
			"destroy_time": original["destroyTime"],
		})
	}
	return transformed
}

func init() {
	registry.Schema{
		Name:        "google_secret_manager_secret_versions",
		ProductName: "secretmanager",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceSecretManagerSecretVersions(),
	}.Register()
}
