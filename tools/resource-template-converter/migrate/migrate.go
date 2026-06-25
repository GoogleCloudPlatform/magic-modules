package migrate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// PATCH-START: existing structs
// PATCH-END: existing structs

func MigrateFile(filePath, serviceName string, onlyMigration, onlyFormat bool) error {
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

	if !onlyFormat {
		// Transformation 1: `examples` -> `samples`
		var newSamplesNodes []*yaml.Node
		if examplesNode != nil && examplesNode.Kind == yaml.SequenceNode {
			for _, exampleMapNode := range examplesNode.Content {
				if exampleMapNode.Kind != yaml.MappingNode {
					continue
				}

				var nameVal, configPathVal string
				for i := 0; i < len(exampleMapNode.Content); i += 2 {
					k := exampleMapNode.Content[i].Value
					v := exampleMapNode.Content[i+1].Value
					if k == "name" {
						nameVal = v
					} else if k == "config_path" {
						configPathVal = v
					}
				}

				var newConfigPath string
				if configPathVal != "" && serviceName != "" {
					templateName := filepath.Base(configPathVal)
					calculatedPath := path.Join("templates/terraform/samples/services", serviceName, templateName)
					defaultPath := path.Join("templates/terraform/samples/services", serviceName, fmt.Sprintf("%s.tf.tmpl", nameVal))
					if calculatedPath != defaultPath {
						newConfigPath = calculatedPath
					}
				}

				samplesContent := []*yaml.Node{}
				stepContent := []*yaml.Node{}

				for i := 0; i < len(exampleMapNode.Content); i += 2 {
					keyNode := exampleMapNode.Content[i]
					valNode := exampleMapNode.Content[i+1]

					switch keyNode.Value {
					case "name":
						samplesContent = append(samplesContent, keyNode, valNode)
						stepContent = append(stepContent, cloneNode(keyNode), cloneNode(valNode))

					case "min_version":
						samplesContent = append(samplesContent, keyNode, valNode)
						stepContent = append(stepContent, cloneNode(keyNode), cloneNode(valNode))

					case "config_path":
						if newConfigPath != "" {
							valNode.Value = newConfigPath
							valNode.Style = 0
							valNode.Tag = "!!str"
							stepContent = append(stepContent, keyNode, valNode)
						}

					case "vars":
						keyNode.Value = "resource_id_vars"
						stepContent = append(stepContent, keyNode, valNode)

					case "test_env_vars", "test_vars_overrides", "oics_vars_overrides", "ignore_read_extra", "exclude_import_test":
						stepContent = append(stepContent, keyNode, valNode)

					case "exclude_docs":
						keyNode.Value = "exclude_basic_doc"
						samplesContent = append(samplesContent, keyNode, valNode)

					case "skip_vcr", "skip_test", "skip_func", "exclude_test", "external_providers", "bootstrap_iam", "primary_resource_id", "primary_resource_type", "region_override", "tgc_skip_test":
						samplesContent = append(samplesContent, keyNode, valNode)

					default:
						// Discard unrecognized, deprecated, or typo fields (e.g., primary_resource_name, exclude_from_docs)
						// to ensure the generated YAML strictly adheres to the target Sample/Step structs.
					}
				}

				// Construct Step Mapping Node
				stepMapNode := &yaml.Node{
					Kind:    yaml.MappingNode,
					Tag:     "!!map",
					Content: stepContent,
				}

				// Construct Steps Sequence Node
				stepsSeqNode := &yaml.Node{
					Kind:    yaml.SequenceNode,
					Tag:     "!!seq",
					Content: []*yaml.Node{stepMapNode},
				}

				// Add Steps to Sample
				stepsKeyNode := &yaml.Node{
					Kind:  yaml.ScalarNode,
					Tag:   "!!str",
					Value: "steps",
				}
				samplesContent = append(samplesContent, stepsKeyNode, stepsSeqNode)

				// Construct final Sample Mapping Node
				sampleNode := &yaml.Node{
					Kind:        yaml.MappingNode,
					Tag:         "!!map",
					Content:     samplesContent,
					HeadComment: exampleMapNode.HeadComment,
					LineComment: exampleMapNode.LineComment,
					FootComment: exampleMapNode.FootComment,
				}

				newSamplesNodes = append(newSamplesNodes, sampleNode)
			}
			madeChanges = true
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
	}

	if !onlyMigration {
		sortMappingNode(resourceMapNode, rootKeyOrder)
		stripStringQuotes(&rootNode)
	}

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

func cloneNode(n *yaml.Node) *yaml.Node {
	if n == nil {
		return nil
	}
	res := *n
	if len(n.Content) > 0 {
		res.Content = make([]*yaml.Node, len(n.Content))
		for i, child := range n.Content {
			res.Content[i] = cloneNode(child)
		}
	}
	return &res
}
