package migrate

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// =============================================================================
// Structs for the NEW `samples` format
// =============================================================================

type IamMember struct {
	Member string `yaml:"member"`
	Role   string `yaml:"role"`
}

type Step struct {
	Name                 string            `yaml:"name,omitempty"`
	ConfigPath           string            `yaml:"config_path,omitempty"`
	MinVersion           string            `yaml:"min_version,omitempty"`
	GenerateDoc          bool              `yaml:"generate_doc,omitempty"`
	PrefixedVars         map[string]string `yaml:"prefixed_vars,omitempty"`
	Vars                 map[string]string `yaml:"vars,omitempty"`
	TestEnvVars          map[string]string `yaml:"test_env_vars,omitempty"`
	TestVarsOverrides    map[string]string `yaml:"test_vars_overrides,omitempty"`
	OicsVarsOverrides    map[string]string `yaml:"oics_vars_overrides,omitempty"`
	IgnoreReadExtra      []string          `yaml:"ignore_read_extra,omitempty"`
	ExcludeImportTest    bool              `yaml:"exclude_import_test,omitempty"`
	ExcludeDocs          bool              `yaml:"exclude_docs,omitempty"`
	DocumentationHCLText string            `yaml:"-"`
	TestHCLText          string            `yaml:"-"`
	OicsHCLText          string            `yaml:"-"`
	PrimaryResourceId    string            `yaml:"-"`
	ProductName          string            `yaml:"-"`
}

type Sample struct {
	Name                string      `yaml:"name"`
	SkipVcr             bool        `yaml:"skip_vcr,omitempty"`
	SkipTest            string      `yaml:"skip_test,omitempty"`
	ExternalProviders   []string    `yaml:"external_providers,omitempty"`
	BootstrapIam        []IamMember `yaml:"bootstrap_iam,omitempty"`
	MinVersion          string      `yaml:"min_version,omitempty"`
	TargetVersionName   string      `yaml:"-"`
	PrimaryResourceId   string      `yaml:"primary_resource_id"`
	PrimaryResourceType string      `yaml:"primary_resource_type,omitempty"`
	PrimaryResourceName string      `yaml:"primary_resource_name,omitempty"`
	ExcludeTest         bool        `yaml:"exclude_test,omitempty"`
	Steps               []Step      `yaml:"steps"`
	NewConfigFuncs      []Step      `yaml:"-"`
	RegionOverride      string      `yaml:"region_override,omitempty"`
	TGCSkipTest         string      `yaml:"tgc_skip_test,omitempty"`
}

// =============================================================================
// Structs for the OLD `examples` format
// =============================================================================

type OldExample struct {
	Name                string            `yaml:"name"`
	PrimaryResourceId   string            `yaml:"primary_resource_id"`
	PrimaryResourceType string            `yaml:"primary_resource_type,omitempty"`
	BootstrapIam        []IamMember       `yaml:"bootstrap_iam,omitempty"`
	Vars                map[string]string `yaml:"vars"`
	TestEnvVars         map[string]string `yaml:"test_env_vars,omitempty"`
	TestVarsOverrides   map[string]string `yaml:"test_vars_overrides,omitempty"`
	OicsVarsOverrides   map[string]string `yaml:"oics_vars_overrides,omitempty"`
	MinVersion          string            `yaml:"min_version,omitempty"`
	IgnoreReadExtra     []string          `yaml:"ignore_read_extra,omitempty"`
	ExcludeTest         bool              `yaml:"exclude_test,omitempty"`
	ExcludeDocs         bool              `yaml:"exclude_docs,omitempty"`
	ExcludeImportTest   bool              `yaml:"exclude_import_test,omitempty"`
	PrimaryResourceName string            `yaml:"primary_resource_name,omitempty"`
	RegionOverride      string            `yaml:"region_override,omitempty"`
	ConfigPath          string            `yaml:"config_path,omitempty"`
	SkipVcr             bool              `yaml:"skip_vcr,omitempty"`
	SkipTest            string            `yaml:"skip_test,omitempty"`
	ExternalProviders   []string          `yaml:"external_providers,omitempty"`
	TGCSkipTest         string            `yaml:"tgc_skip_test,omitempty"`
}

func MigrateFile(filePath, serviceName string) error {
	originalBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(originalBytes)
	madeChanges := false

	// Transformation 1: `examples` -> `samples`
	var resourceMap map[string]interface{}
	if err := yaml.Unmarshal(originalBytes, &resourceMap); err != nil {
		return fmt.Errorf("failed to unmarshal YAML from %s: %w", filePath, err)
	}

	if examplesData, ok := resourceMap["examples"]; ok {
		examplesBytes, err := yaml.Marshal(examplesData)
		if err != nil {
			return fmt.Errorf("failed to re-marshal examples block from %s: %w", filePath, err)
		}
		var oldExamples []OldExample
		if err := yaml.Unmarshal(examplesBytes, &oldExamples); err != nil {
			return fmt.Errorf("failed to unmarshal examples into structured format from %s: %w", filePath, err)
		}
		newSamples := transformExamplesToSamples(oldExamples, filePath, serviceName)

		newSamplesBytes, err := yaml.Marshal(newSamples)
		if err != nil {
			return fmt.Errorf("failed to marshal new samples data for %s: %w", filePath, err)
		}

		lines := strings.Split(contentStr, "\n")
		startLineIndex := -1
		var initialIndent string
		for i, line := range lines {
			if strings.HasPrefix(strings.TrimLeft(line, " \t"), "examples:") {
				startLineIndex = i
				indentation := len(line) - len(strings.TrimLeft(line, " \t"))
				initialIndent = line[:indentation]
				break
			}
		}

		if startLineIndex != -1 {
			var contentIndent string
			for i := startLineIndex + 1; i < len(lines); i++ {
				if strings.TrimSpace(lines[i]) != "" {
					line := lines[i]
					indentation := len(line) - len(strings.TrimLeft(line, " \t"))
					contentIndent = line[:indentation]
					break
				}
			}
			if contentIndent == "" {
				contentIndent = initialIndent + "  " // Default indent if block is empty
			}

			endLineIndex := startLineIndex + 1
			for ; endLineIndex < len(lines); endLineIndex++ {
				line := lines[endLineIndex]
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
					continue
				}
				indentation := len(line) - len(strings.TrimLeft(line, " \t"))
				if indentation <= len(initialIndent) {
					break
				}
			}

			newSamplesStr := string(newSamplesBytes)
			newSamplesContentLines := strings.Split(strings.TrimRight(newSamplesStr, "\n"), "\n")
			var newBlockLines []string
			newBlockLines = append(newBlockLines, initialIndent+"samples:")
			for _, line := range newSamplesContentLines {
				newBlockLines = append(newBlockLines, contentIndent+line)
			}

			var finalLines []string
			finalLines = append(finalLines, lines[:startLineIndex]...)
			finalLines = append(finalLines, newBlockLines...)
			if endLineIndex < len(lines) {
				finalLines = append(finalLines, lines[endLineIndex:]...)
			}

			outputContent := strings.Join(finalLines, "\n")
			if strings.HasSuffix(contentStr, "\n") && !strings.HasSuffix(outputContent, "\n") {
				outputContent += "\n"
			}

			if outputContent != contentStr {
				contentStr = outputContent
				madeChanges = true
				// fmt.Printf("Migrated 'examples' block in %s\n", filePath)
			}
		}
	}

	// Transformation 2: Update `example_config_body` and `iam_policy` paths.
	// This is done with line-by-line processing to handle different replacement rules.
	lines := strings.Split(contentStr, "\n")
	iamMadeChanges := false

	oldIAMPath := "templates/terraform/iam/iam_attributes.go.tmpl"
	newIAMPath := "templates/terraform/iam/iam_attributes_sample.go.tmpl"

	oldIAMKey := "example_config_body:"
	newIAMKey := "sample_config_body:"

	for i, line := range lines {
		// Check if the line contains the key we need to change.
		if strings.Contains(line, oldIAMKey) {
			var newLine string
			// Case 1: Standard config body path. Change both key and value.
			if strings.Contains(line, oldIAMPath) {
				tempLine := strings.Replace(line, oldIAMKey, newIAMKey, 1)
				newLine = strings.Replace(tempLine, oldIAMPath, newIAMPath, 1)
				// Case 2: Custom config body path. Change only the key.
			} else {
				newLine = strings.Replace(line, oldIAMKey, newIAMKey, 1)
			}

			if newLine != line {
				lines[i] = newLine
				iamMadeChanges = true
			}
		}
	}

	if iamMadeChanges {
		newContentStr := strings.Join(lines, "\n")
		// The overall `madeChanges` flag will be set if this transform did anything.
		if newContentStr != contentStr {
			contentStr = newContentStr
			madeChanges = true
			// fmt.Printf("Updated 'iam_policy' / 'example_config_body' in %s\n", filePath)
		}
	}

	if madeChanges {
		if err := ioutil.WriteFile(filePath, []byte(contentStr), 0644); err != nil {
			return fmt.Errorf("failed to write updated file %s: %w", filePath, err)
		}
	}

	return nil
}

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
			SkipTest:            old.SkipTest,
			ExternalProviders:   old.ExternalProviders,
			BootstrapIam:        old.BootstrapIam,
			MinVersion:          old.MinVersion,
			PrimaryResourceId:   old.PrimaryResourceId,
			PrimaryResourceType: old.PrimaryResourceType,
			PrimaryResourceName: old.PrimaryResourceName,
			ExcludeTest:         old.ExcludeTest,
			RegionOverride:      old.RegionOverride,
			TGCSkipTest:         old.TGCSkipTest,
			Steps: []Step{
				{
					Name:              old.Name,
					ConfigPath:        newConfigPath,
					PrefixedVars:      old.Vars,
					TestEnvVars:       old.TestEnvVars,
					TestVarsOverrides: old.TestVarsOverrides,
					OicsVarsOverrides: old.OicsVarsOverrides,
					IgnoreReadExtra:   old.IgnoreReadExtra,
					ExcludeImportTest: old.ExcludeImportTest,
					ExcludeDocs:       old.ExcludeDocs,
				},
			},
		}
	}
	return newSamples
}
