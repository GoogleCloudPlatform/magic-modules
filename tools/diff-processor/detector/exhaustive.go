package detector

import (
	"fmt"
	"sort"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TransitionType string

const (
	TransitionSet    TransitionType = "SET"
	TransitionModify TransitionType = "MODIFY"
	TransitionUnset  TransitionType = "UNSET"
)

type TransitionStatus string

const (
	StatusCovered TransitionStatus = "COVERED"
	StatusMissing TransitionStatus = "MISSING"
	StatusExempt  TransitionStatus = "EXEMPT"
)

type UnsetMode string

const (
	UnsetModeNull    UnsetMode = "null"
	UnsetModeOmitted UnsetMode = "omitted"
)

type TransitionResult struct {
	Status     TransitionStatus `json:"status"`
	CoveredBy  string           `json:"coveredBy,omitempty"`
	UnsetMode  UnsetMode        `json:"unsetMode,omitempty"`
	ExemptCode string           `json:"exemptCode,omitempty"`
}

type ExhaustiveFieldStatus struct {
	FieldPath     string           `json:"fieldPath"`
	IsNestedBlock bool             `json:"isNestedBlock"`
	Set           TransitionResult `json:"set"`
	Modify        TransitionResult `json:"modify"`
	Unset         TransitionResult `json:"unset"`
}

type ExhaustiveMissingTestInfo struct {
	Fields            map[string]*ExhaustiveFieldStatus `json:"fields"`
	IncompleteFields  []string                          `json:"incompleteFields"`
	SuggestedTestNote string                            `json:"suggestedTestNote,omitempty"`
}

// DetectExhaustiveMissingTests checks whether newly added fields (and optionally
// modified existing fields when includeChangedFields is true) are exhaustively
// tested across SET, MODIFY, and UNSET transitions at every level of nesting.
func DetectExhaustiveMissingTests(schemaDiff diff.SchemaDiff, allTests []*reader.Test, includeChangedFields bool) (map[string]*ExhaustiveMissingTestInfo, error) {
	result := make(map[string]*ExhaustiveMissingTestInfo)

	resourceNames := make([]string, 0, len(schemaDiff))
	for resourceName := range schemaDiff {
		resourceNames = append(resourceNames, resourceName)
	}
	sort.Strings(resourceNames)

	for _, resourceName := range resourceNames {
		resourceDiff := schemaDiff[resourceName]
		fullSchema := buildFullSchema(resourceDiff)
		targetFields := getExhaustiveTargetFields(resourceName, resourceDiff, includeChangedFields)
		if len(targetFields) == 0 {
			continue
		}

		evaluateResourceTests(resourceName, resourceDiff, fullSchema, targetFields, allTests)

		var incompleteFields []string
		for _, fieldPath := range sortedFieldPaths(targetFields) {
			status := targetFields[fieldPath]
			if status.Set.Status == StatusMissing || status.Modify.Status == StatusMissing || status.Unset.Status == StatusMissing {
				incompleteFields = append(incompleteFields, fieldPath)
			}
		}

		if len(incompleteFields) > 0 {
			result[resourceName] = &ExhaustiveMissingTestInfo{
				Fields:            targetFields,
				IncompleteFields:  incompleteFields,
				SuggestedTestNote: buildSuggestedTestNote(resourceName, targetFields, incompleteFields),
			}
		}
	}

	return result, nil
}

// buildFullSchema returns a flattened map of all schema paths on the resource,
// falling back to resourceDiff.Fields (and walking nested *schema.Resource
// elements) when unit tests construct a ResourceDiff without FlattenedSchema.New.
func buildFullSchema(resourceDiff diff.ResourceDiff) map[string]*schema.Schema {
	fullSchema := make(map[string]*schema.Schema)
	for k, v := range resourceDiff.FlattenedSchema.New {
		if v != nil {
			fullSchema[k] = v
		}
	}
	for k, fd := range resourceDiff.Fields {
		if fd.New != nil {
			fullSchema[k] = fd.New
			populateNestedResourceSchema(k, fd.New, fullSchema)
		}
	}
	return fullSchema
}

func populateNestedResourceSchema(prefix string, sch *schema.Schema, fullSchema map[string]*schema.Schema) {
	if sch == nil {
		return
	}
	res, ok := sch.Elem.(*schema.Resource)
	if !ok || res == nil {
		return
	}
	for childName, childSch := range res.Schema {
		if childSch == nil {
			continue
		}
		childPath := prefix + "." + childName
		if _, exists := fullSchema[childPath]; !exists {
			fullSchema[childPath] = childSch
		}
		populateNestedResourceSchema(childPath, childSch, fullSchema)
	}
}

func isIAMResource(resourceName string) bool {
	return strings.HasSuffix(resourceName, "_iam_member") ||
		strings.HasSuffix(resourceName, "_iam_policy") ||
		strings.HasSuffix(resourceName, "_iam_binding")
}

func getExhaustiveTargetFields(
	resourceName string,
	resourceDiff diff.ResourceDiff,
	includeChangedFields bool,
) map[string]*ExhaustiveFieldStatus {
	targets := make(map[string]*ExhaustiveFieldStatus)
	isIAM := isIAMResource(resourceName)

	for fieldName, fieldDiff := range resourceDiff.Fields {
		newField := fieldDiff.New
		if newField == nil {
			continue
		}
		if !includeChangedFields && fieldDiff.Old != nil {
			continue
		}
		if fieldName == "project" || fieldName == "deletion_policy" {
			continue
		}
		if isIAM && (fieldName == "condition" || strings.HasPrefix(fieldName, "condition.")) {
			continue
		}
		if newField.Computed && !newField.Optional && !newField.Required {
			continue
		}

		_, isNestedBlock := newField.Elem.(*schema.Resource)
		targets[fieldName] = &ExhaustiveFieldStatus{
			FieldPath:     fieldName,
			IsNestedBlock: isNestedBlock,
			Set:           TransitionResult{Status: StatusMissing},
			Modify:        initialModifyStatus(newField, resourceDiff),
			Unset:         initialUnsetStatus(fieldName, newField),
		}
	}

	return targets
}

type stepBlockRecord struct {
	testName  string
	stepIndex int // 1-based step number
	block     reader.Block
}

func evaluateResourceTests(
	resourceName string,
	resourceDiff diff.ResourceDiff,
	fullSchema map[string]*schema.Schema,
	targetFields map[string]*ExhaustiveFieldStatus,
	allTests []*reader.Test,
) {
	sortedTargets := sortedFieldPaths(targetFields)

	omissionUnset := make(map[string]bool, len(sortedTargets))
	for _, fieldPath := range sortedTargets {
		omissionUnset[fieldPath] = usesOmissionUnset(fullSchema[fieldPath], resourceDiff)
	}

	var allResourceSteps []stepBlockRecord

	for _, test := range allTests {
		// Group steps by resource instance label so multi-step transitions follow the same instance.
		instanceSteps := make(map[string][]stepBlockRecord)
		var instanceOrder []string

		for stepIdx, step := range test.Steps {
			resourceBlocks, ok := step[reader.ResourceBlock][resourceName]
			if !ok {
				continue
			}
			labels := make([]string, 0, len(resourceBlocks))
			for label := range resourceBlocks {
				labels = append(labels, label)
			}
			sort.Strings(labels)

			for _, label := range labels {
				block := resourceBlocks[label]
				rec := stepBlockRecord{
					testName:  test.Name,
					stepIndex: stepIdx + 1,
					block:     block,
				}
				if _, seen := instanceSteps[label]; !seen {
					instanceOrder = append(instanceOrder, label)
				}
				instanceSteps[label] = append(instanceSteps[label], rec)
				allResourceSteps = append(allResourceSteps, rec)

				// Check SET coverage on this step.
				for _, fieldPath := range sortedTargets {
					st := targetFields[fieldPath]
					if isFieldPresent(block, fieldPath) {
						if st.Set.Status == StatusMissing {
							st.Set = TransitionResult{
								Status:    StatusCovered,
								CoveredBy: fmt.Sprintf("%s (step %d)", test.Name, stepIdx+1),
							}
						}
					}
				}
			}
		}

		if resourceDiff.NoUpdate {
			continue
		}

		// Evaluate step-to-step MODIFY and in-place UNSET transitions per instance on updatable resources.
		for _, label := range instanceOrder {
			recs := instanceSteps[label]
			for i := 1; i < len(recs); i++ {
				prev := recs[i-1]
				curr := recs[i]

				transitionLabel := fmt.Sprintf("%s (step %d -> %d)", test.Name, prev.stepIndex, curr.stepIndex)

				for _, fieldPath := range sortedTargets {
					st := targetFields[fieldPath]
					prevPresent := isFieldPresent(prev.block, fieldPath)
					currPresent := isFieldPresent(curr.block, fieldPath)

					// Check MODIFY.
					if st.Modify.Status == StatusMissing && prevPresent && currPresent {
						if isFieldModified(prev.block, curr.block, fieldPath) {
							st.Modify = TransitionResult{
								Status:    StatusCovered,
								CoveredBy: transitionLabel,
							}
						}
					}

					// Check in-place UNSET for mutable fields/blocks. For nested fields like
					// `foo.bar.baz`, `foo.bar` must remain present so removing `foo` altogether
					// does not count as unsetting `baz`.
					if st.Unset.Status == StatusMissing && !omissionUnset[fieldPath] &&
						prevPresent && !currPresent && isParentPresent(curr.block, fieldPath) {
						st.Unset = TransitionResult{
							Status:    StatusCovered,
							CoveredBy: transitionLabel,
							UnsetMode: detectUnsetMode(curr.block, fieldPath),
						}
					}
				}
			}
		}
	}

	// Evaluate omission UNSET across all valid apply steps.
	for _, rec := range allResourceSteps {
		stepLabel := fmt.Sprintf("%s (step %d)", rec.testName, rec.stepIndex)
		for _, fieldPath := range sortedTargets {
			st := targetFields[fieldPath]
			if !omissionUnset[fieldPath] || st.Unset.Status != StatusMissing {
				continue
			}
			if isFieldPresent(rec.block, fieldPath) {
				continue
			}
			if isParentPresent(rec.block, fieldPath) {
				st.Unset = TransitionResult{
					Status:    StatusCovered,
					CoveredBy: stepLabel,
					UnsetMode: detectUnsetMode(rec.block, fieldPath),
				}
			}
		}
	}
}

// isFieldPresent reports whether fieldPath is configured with a non-null value
// in block.
//
// Nested HCL blocks cannot be assigned `= null`—they are only unset by omitting
// the block altogether. Writing `bar {}` (recorded as `"{}"`) or
// `bar { foo = null }` (recorded as `"bar.foo": "null"`) still declares the
// `bar` block in HCL while leaving `bar.foo` unset, so `isFieldPresent` returns
// true for `"bar"` and false for `"bar.foo"`.
func isFieldPresent(block reader.Block, fieldPath string) bool {
	if val, ok := block[fieldPath]; ok && !reader.IsNullValue(val) {
		return true
	}
	prefix := fieldPath + "."
	for k := range block {
		if strings.HasPrefix(k, prefix) {
			return true
		}
	}
	return false
}

// detectUnsetMode reports whether a field already known to be unset in block
// was cleared via an explicit `= null` assignment or by omitting it.
func detectUnsetMode(block reader.Block, fieldPath string) UnsetMode {
	if val, ok := block[fieldPath]; ok && reader.IsNullValue(val) {
		return UnsetModeNull
	}
	return UnsetModeOmitted
}

// nonNullValue strips explicit `"null"` entries from a flattened Block value and
// reports whether any non-null expression remains.
//
// When a resource repeats blocks of the same type, test-reader joins their
// attribute expressions with "\x00" (e.g. `"\"a\"\x00null"`). Filtering out
// `"null"` parts treats `foo = null` inside a repeated block the same as
// omitting `foo` on that block, and returns ("", false) when every instance is
// null.
func nonNullValue(v any) (string, bool) {
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v), true
	}
	if s == "" {
		return "", true
	}
	var kept []string
	for _, p := range strings.Split(s, "\x00") {
		if strings.TrimSpace(p) != reader.NullValue {
			kept = append(kept, p)
		}
	}
	if len(kept) == 0 {
		return "", false
	}
	return strings.Join(kept, "\x00"), true
}

func isParentPresent(block reader.Block, fieldPath string) bool {
	parent := parentPath(fieldPath)
	if parent == "" {
		return true
	}
	return isFieldPresent(block, parent)
}

func parentPath(fieldPath string) string {
	idx := strings.LastIndexByte(fieldPath, '.')
	if idx == -1 {
		return ""
	}
	return fieldPath[:idx]
}

// getFieldTreeValue returns a deterministic snapshot of all non-null values at
// or beneath fieldPath so step-to-step comparisons can detect modifications at
// any nesting depth. A present block with no non-null children (either `bar {}`
// or `bar { foo = null }`) normalizes to `fieldPath + "={}\n"`.
func getFieldTreeValue(block reader.Block, fieldPath string) string {
	prefix := fieldPath + "."
	var keys []string
	vals := make(map[string]string)
	for k, v := range block {
		if k == fieldPath || strings.HasPrefix(k, prefix) {
			if nv, ok := nonNullValue(v); ok {
				keys = append(keys, k)
				vals[k] = nv
			}
		}
	}
	if len(keys) == 0 {
		if isFieldPresent(block, fieldPath) {
			return fieldPath + "={}\n"
		}
		return ""
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&b, "%s=%s\n", k, vals[k])
	}
	return b.String()
}

func isFieldModified(prevBlock, currBlock reader.Block, fieldPath string) bool {
	return getFieldTreeValue(prevBlock, fieldPath) != getFieldTreeValue(currBlock, fieldPath)
}

func sortedFieldPaths(fields map[string]*ExhaustiveFieldStatus) []string {
	paths := make([]string, 0, len(fields))
	for p := range fields {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}
