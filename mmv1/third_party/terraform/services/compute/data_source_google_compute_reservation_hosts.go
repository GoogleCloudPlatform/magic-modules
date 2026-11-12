package compute

import (
	"fmt"
	neturl "net/url"
	"strings"

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
						"status": computeReservationHostsStatusSchema(),
					},
				},
			},
		},
	}
}

func computeReservationHostsStatusSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Computed:    true,
		Description: "The status of the host.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"physical_topology": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "The physical topology of the host, if present.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"cluster": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The cluster name of the reservation sub-block.",
							},
							"block": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The unique identifier of the capacity block within the cluster.",
							},
							"sub_block": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The unique identifier of the capacity sub-block within the capacity block.",
							},
							"host": {
								Type:        schema.TypeString,
								Computed:    true,
								Description: "The unique identifier of the capacity host within the capacity sub-block.",
							},
						},
					},
				},
				"running_instances": {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "The URIs of the instances currently running on this host.",
					Elem: &schema.Schema{
						Type: schema.TypeString,
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

	association, err := buildComputeReservationHostsAssociation(d)
	if err != nil {
		return err
	}

	baseListURL := buildComputeReservationHostsURL(transport_tpg.BaseUrl(Product, config), project, zone, association)
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
				hosts = append(hosts, flattenComputeReservationHostsItem(host))
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

	d.SetId(buildComputeReservationHostsID(project, zone, association))
	return nil
}

func buildComputeReservationHostsAssociation(d *schema.ResourceData) (string, error) {
	var reservation, reservationBlock, reservationSubBlock string
	if v, ok := d.GetOk("reservation"); ok {
		reservation = tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(v.(string)))
	}
	if v, ok := d.GetOk("reservation_block"); ok {
		reservationBlock = tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(v.(string)))
	}
	if v, ok := d.GetOk("reservation_sub_block"); ok {
		reservationSubBlock = tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(v.(string)))
	}

	return formatComputeReservationHostsAssociation(reservation, reservationBlock, reservationSubBlock)
}

func formatComputeReservationHostsAssociation(reservation, reservationBlock, reservationSubBlock string) (string, error) {
	if reservationBlock != "" && reservation == "" {
		return "", fmt.Errorf("reservation must be specified when reservation_block is set")
	}
	if reservationSubBlock != "" && (reservation == "" || reservationBlock == "") {
		return "", fmt.Errorf("reservation and reservation_block must be specified when reservation_sub_block is set")
	}

	if reservation == "" {
		return "", nil
	}
	parts := []string{"reservations", reservation}
	if reservationBlock != "" {
		parts = append(parts, "reservationBlocks", reservationBlock)
	}
	if reservationSubBlock != "" {
		parts = append(parts, "reservationSubBlocks", reservationSubBlock)
	}
	return strings.Join(parts, "/"), nil
}

func buildComputeReservationHostsURL(baseURL, project, zone, association string) string {
	escapedProject := neturl.PathEscape(project)
	escapedZone := neturl.PathEscape(zone)
	if association == "" {
		return fmt.Sprintf("%sprojects/%s/zones/%s/hosts", baseURL, escapedProject, escapedZone)
	}
	return fmt.Sprintf("%sprojects/%s/zones/%s/%s/hosts", baseURL, escapedProject, escapedZone, neturl.PathEscape(association))
}

func buildComputeReservationHostsID(project, zone, association string) string {
	if association == "" {
		return fmt.Sprintf("projects/%s/zones/%s/hosts", project, zone)
	}
	return fmt.Sprintf("projects/%s/zones/%s/%s/hosts", project, zone, association)
}

func flattenComputeReservationHostsItem(host map[string]interface{}) map[string]interface{} {
	zoneName := ""
	if zoneStr, ok := host["zone"].(string); ok && zoneStr != "" {
		zoneName = tpgresource.GetResourceNameFromSelfLink(zoneStr)
	}

	return map[string]interface{}{
		"name":               host["name"],
		"kind":               host["kind"],
		"resource_id":        flattenComputeReservationHostsResourceID(host["id"]),
		"creation_timestamp": host["creationTimestamp"],
		"zone":               zoneName,
		"self_link":          host["selfLink"],
		"self_link_with_id":  host["selfLinkWithId"],
		"alias_links":        flattenComputeReservationHostsStringList(host["aliasLinks"]),
		"state":              host["state"],
		"status":             flattenComputeReservationHostsStatus(host["status"]),
	}
}

func flattenComputeReservationHostsResourceID(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func flattenComputeReservationHostsStringList(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	switch items := v.(type) {
	case []string:
		out := make([]string, len(items))
		copy(out, items)
		return out
	case []interface{}:
		out := make([]string, 0, len(items))
		for _, item := range items {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	default:
		return []string{}
	}
}

func flattenComputeReservationHostsStatus(v interface{}) []map[string]interface{} {
	statusMap, ok := v.(map[string]interface{})
	if !ok || statusMap == nil {
		return nil
	}

	var topologyList []map[string]interface{}
	if ptMap, ok := statusMap["physicalTopology"].(map[string]interface{}); ok && ptMap != nil {
		topologyList = []map[string]interface{}{
			{
				"cluster":   ptMap["cluster"],
				"block":     ptMap["block"],
				"sub_block": ptMap["subBlock"],
				"host":      ptMap["host"],
			},
		}
	}

	return []map[string]interface{}{
		{
			"physical_topology": topologyList,
			"running_instances": flattenComputeReservationHostsStringList(statusMap["runningInstances"]),
		},
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
