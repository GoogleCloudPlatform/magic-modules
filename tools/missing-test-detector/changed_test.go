package main

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceMapChanges(t *testing.T) {
	for _, test := range []struct {
		name                  string
		oldResourceMap        map[string]*schema.Resource
		newResourceMap        map[string]*schema.Resource
		expectedChangedFields map[string]FieldCoverage
	}{
		{
			name:                  "empty-maps",
			oldResourceMap:        map[string]*schema.Resource{},
			newResourceMap:        map[string]*schema.Resource{},
			expectedChangedFields: map[string]FieldCoverage{},
		},
		{
			name:           "empty-resources",
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {},
				"google_service_one_resource_two": {},
			},
			expectedChangedFields: map[string]FieldCoverage{},
		},
		{
			name: "new-nested-field",
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
				"google_service_one_resource_two": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
				"google_service_one_resource_two": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
									"field_four": {
										Type: schema.TypeInt,
									},
								},
							},
						},
					},
				},
			},
			expectedChangedFields: map[string]FieldCoverage{
				"google_service_one_resource_two": {"field_two": FieldCoverage{"field_four": false}},
			},
		},
		{
			name: "new-field-in-two-resources",
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
				"google_service_one_resource_two": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
								},
							},
						},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
									"field_four": {
										Type: schema.TypeInt,
									},
								},
							},
						},
					},
				},
				"google_service_one_resource_two": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
						"field_two": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field_three": {
										Type: schema.TypeString,
									},
									"field_four": {
										Type: schema.TypeInt,
									},
								},
							},
						},
					},
				},
			},
			expectedChangedFields: map[string]FieldCoverage{
				"google_service_one_resource_one": {"field_two": FieldCoverage{"field_four": false}},
				"google_service_one_resource_two": {"field_two": FieldCoverage{"field_four": false}},
			},
		},
		{
			name: "deleted-field",
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{},
				},
			},
			expectedChangedFields: map[string]FieldCoverage{},
		},
		{
			name: "deleted-resource",
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
					},
				},
			},
			expectedChangedFields: map[string]FieldCoverage{},
		},
		{
			name: "new-resource",
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
					},
				},
			},
			expectedChangedFields: map[string]FieldCoverage{"google_service_one_resource_one": {"field_one": false}},
		},
	} {
		changedFields := resourceMapChanges(test.oldResourceMap, test.newResourceMap)
		if !reflect.DeepEqual(changedFields, test.expectedChangedFields) {
			t.Errorf("%s test failed: unexpected changed resources: %v, expected %v", test.name, changedFields, test.expectedChangedFields)
		}
	}
}
