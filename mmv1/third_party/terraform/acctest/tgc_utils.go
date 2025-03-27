package acctest

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// Hardcode the Terraform resource name -> API service name mapping temporarily.
// TODO: [tgc] read the mapping from the resource metadata files.
var CaiProductBackendNames = map[string]string{
	"google_compute_instance": "compute",
	"google_project":          "cloudresourcemanager",
}

// Gets the test metadata for tgc:
//   - test config
//   - cai asset name
//     For example: //compute.googleapis.com/projects/ci-test-188019/zones/us-central1-a/instances/tf-test-mi3fqaucf8
func GetTestMetadataForTgc(service, resourceType, resourceName, config string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		address := fmt.Sprintf("%s.%s", resourceType, resourceName)
		rState := s.RootModule().Resources[address]
		if rState == nil || rState.Primary == nil {
			return fmt.Errorf("The resource state is unavailable. Please check if the address %s.%s is correct.", resourceType, resourceName)
		}

		// Convert the resource ID into CAI asset name
		// and then print out the CAI asset name in the logs
		if productName, ok := CaiProductBackendNames[resourceType]; !ok {
			return fmt.Errorf("The Cai product backend name for resource %s doesn't exist.", resourceType)
		} else {
			var rName string
			switch resourceType {
			case "google_project":
				rName = fmt.Sprintf("projects/%s", rState.Primary.Attributes["number"])
			default:
				rName = rState.Primary.ID
			}

			caiAssetName := fmt.Sprintf("//%s.googleapis.com/%s", productName, rName)
			log.Printf("[DEBUG]TGC CAI asset name: %s", caiAssetName)
		}

		log.Printf("[DEBUG]TGC Terraform service: %s", service)
		log.Printf("[DEBUG]TGC Terraform resource: %s", resourceType)

		re := regexp.MustCompile(`\"(tf[-_]?test[-_]?.*?)([a-z0-9]+)\"`)
		config = re.ReplaceAllString(config, `"${1}tgc"`)

		// Replace resource name with the resource's real name,
		// which is used to get the main resource object by checking the address after parsing raw config.
		// For example, replace `"google_compute_instance" "foobar"` with `"google_compute_instance" "tf-test-mi3fqaucf8"`
		n := tpgresource.GetResourceNameFromSelfLink(rState.Primary.ID)
		old := fmt.Sprintf(`"%s" "%s"`, resourceType, resourceName)
		new := fmt.Sprintf(`"%s" "%s"`, resourceType, n)
		config = strings.Replace(config, old, new, 1)

		log.Printf("[DEBUG]TGC raw_config starts %sEnd of TGC raw_config", config)
		return nil
	}
}
