package dataflow

import (
	"log"

	dataflow "google.golang.org/api/dataflow/v1b3"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *dataflow.Service {
	dataflowClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Google Dataflow client for path %s", dataflowClientBasePath)
	clientDataflow, err := dataflow.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client dataflow: %s", err)
		return nil
	}
	clientDataflow.UserAgent = userAgent
	clientDataflow.BasePath = dataflowClientBasePath

	return clientDataflow
}
