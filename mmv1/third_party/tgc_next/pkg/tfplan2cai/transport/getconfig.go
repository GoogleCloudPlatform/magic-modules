package transport

import (
	"context"

	"net/http"

	"github.com/pkg/errors"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/envvar"
	transport_tpg "github.com/GoogleCloudPlatform/terraform-google-conversion/v7/pkg/transport"
)

func NewConfig(ctx context.Context, project, zone, region string, offline bool, userAgent string, client *http.Client) (*transport_tpg.Config, error) {
	cfg := &transport_tpg.Config{
		Project:   project,
		Zone:      zone,
		Region:    region,
		UserAgent: userAgent,
	}

	if cfg.Project == "" {
		cfg.Project = envvar.MultiEnvSearch([]string{
			"GOOGLE_PROJECT",
			"GOOGLE_CLOUD_PROJECT",
			"GCLOUD_PROJECT",
			"CLOUDSDK_CORE_PROJECT",
		})
	}

	// Search for default credentials
	cfg.Credentials = envvar.MultiEnvSearch([]string{
		"GOOGLE_CREDENTIALS",
		"GOOGLE_CLOUD_KEYFILE_JSON",
		"GCLOUD_KEYFILE_JSON",
	})

	cfg.AccessToken = envvar.MultiEnvSearch([]string{
		"GOOGLE_OAUTH_ACCESS_TOKEN",
	})

	cfg.ImpersonateServiceAccount = envvar.MultiEnvSearch([]string{
		"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
	})

	transport_tpg.ConfigureBasePaths(cfg)
	if !offline {
		if client != nil {
			cfg.Client = client
		} else {
			if err := cfg.LoadAndValidate(ctx); err != nil {
				return nil, errors.Wrap(err, "load and validate config")
			}
		}
	}

	return cfg, nil
}
