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

	if *dirPath == "" {
		log.Fatal("Please provide a directory path using -dir flag")
	}

	allTests, errs := reader.ReadAllTests(*dirPath)
	if errs != nil {
		logger.Printf("Warning: Errors encountered while reading tests: %v", errs)
	}
	logger.Printf("Found %d tests", len(allTests))

	resourceMap, err := processServicesDirectory(*dirPath, logger)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}
	logger.Printf("Processed %d resources from directory", len(resourceMap))

	regions := map[string][]string{}

	// In the main processing loop:
	for testIndex, test := range allTests {
		logger.Printf("Processing test %d", testIndex)
		for stepIndex, config := range test.Steps {
			for resourceName, resourceConfig := range config {
				logger.Printf("Found resource config in test %d, step %d", testIndex, stepIndex)
				for _, fields := range resourceConfig {
					logger.Printf("resourceConfig: %v", fields)

					_, hasMmv1Resource := resourceMap[resourceName]
					if !hasMmv1Resource {
						logger.Printf("Resource %s not found in MMv1 resources, skipping", resourceName)
						continue
					}

					logger.Printf("fields %v", fields)

					region, hasRegion := fields["region"]
					if !hasRegion {
						region, hasRegion = fields["location"]
						if !hasRegion {
							logger.Printf("No region/location found for resource %s", resourceName)
							continue
						}
					}

					if region != "" {
						// Clean and resolve any region references
						regionStr := strings.Trim(region.(string), `"`)
						if strings.Contains(regionStr, ".") {
							regionStr = resolveRegionReference(regionStr, config, nil, logger)
						}

						// Skip empty regions after cleanup
						if regionStr != "" && regionStr != "true" && !strings.Contains(regionStr, "%") && !strings.Contains(regionStr, "/") && !strings.Contains(regionStr, ".") {
							logger.Printf("Adding region %v for resource %s", regionStr, resourceName)
							regions[resourceName] = append(regions[resourceName], regionStr)
						}
					}
				}
			}
		}
	}

	logger.Printf("Found regions for %d resources", len(regions))

	for resourceName, resourceRegions := range regions {
		logger.Printf("Processing regions for resource: %s", resourceName)

		uniqueRegions := make(map[string]bool)
		for _, region := range resourceRegions {
			uniqueRegions[region] = true
		}

		var finalRegions []string
		for region := range uniqueRegions {
			finalRegions = append(finalRegions, region)
		}
		logger.Printf("Found %d unique regions for resource %s: %v",
			len(finalRegions), resourceName, finalRegions)

		mmv1Resource, exists := resourceMap[resourceName]
		if !exists {
			logger.Printf("Resource %s not found in resource map, skipping", resourceName)
			continue
		}

		defaultRegion := []string{"us-central1"}
		logger.Printf("Comparing regions for %s - Current: %v, Default: %v, Sweeper: %v",
			resourceName, finalRegions, defaultRegion, mmv1Resource.Sweeper.Regions)

		if !equalRegions(finalRegions, defaultRegion) && !equalRegions(finalRegions, mmv1Resource.Sweeper.Regions) {
			logger.Printf("Regions differ for resource %s, updating YAML", resourceName)

			// Print what we're going to do
			fmt.Printf("Resource: %s\nsweepers:\n  regions:\n", resourceName)
			for _, region := range finalRegions {
				fmt.Printf("    - \"%s\"\n", region)
			}
			fmt.Println()

			// Update the YAML file
			if err := updateYamlFile(mmv1Resource.SourceYamlFile, finalRegions); err != nil {
				logger.Printf("Error updating YAML file: %v", err)
			}
		} else {
			logger.Printf("Regions match defaults for resource %s, skipping", resourceName)
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
			logger.Printf(" shouldgen: %v", resourceModel.ShouldGenerateSweepers())

			if !resourceModel.ShouldGenerateSweepers() {
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

func equalRegions(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]bool)
	for _, val := range a {
		aMap[val] = true
	}

	for _, val := range b {
		if !aMap[val] {
			return false
		}
	}

	return true
}

func resolveRegionReference(regionRef string, config reader.Step, seen map[string]bool, logger *debugLogger) string {
	if seen == nil {
		seen = make(map[string]bool)
	}

	// Prevent infinite recursion
	if seen[regionRef] {
		logger.Printf("Circular reference detected for: %s", regionRef)
		return regionRef
	}
	seen[regionRef] = true

	parts := strings.Split(strings.Trim(regionRef, `"`), ".")
	if len(parts) != 3 {
		return regionRef // Not a reference, return as is
	}

	resourceType := parts[0] // e.g., "google_compute_subnetwork"
	resourceName := parts[1] // e.g., "subnetwork"
	field := parts[2]        // e.g., "region"

	// Look for the referenced resource in the config
	if resources, exists := config[resourceType]; exists {
		if resource, exists := resources[resourceName]; exists {
			if value, exists := resource[field]; exists {
				logger.Printf("Found value for reference %s: %v", regionRef, value)

				// If the value is itself a reference, resolve it recursively
				if strValue, ok := value.(string); ok {
					strValue = strings.Trim(strValue, `"`)
					if strings.Contains(strValue, ".") && strings.Count(strValue, ".") == 2 {
						return resolveRegionReference(strValue, config, seen, logger)
					}
					return strValue
				}
				return fmt.Sprintf("%v", value)
			}
		}
	}

	logger.Printf("Could not resolve region reference: %s", regionRef)
	return ""
}
