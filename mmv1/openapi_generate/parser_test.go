package openapi_generate

import (
	_ "embed"
	"path/filepath"
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

func TestAttachStandardFunctionality(t *testing.T) {
	resource := api.Resource{
		Name:     "TestResource",
		SelfLink: "test/self/link",
	}

	result := attachStandardFunctionality(resource)

	if len(result.Examples) != 0 {
		t.Errorf("Expected result.Examples to be empty, got: %d examples", len(result.Examples))
	}

	if len(result.Samples) != 1 {
		t.Fatalf("Expected result.Samples to have length 1, got: %d", len(result.Samples))
	}

	sample := result.Samples[0]
	if sample.Name != "name_of_sample_file" {
		t.Errorf("Expected sample.Name to be 'name_of_sample_file', got: %q", sample.Name)
	}

	if sample.PrimaryResourceId != "example" {
		t.Errorf("Expected sample.PrimaryResourceId to be 'example', got: %q", sample.PrimaryResourceId)
	}

	if len(sample.Steps) != 1 {
		t.Fatalf("Expected sample.Steps to have length 1, got: %d", len(sample.Steps))
	}

	step := sample.Steps[0]
	if step.Name != "name_of_sample_file" {
		t.Errorf("Expected step.Name to be 'name_of_sample_file', got: %q", step.Name)
	}

	if val, ok := step.Vars["resource_name"]; !ok || val != "test-resource" {
		t.Errorf("Expected step.Vars['resource_name'] to be 'test-resource', got: %q (present=%t)", val, ok)
	}
}

func TestDeprecatedParameterSkip(t *testing.T) {
	yamlData := []byte(`
openapi: "3.0.0"
info:
  version: 1.0.0
  title: Test API
paths:
  /test:
    get:
      parameters:
        - name: activeParam
          in: query
          description: Active parameter
          schema:
            type: string
        - name: deprecatedParam
          in: query
          description: Deprecated. Use activeParam instead.
          schema:
            type: string
        - name: optionalDeprecatedParam
          in: query
          description: Optional. Deprecated. Skip this too.
          schema:
            type: string
      responses:
        default:
          description: OK
`)
	ctx := t.Context()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromData(yamlData)
	if err != nil {
		t.Fatalf("Could not load data %s", err)
	}
	err = doc.Validate(ctx)
	if err != nil {
		t.Fatalf("Could not validate data %s", err)
	}

	op := doc.Paths.Map()["/test"].Get
	resArray := parseOpenApi("/test", "Test", op)
	parameters := resArray[0].([]*api.Type)

	var activeFound, deprecatedFound, optionalDeprecatedFound bool
	for _, p := range parameters {
		if p.Name == "activeParam" {
			activeFound = true
		}
		if p.Name == "deprecatedParam" {
			deprecatedFound = true
		}
		if p.Name == "optionalDeprecatedParam" {
			optionalDeprecatedFound = true
		}
	}

	if !activeFound {
		t.Error("Expected to find activeParam in parameters list")
	}
	if deprecatedFound {
		t.Error("Expected deprecatedParam to be skipped, but it was found")
	}
	if optionalDeprecatedFound {
		t.Error("Expected optionalDeprecatedParam to be skipped, but it was found")
	}
}

func TestBuildProduct(t *testing.T) {
	tests := []struct {
		name           string
		spec           string
		expectedBase   string
		expectedRepUrl string
	}{
		{
			name: "GlobalAndRegionalEndpoints",
			spec: `
openapi: "3.0.3"
info:
  version: v1
  title: Test Service API
servers:
  - url: https://testservice.googleapis.com
    description: Global Endpoint
  - url: https://asia-east1-testservice.googleapis.com/
    description: Locational Endpoint
    x-google-endpoint-location: asia-east1
  - url: https://testservice.africa-south1.rep.googleapis.com/
    description: Regional Endpoint
    x-google-deprecated: true
    x-google-endpoint-location: africa-south1
  - url: https://testservice.asia-east1.rep.googleapis.com/
    description: Regional Endpoint
    x-google-endpoint-location: asia-east1
paths: {}
`,
			expectedBase:   "https://testservice.googleapis.com/v1/",
			expectedRepUrl: "https://testservice.{{location}}.rep.googleapis.com/v1/",
		},
		{
			name: "GlobalOnlyNoRegionalEndpoints",
			spec: `
openapi: "3.0.3"
info:
  version: v1
  title: Test Service API
servers:
  - url: https://testservice.googleapis.com
    description: Global Endpoint
  - url: https://asia-east1-testservice.googleapis.com/
    description: Locational Endpoint
    x-google-endpoint-location: asia-east1
paths: {}
`,
			expectedBase:   "https://testservice.googleapis.com/v1/",
			expectedRepUrl: "",
		},
		{
			name: "ServiceNameContainsLocationSubstring",
			spec: `
openapi: "3.0.3"
info:
  version: v1
  title: Test Service API
servers:
  - url: https://lustre.googleapis.com
    description: Global Endpoint
  - url: https://lustre.us.rep.googleapis.com/
    description: Regional Endpoint
    x-google-endpoint-location: us
paths: {}
`,
			expectedBase:   "https://lustre.googleapis.com/v1/",
			expectedRepUrl: "https://lustre.{{location}}.rep.googleapis.com/v1/",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()
			loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
			doc, err := loader.LoadFromData([]byte(tc.spec))
			if err != nil {
				t.Fatalf("Could not load data: %v", err)
			}
			if err := doc.Validate(ctx); err != nil {
				t.Fatalf("Could not validate data: %v", err)
			}

			outDir := t.TempDir()
			productPath := buildProduct("/fake/testservice_v1.json", outDir, doc, header)

			var product api.Product
			api.Compile(filepath.Join(productPath, "product.yaml"), &product)
			product.Validate()

			if product.Name != "TestService" {
				t.Errorf("Expected product.Name to be 'TestService', got %q", product.Name)
			}
			if product.DisplayName != "Test Service" {
				t.Errorf("Expected product.DisplayName to be 'Test Service', got %q", product.DisplayName)
			}
			if len(product.Versions) != 1 {
				t.Fatalf("Expected 1 version in product.Versions, got %d", len(product.Versions))
			}

			v := product.Versions[0]
			if v.Name != "ga" {
				t.Errorf("Expected version.Name to be 'ga', got %q", v.Name)
			}
			if v.BaseUrl != tc.expectedBase {
				t.Errorf("Expected version.BaseUrl to be %q, got %q", tc.expectedBase, v.BaseUrl)
			}
			if v.RepUrl != tc.expectedRepUrl {
				t.Errorf("Expected version.RepUrl to be %q, got %q", tc.expectedRepUrl, v.RepUrl)
			}
		})
	}
}
