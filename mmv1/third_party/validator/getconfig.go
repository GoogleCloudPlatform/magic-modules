package google

import (
	"context"

	"github.com/pkg/errors"
)

func GetConfig(ctx context.Context, project string, offline bool) (*Config, error) {
	cfg := &Config{
		Project: project,
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

	if !offline {
		ConfigureBasePaths(cfg)
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
	}

	return cfg, nil
}
