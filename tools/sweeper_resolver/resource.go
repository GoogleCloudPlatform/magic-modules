package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api/resource"
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

func shouldSkipUpdate(pairs []LocationPair, existingSweeper resource.Sweeper) bool {
	// If there are no pairs to add, skip the update
	if len(pairs) == 0 {
		return true
	}

	// For backward compatibility: Check if the only region present is us-central1
	if len(pairs) == 1 && pairs[0].Region == "us-central1" && pairs[0].Zone == "" {
		// Check if the existing sweeper has only us-central1 or is empty
		if len(existingSweeper.Regions) == 0 || (len(existingSweeper.Regions) == 1 && existingSweeper.Regions[0] == "us-central1") {
			return true
		}
	}

	// Check existing URL substitutions
	if len(existingSweeper.URLSubstitutions) > 0 {
		// Create a map for easier lookup of existing pairs
		existingPairs := make(map[string]map[string]bool)
		for _, sub := range existingSweeper.URLSubstitutions {
			if existingPairs[sub.Region] == nil {
				existingPairs[sub.Region] = make(map[string]bool)
			}
			existingPairs[sub.Region][sub.Zone] = true
		}

		// If all new pairs already exist and there are no extra pairs, skip the update
		allPairsExist := true
		for _, pair := range pairs {
			if zones, hasRegion := existingPairs[pair.Region]; !hasRegion {
				allPairsExist = false
				break
			} else if !zones[pair.Zone] {
				allPairsExist = false
				break
			}
		}

		if allPairsExist && len(pairs) == len(existingSweeper.URLSubstitutions) {
			return true
		}
	}

	return false
}

func getExistingURLSubstitutions(lines []string) []LocationPair {
	var pairs []LocationPair
	var inSweeper, inURLSubs bool
	var currentPair LocationPair

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "sweeper:" {
			inSweeper = true
			continue
		}

		if inSweeper && trimmed == "url_substitutions:" {
			inURLSubs = true
			continue
		}

		if inURLSubs {
			if strings.HasPrefix(trimmed, "- ") {
				// If we find a new item, save the previous pair if it exists
				if currentPair.Region != "" || currentPair.Zone != "" {
					pairs = append(pairs, currentPair)
					currentPair = LocationPair{}
				}
			} else if !strings.HasPrefix(trimmed, "region:") && !strings.HasPrefix(trimmed, "zone:") {
				inURLSubs = false
				if currentPair.Region != "" || currentPair.Zone != "" {
					pairs = append(pairs, currentPair)
				}
				continue
			}

			if strings.HasPrefix(trimmed, "region:") {
				currentPair.Region = strings.Trim(strings.TrimPrefix(trimmed, "region:"), " \"")
			} else if strings.HasPrefix(trimmed, "zone:") {
				currentPair.Zone = strings.Trim(strings.TrimPrefix(trimmed, "zone:"), " \"")
			}
		}
	}

	// Add the last pair if it exists
	if inURLSubs && (currentPair.Region != "" || currentPair.Zone != "") {
		pairs = append(pairs, currentPair)
	}

	return pairs
}

func isUniquePair(pair LocationPair, existing []LocationPair) bool {
	for _, p := range existing {
		if p.Region == pair.Region && p.Zone == pair.Zone {
			return false
		}
	}
	return true
}

func updateYamlFile(filePath string, locationPairs []LocationPair, existingSweeper resource.Sweeper) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	var sweepersFound bool
	baseIndent := "  " // Default indentation

	// First pass - find sweeper section and get base indentation
	for _, line := range lines {
		if strings.TrimSpace(line) == "sweeper:" {
			sweepersFound = true
			baseIndent = strings.Repeat(" ", len(line)-len(strings.TrimLeft(line, " ")))
			break
		}
	}

	// Second pass - rebuild the file
	var inSweeper bool
	var skipToNextSection bool
	var urlSubsAdded bool

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)
		indent := len(line) - len(strings.TrimLeft(line, " "))

		// Handle new sweeper block insertion before major sections
		if !sweepersFound && (trimmed == "examples:" || trimmed == "parameters:" || trimmed == "properties:") {
			// Add sweeper block before major section
			// Check if we're after custom_code
			if len(newLines) > 0 {
				lastLine := strings.TrimSpace(newLines[len(newLines)-1])
				if strings.HasSuffix(lastLine, ".tmpl") || lastLine == "custom_code:" {
					// No extra newline needed after custom_code
					addSweeperBlock(baseIndent, locationPairs, &newLines)
				} else {
					// Add newline before sweeper for other sections
					newLines = append(newLines, "")
					addSweeperBlock(baseIndent, locationPairs, &newLines)
				}
			}
			sweepersFound = true
		}

		// Track sweeper section
		if trimmed == "sweeper:" {
			inSweeper = true
			newLines = append(newLines, line)

			// Add url_substitutions immediately after sweeper:
			if len(locationPairs) > 0 && !urlSubsAdded {
				newLines = append(newLines, baseIndent+"  url_substitutions:")
				addLocationPairs(baseIndent+"    ", locationPairs, &newLines)
				urlSubsAdded = true
			}
			continue
		}

		// Skip existing url_substitutions sections
		if inSweeper && indent == len(baseIndent)+2 {
			if trimmed == "url_substitutions:" {
				skipToNextSection = true
				continue
			} else if trimmed != "" {
				skipToNextSection = false
				inSweeper = indent <= len(baseIndent)
			}
		}

		if skipToNextSection {
			nextLineIdx := i + 1
			if nextLineIdx < len(lines) {
				nextIndent := len(lines[nextLineIdx]) - len(strings.TrimLeft(lines[nextLineIdx], " "))
				if nextIndent <= len(baseIndent)+2 && strings.TrimSpace(lines[nextLineIdx]) != "" {
					skipToNextSection = false
				}
			}
			continue
		}

		newLines = append(newLines, line)
	}

	// If we still need to add the sweeper section and haven't yet
	if !sweepersFound && len(locationPairs) > 0 {
		// Check if we're after custom_code
		if len(newLines) > 0 {
			// No extra newline needed after custom_code
			addSweeperBlock(baseIndent, locationPairs, &newLines)
		}
	}

	// Ensure file ends with a newline
	if len(newLines) > 0 && newLines[len(newLines)-1] != "" {
		newLines = append(newLines, "")
	}

	return os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644)
}

func addSweeperBlock(baseIndent string, pairs []LocationPair, lines *[]string) {
	// Don't add extra newline, just add the sweeper block
	*lines = append(*lines, "sweeper:")
	*lines = append(*lines, baseIndent+"  url_substitutions:")
	addLocationPairs(baseIndent+"    ", pairs, lines)
}

func addLocationPairs(indent string, pairs []LocationPair, lines *[]string) {
	for _, pair := range pairs {
		if pair.Region != "" && pair.Zone != "" {
			*lines = append(*lines,
				fmt.Sprintf("%s- region: \"%s\"", indent, pair.Region),
				fmt.Sprintf("%s  zone: \"%s\"", indent, pair.Zone))
		} else if pair.Region != "" {
			*lines = append(*lines,
				fmt.Sprintf("%s- region: \"%s\"", indent, pair.Region))
		} else if pair.Zone != "" {
			*lines = append(*lines,
				fmt.Sprintf("%s- zone: \"%s\"", indent, pair.Zone))
		}
	}
}
