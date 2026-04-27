package spanner

import (
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/spanner/v1"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *spanner.Service {
	spannerClientBasePath := transport_tpg.RemoveBasePathVersion(c.SpannerBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Spanner client for path %s", spannerClientBasePath)
	clientSpanner, err := spanner.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client source repo: %s", err)
		return nil
	}
	clientSpanner.UserAgent = userAgent
	clientSpanner.BasePath = spannerClientBasePath

	return clientSpanner
}
