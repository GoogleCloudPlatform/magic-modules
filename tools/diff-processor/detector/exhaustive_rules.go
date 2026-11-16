package detector

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func initialModifyStatus(field *schema.Schema, hasUpdate bool) TransitionResult {
	if !hasUpdate {
		return TransitionResult{Status: StatusExempt, ExemptCode: "resource_not_updatable"}
	}
	if field.ForceNew {
		return TransitionResult{Status: StatusExempt, ExemptCode: "force_new"}
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
