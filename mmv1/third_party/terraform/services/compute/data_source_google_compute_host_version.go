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

func DataSourceGoogleComputeHostVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeHostVersionRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

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
			"sbom_selections": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The SBOM selections to return. Possible values include `SBOM_SELECTION_CURRENT` and `SBOM_SELECTION_TARGET`. See https://cloud.google.com/compute/docs/reference/rest/v1/hosts/getVersion for details.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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

func dataSourceGoogleComputeHostVersionRead(d *schema.ResourceData, meta interface{}) error {
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

	rawURL := buildComputeHostsURL(transport_tpg.BaseUrl(Product, config), project, zone, association, "/"+neturl.PathEscape(hostName)+"/getVersion")
	body := expandComputeSbomSelectionsBody(d)

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
		return fmt.Errorf("Error getting version for Host %s: %w", hostName, err)
	}

	currentVersions, targetVersions, err := waitForComputeZonalGetVersionOperation(
		config,
		op,
		project,
		zone,
		fmt.Sprintf("getting version for Host %s", hostName),
		userAgent,
		d.Timeout(schema.TimeoutRead),
	)
	if err != nil {
		return fmt.Errorf("Error waiting to get version for Host %s: %w", hostName, err)
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

	d.SetId(buildComputeHostsID(project, zone, association, "/"+hostName+"/getVersion"))
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_host_version",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeHostVersion(),
	}.Register()
}
