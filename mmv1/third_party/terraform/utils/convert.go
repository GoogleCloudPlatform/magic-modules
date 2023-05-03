package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// Convert between two types by converting to/from JSON. Intended to switch
// between multiple API versions, as they are strict supersets of one another.
// item and out are pointers to structs
func Convert(item, out interface{}) error {
	return tpgresource.Convert(item, out)
}

// When converting to a map, we can't use setOmittedFields because FieldByName
// fails. Luckily, we don't use the omitted fields anymore with generated
// resources, and this function is used to bridge from handwritten -> generated.
// Since this is a known type, we can create it inline instead of needing to
// pass an object in.
func ConvertToMap(item interface{}) (map[string]interface{}, error) {
	return tpgresource.ConvertToMap(item)
}
