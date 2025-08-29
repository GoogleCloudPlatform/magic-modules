package utils_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/cai2hcl/converters/utils"
	tpg_provider "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/provider"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/tpgresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
)

func TestSubsetOfFieldsMapsToCtyValue(t *testing.T) {
	schema := createSchema("google_compute_instance")

	outputMap := map[string]interface{}{
		"name": "forwarding-rule-1",
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, schema)

	assert.Nil(t, err)
	assert.Equal(t, "forwarding-rule-1", val.GetAttr("name").AsString())
}

func TestWrongFieldTypeBreaksConversion(t *testing.T) {
	resourceSchema := createSchema("google_compute_instance")
	outputMap := map[string]interface{}{
		"name":        "fr-1",
		"description": []string{"unknownValue"}, // string is required, not array.
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.True(t, val.IsNull())
	assert.Contains(t, err.Error(), "string is required")
}

func TestNilValue(t *testing.T) {
	resourceSchema := createSchema("google_compute_instance")
	outputMap := map[string]interface{}{
		"name":        "fr-1",
		"description": nil,
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)
	assert.Equal(t, cty.Value(cty.StringVal("fr-1")), val.GetAttr("name"))
	assert.Equal(t, cty.Value(cty.NullVal(cty.String)), val.GetAttr("description"))
}

func TestNilValueInRequiredField(t *testing.T) {
	resourceSchema := createSchema("google_compute_instance")
	outputMap := map[string]interface{}{
		"name": nil,
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, resourceSchema)

	// In future we may want to fail in this case.
	assert.Nil(t, err)
	assert.Equal(t, cty.Value(cty.NullVal(cty.String)), val.GetAttr("name"))
}

func TestFieldsWithTypeSlice(t *testing.T) {
	resourceSchema := createSchema("google_compute_instance")
	outputMap := map[string]interface{}{
		"name":              "fr-1",
		"resource_policies": []string{"test"},
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)

	assert.Equal(t, []cty.Value{cty.StringVal("test")}, val.GetAttr("resource_policies").AsValueSlice())
}

func TestMissingFieldDoesNotBreakConversionConversion(t *testing.T) {
	resourceSchema := createSchema("google_compute_instance")
	outputMap := map[string]interface{}{
		"name":         "fr-1",
		"unknownField": "unknownValue",
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)

	assert.True(t, val.Type().HasAttribute("name"))
	assert.Equal(t, "fr-1", val.GetAttr("name").AsString())
}

func TestFieldWithTypeSchemaSet(t *testing.T) {
	resourceSchema := createSchema("google_compute_instance")
	outputMap := map[string]interface{}{
		"name":              "fr-1",
		"resource_policies": schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface([]string{"test"})),
	}

	val, err := utils.MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)
	assert.Equal(t, []cty.Value{cty.StringVal("test")}, val.GetAttr("resource_policies").AsValueSlice())
}

func TestFieldWithTypeSchemaListAndNestedObject(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"list": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"nested_key": {
						Type: schema.TypeString,
					},
				},
			},
		},
	}
	flattenedMap := map[string]interface{}{
		"list": []interface{}{
			map[string]interface{}{
				"nested_key":         "value",
				"nested_unknown_key": "unknown_key_value",
			},
		},
	}

	val, err := utils.MapToCtyValWithSchema(flattenedMap, resourceSchema)

	assert.Nil(t, err)
	assert.Equal(t,
		[]cty.Value{
			cty.ObjectVal(
				map[string]cty.Value{
					"nested_key": cty.StringVal("value"),
				},
			),
		},
		val.GetAttr("list").AsValueSlice(),
	)
}

func TestFieldWithTypeSchemaSetAndNestedObject(t *testing.T) {
	nestedResource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nested_key": {
				Type: schema.TypeString,
			},
		},
	}
	resourceSchema := map[string]*schema.Schema{
		"list": {
			Type: schema.TypeSet,
			Elem: nestedResource,
		},
	}

	flattenedMap := map[string]interface{}{
		"list": schema.NewSet(schema.HashResource(nestedResource), []interface{}{
			map[string]interface{}{
				"nested_key":         "value",
				"nested_unknown_key": "unknown_key_value",
			},
		}),
	}

	val, err := utils.MapToCtyValWithSchema(flattenedMap, resourceSchema)

	assert.Nil(t, err)
	assert.Equal(t,
		[]cty.Value{
			cty.ObjectVal(
				map[string]cty.Value{
					"nested_key": cty.StringVal("value"),
				},
			)},
		val.GetAttr("list").AsValueSlice())
}

func createSchema(name string) map[string]*schema.Schema {
	provider := tpg_provider.Provider()

	return provider.ResourcesMap[name].Schema
}

func TestParseUrlParamValuesFromAssetName(t *testing.T) {
	compareMaps := func(m1, m2 map[string]any) error {
		if diff := cmp.Diff(m1, m2, cmpopts.SortMaps(func(k1, k2 string) bool { return k1 < k2 })); diff != "" {
			return fmt.Errorf("maps are not equal (-got +want):\n%s", diff)
		}
		return nil
	}

	// Test cases for different scenarios
	testCases := []struct {
		name         string
		template     string
		assetName    string
		outputFields map[string]struct{}
		want         map[string]any
	}{
		{
			name:         "ComputeUrlmap",
			template:     "//compute.googleapis.com/projects/{{project}}/global/urlMaps/{{name}}",
			assetName:    "//compute.googleapis.com/projects/my-project/global/urlMaps/urlmapibgtchooyo",
			outputFields: make(map[string]struct{}),
			want:         map[string]any{"project": "my-project", "name": "urlmapibgtchooyo"},
		},
		{
			name:         "BigQueryDataset",
			template:     "//bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}",
			assetName:    "//bigquery.googleapis.com/projects/my-project/datasets/my-dataset",
			outputFields: make(map[string]struct{}),
			want:         map[string]any{"project": "my-project", "dataset_id": "my-dataset"},
		},
		{
			name:         "AlloyDBInstance",
			template:     "//alloydb.googleapis.com/{{cluster}}/instances/{{instance_id}}",
			assetName:    "//alloydb.googleapis.com/projects/ci-test/locations/us-central1/clusters/tf-test-cluster/instances/tf-test-instance",
			outputFields: make(map[string]struct{}),
			want:         map[string]any{"cluster": "projects/ci-test/locations/us-central1/clusters/tf-test-cluster", "instance_id": "tf-test-instance"},
		},
		{
			name:         "WithOutputFieldsIgnored",
			template:     "//bigquery.googleapis.com/projects/{{project}}/location/{{location}}/datasets/{{dataset_id}}",
			assetName:    "//bigquery.googleapis.com/projects/my-project/location/abc/datasets/my-dataset",
			outputFields: map[string]struct{}{"location": {}}, // 'location' should be ignored
			want:         map[string]any{"project": "my-project", "dataset_id": "my-dataset"},
		},
		{
			name:         "WithMissingSuffix",
			template:     "//bigquery.googleapis.com/projects/{{project/datasets/{{dataset_id}}",
			assetName:    "//bigquery.googleapis.com/projects/my-project/datasets/my-dataset",
			outputFields: make(map[string]struct{}),
			want:         map[string]any{"dataset_id": "my-dataset"},
		},
		{
			name:         "EmptyTemplate",
			template:     "",
			assetName:    "//bigquery.googleapis.com/projects/my-project/datasets/my-dataset",
			outputFields: make(map[string]struct{}),
			want:         map[string]any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hclData := make(map[string]any)
			utils.ParseUrlParamValuesFromAssetName(tc.assetName, tc.template, tc.outputFields, hclData)

			if err := compareMaps(hclData, tc.want); err != nil {
				t.Fatalf("map mismatch: %v", err)
			}
		})
	}
}
