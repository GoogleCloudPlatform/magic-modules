package openapi_generate

import (
	"context"
	"github.com/getkin/kin-openapi/openapi3"
	"testing"
)

func TestMapType(t *testing.T) {
	_ = NewOpenapiParser("/fake", "/fake")
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, _ := loader.LoadFromFile("./test_data/test_api.yaml")
	_ = doc.Validate(ctx)

	petSchema := doc.Paths.Map()["/pets"].Post.Parameters[0].Value.Schema
	mmObject := WriteObject("pet", petSchema, propType(petSchema), false)
	if mmObject.KeyName == "" || mmObject.Type != "Map" {
		t.Error("Failed to parse map type")
	}
	if len(mmObject.ValueType.Properties) != 4 {
		t.Errorf("Expected 4 properties, found %d", len(mmObject.ValueType.Properties))
	}
}
