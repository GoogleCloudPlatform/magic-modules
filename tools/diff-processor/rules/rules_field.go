package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// FieldRule provides structure for rules
// regarding field attribute changes
type FieldRule struct {
	name        string
	definition  string
	message     string
	identifier  string
	isRuleBreak func(old, new *schema.Schema, mc MessageContext) *BreakingChange
}

// FieldRules is a list of FieldRule
// guarding against provider breaking changes
var FieldRules = []FieldRule{
	fieldRule_ChangingType,
	fieldRule_BecomingRequired,
	fieldRule_BecomingComputedOnly,
	fieldRule_OptionalComputedToOptional,
	fieldRule_DefaultModification,
	fieldRule_GrowingMin,
	fieldRule_ShrinkingMax,
	fieldRule_RemovingDiffSuppress,
	fieldRule_AddingSubfieldToConfigModeAttr,
	fieldRule_ChangingFieldDataFormat,
}

var fieldRule_ChangingFieldDataFormat = FieldRule{
	name:       "Changing field data format",
	definition: "Modification of the data format (either by the API or manually) will cause a diff in subsequent plans if that field is not Computed. This results in a breakage. API breaking changes are out of scope with respect to provider responsibility but we may make changes in response to API breakages in some instances to provide more customer stability.",
	identifier: "field-changing-data-format",
}

var fieldRule_ChangingType = FieldRule{
	name:        "Changing Field Type",
	definition:  "While certain Field Type migrations may be supported at a technical level, it's a practice that we highly discourage. We see little value for these transitions vs the risk they impose.",
	message:     "Field {{field}} changed from {{oldType}} to {{newType}} on {{resource}}",
	identifier:  "field-changing-type",
	isRuleBreak: fieldRule_ChangingType_func,
}

func fieldRule_ChangingType_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// Type change doesn't matter for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	if old.Type != new.Type {
		oldType := getValueType(old.Type)
		newType := getValueType(new.Type)
		message = strings.ReplaceAll(message, "{{oldType}}", oldType)
		message = strings.ReplaceAll(message, "{{newType}}", newType)
		return populateMessageContext(message, mc)
	}

	oldCasted, _ := old.Elem.(*schema.Schema)
	newCasted, _ := new.Elem.(*schema.Schema)
	if oldCasted != nil && newCasted != nil && oldCasted.Type != newCasted.Type {
		oldType := getValueType(old.Type) + "." + getValueType(oldCasted.Type)
		newType := getValueType(new.Type) + "." + getValueType(newCasted.Type)
		message = strings.ReplaceAll(message, "{{oldType}}", oldType)
		message = strings.ReplaceAll(message, "{{newType}}", newType)
		return populateMessageContext(message, mc)
	}

	return nil
}

var fieldRule_BecomingRequired = FieldRule{
	name:        "Field becoming Required Field",
	definition:  "A field cannot become required as existing configs may not have this field defined. Thus, breaking configs in sequential plan or applies. If you are adding Required to a field so a block won't remain empty, this can cause two issues. First if it's a singular nested field the block may gain more fields later and it's not clear whether the field is actually required so it may be misinterpreted by future contributors. Second if users are defining empty blocks in existing configurations this change will break them. Consider these points in admittance of this type of change.",
	message:     "Field {{field}} changed from optional to required on {{resource}}",
	identifier:  "field-optional-to-required",
	isRuleBreak: fieldRule_BecomingRequired_func,
}

func fieldRule_BecomingRequired_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// Ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	if !old.Required && new.Required {
		return populateMessageContext(message, mc)
	}

	return nil
}

var fieldRule_BecomingComputedOnly = FieldRule{
	name:        "Becoming a Computed only Field",
	definition:  "While a field can go from Optional to Optional+Computed it cannot go from Required or Optional to only Computed. This transition would effectively make the field read-only thus breaking configs in sequential plan or applies where this field is defined in a configuration.",
	message:     "Field {{field}} became Computed only on {{resource}}",
	identifier:  "field-becoming-computed",
	isRuleBreak: fieldRule_BecomingComputedOnly_func,
}

func fieldRule_BecomingComputedOnly_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	// if the field is computed only already
	// this rule doesn't apply
	if old.Computed && !old.Optional {
		return nil
	}

	if new.Computed && !new.Optional {
		return populateMessageContext(message, mc)
	}
	return nil
}

var fieldRule_OptionalComputedToOptional = FieldRule{
	name:        "Optional and Computed to Optional",
	definition:  "A field cannot go from Computed + Optional to Optional. On a sequential `apply` the terraform state will have the previously computed value. The value won't be present in the config, thus ultimately causing a diff.",
	message:     "Field {{field}} transitioned from optional+computed to optional {{resource}}",
	identifier:  "field-oc-to-c",
	isRuleBreak: fieldRule_OptionalComputedToOptional_func,
}

func fieldRule_OptionalComputedToOptional_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	if (old.Computed && old.Optional) && (new.Optional && !new.Computed) {
		return populateMessageContext(message, mc)
	}
	return nil
}

var fieldRule_DefaultModification = FieldRule{
	name:        "Adding or Changing a Default Value",
	definition:  "Adding a default value where one was not previously declared can work in a very limited subset of scenarios but is an all around 'not good' practice to engage in. Changing a default value will absolutely cause a breakage. The mechanism of break for both scenarios is current terraform deployments now gain a diff with sequential applies where the diff is the new or changed default value.",
	message:     "Field {{field}} default value changed from {{oldDefault}} to {{newDefault}} on {{resource}}",
	identifier:  "field-changing-default-value",
	isRuleBreak: fieldRule_DefaultModification_func,
}

func fieldRule_DefaultModification_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	if old.Default != new.Default {
		oldDefault := fmt.Sprintf("%v", old.Default)
		newDefault := fmt.Sprintf("%v", new.Default)
		message = strings.ReplaceAll(message, "{{oldDefault}}", oldDefault)
		message = strings.ReplaceAll(message, "{{newDefault}}", newDefault)
		return populateMessageContext(message, mc)
	}

	return nil
}

var fieldRule_GrowingMin = FieldRule{
	name:        "Growing Minimum Items",
	definition:  "MinItems cannot grow. Otherwise existing terraform configurations that don't satisfy this rule will break.",
	message:     "Field {{field}} MinItems went from {{oldMin}} to {{newMin}} on {{resource}}",
	identifier:  "field-growing-min",
	isRuleBreak: fieldRule_GrowingMin_func,
}

func fieldRule_GrowingMin_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	if old.MinItems < new.MinItems || old.MinItems == 0 && new.MinItems > 0 {
		oldMin := fmt.Sprint(old.MinItems)
		if old.MinItems == 0 {
			oldMin = "unset"
		}
		newMin := fmt.Sprint(new.MinItems)
		message = strings.ReplaceAll(message, "{{oldMin}}", oldMin)
		message = strings.ReplaceAll(message, "{{newMin}}", newMin)
		return populateMessageContext(message, mc)
	}
	return nil
}

var fieldRule_ShrinkingMax = FieldRule{
	name:        "Shrinking Maximum Items",
	definition:  "MaxItems cannot shrink. Otherwise existing terraform configurations that don't satisfy this rule will break.",
	message:     "Field {{field}} MinItems went from {{oldMax}} to {{newMax}} on {{resource}}",
	identifier:  "field-shrinking-max",
	isRuleBreak: fieldRule_ShrinkingMax_func,
}

func fieldRule_ShrinkingMax_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	message := mc.message
	if old.MaxItems > new.MaxItems || old.MaxItems == 0 && new.MaxItems > 0 {
		oldMax := fmt.Sprint(old.MaxItems)
		if old.MaxItems == 0 {
			oldMax = "unset"
		}
		newMax := fmt.Sprint(new.MaxItems)
		message = strings.ReplaceAll(message, "{{oldMax}}", oldMax)
		message = strings.ReplaceAll(message, "{{newMax}}", newMax)
		return populateMessageContext(message, mc)
	}
	return nil
}

var fieldRule_RemovingDiffSuppress = FieldRule{
	name:        "Removing Diff Suppress Function",
	definition:  "Diff suppress functions cannot be removed. Otherwise terraform configurations that previously had no diffs would show diffs.",
	message:     "Field {{field}} lost its diff suppress function",
	identifier:  "field-removing-diff-suppress",
	isRuleBreak: fieldRule_RemovingDiffSuppress_func,
}

func fieldRule_RemovingDiffSuppress_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	// ignore for added / removed fields
	if old == nil || new == nil {
		return nil
	}
	if old.DiffSuppressFunc != nil && new.DiffSuppressFunc == nil {
		return populateMessageContext(mc.message, mc)
	}
	return nil
}

var fieldRule_AddingSubfieldToConfigModeAttr = FieldRule{
	name:        "Adding a subfield to a SchemaConfigModeAttr field",
	definition:  "Subfields cannot be added to fields with SchemaConfigModeAttr because they will be treated as required even if optional.",
	message:     "Field {{field}} gained a subfield {{subfield}} when it has SchemaConfigModeAttr",
	identifier:  "field-adding-subfield-to-config-mode-attr",
	isRuleBreak: fieldRule_AddingSubfieldToConfigModeAttr_func,
}

func fieldRule_AddingSubfieldToConfigModeAttr_func(old, new *schema.Schema, mc MessageContext) *BreakingChange {
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
		message := mc.message
		for fieldName := range newObj.Schema {
			if _, ok := oldObj.Schema[fieldName]; !ok {
				message = strings.ReplaceAll(message, "{{subfield}}", fieldName)
				return populateMessageContext(message, mc)
			}
		}
	}
	return nil
}

func fieldRulesToRuleArray(frs []FieldRule) []Rule {
	var rules []Rule
	for _, fr := range frs {
		rules = append(rules, fr)
	}
	return rules
}

// Name - a human readable name for the rule
func (fr FieldRule) Name() string {
	return fr.name
}

// Definition - a definition for the rule
func (fr FieldRule) Definition() string {
	return fr.definition
}

// Identifier - a navigation oriented name for the rule
func (fr FieldRule) Identifier() string {
	return fr.identifier
}

// IsRuleBreak - compares the fields and returns
// a string defining the rule breakage if detected
func (fr FieldRule) IsRuleBreak(old, new *schema.Schema, mc MessageContext) *BreakingChange {
	if fr.isRuleBreak == nil {
		return nil
	}
	mc.identifier = fr.identifier
	mc.message = fr.message
	return fr.isRuleBreak(old, new, mc)
}

// Undetectable - informs if there are functions in place
// to detect this rule.
func (fr FieldRule) Undetectable() bool {
	return fr.isRuleBreak == nil
}
