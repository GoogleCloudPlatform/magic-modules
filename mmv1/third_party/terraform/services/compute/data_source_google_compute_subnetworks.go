package compute

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeSubnetworks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeSubnetworksRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnetworks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_cidr_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_self_link": {
							Type:     schema.TypeString,
							Computed: true,
							// TODO: remove in next major release (7.0.0) also from docs and implementation below
							Deprecated: "Use `network_name` instead. This field will be removed in a future major release.",
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_google_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeSubnetworksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for subnetwork: %s", err)
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching region for subnetwork: %s", err)
	}

	filter := d.Get("filter").(string)

	subnetworks := make([]map[string]interface{}, 0)

	baseURL, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/regions/{{region}}/subnetworks")
	if err != nil {
		return err
	}

	params := map[string]string{}
	if filter != "" {
		params["filter"] = filter
	}
	url, err := transport_tpg.AddQueryParams(baseURL, params)
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Subnetworks : %s %s", project, region))
	}

	if items, ok := res["items"].([]interface{}); ok {
		for _, item := range items {
			subnet, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			network, _ := subnet["network"].(string)
			subnetworks = append(subnetworks, map[string]interface{}{
				"description":              subnet["description"],
				"ip_cidr_range":            subnet["ipCidrRange"],
				"name":                     subnet["name"],
				"network_self_link":        filepath.Base(network),
				"network":                  subnet["network"],
				"network_name":             filepath.Base(network),
				"private_ip_google_access": subnet["privateIpGoogleAccess"],
				"self_link":                subnet["selfLink"],
			})
		}
	}

	if err := d.Set("subnetworks", subnetworks); err != nil {
		return fmt.Errorf("Error retrieving subnetworks: %s", err)
	}

	d.SetId(fmt.Sprintf(
		"projects/%s/regions/%s/subnetworks",
		project,
		region,
	))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_subnetworks",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeSubnetworks(),
	}.Register()
}
