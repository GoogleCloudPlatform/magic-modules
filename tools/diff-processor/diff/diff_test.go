package diff

import (
	"strings"
	"testing"

	newProvider "google/provider/new/google/provider"

	newVerify "google/provider/new/google/verify"
	oldProvider "google/provider/old/google/provider"
	oldVerify "google/provider/old/google/verify"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestNewProviderOldProviderChanges(t *testing.T) {
	changes := ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())

	for resource, resourceDiff := range changes {
		if resourceDiff.ResourceConfig.Old == nil {
			t.Logf("%s is added", resource)
			continue
		}
		if resourceDiff.ResourceConfig.New == nil {
			t.Logf("%s is removed", resource)
			continue
		}
		t.Logf("%s is modified", resource)
		if diff := cmp.Diff(resourceDiff.ResourceConfig.Old, resourceDiff.ResourceConfig.New); diff != "" {
			t.Logf("%s config changes (-old, +new):\n%s", resource, diff)
		}
		for field, fieldDiff := range resourceDiff.Fields {
			if fieldDiff.Old == nil {
				t.Logf("%s.%s is added", resource, field)
				continue
			}
			if fieldDiff.New == nil {
				t.Logf("%s.%s is removed", resource, field)
				continue
			}
			t.Logf("%s.%s is modified", resource, field)
			if diff := cmp.Diff(fieldDiff.Old, fieldDiff.New); diff != "" {
				t.Logf("%s.%s changes (-old, +new):\n%s", resource, field, diff)
			}
		}

	}
}

func TestFlattenSchema(t *testing.T) {
	cases := map[string]struct {
		resourceSchema  map[string]*schema.Schema
		expectFlattened map[string]*schema.Schema
	}{
		"primitive fields": {
			resourceSchema: map[string]*schema.Schema{
				"bool": {
					Type: schema.TypeBool,
				},
				"int": {
					Type: schema.TypeInt,
				},
				"float": {
					Type: schema.TypeFloat,
				},
				"string": {
					Type: schema.TypeString,
				},
			},
			expectFlattened: map[string]*schema.Schema{
				"bool": {
					Type: schema.TypeBool,
				},
				"int": {
					Type: schema.TypeInt,
				},
				"float": {
					Type: schema.TypeFloat,
				},
				"string": {
					Type: schema.TypeString,
				},
			},
		},
		"map field": {
			resourceSchema: map[string]*schema.Schema{
				"map": {
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			expectFlattened: map[string]*schema.Schema{
				"map": {
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		"simple list field": {
			resourceSchema: map[string]*schema.Schema{
				"list": {
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			expectFlattened: map[string]*schema.Schema{
				"list": {
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		"simple set field": {
			resourceSchema: map[string]*schema.Schema{
				"set": {
					Type: schema.TypeSet,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			expectFlattened: map[string]*schema.Schema{
				"set": {
					Type: schema.TypeSet,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		"nested list field": {
			resourceSchema: map[string]*schema.Schema{
				"list": {
					Type: schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"nested_string": {
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			expectFlattened: map[string]*schema.Schema{
				"list": {
					Type: schema.TypeList,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"nested_string": {
								Type: schema.TypeString,
							},
						},
					},
				},
				"list.nested_string": {
					Type: schema.TypeString,
				},
			},
		},
		"nested set field": {
			resourceSchema: map[string]*schema.Schema{
				"set": {
					Type: schema.TypeSet,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"nested_string": {
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			expectFlattened: map[string]*schema.Schema{
				"set": {
					Type: schema.TypeSet,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"nested_string": {
								Type: schema.TypeString,
							},
						},
					},
				},
				"set.nested_string": {
					Type: schema.TypeString,
				},
			},
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			flattened := flattenSchema("", tc.resourceSchema)
			assert.Equal(t, tc.expectFlattened, flattened)
		})
	}
}

func testDefaultFunc1() (interface{}, error) {
	return "default1", nil
}
func testDefaultFunc2() (interface{}, error) {
	return "default2", nil
}
func testStateFunc1(interface{}) string {
	return "state1"
}
func testStateFunc2(interface{}) string {
	return "state2"
}
func testValidateDiagFunc1(v interface{}, p cty.Path) diag.Diagnostics {
	return diag.Diagnostics{}
}
func testValidateDiagFunc2(v interface{}, p cty.Path) diag.Diagnostics {
	return diag.Diagnostics{}
}

func TestFieldChanged(t *testing.T) {
	cases := map[string]struct {
		oldField      *schema.Schema
		newField      *schema.Schema
		expectChanged bool
	}{
		"both nil": {
			oldField:      nil,
			newField:      nil,
			expectChanged: false,
		},
		"old nil": {
			oldField: nil,
			newField: &schema.Schema{
				Type: schema.TypeString,
			},
			expectChanged: true,
		},
		"new nil": {
			oldField: &schema.Schema{
				Type: schema.TypeString,
			},
			newField:      nil,
			expectChanged: true,
		},
		"not changed": {
			oldField:      &schema.Schema{},
			newField:      &schema.Schema{},
			expectChanged: false,
		},
		"Type changed": {
			oldField: &schema.Schema{
				Type: schema.TypeString,
			},
			newField: &schema.Schema{
				Type: schema.TypeInt,
			},
			expectChanged: true,
		},
		"ConfigMode changed": {
			oldField: &schema.Schema{
				ConfigMode: schema.SchemaConfigModeAttr,
			},
			newField: &schema.Schema{
				ConfigMode: schema.SchemaConfigModeBlock,
			},
			expectChanged: true,
		},
		"Required changed": {
			oldField: &schema.Schema{
				Required: false,
			},
			newField: &schema.Schema{
				Required: true,
			},
			expectChanged: true,
		},
		"Optional changed": {
			oldField: &schema.Schema{
				Optional: false,
			},
			newField: &schema.Schema{
				Optional: true,
			},
			expectChanged: true,
		},
		"Computed changed": {
			oldField: &schema.Schema{
				Computed: false,
			},
			newField: &schema.Schema{
				Computed: true,
			},
			expectChanged: true,
		},
		"ForceNew changed": {
			oldField: &schema.Schema{
				ForceNew: false,
			},
			newField: &schema.Schema{
				ForceNew: true,
			},
			expectChanged: true,
		},
		"DiffSuppressOnRefresh changed": {
			oldField: &schema.Schema{
				DiffSuppressOnRefresh: false,
			},
			newField: &schema.Schema{
				DiffSuppressOnRefresh: true,
			},
			expectChanged: true,
		},
		"Default changed": {
			oldField: &schema.Schema{
				Default: 10,
			},
			newField: &schema.Schema{
				Default: 20,
			},
			expectChanged: true,
		},
		"Description changed": {
			oldField: &schema.Schema{
				Description: "Hello",
			},
			newField: &schema.Schema{
				Description: "Goodbye",
			},
			expectChanged: true,
		},
		"InputDefault changed": {
			oldField: &schema.Schema{
				InputDefault: "Hello",
			},
			newField: &schema.Schema{
				InputDefault: "Goodbye",
			},
			expectChanged: true,
		},
		"MaxItems changed": {
			oldField: &schema.Schema{
				MaxItems: 10,
			},
			newField: &schema.Schema{
				MaxItems: 20,
			},
			expectChanged: true,
		},
		"MinItems changed": {
			oldField: &schema.Schema{
				MinItems: 10,
			},
			newField: &schema.Schema{
				MinItems: 20,
			},
			expectChanged: true,
		},
		"Deprecated changed": {
			oldField: &schema.Schema{
				Deprecated: "Hello",
			},
			newField: &schema.Schema{
				Deprecated: "Goodbye",
			},
			expectChanged: true,
		},
		"Sensitive changed": {
			oldField: &schema.Schema{
				Sensitive: false,
			},
			newField: &schema.Schema{
				Sensitive: true,
			},
			expectChanged: true,
		},
		"ConflictsWith reordered": {
			oldField: &schema.Schema{
				ConflictsWith: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				ConflictsWith: []string{"field_two", "field_one"},
			},
			expectChanged: false,
		},
		"ConflictsWith changed": {
			oldField: &schema.Schema{
				ConflictsWith: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				ConflictsWith: []string{"field_two", "field_three"},
			},
			expectChanged: true,
		},
		"ExactlyOneOf reordered": {
			oldField: &schema.Schema{
				ExactlyOneOf: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				ExactlyOneOf: []string{"field_two", "field_one"},
			},
			expectChanged: false,
		},
		"ExactlyOneOf changed": {
			oldField: &schema.Schema{
				ExactlyOneOf: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				ExactlyOneOf: []string{"field_two", "field_three"},
			},
			expectChanged: true,
		},
		"AtLeastOneOf reordered": {
			oldField: &schema.Schema{
				AtLeastOneOf: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				AtLeastOneOf: []string{"field_two", "field_one"},
			},
			expectChanged: false,
		},
		"AtLeastOneOf changed": {
			oldField: &schema.Schema{
				AtLeastOneOf: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				AtLeastOneOf: []string{"field_two", "field_three"},
			},
			expectChanged: true,
		},
		"RequiredWith reordered": {
			oldField: &schema.Schema{
				RequiredWith: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				RequiredWith: []string{"field_two", "field_one"},
			},
			expectChanged: false,
		},
		"RequiredWith changed": {
			oldField: &schema.Schema{
				RequiredWith: []string{"field_one", "field_two"},
			},
			newField: &schema.Schema{
				RequiredWith: []string{"field_two", "field_three"},
			},
			expectChanged: true,
		},
		"simple Elem unset -> set": {
			oldField: &schema.Schema{},
			newField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeInt},
			},
			expectChanged: true,
		},
		"simple Elem set -> unset": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeInt},
			},
			newField:      &schema.Schema{},
			expectChanged: true,
		},
		"simple Elem unchanged": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			expectChanged: false,
		},
		"simple Elem changed": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeInt},
			},
			expectChanged: true,
		},
		"nested Elem unset -> set": {
			oldField: &schema.Schema{},
			newField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			expectChanged: true,
		},
		"nested Elem set -> unset": {
			oldField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			newField:      &schema.Schema{},
			expectChanged: true,
		},
		"nested Elem unchanged": {
			oldField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			expectChanged: false,
		},
		"nested Elem changing is ignored": {
			oldField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"barbaz": {
							Type: schema.TypeString,
						},
					},
				},
			},
			expectChanged: false,
		},
		"Elem simple -> nested": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			newField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			expectChanged: true,
		},
		"Elem nested -> simple": {
			oldField: &schema.Schema{
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"foobar": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			expectChanged: true,
		},

		"DefaultFunc added": {
			oldField: &schema.Schema{},
			newField: &schema.Schema{
				DefaultFunc: testDefaultFunc1,
			},
			expectChanged: true,
		},
		"DefaultFunc removed": {
			oldField: &schema.Schema{
				DefaultFunc: testDefaultFunc1,
			},
			newField:      &schema.Schema{},
			expectChanged: true,
		},
		"DefaultFunc remains set": {
			oldField: &schema.Schema{
				DefaultFunc: testDefaultFunc1,
			},
			newField: &schema.Schema{
				DefaultFunc: testDefaultFunc2,
			},
			expectChanged: false,
		},
		"Elem DefaultFunc added": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					DefaultFunc: testDefaultFunc2,
				},
			},
			expectChanged: true,
		},
		"Elem DefaultFunc removed": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					DefaultFunc: testDefaultFunc1,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			expectChanged: true,
		},
		"Elem DefaultFunc remains set": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					DefaultFunc: testDefaultFunc1,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					DefaultFunc: testDefaultFunc2,
				},
			},
			expectChanged: false,
		},

		"StateFunc added": {
			oldField: &schema.Schema{},
			newField: &schema.Schema{
				StateFunc: testStateFunc1,
			},
			expectChanged: true,
		},
		"StateFunc removed": {
			oldField: &schema.Schema{
				StateFunc: testStateFunc1,
			},
			newField:      &schema.Schema{},
			expectChanged: true,
		},
		"StateFunc remains set": {
			oldField: &schema.Schema{
				StateFunc: testStateFunc1,
			},
			newField: &schema.Schema{
				StateFunc: testStateFunc2,
			},
			expectChanged: false,
		},
		"Elem StateFunc added": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					StateFunc: testStateFunc2,
				},
			},
			expectChanged: true,
		},
		"Elem StateFunc removed": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					StateFunc: testStateFunc1,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			expectChanged: true,
		},
		"Elem StateFunc remains set": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					StateFunc: testStateFunc1,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:      schema.TypeString,
					StateFunc: testStateFunc2,
				},
			},
			expectChanged: false,
		},

		"ValidateFunc added": {
			oldField: &schema.Schema{},
			newField: &schema.Schema{
				ValidateFunc: newVerify.ValidateBase64String,
			},
			expectChanged: true,
		},
		"ValidateFunc removed": {
			oldField: &schema.Schema{
				ValidateFunc: oldVerify.ValidateBase64String,
			},
			newField:      &schema.Schema{},
			expectChanged: true,
		},
		"ValidateFunc remains set": {
			oldField: &schema.Schema{
				ValidateFunc: oldVerify.ValidateBase64String,
			},
			newField: &schema.Schema{
				ValidateFunc: newVerify.ValidateBase64String,
			},
			expectChanged: false,
		},
		"Elem ValidateFunc added": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: newVerify.ValidateBase64String,
				},
			},
			expectChanged: true,
		},
		"Elem ValidateFunc removed": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: oldVerify.ValidateBase64String,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			expectChanged: true,
		},
		"Elem ValidateFunc remains set": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: oldVerify.ValidateBase64String,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: newVerify.ValidateBase64String,
				},
			},
			expectChanged: false,
		},

		"ValidateDiagFunc added": {
			oldField: &schema.Schema{},
			newField: &schema.Schema{
				ValidateDiagFunc: testValidateDiagFunc1,
			},
			expectChanged: true,
		},
		"ValidateDiagFunc removed": {
			oldField: &schema.Schema{
				ValidateDiagFunc: testValidateDiagFunc1,
			},
			newField:      &schema.Schema{},
			expectChanged: true,
		},
		"ValidateDiagFunc remains set": {
			oldField: &schema.Schema{
				ValidateDiagFunc: testValidateDiagFunc1,
			},
			newField: &schema.Schema{
				ValidateDiagFunc: testValidateDiagFunc2,
			},
			expectChanged: false,
		},
		"Elem ValidateDiagFunc added": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: testValidateDiagFunc2,
				},
			},
			expectChanged: true,
		},
		"Elem ValidateDiagFunc removed": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: testValidateDiagFunc1,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			expectChanged: true,
		},
		"Elem ValidateDiagFunc remains set": {
			oldField: &schema.Schema{
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: testValidateDiagFunc1,
				},
			},
			newField: &schema.Schema{
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: testValidateDiagFunc2,
				},
			},
			expectChanged: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			_, _, changed := diffFields(tc.oldField, tc.newField, "")
			if changed != tc.expectChanged {
				if diff := cmp.Diff(tc.oldField, tc.newField); diff != "" {
					t.Errorf("want %t; got %t.\nField diff (-old, +new):\n%s",
						tc.expectChanged,
						changed,
						diff,
					)
				} else {
					t.Errorf("want %t; got %t. No field diff.\nOld field: %s\nNew field: %s\n",
						tc.expectChanged,
						changed,
						spew.Sdump(tc.oldField),
						spew.Sdump(tc.newField),
					)
				}
			}
		})
	}
}

func TestComputeSchemaDiff(t *testing.T) {
	cases := map[string]struct {
		oldResourceMap     map[string]*schema.Resource
		newResourceMap     map[string]*schema.Resource
		expectedSchemaDiff SchemaDiff
	}{
		"empty-maps": {
			oldResourceMap:     map[string]*schema.Resource{},
			newResourceMap:     map[string]*schema.Resource{},
			expectedSchemaDiff: SchemaDiff{},
		},
		"empty-resources": {
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {},
				"google_service_one_resource_two": {},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {},
				"google_service_one_resource_two": {},
			},
			expectedSchemaDiff: SchemaDiff{},
		},
		"unchanged-nested-field": {
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
			expectedSchemaDiff: SchemaDiff{},
		},
		"new-nested-field": {
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
							ConflictsWith: []string{
								"field_two",
							},
						},
						"field_two": {
							Type: schema.TypeList,
							ConflictsWith: []string{
								"field_one",
							},
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
							ConflictsWith: []string{
								"field_two",
							},
						},
						"field_two": {
							Type: schema.TypeList,
							ConflictsWith: []string{
								"field_one",
							},
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
										ConflictsWith: []string{
											"field_two.0.field_four",
										},
									},
									"field_four": {
										Type: schema.TypeInt,
										ConflictsWith: []string{
											"field_two.0.field_three",
										},
									},
								},
							},
						},
					},
				},
			},
			expectedSchemaDiff: SchemaDiff{
				"google_service_one_resource_two": ResourceDiff{
					ResourceConfig: ResourceConfigDiff{
						Old: &schema.Resource{},
						New: &schema.Resource{},
					},
					FlattenedSchema: FlattenedSchemaRaw{
						Old: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
							"field_two": {
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"field_three": {Type: schema.TypeString},
									},
								},
							},
							"field_two.field_three": {Type: schema.TypeString},
						},
						New: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
							"field_two": {
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"field_three": {
											Type:          schema.TypeString,
											ConflictsWith: []string{"field_two.0.field_four"},
										},
										"field_four": {
											Type:          schema.TypeInt,
											ConflictsWith: []string{"field_two.0.field_three"},
										},
									},
								},
							},
							"field_two.field_three": {
								Type:          schema.TypeString,
								ConflictsWith: []string{"field_two.0.field_four"},
							},
							"field_two.field_four": {
								Type:          schema.TypeInt,
								ConflictsWith: []string{"field_two.0.field_three"},
							},
						},
					},
					Fields: map[string]FieldDiff{
						"field_two.field_three": FieldDiff{
							Old: &schema.Schema{
								Type: schema.TypeString,
							},
							New: &schema.Schema{
								Type: schema.TypeString,
								ConflictsWith: []string{
									"field_two.0.field_four",
								},
							},
						},
						"field_two.field_four": FieldDiff{
							Old: nil,
							New: &schema.Schema{
								Type: schema.TypeInt,
								ConflictsWith: []string{
									"field_two.0.field_three",
								},
							},
						},
					},
					FieldSets: ResourceFieldSetsDiff{
						Old: ResourceFieldSets{},
						New: ResourceFieldSets{
							ConflictsWith: map[string]FieldSet{
								"field_two.field_four,field_two.field_three": FieldSet{
									"field_two.field_three": {},
									"field_two.field_four":  {},
								},
							},
						},
					},
				},
			},
		},
		"new-field-in-two-resources": {
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
			expectedSchemaDiff: SchemaDiff{
				"google_service_one_resource_one": ResourceDiff{
					ResourceConfig: ResourceConfigDiff{
						Old: &schema.Resource{},
						New: &schema.Resource{},
					},
					FlattenedSchema: FlattenedSchemaRaw{
						Old: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
							"field_two": {
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"field_three": {Type: schema.TypeString},
									},
								},
							},
							"field_two.field_three": {Type: schema.TypeString},
						},
						New: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
							"field_two": {
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"field_three": {Type: schema.TypeString},
										"field_four":  {Type: schema.TypeInt},
									},
								},
							},
							"field_two.field_three": {Type: schema.TypeString},
							"field_two.field_four":  {Type: schema.TypeInt},
						},
					},
					Fields: map[string]FieldDiff{
						"field_two.field_four": FieldDiff{
							Old: nil,
							New: &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
				"google_service_one_resource_two": ResourceDiff{
					ResourceConfig: ResourceConfigDiff{
						Old: &schema.Resource{},
						New: &schema.Resource{},
					},
					FlattenedSchema: FlattenedSchemaRaw{
						Old: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
							"field_two": {
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"field_three": {Type: schema.TypeString},
									},
								},
							},
							"field_two.field_three": {Type: schema.TypeString},
						},
						New: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
							"field_two": {
								Type: schema.TypeList,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"field_three": {Type: schema.TypeString},
										"field_four":  {Type: schema.TypeInt},
									},
								},
							},
							"field_two.field_three": {Type: schema.TypeString},
							"field_two.field_four":  {Type: schema.TypeInt},
						},
					},
					Fields: map[string]FieldDiff{
						"field_two.field_four": FieldDiff{
							Old: nil,
							New: &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
		},
		"deleted-field": {
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
			expectedSchemaDiff: SchemaDiff{
				"google_service_one_resource_one": ResourceDiff{
					ResourceConfig: ResourceConfigDiff{
						Old: &schema.Resource{},
						New: &schema.Resource{},
					},
					FlattenedSchema: FlattenedSchemaRaw{
						Old: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
						},
						New: map[string]*schema.Schema{},
					},
					Fields: map[string]FieldDiff{
						"field_one": FieldDiff{
							Old: &schema.Schema{Type: schema.TypeString},
							New: nil,
						},
					},
				},
			},
		},
		"deleted-resource": {
			oldResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
					},
				},
			},
			expectedSchemaDiff: SchemaDiff{
				"google_service_one_resource_one": ResourceDiff{
					ResourceConfig: ResourceConfigDiff{
						Old: &schema.Resource{},
						New: nil,
					},
					FlattenedSchema: FlattenedSchemaRaw{
						Old: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
						},
						New: nil,
					},
					Fields: map[string]FieldDiff{
						"field_one": FieldDiff{
							Old: &schema.Schema{Type: schema.TypeString},
							New: nil,
						},
					},
				},
			},
		},
		"new-resource": {
			newResourceMap: map[string]*schema.Resource{
				"google_service_one_resource_one": {
					Schema: map[string]*schema.Schema{
						"field_one": {
							Type: schema.TypeString,
						},
					},
				},
			},
			expectedSchemaDiff: SchemaDiff{
				"google_service_one_resource_one": ResourceDiff{
					ResourceConfig: ResourceConfigDiff{
						Old: nil,
						New: &schema.Resource{},
					},
					FlattenedSchema: FlattenedSchemaRaw{
						Old: nil,
						New: map[string]*schema.Schema{
							"field_one": {Type: schema.TypeString},
						},
					},
					Fields: map[string]FieldDiff{
						"field_one": FieldDiff{
							Old: nil,
							New: &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			schemaDiff := ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
			if diff := cmp.Diff(tc.expectedSchemaDiff, schemaDiff); diff != "" {
				t.Errorf("schema diff not equal (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestIsNewResource(t *testing.T) {
	cases := map[string]struct {
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		resourceName   string
		expected       bool
	}{
		"resource exists in both maps": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {Schema: map[string]*schema.Schema{}},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {Schema: map[string]*schema.Schema{}},
			},
			resourceName: "google_resource",
			expected:     false,
		},
		"resource only in new map": {
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {Schema: map[string]*schema.Schema{}},
			},
			resourceName: "google_resource",
			expected:     true,
		},
		"resource only in old map": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {Schema: map[string]*schema.Schema{}},
			},
			newResourceMap: map[string]*schema.Resource{},
			resourceName:   "google_resource",
			expected:       false, // ResourceConfig.New would be nil
		},
		"resource not in diff because it has no changes": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {Schema: map[string]*schema.Schema{}},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {Schema: map[string]*schema.Schema{}},
			},
			resourceName: "non_existent_resource",
			expected:     false, // Resource isn't in the diff at all
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			schemaDiff := ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)
			resourceConfig, _ := schemaDiff[tc.resourceName]
			result := resourceConfig.IsNewResource()
			if result != tc.expected {
				t.Errorf("IsNewResource(%q) = %v, want %v", tc.resourceName, result, tc.expected)
			}
		})
	}
}

func TestIsFieldInNewNestedStructure(t *testing.T) {
	cases := map[string]struct {
		oldResourceMap map[string]*schema.Resource
		newResourceMap map[string]*schema.Resource
		resourceName   string
		fieldPath      string
		expected       bool
	}{
		"top-level field in existing resource": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"old_field": {Type: schema.TypeString},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"old_field": {Type: schema.TypeString},
						"new_field": {Type: schema.TypeString},
					},
				},
			},
			resourceName: "google_resource",
			fieldPath:    "new_field",
			expected:     false, // Top-level field, not in a nested structure
		},
		"field in existing nested structure": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"nested": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"existing_field": {Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"nested": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"existing_field": {Type: schema.TypeString},
									"new_field":      {Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
			resourceName: "google_resource",
			fieldPath:    "nested.new_field",
			expected:     false, // Parent "nested" exists in old schema
		},
		"field in new nested structure": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"old_field": {Type: schema.TypeString},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"old_field": {Type: schema.TypeString},
						"new_nested": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"new_field": {Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
			resourceName: "google_resource",
			fieldPath:    "new_nested.new_field",
			expected:     true, // Parent "new_nested" doesn't exist in old schema
		},
		"field in new deeply nested structure": {
			oldResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"existing_nested": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"existing_field": {Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"existing_nested": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"existing_field": {Type: schema.TypeString},
									"new_nested": {
										Type: schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"new_field": {Type: schema.TypeString},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			resourceName: "google_resource",
			fieldPath:    "existing_nested.new_nested.new_field",
			expected:     true, // Parent "existing_nested.new_nested" doesn't exist in old schema
		},
		"field in new resource": {
			oldResourceMap: map[string]*schema.Resource{},
			newResourceMap: map[string]*schema.Resource{
				"google_resource": {
					Schema: map[string]*schema.Schema{
						"nested": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {Type: schema.TypeString},
								},
							},
						},
					},
				},
			},
			resourceName: "google_resource",
			fieldPath:    "nested.field",
			expected:     true, // New resource, so all fields are in new structures
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			schemaDiff := ComputeSchemaDiff(tc.oldResourceMap, tc.newResourceMap)

			// Verify that FlattenedSchema was properly populated
			if rd, ok := schemaDiff[tc.resourceName]; ok {
				// Debug information for test verification
				if tc.expected {
					// If we expect the field to be in a new nested structure
					// The parent path should not exist in the old schema but should exist in the new schema
					lastDotIndex := strings.LastIndex(tc.fieldPath, ".")
					if lastDotIndex != -1 {
						parentPath := tc.fieldPath[:lastDotIndex]
						_, parentInOld := rd.FlattenedSchema.Old[parentPath]
						_, parentInNew := rd.FlattenedSchema.New[parentPath]

						// Log the verification for debugging
						t.Logf("For %s: Parent path '%s' exists in old schema: %v, exists in new schema: %v",
							tc.fieldPath, parentPath, parentInOld, parentInNew)

						// This should match our expectation
						if parentInOld || !parentInNew {
							t.Errorf("For field %s: Expected parent path %s to not exist in old schema and exist in new schema, but got old: %v, new: %v",
								tc.fieldPath, parentPath, parentInOld, parentInNew)
						}
					}
				}
			}

			// Now test the actual method
			resourceConfig := schemaDiff[tc.resourceName]
			result := resourceConfig.IsFieldInNewNestedStructure(tc.fieldPath)
			if result != tc.expected {
				t.Errorf("IsFieldInNewNestedStructure(%q, %q) = %v, want %v",
					tc.resourceName, tc.fieldPath, result, tc.expected)
			}
		})
	}
}
