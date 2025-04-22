package dataplex_test

import (
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	dataplex "github.com/hashicorp/terraform-provider-google/google/services/dataplex"
)

func TestProjectNumberValidation(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		fieldName   string
		expectError bool
		errorMsg    string
	}{
		{"valid input", "projects/1234567890/locations/us-central1", "resourceName", false, ""},
		{"valid input with only number", "projects/987/stuff", "name", false, ""},
		{"valid input with trailing slash content", "projects/1/a/b/c", "parent", false, ""},
		{"valid input minimal", "projects/1/a", "field_a", false, ""},
		{"valid input trailing slash only", "projects/555/", "field_b", false, ""},
		{"invalid type - int", 123, "resourceName", true, `expected type of field "resourceName" to be string, but got int`},
		{"invalid type - nil", nil, "resourceName", true, `expected type of field "resourceName" to be string, but got <nil>`},
		{"invalid format - missing 'projects/' prefix", "12345/locations/us", "resourceName", true, "has an invalid format"},
		{"invalid format - project number starts with 0", "projects/0123/data", "resourceName", true, "has an invalid format"},
		{"invalid format - no project number", "projects//data", "resourceName", true, "has an invalid format"},
		{"invalid format - letters instead of number", "projects/abc/data", "resourceName", true, "has an invalid format"},
		{"invalid format - missing content after number/", "projects/123", "resourceName", true, "has an invalid format"},
		{"invalid format - empty string", "", "resourceName", true, "has an invalid format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.ProjectNumberValidation(tc.input, tc.fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("ProjectNumberValidation() error expectation mismatch: got error = %v (%v), want error = %v", hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("ProjectNumberValidation() expected error containing %q, but got: %v", tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestAspectProjectNumberValidation(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		fieldName   string
		expectError bool
		errorMsg    string
	}{
		{"valid input", "1234567890.compute.googleapis.com/Disk", "aspect_type", false, ""},
		{"valid input minimal", "1.a", "aspect_type", false, ""},
		{"valid input trailing dot only", "987.", "aspect_type", false, ""},
		{"invalid type - int", 456, "aspect_type", true, `expected type of field "aspect_type" to be string, but got int`},
		{"invalid type - nil", nil, "aspect_type", true, `expected type of field "aspect_type" to be string, but got <nil>`},
		{"invalid format - missing number", ".compute.googleapis.com/Disk", "aspect_type", true, "has an invalid format"},
		{"invalid format - number starts with 0", "0123.compute.googleapis.com/Disk", "aspect_type", true, "has an invalid format"},
		{"invalid format - missing dot", "12345compute", "aspect_type", true, "has an invalid format"},
		{"invalid format - letters instead of number", "abc.compute.googleapis.com/Disk", "aspect_type", true, "has an invalid format"},
		{"invalid format - missing content after dot", "12345", "aspect_type", true, "has an invalid format"},
		{"invalid format - empty string", "", "aspect_type", true, "has an invalid format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := dataplex.AspectProjectNumberValidation(tc.input, tc.fieldName)
			hasError := len(errors) > 0

			if hasError != tc.expectError {
				t.Fatalf("AspectProjectNumberValidation() error expectation mismatch: got error = %v (%v), want error = %v", hasError, errors, tc.expectError)
			}

			if tc.expectError && tc.errorMsg != "" {
				found := false
				for _, err := range errors {
					if strings.Contains(err.Error(), tc.errorMsg) { // Check if error message contains the expected substring
						found = true
						break
					}
				}
				if !found {
					t.Errorf("AspectProjectNumberValidation() expected error containing %q, but got: %v", tc.errorMsg, errors)
				}
			}
		})
	}
}

func TestFilterAspects(t *testing.T) {
	testCases := []struct {
		name             string
		aspectKeySet     map[string]struct{}
		resInput         map[string]interface{}
		expectedAspects  map[string]interface{}
		expectNilAspects bool
	}{
		{
			name:             "aspects is nil",
			aspectKeySet:     map[string]struct{}{"keep": {}},
			resInput:         map[string]interface{}{"otherKey": "value"},
			expectedAspects:  nil,
			expectNilAspects: true,
		},
		{
			name:            "empty aspectKeySet",
			aspectKeySet:    map[string]struct{}{},
			resInput:        map[string]interface{}{"aspects": map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}},
			expectedAspects: map[string]interface{}{},
		},
		{
			name:            "keep all aspects",
			aspectKeySet:    map[string]struct{}{"one": {}, "two": {}},
			resInput:        map[string]interface{}{"aspects": map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}},
			expectedAspects: map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}},
		},
		{
			name:            "keep some aspects",
			aspectKeySet:    map[string]struct{}{"two": {}, "three_not_present": {}},
			resInput:        map[string]interface{}{"aspects": map[string]interface{}{"one": map[string]interface{}{"data": 1}, "two": map[string]interface{}{"data": 2}}},
			expectedAspects: map[string]interface{}{"two": map[string]interface{}{"data": 2}},
		},
		{
			name:            "input aspects map is empty",
			aspectKeySet:    map[string]struct{}{"keep": {}},
			resInput:        map[string]interface{}{"aspects": map[string]interface{}{}},
			expectedAspects: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCopy := deepCopyMap(tc.resInput)
			dataplex.FilterAspects(tc.aspectKeySet, resCopy)

			actualAspectsRaw, aspectsKeyExists := resCopy["aspects"]

			if tc.expectNilAspects {
				if aspectsKeyExists && actualAspectsRaw != nil {
					t.Errorf("Expected 'aspects' to be nil or absent, but got: %v", actualAspectsRaw)
				}
				return
			}

			if !aspectsKeyExists {
				t.Fatalf("Expected 'aspects' key to exist, but it was absent")
			}

			actualAspects, ok := actualAspectsRaw.(map[string]interface{})
			if !ok {
				t.Fatalf("Expected 'aspects' to be a map[string]interface{}, but got %T", actualAspectsRaw)
			}

			if !reflect.DeepEqual(actualAspects, tc.expectedAspects) {
				t.Errorf("FilterAspects() result mismatch:\ngot:  %#v\nwant: %#v", actualAspects, tc.expectedAspects)
			}
		})
	}
}

func TestAddAspectsToSet(t *testing.T) {
	testCases := []struct {
		name         string
		initialSet   map[string]struct{}
		aspectsInput interface{}
		expectedSet  map[string]struct{}
		expectPanic  bool
	}{
		{"add to empty set", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_key": "key1"}, map[string]interface{}{"aspect_key": "key2"}}, map[string]struct{}{"key1": {}, "key2": {}}, false},
		{"add to existing set", map[string]struct{}{"existing": {}}, []interface{}{map[string]interface{}{"aspect_key": "key1"}}, map[string]struct{}{"existing": {}, "key1": {}}, false},
		{"add duplicate keys", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_key": "key1"}, map[string]interface{}{"aspect_key": "key1"}, map[string]interface{}{"aspect_key": "key2"}}, map[string]struct{}{"key1": {}, "key2": {}}, false},
		{"input aspects is empty slice", map[string]struct{}{"existing": {}}, []interface{}{}, map[string]struct{}{"existing": {}}, false},
		{"input aspects is nil", map[string]struct{}{}, nil, map[string]struct{}{}, true},
		{"input aspects is wrong type", map[string]struct{}{}, "not a slice", map[string]struct{}{}, true},
		{"item in slice is not a map", map[string]struct{}{}, []interface{}{"not a map"}, map[string]struct{}{}, true},
		{"item map missing aspect_key", map[string]struct{}{}, []interface{}{map[string]interface{}{"wrong_key": "key1"}}, map[string]struct{}{}, true},
		{"aspect_key is not a string", map[string]struct{}{}, []interface{}{map[string]interface{}{"aspect_key": 123}}, map[string]struct{}{}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			currentSet := make(map[string]struct{})
			for k, v := range tc.initialSet {
				currentSet[k] = v
			}

			defer func() {
				r := recover()
				if tc.expectPanic && r == nil {
					t.Errorf("Expected a panic, but AddAspectsToSet did not panic")
				} else if !tc.expectPanic && r != nil {
					t.Errorf("AddAspectsToSet panicked unexpectedly: %v", r)
				}

				if !tc.expectPanic {
					if !reflect.DeepEqual(currentSet, tc.expectedSet) {
						t.Errorf("AddAspectsToSet() result mismatch:\ngot:  %v\nwant: %v", currentSet, tc.expectedSet)
					}
				}
			}()

			dataplex.AddAspectsToSet(currentSet, tc.aspectsInput)
		})
	}
}

func sortAspectSlice(slice []interface{}) {
	sort.SliceStable(slice, func(i, j int) bool {
		mapI, okI := slice[i].(map[string]interface{})
		mapJ, okJ := slice[j].(map[string]interface{})
		if !okI || !okJ {
			return false
		} // Should not happen in valid tests

		keyI, okI := mapI["aspectKey"].(string)
		keyJ, okJ := mapJ["aspectKey"].(string)
		if !okI || !okJ {
			return false
		} // Should not happen in valid tests

		return keyI < keyJ
	})
}

func TestInverseTransformAspects(t *testing.T) {
	testCases := []struct {
		name             string
		resInput         map[string]interface{}
		expectedAspects  []interface{}
		expectNilAspects bool
		expectPanic      bool
	}{
		{"aspects is nil", map[string]interface{}{"otherKey": "value"}, nil, true, false},
		{"aspects is empty map", map[string]interface{}{"aspects": map[string]interface{}{}}, []interface{}{}, false, false},
		{"aspects with one entry", map[string]interface{}{"aspects": map[string]interface{}{"key1": map[string]interface{}{"data": "value1"}}}, []interface{}{map[string]interface{}{"aspectKey": "key1", "aspectValue": map[string]interface{}{"data": "value1"}}}, false, false},
		{"aspects with multiple entries", map[string]interface{}{"aspects": map[string]interface{}{"key2": map[string]interface{}{"data": "value2"}, "key1": map[string]interface{}{"data": "value1"}}}, []interface{}{map[string]interface{}{"aspectKey": "key1", "aspectValue": map[string]interface{}{"data": "value1"}}, map[string]interface{}{"aspectKey": "key2", "aspectValue": map[string]interface{}{"data": "value2"}}}, false, false},
		{"aspects is wrong type (not map)", map[string]interface{}{"aspects": "not a map"}, nil, false, true},
		{"aspect value is not a map", map[string]interface{}{"aspects": map[string]interface{}{"key1": "not a map"}}, nil, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resCopy := deepCopyMap(tc.resInput)

			defer func() {
				r := recover()
				if tc.expectPanic && r == nil {
					t.Errorf("Expected a panic, but InverseTransformAspects did not panic")
				} else if !tc.expectPanic && r != nil {
					t.Errorf("InverseTransformAspects panicked unexpectedly: %v", r)
				}

				if !tc.expectPanic {
					actualAspectsRaw, aspectsKeyExists := resCopy["aspects"]

					if tc.expectNilAspects {
						if aspectsKeyExists && actualAspectsRaw != nil {
							t.Errorf("Expected 'aspects' to be nil or absent, but got: %v", actualAspectsRaw)
						}
						return
					}

					if !aspectsKeyExists && !tc.expectNilAspects { // Should exist if not expecting nil
						t.Fatalf("Expected 'aspects' key in result map, but it was missing")
					}

					actualAspects, ok := actualAspectsRaw.([]interface{})
					if !ok && !tc.expectNilAspects { // Type check only if we didn't expect nil and key exists
						t.Fatalf("Expected 'aspects' to be []interface{}, but got %T", actualAspectsRaw)
					}

					sortAspectSlice(actualAspects)
					sortAspectSlice(tc.expectedAspects) // Ensure expected is sorted if non-nil

					if !reflect.DeepEqual(actualAspects, tc.expectedAspects) {
						t.Errorf("InverseTransformAspects() result mismatch:\ngot:  %#v\nwant: %#v", actualAspects, tc.expectedAspects)
					}
				}
			}()

			dataplex.InverseTransformAspects(resCopy)
		})
	}
}

func TestTransformAspects(t *testing.T) {
	testCases := []struct {
		name             string
		objInput         map[string]interface{}
		expectedAspects  map[string]interface{}
		expectNilAspects bool
		expectPanic      bool
	}{
		{"aspects is nil", map[string]interface{}{"otherKey": "value"}, nil, true, false},
		{"aspects is empty slice", map[string]interface{}{"aspects": []interface{}{}}, map[string]interface{}{}, false, false},
		{"aspects with one item", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspectValue": map[string]interface{}{"data": "value1"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": "value1"}}, false, false},
		{"aspects with multiple items", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspectValue": map[string]interface{}{"data": "value1"}}, map[string]interface{}{"aspectKey": "key2", "aspectValue": map[string]interface{}{"data": "value2"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": "value1"}, "key2": map[string]interface{}{"data": "value2"}}, false, false},
		{"aspects with duplicate aspectKey", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspectValue": map[string]interface{}{"data": "value_first"}}, map[string]interface{}{"aspectKey": "key2", "aspectValue": map[string]interface{}{"data": "value2"}}, map[string]interface{}{"aspectKey": "key1", "aspectValue": map[string]interface{}{"data": "value_last"}}}}, map[string]interface{}{"key1": map[string]interface{}{"data": "value_last"}, "key2": map[string]interface{}{"data": "value2"}}, false, false},
		{"aspects is wrong type (not slice)", map[string]interface{}{"aspects": "not a slice"}, nil, false, true},
		{"item in slice is not a map", map[string]interface{}{"aspects": []interface{}{"not a map"}}, nil, false, true},
		{"item map missing aspectKey", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"wrongKey": "k1", "aspectValue": map[string]interface{}{}}}}, nil, false, true},
		{"aspectKey is not a string", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": 123, "aspectValue": map[string]interface{}{}}}}, nil, false, true},
		{"item map missing aspectValue", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1"}}}, nil, false, true},
		{"aspectValue is not a map", map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "key1", "aspectValue": "not a map"}}}, nil, false, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			objCopy := deepCopyMap(tc.objInput)

			defer func() {
				r := recover()
				if tc.expectPanic && r == nil {
					t.Errorf("Expected a panic, but TransformAspects did not panic")
				} else if !tc.expectPanic && r != nil {
					t.Errorf("TransformAspects panicked unexpectedly: %v", r)
				}

				if !tc.expectPanic {
					actualAspectsRaw, aspectsKeyExists := objCopy["aspects"]

					if tc.expectNilAspects {
						if aspectsKeyExists && actualAspectsRaw != nil {
							t.Errorf("Expected 'aspects' to be nil or absent, but got: %v", actualAspectsRaw)
						}
						return
					}

					if !aspectsKeyExists && !tc.expectNilAspects {
						t.Fatalf("Expected 'aspects' key in result map, but it was missing")
					}

					actualAspects, ok := actualAspectsRaw.(map[string]interface{})
					if !ok && !tc.expectNilAspects {
						t.Fatalf("Expected 'aspects' to be map[string]interface{}, but got %T", actualAspectsRaw)
					}

					if !reflect.DeepEqual(actualAspects, tc.expectedAspects) {
						t.Errorf("TransformAspects() result mismatch:\ngot:  %#v\nwant: %#v", actualAspects, tc.expectedAspects)
					}
				}
			}()

			dataplex.TransformAspects(objCopy)
		})
	}
}

func deepCopyMap(original map[string]interface{}) map[string]interface{} {
	if original == nil {
		return nil
	}
	copyMap := make(map[string]interface{}, len(original))
	for key, value := range original {
		copyMap[key] = deepCopyValue(value)
	}
	return copyMap
}

func deepCopySlice(original []interface{}) []interface{} {
	if original == nil {
		return nil
	}
	copySlice := make([]interface{}, len(original))
	for i, value := range original {
		copySlice[i] = deepCopyValue(value)
	}
	return copySlice
}

func deepCopyValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case map[string]interface{}:
		return deepCopyMap(v)
	case []interface{}:
		return deepCopySlice(v)
	default:
		return v
	}
}

func TestAccDataplexEntry_dataplexEntryUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckDataplexEntryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexEntry_dataplexEntryFullUpdatePreapre(context),
			},
			{
				ResourceName:            "google_dataplex_entry.test_entry_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspects", "entry_group_id", "entry_id", "location"},
			},
			{
				Config: testAccDataplexEntry_dataplexEntryUpdate(context),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_dataplex_entry.test_entry_full", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_dataplex_entry.test_entry_full",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"aspects", "entry_group_id", "entry_id", "location"},
			},
		},
	})
}

func testAccDataplexEntry_dataplexEntryFullUpdatePreapre(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-one" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-one"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-two" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-two"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "story",
      "type": "enum",
      "annotations": {
        "displayName": "Story",
        "description": "Specifies the story of an entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "SEQUENCE",
          "index": 1
        },
        {
          "name": "DESERT_ISLAND",
          "index": 2
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_entry_group" "tf-test-entry-group-full%{random_suffix}" {
  entry_group_id = "tf-test-entry-group-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "tf-test-entry-type-full%{random_suffix}" {
  entry_type_id = "tf-test-entry-type-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform entry type"
  description = "entry type created by Terraform"

  type_aliases = ["TABLE", "DATABASE"]
  platform = "GCS"
  system = "CloudSQL"

  required_aspects {
    type = google_dataplex_aspect_type.tf-test-aspect-type-full%{random_suffix}-one.name
  }
}

resource "google_dataplex_entry" "test_entry_full" {
  entry_group_id = google_dataplex_entry_group.tf-test-entry-group-full%{random_suffix}.entry_group_id
  project = "%{project_number}"
  location = "us-central1"
  entry_id = "tf-test-entry-full%{random_suffix}"
  entry_type = google_dataplex_entry_type.tf-test-entry-type-full%{random_suffix}.name
  fully_qualified_name = "bigquery:%{project_number}.test-dataset"
  parent_entry = "projects/%{project_number}/locations/us-central1/entryGroups/tf-test-entry-group-full%{random_suffix}/entries/some-other-entry"
  entry_source {
    resource = "bigquery:%{project_number}.test-dataset"
    system = "System III"
    platform = "BigQuery"
    display_name = "Human readable name"
    description = "Description from source system"
    labels = {
      "old-label": "old-value"
      "some-label": "some-value"
    }

    ancestors {
      name = "ancestor-one"
      type = "type-one"
    }

    ancestors {
      name = "ancestor-two"
      type = "type-two"
    }

    create_time = "2023-08-03T19:19:00.094Z"
    update_time = "2023-08-03T20:19:00.094Z"
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-one"
    aspect_value {
      data = <<EOF
          {"type": "VIEW"    }
        EOF
    }
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-two"
    aspect_value {
      data = <<EOF
          {"story": "SEQUENCE"    }
        EOF
    }
  }
}
`, context)
}

func testAccDataplexEntry_dataplexEntryUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-one" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-one"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "type",
      "type": "enum",
      "annotations": {
        "displayName": "Type",
        "description": "Specifies the type of view represented by the entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "VIEW",
          "index": 1
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_aspect_type" "tf-test-aspect-type-full%{random_suffix}-two" {
  aspect_type_id         = "tf-test-aspect-type-full%{random_suffix}-two"
  location     = "us-central1"
  project      = "%{project_number}"

  metadata_template = <<EOF
{
  "name": "tf-test-template",
  "type": "record",
  "recordFields": [
    {
      "name": "story",
      "type": "enum",
      "annotations": {
        "displayName": "Story",
        "description": "Specifies the story of an entry."
      },
      "index": 1,
      "constraints": {
        "required": true
      },
      "enumValues": [
        {
          "name": "SEQUENCE",
          "index": 1
        },
        {
          "name": "DESERT_ISLAND",
          "index": 2
        }
      ]
    }
  ]
}
EOF
}

resource "google_dataplex_entry_group" "tf-test-entry-group-full%{random_suffix}" {
  entry_group_id = "tf-test-entry-group-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"
}

resource "google_dataplex_entry_type" "tf-test-entry-type-full%{random_suffix}" {
  entry_type_id = "tf-test-entry-type-full%{random_suffix}"
  project = "%{project_number}"
  location = "us-central1"

  labels = { "tag": "test-tf" }
  display_name = "terraform entry type"
  description = "entry type created by Terraform"

  type_aliases = ["TABLE", "DATABASE"]
  platform = "GCS"
  system = "CloudSQL"

  required_aspects {
    type = google_dataplex_aspect_type.tf-test-aspect-type-full%{random_suffix}-one.name
  }
}

resource "google_dataplex_entry" "test_entry_full" {
  entry_group_id = google_dataplex_entry_group.tf-test-entry-group-full%{random_suffix}.entry_group_id
  project = "%{project_number}"
  location = "us-central1"
  entry_id = "tf-test-entry-full%{random_suffix}"
  entry_type = google_dataplex_entry_type.tf-test-entry-type-full%{random_suffix}.name
  fully_qualified_name = "bigquery:%{project_number}.test-dataset-modified"
  parent_entry = "projects/%{project_number}/locations/us-central1/entryGroups/tf-test-entry-group-full%{random_suffix}/entries/some-other-entry"
  entry_source {
    resource = "bigquery:%{project_number}.test-dataset-modified"
    system = "System III - modified"
    platform = "BigQuery-modified"
    display_name = "Human readable name-modified"
    description = "Description from source system-modified"
    labels = {
      "some-label": "some-value-modified"
      "new-label": "new-value"
    }

    ancestors {
      name = "ancestor-one"
      type = "type-one"
    }

    ancestors {
      name = "ancestor-two"
      type = "type-two"
    }

    create_time = "2024-08-03T19:19:00.094Z"
    update_time = "2024-08-03T20:19:00.094Z"
  }

  aspects {
    aspect_key = "%{project_number}.us-central1.tf-test-aspect-type-full%{random_suffix}-one"
    aspect_value {
      data = <<EOF
     {"type": "VIEW"    }
        EOF
    }
  }
}
`, context)
}
