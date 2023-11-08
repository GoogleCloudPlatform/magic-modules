package rules

import (
	"fmt"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func documentationReference(identifier string) string {
	return fmt.Sprintf(" - [reference](%s)", constants.GetFileUrl(identifier))
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
