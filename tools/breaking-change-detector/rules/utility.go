package rules

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/.ci/breaking-change-detector/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func documentationReference(version, identifier string) string {
	return fmt.Sprintf(" - [reference](%s)", constants.GetFileUrl(version, identifier))
}

func getValueType(valueType schema.ValueType) string {
	switch valueType {
	case schema.TypeBool:
		return "TypeBool"
	case schema.TypeInt:
		return "TypeInt"
	case schema.TypeFloat:
		return "TypeFloat"
	case schema.TypeString:
		return "TypeString"
	case schema.TypeList:
		return "TypeList"
	case schema.TypeMap:
		return "TypeMap"
	case schema.TypeSet:
		return "TypeSet"
	}
	return "TypeUndefined"
}
