package storagecontrol

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageControlOrganizationIntelligenceFindingsSummary() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleStorageControlOrganizationIntelligenceFindingsSummaryRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The filter expression. Supports filtering by FindingType.`,
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "global",
				Description: `The location of the intelligence findings summary. Currently default value is global and users cannot use for input for now.`,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Google Cloud Organization.`,
			},
			"resource_scope": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PARENT",
				ValidateFunc: validation.StringInSlice([]string{"PARENT", "PROJECT"}, false),
				Description:  `Determines the granularity of the findings when the parent is an organization or folder. Possible values are PARENT and PROJECT. Default value is PARENT.`,
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
							Description: `The type of finding.`,
						},
						"category": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The category of the finding.`,
						},
						"target_resource": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The fully qualified Cloud resource name for which this summary was generated.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the earliest finding that this summary is based on.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time of the most recent update among all the findings that this summary is based on.`,
						},
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Severity of the finding.`,
						},
						"summary_details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The SummaryDetails resources.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The count of impacted resources.`,
									},
									"percentage": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The percentage of impacted resources.`,
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of Cloud resource this summary detail applies to.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `A short description about the FindingSummary.`,
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

func dataSourceGoogleStorageControlOrganizationIntelligenceFindingsSummaryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	organization := d.Get("organization").(string)
	location := d.Get("location").(string)

	params := make(map[string]string)
	if v, ok := d.GetOk("filter"); ok {
		params["filter"] = v.(string)
	}
	if v, ok := d.GetOk("resource_scope"); ok {
		params["resourceScope"] = v.(string)
	}

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf(transport_tpg.BaseUrl(Product, config)+"organizations/%s/locations/%s/intelligenceFindings:summarize", organization, location))
	if err != nil {
		return fmt.Errorf("Error formatting url for organization intelligence findings summary: %s", err)
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
		return fmt.Errorf("Error retrieving organization intelligence findings summary: %s", err)
	}

	if err := d.Set("finding_summaries", flattenStorageControlFindingSummaries(res["findingSummaries"])); err != nil {
		return fmt.Errorf("Error setting finding_summaries: %s", err)
	}

	d.SetId(fmt.Sprintf("organizations/%s/locations/%s/intelligenceFindingsSummary", organization, location))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_storage_control_organization_intelligence_findings_summary",
		ProductName: "storagecontrol",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleStorageControlOrganizationIntelligenceFindingsSummary(),
	}.Register()
}
