package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCheckDiffSuppressFunc(t *testing.T) {
	// Define a dummy schema with DiffSuppressFunc
	resourceSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"simple_string": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new == "suppress_me"
				},
			},
			"nested_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inner_string": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return new == "suppress_nested"
							},
						},
					},
				},
			},
			"map_with_suppress": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress if key is "ignored_key"
					// k will be "map_with_suppress.ignored_key"
					return k == "map_with_suppress.ignored_key"
				},
			},
			"nested_schema_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return new == "suppress_primitive"
					},
				},
			},
		},
	}

	provider := &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"dummy_resource": resourceSchema,
		},
	}

	tests := []struct {
		name         string
		key          string
		val          any
		resourceType string
		want         bool
	}{
		{
			name:         "Simple string suppressed",
			key:          "simple_string",
			val:          "suppress_me",
			resourceType: "dummy_resource",
			want:         true,
		},
		{
			name:         "Simple string not suppressed",
			key:          "simple_string",
			val:          "keep_me",
			resourceType: "dummy_resource",
			want:         false,
		},
		{
			name:         "Nested list string suppressed",
			key:          "nested_list.0.inner_string",
			val:          "suppress_nested",
			resourceType: "dummy_resource",
			want:         true,
		},
		{
			name:         "Map key suppressed by parent map",
			key:          "map_with_suppress.ignored_key",
			val:          "any_val",
			resourceType: "dummy_resource",
			want:         true,
		},
		{
			name:         "Map key not suppressed",
			key:          "map_with_suppress.other_key",
			val:          "any_val",
			resourceType: "dummy_resource",
			want:         false,
		},
		{
			name:         "Primitive list suppressed",
			key:          "nested_schema_list.0",
			val:          "suppress_primitive",
			resourceType: "dummy_resource",
			want:         true,
		},
		{
			name:         "Unknown resource",
			key:          "simple_string",
			val:          "suppress_me",
			resourceType: "unknown_resource",
			want:         false,
		},
		{
			name:         "Unknown field",
			key:          "unknown_field",
			val:          "val",
			resourceType: "dummy_resource",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := checkDiffSuppressFunc(tt.key, tt.val, tt.resourceType, provider)
			if got != tt.want {
				t.Errorf("checkDiffSuppressFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
