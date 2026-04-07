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

// ... (Keep the existing structs IamMember, Step, Sample, OldExample) ...
// PATCH-START: existing structs
type IamMember struct {
	Member string `yaml:"member"`
	Role   string `yaml:"role"`
}

type Step struct {
	Name                  string            `yaml:"name,omitempty"`
	ConfigPath            string            `yaml:"config_path,omitempty"`
	MinVersion            string            `yaml:"min_version,omitempty"`
	GenerateDoc           bool              `yaml:"generate_doc,omitempty"`
	ResourceIdVars        map[string]string `yaml:"resource_id_vars,omitempty"`
	Vars                  map[string]string `yaml:"vars,omitempty"`
	TestEnvVars           map[string]string `yaml:"test_env_vars,omitempty"`
	TestVarsOverrides     map[string]string `yaml:"test_vars_overrides,omitempty"`
	OicsVarsOverrides     map[string]string `yaml:"oics_vars_overrides,omitempty"`
	IgnoreReadExtra       []string          `yaml:"ignore_read_extra,omitempty"`
	ExcludeIdentityImport bool              `yaml:"exclude_identity_import,omitempty"`
	ExcludeImportTest     bool              `yaml:"exclude_import_test,omitempty"`
	IncludeStepDoc        bool              `yaml:"include_step_doc,omitempty"` // Opt-in for ANY step
	DocumentationHCLText  string            `yaml:"-"`
	TestHCLText           string            `yaml:"-"`
	OicsHCLText           string            `yaml:"-"`
	PrimaryResourceId     string            `yaml:"-"`
	ProductName           string            `yaml:"-"`
}

type Sample struct {
	Name                string      `yaml:"name"`
	SkipVcr             bool        `yaml:"skip_vcr,omitempty"`
	SkipFunc            string      `yaml:"skip_func,omitempty"`
	SkipTest            string      `yaml:"skip_test,omitempty"`
	ExternalProviders   []string    `yaml:"external_providers,omitempty"`
	BootstrapIam        []IamMember `yaml:"bootstrap_iam,omitempty"`
	MinVersion          string      `yaml:"min_version,omitempty"`
	TargetVersionName   string      `yaml:"-"`
	PrimaryResourceId   string      `yaml:"primary_resource_id"`
	PrimaryResourceType string      `yaml:"primary_resource_type,omitempty"`
	ExcludeTest         bool        `yaml:"exclude_test,omitempty"`
	ExcludeBasicDoc     bool        `yaml:"exclude_basic_doc,omitempty"` // Opt-out for FIRST step
	Steps               []Step      `yaml:"steps"`
	NewConfigFuncs      []Step      `yaml:"-"`
	RegionOverride      string      `yaml:"region_override,omitempty"`
	TGCSkipTest         string      `yaml:"tgc_skip_test,omitempty"`
}

type OldExample struct {
	Name                  string            `yaml:"name"`
	PrimaryResourceId     string            `yaml:"primary_resource_id"`
	PrimaryResourceType   string            `yaml:"primary_resource_type,omitempty"`
	BootstrapIam          []IamMember       `yaml:"bootstrap_iam,omitempty"`
	Vars                  map[string]string `yaml:"vars"`
	TestEnvVars           map[string]string `yaml:"test_env_vars,omitempty"`
	TestVarsOverrides     map[string]string `yaml:"test_vars_overrides,omitempty"`
	OicsVarsOverrides     map[string]string `yaml:"oics_vars_overrides,omitempty"`
	MinVersion            string            `yaml:"min_version,omitempty"`
	IgnoreReadExtra       []string          `yaml:"ignore_read_extra,omitempty"`
	ExcludeTest           bool              `yaml:"exclude_test,omitempty"`
	ExcludeDocs           bool              `yaml:"exclude_docs,omitempty"`
	ExcludeIdentityImport bool              `yaml:"exclude_identity_import,omitempty"`
	ExcludeImportTest     bool              `yaml:"exclude_import_test,omitempty"`
	RegionOverride        string            `yaml:"region_override,omitempty"`
	ConfigPath            string            `yaml:"config_path,omitempty"`
	SkipVcr               bool              `yaml:"skip_vcr,omitempty"`
	SkipFunc              string            `yaml:"skip_func,omitempty"`
	SkipTest              string            `yaml:"skip_test,omitempty"`
	ExternalProviders     []string          `yaml:"external_providers,omitempty"`
	TGCSkipTest           string            `yaml:"tgc_skip_test,omitempty"`
}

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
		var oldExamples []OldExample
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

	if madeChanges {
		var buf bytes.Buffer
		yamlEncoder := yaml.NewEncoder(&buf)
		yamlEncoder.SetIndent(2) // v3 supports SetIndent
		if err := yamlEncoder.Encode(&rootNode); err != nil {
			return fmt.Errorf("failed to marshal updated root node for %s: %w", filePath, err)
		}

		newYAMLContent := buf.String()
		finalContent := newYAMLContent
		if headerMatch != "" {
			finalContent = headerMatch + newYAMLContent
		}

		if err := ioutil.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
			return fmt.Errorf("failed to write updated file %s: %w", filePath, err)
		}
	}

	return nil
}

// ... (Keep the existing transformExamplesToSamples function) ...
// PATCH-START: transformExamplesToSamples
func transformExamplesToSamples(oldExamples []OldExample, filePath, serviceName string) []Sample {
	newSamples := make([]Sample, len(oldExamples))
	for i, old := range oldExamples {
		var newConfigPath string

		if serviceName != "" {
			var templateName string
			if old.ConfigPath != "" {
				templateName = filepath.Base(old.ConfigPath)
				newConfigPath = path.Join("templates/terraform/samples/services", serviceName, templateName)
			}
		} else {
			newConfigPath = old.ConfigPath
		}
		newSamples[i] = Sample{
			Name:                old.Name,
			SkipVcr:             old.SkipVcr,
			SkipFunc:            old.SkipFunc,
			SkipTest:            old.SkipTest,
			ExternalProviders:   old.ExternalProviders,
			BootstrapIam:        old.BootstrapIam,
			MinVersion:          old.MinVersion,
			PrimaryResourceId:   old.PrimaryResourceId,
			PrimaryResourceType: old.PrimaryResourceType,
			ExcludeBasicDoc:     old.ExcludeDocs,
			ExcludeTest:         old.ExcludeTest,
			RegionOverride:      old.RegionOverride,
			TGCSkipTest:         old.TGCSkipTest,
			Steps: []Step{
				{
					Name:                  old.Name,
					ConfigPath:            newConfigPath,
					ResourceIdVars:        old.Vars,
					TestEnvVars:           old.TestEnvVars,
					TestVarsOverrides:     old.TestVarsOverrides,
					OicsVarsOverrides:     old.OicsVarsOverrides,
					IgnoreReadExtra:       old.IgnoreReadExtra,
					ExcludeImportTest:     old.ExcludeImportTest,
					ExcludeIdentityImport: old.ExcludeIdentityImport,
				},
			},
		}
	}
	return newSamples
}

// PATCH-END: transformExamplesToSamples
