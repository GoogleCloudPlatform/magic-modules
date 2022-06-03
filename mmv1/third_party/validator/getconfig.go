package google

import (
	"context"

	"github.com/pkg/errors"
)

// Return the value of the private userAgent field
func (c *Config) UserAgent() string {
	return c.userAgent
}

func GetConfig(ctx context.Context, project string, offline bool, userAgent string) (*Config, error) {
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
	}

	return cfg, nil
}
