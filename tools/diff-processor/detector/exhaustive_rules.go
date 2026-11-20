package detector

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	enumProbeKey = "__diff_processor_enum_probe__"
	enumProbeVal = "\x00__diff_processor_enum_probe__\x00"
	enumErrStart = "expected " + enumProbeKey + " to be one of ["
	enumErrEnd   = "], got " + enumProbeVal
)

func initialModifyStatus(field *schema.Schema, hasUpdate bool) TransitionResult {
	if !hasUpdate {
		return TransitionResult{Status: StatusExempt, ExemptCode: "resource_not_updatable"}
	}
	if field.ForceNew {
		return TransitionResult{Status: StatusExempt, ExemptCode: "force_new"}
	}
	// Exempt 1-value enums only when Required or backed by a non-empty schema
	// Default/DefaultFunc (only one value is reachable). Optional 1-value enums
	// without a schema default have two distinct states—unset and the enum
	// value—so MODIFY remains required and can be satisfied by an in-place
	// unset -> value transition.
	if isSingleValueEnum(field) && (field.Required || hasSchemaDefault(field)) {
		return TransitionResult{Status: StatusExempt, ExemptCode: "single_value_enum"}
	}
	return TransitionResult{Status: StatusMissing}
}

func initialUnsetStatus(fieldPath string, field *schema.Schema) TransitionResult {
	if field.Required {
		if !strings.Contains(fieldPath, ".") {
			return TransitionResult{Status: StatusExempt, ExemptCode: "required"}
		}
		return TransitionResult{Status: StatusExempt, ExemptCode: "required_in_parent"}
	}
	if field.Optional && field.Computed {
		return TransitionResult{Status: StatusExempt, ExemptCode: "optional_computed"}
	}
	return TransitionResult{Status: StatusMissing}
}

// usesOmissionUnset returns true when a field or block cannot be unset via an in-place
// update step (e.g. non-updatable resource or ForceNew field) and therefore satisfies
// UNSET by being omitted in at least one valid apply step where its parent is present.
func usesOmissionUnset(field *schema.Schema, hasUpdate bool) bool {
	if !hasUpdate {
		return true
	}
	if field != nil && field.ForceNew {
		return true
	}
	return false
}

// hasUpdate reports whether a resource supports in-place updates—requiring both
// an Update/UpdateContext callback and at least one mutable configurable field
// (MMv1 generates a stub Update function on otherwise immutable resources solely
// to persist client-side `deletion_policy` changes in Terraform state).
//
// Synthetic unit-test fixtures that pass a nil resource or omit CRUD
// callbacks altogether are treated as supporting updates.
func hasUpdate(resource *schema.Resource, fullSchema map[string]*schema.Schema) bool {
	if resource == nil {
		return true
	}
	hasCRUD := resource.Create != nil || resource.Read != nil || resource.CreateContext != nil || resource.ReadContext != nil
	if !hasCRUD {
		return true
	}
	if resource.Update == nil && resource.UpdateContext == nil {
		return false
	}
	return hasMutableConfigurableField(fullSchema)
}

func hasMutableConfigurableField(fullSchema map[string]*schema.Schema) bool {
	hasAnyConfigurable := false
	for path, sch := range fullSchema {
		if sch == nil || path == "project" || path == "deletion_policy" {
			continue
		}
		// Skip non-empty container blocks (whose child fields are evaluated
		// directly), while counting empty nested blocks (`basic {}`) as
		// configurable leaf-like blocks.
		if res, isRes := sch.Elem.(*schema.Resource); isRes && res != nil && len(res.Schema) > 0 {
			continue
		}
		if !sch.Optional && !sch.Required {
			continue
		}
		hasAnyConfigurable = true
		if !sch.ForceNew && !hasForceNewAncestor(path, fullSchema) {
			return true
		}
	}
	return !hasAnyConfigurable
}

func hasForceNewAncestor(fieldPath string, fullSchema map[string]*schema.Schema) bool {
	parts := strings.Split(fieldPath, ".")
	for i := 1; i < len(parts); i++ {
		ancestor := strings.Join(parts[:i], ".")
		if sch, ok := fullSchema[ancestor]; ok && sch != nil && sch.ForceNew {
			return true
		}
	}
	return false
}

func hasSchemaDefault(field *schema.Schema) bool {
	if field == nil {
		return false
	}
	if field.DefaultFunc != nil {
		return true
	}
	return field.Default != nil && field.Default != ""
}

// isSingleValueEnum probes field.ValidateFunc to determine whether a TypeString
// attribute uses validation.StringInSlice (including MMv1's verify.ValidateEnum
// and verify.OrEmpty) with exactly one non-empty allowed value.
func isSingleValueEnum(field *schema.Schema) (ok bool) {
	if field == nil || field.Type != schema.TypeString || field.ValidateFunc == nil {
		return false
	}
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	_, errs := field.ValidateFunc(enumProbeVal, enumProbeKey)
	if len(errs) != 1 || errs[0] == nil {
		return false
	}
	inner, hasPrefix := strings.CutPrefix(errs[0].Error(), enumErrStart)
	if !hasPrefix {
		return false
	}
	inner, hasSuffix := strings.CutSuffix(inner, enumErrEnd)
	if !hasSuffix {
		return false
	}
	nonEmpty := 0
	rest := strings.TrimSpace(inner)
	for rest != "" {
		quoted, err := strconv.QuotedPrefix(rest)
		if err != nil {
			return false
		}
		val, err := strconv.Unquote(quoted)
		if err != nil {
			return false
		}
		if val != "" {
			nonEmpty++
		}
		rest = strings.TrimSpace(rest[len(quoted):])
	}
	return nonEmpty == 1
}
