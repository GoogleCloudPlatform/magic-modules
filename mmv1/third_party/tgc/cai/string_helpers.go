package cai

func ConvertInterfaceToStringArray(values []interface{}) []string {
	stringArray := make([]string, len(values))
	for i, v := range values {
		stringArray[i] = v.(string)
	}
	return stringArray
}
