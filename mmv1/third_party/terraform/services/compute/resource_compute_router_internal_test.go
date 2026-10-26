package compute

import (
	"reflect"
	"testing"
)

func TestUnitComputeRouter_mergeMd5AuthenticationKeys(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		current  []interface{}
		oldName  interface{}
		newKey   map[string]interface{}
		expected []interface{}
	}{
		"no keys on router": {
			current: nil,
			newKey:  map[string]interface{}{"name": "router-key", "key": "secret"},
			expected: []interface{}{
				map[string]interface{}{"name": "router-key", "key": "secret"},
			},
		},
		"peer keys are kept by name": {
			current: []interface{}{
				map[string]interface{}{"name": "peer-key-1"},
				map[string]interface{}{"name": "peer-key-2"},
			},
			newKey: map[string]interface{}{"name": "router-key", "key": "secret"},
			expected: []interface{}{
				map[string]interface{}{"name": "peer-key-1"},
				map[string]interface{}{"name": "peer-key-2"},
				map[string]interface{}{"name": "router-key", "key": "secret"},
			},
		},
		"rotating the router key replaces it": {
			current: []interface{}{
				map[string]interface{}{"name": "router-key"},
				map[string]interface{}{"name": "peer-key"},
			},
			oldName: "router-key",
			newKey:  map[string]interface{}{"name": "router-key", "key": "rotated"},
			expected: []interface{}{
				map[string]interface{}{"name": "peer-key"},
				map[string]interface{}{"name": "router-key", "key": "rotated"},
			},
		},
		"renaming the router key removes the old one": {
			current: []interface{}{
				map[string]interface{}{"name": "old-router-key"},
				map[string]interface{}{"name": "peer-key"},
			},
			oldName: "old-router-key",
			newKey:  map[string]interface{}{"name": "new-router-key", "key": "secret"},
			expected: []interface{}{
				map[string]interface{}{"name": "peer-key"},
				map[string]interface{}{"name": "new-router-key", "key": "secret"},
			},
		},
		"taking over a key a peer already created": {
			current: []interface{}{
				map[string]interface{}{"name": "peer-key-1"},
				map[string]interface{}{"name": "peer-key-2"},
			},
			newKey: map[string]interface{}{"name": "peer-key-1", "key": "secret"},
			expected: []interface{}{
				map[string]interface{}{"name": "peer-key-2"},
				map[string]interface{}{"name": "peer-key-1", "key": "secret"},
			},
		},
		"only names of current keys are sent": {
			current: []interface{}{
				map[string]interface{}{"name": "peer-key", "key": "unexpected"},
			},
			newKey: map[string]interface{}{"name": "router-key", "key": "secret"},
			expected: []interface{}{
				map[string]interface{}{"name": "peer-key"},
				map[string]interface{}{"name": "router-key", "key": "secret"},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := mergeComputeRouterMd5AuthenticationKeys(tc.current, tc.oldName, tc.newKey)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
