package cloudrunv2

import (
	"log"

	"google.golang.org/api/option"
	runadminv2 "google.golang.org/api/run/v2"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *runadminv2.Service {
	runAdminV2ClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c)))
	log.Printf("[INFO] Instantiating Google Cloud Run Admin v2 client for path %s", runAdminV2ClientBasePath)
	clientRunAdminV2, err := runadminv2.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client run admin: %s", err)
		return nil
	}
	clientRunAdminV2.UserAgent = userAgent
	clientRunAdminV2.BasePath = runAdminV2ClientBasePath

	return clientRunAdminV2
}
