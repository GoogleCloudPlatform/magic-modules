package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	hashicorpcty "github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// ParseFieldValue extracts named part from resource url.
func ParseFieldValue(url string, name string) string {
	fragments := strings.Split(url, "/")
	for ix, item := range fragments {
		if item == name && ix+1 < len(fragments) {
			return fragments[ix+1]
		}
	}
	return ""
}

/*
	ParseUrlParamValuesFromAssetName uses CaiAssetNameTemplate to parse hclData from assetName, filtering out all outputFields

template: //bigquery.googleapis.com/projects/{{project}}/datasets/{{dataset_id}}
assetName: //bigquery.googleapis.com/projects/my-project/datasets/my-dataset
hclData: [project:my-project dataset_id:my-dataset]
*/
func ParseUrlParamValuesFromAssetName(assetName, template string, outputFields map[string]struct{}, hclData map[string]any) {
	fragments := strings.Split(template, "/")
	if len(fragments) < 2 {
		// We need a field and a prefix.
		return
	}
	fields := make(map[string]string) // keys are prefixes in URI, values are names of fields
	for ix, item := range fragments[1:] {
		if trimmed, ok := strings.CutPrefix(item, "{{"); ok {
			if trimmed, ok = strings.CutSuffix(trimmed, "}}"); ok {
				fields[fragments[ix]] = trimmed // ix is relative to the subslice
			}
		}
	}
	fragments = strings.Split(assetName, "/")
	for ix, item := range fragments[:len(fragments)-1] {
		if fieldName, ok := fields[item]; ok {
			if _, isOutput := outputFields[fieldName]; !isOutput {
				hclData[fieldName] = fragments[ix+1]
			}
		}
	}
}

// DecodeJSON decodes the map object into the target struct.
func DecodeJSON(data map[string]interface{}, v interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	return nil
}

// MapToCtyValWithSchema normalizes and converts resource from untyped map format to TF JSON.
//
// Normalization is a post-processing of the output map, which does the following:
// * Converts unmarshallable "schema.Set" to marshallable counterpart.
// * Strips out properties, which are not part ofthe resource TF schema.
func MapToCtyValWithSchema(m map[string]interface{}, s map[string]*schema.Schema) (cty.Value, error) {
	m = normalizeFlattenedObj(m, s).(map[string]interface{})

	b, err := json.Marshal(&m)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error marshaling map as JSON: %v", err)
	}

	ty, err := hashicorpCtyTypeToZclconfCtyType(schema.InternalMap(s).CoreConfigSchema().ImpliedType())
	if err != nil {
		return cty.NilVal, fmt.Errorf("error casting type: %v", err)
	}
	ret, err := ctyjson.Unmarshal(b, ty)
	if err != nil {
		return cty.NilVal, fmt.Errorf("error unmarshaling JSON as cty.Value: %v", err)
	}
	return ret, nil
}

func hashicorpCtyTypeToZclconfCtyType(t hashicorpcty.Type) (cty.Type, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return cty.NilType, err
	}
	var ret cty.Type
	if err := json.Unmarshal(b, &ret); err != nil {
		return cty.NilType, err
	}
	return ret, nil
}

// normalizeFlattenedObj traverses the output map recursively, removes fields which are
// not a part of TF schema and converts unmarshallable "schema.Set" objects to arrays.
func normalizeFlattenedObj(obj interface{}, schemaPerProp map[string]*schema.Schema) interface{} {
	obj = convertToMarshallableObj(obj)

	if schemaPerProp == nil {
		// Schema for leaf nodes was already checked.
		return obj
	}

	switch obj.(type) {
	case map[string]interface{}:
		objMap := obj.(map[string]interface{})
		objMapNew := map[string]interface{}{}

		for property, propertySchema := range schemaPerProp {
			propertyValue := objMap[property]

			switch propertySchema.Elem.(type) {
			case *schema.Resource:
				objMapNew[property] = normalizeFlattenedObj(propertyValue, propertySchema.Elem.(*schema.Resource).Schema)
			case *schema.ValueType:
			default:
				objMapNew[property] = normalizeFlattenedObj(propertyValue, nil)
			}
		}
		return objMapNew
	case []interface{}:
		arr := obj.([]interface{})
		arrNew := make([]interface{}, len(arr))

		for i := range arr {
			arrNew[i] = normalizeFlattenedObj(arr[i], schemaPerProp)
		}

		return arrNew
	default:
		return obj
	}
}

func convertToMarshallableObj(node interface{}) interface{} {
	switch node.(type) {
	case *schema.Set:
		nodeSet := node.(*schema.Set)

		return nodeSet.List()
	default:
		return node
	}
}
