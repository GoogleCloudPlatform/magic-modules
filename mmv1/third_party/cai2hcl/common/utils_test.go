package common

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tpg_provider "github.com/hashicorp/terraform-provider-google-beta/google-beta/provider"
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

func createSchema(name string) map[string]*schema.Schema {
	provider := tpg_provider.Provider()

	return provider.ResourcesMap[name].Schema
}
