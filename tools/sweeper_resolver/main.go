package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/mmv1/api"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"gopkg.in/yaml.v2"
)

// Custom logger that only prints when debug mode is enabled
type debugLogger struct {
	debug bool
}

func (l *debugLogger) Printf(format string, v ...interface{}) {
	if l.debug {
		log.Printf(format, v...)
	}
}

type MetaData struct {
	Resource       string `yaml:"resource"`
	GenerationType string `yaml:"generation_type"`
	SourceFile     string `yaml:"source_file,omitempty"`
}

// LocationPair represents a correlated region-zone pair from a test configuration
type LocationPair struct {
	Region string
	Zone   string
}

// ResourceLocations stores the unique location pairs for a resource
type ResourceLocations struct {
	LocationPairs []LocationPair
}

func main() {
	// Define command line flags
	dirPath := flag.String("dir", "", "Path to services directory")
	debug := flag.Bool("debug", false, "Enable debug logging")
	flag.Parse()

	// Initialize debug logger
	logger := &debugLogger{debug: *debug}

	// Configure logging format if debug is enabled
	if *debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	logger.Printf("Starting program with directory path: %s", *dirPath)

	// Check for empty directory path after trimming spaces
	if strings.TrimSpace(*dirPath) == "" {
		log.Fatal("Please provide a directory path using -dir flag")
	}

	// Clean the path and verify it exists
	cleanPath := filepath.Clean(*dirPath)
	if _, err := os.Stat(cleanPath); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", cleanPath)
	}

	allTests, errs := reader.ReadAllTests(cleanPath)
	if errs != nil {
		logger.Printf("Warning: Errors encountered while reading tests: %v", errs)
	}
	logger.Printf("Found %d tests", len(allTests))

	resourceMap, err := processServicesDirectory(cleanPath, logger)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}
	logger.Printf("Processed %d resources from directory", len(resourceMap))

	// Track location pairs for each resource
	locations := map[string]*ResourceLocations{}

	// Process tests to find region-zone pairs
	for testIndex, test := range allTests {
		logger.Printf("Processing test %d", testIndex)
		for stepIndex, config := range test.Steps {
			for resourceName, resourceConfig := range config {
				logger.Printf("Found resource config in test %d, step %d", testIndex, stepIndex)

				_, hasMmv1Resource := resourceMap[resourceName]
				if !hasMmv1Resource {
					logger.Printf("Resource %s not found in MMv1 resources, skipping", resourceName)
					continue
				}

				// Initialize locations for this resource if not exists
				if _, exists := locations[resourceName]; !exists {
					locations[resourceName] = &ResourceLocations{
						LocationPairs: []LocationPair{},
					}
				}

				// Process each configuration block
				for _, fields := range resourceConfig {
					logger.Printf("Processing fields: %v", fields)

					// Extract and process region/location
					var regionStr string
					region, hasRegion := fields["region"]
					if !hasRegion {
						region, hasRegion = fields["location"]
					}
					if hasRegion && region != "" {
						regionStr = processLocationString(region.(string), config, nil, logger)
					}

					// Extract and process zone
					var zoneStr string
					zone, hasZone := fields["zone"]
					if hasZone && zone != "" {
						zoneStr = processLocationString(zone.(string), config, nil, logger)
					}

					// If we have either a region or zone, create a location pair
					if regionStr != "" || zoneStr != "" {
						pair := LocationPair{
							Region: regionStr,
							Zone:   zoneStr,
						}

						// Only add if this pair doesn't already exist
						if !containsLocationPair(locations[resourceName].LocationPairs, pair) {
							locations[resourceName].LocationPairs = append(locations[resourceName].LocationPairs, pair)
							logger.Printf("Added location pair for resource %s: region=%s, zone=%s",
								resourceName, pair.Region, pair.Zone)
						}
					}
				}
			}
		}
	}

	logger.Printf("Found locations for %d resources", len(locations))

	// Process and update resources
	for resourceName, resourceLocations := range locations {
		logger.Printf("Processing locations for resource: %s", resourceName)

		mmv1Resource, exists := resourceMap[resourceName]
		if !exists {
			logger.Printf("Resource %s not found in resource map, skipping", resourceName)
			continue
		}

		if len(resourceLocations.LocationPairs) > 0 {
			logger.Printf("Found location pairs for resource %s", resourceName)

			// Check if we need to update
			if !shouldSkipUpdate(resourceLocations.LocationPairs, mmv1Resource.Sweeper) {
				fmt.Printf("Resource: %s\n", resourceName)
				fmt.Printf("sweeper:\n  url_substitutions:\n")
				for _, pair := range resourceLocations.LocationPairs {
					if pair.Region != "" && pair.Zone != "" {
						fmt.Printf("    - region: \"%s\"\n      zone: \"%s\"\n", pair.Region, pair.Zone)
					} else if pair.Region != "" {
						fmt.Printf("    - region: \"%s\"\n", pair.Region)
					} else if pair.Zone != "" {
						fmt.Printf("    - zone: \"%s\"\n", pair.Zone)
					}
				}
				fmt.Println()

				// Update the YAML file
				if err := updateYamlFile(mmv1Resource.SourceYamlFile, resourceLocations.LocationPairs, mmv1Resource.Sweeper); err != nil {
					logger.Printf("Error updating YAML file: %v", err)
				}
			} else {
				logger.Printf("Skipping update for resource %s (matches existing configuration)", resourceName)
			}
		} else {
			logger.Printf("No location pairs found for resource %s, skipping", resourceName)
		}
	}
}

func processServicesDirectory(dirPath string, logger *debugLogger) (map[string]*api.Resource, error) {
	logger.Printf("Processing services directory: %s", dirPath)
	resourceMap := make(map[string]*api.Resource)

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Printf("Error accessing path %s: %v", path, err)
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), "_generated_meta.yaml") {
			return nil
		}

		logger.Printf("Processing YAML file: %s", path)

		metadata, err := processYAMLFile(path)
		if err != nil {
			logger.Printf("Warning: Error processing file %s: %v", path, err)
			return nil
		}

		if metadata.GenerationType == "mmv1" {
			if metadata.SourceFile == "" {
				logger.Printf("Warning: MMv1 resource %s missing source_file", metadata.Resource)
				return nil
			}

			resourceModel := InitializeResource(metadata.SourceFile)
			logger.Printf("shouldgen: %v", resourceModel.ShouldGenerateSweepers())

			if !resourceModel.ShouldGenerateSweepers() {
				fmt.Printf("%s : %s\n", resourceModel.TerraformName(), resourceModel.ListUrlTemplate())
				logger.Printf("Resource %s should not generate sweepers, skipping", metadata.Resource)
				return nil
			}

			resourceMap[metadata.Resource] = resourceModel
			logger.Printf("Added resource %s to map", metadata.Resource)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	logger.Printf("Finished processing directory, found %d resources", len(resourceMap))
	return resourceMap, nil
}

func processYAMLFile(filePath string) (*MetaData, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var metadata MetaData
	err = yaml.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML: %v", err)
	}

	if metadata.Resource == "" {
		return nil, fmt.Errorf("missing required field 'resource'")
	}
	if metadata.GenerationType == "" {
		return nil, fmt.Errorf("missing required field 'generation_type'")
	}

	return &metadata, nil
}

func containsLocationPair(pairs []LocationPair, newPair LocationPair) bool {
	for _, pair := range pairs {
		if pair.Region == newPair.Region && pair.Zone == newPair.Zone {
			return true
		}
	}
	return false
}

func processLocationString(location string, config reader.Step, seen map[string]bool, logger *debugLogger) string {
	locationStr := strings.Trim(location, `"`)
	if strings.Contains(locationStr, ".") {
		locationStr = resolveLocationReference(locationStr, config, seen, logger)
	}

	// Skip invalid locations
	if locationStr != "" && locationStr != "true" &&
		!strings.Contains(locationStr, "%") &&
		!strings.Contains(locationStr, "/") &&
		!strings.Contains(locationStr, ".") {
		return locationStr
	}
	return ""
}

func resolveLocationReference(ref string, config reader.Step, seen map[string]bool, logger *debugLogger) string {
	if seen == nil {
		seen = make(map[string]bool)
	}

	// Prevent infinite recursion
	if seen[ref] {
		logger.Printf("Circular reference detected for: %s", ref)
		return ref
	}
	seen[ref] = true

	parts := strings.Split(strings.Trim(ref, `"`), ".")
	if len(parts) != 3 {
		return ref // Not a reference, return as is
	}

	resourceType := parts[0]
	resourceName := parts[1]
	field := parts[2]

	if resources, exists := config[resourceType]; exists {
		if resource, exists := resources[resourceName]; exists {
			if value, exists := resource[field]; exists {
				logger.Printf("Found value for reference %s: %v", ref, value)

				if strValue, ok := value.(string); ok {
					strValue = strings.Trim(strValue, `"`)
					if strings.Contains(strValue, ".") && strings.Count(strValue, ".") == 2 {
						return resolveLocationReference(strValue, config, seen, logger)
					}
					return strValue
				}
				return fmt.Sprintf("%v", value)
			}
		}
	}

	logger.Printf("Could not resolve location reference: %s", ref)
	return ""
}
