package google

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

// Compatibility shim to let this change happen in two commits.
// NewConfig better matches golang best practices.
func GetConfig(ctx context.Context, project string, offline bool, userAgent string) (*Config, error) {
	return NewConfig(ctx, project, offline, userAgent, nil)
}

// Return the value of the private userAgent field
func (c *Config) UserAgent() string {
	return c.userAgent
}

// Return the value of the private client field
func (c *Config) Client() *http.Client {
	return c.client
}

func NewConfig(ctx context.Context, project string, offline bool, userAgent string, client *http.Client) (*Config, error) {
	cfg := &Config{
		Project:   project,
		userAgent: userAgent,
	}

	// Search for default credentials
	cfg.Credentials = multiEnvSearch([]string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_CLOUD_KEYFILE_JSON",
		"GCLOUD_KEYFILE_JSON",
	})

	cfg.AccessToken = multiEnvSearch([]string{
		"GOOGLE_OAUTH_ACCESS_TOKEN",
	})

	cfg.ImpersonateServiceAccount = multiEnvSearch([]string{
		"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
	})

	cfg.Zone = multiEnvSearch([]string{
		"GOOGLE_ZONE",
		"GCLOUD_ZONE",
		"CLOUDSDK_COMPUTE_ZONE",
	})

	cfg.Region = multiEnvSearch([]string{
		"GOOGLE_REGION",
		"GCLOUD_REGION",
		"CLOUDSDK_COMPUTE_REGION",
	})

	if !offline {
		ConfigureBasePaths(cfg)
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
		if client != nil {
			cfg.client = client
		}
	}

	return cfg, nil
}
