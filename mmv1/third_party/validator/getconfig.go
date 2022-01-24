package google

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/terraform-validator/version"
)

// Return the value of the private userAgent field
func (c *Config) UserAgent() string {
	return c.userAgent
}

func GetConfig(ctx context.Context, project string, offline bool) (*Config, error) {
	cfg := &Config{
		Project:   project,
		userAgent: fmt.Sprintf("config-validator-tf/%s", version.BuildVersion()),
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

	// opt in extension for adding to the User-Agent header
	if ext := os.Getenv("GOOGLE_TERRAFORM_VALIDATOR_USERAGENT_EXTENSION"); ext != "" {
		ua := cfg.userAgent
		cfg.userAgent = fmt.Sprintf("%s %s", ua, ext)
	}

	if !offline {
		ConfigureBasePaths(cfg)
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
	}

	return cfg, nil
}
