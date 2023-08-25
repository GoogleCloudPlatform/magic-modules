package main

import (
	newProvider "google/provider/new/google"
	oldProvider "google/provider/old/google"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/.ci/breaking-change-detector/rules"
	"github.com/golang/glog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func compare() []string {
	resourceMapOld := oldProvider.ResourceMap()
	resourceMapNew := newProvider.ResourceMap()

	return compareResourceMaps(resourceMapOld, resourceMapNew)
}

func compareResourceMaps(old, new map[string]*schema.Resource) []string {
	messages := []string{}

	for _, rule := range rules.ResourceInventoryRules {
		violatingResources := rule.IsRuleBreak(old, new)
		if len(violatingResources) > 0 {
			for _, resourceName := range violatingResources {
				newMessage := rule.Message(resourceName)
				messages = append(messages, newMessage)
			}
		}

	}

	for resourceName, resource := range new {
		oldResource, ok := old[resourceName]
		if ok {
			newMessages := compareResourceSchema(resourceName, oldResource.Schema, resource.Schema)
			messages = append(messages, newMessages...)
		}
	}

	return messages
}

func compareResourceSchema(resourceName string, old, new map[string]*schema.Schema) []string {
	messages := []string{}
	oldCompressed := flattenSchema(old)
	newCompressed := flattenSchema(new)

	for _, rule := range rules.ResourceSchemaRules {
		violatingFields := rule.IsRuleBreak(oldCompressed, newCompressed)
		if len(violatingFields) > 0 {
			for _, fieldName := range violatingFields {
				newMessage := rule.Message(resourceName, fieldName)
				messages = append(messages, newMessage)
			}
		}
	}

	for fieldName, field := range newCompressed {
		oldField, ok := oldCompressed[fieldName]
		if ok {
			newMessages := compareField(resourceName, fieldName, oldField, field)
			messages = append(messages, newMessages...)
		}
	}

	return messages
}

func compareField(resourceName, fieldName string, old, new *schema.Schema) []string {
	messages := []string{}
	fieldRules := rules.FieldRules

	for _, rule := range fieldRules {
		breakageMessage := rule.IsRuleBreak(
			old,
			new,
			rules.MessageContext{
				Resource: resourceName,
				Field:    fieldName,
			},
		)
		if breakageMessage != "" {
			messages = append(messages, breakageMessage)
		}
	}
	return messages
}

func flattenSchema(schemaObj map[string]*schema.Schema) map[string]*schema.Schema {
	return flattenSchemaRecursive(nil, schemaObj)
}

func flattenSchemaRecursive(parentLineage []string, schemaObj map[string]*schema.Schema) map[string]*schema.Schema {
	compressed := make(map[string]*schema.Schema)

	// prepare prefix to bring nested entries up
	parentPrefix := strings.Join(parentLineage, ".")
	if len(parentPrefix) > 0 {
		parentPrefix += "."
	}

	// add entry to output and call
	// flattenSchemaRecursive for any children
	for fieldName, field := range schemaObj {
		compressed[parentPrefix+fieldName] = field
		casted, typeConverted := field.Elem.(*schema.Resource)
		if field.Elem != nil && typeConverted {
			newLineage := append([]string{}, parentLineage...)
			newLineage = append(newLineage, fieldName)
			compressedChild := flattenSchemaRecursive(newLineage, casted.Schema)
			compressed = mergeSchemaMaps(compressed, compressedChild)
		}
	}

	return compressed
}

func mergeSchemaMaps(map1, map2 map[string]*schema.Schema) map[string]*schema.Schema {
	merged := make(map[string]*schema.Schema)
	for key, value := range map1 {
		merged[key] = value
	}

	for key, value := range map2 {
		if _, alreadyExists := merged[key]; alreadyExists {
			glog.Errorf("error when trying to merge maps key " + key + " was found in both maps.. please ensure the children you are merging up have the prefix on the key names.")
		}
		merged[key] = value
	}

	return merged
}
