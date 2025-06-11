package transport

import (
	"context"

	"github.com/pkg/errors"

	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v6/pkg/transport"
)

func NewConfig(ctx context.Context, project, zone, region string, offline bool, userAgent string) (*transport_tpg.Config, error) {
	cfg := &transport_tpg.Config{
		Project:   project,
		Zone:      zone,
		Region:    region,
		UserAgent: userAgent,
	}

	if cfg.Project == "" {
		cfg.Project = transport_tpg.MultiEnvSearch([]string{
			"GOOGLE_PROJECT",
			"GOOGLE_CLOUD_PROJECT",
			"GCLOUD_PROJECT",
			"CLOUDSDK_CORE_PROJECT",
		})
	}

	// Search for default credentials
	cfg.Credentials = transport_tpg.MultiEnvSearch([]string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_CLOUD_KEYFILE_JSON",
		"GCLOUD_KEYFILE_JSON",
	})

	cfg.AccessToken = transport_tpg.MultiEnvSearch([]string{
		"GOOGLE_OAUTH_ACCESS_TOKEN",
	})

	cfg.ImpersonateServiceAccount = transport_tpg.MultiEnvSearch([]string{
		"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
	})

	transport_tpg.ConfigureBasePaths(cfg)
	if !offline {
		if err := cfg.LoadAndValidate(ctx); err != nil {
			return nil, errors.Wrap(err, "load and validate config")
		}
	}

	return cfg, nil
}
