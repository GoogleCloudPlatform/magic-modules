package storagecontrol

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageControlProjectIntelligenceFindingRevision() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageControlProjectIntelligenceFindingRevisionRead,
		Schema: map[string]*schema.Schema{
			"finding_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the intelligence finding.`,
			},
			"revision_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the finding revision.`,
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				Description: `The location of the intelligence finding. Currently default value is global and users cannot use for input for now.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource name of the finding revision.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the finding revision was created.`,
			},
			"snapshot": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The snapshot of the finding at revision creation time.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource name of the finding.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `A short description of the finding.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of this finding.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The category of the finding.`,
						},
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The severity of the finding.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the finding was created.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the finding was last updated.`,
						},
						"target_resource": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The fully qualified resource name of the resource that this IntelligenceFinding applies to.`,
						},
						"associated_resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Google Cloud resource names that are relevant to the IntelligenceFinding. This list also includes the targetResource.`,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"observation_period": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The time interval from which the underlying data generated this IntelligenceFinding was observed.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"coldline_and_archival_storage_operations_spike": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A finding about a spike in Class A or Class B operations on Coldline or Archive Cloud Storage objects.",
							Elem: &schema.Resource{
								Schema: storageControlColdlineSpikeSchema(),
							},
						},
						"throttled_requests_spike": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A finding about a spike in throttled requests (429 errors) within a project.",
							Elem: &schema.Resource{
								Schema: storageControlThrottledRequestsSpikeSchema(),
							},
						},
						"cross_region_egress_spike": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A finding about a spike in cross-region egress from Cloud Storage.",
							Elem: &schema.Resource{
								Schema: storageControlCrossRegionEgressSpikeSchema(),
							},
						},
						"storage_growth_above_trend": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A finding about a spike in storage growth (bytes or object count) that is outside the normal historical trend.",
							Elem: &schema.Resource{
								Schema: storageControlStorageGrowthSpikeSchema(),
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleStorageControlProjectIntelligenceFindingRevisionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for intelligence finding revision: %s", err)
	}
	location := d.Get("location").(string)
	findingId := d.Get("finding_id").(string)
	revisionId := d.Get("revision_id").(string)

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf(transport_tpg.BaseUrl(Product, config)+"projects/%s/locations/%s/intelligenceFindings/%s/revisions/%s", project, location, findingId, revisionId))
	if err != nil {
		return fmt.Errorf("Error formatting url for intelligence finding revision: %s", err)
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, "StorageControlProjectIntelligenceFindingRevision", fmt.Sprintf("StorageControlProjectIntelligenceFindingRevision %s/revisions/%s", findingId, revisionId))
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("create_time", res["createTime"]); err != nil {
		return fmt.Errorf("Error setting create_time: %s", err)
	}
	if err := d.Set("snapshot", flattenStorageControlFindingSnapshot(res["snapshot"])); err != nil {
		return fmt.Errorf("Error setting snapshot: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/intelligenceFindings/%s/revisions/%s", project, location, findingId, revisionId))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_storage_control_project_intelligence_finding_revision",
		ProductName: "storagecontrol",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleStorageControlProjectIntelligenceFindingRevision(),
	}.Register()
}
