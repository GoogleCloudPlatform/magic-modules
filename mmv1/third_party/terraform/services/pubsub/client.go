package pubsub

import (
	"log"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/option"
	pubsub "google.golang.org/api/pubsub/v1"
)

func NewClient(c *transport_tpg.Config, userAgent string) *pubsub.Service {
	pubsubClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Google Pubsub client for path %s", pubsubClientBasePath)
	wrappedPubsubClient := transport_tpg.ClientWithAdditionalRetries(c.Client, transport_tpg.PubsubTopicProjectNotReady)
	clientPubsub, err := pubsub.NewService(c.Context, option.WithHTTPClient(wrappedPubsubClient))
	if err != nil {
		log.Printf("[WARN] Error creating client pubsub: %s", err)
		return nil
	}
	clientPubsub.UserAgent = userAgent
	clientPubsub.BasePath = pubsubClientBasePath

	return clientPubsub
}
