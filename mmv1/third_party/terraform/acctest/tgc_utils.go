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

var CaiProductBackendNames = map[string]string{
	"compute":         "compute",
	"resourcemanager": "cloudresourcemanager",
}

// Gets the test metadata for tgc:
//   - test config
//   - cai asset name
//     For example: //compute.googleapis.com/projects/terraform-dev-zhenhuali/zones/us-central1-a/instances/tf-test-mi3fqaucf8
func GetTestMetadataForTgc(service, resourceType, resourceName, config string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		address := fmt.Sprintf("%s.%s", resourceType, resourceName)
		rState := s.RootModule().Resources[address]

		var rId string
		switch resourceType {
		case "google_project":
			rId = fmt.Sprintf("projects/%s", rState.Primary.Attributes["number"])
		default:
			rId = rState.Primary.ID
		}
		log.Printf("[DEBUG]TGC resource ID in Terraform state: %s", rId)
		// Convert the resource ID into CAI asset name
		// and then print out the CAI asset name in the logs
		if productName, ok := CaiProductBackendNames[service]; !ok {
			return fmt.Errorf("The Cai product backend name for service %s doesn't exist.", service)
		} else {
			caiAssetName := fmt.Sprintf("//%s.googleapis.com/%s", productName, rId)
			log.Printf("[DEBUG]TGC CAI asset name: %s", caiAssetName)
		}

		log.Printf("[DEBUG]TGC Terraform service: %s", service)
		log.Printf("[DEBUG]TGC Terraform resource: %s", resourceType)

		re := regexp.MustCompile(`\"(tf[-_]?test[-_]?.*?)([a-z0-9]+)\"`)
		config = re.ReplaceAllString(config, `"${1}tgc"`)

		// Replace resource name with the resource's real name
		// For example, replace `"google_compute_instance" "foobar"` with `"google_compute_instance" "tf-test-mi3fqaucf8"`
		n := tpgresource.GetResourceNameFromSelfLink(rId)
		old := fmt.Sprintf(`"%s" "%s"`, resourceType, resourceName)
		new := fmt.Sprintf(`"%s" "%s"`, resourceType, n)
		config = strings.Replace(config, old, new, 1)

		log.Printf("[DEBUG]TGC raw_config starts %s", config)
		log.Printf("[DEBUG]End of TGC raw_config")
		return nil
	}
}
