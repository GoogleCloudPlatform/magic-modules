package compute

import (
	"fmt"
	neturl "net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeReservationHost() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeReservationHostRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or resource ID of the host.",
			},
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
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The zone where the host resides.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project in which the resource belongs.",
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
	}
}

func computeHostStatusSchema() *schema.Schema {
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

func dataSourceGoogleComputeReservationHostRead(d *schema.ResourceData, meta interface{}) error {
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

	hostName := tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(d.Get("name").(string)))
	if hostName == "" {
		return fmt.Errorf("name must not be empty")
	}

	rawURL := buildComputeHostsURL(transport_tpg.BaseUrl(Product, config), project, zone, association, "/"+neturl.PathEscape(hostName))

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    rawURL,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading Host %s: %w", hostName, err)
	}
	if res == nil {
		return fmt.Errorf("Host %s not found", hostName)
	}

	if resource, ok := res["resource"].(map[string]interface{}); ok {
		res = resource
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("kind", res["kind"]); err != nil {
		return fmt.Errorf("Error setting kind: %s", err)
	}
	if err := d.Set("resource_id", flattenComputeHostResourceID(res["id"])); err != nil {
		return fmt.Errorf("Error setting resource_id: %s", err)
	}
	if err := d.Set("creation_timestamp", res["creationTimestamp"]); err != nil {
		return fmt.Errorf("Error setting creation_timestamp: %s", err)
	}
	if err := d.Set("self_link", res["selfLink"]); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("self_link_with_id", res["selfLinkWithId"]); err != nil {
		return fmt.Errorf("Error setting self_link_with_id: %s", err)
	}
	if err := d.Set("alias_links", flattenComputeStringList(res["aliasLinks"])); err != nil {
		return fmt.Errorf("Error setting alias_links: %s", err)
	}
	if err := d.Set("state", res["state"]); err != nil {
		return fmt.Errorf("Error setting state: %s", err)
	}
	if err := d.Set("status", flattenComputeHostStatus(res["status"])); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}

	d.SetId(buildComputeHostsID(project, zone, association, "/"+hostName))
	return nil
}

func buildComputeHostAssociation(d *schema.ResourceData) (string, error) {
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

	return formatComputeHostAssociation(reservation, reservationBlock, reservationSubBlock)
}

func formatComputeHostAssociation(reservation, reservationBlock, reservationSubBlock string) (string, error) {
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

func buildComputeHostsURL(baseURL, project, zone, association, suffix string) string {
	escapedProject := neturl.PathEscape(project)
	escapedZone := neturl.PathEscape(zone)
	if association == "" {
		return fmt.Sprintf("%sprojects/%s/zones/%s/hosts%s", baseURL, escapedProject, escapedZone, suffix)
	}
	return fmt.Sprintf("%sprojects/%s/zones/%s/%s/hosts%s", baseURL, escapedProject, escapedZone, neturl.PathEscape(association), suffix)
}

func buildComputeHostsID(project, zone, association, suffix string) string {
	if association == "" {
		return fmt.Sprintf("projects/%s/zones/%s/hosts%s", project, zone, suffix)
	}
	return fmt.Sprintf("projects/%s/zones/%s/%s/hosts%s", project, zone, association, suffix)
}

func flattenComputeHostResourceID(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func flattenComputeStringList(v interface{}) []string {
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

func flattenComputeHostStatus(v interface{}) []map[string]interface{} {
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
			"running_instances": flattenComputeStringList(statusMap["runningInstances"]),
		},
	}
}

func expandComputeSbomSelectionsBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{}
	if v, ok := d.GetOk("sbom_selections"); ok {
		rawList := v.([]interface{})
		selections := make([]string, 0, len(rawList))
		for _, item := range rawList {
			if s, ok := item.(string); ok && s != "" {
				selections = append(selections, s)
			}
		}
		if len(selections) > 0 {
			body["sbomSelections"] = selections
		}
	}
	return body
}

func waitForComputeZonalGetVersionOperation(config *transport_tpg.Config, rawOp map[string]interface{}, project, zone, activity, userAgent string, timeout time.Duration) (map[string]interface{}, map[string]interface{}, error) {
	if rawOp == nil {
		return nil, nil, fmt.Errorf("empty operation response while %s", activity)
	}
	if innerOp, ok := rawOp["operation"].(map[string]interface{}); ok && innerOp != nil {
		rawOp = innerOp
	}
	if z, ok := rawOp["zone"].(string); !ok || z == "" {
		rawOp["zone"] = zone
	}

	w := &ComputeOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Context:   config.Context,
		Op:        rawOp,
		Project:   project,
	}
	if err := w.SetOp(rawOp); err != nil {
		return nil, nil, err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return nil, nil, err
	}

	return flattenComputeGetVersionOperationMetadata(w.Op)
}

func flattenComputeGetVersionOperationMetadata(op map[string]interface{}) (map[string]interface{}, map[string]interface{}, error) {
	currentVersions := map[string]interface{}{}
	targetVersions := map[string]interface{}{}

	if op == nil {
		return currentVersions, targetVersions, nil
	}
	metaMap, ok := op["getVersionOperationMetadata"].(map[string]interface{})
	if !ok || metaMap == nil {
		return currentVersions, targetVersions, nil
	}
	sbomInfo, ok := metaMap["inlineSbomInfo"].(map[string]interface{})
	if !ok || sbomInfo == nil {
		return currentVersions, targetVersions, nil
	}
	if cur, ok := sbomInfo["currentComponentVersions"].(map[string]interface{}); ok && cur != nil {
		for k, v := range cur {
			currentVersions[k] = fmt.Sprintf("%v", v)
		}
	}
	if tgt, ok := sbomInfo["targetComponentVersions"].(map[string]interface{}); ok && tgt != nil {
		for k, v := range tgt {
			targetVersions[k] = fmt.Sprintf("%v", v)
		}
	}
	return currentVersions, targetVersions, nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_reservation_host",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeReservationHost(),
	}.Register()
}
