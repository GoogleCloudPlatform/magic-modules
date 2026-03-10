package openapi_generate

import (
	_ "embed"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed test_data/test_api.yaml
var testData []byte

func TestMapType(t *testing.T) {
	_ = NewOpenapiParser("/fake", "/fake")
	ctx := t.Context()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromData(testData)
	if err != nil {
		t.Fatalf("Could not load data %s", err)
	}
	err = doc.Validate(ctx)
	if err != nil {
		t.Fatalf("Could not validate data %s", err)
	}

	petSchema := doc.Paths.Map()["/pets"].Post.RequestBody.Value.Content["application/json"].Schema
	mmObject := WriteObject("pet", petSchema, propType(petSchema), false)
	if mmObject.KeyName == "" || mmObject.Type != "Map" {
		t.Error("Failed to parse map type")
	}
	if len(mmObject.ValueType.Properties) != 4 {
		t.Errorf("Expected 4 properties, found %d", len(mmObject.ValueType.Properties))
	}
}

func TestFindResources(t *testing.T) {
	ctx := t.Context()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromData(testData)
	if err != nil {
		t.Fatalf("Could not load data %s", err)
	}
	err = doc.Validate(ctx)
	if err != nil {
		t.Fatalf("Could not validate data %s", err)
	}
	res := findResources(doc)
	if len(res) != 3 {
		t.Fatalf("Expected 2 resources, found: %d", len(res))
	}
	if !res["Food"].create.async {
		t.Error("Food resource is supposed to be detected as async and is not")
	}
	if res["Pet"].create.async {
		t.Error("Pet resource is not supposed to be detected as async")
	}
	if res["Breeds"].update == nil {
		t.Error("Singleton update should be found")
	}
}
