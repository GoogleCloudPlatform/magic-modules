package acctest

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Hardcode the Terraform resource name -> API service name mapping temporarily.
// TODO: [tgc] read the mapping from the resource metadata files.
var ApiServiceNames = map[string]string{
	"google_compute_instance": "compute.googleapis.com",
	"google_project":          "cloudresourcemanager.googleapis.com",
}

// Gets the test metadata for tgc:
//   - test config
//   - cai asset name
//     For example: //compute.googleapis.com/projects/ci-test-188019/zones/us-central1-a/instances/tf-test-mi3fqaucf8
func GetTestMetadataForTgc(service, address, rawConfig string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		splits := strings.Split(address, ".")
		if splits == nil || len(splits) < 2 {
			return fmt.Errorf("The resource address %s is invalid.", address)
		}
		resourceType := splits[0]
		resourceName := splits[1]

		rState := s.RootModule().Resources[address]
		if rState == nil || rState.Primary == nil {
			return fmt.Errorf("The resource state is unavailable. Please check if the address %s.%s is correct.", resourceType, resourceName)
		}

		// Convert the resource ID into CAI asset name
		// and then print out the CAI asset name in the logs
		if apiServiceName, ok := ApiServiceNames[resourceType]; !ok {
			unknownType := "unkown"
			log.Printf("[DEBUG]TGC CAI asset names start\n%s\nEnd of TGC CAI asset names", unknownType)
		} else {
			var rName string
			switch resourceType {
			case "google_project":
				rName = fmt.Sprintf("projects/%s", rState.Primary.Attributes["number"])
			default:
				rName = rState.Primary.ID
			}
			caiAssetName := fmt.Sprintf("//%s/%s", apiServiceName, rName)
			log.Printf("[DEBUG]TGC CAI asset names start\n%s\nEnd of TGC CAI asset names", caiAssetName)
		}

		// The acceptance tests names will be also used for the tgc tests.
		// "service" is logged and will be used to put the tgc tests into specific service packages.
		log.Printf("[DEBUG]TGC Terraform service: %s", service)
		log.Printf("[DEBUG]TGC Terraform resource: %s", address)

		re := regexp.MustCompile(`\"(tf[-_]?test[-_]?.*?)([a-z0-9]+)\"`)
		rawConfig = re.ReplaceAllString(rawConfig, `"${1}tgc"`)
		log.Printf("[DEBUG]TGC raw_config starts %sEnd of TGC raw_config", rawConfig)
		return nil
	}
}

// parseResources extracts all resources from a Terraform configuration string
func parseResources(config string) []string {
	// This regex matches resource blocks in Terraform configurations
	resourceRegex := regexp.MustCompile(`resource\s+"([^"]+)"\s+"([^"]+)"`)
	matches := resourceRegex.FindAllStringSubmatch(config, -1)

	var resources []string
	for _, match := range matches {
		if len(match) >= 3 {
			// Combine resource type and name to form the address
			resources = append(resources, fmt.Sprintf("%s.%s", match[1], match[2]))
		}
	}

	return resources
}

// getServicePackage determines the service package for a resource type
func getServicePackage(resourceType string) string {
	var ServicePackages = map[string]string{
		"google_compute_":   "compute",
		"google_storage_":   "storage",
		"google_sql_":       "sql",
		"google_container_": "container",
		"google_bigquery_":  "bigquery",
		"google_project":    "resourcemanager",
		"google_cloud_run_": "cloudrun",
	}

	// Check for exact matches first
	if service, ok := ServicePackages[resourceType]; ok {
		return service
	}

	// Check for prefix matches
	for prefix, service := range ServicePackages {
		if strings.HasPrefix(resourceType, prefix) {
			return service
		}
	}

	// Default to "unknown" if no match found
	return "unknown"
}

// extendWithTGCData adds TGC metadata check functions to each TestStep
func extendWithTGCData(c resource.TestCase) resource.TestCase {
	var updatedSteps []resource.TestStep

	for _, step := range c.Steps {
		if step.Config != "" {
			// Parse resources from the config
			resources := parseResources(step.Config)

			// Create TGC metadata checks for each resource
			var tgcChecks []resource.TestCheckFunc
			for _, res := range resources {
				parts := strings.Split(res, ".")
				if len(parts) >= 2 {
					resourceType := parts[0]
					service := getServicePackage(resourceType)
					tgcChecks = append(tgcChecks, GetTestMetadataForTgc(service, res, step.Config))
				}
			}

			// If there are TGC checks to add
			if len(tgcChecks) > 0 {
				// If there's an existing check function, wrap it with ours
				if step.Check != nil {
					existingCheck := step.Check
					step.Check = resource.ComposeTestCheckFunc(
						existingCheck,
						resource.ComposeTestCheckFunc(tgcChecks...),
					)
				} else {
					// Otherwise, just use our TGC checks
					step.Check = resource.ComposeTestCheckFunc(tgcChecks...)
				}
			}
		}

		updatedSteps = append(updatedSteps, step)
	}

	c.Steps = updatedSteps
	return c
}
