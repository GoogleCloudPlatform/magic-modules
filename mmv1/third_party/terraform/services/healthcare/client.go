package healthcare

import (
	"log"

	healthcare "google.golang.org/api/healthcare/v1"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *healthcare.Service {
	healthcareClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Google Cloud Healthcare client for path %s", healthcareClientBasePath)
	clientHealthcare, err := healthcare.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client healthcare: %s", err)
		return nil
	}
	clientHealthcare.UserAgent = userAgent
	clientHealthcare.BasePath = healthcareClientBasePath

	return clientHealthcare
}
