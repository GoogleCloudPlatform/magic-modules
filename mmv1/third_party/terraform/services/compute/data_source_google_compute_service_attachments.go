package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeServiceAttachments() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleComputeServiceAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_preference": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nat_subnets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enable_proxy_protocol": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"domain_names": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"fingerprint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleComputeServiceAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	if filter, ok := d.GetOk("filter"); ok {
		params["filter"] = filter.(string)
	}

	serviceAttachments := make([]map[string]interface{}, 0)

	for {
		url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/serviceAttachments")
		if err != nil {
			return err
		}

		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("error retrieving service attachments: %s", err)
		}

		pageItems := flattenDatasourceGoogleComputeServiceAttachmentsList(res["items"])
		serviceAttachments = append(serviceAttachments, pageItems...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("service_attachments", serviceAttachments); err != nil {
		return fmt.Errorf("error setting service_attachments: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/serviceAttachments", project, region))

	return nil
}

func flattenDatasourceGoogleComputeServiceAttachmentsList(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	serviceAttachments := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		p := raw.(map[string]interface{})

		var natSubnets []interface{}
		if ns, ok := p["natSubnets"]; ok && ns != nil {
			natSubnets = ns.([]interface{})
		}

		var domainNames []interface{}
		if dn, ok := p["domainNames"]; ok && dn != nil {
			domainNames = dn.([]interface{})
		}

		var enableProxyProtocol bool
		if epp, ok := p["enableProxyProtocol"]; ok && epp != nil {
			enableProxyProtocol = epp.(bool)
		}

		serviceAttachments = append(serviceAttachments, map[string]interface{}{
			"name":                  p["name"],
			"description":           p["description"],
			"self_link":             p["selfLink"],
			"target_service":        p["targetService"],
			"connection_preference": p["connectionPreference"],
			"nat_subnets":           natSubnets,
			"enable_proxy_protocol": enableProxyProtocol,
			"domain_names":          domainNames,
			"fingerprint":           p["fingerprint"],
			"region":                p["region"],
		})
	}

	return serviceAttachments
}

func init() {
	registry.Schema{
		Name:        "google_compute_service_attachments",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeServiceAttachments(),
	}.Register()
}
