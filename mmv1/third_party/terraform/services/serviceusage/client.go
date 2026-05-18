package serviceusage

import (
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/serviceusage/v1"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *serviceusage.Service {
	serviceUsageClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Google Cloud Service Usage client for path %s", serviceUsageClientBasePath)
	clientServiceUsage, err := serviceusage.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client service usage: %s", err)
		return nil
	}
	clientServiceUsage.UserAgent = userAgent
	clientServiceUsage.BasePath = serviceUsageClientBasePath

	return clientServiceUsage
}
