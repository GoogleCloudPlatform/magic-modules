package logging

import (
	"log"

	cloudlogging "google.golang.org/api/logging/v2"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *cloudlogging.Service {
	loggingClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Google Stackdriver Logging client for path %s", loggingClientBasePath)
	clientLogging, err := cloudlogging.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client logging: %s", err)
		return nil
	}
	clientLogging.UserAgent = userAgent
	clientLogging.BasePath = loggingClientBasePath

	return clientLogging
}
