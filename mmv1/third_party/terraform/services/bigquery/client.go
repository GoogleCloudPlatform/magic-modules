package bigquery

import (
	"log"

	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *bigquery.Service {
	bigQueryClientBasePath := transport_tpg.BaseUrl(Product, c)
	log.Printf("[INFO] Instantiating Google Cloud BigQuery client for path %s", bigQueryClientBasePath)
	wrappedBigQueryClient := transport_tpg.ClientWithAdditionalRetries(c.Client, transport_tpg.IamMemberMissing)
	clientBigQuery, err := bigquery.NewService(c.Context, option.WithHTTPClient(wrappedBigQueryClient))
	if err != nil {
		log.Printf("[WARN] Error creating client big query: %s", err)
		return nil
	}
	clientBigQuery.UserAgent = userAgent
	clientBigQuery.BasePath = bigQueryClientBasePath

	return clientBigQuery
}
