package main

import (
	"encoding/json"
	"reflect"
	"testing"

	newProvider "google/provider/new/google-beta/provider"
	oldProvider "google/provider/old/google-beta/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestNewProviderOldProviderChanges(t *testing.T) {
	changes := resourceMapChanges(oldProvider.ResourceMap(), newProvider.ResourceMap())

	jsonChanges, err := json.MarshalIndent(changes, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling resource map changes to json: %s", err)
	}

	t.Logf("Changes between old and new providers: %s", jsonChanges)
}

func TestResourceMapChanges(t *testing.T) {
	for _, test := range []struct {
		name                  string
		oldResourceMap        map[string]*schema.Resource
		newResourceMap        map[string]*schema.Resource
		expectedChangedFields map[string]ResourceChanges
	}{
		{
			name:                  "empty-maps",
			oldResourceMap:        map[string]*schema.Resource{},
			newResourceMap:        map[string]*schema.Resource{},
			expectedChangedFields: map[string]ResourceChanges{},
		},
		{
			name:           "empty-resources",
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {},
				"google_service_one_resource_two": {},
			},
			expectedChangedFields: map[string]ResourceChanges{},
		},
		{
			name: "unchanged-nested-field",
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
			},
			expectedChangedFields: map[string]ResourceChanges{},
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
			expectedChangedFields: map[string]ResourceChanges{
				"google_service_one_resource_two": {
					"field_two": ResourceChanges{
						"field_four": &Field{Added: true},
					},
				},
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
			expectedChangedFields: map[string]ResourceChanges{
				"google_service_one_resource_one": {
					"field_two": ResourceChanges{
						"field_four": &Field{Added: true},
					},
				},
				"google_service_one_resource_two": {
					"field_two": ResourceChanges{
						"field_four": &Field{Added: true},
					},
				},
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
			expectedChangedFields: map[string]ResourceChanges{},
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
			expectedChangedFields: map[string]ResourceChanges{},
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
			expectedChangedFields: map[string]ResourceChanges{
				"google_service_one_resource_one": {
					"field_one": &Field{Added: true},
				},
			},
		},
	} {
		changedFields := resourceMapChanges(test.oldResourceMap, test.newResourceMap)
		if !reflect.DeepEqual(changedFields, test.expectedChangedFields) {
			t.Errorf("%s test failed: unexpected changed resources: %v, expected %v", test.name, changedFields, test.expectedChangedFields)
		}
	}
}
