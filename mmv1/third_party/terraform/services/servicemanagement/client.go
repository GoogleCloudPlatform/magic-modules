package servicemanagement

import (
	"log"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"
)

func NewClient(c *transport_tpg.Config, userAgent string) *servicemanagement.APIService {
	serviceManagementClientBasePath := transport_tpg.RemoveBasePathVersion(c.ServiceManagementBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Service Management client for path %s", serviceManagementClientBasePath)
	clientServiceMan, err := servicemanagement.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client service management: %s", err)
		return nil
	}
	clientServiceMan.UserAgent = userAgent
	clientServiceMan.BasePath = serviceManagementClientBasePath

	return clientServiceMan
}
