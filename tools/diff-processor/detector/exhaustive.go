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
	Fields             map[string]*ExhaustiveFieldStatus `json:"fields"`
	IncompleteFields   []string                          `json:"incompleteFields"`
	SkippedUpdateSteps []string                          `json:"skippedUpdateSteps,omitempty"`
	SuggestedTestNote  string                            `json:"suggestedTestNote,omitempty"`
}

// DetectExhaustiveMissingTests checks whether newly added fields (and optionally
// modified existing fields when includeChangedFields is true) are exhaustively
// tested across SET, MODIFY, and UNSET transitions at every level of nesting.
func DetectExhaustiveMissingTests(schemaDiff diff.SchemaDiff, resourceMap map[string]*schema.Resource, allTests []*reader.Test, includeChangedFields bool) (map[string]*ExhaustiveMissingTestInfo, error) {
	result := make(map[string]*ExhaustiveMissingTestInfo)

	resourceNames := make([]string, 0, len(schemaDiff))
	for resourceName := range schemaDiff {
		resourceNames = append(resourceNames, resourceName)
	}
	sort.Strings(resourceNames)

	for _, resourceName := range resourceNames {
		resourceDiff := schemaDiff[resourceName]
		fullSchema := buildFullSchema(resourceDiff)
		updatable := hasUpdate(resourceMap[resourceName], fullSchema)
		targetFields := getExhaustiveTargetFields(resourceName, resourceDiff, updatable, includeChangedFields)
		if len(targetFields) == 0 {
			continue
		}

		skippedUpdateSteps := evaluateResourceTests(resourceName, updatable, fullSchema, targetFields, allTests)

		var incompleteFields []string
		for _, fieldPath := range sortedFieldPaths(targetFields) {
			status := targetFields[fieldPath]
			if status.Set.Status == StatusMissing || status.Modify.Status == StatusMissing || status.Unset.Status == StatusMissing {
				incompleteFields = append(incompleteFields, fieldPath)
			}
		}

		if len(incompleteFields) > 0 {
			result[resourceName] = &ExhaustiveMissingTestInfo{
				Fields:             targetFields,
				IncompleteFields:   incompleteFields,
				SkippedUpdateSteps: skippedUpdateSteps,
				SuggestedTestNote:  buildSuggestedTestNote(resourceName, targetFields, incompleteFields, skippedUpdateSteps),
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
	hasUpdate bool,
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
			Modify:        initialModifyStatus(newField, hasUpdate),
			Unset:         initialUnsetStatus(fieldName, newField),
		}
	}

	return targets
}

type stepBlockRecord struct {
	testName   string
	stepIndex  int // 1-based step number
	block      reader.Block
	planChecks *reader.StepPlanChecks
}

type skippedUpdateTransition struct {
	note             string
	skippedModifyFor []string
	skippedUnsetFor  []string
}

func evaluateResourceTests(
	resourceName string,
	hasUpdate bool,
	fullSchema map[string]*schema.Schema,
	targetFields map[string]*ExhaustiveFieldStatus,
	allTests []*reader.Test,
) []string {
	var skippedUpdates []skippedUpdateTransition
	sortedTargets := sortedFieldPaths(targetFields)

	omissionUnset := make(map[string]bool, len(sortedTargets))
	for _, fieldPath := range sortedTargets {
		omissionUnset[fieldPath] = usesOmissionUnset(fullSchema[fieldPath], hasUpdate)
	}

	var allResourceSteps []stepBlockRecord

	for _, test := range orderOwnResourceTestsFirst(resourceName, allTests) {
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
					testName:   test.Name,
					stepIndex:  stepIdx + 1,
					block:      block,
					planChecks: stepPlanChecksAt(test.PlanChecks, stepIdx),
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

		if !hasUpdate {
			continue
		}

		// Evaluate step-to-step MODIFY and in-place UNSET transitions per instance on updatable resources.
		for _, label := range instanceOrder {
			recs := instanceSteps[label]
			for i := 1; i < len(recs); i++ {
				prev := recs[i-1]
				curr := recs[i]

				if ok, skipReason := verifyInPlaceUpdateTransition(resourceName, label, prev, curr); !ok {
					var skippedModify, skippedUnset []string
					for _, fieldPath := range sortedTargets {
						st := targetFields[fieldPath]
						prevPresent := isFieldPresent(prev.block, fieldPath)
						currPresent := isFieldPresent(curr.block, fieldPath)
						if st.Modify.Status != StatusExempt && isModifiedInUpdateStep(prev.block, curr.block, prevPresent, currPresent, fieldPath, fullSchema) {
							skippedModify = append(skippedModify, fieldPath)
						}
						if st.Unset.Status != StatusExempt && !omissionUnset[fieldPath] && prevPresent && !currPresent && isParentOrSoleChildAncestorPresent(curr.block, fieldPath, fullSchema) {
							skippedUnset = append(skippedUnset, fieldPath)
						}
					}
					if len(skippedModify) > 0 || len(skippedUnset) > 0 {
						skippedUpdates = append(skippedUpdates, skippedUpdateTransition{
							note:             fmt.Sprintf("%s (step %d -> %d: %s)", test.Name, prev.stepIndex, curr.stepIndex, skipReason),
							skippedModifyFor: skippedModify,
							skippedUnsetFor:  skippedUnset,
						})
					}
					continue
				}

				transitionLabel := fmt.Sprintf("%s (step %d -> %d)", test.Name, prev.stepIndex, curr.stepIndex)

				for _, fieldPath := range sortedTargets {
					st := targetFields[fieldPath]
					prevPresent := isFieldPresent(prev.block, fieldPath)
					currPresent := isFieldPresent(curr.block, fieldPath)

					// Check MODIFY.
					if st.Modify.Status == StatusMissing && isModifiedInUpdateStep(prev.block, curr.block, prevPresent, currPresent, fieldPath, fullSchema) {
						st.Modify = TransitionResult{
							Status:    StatusCovered,
							CoveredBy: transitionLabel,
						}
					}

					// Check in-place UNSET for mutable fields/blocks. For nested fields like
					// `foo.bar.baz`, `foo.bar` must remain present (unless `baz` is the sole
					// configurable child in `foo.bar`'s schema) so removing `foo` altogether
					// does not skip unsetting `baz` on multi-child blocks.
					if st.Unset.Status == StatusMissing && !omissionUnset[fieldPath] &&
						prevPresent && !currPresent && isParentOrSoleChildAncestorPresent(curr.block, fieldPath, fullSchema) {
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
			if isParentOrSoleChildAncestorPresent(rec.block, fieldPath, fullSchema) {
				st.Unset = TransitionResult{
					Status:    StatusCovered,
					CoveredBy: stepLabel,
					UnsetMode: detectUnsetMode(rec.block, fieldPath),
				}
			}
		}
	}

	var relevantSkipped []string
	seenSkipped := make(map[string]bool)
	for _, su := range skippedUpdates {
		if seenSkipped[su.note] {
			continue
		}
		stillMissing := false
		for _, fp := range su.skippedModifyFor {
			if targetFields[fp].Modify.Status == StatusMissing {
				stillMissing = true
				break
			}
		}
		if !stillMissing {
			for _, fp := range su.skippedUnsetFor {
				if targetFields[fp].Unset.Status == StatusMissing {
					stillMissing = true
					break
				}
			}
		}
		if stillMissing {
			seenSkipped[su.note] = true
			relevantSkipped = append(relevantSkipped, su.note)
		}
	}
	return relevantSkipped
}

func stepPlanChecksAt(planChecks []reader.StepPlanChecks, stepIdx int) *reader.StepPlanChecks {
	if stepIdx >= 0 && stepIdx < len(planChecks) {
		return &planChecks[stepIdx]
	}
	return nil
}

// verifyInPlaceUpdateTransition reports whether the transition from prev to
// curr is backed by a PreApply plan check proving the resource was updated in
// place rather than destroyed and recreated.
//
// ReadTestFiles populates Test.PlanChecks 1:1 with Test.Steps, so
// curr.planChecks is non-nil for every test parsed from Go source and nil only
// for synthetic unit-test fixtures that omit Test.PlanChecks.
func verifyInPlaceUpdateTransition(resourceName, label string, prev, curr stepBlockRecord) (bool, string) {
	addr := resourceName + "." + label
	if curr.stepIndex != prev.stepIndex+1 {
		return false, fmt.Sprintf("%s was absent in an intervening step", addr)
	}
	pc := curr.planChecks
	if pc == nil {
		return true, ""
	}
	if pc.ExpectEmptyPlan {
		return false, "plancheck expects empty plan"
	}
	if action, ok := pc.ResourceActions[addr]; ok {
		if action == "ResourceActionUpdate" {
			return true, ""
		}
		return false, fmt.Sprintf("plancheck expects %s on %s", action, addr)
	}
	if pc.ExpectNoDelete {
		return true, ""
	}
	return false, fmt.Sprintf("missing plancheck.ResourceActionUpdate on %s", addr)
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

// orderOwnResourceTestsFirst returns a deterministically ordered copy of
// allTests with the resource's own TestAcc<Resource>_* tests placed ahead of
// other tests, ordered within each group by shorter test name first and then
// alphabetically so concise tests (such as _basic or _full) are preferred and
// CoveredBy is stable regardless of map iteration order in test-reader.
func orderOwnResourceTestsFirst(resourceName string, allTests []*reader.Test) []*reader.Test {
	ordered := make([]*reader.Test, len(allTests))
	copy(ordered, allTests)
	sort.SliceStable(ordered, func(i, j int) bool {
		iOwn := isOwnResourceTest(ordered[i].Name, resourceName)
		jOwn := isOwnResourceTest(ordered[j].Name, resourceName)
		if iOwn != jOwn {
			return iOwn
		}
		if len(ordered[i].Name) != len(ordered[j].Name) {
			return len(ordered[i].Name) < len(ordered[j].Name)
		}
		return ordered[i].Name < ordered[j].Name
	})
	return ordered
}

func isOwnResourceTest(testName, resourceName string) bool {
	trimmed := strings.TrimPrefix(resourceName, "google_")
	prefix := "testacc" + strings.ReplaceAll(trimmed, "_", "")
	lower := strings.ToLower(testName)
	return lower == prefix || strings.HasPrefix(lower, prefix+"_")
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

// isModifiedInUpdateStep reports whether an in-place update step modifies
// fieldPath. Multi-value fields require two distinct non-null values
// (prevPresent && currPresent), whereas optional single-value enums have only
// two distinct states (unset and their sole value) and therefore satisfy MODIFY
// via an in-place unset -> value transition.
func isModifiedInUpdateStep(
	prevBlock, currBlock reader.Block,
	prevPresent, currPresent bool,
	fieldPath string,
	fullSchema map[string]*schema.Schema,
) bool {
	if prevPresent && currPresent {
		return isFieldModified(prevBlock, currBlock, fieldPath)
	}
	if !prevPresent && currPresent && isSingleValueEnum(fullSchema[fieldPath]) {
		return isParentOrSoleChildAncestorPresent(prevBlock, fieldPath, fullSchema)
	}
	return false
}

// isParentOrSoleChildAncestorPresent reports whether fieldPath's immediate
// parent is present in block, or a purely Optional (non-Computed) ancestor
// whose only configurable schema subtree is a defaultless fieldPath was omitted
// while that ancestor's parent remained present (avoiding an invalid empty
// ancestor {} block when unsetting its sole optional child).
//
// Two cases are excluded because `ancestor {}` is distinct from omitting
// `ancestor`:
//   - Optional+Computed (Computed: true) ancestors: omitting `ancestor` retains
//     prior state in SDKv2 without clearing child attributes, whereas writing
//     `ancestor {}` clears the child.
//   - Children with a schema Default/DefaultFunc: writing `ancestor {}`
//     populates the child's default value (making `ancestor {}` non-empty in
//     SDKv2), whereas omitting `ancestor` removes the block without testing the
//     default.
func isParentOrSoleChildAncestorPresent(block reader.Block, fieldPath string, fullSchema map[string]*schema.Schema) bool {
	if isParentPresent(block, fieldPath) {
		return true
	}
	if fieldSchema := fullSchema[fieldPath]; fieldSchema != nil && (fieldSchema.Default != nil || fieldSchema.DefaultFunc != nil) {
		return false
	}
	for anc := parentPath(fieldPath); anc != ""; anc = parentPath(anc) {
		ancSchema := fullSchema[anc]
		if ancSchema != nil {
			if ancSchema.Computed {
				return false
			}
			if ancSchema.Required {
				continue
			}
		}
		if !isFieldPresent(block, anc) && isParentPresent(block, anc) &&
			isOnlyConfigurableChildInSchema(anc, fieldPath, fullSchema) {
			return true
		}
	}
	return false
}

// isOnlyConfigurableChildInSchema reports whether ancestorPath has no Required
// or Optional schema fields outside targetPath's lineage, meaning ancestorPath
// cannot remain non-empty in HCL once targetPath is omitted.
func isOnlyConfigurableChildInSchema(ancestorPath, targetPath string, fullSchema map[string]*schema.Schema) bool {
	ancPrefix := ancestorPath + "."
	targetPrefix := targetPath + "."
	for k, sch := range fullSchema {
		if !strings.HasPrefix(k, ancPrefix) {
			continue
		}
		if k == targetPath || strings.HasPrefix(k, targetPrefix) {
			continue
		}
		if strings.HasPrefix(targetPath, k+".") {
			continue
		}
		if sch != nil && (sch.Optional || sch.Required) {
			return false
		}
	}
	return true
}

func sortedFieldPaths(fields map[string]*ExhaustiveFieldStatus) []string {
	paths := make([]string, 0, len(fields))
	for p := range fields {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return paths
}
