package migrate

import (
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"gopkg.in/yaml.v3"
)

func stripStringQuotes(node *yaml.Node) {
	if node.Kind == yaml.ScalarNode && node.Tag == "!!str" {
		if isKeyword(node.Value) {
			node.Style = yaml.DoubleQuotedStyle
		} else if isNumber(node.Value) {
			if node.Style == 0 {
				node.Style = yaml.DoubleQuotedStyle
			}
		} else {
			node.Style = 0
		}
	}

	for _, child := range node.Content {
		stripStringQuotes(child)
	}
}

func isKeyword(val string) bool {
	lower := strings.ToLower(val)
	switch lower {
	case "true", "false", "null", "yes", "no", "on", "off", "y", "n":
		return true
	}
	return false
}

func isNumber(val string) bool {
	_, err := strconv.ParseFloat(val, 64)
	return err == nil
}

var rootKeyOrder []string
var nestedKeyOrders = make(map[string][]string)

func init() {
	// Prepend "name" at the top as the visual baseline for resource configurations
	rootKeyOrder = append([]string{"name"}, getYamlStructFieldOrder(api.Resource{})...)

	nestedKeyOrders["timeouts"] = getYamlStructFieldOrder(api.Timeouts{})
	nestedKeyOrders["sweeper"] = getYamlStructFieldOrder(resource.Sweeper{})
	nestedKeyOrders["ensure_value"] = getYamlStructFieldOrder(resource.EnsureValue{})
	nestedKeyOrders["parent"] = getYamlStructFieldOrder(resource.ParentResource{})
	nestedKeyOrders["async"] = getYamlStructFieldOrder(api.Async{})
	nestedKeyOrders["operation"] = getYamlStructFieldOrder(api.Operation{})
	nestedKeyOrders["iam_policy"] = getYamlStructFieldOrder(resource.IamPolicy{})
	nestedKeyOrders["custom_code"] = getYamlStructFieldOrder(resource.CustomCode{})
}

func sortMappingNode(mappingNode *yaml.Node, keyOrder []string) {
	if mappingNode.Kind != yaml.MappingNode {
		return
	}

	type pair struct {
		key   *yaml.Node
		value *yaml.Node
	}
	var pairs []pair
	for i := 0; i < len(mappingNode.Content); i += 2 {
		keyNode := mappingNode.Content[i]
		valNode := mappingNode.Content[i+1]

		// Recursively sort child structures
		if order, exists := nestedKeyOrders[keyNode.Value]; exists {
			sortMappingNode(valNode, order)
		} else if keyNode.Value == "operation" {
			sortMappingNode(valNode, nestedKeyOrders["operation"])
		} else if keyNode.Value == "parent" {
			sortMappingNode(valNode, nestedKeyOrders["parent"])
		} else if keyNode.Value == "ensure_value" {
			sortMappingNode(valNode, nestedKeyOrders["ensure_value"])
		}

		pairs = append(pairs, pair{key: keyNode, value: valNode})
	}

	if len(keyOrder) == 0 {
		return
	}

	orderMap := make(map[string]int)
	for idx, keyName := range keyOrder {
		orderMap[keyName] = idx
	}

	slices.SortFunc(pairs, func(a, b pair) int {
		idxA, okA := orderMap[a.key.Value]
		idxB, okB := orderMap[b.key.Value]

		if okA && okB {
			return idxA - idxB
		}
		if okA {
			return -1
		}
		if okB {
			return 1
		}
		return strings.Compare(a.key.Value, b.key.Value)
	})

	newContent := make([]*yaml.Node, 0, len(mappingNode.Content))
	for _, p := range pairs {
		newContent = append(newContent, p.key, p.value)
	}
	mappingNode.Content = newContent
}

func getYamlStructFieldOrder(structObj interface{}) []string {
	t := reflect.TypeOf(structObj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}

	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Recursively handle embedded structs (e.g., TGCResource)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			embeddedFields := getYamlStructFieldOrder(reflect.New(field.Type).Elem().Interface())
			for _, ef := range embeddedFields {
				if !slices.Contains(fields, ef) {
					fields = append(fields, ef)
				}
			}
			continue
		}

		tag := field.Tag.Get("yaml")
		var keyName string
		if tag == "" || tag == "-" {
			if tag == "-" {
				continue
			}
			// If the field is exported and has no tag, Go yaml defaults to lowercase field name
			if field.PkgPath == "" {
				keyName = strings.ToLower(field.Name)
			} else {
				continue
			}
		} else {
			parts := strings.Split(tag, ",")
			keyName = parts[0]
		}

		if keyName == "" {
			continue // Ignore tags like `,inline`
		}

		if !slices.Contains(fields, keyName) {
			fields = append(fields, keyName)
		}
	}
	return fields
}
