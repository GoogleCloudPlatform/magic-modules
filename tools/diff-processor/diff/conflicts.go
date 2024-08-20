package diff

import (
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	fieldSeparator        = ","
	conflictSetSeparator  = ";"
	conflictTypeSeparator = "/"
	numFieldConflictTypes = 4
) // separator characters for key strings

// ConflictSet is a single set of dot-separated field names which conflict with each other
type ConflictSet struct {
	Fields map[string]struct{}
	Key    string // sorted string representations of the set
}

// SetOfConflictSets represents a set of ConflictSets of a particular type
// Possible types are ConflictsWith, ExactlyOneOf, AtLeastOneOf, RequiredWith
type SetOfConflictSets struct {
	Sets    map[string]ConflictSet
	SetKeys []string // sorted array
	Key     string   // sorted string representations of all sets
}

// A set of all sets of conflicting fields in the field or resource
type FieldConflictSets struct {
	ConflictsWiths *SetOfConflictSets
	ExactlyOneOfs  *SetOfConflictSets
	AtLeastOneOfs  *SetOfConflictSets
	RequiredWiths  *SetOfConflictSets
	Key            string // sorted string representation of all sets, used for quick comparison between fields
}

type FieldConflictSetsDiff struct {
	Old *FieldConflictSets
	New *FieldConflictSets
}

func diffFieldConflictSets(ofcs, nfcs *FieldConflictSets) *FieldConflictSetsDiff {
	if ofcs == nil && nfcs == nil {
		return nil
	}

	fcsd := &FieldConflictSetsDiff{
		Old: ofcs,
		New: nfcs,
	}

	if ofcs == nil || nfcs == nil {
		return fcsd
	}

	if ofcs.Key == nfcs.Key {
		return nil
	}

	return fcsd
}

func (fcsd *FieldConflictSetsDiff) Merge(other *FieldConflictSetsDiff) {
	if fcsd.New == nil {
		fcsd.New = other.New
	}
	if fcsd.Old == nil {
		fcsd.Old = other.Old
	}
	if fcsd.New != nil {
		fcsd.New.Merge(other.New)
	}
	if fcsd.Old != nil {
		fcsd.Old.Merge(other.Old)
	}
}

// Returns a FieldConflictSets for the given field schema, or nil if the schema has no conflicts
// This object contains conflict sets for all four kinds of field conflicts
// potentially belonging to the field.
func makeFieldConflictSetsFromSchema(field *schema.Schema) *FieldConflictSets {
	if field == nil {
		return nil
	}
	conflictsWiths := makeSetOfConflictSetsFromConflictSet(makeConflictSetFromRawFieldNames(field.ConflictsWith))
	exactlyOneOfs := makeSetOfConflictSetsFromConflictSet(makeConflictSetFromRawFieldNames(field.ExactlyOneOf))
	atLeastOneOfs := makeSetOfConflictSetsFromConflictSet(makeConflictSetFromRawFieldNames(field.AtLeastOneOf))
	requiredWiths := makeSetOfConflictSetsFromConflictSet(makeConflictSetFromRawFieldNames(field.RequiredWith))

	fcs := &FieldConflictSets{
		ConflictsWiths: conflictsWiths,
		ExactlyOneOfs:  exactlyOneOfs,
		AtLeastOneOfs:  atLeastOneOfs,
		RequiredWiths:  requiredWiths,
	}

	conflictTypeKeys := make([]string, numFieldConflictTypes)

	for i, socs := range []*SetOfConflictSets{
		conflictsWiths,
		exactlyOneOfs,
		atLeastOneOfs,
		requiredWiths,
	} {
		if socs == nil {
			conflictTypeKeys[i] = ""
		} else {
			conflictTypeKeys[i] = socs.Key
		}
	}
	fcs.Key = strings.Join(conflictTypeKeys, conflictTypeSeparator)
	if len(fcs.Key) < numFieldConflictTypes {
		// There are no field conflicts.
		return nil
	}
	return fcs
}

func makeFieldConflictSetsFromKey(key string) *FieldConflictSets {
	keyParts := strings.Split(key, conflictTypeSeparator)
	if len(keyParts) < numFieldConflictTypes {
		return nil
	}

	return &FieldConflictSets{
		ConflictsWiths: makeSetOfConflictSetsFromKey(keyParts[0]),
		ExactlyOneOfs:  makeSetOfConflictSetsFromKey(keyParts[1]),
		AtLeastOneOfs:  makeSetOfConflictSetsFromKey(keyParts[2]),
		RequiredWiths:  makeSetOfConflictSetsFromKey(keyParts[3]),
		Key:            key,
	}
}

// Make a single conflict set into a set of sets with one element
func makeSetOfConflictSetsFromConflictSet(conflictSet ConflictSet) *SetOfConflictSets {
	if conflictSet.Key == "" {
		return nil
	}
	return &SetOfConflictSets{
		Sets: map[string]ConflictSet{
			conflictSet.Key: conflictSet,
		},
		SetKeys: []string{conflictSet.Key},
		Key:     conflictSet.Key,
	}
}

func makeSetOfConflictSetsFromKey(key string) *SetOfConflictSets {
	if key == "" {
		return nil
	}
	conflictSetKeys := strings.Split(key, conflictSetSeparator)
	conflictSets := make(map[string]ConflictSet, len(conflictSetKeys))
	for _, csk := range conflictSetKeys {
		if csk == "" {
			continue
		}
		conflictSets[csk] = makeConflictSetFromKey(csk)
	}
	return &SetOfConflictSets{
		Sets:    conflictSets,
		SetKeys: conflictSetKeys,
		Key:     key,
	}
}

// Returns a ConflictSet from the given conflict schema values
// Calls removeZeroPading on each value before creating ConflictSet
func makeConflictSetFromRawFieldNames(rawFieldNames []string) ConflictSet {
	if len(rawFieldNames) == 0 {
		return ConflictSet{}
	}
	fields := make(map[string]struct{}, len(rawFieldNames))
	keyParts := make([]string, len(rawFieldNames))
	for i, field := range rawFieldNames {
		trimmedField := removeZeroPadding(field)
		fields[trimmedField] = struct{}{}
		keyParts[i] = trimmedField
	}
	sort.Strings(keyParts)
	return ConflictSet{
		Fields: fields,
		Key:    strings.Join(keyParts, fieldSeparator),
	}
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

func makeConflictSetFromKey(key string) ConflictSet {
	fieldNames := strings.Split(key, fieldSeparator)
	fieldSet := make(map[string]struct{}, len(fieldNames))
	for _, fieldName := range fieldNames {
		fieldSet[fieldName] = struct{}{}
	}
	return ConflictSet{
		Fields: fieldSet,
		Key:    key,
	}
}

// Merge adds all field conflicts of all types from the given field conflict sets to this set of sets.
func (fcs *FieldConflictSets) Merge(other *FieldConflictSets) {
	if other == nil {
		return
	}

	fcs.ConflictsWiths = mergeSetsOfConflictSets(fcs.ConflictsWiths, other.ConflictsWiths)
	fcs.ExactlyOneOfs = mergeSetsOfConflictSets(fcs.ExactlyOneOfs, other.ExactlyOneOfs)
	fcs.AtLeastOneOfs = mergeSetsOfConflictSets(fcs.AtLeastOneOfs, other.AtLeastOneOfs)
	fcs.RequiredWiths = mergeSetsOfConflictSets(fcs.RequiredWiths, other.RequiredWiths)

	allSocs := []*SetOfConflictSets{
		fcs.ConflictsWiths,
		fcs.ExactlyOneOfs,
		fcs.AtLeastOneOfs,
		fcs.RequiredWiths,
	}
	mergedKeys := make([]string, len(allSocs))
	for i, socs := range allSocs {
		if socs == nil {
			mergedKeys[i] = ""
		} else {
			mergedKeys[i] = socs.Key
		}
	}
	fcs.Key = strings.Join(mergedKeys, conflictTypeSeparator)
}

// merges two sets of conflict sets, either of which can be nil
func mergeSetsOfConflictSets(socs1, socs2 *SetOfConflictSets) *SetOfConflictSets {
	if socs1 == nil {
		return socs2
	}
	if socs2 == nil {
		return socs1
	}
	socs1.Merge(socs2)
	return socs1
}

// Merge adds all conflict sets from the given set to this one.
func (socs *SetOfConflictSets) Merge(other *SetOfConflictSets) {
	if other == nil {
		return
	}

	socs.Sets = union(socs.Sets, other.Sets)
	socs.SetKeys = sortedUnion(socs.SetKeys, other.SetKeys)
	socs.Key = strings.Join(socs.SetKeys, conflictSetSeparator)
}
