package google

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

type configAttrGetter func(cfg *transport_tpg.Config) string

func getCredentials(cfg *transport_tpg.Config) string {
	return cfg.Credentials
}

func getAccessToken(cfg *transport_tpg.Config) string {
	return cfg.AccessToken
}

func getImpersonateServiceAccount(cfg *transport_tpg.Config) string {
	return cfg.ImpersonateServiceAccount
}

func TestNewConfigExtractsEnvVars(t *testing.T) {
	ctx := context.Background()
	offline := true
	cases := []struct {
		name           string
		envKey         string
		unsetKeys      []string // environment variables that must be unset before constructing config
		envValue       string
		expected       string
		getConfigValue configAttrGetter
	}{
		{
			name:           "GOOGLE_CREDENTIALS",
			envKey:         "GOOGLE_CREDENTIALS",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getCredentials,
		},
		{
			name:           "GOOGLE_CLOUD_KEYFILE_JSON",
			envKey:         "GOOGLE_CLOUD_KEYFILE_JSON",
			unsetKeys:      []string{"GOOGLE_CREDENTIALS"},
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getCredentials,
		},
		{
			name:           "GCLOUD_KEYFILE_JSON",
			envKey:         "GCLOUD_KEYFILE_JSON",
			unsetKeys:      []string{"GOOGLE_CREDENTIALS", "GOOGLE_CLOUD_KEYFILE_JSON"},
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getCredentials,
		},
		{
			name:           "GOOGLE_OAUTH_ACCESS_TOKEN",
			envKey:         "GOOGLE_OAUTH_ACCESS_TOKEN",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getAccessToken,
		},
		{
			name:           "GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
			envKey:         "GOOGLE_IMPERSONATE_SERVICE_ACCOUNT",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getImpersonateServiceAccount,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Store existing state of environment variables.
			existingEnv := make(map[string]string, len(c.unsetKeys)+1)
			for _, key := range c.unsetKeys {
				if originalValue, isSet := os.LookupEnv(key); isSet {
					existingEnv[key] = originalValue
					// Unset variables that would interfere with test.
					err := os.Unsetenv(key)
					if err != nil {
						t.Fatalf("error unsetting env var %s: %s", key, err)
					}
				}
			}
			originalValue, isSet := os.LookupEnv(c.envKey)
			if isSet {
				existingEnv[c.envKey] = originalValue
			}
			err := os.Setenv(c.envKey, c.envValue)
			if err != nil {
				t.Fatalf("error setting env var %s=%s: %s", c.envKey, c.envValue, err)
			}

			cfg, err := NewConfig(ctx, "project", "", "", offline, "", nil)
			if err != nil {
				t.Fatalf("error building converter: %s", err)
			}

			assert.Equal(t, c.expected, c.getConfigValue(cfg))

			// Restore previous states of environment variables.
			if !isSet {
				// c.envKey was previously unset.
				err = os.Unsetenv(c.envKey)
				if err != nil {
					t.Fatalf("error unsetting env var %s: %s", c.envKey, err)
				}
			}
			for key, value := range existingEnv {
				err = os.Setenv(key, value)
				if err != nil {
					t.Fatalf("error setting env var %s=%s: %s", key, value, err)
				}
			}
		})
	}
}

func TestNewConfigUserAgent(t *testing.T) {
	ctx := context.Background()
	offline := true
	cases := []struct {
		userAgent string
		expected  string
	}{
		{
			userAgent: "",
			expected:  "",
		},
		{
			userAgent: "whatever",
			expected:  "whatever",
		},
	}

	for _, c := range cases {
		t.Run(c.userAgent, func(t *testing.T) {
			cfg, err := NewConfig(ctx, "project", "", "", offline, c.userAgent, nil)
			if err != nil {
				t.Fatalf("error building config: %s", err)
			}

			assert.Equal(t, c.expected, cfg.UserAgent)
		})
	}
}

func TestNewConfigUserAgent_nilClientUsesDefault(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}
	ctx := context.Background()
	offline := false
	cfg, err := NewConfig(ctx, "project", "", "", offline, "", nil)
	if err != nil {
		t.Fatalf("error building config: %s", err)
	}

	assert.NotEmpty(t, cfg.Client)
}

func TestNewConfigUserAgent_usesPassedClient(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode.")
		return
	}
	ctx := context.Background()
	offline := false
	client := &http.Client{}
	cfg, err := NewConfig(ctx, "project", "", "", offline, "", client)
	if err != nil {
		t.Fatalf("error building config: %s", err)
	}

	assert.Exactly(t, client, cfg.Client)
}
