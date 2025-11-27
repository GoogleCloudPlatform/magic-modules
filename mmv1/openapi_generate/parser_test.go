package openapi_generate

import (
	_ "embed"
	"testing"

	"google3/third_party/golang/kin_openapi/current/openapi3/openapi3"
)

//go:embed test_data/test_api.yaml
var testData []byte

func TestMapType(t *testing.T) {
	_ = NewOpenapiParser("/fake", "/fake")
	ctx := t.Context()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, _ := loader.LoadFromData(testData)
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
