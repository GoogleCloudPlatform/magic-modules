// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudquotas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudQuotasQuotaInfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudQuotasQuotaInfosRead,

		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"page_token": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"page_size": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "100",
			},
			"next_page_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"quota_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"quota_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_precise": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"refresh_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"metric_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_increase_eligibility": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_eligible": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"ineligibility_reason": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"is_fixed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"dimensions_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimensions": {
										Type:     schema.TypeMap,
										Computed: true,
									},
									"details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"applicable_locations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"is_concurrent": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"service_request_quota_uri": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func dataSourceGoogleCloudQuotasQuotaInfosRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{CloudQuotasBasePath}}{{parent}}/locations/global/services/{{service}}/quotaInfos?pageSize={{page_size}}&pageToken={{page_token}}")
	if err != nil {
		return fmt.Errorf("error setting api endpoint")
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("CloudQuotasQuotaInfo %q", d.Id()))
	}

	if err := d.Set("next_page_token", res["nextPageToken"]); err != nil {
		return fmt.Errorf("error reading next page token: %s", err)
	}
	if err := d.Set("quota_infos", flattenCloudQuotasQuotaInfos(res["quotaInfos"].([]interface{}), d, config)); err != nil {
		return fmt.Errorf("error reading quota infos : %s", err)
	}

	d.SetId(url)
	return nil
}

func flattenCloudQuotasQuotaInfos(fetchedQuotaInfos []interface{}, d *schema.ResourceData, config *transport_tpg.Config) []interface{} {
	if fetchedQuotaInfos == nil {
		return make([]interface{}, 0)
	}

	quotaInfos := make([]interface{}, 0, len(fetchedQuotaInfos))
	for _, rawQuotaInfo := range fetchedQuotaInfos {
		quotaInfo := make(map[string]interface{})
		rawQuotaInfo := rawQuotaInfo.(map[string]interface{})

		quotaInfo["name"] = rawQuotaInfo["name"]
		quotaInfo["quota_id"] = rawQuotaInfo["quotaId"]
		quotaInfo["metric"] = rawQuotaInfo["metric"]
		quotaInfo["service"] = rawQuotaInfo["service"]
		quotaInfo["is_precise"] = rawQuotaInfo["isPrecise"]
		quotaInfo["refresh_interval"] = rawQuotaInfo["refreshInterval"]
		quotaInfo["container_type"] = rawQuotaInfo["containerType"]
		quotaInfo["dimensions"] = rawQuotaInfo["dimensions"]
		quotaInfo["metric_display_name"] = rawQuotaInfo["metricDisplayName"]
		quotaInfo["quota_display_name"] = rawQuotaInfo["quotaDisplayName"]
		quotaInfo["metric_unit"] = rawQuotaInfo["metricUnit"]
		quotaInfo["quota_increase_eligibility"] = flattenCloudQuotasQuotaInfoQuotaIncreaseEligibility(rawQuotaInfo["quotaIncreaseEligibility"], d, config)
		quotaInfo["is_fixed"] = rawQuotaInfo["isFixed"]
		quotaInfo["dimensions_infos"] = flattenCloudQuotasQuotaInfoDimensionsInfos(rawQuotaInfo["dimensionsInfos"], d, config)
		quotaInfo["is_concurrent"] = rawQuotaInfo["isConcurrent"]
		quotaInfo["service_request_quota_uri"] = rawQuotaInfo["serviceRequestQuotaUri"]

		quotaInfos = append(quotaInfos, quotaInfo)
	}
	return quotaInfos
}
