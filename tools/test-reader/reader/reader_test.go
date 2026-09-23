package reader

import (
	"go/ast"
	"go/parser"
	"os"
	"reflect"
	"testing"
)

// This test only ensures there isn't a panic reading tests in the provider.
func TestReadAllTests(t *testing.T) {
	if servicesDir := os.Getenv("SERVICES_DIR"); servicesDir != "" {
		_, errs := ReadAllTests(servicesDir)
		for path, err := range errs {
			t.Logf("path: %s, err: %v", path, err)
		}
	} else {
		t.Log("no services directory provided, skipping TestReadAllTests")
	}
}

func TestReadCoveredResourceTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/covered_resource_test.go"})
	if err != nil {
		t.Fatalf("error reading covered resource test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if len(tests[0].Steps) != 2 {
		t.Fatalf("unexpected number of test steps: %d, expected 2", len(tests[0].Steps))
	}
	if coveredResources, ok := tests[0].Steps[0][ResourceBlock]["covered_resource"]; !ok {
		t.Errorf("did not find covered_resource in %v", tests[0].Steps[0])
	} else if coveredResource, ok := coveredResources["resource"]; !ok {
		t.Errorf("did not find a covered resource in %v", coveredResources)
	} else if expectedResource := (Block{
		"field_four.field_five.field_six": "true",
		"field_one":                       "\"value-one\"",
		"field_seven":                     "true",
	}); !reflect.DeepEqual(coveredResource, expectedResource) {
		t.Errorf("found wrong fields in covered resource config: %#v, expected %#v", coveredResource, expectedResource)
	}
}

func TestReadConfigVariableTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/config_variable_test.go"})
	if err != nil {
		t.Fatalf("error reading config variable test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if len(tests[0].Steps) != 1 {
		t.Fatalf("unexpected number of test steps: %d, expected 1", len(tests[0].Steps))
	}
	if configVariableResources, ok := tests[0].Steps[0][ResourceBlock]["config_variable"]; !ok {
		t.Errorf("did not find config_variable in %v", tests[0].Steps[0])
	} else if configVariableResource, ok := configVariableResources["basic"]; !ok {
		t.Errorf("did not find a resource in %v", configVariableResources)
	} else if expectedResource := (Block{"field_one": "\"value-one\""}); !reflect.DeepEqual(configVariableResource, expectedResource) {
		t.Errorf("found wrong fields in config variable config: %#v, expected %#v", configVariableResource, expectedResource)
	}
}

func TestReadMultipleResourcesTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/multiple_resource_test.go"})
	if err != nil {
		t.Fatalf("error reading multiple resources test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if expectedSteps := []Step{
		{
			ResourceBlock: {
				"resource_one": {
					"instace_two":  {"field_one": "\"value-one\""},
					"instance_one": {"field_one": "\"value-one\""},
				},
				"resource_two": {
					"instace_one": {"field_one": "\"value-one\""},
					"instace_two": {"field_one": "\"value-one\""},
				},
			},
		},
		{
			ResourceBlock: {
				"resource_one": {
					"instace_two":  {"field_one": "\"value-two\""},
					"instance_one": {"field_one": "\"value-two\""},
				},
				"resource_two": {
					"instace_one": {"field_one": "\"value-two\""},
					"instace_two": {"field_one": "\"value-two\""},
				},
			},
		},
	}; !reflect.DeepEqual(tests[0].Steps, expectedSteps) {
		t.Errorf("found unexpected test steps for multiple resources: %#v, expected %#v", tests[0].Steps, expectedSteps)
	}
}

func TestReadSerialResourceTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/serial_resource_test.go"})
	if err != nil {
		t.Fatalf("error reading serial resource test file: %v", err)
	}
	if len(tests) != 2 {
		t.Fatalf("unexpected number of tests: %d, expected 2", len(tests))
	}
	if expectedTests := []*Test{
		{
			Name: "testAccSerialResource1",
			Steps: []Step{
				{
					ResourceBlock: {
						"serial_resource": {
							"resource": {"field_one": "\"value-one\""},
						},
					},
				},
			},
		},
		{
			Name: "testAccSerialResource2",
			Steps: []Step{
				{
					ResourceBlock: {
						"serial_resource": {
							"resource": {
								"field_two.field_three": "\"value-two\"",
							},
						},
					},
				},
			},
		},
	}; !reflect.DeepEqual(tests, expectedTests) {
		t.Errorf("found unexpected serialized tests: %v, expected %v", tests, expectedTests)
	}

}

func TestReadCrossFileTests(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/cross_file_1_test.go", "testdata/service/cross_file_2_test.go"})
	if err != nil {
		t.Fatalf("error reading cross file tests: %v", err)
	}

	expectedTests := []*Test{
		{
			Name: "testAccCrossFile1",
			Steps: []Step{
				{
					ResourceBlock: {
						"serial_resource": {
							"resource": {"field_one": "\"value-one\""},
						},
					},
				},
			},
		},
		{
			Name: "testAccCrossFile2",
			Steps: []Step{
				{
					ResourceBlock: {
						"serial_resource": {
							"resource": {
								"field_two.field_three": "\"value-two\"",
							},
						},
					},
				},
			},
		},
	}

	if len(tests) != len(expectedTests) {
		t.Fatalf("unexpected number of tests: %d, expected %d", len(tests), len(expectedTests))
	}

	if !reflect.DeepEqual(tests, expectedTests) {
		t.Errorf("found unexpected cross file tests: %v, expected %v", tests, expectedTests)
	}

}

func TestReadHelperFunctionCall(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/function_call_test.go"})
	if err != nil {
		t.Fatalf("error reading function call test: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	expectedTest := &Test{
		Name: "TestAccFunctionCallResource",
		Steps: []Step{
			{
				ResourceBlock: TypeLabels{
					"helped_resource": Blocks{
						"primary": Block{
							"field_one": "\"value-one\"",
						},
					},
					"helper_resource": Blocks{
						"default": Block{
							"field_one": "\"value-one\"",
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(tests[0], expectedTest) {
		t.Errorf("found unexpected tests using helper function: %v, expected %v", tests[0], expectedTest)
	}
}

func TestReadBlockTypesTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/block_types_test.go"})
	if err != nil {
		t.Fatalf("error reading block types test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if len(tests[0].Steps) != 1 {
		t.Fatalf("unexpected number of test steps: %d, expected 1", len(tests[0].Steps))
	}
	step := tests[0].Steps[0]
	// Blocks with a type and a name label are read regardless of the block
	// type, blocks with any other number of labels are ignored, and each block
	// type gets its own bucket even when they all share a type label.
	if expectedStep := (Step{
		ResourceBlock: {
			"block_types_resource": {"resource": {"field_one": "var.var_one"}},
		},
		DataBlock: {
			"block_types_resource": {"data": {"field_two": "\"value-two\""}},
		},
		EphemeralBlock: {
			"block_types_resource": {"ephemeral": {"field_three": "\"value-three\""}},
		},
		ListBlock: {
			"block_types_resource": {"list_query": {"provider": "google"}},
		},
	}); !reflect.DeepEqual(step, expectedStep) {
		t.Errorf("found unexpected step: %#v, expected %#v", step, expectedStep)
	}
}

func TestReadWholeLineSubstitutionTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/whole_line_substitution_test.go"})
	if err != nil {
		t.Fatalf("error reading whole line substitution test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	if len(tests[0].Steps) != 1 {
		t.Fatalf("unexpected number of test steps: %d, expected 1", len(tests[0].Steps))
	}
	step := tests[0].Steps[0]
	// A substitution alone on its line is dropped, so neither the setup block
	// it would have injected at the top level nor the field it would have
	// injected inside a block body is recorded. Everything else still is.
	if expectedStep := (Step{
		ResourceBlock: {
			"whole_line_substitution": {
				"resource": {
					"field_one":              "\"value-one\"",
					"field_three.field_five": "\"value-five\"",
				},
			},
		},
	}); !reflect.DeepEqual(step, expectedStep) {
		t.Errorf("found unexpected step: %#v, expected %#v", step, expectedStep)
	}
}

func TestReadFormattingCallTestFile(t *testing.T) {
	tests, err := ReadTestFiles([]string{"testdata/service/formatting_call_test.go"})
	if err != nil {
		t.Fatalf("error reading formatting call test file: %v", err)
	}
	if len(tests) != 1 {
		t.Fatalf("unexpected number of tests: %d, expected 1", len(tests))
	}
	// A config is read the same way wherever it is assembled: inline in the
	// step or in a config func, from a literal, a shared base config, a
	// concatenation of the two, or a call to another config func.
	if expectedSteps := []Step{
		{
			ResourceBlock: {
				"formatting_call_inline": {"inline": {"field_two": "\"true\""}},
			},
		},
		{
			ResourceBlock: {
				"formatting_call_base":   {"base": {"field_one": "\"value-one\""}},
				"formatting_call_concat": {"concat": {"field_three": "\"true\""}},
			},
		},
		{
			// The string passed to the nested config func is a value to
			// interpolate, not a config.
			ResourceBlock: {
				"formatting_call_string_arg": {"string_arg": {"field_four": "\"true\""}},
			},
		},
		{
			ResourceBlock: {
				"formatting_call_base":       {"base": {"field_one": "\"value-one\""}},
				"formatting_call_string_arg": {"string_arg": {"field_four": "\"true\""}},
			},
		},
		{
			ResourceBlock: {
				"formatting_call_literal": {"literal": {"field_six": "\"value-six\""}},
			},
		},
	}; !reflect.DeepEqual(tests[0].Steps, expectedSteps) {
		t.Errorf("found unexpected steps: %#v, expected %#v", tests[0].Steps, expectedSteps)
	}
}

func TestReadConfigCallExpr(t *testing.T) {
	for _, tc := range []struct {
		name     string
		expr     string
		expected string
		wantErr  bool
	}{
		{
			name:     "sprintf-template",
			expr:     "fmt.Sprintf(`resource \"a\" \"b\" {}`, name)",
			expected: "resource \"a\" \"b\" {}",
		},
		{
			name:     "nprintf-template",
			expr:     "acctest.Nprintf(`resource \"a\" \"b\" {}`, context)",
			expected: "resource \"a\" \"b\" {}",
		},
		{
			// Sprint has no template: every argument is part of the config.
			name:     "sprint-concatenation",
			expr:     "fmt.Sprint(`resource \"a\" \"b\" {}`, `resource \"c\" \"d\" {}`)",
			expected: "resource \"a\" \"b\" {}resource \"c\" \"d\" {}",
		},
		{
			// Only a formatting function is known to return its first
			// argument's config. Reading the first argument of anything else
			// would report fields as covered that the test may never apply.
			name:    "non-formatting-call",
			expr:    "acctest.EchoResourceConfig(`resource \"a\" \"b\" {}`, \"echo\")",
			wantErr: true,
		},
		{
			name:    "undeclared-config-func",
			expr:    "testAccUndeclared(\"name\")",
			wantErr: true,
		},
		{
			name:    "formatting-call-without-arguments",
			expr:    "fmt.Sprintf()",
			wantErr: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tc.expr)
			if err != nil {
				t.Fatalf("error parsing %s: %v", tc.expr, err)
			}
			callExpr, ok := expr.(*ast.CallExpr)
			if !ok {
				t.Fatalf("%s is not a call expression", tc.expr)
			}
			configStr, err := readConfigCallExpr(callExpr, map[string]*ast.FuncDecl{}, map[string]*ast.BasicLit{})
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected an error reading %s, read %q", tc.expr, configStr)
				}
				return
			}
			if err != nil {
				t.Fatalf("error reading %s: %v", tc.expr, err)
			}
			if configStr != tc.expected {
				t.Errorf("read %q from %s, expected %q", configStr, tc.expr, tc.expected)
			}
		})
	}
}

func TestFlattenBlock(t *testing.T) {
	for _, tc := range []struct {
		name        string
		unflattened Block
		flattened   Block
	}{
		{
			name:        "empty-resource",
			unflattened: Block{},
			flattened:   Block{},
		},
		{
			name: "no-nested-fields",
			unflattened: Block{
				"a": "b",
				"c": "d",
			},
			flattened: Block{
				"a": "b",
				"c": "d",
			},
		},
		{
			name: "nested-fields",
			unflattened: Block{
				"a": Block{
					"b": Block{
						"c": "d",
					},
				},
				"e": "f",
			},
			flattened: Block{
				"a.b.c": "d",
				"e":     "f",
			},
		},
	} {
		if got := flattenBlock(tc.unflattened, ""); !reflect.DeepEqual(got, tc.flattened) {
			t.Errorf("unexpected result of flattening in test %s, expected %v, got %v", tc.name, tc.flattened, got)
		}
	}
}
