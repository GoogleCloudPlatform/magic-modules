package certificatemanager

import (
	"log"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/certificatemanager/v1"
	"google.golang.org/api/option"
)

func NewClient(c *transport_tpg.Config, userAgent string) *certificatemanager.Service {
	certificateManagerClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Certificate Manager client for path %s", certificateManagerClientBasePath)
	clientCertificateManager, err := certificatemanager.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client certificate manager: %s", err)
		return nil
	}
	clientCertificateManager.UserAgent = userAgent
	clientCertificateManager.BasePath = certificateManagerClientBasePath

	return clientCertificateManager
}
