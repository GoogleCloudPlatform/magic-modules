package common

import (
	"encoding/json"
	"fmt"
	"strings"

	hashicorpcty "github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
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

func NewConfig() *transport_tpg.Config {
	return &transport_tpg.Config{}
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
