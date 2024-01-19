package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	"github.com/zclconf/go-cty/cty"
)

func TestSubsetOfFieldsMapsToCtyValue(t *testing.T) {
	schema := createSchema("google_compute_forwarding_rule")

	outputMap := map[string]interface{}{
		"name": "forwarding-rule-1",
	}

	val, err := MapToCtyValWithSchema(outputMap, schema)

	assert.Nil(t, err)
	assert.Equal(t, "forwarding-rule-1", val.GetAttr("name").AsString())
}

func TestWrongFieldTypeBreaksConversion(t *testing.T) {
	resourceSchema := createSchema("google_compute_backend_service")
	outputMap := map[string]interface{}{
		"name":        "fr-1",
		"description": []string{"unknownValue"}, // string is required, not array.
	}

	val, err := MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.True(t, val.IsNull())
	assert.Contains(t, err.Error(), "string is required")
}

func TestNilValue(t *testing.T) {
	resourceSchema := createSchema("google_compute_forwarding_rule")
	outputMap := map[string]interface{}{
		"name":        "fr-1",
		"description": nil,
	}

	val, err := MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)
	assert.Equal(t, cty.Value(cty.StringVal("fr-1")), val.GetAttr("name"))
	assert.Equal(t, cty.Value(cty.NullVal(cty.String)), val.GetAttr("description"))
}

func TestNilValueInRequiredField(t *testing.T) {
	resourceSchema := createSchema("google_compute_forwarding_rule")
	outputMap := map[string]interface{}{
		"name": nil,
	}

	val, err := MapToCtyValWithSchema(outputMap, resourceSchema)

	// In future we may want to fail in this case.
	assert.Nil(t, err)
	assert.Equal(t, cty.Value(cty.NullVal(cty.String)), val.GetAttr("name"))
}

func TestFieldsWithTypeSlice(t *testing.T) {
	resourceSchema := createSchema("google_compute_forwarding_rule")
	outputMap := map[string]interface{}{
		"name":  "fr-1",
		"ports": []string{"80"},
	}

	val, err := MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)

	assert.Equal(t, []cty.Value{cty.StringVal("80")}, val.GetAttr("ports").AsValueSlice())
}

func TestMissingFieldDoesNotBreakConversionConversionWhenOutputNormalized(t *testing.T) {
	resourceSchema := createSchema("google_compute_forwarding_rule")
	outputMap := map[string]interface{}{
		"name":         "fr-1",
		"unknownField": "unknownValue",
	}

	val, err := MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)

	assert.True(t, val.Type().HasAttribute("name"))
	assert.Equal(t, "fr-1", val.GetAttr("name").AsString())
}

func TestFieldWithTypeSchemaSetDoesNotBreakConversionWhenOutputNormalized(t *testing.T) {
	resourceSchema := createSchema("google_compute_forwarding_rule")
	outputMap := map[string]interface{}{
		"name":  "fr-1",
		"ports": schema.NewSet(schema.HashString, tpgresource.ConvertStringArrToInterface([]string{"80"})),
	}

	val, err := MapToCtyValWithSchema(outputMap, resourceSchema)

	assert.Nil(t, err)
	assert.Equal(t, []cty.Value{cty.StringVal("80")}, val.GetAttr("ports").AsValueSlice())
}

func createSchema(name string) map[string]*schema.Schema {
	provider := tpg_provider.Provider()

	return provider.ResourcesMap[name].Schema
}
