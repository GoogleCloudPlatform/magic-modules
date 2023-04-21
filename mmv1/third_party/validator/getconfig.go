package google

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v2/tfplan2cai/converters/google/resources/transport"
)

// Return the value of the private userAgent field
func (c *transport_tpg.Config) GetUserAgent() string {
	return c.UserAgent
}

// Return the value of the private client field
func (c *transport_tpg.Config) GetClient() *http.Client {
	return c.Client
}

func NewConfig(ctx context.Context, project, zone, region string, offline bool, userAgent string, client *http.Client) (*transport_tpg.Config, error) {
	cfg := &Config{
		Project:   project,
		Zone:      zone,
		Region:    region,
		UserAgent: userAgent,
	}

	// Search for default credentials
	cfg.Credentials = MultiEnvSearch([]string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_CLOUD_KEYFILE_JSON",
		"GCLOUD_KEYFILE_JSON",
	})

	cfg.AccessToken = MultiEnvSearch([]string{
		"GOOGLE_OAUTH_ACCESS_TOKEN",
	})

	cfg.ImpersonateServiceAccount = MultiEnvSearch([]string{
		"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
	})

	if !offline {
		ConfigureBasePaths(cfg)
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
		if client != nil {
			cfg.Client = client
		}
	}

	return cfg, nil
}
