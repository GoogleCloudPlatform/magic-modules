package resourcemanager

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/sweeper"
)

func init() {
	sweeper.AddTestSweepersLegacy("GoogleFolder", testSweepFolder)
}
func testSweepFolder(region string) error {
	config, err := sweeper.SharedConfigForRegion(region)
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error getting shared config for region: %s", err)
		return err
	}

	err = config.LoadAndValidate(context.Background())
	if err != nil {
		log.Printf("[INFO][SWEEPER_LOG] error loading: %s", err)
		return err
	}

	org := envvar.UnsafeGetTestOrgFromEnv()
	log.Printf("[DEBUG] org %s", org)

	if org == "" {
		log.Printf("[INFO][SWEEPER_LOG] no organization set, failing folder sweeper")
		return fmt.Errorf("no organization set")
	}

	parent := "organizations/" + org

	token := ""
	for paginate := true; paginate; {
		// Filter for folders with test prefix
		// filter := fmt.Sprintf("id:\"%s*\" -lifecycleState:DELETE_REQUESTED parent.id:%v", TestPrefix, org)
		found, err := config.NewResourceManagerV3Client(config.UserAgent).Folders.List().Parent(parent).PageToken(token).Do()
		if err != nil {
			log.Printf("[INFO][SWEEPER_LOG] error listing folders: %s", err)
			return nil
		}

		for _, folder := range found.Folders {
			if !strings.HasPrefix(folder.DisplayName, TestPrefix) {
				continue
			}
			log.Printf("[INFO][SWEEPER_LOG] Sweeping Folder id: %s, name: %s", folder.Name, folder.DisplayName)
			_, err := config.NewResourceManagerV3Client(config.UserAgent).Folders.Delete(folder.Name).Do()
			if err != nil {
				log.Printf("[INFO][SWEEPER_LOG] Error, failed to delete folder %s: %s", folder.Name, err)
				continue
			}
		}
		token = found.NextPageToken
		paginate = token != ""
	}

	return nil
}
