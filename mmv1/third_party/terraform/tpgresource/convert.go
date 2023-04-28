package tpgresource

import (
	"encoding/json"
)

// When converting to a map, we can't use setOmittedFields because FieldByName
// fails. Luckily, we don't use the omitted fields anymore with generated
// resources, and this function is used to bridge from handwritten -> generated.
// Since this is a known type, we can create it inline instead of needing to
// pass an object in.
func ConvertToMap(item interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	bytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
