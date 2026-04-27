package dataproc

import (
	"log"

	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *dataproc.Service {
	dataprocClientBasePath := transport_tpg.RemoveBasePathVersion(c.DataprocBasePath)
	log.Printf("[INFO] Instantiating Google Cloud Dataproc client for path %s", dataprocClientBasePath)
	clientDataproc, err := dataproc.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataproc: %s", err)
		return nil
	}
	clientDataproc.UserAgent = userAgent
	clientDataproc.BasePath = dataprocClientBasePath

	return clientDataproc
}
