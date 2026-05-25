package storagecontrol

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageControlProjectIntelligenceFindingsSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageControlProjectIntelligenceFindingsSummaryRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The filter expression to apply.`,
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				Description: `The location of the intelligence findings summary. Currently default value is global and users cannot use for input for now.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
			"resource_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PARENT",
				ValidateFunc: validation.StringInSlice([]string{"PARENT", "PROJECT"}, false),
				Description:  `The scope of the resources to include in the summary. Possible values are PARENT and PROJECT. Default value is PARENT.`,
			},
			"finding_summaries": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `A list of summaries for individual finding types.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The finding type.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The category of the finding.`,
						},
						"target_resource": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The target resource of the finding summary.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the finding summary.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last update time of the finding summary.`,
						},
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The severity of the finding.`,
						},
						"summary_details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Detailed summaries for the finding type.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The total count of findings for this summary slice.`,
									},
									"percentage": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The percentage magnitude associated with this summary slice.`,
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource type associated with the summary slice.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `A description explaining the summary slice.`,
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

func dataSourceGoogleStorageControlProjectIntelligenceFindingsSummaryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for intelligence findings summary: %s", err)
	}
	location := d.Get("location").(string)

	params := make(map[string]string)
	if v, ok := d.GetOk("filter"); ok {
		params["filter"] = v.(string)
	}
	if v, ok := d.GetOk("resource_scope"); ok {
		params["resourceScope"] = v.(string)
	}

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf(transport_tpg.BaseUrl(Product, config)+"projects/%s/locations/%s/intelligenceFindings:summarize", project, location))
	if err != nil {
		return fmt.Errorf("Error formatting url for project intelligence findings summary: %s", err)
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
		return fmt.Errorf("Error retrieving project intelligence findings summary: %s", err)
	}

	if err := d.Set("finding_summaries", flattenStorageControlFindingSummaries(res["findingSummaries"])); err != nil {
		return fmt.Errorf("Error setting finding_summaries: %s", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/locations/%s/intelligenceFindingsSummary", project, location))

	return nil
}

func flattenStorageControlFindingSummaries(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ls, ok := v.([]interface{})
	if !ok {
		return nil
	}
	summaries := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		o, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		summary := map[string]interface{}{
			"type":            o["type"],
			"category":        o["category"],
			"target_resource": o["targetResource"],
			"create_time":     o["createTime"],
			"update_time":     o["updateTime"],
			"severity":        o["severity"],
			"summary_details": flattenStorageControlSummaryDetails(o["summaryDetails"]),
		}
		summaries = append(summaries, summary)
	}
	return summaries
}

func flattenStorageControlSummaryDetails(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	ls, ok := v.([]interface{})
	if !ok {
		return nil
	}
	details := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		o, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		detail := map[string]interface{}{
			"count":         o["count"],
			"percentage":    flattenStorageControlFloat(o["percentage"]),
			"resource_type": o["resourceType"],
			"description":   o["description"],
		}
		details = append(details, detail)
	}
	return details
}

func init() {
	registry.Schema{
		Name:        "google_storage_control_project_intelligence_findings_summary",
		ProductName: "storagecontrol",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleStorageControlProjectIntelligenceFindingsSummary(),
	}.Register()
}
