package tpgresource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func EffectiveLabelsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeMap,
		Computed:    true,
		Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}
}

func FlattenLabels(labels map[string]string, d *schema.ResourceData) map[string]interface{} {
	transformed := make(map[string]interface{})

	if v, ok := d.GetOk("labels"); ok {
		if labels != nil {
			for k, _ := range v.(map[string]interface{}) {
				transformed[k] = labels[k]
			}
		}
	}

	return transformed
}
