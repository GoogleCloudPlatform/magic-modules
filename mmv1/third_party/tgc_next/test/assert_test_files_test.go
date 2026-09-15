package test

import (
	"testing"
)

func TestFindMissingKeys(t *testing.T) {
	tests := []struct {
		name          string
		map1          map[string]any
		map2          map[string]any
		ignoredFields map[string]any
		want          []string
	}{
		{
			name: "identical maps",
			map1: map[string]any{"foo": "bar", "num": 123},
			map2: map[string]any{"foo": "bar", "num": 123},
			want: nil,
		},
		{
			name: "missing key in map2",
			map1: map[string]any{"foo": "bar", "num": 123, "missing": "yes"},
			map2: map[string]any{"foo": "bar", "num": 123},
			want: []string{"missing"},
		},
		{
			name:          "missing key ignored",
			map1:          map[string]any{"foo": "bar", "missing": "yes"},
			map2:          map[string]any{"foo": "bar"},
			ignoredFields: map[string]any{"missing": struct{}{}},
			want:          nil,
		},
		{
			name: "nested map missing key",
			map1: map[string]any{
				"foo": "bar",
				"nested": map[string]any{
					"inner":   "val",
					"missing": "yes",
				},
			},
			map2: map[string]any{
				"foo": "bar",
				"nested": map[string]any{
					"inner": "val",
				},
			},
			want: []string{"nested.missing"},
		},
		{
			name: "nested map missing key ignored",
			map1: map[string]any{
				"foo": "bar",
				"nested": map[string]any{
					"inner":   "val",
					"missing": "yes",
				},
			},
			map2: map[string]any{
				"foo": "bar",
				"nested": map[string]any{
					"inner": "val",
				},
			},
			ignoredFields: map[string]any{"nested.missing": struct{}{}},
			want:          nil,
		},
		{
			name: "slice of maps missing key",
			map1: map[string]any{
				"list": []any{
					map[string]any{"key": "val", "missing": "yes"},
				},
			},
			map2: map[string]any{
				"list": []any{
					map[string]any{"key": "val"},
				},
			},
			want: []string{"list.0.missing"},
		},
		{
			name: "slice of maps missing key partial ignore",
			map1: map[string]any{
				"list": []any{
					map[string]any{"key": "val", "missing": "yes"},
				},
			},
			map2: map[string]any{
				"list": []any{
					map[string]any{"key": "val"},
				},
			},
			ignoredFields: map[string]any{"list.missing": struct{}{}},
			want:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findMissingKeys(tt.map1, tt.map2, "", tt.ignoredFields)
			if len(got) != len(tt.want) {
				t.Errorf("findMissingKeys() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("findMissingKeys() = %v, want %v", got, tt.want)
					break
				}
			}
		})
	}
}
