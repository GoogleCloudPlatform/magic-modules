package servicenetworking

import (
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/servicenetworking/v1"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *servicenetworking.APIService {
	serviceNetworkingClientBasePath := transport_tpg.RemoveBasePathVersion(c.ServiceNetworkingBasePath)
	log.Printf("[INFO] Instantiating Service Networking client for path %s", serviceNetworkingClientBasePath)
	clientServiceNetworking, err := servicenetworking.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client service networking: %s", err)
		return nil
	}
	clientServiceNetworking.UserAgent = userAgent
	clientServiceNetworking.BasePath = serviceNetworkingClientBasePath

	return clientServiceNetworking
}
