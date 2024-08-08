package breaking_changes

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
)

// FieldDiffRule provides structure for rules
// regarding field attribute changes
type FieldDiffRule struct {
	Identifier string
	Messages func(resource, field string, fieldDiff diff.FieldDiff) []string
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

func FieldChangingTypeMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// Type change doesn't matter for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	tmpl := "Field `%s` changed from %s to %s on `%s`"
	if fieldDiff.Old.Type != fieldDiff.New.Type {
		oldType := getValueType(fieldDiff.Old.Type)
		newType := getValueType(fieldDiff.New.Type)
		return []string{fmt.Sprintf(tmpl, field, oldType, newType, resource)}
	}

	oldCasted, _ := fieldDiff.Old.Elem.(*schema.Schema)
	newCasted, _ := fieldDiff.New.Elem.(*schema.Schema)
	if oldCasted != nil && newCasted != nil && oldCasted.Type != newCasted.Type {
		oldType := getValueType(fieldDiff.Old.Type) + "." + getValueType(oldCasted.Type)
		newType := getValueType(fieldDiff.New.Type) + "." + getValueType(newCasted.Type)
		return []string{fmt.Sprintf(tmpl, field, oldType, newType, resource)}
	}

	return nil
}

var FieldBecomingRequired = FieldDiffRule{
	Identifier: "field-optional-to-required",
	Messages:   FieldBecomingRequiredMessages,
}

func FieldBecomingRequiredMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// Ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	tmpl := "Field `%s` changed from optional to required on `%s`"
	if !fieldDiff.Old.Required && fieldDiff.New.Required {
		return []string{fmt.Sprintf(tmpl, field, resource)}
	}

	return nil
}

var FieldBecomingComputedOnly = FieldDiffRule{
	Identifier: "field-becoming-computed",
	Messages:   FieldBecomingComputedOnlyMessages,
}

func FieldBecomingComputedOnlyMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	// if the field is computed only already
	// this rule doesn't apply
	if fieldDiff.Old.Computed && !fieldDiff.Old.Optional {
		return nil
	}

	tmpl := "Field `%s` became Computed only on `%s`"
	if fieldDiff.New.Computed && !fieldDiff.New.Optional {
		return []string{fmt.Sprintf(tmpl, field, resource)}
	}
	return nil
}

var FieldOptionalComputedToOptional = FieldDiffRule{
	Identifier: "field-oc-to-c",
	Messages:   FieldOptionalComputedToOptionalMessages,
}

func FieldOptionalComputedToOptionalMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	tmpl := "Field `%s` transitioned from optional+computed to optional `%s`"
	if (fieldDiff.Old.Computed && fieldDiff.Old.Optional) && (fieldDiff.New.Optional && !fieldDiff.New.Computed) {
		return []string{fmt.Sprintf(tmpl, field, resource)}
	}
	return nil
}

var FieldDefaultModification = FieldDiffRule{
	Identifier: "field-changing-default-value",
	Messages:   FieldDefaultModificationMessages,
}

func FieldDefaultModificationMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	tmpl := "Field `%s` default value changed from %s to %s on `%s`"
	if fieldDiff.Old.Default != fieldDiff.New.Default {
		oldDefault := fmt.Sprintf("%v", fieldDiff.Old.Default)
		newDefault := fmt.Sprintf("%v", fieldDiff.New.Default)
		return []string{fmt.Sprintf(tmpl, field, oldDefault, newDefault, resource)}
	}

	return nil
}

var FieldGrowingMin = FieldDiffRule{
	Identifier: "field-growing-min",
	Messages:   FieldGrowingMinMessages,
}

func FieldGrowingMinMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	tmpl := "Field `%s` MinItems went from %s to %s on `%s`"
	if fieldDiff.Old.MinItems < fieldDiff.New.MinItems || fieldDiff.Old.MinItems == 0 && fieldDiff.New.MinItems > 0 {
		oldMin := strconv.Itoa(fieldDiff.Old.MinItems)
		if fieldDiff.Old.MinItems == 0 {
			oldMin = "unset"
		}
		newMin := strconv.Itoa(fieldDiff.New.MinItems)
		return []string{fmt.Sprintf(tmpl, field, oldMin, newMin, resource)}
	}
	return nil
}

var FieldShrinkingMax = FieldDiffRule{
	Identifier: "field-shrinking-max",
	Messages:   FieldShrinkingMaxMessages,
}

func FieldShrinkingMaxMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	tmpl := "Field `%s` MinItems went from %s to %s on `%s`"
	if fieldDiff.Old.MaxItems > fieldDiff.New.MaxItems || fieldDiff.Old.MaxItems == 0 && fieldDiff.New.MaxItems > 0 {
		oldMax := strconv.Itoa(fieldDiff.Old.MaxItems)
		if fieldDiff.Old.MaxItems == 0 {
			oldMax = "unset"
		}
		newMax := strconv.Itoa(fieldDiff.New.MaxItems)
		return []string{fmt.Sprintf(tmpl, field, oldMax, newMax, resource)}
	}
	return nil
}

var FieldRemovingDiffSuppress = FieldDiffRule{
	Identifier: "field-removing-diff-suppress",
	Messages:   FieldRemovingDiffSuppressMessages,
}

func FieldRemovingDiffSuppressMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	// ignore for added / removed fields
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	// TODO: Add resource to this message
	tmpl := "Field `%s` lost its diff suppress function"
	if fieldDiff.Old.DiffSuppressFunc != nil && fieldDiff.New.DiffSuppressFunc == nil {
		return []string{fmt.Sprintf(tmpl, field)}
	}
	return nil
}

var FieldAddingSubfieldToConfigModeAttr = FieldDiffRule{
	Identifier: "field-adding-subfield-to-config-mode-attr",
	Messages:   FieldAddingSubfieldToConfigModeAttrMessages,
}

func FieldAddingSubfieldToConfigModeAttrMessages(resource, field string, fieldDiff diff.FieldDiff) []string {
	if fieldDiff.Old == nil || fieldDiff.New == nil {
		return nil
	}
	if fieldDiff.New.ConfigMode == schema.SchemaConfigModeAttr {
		newObj, ok := fieldDiff.New.Elem.(*schema.Resource)
		if !ok {
			return nil
		}
		oldObj, ok := fieldDiff.Old.Elem.(*schema.Resource)
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
