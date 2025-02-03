package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
)

const relativeMMV1Path = "../../mmv1"

func InitializeResource(resourcePath string) *api.Resource {
	// Get the product directory from the resource path
	productDir := filepath.Dir(resourcePath)
	productYamlPath := path.Join(relativeMMV1Path, productDir, "product.yaml")

	// Initialize and compile the product
	productApi := &api.Product{}
	api.Compile(productYamlPath, productApi, "")

	// Initialize and compile the resource
	resource := &api.Resource{}
	absResourcePath := path.Join(relativeMMV1Path, resourcePath)
	api.Compile(absResourcePath, resource, "")

	// Set source file for reference
	resource.SourceYamlFile = absResourcePath

	// Set up the resource within the product context
	resource.Properties = resource.AddLabelsRelatedFields(resource.PropertiesWithExcluded(), nil)
	resource.TargetVersionName = "beta"
	resource.SetDefault(productApi)

	// Validate the resource
	resource.Validate()
	return resource
}

func updateYamlFile(filePath string, regions []string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string

	sweepersFound := false
	indent := "  " // Default indentation

	// First pass - check for existing sweepers
	for _, line := range lines {
		if strings.TrimSpace(line) == "sweepers:" {
			sweepersFound = true
			indent = strings.Repeat(" ", len(line)-len(strings.TrimLeft(line, " ")))
			break
		}
	}

	// Second pass - build new content
	inSweepers := false
	inRegions := false
	addedNewBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if we need to insert new sweepers block
		if !sweepersFound && !addedNewBlock {
			if trimmed == "examples:" || trimmed == "parameters:" || trimmed == "properties:" {
				// Add the new sweepers block without trailing newline
				newLines = append(newLines, "sweeper:")
				newLines = append(newLines, "  regions:")
				for _, region := range regions {
					newLines = append(newLines, fmt.Sprintf("    - \"%s\"", region))
				}
				addedNewBlock = true
			}
		}

		// Handle existing sweepers section
		if trimmed == "sweepers:" {
			inSweepers = true
			newLines = append(newLines, line)
			continue
		} else if inSweepers && trimmed == "regions:" {
			inRegions = true
			newLines = append(newLines, line)
			// Add all regions
			for _, region := range regions {
				newLines = append(newLines, fmt.Sprintf("%s- \"%s\"", indent+"  ", region))
			}
			continue
		} else if inRegions {
			if strings.HasPrefix(trimmed, "- ") {
				continue // Skip existing regions
			}
			inRegions = false
		}

		newLines = append(newLines, line)
	}

	// Write the file
	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}
