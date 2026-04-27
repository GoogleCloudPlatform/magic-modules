package storage

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func storageCustomHeaders() *schema.Schema {
	element := &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: `User-provided custom headers, in key/value pairs.`,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}

	return element
}

func expandCustomHeaders(v interface{}) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
