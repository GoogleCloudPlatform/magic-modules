package diff

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
