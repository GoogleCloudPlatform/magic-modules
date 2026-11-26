package compute

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceGoogleComputeRoutes returns the schema of the google_compute_routes data source.
func DataSourceGoogleComputeRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeRoutesRead,

		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"routes": {
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
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dest_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"next_hop_gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_instance": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_vpn_tunnel": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_ilb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_peering": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_hub": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_interconnect_attachment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_origin": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// next_hop_med and next_hop_inter_region_cost are strings to match google_compute_route; the API declares integers.
						"next_hop_med": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop_inter_region_cost": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceGoogleComputeRoutesRead lists every route in the project, following pages and applying the optional filter.
func dataSourceGoogleComputeRoutesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	baseURL, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/routes")
	if err != nil {
		return err
	}

	params := map[string]string{}
	if filter := d.Get("filter").(string); filter != "" {
		params["filter"] = filter
	}

	routes := make([]map[string]interface{}, 0)
	for {
		url, err := transport_tpg.AddQueryParams(baseURL, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error retrieving list of routes: %w", err)
		}

		items, _ := res["items"].([]interface{})
		for _, item := range items {
			route, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Error reading list of routes: expected an object, got %T", item)
			}
			routes = append(routes, flattenComputeRoutesItem(route))
		}

		pageToken, _ := res["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
		params["pageToken"] = pageToken
	}

	if err := d.Set("routes", routes); err != nil {
		return fmt.Errorf("Error setting routes: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/global/routes", project))

	return nil
}

// flattenComputeRoutesItem maps one route from the API response to the schema of a routes entry.
func flattenComputeRoutesItem(route map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":                             route["name"],
		"description":                      route["description"],
		"network":                          route["network"],
		"dest_range":                       route["destRange"],
		"priority":                         route["priority"],
		"tags":                             route["tags"],
		"next_hop_gateway":                 route["nextHopGateway"],
		"next_hop_instance":                route["nextHopInstance"],
		"next_hop_ip":                      route["nextHopIp"],
		"next_hop_vpn_tunnel":              route["nextHopVpnTunnel"],
		"next_hop_ilb":                     route["nextHopIlb"],
		"next_hop_network":                 route["nextHopNetwork"],
		"next_hop_peering":                 route["nextHopPeering"],
		"next_hop_hub":                     route["nextHopHub"],
		"next_hop_interconnect_attachment": route["nextHopInterconnectAttachment"],
		"next_hop_origin":                  route["nextHopOrigin"],
		"next_hop_med":                     flattenComputeRoutesUint32(route["nextHopMed"]),
		"next_hop_inter_region_cost":       flattenComputeRoutesUint32(route["nextHopInterRegionCost"]),
		"route_type":                       route["routeType"],
		"route_status":                     route["routeStatus"],
		"self_link":                        route["selfLink"],
		"creation_timestamp":               route["creationTimestamp"],
	}
}

// flattenComputeRoutesUint32 returns a 32-bit unsigned value as decimal text, whether the API sent a number or a string.
func flattenComputeRoutesUint32(v interface{}) interface{} {
	if f, ok := v.(float64); ok {
		return strconv.FormatFloat(f, 'f', 0, 64)
	}
	return v
}

func init() {
	registry.Schema{
		Name:        "google_compute_routes",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeRoutes(),
	}.Register()
}
