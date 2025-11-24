package diff

import (
	"sort"
	"strings"
)

// FieldSet is a set of strings representing fields.
type FieldSet map[string]struct{}

// Difference returns the fields in s that are not in other.
func (s FieldSet) Difference(other FieldSet) FieldSet {
	diff := make(FieldSet)
	for k := range s {
		if _, ok := other[k]; !ok {
			diff[k] = struct{}{}
		}
	}
	return diff
}

// IsSubsetOf returns true if s is a subset of other.
func (s FieldSet) IsSubsetOf(other FieldSet) bool {
	for k := range s {
		if _, ok := other[k]; !ok {
			return false
		}
	}
	return true
}

// Intersection returns the fields that are in both s and other.
func (s FieldSet) Intersection(other FieldSet) FieldSet {
	intersection := make(FieldSet)
	for k := range s {
		if _, ok := other[k]; ok {
			intersection[k] = struct{}{}
		}
	}
	return intersection
}

func sliceToSet(slice []string) FieldSet {
	set := make(FieldSet)
	for _, s := range slice {
		if s != "" {
			set[s] = struct{}{}
		}
	}
	return set
}

func sliceToSetRemoveZeroPadding(slice []string) FieldSet {
	set := make(FieldSet)
	for _, s := range slice {
		if s != "" {
			set[strings.ReplaceAll(s, ".0", "")] = struct{}{}
		}
	}
	return set
}

func setToSortedSlice(set FieldSet) []string {
	slice := make([]string, 0, len(set))
	for k := range set {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	return slice
}

func union[T any](a, b map[string]T) map[string]struct{} {
	c := make(map[string]struct{})
	for k := range a {
		c[k] = struct{}{}
	}
	for k := range b {
		c[k] = struct{}{}
	}
	return c
}
