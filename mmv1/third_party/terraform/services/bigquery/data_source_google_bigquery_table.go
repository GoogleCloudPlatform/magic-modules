package bigquery

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBigQueryTable() *schema.Resource {
	fieldSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The field name. The name must contain only letters (a-z, A-Z), numbers (0-9), or underscores (_), and must start with a letter or underscore. The maximum length is 300 characters",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The field data type.",
		},
		"mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "The field mode (NULLABLE, REQUIRED, or REPEATED).",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Field description. The maximum length is 1,024 characters.",
		},
		"fields": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Describes the nested schema fields if the type property is set to RECORD.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{},
			},
		},
		"policy_tags": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Policy tag list for this field.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"names": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
		"max_length": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Maximum length of values of this field for STRINGS or BYTES.",
		},
		"precision": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Precision (maximum number of total digits) for NUMERIC or BIGNUMERIC.",
		},
		"scale": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Scale (maximum number of digits in the fractional part) for NUMERIC or BIGNUMERIC.",
		},
		"rounding_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Rounding mode for NUMERIC or BIGNUMERIC.",
		},
		"collation": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Collation specification of the field.",
		},
		"default_value_expression": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Default value expression for this field.",
		},
		"range_element_type": {
			Type:        schema.TypeList,
			Computed:    true,
			Optional:    true,
			Description: "Element type for RANGE type fields.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	// fieldSchema["fields"].Elem.(*schema.Resource).Schema = fieldSchema

	dsSchema := map[string]*schema.Schema{
		"dataset_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the dataset containing the table.",
		},
		"table_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the table.",
		},
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the project in which the table is located. If it is not provided, the provider project is used.",
		},
		"table": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"schema": {
						Type:     schema.TypeList,
						Computed: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"fields": {
									Type:     schema.TypeList,
									Computed: true,
									Elem: &schema.Resource{
										Schema: fieldSchema,
									},
								},
							},
						},
					},
					// TODO(ramon) add other properties
					"table_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	return &schema.Resource{
		ReadContext: dataSourceGoogleBigQueryTableRead,
		Schema:      dsSchema,
	}
}

func dataSourceGoogleBigQueryTableRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return diag.FromErr(err)
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}")
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	log.Printf("[RAMON][DEBUG] BigQuery response: %s", res)
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error retrieving table: %s", err))
	}

	return nil
}
