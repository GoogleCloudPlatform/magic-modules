package openapi_generate

import (
	_ "embed"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
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
	mmObject := WriteObject("pet", petSchema, propType(petSchema), false, make(map[string]bool), make(map[*openapi3.Schema]bool))
	if mmObject.KeyName == "" || mmObject.Type != "Map" {
		t.Error("Failed to parse map type")
	}
	if len(mmObject.ValueType.Properties) != 6 {
		t.Errorf("Expected 6 properties, found %d", len(mmObject.ValueType.Properties))
	}

	var recursivePet *api.Type
	for _, p := range mmObject.ValueType.Properties {
		if p.Name == "recursivePet" {
			recursivePet = p
			break
		}
	}
	if recursivePet == nil {
		t.Error("Failed to find recursivePet property")
	}

	var secondRecursivePet *api.Type
	for _, p := range recursivePet.Properties {
		if p.Name == "recursivePet" {
			secondRecursivePet = p
			break
		}
	}
	if secondRecursivePet == nil {
		t.Error("Failed to find second-level recursivePet property")
	} else if secondRecursivePet.Type != "String" || secondRecursivePet.CustomExpand == "" {
		t.Errorf("Expected second-level recursivePet to be a JSON field (String with templates), found Type: %s", secondRecursivePet.Type)
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

func TestReadOnlyPropagation(t *testing.T) {
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
	mmObject := WriteObject("pet", petSchema, propType(petSchema), false, make(map[string]bool), make(map[*openapi3.Schema]bool))

	var readOnlyObject *api.Type
	for _, p := range mmObject.ValueType.Properties {
		if p.Name == "readOnlyObject" {
			readOnlyObject = p
			break
		}
	}

	if readOnlyObject == nil {
		t.Fatal("Failed to find readOnlyObject property")
	}

	if !readOnlyObject.Output {
		t.Error("Expected readOnlyObject to have Output=true")
	}

	// Check nested string property
	var nameProp *api.Type
	var nestedProp *api.Type
	var arrayNestedProp *api.Type
	for _, p := range readOnlyObject.Properties {
		switch p.Name {
		case "name":
			nameProp = p
		case "nested":
			nestedProp = p
		case "arrayNested":
			arrayNestedProp = p
		}
	}

	if nameProp == nil || !nameProp.Output {
		t.Errorf("Expected property 'name' under 'readOnlyObject' to be parsed and have Output=true, got: %+v", nameProp)
	}

	if nestedProp == nil || !nestedProp.Output {
		t.Errorf("Expected property 'nested' under 'readOnlyObject' to be parsed and have Output=true, got: %+v", nestedProp)
	} else {
		var subNameProp *api.Type
		for _, p := range nestedProp.Properties {
			if p.Name == "subName" {
				subNameProp = p
				break
			}
		}
		if subNameProp == nil || !subNameProp.Output {
			t.Errorf("Expected 'subName' under 'nested' to have Output=true, got: %+v", subNameProp)
		}
	}

	if arrayNestedProp == nil || !arrayNestedProp.Output {
		t.Errorf("Expected property 'arrayNested' under 'readOnlyObject' to be parsed and have Output=true, got: %+v", arrayNestedProp)
	} else {
		if arrayNestedProp.ItemType == nil || !arrayNestedProp.ItemType.Output {
			t.Errorf("Expected 'ItemType' of 'arrayNested' to have Output=true, got: %+v", arrayNestedProp.ItemType)
		} else {
			var itemSubNameProp *api.Type
			for _, p := range arrayNestedProp.ItemType.Properties {
				if p.Name == "itemSubName" {
					itemSubNameProp = p
					break
				}
			}
			if itemSubNameProp == nil || !itemSubNameProp.Output {
				t.Errorf("Expected 'itemSubName' under 'arrayNested.ItemType' to have Output=true, got: %+v", itemSubNameProp)
			}
		}
	}
}

