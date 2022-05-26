package google

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type configAttrGetter func(cfg *Config) string

func getCredentials(cfg *Config) string {
	return cfg.Credentials
}
func getAccessToken(cfg *Config) string {
	return cfg.AccessToken
}
func getImpersonateServiceAccount(cfg *Config) string {
	return cfg.ImpersonateServiceAccount
}
func getUserAgent(cfg *Config) string {
	return cfg.UserAgent()
}
func getZoneValue(cfg *Config) string {
	return cfg.Zone
}
func getRegionValue(cfg *Config) string {
	return cfg.Region
}

func TestGetConfigExtractsEnvVars(t *testing.T) {
	ctx := context.Background()
	offline := true
	cases := []struct {
		name           string
		envKey         string
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
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getCredentials,
		},
		{
			name:           "GCLOUD_KEYFILE_JSON",
			envKey:         "GCLOUD_KEYFILE_JSON",
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
		{
			name:            "GOOGLE_ZONE",
		        envKey:          "GOOGLE_ZONE",
		        envValue:        "whatever",
		        expected:        "whatever",
		        getConfigValue:  getZoneValue,
	        },
		{
			name:           "GCLOUD_ZONE",
			envKey:         "GCLOUD_ZONE",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getZoneValue,
		},
		{
			name:           "CLOUDSDK_COMPUTE_ZONE",
			envKey:         "CLOUDSDK_COMPUTE_ZONE",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getZoneValue,
		},
		{
			name:           "GOOGLE_REGION",
			envKey:         "GOOGLE_REGION",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getRegionValue,
		},
		{
			name:           "GCLOUD_REGION",
			envKey:         "GCLOUD_REGION",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getRegionValue,
		},
		{
			name:           "CLOUDSDK_COMPUTE_REGION",
			envKey:         "CLOUDSDK_COMPUTE_REGION",
			envValue:       "whatever",
			expected:       "whatever",
			getConfigValue: getRegionValue,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			originalValue, isSet := os.LookupEnv(c.envKey)
			err := os.Setenv(c.envKey, c.envValue)
			if err != nil {
				t.Fatalf("error setting env var %s=%s: %s", c.envKey, c.envValue, err)
			}

			cfg, err := GetConfig(ctx, "project", offline, "")
			if err != nil {
				t.Fatalf("error building converter: %s", err)
			}

			assert.Equal(t, c.expected, c.getConfigValue(cfg))

			if isSet {
				err = os.Setenv(c.envKey, originalValue)
				if err != nil {
					t.Fatalf("error setting env var %s=%s: %s", c.envKey, originalValue, err)
				}
			} else {
				err = os.Unsetenv(c.envKey)
				if err != nil {
					t.Fatalf("error unsetting env var %s: %s", c.envKey, err)
				}
			}
		})
	}
}

func TestGetConfigUserAgent(t *testing.T) {
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
			cfg, err := GetConfig(ctx, "project", offline, c.userAgent)
			if err != nil {
				t.Fatalf("error building converter: %s", err)
			}

			assert.Equal(t, c.expected, getUserAgent(cfg))
		})
	}
}
