package firebaserules

import (
	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"time"
)

func NewDCLFirebaserulesClient(config *transport_tpg.Config, userAgent, billingProject string, timeout time.Duration) *Client {
	configOptions := []dcl.ConfigOption{
		dcl.WithHTTPClient(config.Client),
		dcl.WithUserAgent(userAgent),
		dcl.WithLogger(dcl.DCLLogger{}),
		dcl.WithBasePath(config.FirebaserulesBasePath),
	}

	if timeout != 0 {
		configOptions = append(configOptions, dcl.WithTimeout(timeout))
	}

	if config.UserProjectOverride {
		configOptions = append(configOptions, dcl.WithUserProjectOverride())
		if billingProject != "" {
			configOptions = append(configOptions, dcl.WithBillingProject(billingProject))
		}
	}

	dclConfig := dcl.NewConfig(configOptions...)
	return NewClient(dclConfig)
}
