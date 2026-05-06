package backupdr

import (
	"log"

	backupdr "google.golang.org/api/backupdr/v1"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *backupdr.Service {
	backupdrClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c)))
	log.Printf("[INFO] Instantiating Google SqlAdmin client for path %s", backupdrClientBasePath)
	clientBackupdrAdmin, err := backupdr.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage: %s", err)
		return nil
	}
	clientBackupdrAdmin.UserAgent = userAgent
	clientBackupdrAdmin.BasePath = backupdrClientBasePath

	return clientBackupdrAdmin
}
