package google

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("AppEngineFlexibleAppVersion", &resource.Sweeper{
		Name: "AppEngineFlexibleAppVersion",
		F:    testSweepAppEngineFlexibleAppVersion,
	})
}

// At the time of writing, the CI only passes us-central1 as the region
func testSweepAppEngineFlexibleAppVersion(region string) error {
	resourceName := "AppEngineFlexibleAppVersion"
	log.Printf("[INFO][SWEEPER_LOG] Starting sweeper for %s", resourceName)

	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	// Setup variables to replace in list template
	d := &ResourceDataMock{
		FieldsInSchema: map[string]interface{}{
			"project":  config.Project,
			"region":   region,
			"location": region,
			"zone":     "-",
		},
	}

	deleteTemplate := "https://appengine.googleapis.com/v1/apps/{{project}}/services/{{service}}"
	deleteUrl, err := replaceVars(d, config, deleteTemplate)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error preparing delete url: %s", err)
		return nil
	}

	// Don't wait on operations as we may have a lot to delete
	_, err = sendRequest(config, "DELETE", config.Project, deleteUrl, nil)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] Error deleting for url %s : %s", deleteUrl, err)
	} else {
		log.Printf("[INFO][SWEEPER_LOG] Sent delete request for %s resource: %s", resourceName, d.Get("service"))
	}

	return nil
}
