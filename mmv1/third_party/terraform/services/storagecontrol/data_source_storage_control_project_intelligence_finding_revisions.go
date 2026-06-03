package storagecontrol

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageControlProjectIntelligenceFindingRevisions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageControlProjectIntelligenceFindingRevisionsRead,
		Schema: map[string]*schema.Schema{
			"finding_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the intelligence finding.`,
			},
			"page_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: `The maximum number of IntelligenceFindingRevision resources to return.`,
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
			"revisions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of intelligence finding revisions.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}

func dataSourceGoogleStorageControlProjectIntelligenceFindingRevisionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for intelligence finding revisions: %s", err)
	}
	location := d.Get("location").(string)
	findingId := d.Get("finding_id").(string)

	params := make(map[string]string)
	if v, ok := d.GetOk("page_size"); ok {
		params["pageSize"] = strconv.Itoa(v.(int))
	}

	revisions := make([]map[string]interface{}, 0)

	for {
		url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf(transport_tpg.BaseUrl(Product, config)+"projects/%s/locations/%s/intelligenceFindings/%s/revisions", project, location, findingId))
		if err != nil {
			return fmt.Errorf("Error formatting url for intelligence finding revisions: %s", err)
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
			return fmt.Errorf("Error retrieving intelligence finding revisions: %s", err)
		}

		var items interface{}
		if v, ok := res["intelligenceFindingRevisions"]; ok {
			items = v
		}

		pageRevisions := flattenStorageControlIntelligenceFindingRevisionsList(items)
		revisions = append(revisions, pageRevisions...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("revisions", revisions); err != nil {
		return fmt.Errorf("Error setting revisions: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/intelligenceFindings/%s/revisions", project, location, findingId))

	return nil
}

func flattenStorageControlIntelligenceFindingRevisionsList(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls, ok := v.([]interface{})
	if !ok {
		return make([]map[string]interface{}, 0)
	}
	revisions := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		o, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}

		revision := map[string]interface{}{
			"name":        o["name"],
			"create_time": o["createTime"],
			"snapshot":    flattenStorageControlFindingSnapshot(o["snapshot"]),
		}
		revisions = append(revisions, revision)
	}

	return revisions
}

func flattenStorageControlFindingSnapshot(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	o, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}

	finding := map[string]interface{}{
		"name":                 o["name"],
		"description":          o["description"],
		"type":                 o["type"],
		"category":             o["category"],
		"severity":             o["severity"],
		"create_time":          o["createTime"],
		"update_time":          o["updateTime"],
		"target_resource":      o["targetResource"],
		"associated_resources": flattenStorageControlStringList(o["associatedResources"]),
		"observation_period":   flattenStorageControlObservationPeriod(o["observationPeriod"]),
		"coldline_and_archival_storage_operations_spike": flattenStorageControlColdlineSpike(o["coldlineAndArchivalStorageOperationsSpike"]),
		"throttled_requests_spike":                       flattenStorageControlThrottledRequestsSpike(o["throttledRequestsSpike"]),
		"cross_region_egress_spike":                      flattenStorageControlCrossRegionEgressSpike(o["crossRegionEgressSpike"]),
		"storage_growth_above_trend":                     flattenStorageControlStorageGrowthSpike(o["storageGrowthAboveTrend"]),
	}

	return []map[string]interface{}{finding}
}

func init() {
	registry.Schema{
		Name:        "google_storage_control_project_intelligence_finding_revisions",
		ProductName: "storagecontrol",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleStorageControlProjectIntelligenceFindingRevisions(),
	}.Register()
}
