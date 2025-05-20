package bigquery

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// see https://cloud.google.com/bigquery/docs/nested-repeated#limitations
const maxNestingLevel = 15

func DataSourceGoogleBigQueryTable() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceBigQueryTableRead,
		Schema: getDataSourceSchema(),
	}
}

func getDataSourceSchema() map[string]*schema.Schema {
	fieldSchema := buildFieldSchema(maxNestingLevel)

	return map[string]*schema.Schema{
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
			Description: "The ID of the project in which the table is located.",
		},
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
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func buildFieldSchema(depth int) map[string]*schema.Schema {
	baseSchema := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The field name.",
		},
		"type": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The field data type.",
		},
		"mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The field mode (NULLABLE, REQUIRED, or REPEATED).",
		},
		"description": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Field description.",
		},
		"policy_tags": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Policy tag list for this field.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"names": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
		"max_length": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Maximum length of values of this field for STRINGS or BYTES.",
		},
		"precision": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Precision (maximum number of total digits) for NUMERIC or BIGNUMERIC.",
		},
		"scale": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Scale (maximum number of digits in the fractional part) for NUMERIC or BIGNUMERIC.",
		},
		"rounding_mode": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Rounding mode for NUMERIC or BIGNUMERIC.",
		},
		"collation": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Collation specification of the field.",
		},
		"default_value_expression": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Default value expression for this field.",
		},
		"range_element_type": {
			Type:        schema.TypeList,
			Computed:    true,
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

	if depth > 0 {
		nestedSchema := make(map[string]*schema.Schema)
		for k, v := range baseSchema {
			nestedSchema[k] = v
		}
		nestedSchema["fields"] = &schema.Schema{
			Type:        schema.TypeList,
			Computed:    true,
			Description: "Nested fields for RECORD type fields.",
			Elem: &schema.Resource{
				Schema: buildFieldSchema(depth - 1),
			},
		}
		return nestedSchema
	}

	return baseSchema
}

func dataSourceBigQueryTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}

	datasetID := d.Get("dataset_id").(string)
	tableID := d.Get("table_id").(string)

	url, err := tpgresource.ReplaceVars(d, config, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}")
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
		return fmt.Errorf("Error retrieving table: %s", err)
	}

	if schemaData, ok := res["schema"].(map[string]interface{}); ok {
		if fields, ok := schemaData["fields"].([]interface{}); ok {
			if err := d.Set("schema", []map[string]interface{}{
				{"fields": flattenSchemaFields(fields, 0)},
			}); err != nil {
				return fmt.Errorf("Error setting schema: %s", err)
			}
		}
	}

	d.SetId(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", project, datasetID, tableID))
	return nil
}

func flattenSchemaFields(fields []interface{}, currentLevel int) []interface{} {
	if currentLevel > maxNestingLevel {
		return nil
	}

	var result []interface{}
	for _, f := range fields {
		field := f.(map[string]interface{})
		flattened := map[string]interface{}{
			"name":                     field["name"],
			"type":                     field["type"],
			"mode":                     field["mode"],
			"description":              field["description"],
			"max_length":               field["maxLength"],
			"precision":                field["precision"],
			"scale":                    field["scale"],
			"rounding_mode":            field["roundingMode"],
			"collation":                field["collation"],
			"default_value_expression": field["defaultValueExpression"],
		}

		if policyTags, ok := field["policyTags"].(map[string]interface{}); ok {
			if names, ok := policyTags["names"].([]interface{}); ok {
				flattened["policy_tags"] = []interface{}{
					map[string]interface{}{"names": names},
				}
			}
		}

		if rangeElementType, ok := field["rangeElementType"].(map[string]interface{}); ok {
			flattened["range_element_type"] = []interface{}{
				map[string]interface{}{"type": rangeElementType["type"]},
			}
		}

		if field["type"] == "RECORD" {
			if nestedFields, ok := field["fields"].([]interface{}); ok {
				// a RECORD has nested fields, therefore recursion is applied here
				flattened["fields"] = flattenSchemaFields(nestedFields, currentLevel+1)
			}
		}

		result = append(result, flattened)
	}
	return result
}
