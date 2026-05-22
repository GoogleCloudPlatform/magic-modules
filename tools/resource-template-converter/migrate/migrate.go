package migrate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
	"gopkg.in/yaml.v3"
)

// PATCH-START: existing structs
// PATCH-END: existing structs

func MigrateFile(filePath, serviceName string) error {
	originalBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	contentStr := string(originalBytes)

	headerRegex := regexp.MustCompile(`^(#.*\n)*\n---\n`)
	headerMatch := headerRegex.FindString(contentStr)

	yamlContent := contentStr
	if headerMatch != "" {
		yamlContent = strings.TrimPrefix(contentStr, headerMatch)
	}

	var rootNode yaml.Node
	if err := yaml.Unmarshal([]byte(yamlContent), &rootNode); err != nil {
		return fmt.Errorf("failed to unmarshal YAML from %s: %w", filePath, err)
	}

	if len(rootNode.Content) == 0 || rootNode.Content[0].Kind != yaml.MappingNode {
		return fmt.Errorf("expected root to be a mapping node")
	}
	resourceMapNode := rootNode.Content[0]

	var examplesNode *yaml.Node
	var samplesNode *yaml.Node
	var iamPolicyNode *yaml.Node
	var samplesNodeIndex = -1

	for i := 0; i < len(resourceMapNode.Content); i += 2 {
		keyNode := resourceMapNode.Content[i]
		valueNode := resourceMapNode.Content[i+1]

		switch keyNode.Value {
		case "examples":
			examplesNode = valueNode
		case "samples":
			samplesNode = valueNode
			samplesNodeIndex = i
		case "iam_policy":
			iamPolicyNode = valueNode
		}
	}

	madeChanges := false

	// Transformation 1: `examples` -> `samples`
	var newSamplesNodes []*yaml.Node
	if examplesNode != nil {
		var oldExamples []*resource.Examples
		if err := examplesNode.Decode(&oldExamples); err != nil {
			return fmt.Errorf("failed to decode examples: %w", err)
		}

		if len(oldExamples) > 0 {
			migratedSamples := transformExamplesToSamples(oldExamples, filePath, serviceName)
			for _, sample := range migratedSamples {
				sampleBytes, err := yaml.Marshal(sample)
				if err != nil {
					return fmt.Errorf("failed to marshal migrated sample: %w", err)
				}
				var sampleNode yaml.Node
				if err := yaml.Unmarshal(sampleBytes, &sampleNode); err != nil {
					return fmt.Errorf("failed to unmarshal sample back to node: %w", err)
				}
				// The result of unmarshaling a map is a DocumentNode whose first content is the MappingNode
				if len(sampleNode.Content) > 0 {
					newSamplesNodes = append(newSamplesNodes, sampleNode.Content[0])
				}
			}
			madeChanges = true
		}
	}

	if madeChanges { // This means examples were found and processed
		if samplesNodeIndex != -1 { // samples block EXISTS
			newContent := []*yaml.Node{}
			for i := 0; i < len(resourceMapNode.Content); i += 2 {
				keyNode := resourceMapNode.Content[i]
				if keyNode.Value == "examples" {
					continue
				}
				if keyNode.Value == "samples" {
					newContent = append(newContent, keyNode) // samples key
					combinedSamplesNode := yaml.Node{Kind: yaml.SequenceNode}
					combinedSamplesNode.Content = append(combinedSamplesNode.Content, newSamplesNodes...)
					combinedSamplesNode.Content = append(combinedSamplesNode.Content, samplesNode.Content...)
					newContent = append(newContent, &combinedSamplesNode)
				} else {
					newContent = append(newContent, keyNode, resourceMapNode.Content[i+1])
				}
			}
			resourceMapNode.Content = newContent
		} else { // samples block DOES NOT EXIST
			samplesKeyNode := yaml.Node{Kind: yaml.ScalarNode, Value: "samples"}
			combinedSamplesNode := yaml.Node{Kind: yaml.SequenceNode, Content: newSamplesNodes}
			tempContent := []*yaml.Node{}
			for i := 0; i < len(resourceMapNode.Content); i += 2 {
				keyNode := resourceMapNode.Content[i]
				if keyNode.Value == "examples" {
					tempContent = append(tempContent, &samplesKeyNode, &combinedSamplesNode)
				} else {
					tempContent = append(tempContent, keyNode, resourceMapNode.Content[i+1])
				}
			}
			resourceMapNode.Content = tempContent
		}
	}

	// Transformation 2: Update `iam_policy` paths
	if iamPolicyNode != nil && iamPolicyNode.Kind == yaml.MappingNode {
		for i := 0; i < len(iamPolicyNode.Content); i += 2 {
			keyNode := iamPolicyNode.Content[i]
			if keyNode.Value == "example_config_body" {
				keyNode.Value = "sample_config_body" // Rename the key
				valueNode := iamPolicyNode.Content[i+1]
				oldIAMPath := "templates/terraform/iam/iam_attributes.go.tmpl"
				newIAMPath := "templates/terraform/iam/iam_attributes_sample.go.tmpl"
				if strings.Contains(valueNode.Value, oldIAMPath) {
					valueNode.Value = strings.Replace(valueNode.Value, oldIAMPath, newIAMPath, 1)
				}
				madeChanges = true
				break
			}
		}
	}

	sortMappingNode(resourceMapNode, rootKeyOrder)

	stripStringQuotes(&rootNode)

	var buf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buf)
	yamlEncoder.SetIndent(2) // v3 supports SetIndent
	if err := yamlEncoder.Encode(&rootNode); err != nil {
		return fmt.Errorf("failed to marshal updated root node for %s: %w", filePath, err)
	}

	newYAMLContent := buf.String()
	
	var finalContent string
	if headerMatch != "" {
		if !strings.Contains(headerMatch, "---") {
			finalContent = headerMatch + "---\n" + newYAMLContent
		} else {
			finalContent = headerMatch + newYAMLContent
		}
	} else {
		finalContent = "---\n" + newYAMLContent
	}

	newlineRegex := regexp.MustCompile(`\n{3,}`)
	finalContent = newlineRegex.ReplaceAllString(finalContent, "\n\n")

	finalContent = strings.TrimRight(finalContent, " \t\r\n") + "\n"

	if err := ioutil.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated file %s: %w", filePath, err)
	}

	return nil
}

// ... (Keep the existing transformExamplesToSamples function) ...
// PATCH-START: transformExamplesToSamples
func transformExamplesToSamples(oldExamples []*resource.Examples, filePath, serviceName string) []*resource.Sample {
	newSamples := make([]*resource.Sample, len(oldExamples))
	for i, old := range oldExamples {
		var newConfigPath string

		if serviceName != "" {
			var templateName string
			if old.ConfigPath != "" {
				templateName = filepath.Base(old.ConfigPath)
				calculatedPath := path.Join("templates/terraform/samples/services", serviceName, templateName)
				
				// If the calculated path differs from the default convention, set it explicitly
				defaultPath := path.Join("templates/terraform/samples/services", serviceName, fmt.Sprintf("%s.tf.tmpl", old.Name))
				if calculatedPath != defaultPath {
					newConfigPath = calculatedPath
				}
			}
		} else {
			newConfigPath = old.ConfigPath
		}
		steps := []*resource.Step{
			{
				Name:              old.Name,
				ConfigPath:        newConfigPath,
				ResourceIdVars:    old.Vars,
				TestEnvVars:       old.TestEnvVars,
				TestVarsOverrides: old.TestVarsOverrides,
				OicsVarsOverrides: old.OicsVarsOverrides,
				MinVersion:        old.MinVersion,
				IgnoreReadExtra:   old.IgnoreReadExtra,
				ExcludeImportTest: old.ExcludeImportTest,
			},
		}
		newSamples[i] = &resource.Sample{
			Name:                old.Name,
			SkipVcr:             old.SkipVcr,
			SkipTest:            old.SkipTest,
			SkipFunc:            old.SkipFunc,
			ExcludeTest:         old.ExcludeTest,
			ExcludeBasicDoc:     old.ExcludeDocs,
			ExternalProviders:   old.ExternalProviders,
			BootstrapIam:        old.BootstrapIam,
			MinVersion:          old.MinVersion,
			PrimaryResourceId:   old.PrimaryResourceId,
			PrimaryResourceType: old.PrimaryResourceType,
			RegionOverride:      old.RegionOverride,
			TGCSkipTest:         old.TGCSkipTest,
			Steps:               steps,
		}
	}
	return newSamples
}

// PATCH-END: transformExamplesToSamples

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
