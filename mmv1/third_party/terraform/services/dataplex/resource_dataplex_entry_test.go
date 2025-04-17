package dataplex

import (
	"reflect"
	"sort"
	"testing"
)

func TestAspectTypeProjectNumberDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"different project identifiers": {
			Old:                "123123131.us-central1.some-aspect",
			New:                "some-project.us-central1.some-aspect",
			ExpectDiffSuppress: true,
		},
		"different resources": {
			Old:                "123123131.us-central1.some-aspect",
			New:                "some-project.us-central1.some-other-aspect",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if AspectTypeProjectNumberDiffSuppress("diffSuppress", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestTransformAspects(t *testing.T) {
	cases := map[string]struct {
		InputObj    map[string]interface{}
		ExpectedObj map[string]interface{}
	}{
		"Empty aspects slice": {
			InputObj: map[string]interface{}{
				"aspects": []interface{}{},
				"other":   "data",
			},
			ExpectedObj: map[string]interface{}{
				"aspects": map[string]interface{}{},
				"other":   "data",
			},
		},
		"Single aspect": {
			InputObj: map[string]interface{}{
				"aspects": []interface{}{
					map[string]interface{}{"aspectKey": "aspect1", "data": "value1"},
				},
			},
			ExpectedObj: map[string]interface{}{
				"aspects": map[string]interface{}{
					"aspect1": map[string]interface{}{"data": "value1"},
				},
			},
		},
		"Multiple unique aspects": {
			InputObj: map[string]interface{}{
				"aspects": []interface{}{
					map[string]interface{}{"aspectKey": "aspectA", "value": 123, "enabled": true},
					map[string]interface{}{"aspectKey": "aspectB", "config": "settings"},
				},
				"id": "test1",
			},
			ExpectedObj: map[string]interface{}{
				"aspects": map[string]interface{}{
					"aspectA": map[string]interface{}{"value": 123, "enabled": true},
					"aspectB": map[string]interface{}{"config": "settings"},
				},
				"id": "test1",
			},
		},
		"Aspect with only aspectKey": {
			InputObj: map[string]interface{}{
				"aspects": []interface{}{
					map[string]interface{}{"aspectKey": "minimal"},
				},
			},
			ExpectedObj: map[string]interface{}{
				"aspects": map[string]interface{}{
					"minimal": map[string]interface{}{},
				},
			},
		},
		"Duplicate aspectKeys (last one wins)": {
			InputObj: map[string]interface{}{
				"aspects": []interface{}{
					map[string]interface{}{"aspectKey": "dupKey", "data": "first"},
					map[string]interface{}{"aspectKey": "otherKey", "data": "unique"},
					map[string]interface{}{"aspectKey": "dupKey", "data": "second"},
				},
			},
			ExpectedObj: map[string]interface{}{
				"aspects": map[string]interface{}{
					"dupKey":   map[string]interface{}{"data": "second"},
					"otherKey": map[string]interface{}{"data": "unique"},
				},
			},
		},
	}

	for name, tc := range cases {
		inputCopy := make(map[string]interface{}, len(tc.InputObj))
		for k, v := range tc.InputObj {
			inputCopy[k] = v
		}

		t.Run(name, func(t *testing.T) {
			TransformAspects(inputCopy)

			if !reflect.DeepEqual(inputCopy, tc.ExpectedObj) {
				t.Errorf("Test case %q failed:\nInput (after transform):\n%#v\nExpected:\n%#v",
					name, inputCopy, tc.ExpectedObj)
			}
		})
	}
}

func TestInverseTransformAspectsConcise(t *testing.T) {
	cases := map[string]struct {
		InputObj    map[string]interface{}
		ExpectedObj map[string]interface{}
	}{
		"Empty": {
			InputObj:    map[string]interface{}{"aspects": map[string]interface{}{}, "id": 1},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{}, "id": 1},
		},
		"Simple": {
			InputObj:    map[string]interface{}{"aspects": map[string]interface{}{"k1": map[string]interface{}{"d": "v1"}}},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "k1", "d": "v1"}}},
		},
		"Multiple": {
			InputObj: map[string]interface{}{"aspects": map[string]interface{}{
				"k2": map[string]interface{}{"v": 2},
				"k1": map[string]interface{}{"v": 1},
			}},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{
				map[string]interface{}{"aspectKey": "k1", "v": 1},
				map[string]interface{}{"aspectKey": "k2", "v": 2},
			}},
		},
		"InnerMapEmpty": {
			InputObj:    map[string]interface{}{"aspects": map[string]interface{}{"k_empty": map[string]interface{}{}}},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "k_empty"}}},
		},
		"InnerMapOverwritesKey": {
			InputObj:    map[string]interface{}{"aspects": map[string]interface{}{"outer": map[string]interface{}{"aspectKey": "inner", "data": "stuff"}}},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{map[string]interface{}{"aspectKey": "outer", "data": "stuff"}}},
		},
		"ComplexInnerMap": {
			InputObj: map[string]interface{}{"aspects": map[string]interface{}{
				"complex_key": map[string]interface{}{"field_a": 123, "field_b": true, "field_c": "hello world"},
			}},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{
				map[string]interface{}{"aspectKey": "complex_key", "field_a": 123, "field_b": true, "field_c": "hello world"},
			}},
		},
		"MixedInnerMaps": {
			InputObj: map[string]interface{}{"aspects": map[string]interface{}{
				"normal":    map[string]interface{}{"data": 1},
				"empty":     map[string]interface{}{},
				"overwrite": map[string]interface{}{"aspectKey": "ignored", "val": true},
			}, "other_data": "preserved"},
			ExpectedObj: map[string]interface{}{"aspects": []interface{}{
				map[string]interface{}{"aspectKey": "empty"}, // Alphabetical for stable expected def
				map[string]interface{}{"aspectKey": "normal", "data": 1},
				map[string]interface{}{"aspectKey": "overwrite", "val": true},
			}, "other_data": "preserved"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			InverseTransformAspects(tc.InputObj)

			sortSlice := func(slice []interface{}) {
				sort.SliceStable(slice, func(i, j int) bool {
					keyI := slice[i].(map[string]interface{})["aspectKey"].(string)
					keyJ := slice[j].(map[string]interface{})["aspectKey"].(string)
					return keyI < keyJ
				})
			}

			if actualSlice, ok := tc.InputObj["aspects"].([]interface{}); ok {
				sortSlice(actualSlice)
			}
			if expectedSlice, ok := tc.ExpectedObj["aspects"].([]interface{}); ok {
				sortSlice(expectedSlice)
			}

			if !reflect.DeepEqual(tc.InputObj, tc.ExpectedObj) {
				t.Errorf("Test case %q failed:\nActual (sorted slice):\n%#v\nExpected (sorted slice):\n%#v", name, tc.InputObj, tc.ExpectedObj)
			}
		})
	}
}
