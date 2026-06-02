package iamcredentials

import (
	"log"

	iamcredentials "google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func NewClient(c *transport_tpg.Config, userAgent string) *iamcredentials.Service {
	iamCredentialsClientBasePath := transport_tpg.RemoveBasePathVersion(transport_tpg.BaseUrl(Product, c))
	log.Printf("[INFO] Instantiating Google Cloud IAMCredentials client for path %s", iamCredentialsClientBasePath)
	clientIamCredentials, err := iamcredentials.NewService(c.Context, option.WithHTTPClient(c.Client))
	if err != nil {
		log.Printf("[WARN] Error creating client iam credentials: %s", err)
		return nil
	}
	clientIamCredentials.UserAgent = userAgent
	clientIamCredentials.BasePath = iamCredentialsClientBasePath

	return clientIamCredentials
}
