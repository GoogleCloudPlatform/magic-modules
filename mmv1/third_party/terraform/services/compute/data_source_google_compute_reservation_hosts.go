package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeReservationHosts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeReservationHostsRead,

		Schema: map[string]*schema.Schema{
			"reservation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the parent reservation.",
			},
			"reservation_block": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"reservation"},
				Description:  "The name of the parent reservation block. Requires `reservation` to be set.",
			},
			"reservation_sub_block": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"reservation", "reservation_block"},
				Description:  "The name of the parent reservation sub-block. Requires `reservation` and `reservation_block` to be set.",
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A filter expression that filters resources listed in the response.",
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The zone where the hosts reside.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project in which the resource belongs.",
			},
			"hosts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of hosts matching the query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the host.",
						},
						"kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource. Always compute#host for hosts.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique server-generated identifier for the host.",
						},
						"creation_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp in RFC3339 text format.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The zone in which the host resides.",
						},
						"self_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server-defined fully-qualified URL for this resource.",
						},
						"self_link_with_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server-defined URL for this resource with the resource ID.",
						},
						"alias_links": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "All alias resource paths for this host.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the host.",
						},
						"status": computeHostStatusSchema(),
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeReservationHostsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}

	billingProject := project
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	association, err := buildComputeHostAssociation(d)
	if err != nil {
		return err
	}

	baseListURL := buildComputeHostsURL(transport_tpg.BaseUrl(Product, config), project, zone, association, "")
	hosts := make([]map[string]interface{}, 0)
	pageToken := ""

	for {
		params := map[string]string{}
		if filter, ok := d.GetOk("filter"); ok {
			params["filter"] = filter.(string)
		}
		if pageToken != "" {
			params["pageToken"] = pageToken
		}

		url, err := transport_tpg.AddQueryParams(baseListURL, params)
		if err != nil {
			return err
		}

		resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error listing Hosts: %w", err)
		}

		if items, ok := resp["items"].([]interface{}); ok {
			for _, raw := range items {
				host, ok := raw.(map[string]interface{})
				if !ok || host == nil {
					continue
				}
				hosts = append(hosts, flattenComputeHostItem(host))
			}
		}

		pageToken, _ = resp["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("hosts", hosts); err != nil {
		return fmt.Errorf("Error setting hosts: %s", err)
	}

	d.SetId(buildComputeHostsID(project, zone, association, ""))
	return nil
}

func flattenComputeHostItem(host map[string]interface{}) map[string]interface{} {
	zoneName := ""
	if zoneStr, ok := host["zone"].(string); ok && zoneStr != "" {
		zoneName = tpgresource.GetResourceNameFromSelfLink(zoneStr)
	}

	return map[string]interface{}{
		"name":               host["name"],
		"kind":               host["kind"],
		"resource_id":        flattenComputeHostResourceID(host["id"]),
		"creation_timestamp": host["creationTimestamp"],
		"zone":               zoneName,
		"self_link":          host["selfLink"],
		"self_link_with_id":  host["selfLinkWithId"],
		"alias_links":        flattenComputeStringList(host["aliasLinks"]),
		"state":              host["state"],
		"status":             flattenComputeHostStatus(host["status"]),
	}
}

func init() {
	registry.Schema{
		Name:        "google_compute_reservation_hosts",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeReservationHosts(),
	}.Register()
}
