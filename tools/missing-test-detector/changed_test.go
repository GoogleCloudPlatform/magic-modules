package main

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceMapChanges(t *testing.T) {
	for _, test := range []struct {
		oldResourceMap        map[string]*schema.Resource
		newResourceMap        map[string]*schema.Resource
		expectedChangedFields map[string][]string
	}{
		{
			oldResourceMap:        map[string]*schema.Resource{},
			newResourceMap:        map[string]*schema.Resource{},
			expectedChangedFields: map[string][]string{},
		},
		{
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {},
				"google_service_one_resource_two": {},
			},
			expectedChangedFields: map[string][]string{
				"google_service_one_resource_one": {},
				"google_service_one_resource_two": {},
			},
		},
		{
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
			expectedChangedFields: map[string][]string{
				"google_service_one_resource_two": {"field_two.field_four"},
			},
		},
		{
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
			expectedChangedFields: map[string][]string{
				"google_service_one_resource_one": {"field_two.field_four"},
				"google_service_one_resource_two": {"field_two.field_four"},
			},
		},
	} {
		changedFields := resourceMapChanges(test.oldResourceMap, test.newResourceMap)
		if !reflect.DeepEqual(changedFields, test.expectedChangedFields) {
			t.Errorf("unexpected changed resources: %v, expected %v", changedFields, test.expectedChangedFields)
		}
	}
}
