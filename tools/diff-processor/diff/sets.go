package diff

import (
	"sort"
	"strings"
)

// Return the union of two maps, overwriting any shared keys with the second map's values
func union[K comparable, V any](map1, map2 map[K]V) map[K]V {
	if len(map1) == 0 {
		return map2
	}
	if len(map2) == 0 {
		return map1
	}
	merged := make(map[K]V, len(map1)+len(map2))
	for k, v := range map1 {
		merged[k] = v
	}
	for k, v := range map2 {
		merged[k] = v
	}
	return merged
}

func sliceToSetRemoveZeroPadding(slice []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, item := range slice {
		set[removeZeroPadding(item)] = struct{}{}
	}
	return set
}

// field1.0.field2 -> field1.field2
func removeZeroPadding(zeroPadded string) string {
	var trimmed string
	for _, part := range strings.Split(zeroPadded, ".") {
		if part != "0" {
			trimmed += part + "."
		}
	}
	if trimmed == "" {
		return ""
	}
	return trimmed[:len(trimmed)-1]
}

func setToSortedSlice(set map[string]struct{}) []string {
	slice := make([]string, 0, len(set))
	for item := range set {
		slice = append(slice, item)
	}
	sort.Strings(slice)
	return slice
}

func (fs FieldSet) IsSubsetOf(other FieldSet) bool {
	for field := range fs {
		if _, ok := other[field]; !ok {
			return false
		}
	}
	return true
}

func (fs FieldSet) Difference(subset FieldSet) map[string]struct{} {
	diff := make(map[string]struct{})
	for k := range fs {
		if _, ok := subset[k]; !ok {
			diff[k] = struct{}{}
		}
	}
	return diff
}
