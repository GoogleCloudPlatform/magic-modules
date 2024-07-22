package google

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func NewConfig(ctx context.Context, project, zone, region string, offline bool, userAgent string, client *http.Client) (*transport_tpg.Config, error) {
	cfg := &transport_tpg.Config{
		Project:   project,
		Zone:      zone,
		Region:    region,
		UserAgent: userAgent,
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
		if client != nil {
			cfg.Client = client
		}
	}

	return cfg, nil
}
