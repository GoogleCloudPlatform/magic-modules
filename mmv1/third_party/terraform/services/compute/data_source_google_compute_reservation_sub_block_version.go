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

func DataSourceGoogleComputeReservationSubBlockVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeReservationSubBlockVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name or resource ID of the reservation sub-block.",
			},
			"reservation_block": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the parent reservation block.",
			},
			"reservation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the parent reservation.",
			},
			"sbom_selections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The SBOM selections to return. Possible values include `SBOM_SELECTION_CURRENT` and `SBOM_SELECTION_TARGET`. See https://cloud.google.com/compute/docs/reference/rest/v1/reservationSubBlocks/getVersion for details.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The zone where the reservation sub-block resides.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project in which the resource belongs.",
			},
			"current_component_versions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A mapping of components to their currently-applied versions or other appropriate identifiers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target_component_versions": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A mapping of components to their target versions or other appropriate identifiers.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGoogleComputeReservationSubBlockVersionRead(d *schema.ResourceData, meta interface{}) error {
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

	subBlockName := tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(d.Get("name").(string)))
	reservationBlock := tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(d.Get("reservation_block").(string)))
	reservation := tpgresource.GetResourceNameFromSelfLink(strings.TrimSpace(d.Get("reservation").(string)))
	if subBlockName == "" || reservationBlock == "" || reservation == "" {
		return fmt.Errorf("name, reservation_block, and reservation must not be empty")
	}

	parentName := fmt.Sprintf("reservations/%s/reservationBlocks/%s", reservation, reservationBlock)
	rawURL := fmt.Sprintf(
		"%sprojects/%s/zones/%s/%s/reservationSubBlocks/%s/getVersion",
		transport_tpg.BaseUrl(Product, config),
		neturl.PathEscape(project),
		neturl.PathEscape(zone),
		neturl.PathEscape(parentName),
		neturl.PathEscape(subBlockName),
	)
	body := expandComputeReservationSubBlockVersionSbomSelectionsBody(d)

	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    rawURL,
		UserAgent: userAgent,
		Body:      body,
		Timeout:   d.Timeout(schema.TimeoutRead),
	})
	if err != nil {
		return fmt.Errorf("Error getting version for ReservationSubBlock %s: %w", subBlockName, err)
	}

	currentVersions, targetVersions, err := waitForComputeReservationSubBlockVersionOperation(
		config,
		op,
		project,
		zone,
		fmt.Sprintf("getting version for ReservationSubBlock %s", subBlockName),
		userAgent,
		d.Timeout(schema.TimeoutRead),
	)
	if err != nil {
		return fmt.Errorf("Error waiting to get version for ReservationSubBlock %s: %w", subBlockName, err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}
	if err := d.Set("current_component_versions", currentVersions); err != nil {
		return fmt.Errorf("Error setting current_component_versions: %s", err)
	}
	if err := d.Set("target_component_versions", targetVersions); err != nil {
		return fmt.Errorf("Error setting target_component_versions: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s/reservationBlocks/%s/reservationSubBlocks/%s/getVersion", project, zone, reservation, reservationBlock, subBlockName))
	return nil
}

func expandComputeReservationSubBlockVersionSbomSelectionsBody(d *schema.ResourceData) map[string]interface{} {
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

func waitForComputeReservationSubBlockVersionOperation(config *transport_tpg.Config, rawOp map[string]interface{}, project, zone, activity, userAgent string, timeout time.Duration) (map[string]interface{}, map[string]interface{}, error) {
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

	return flattenComputeReservationSubBlockVersionOperationMetadata(w.Op)
}

func flattenComputeReservationSubBlockVersionOperationMetadata(op map[string]interface{}) (map[string]interface{}, map[string]interface{}, error) {
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
		Name:        "google_compute_reservation_sub_block_version",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeReservationSubBlockVersion(),
	}.Register()
}
