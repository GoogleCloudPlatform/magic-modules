package breaking_changes

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FieldDiffRule provides structure for rules
// regarding field attribute changes
type FieldDiffRule struct {
	Identifier string
	// TODO: change signature to take FieldDiff instead of old, new.
	Messages func(resource, field string, old, new *schema.Schema) []string
}

// FieldDiffRules is a list of FieldDiffRule
// guarding against provider breaking changes
var FieldDiffRules = []FieldDiffRule{
	FieldChangingType,
	FieldBecomingRequired,
	FieldBecomingComputedOnly,
	FieldOptionalComputedToOptional,
	FieldDefaultModification,
	FieldGrowingMin,
	FieldShrinkingMax,
	FieldRemovingDiffSuppress,
	FieldAddingSubfieldToConfigModeAttr,
}

var FieldChangingType = FieldDiffRule{
	Identifier: "field-changing-type",
	Messages:   FieldChangingTypeMessages,
}

func FieldChangingTypeMessages(resource, field string, old, new *schema.Schema) []string {
	// Type change doesn't matter for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	tmpl := "Field `%s` changed from %s to %s on `%s`"
	if old.Type != new.Type {
		oldType := getValueType(old.Type)
		newType := getValueType(new.Type)
		return []string{fmt.Sprintf(tmpl, field, oldType, newType, resource)}
	}

	oldCasted, _ := old.Elem.(*schema.Schema)
	newCasted, _ := new.Elem.(*schema.Schema)
	if oldCasted != nil && newCasted != nil && oldCasted.Type != newCasted.Type {
		oldType := getValueType(old.Type) + "." + getValueType(oldCasted.Type)
		newType := getValueType(new.Type) + "." + getValueType(newCasted.Type)
		return []string{fmt.Sprintf(tmpl, field, oldType, newType, resource)}
	}

	return nil
}

var FieldBecomingRequired = FieldDiffRule{
	Identifier: "field-optional-to-required",
	Messages:   FieldBecomingRequiredMessages,
}

func FieldBecomingRequiredMessages(resource, field string, old, new *schema.Schema) []string {
	// Ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	tmpl := "Field `%s` changed from optional to required on `%s`"
	if !old.Required && new.Required {
		return []string{fmt.Sprintf(tmpl, field, resource)}
	}

	return nil
}

var FieldBecomingComputedOnly = FieldDiffRule{
	Identifier: "field-becoming-computed",
	Messages:   FieldBecomingComputedOnlyMessages,
}

func FieldBecomingComputedOnlyMessages(resource, field string, old, new *schema.Schema) []string {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	// if the field is computed only already
	// this rule doesn't apply
	if old.Computed && !old.Optional {
		return nil
	}

	tmpl := "Field `%s` became Computed only on `%s`"
	if new.Computed && !new.Optional {
		return []string{fmt.Sprintf(tmpl, field, resource)}
	}
	return nil
}

var FieldOptionalComputedToOptional = FieldDiffRule{
	Identifier: "field-oc-to-c",
	Messages:   FieldOptionalComputedToOptionalMessages,
}

func FieldOptionalComputedToOptionalMessages(resource, field string, old, new *schema.Schema) []string {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	tmpl := "Field `%s` transitioned from optional+computed to optional `%s`"
	if (old.Computed && old.Optional) && (new.Optional && !new.Computed) {
		return []string{fmt.Sprintf(tmpl, field, resource)}
	}
	return nil
}

var FieldDefaultModification = FieldDiffRule{
	Identifier: "field-changing-default-value",
	Messages:   FieldDefaultModificationMessages,
}

func FieldDefaultModificationMessages(resource, field string, old, new *schema.Schema) []string {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	tmpl := "Field `%s` default value changed from %s to %s on `%s`"
	if old.Default != new.Default {
		oldDefault := fmt.Sprintf("%v", old.Default)
		newDefault := fmt.Sprintf("%v", new.Default)
		return []string{fmt.Sprintf(tmpl, field, oldDefault, newDefault, resource)}
	}

	return nil
}

var FieldGrowingMin = FieldDiffRule{
	Identifier: "field-growing-min",
	Messages:   FieldGrowingMinMessages,
}

func FieldGrowingMinMessages(resource, field string, old, new *schema.Schema) []string {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	tmpl := "Field `%s` MinItems went from %s to %s on `%s`"
	if old.MinItems < new.MinItems || old.MinItems == 0 && new.MinItems > 0 {
		oldMin := strconv.Itoa(old.MinItems)
		if old.MinItems == 0 {
			oldMin = "unset"
		}
		newMin := strconv.Itoa(new.MinItems)
		return []string{fmt.Sprintf(tmpl, field, oldMin, newMin, resource)}
	}
	return nil
}

var FieldShrinkingMax = FieldDiffRule{
	Identifier: "field-shrinking-max",
	Messages:   FieldShrinkingMaxMessages,
}

func FieldShrinkingMaxMessages(resource, field string, old, new *schema.Schema) []string {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	tmpl := "Field `%s` MinItems went from %s to %s on `%s`"
	if old.MaxItems > new.MaxItems || old.MaxItems == 0 && new.MaxItems > 0 {
		oldMax := strconv.Itoa(old.MaxItems)
		if old.MaxItems == 0 {
			oldMax = "unset"
		}
		newMax := strconv.Itoa(new.MaxItems)
		return []string{fmt.Sprintf(tmpl, field, oldMax, newMax, resource)}
	}
	return nil
}

var FieldRemovingDiffSuppress = FieldDiffRule{
	Identifier: "field-removing-diff-suppress",
	Messages:   FieldRemovingDiffSuppressMessages,
}

func FieldRemovingDiffSuppressMessages(resource, field string, old, new *schema.Schema) []string {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	// TODO: Add resource to this message
	tmpl := "Field `%s` lost its diff suppress function"
	if old.DiffSuppressFunc != nil && new.DiffSuppressFunc == nil {
		return []string{fmt.Sprintf(tmpl, field)}
	}
	return nil
}

var FieldAddingSubfieldToConfigModeAttr = FieldDiffRule{
	Identifier: "field-adding-subfield-to-config-mode-attr",
	Messages:   FieldAddingSubfieldToConfigModeAttrMessages,
}

func FieldAddingSubfieldToConfigModeAttrMessages(resource, field string, old, new *schema.Schema) []string {
	if old == nil || new == nil {
		return nil
	}
	if new.ConfigMode == schema.SchemaConfigModeAttr {
		newObj, ok := new.Elem.(*schema.Resource)
		if !ok {
			return nil
		}
		oldObj, ok := old.Elem.(*schema.Resource)
		if !ok {
			return nil
		}
		// TODO: Add resource to this message
		tmpl := "Field `%s` gained a subfield `%s` when it has SchemaConfigModeAttr"
		for subfield := range newObj.Schema {
			if _, ok := oldObj.Schema[subfield]; !ok {
				return []string{fmt.Sprintf(tmpl, field, subfield)}
			}
		}
	}
	return nil
}
