package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	Name                 string              `yaml:"name,omitempty"`
	ConfigPath           string              `yaml:"config_path,omitempty"`
	MinVersion           string              `yaml:"min_version,omitempty"`
	GenerateDoc          bool                `yaml:"generate_doc,omitempty"`
	PrefixedVars         map[string]string   `yaml:"prefixed_vars,omitempty"`
	Vars                 map[string]string   `yaml:"vars,omitempty"`
	TestEnvVars          map[string]string   `yaml:"test_env_vars,omitempty"`
	TestVarsOverrides    map[string]string   `yaml:"test_vars_overrides,omitempty"`
	OicsVarsOverrides    map[string]string   `yaml:"oics_vars_overrides,omitempty"`
	IgnoreReadExtra      []string            `yaml:"ignore_read_extra,omitempty"`
	ExcludeImportTest    bool                `yaml:"exclude_import_test,omitempty"`
	ExcludeDocs          bool                `yaml:"exclude_docs,omitempty"`
	DocumentationHCLText string              `yaml:"-"`
	TestHCLText          string              `yaml:"-"`
	OicsHCLText          string              `yaml:"-"`
	PrimaryResourceId    string              `yaml:"-"`
	ProductName          string              `yaml:"-"`
}

type Sample struct {
	Name                 string            `yaml:"name"`
	SkipVcr              bool              `yaml:"skip_vcr,omitempty"`
	SkipTest             string            `yaml:"skip_test,omitempty"`
	ExternalProviders    []string          `yaml:"external_providers,omitempty"`
	BootstrapIam         []IamMember       `yaml:"bootstrap_iam,omitempty"`
	MinVersion           string            `yaml:"min_version,omitempty"`
	TargetVersionName    string            `yaml:"-"`
	PrimaryResourceId    string            `yaml:"primary_resource_id"`
	PrimaryResourceType  string            `yaml:"primary_resource_type,omitempty"`
	PrimaryResourceName  string            `yaml:"primary_resource_name,omitempty"`
	ExcludeTest          bool              `yaml:"exclude_test,omitempty"`
	Steps                []Step            `yaml:"steps"`
	NewConfigFuncs       []Step            `yaml:"-"`
	RegionOverride       string            `yaml:"region_override,omitempty"`
	TGCTestIgnoreExtra   []string          `yaml:"tgc_test_ignore_extra,omitempty"`
	TGCTestIgnoreInAsset []string          `yaml:"tgc_test_ignore_in_asset,omitempty"`
	TGCSkipTest          string            `yaml:"tgc_skip_test,omitempty"`
}

// =============================================================================
// Structs for the OLD `examples` format
// =============================================================================

type OldExample struct {
	Name                 string            `yaml:"name"`
	PrimaryResourceId    string            `yaml:"primary_resource_id"`
	PrimaryResourceType  string            `yaml:"primary_resource_type,omitempty"`
	BootstrapIam         []IamMember       `yaml:"bootstrap_iam,omitempty"`
	Vars                 map[string]string `yaml:"vars"`
	TestEnvVars          map[string]string `yaml:"test_env_vars,omitempty"`
	TestVarsOverrides    map[string]string `yaml:"test_vars_overrides,omitempty"`
	OicsVarsOverrides    map[string]string `yaml:"oics_vars_overrides,omitempty"`
	MinVersion           string            `yaml:"min_version,omitempty"`
	IgnoreReadExtra      []string          `yaml:"ignore_read_extra,omitempty"`
	ExcludeTest          bool              `yaml:"exclude_test,omitempty"`
	ExcludeDocs          bool              `yaml:"exclude_docs,omitempty"`
	ExcludeImportTest    bool              `yaml:"exclude_import_test,omitempty"`
	PrimaryResourceName  string            `yaml:"primary_resource_name,omitempty"`
	RegionOverride       string            `yaml:"region_override,omitempty"`
	ConfigPath           string            `yaml:"config_path,omitempty"`
	SkipVcr              bool              `yaml:"skip_vcr,omitempty"`
	SkipTest             string            `yaml:"skip_test,omitempty"`
	ExternalProviders    []string          `yaml:"external_providers,omitempty"`
	TGCTestIgnoreExtra   []string          `yaml:"tgc_test_ignore_extra,omitempty"`
	TGCTestIgnoreInAsset []string          `yaml:"tgc_test_ignore_in_asset,omitempty"`
	TGCSkipTest          string            `yaml:"tgc_skip_test,omitempty"`
}

func main() {
	basePath := "../.."
	if _, err := os.Stat(filepath.Join(basePath, "mmv1")); os.IsNotExist(err) {
		log.Fatalf("magic-modules directory structure not found. Please ensure this tool is run from 'magic-modules/tools/example-split'.")
	}

	productsPath := filepath.Join(basePath, "mmv1", "products")
	fmt.Printf("Starting migration of YAML files in: %s\n", productsPath)

	err := filepath.Walk(productsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			if err := migrateFile(path); err != nil {
				log.Printf("Failed to migrate file %s: %v\n", path, err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the products path %q: %v\n", productsPath, err)
	}

	fmt.Println("Migration complete.")
}

func migrateFile(filePath string) error {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Step 1: Unmarshal into a generic map to robustly find the 'examples' data,
	// without relying on fragile string parsing for the data itself.
	var resourceMap map[string]interface{}
	if err := yaml.Unmarshal(yamlFile, &resourceMap); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	examplesData, ok := resourceMap["examples"]
	if !ok {
		// No 'examples' block to migrate, so we skip this file.
		return nil
	}

	// Step 2: Transform the extracted data
	examplesBytes, err := yaml.Marshal(examplesData)
	if err != nil {
		return fmt.Errorf("failed to re-marshal examples block: %w", err)
	}
	var oldExamples []OldExample
	if err := yaml.Unmarshal(examplesBytes, &oldExamples); err != nil {
		return fmt.Errorf("failed to unmarshal examples into structured format: %w", err)
	}
	newSamples := transformExamplesToSamples(oldExamples)

	// Step 3: Generate the new 'samples' block as a string.
	newSamplesBytes, err := yaml.Marshal(newSamples)
	if err != nil {
		return fmt.Errorf("failed to marshal new samples data: %w", err)
	}

	// Step 4: Perform a textual replacement on the original file content
	// to preserve all formatting and comments outside the 'examples' block.
	contentStr := string(yamlFile)
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

	if startLineIndex == -1 {
		return fmt.Errorf("consistency error: could not find 'examples:' line for textual replacement")
	}

	// Determine the indentation for the content of the block.
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

	// Find the end of the block to be replaced.
	endLineIndex := startLineIndex + 1
	for ; endLineIndex < len(lines); endLineIndex++ {
		line := lines[endLineIndex]
		if strings.TrimSpace(line) == "" {
			continue // Skip empty lines
		}
		indentation := len(line) - len(strings.TrimLeft(line, " \t"))
		if indentation <= len(initialIndent) {
			break
		}
	}

	// Prepare the new content lines, correctly indented.
	newSamplesStr := string(newSamplesBytes)
	newSamplesContentLines := strings.Split(strings.TrimRight(newSamplesStr, "\n"), "\n")
	var newBlockLines []string
	newBlockLines = append(newBlockLines, initialIndent+"samples:")
	for _, line := range newSamplesContentLines {
		newBlockLines = append(newBlockLines, contentIndent+line)
	}

	// Stitch the file back together.
	var finalLines []string
	finalLines = append(finalLines, lines[:startLineIndex]...)
	finalLines = append(finalLines, newBlockLines...)
	if endLineIndex < len(lines) {
		finalLines = append(finalLines, lines[endLineIndex:]...)
	}

	outputContent := strings.Join(finalLines, "\n")
	// Preserve trailing newline if it existed.
	if strings.HasSuffix(contentStr, "\n") && !strings.HasSuffix(outputContent, "\n") {
		outputContent += "\n"
	}

	// Write the new content back to the original file.
	if err := ioutil.WriteFile(filePath, []byte(outputContent), 0644); err != nil {
		return fmt.Errorf("failed to write updated file: %w", err)
	}

	fmt.Printf("Migrated %s\n", filePath)
	return nil
}

// transformExamplesToSamples converts the old structure to the new one.
func transformExamplesToSamples(oldExamples []OldExample) []Sample {
	newSamples := make([]Sample, len(oldExamples))
	for i, old := range oldExamples {
		// Per the new logic, all fields from the old 'vars' block are moved
		// directly into the new 'prefixed_vars' block. The step-level 'vars'
		// and 'min_version' are ignored.
		newSamples[i] = Sample{
			Name:                 old.Name,
			SkipVcr:              old.SkipVcr,
			SkipTest:             old.SkipTest,
			ExternalProviders:    old.ExternalProviders,
			BootstrapIam:         old.BootstrapIam,
			MinVersion:           old.MinVersion,
			PrimaryResourceId:    old.PrimaryResourceId,
			PrimaryResourceType:  old.PrimaryResourceType,
			PrimaryResourceName:  old.PrimaryResourceName,
			ExcludeTest:          old.ExcludeTest,
			RegionOverride:       old.RegionOverride,
			TGCTestIgnoreExtra:   old.TGCTestIgnoreExtra,
			TGCTestIgnoreInAsset: old.TGCTestIgnoreInAsset,
			TGCSkipTest:          old.TGCSkipTest,
			Steps: []Step{
				{
					Name:              old.Name,
					ConfigPath:        old.ConfigPath,
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

