package storagetransfer

import (
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/storagetransfer/v1"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *storagetransfer.Service {
	storageTransferClientBasePath := transport_tpg.RemoveBasePathVersion(c.StorageTransferBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Storage Transfer client for path %s", storageTransferClientBasePath)
	clientStorageTransfer, err := storagetransfer.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client storage transfer: %s", err)
		return nil
	}
	clientStorageTransfer.UserAgent = userAgent
	clientStorageTransfer.BasePath = storageTransferClientBasePath

	return clientStorageTransfer
}
