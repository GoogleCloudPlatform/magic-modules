// Package client provides a function for creating a resourcemanager client.
// This is a separate package due to a circular dependency with compute.
package client

import (
	"log"

	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *cloudresourcemanager.Service {
	resourceManagerBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(registry.GetProduct("resourcemanager"), c))
	log.Printf("[INFO] Instantiating Google Cloud ResourceManager client for path %s", resourceManagerBasePath)
	clientResourceManager, err := cloudresourcemanager.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client resource manager: %s", err)
		return nil
	}
	clientResourceManager.UserAgent = userAgent
	clientResourceManager.BasePath = resourceManagerBasePath

	return clientResourceManager
}
