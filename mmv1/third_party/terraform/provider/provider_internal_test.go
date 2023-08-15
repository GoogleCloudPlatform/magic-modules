package provider

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/hashicorp/magic-modules/mmv1/third_party/terraform/envvar"
	"github.com/hashicorp/magic-modules/mmv1/third_party/terraform/tpgresource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestProvider_validateCredentials(t *testing.T) {
	cases := map[string]struct {
		ConfigValue      func(t *testing.T) interface{}
		ValueNotProvided bool
		ExpectedWarnings []string
		ExpectedErrors   []error
	}{
		"configuring credentials as a path to a credentials JSON file is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return transport_tpg.TestFakeCredentialsPath // Path to a test fixture
			},
		},
		"configuring credentials as a path to a non-existant file is NOT valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return "./this/path/doesnt/exist.json" // Doesn't exist
			},
			ExpectedErrors: []error{
				// As the file doesn't exist, so the function attempts to parse it as a JSON
				errors.New("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
		"configuring credentials as a credentials JSON string is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				contents, err := ioutil.ReadFile(transport_tpg.TestFakeCredentialsPath)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}
				return string(contents)
			},
		},
		"configuring credentials as an empty string is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return ""
			},
		},
		"leaving credentials unconfigured is valid": {
			ValueNotProvided: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			var configValue interface{}
			if !tc.ValueNotProvided {
				configValue = tc.ConfigValue(t)
			}

			// Act
			// Note: second argument is currently unused by the function but is necessary to fulfill the SchemaValidateFunc type's function signature
			ws, es := validateCredentials(configValue, "")

			// Assert
			if len(ws) != len(tc.ExpectedWarnings) {
				t.Errorf("Expected %d warnings, got %d: %v", len(tc.ExpectedWarnings), len(ws), ws)
			}
			if len(es) != len(tc.ExpectedErrors) {
				t.Errorf("Expected %d errors, got %d: %v", len(tc.ExpectedErrors), len(es), es)
			}

			if len(tc.ExpectedErrors) > 0 {
				if es[0].Error() != tc.ExpectedErrors[0].Error() {
					t.Errorf("Expected first error to be \"%s\", got \"%s\"", tc.ExpectedErrors[0], es[0])
				}
			}
		})
	}
}

// ProviderConfigEnvNames returns a list of all the environment variables that could be set by a user to configure the provider
func ProviderConfigEnvNames() []string {

	envs := []string{}

	// Use existing collections of ENV names
	envVarsSets := [][]string{
		envvar.CredsEnvVars,   // credentials field
		envvar.ProjectEnvVars, // project field
		envvar.RegionEnvVars,  //region field
		envvar.ZoneEnvVars,    // zone field
	}
	for _, set := range envVarsSets {
		envs = append(envs, set...)
	}

	// Add remaining ENVs
	envs = append(envs, "GOOGLE_OAUTH_ACCESS_TOKEN")          // access_token field
	envs = append(envs, "GOOGLE_BILLING_PROJECT")             // billing_project field
	envs = append(envs, "GOOGLE_IMPERSONATE_SERVICE_ACCOUNT") // impersonate_service_account field
	envs = append(envs, "USER_PROJECT_OVERRIDE")              // user_project_override field
	envs = append(envs, "CLOUDSDK_CORE_REQUEST_REASON")       // request_reason field

	return envs
}

// unsetProviderConfigEnvs unsets any ENVs in the test environment that
// configure the provider.
// The testing package will restore the original values after the test
func unsetTestProviderConfigEnvs(t *testing.T) {
	envs := ProviderConfigEnvNames()
	if len(envs) > 0 {
		for _, k := range envs {
			t.Setenv(k, "")
		}
	}
}

func setupTestEnvs(t *testing.T, envValues map[string]string) {
	// Set ENVs
	if len(envValues) > 0 {
		for k, v := range envValues {
			t.Setenv(k, v)
		}
	}
}

// Returns a fake credentials JSON string with the client_email set to a test-specific value
func generateFakeCredentialsJson(testId string) string {
	json := fmt.Sprintf(`{"private_key_id": "foo","private_key": "bar","client_email": "%s@example.com","client_id": "id@foo.com","type": "service_account"}`, testId)
	return json
}

func TestProvider_providerConfigure_credentials(t *testing.T) {

	const pathToMissingFile string = "./this/path/doesnt/exist.json" // Doesn't exist

	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectError         bool
		ExpectFieldUnset    bool
		ExpectedSchemaValue string
		ExpectedConfigValue string
	}{
		"credentials can be configured as a path to a credentials JSON file": {
			ConfigValues: map[string]interface{}{
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables:        map[string]string{},
			ExpectedSchemaValue: transport_tpg.TestFakeCredentialsPath,
			ExpectedConfigValue: transport_tpg.TestFakeCredentialsPath,
		},
		"configuring credentials as a path to a non-existant file results in an error": {
			ConfigValues: map[string]interface{}{
				"credentials": pathToMissingFile,
			},
			ExpectError:         true,
			ExpectedSchemaValue: pathToMissingFile,
			ExpectedConfigValue: pathToMissingFile,
		},
		"credentials set in the config are not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"credentials": generateFakeCredentialsJson("test"),
			},
			EnvVariables: map[string]string{
				"GOOGLE_CREDENTIALS":             generateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
				"GOOGLE_CLOUD_KEYFILE_JSON":      generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: generateFakeCredentialsJson("test"),
			ExpectedConfigValue: generateFakeCredentialsJson("test"),
		},
		"when credentials is unset in the config, environment variables are used: GOOGLE_CREDENTIALS used first": {
			EnvVariables: map[string]string{
				"GOOGLE_CREDENTIALS":             generateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
				"GOOGLE_CLOUD_KEYFILE_JSON":      generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: generateFakeCredentialsJson("GOOGLE_CREDENTIALS"),
		},
		"when credentials is unset in the config, environment variables are used: GOOGLE_CLOUD_KEYFILE_JSON used second": {
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				"GOOGLE_CLOUD_KEYFILE_JSON":      generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: generateFakeCredentialsJson("GOOGLE_CLOUD_KEYFILE_JSON"),
		},
		"when credentials is unset in the config, environment variables are used: GCLOUD_KEYFILE_JSON used third": {
			EnvVariables: map[string]string{
				// GOOGLE_CREDENTIALS not set
				// GOOGLE_CLOUD_KEYFILE_JSON not set
				"GCLOUD_KEYFILE_JSON":            generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
				"GOOGLE_APPLICATION_CREDENTIALS": generateFakeCredentialsJson("GOOGLE_APPLICATION_CREDENTIALS"),
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: generateFakeCredentialsJson("GCLOUD_KEYFILE_JSON"),
		},
		"when credentials is unset in the config (and access_token unset), GOOGLE_APPLICATION_CREDENTIALS is used for auth but not to set values in the config": {
			EnvVariables: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
		},
		// Handling empty strings in config
		"when credentials is set to an empty string in the config (and access_token unset), GOOGLE_APPLICATION_CREDENTIALS is used": {
			ConfigValues: map[string]interface{}{
				"credentials": "",
			},
			EnvVariables: map[string]string{
				"GOOGLE_APPLICATION_CREDENTIALS": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
		},
		// Error states
		// NOTE: these tests can't run in Cloud Build due to ADC locating credentials despite `GOOGLE_APPLICATION_CREDENTIALS` being unset
		// See https://cloud.google.com/docs/authentication/application-default-credentials#search_order
		// Also, when running these tests locally you need to run `gcloud auth application-default revoke` to ensure your machine isn't supplying ADCs
		// "error returned if credentials is set as an empty string and GOOGLE_APPLICATION_CREDENTIALS is unset": {
		// 	ConfigValues: map[string]interface{}{
		// 		"credentials": "",
		// 	},
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": "", // setting to empty string to help test run in CI
		// 	},
		// 	ExpectError: true,
		// },
		// "error returned if neither credentials nor access_token set in the provider config, and GOOGLE_APPLICATION_CREDENTIALS is unset": {
		// 	EnvVariables: map[string]string{
		// 		"GOOGLE_APPLICATION_CREDENTIALS": "", // setting to empty string to help test run in CI
		// 	},
		// 	ExpectError: true,
		// },
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("credentials")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSchemaValue {
						t.Fatalf("expected credentials value set in provider config data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected credentials value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("credentials")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected credentials value to be unset in provider config data")
			}
			if v != tc.ExpectedSchemaValue {
				t.Fatalf("expected credentials value set in provider config data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.Credentials != tc.ExpectedConfigValue {
				t.Fatalf("expected credentials value set in Config struct to be to be %s, got %s", tc.ExpectedConfigValue, config.Credentials)
			}
		})
	}
}

func TestProvider_providerConfigure_accessToken(t *testing.T) {

	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectedSchemaValue string
		ExpectedConfigValue string
		ExpectError         bool
		ExpectFieldUnset    bool
	}{
		"access_token configured in the provider can be invalid without resulting in errors": {
			ConfigValues: map[string]interface{}{
				"access_token": "This is not a valid token string",
			},
			EnvVariables:        map[string]string{},
			ExpectedSchemaValue: "This is not a valid token string",
			ExpectedConfigValue: "This is not a valid token string",
		},
		"access_token set in the provider config is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"access_token": "value-from-config",
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-env",
			},
			ExpectedSchemaValue: "value-from-config",
			ExpectedConfigValue: "value-from-config",
		},
		"when access_token is unset in the config, an environment variable is used but doesn't update the schema data": {
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// access_token unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:         false,
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "",
		},
		// Handle empty strings in config
		"when access_token is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"access_token": "",
				"credentials":  testFakeCredentialsPath,
			},
			ExpectError:         false,
			ExpectFieldUnset:    true,
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "",
		},
		"when access_token is set as an empty string in the config, an environment variable is used but doesn't update the schema data": {
			ConfigValues: map[string]interface{}{
				"access_token": "",
			},
			EnvVariables: map[string]string{
				"GOOGLE_OAUTH_ACCESS_TOKEN": "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
			},
			ExpectedSchemaValue: "",
			ExpectedConfigValue: "value-from-GOOGLE_OAUTH_ACCESS_TOKEN",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("access_token")
				if ok {
					val := v.(string)
					if val != tc.ExpectedSchemaValue {
						t.Fatalf("expected access_token value set in provider config data to be %s, got %s", tc.ExpectedSchemaValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected access_token value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("access_token")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected access_token value to be unset in provider config data")
			}
			if val != tc.ExpectedSchemaValue {
				t.Fatalf("expected access_token value set in provider config data to be %s, got %s", tc.ExpectedSchemaValue, val)
			}
			if config.AccessToken != tc.ExpectedConfigValue {
				t.Fatalf("expected access_token value set in Config struct to be to be %s, got %s", tc.ExpectedConfigValue, config.AccessToken)
			}
		})
	}
}

func TestProvider_providerConfigure_impersonateServiceAccount(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"impersonate_service_account value set in the provider schema is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account": "value-from-config@example.com",
				"credentials":                 transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedValue: "value-from-config@example.com",
		},
		"impersonate_service_account value can be set by environment variable": {
			ConfigValues: map[string]interface{}{
				// impersonate_service_account unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_IMPERSONATE_SERVICE_ACCOUNT": "value-from-env@example.com",
			},
			ExpectedValue: "value-from-env@example.com",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// impersonate_service_account unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		// Handling empty strings in config
		"when impersonate_service_account is set as an empty array the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account": "",
				"credentials":                 testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("impersonate_service_account")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected impersonate_service_account value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected impersonate_service_account value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("impersonate_service_account")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected impersonate_service_account value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected impersonate_service_account value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.ImpersonateServiceAccount != tc.ExpectedValue {
				t.Fatalf("expected impersonate_service_account value in Config struct to be %s, got %s", tc.ExpectedValue, config.ImpersonateServiceAccount)
			}
		})
	}
}

func TestProvider_providerConfigure_impersonateServiceAccountDelegates(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    []string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"impersonate_service_account_delegates value can be set in the provider schema": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account_delegates": []string{
					"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
					"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
				},
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue: []string{
				"projects/-/serviceAccounts/my-service-account-1@example.iam.gserviceaccount.com",
				"projects/-/serviceAccounts/my-service-account-2@example.iam.gserviceaccount.com",
			},
		},
		// No environment variables can be used for impersonate_service_account_delegates
		"when no impersonate_service_account_delegates value is provided via config, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// impersonate_service_account_delegates unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    nil,
		},
		// Handling empty values in config
		"when project is set as an empty array the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"impersonate_service_account_delegates": []string{},
				"credentials":                           testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    nil,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("impersonate_service_account_delegates")
				if ok {
					val := v.([]interface{})
					if tc.ExpectFieldUnset {
						t.Fatalf("expected impersonate_service_account_delegates value to not be set in provider config data, got %#v", val)
					}
					if len(val) != len(tc.ExpectedValue) {
						t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
					}
					for i := 0; i < len(val); i++ {
						if val[i].(string) != tc.ExpectedValue[i] {
							t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
						}
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors
			v, ok := d.GetOk("impersonate_service_account_delegates")
			val := v.([]interface{})
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected impersonate_service_account_delegates value to be unset in provider config data")
			}
			if len(val) != len(tc.ExpectedValue) {
				t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
			}
			for i := 0; i < len(val); i++ {
				if val[i].(string) != tc.ExpectedValue[i] {
					t.Fatalf("expected impersonate_service_account_delegates value set in provider config data to be %#v, got %#v", tc.ExpectedValue, val)
				}
				if config.ImpersonateServiceAccountDelegates[i] != tc.ExpectedValue[i] {
					t.Fatalf("expected impersonate_service_account_delegates value set in Config struct to be to be %#v, got %#v", tc.ExpectedValue, config.ImpersonateServiceAccountDelegates)
				}
			}
		})
	}
}

func TestProvider_providerConfigure_project(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"project value set in the provider schema is not overridden by environment variables": {
			ConfigValues: map[string]interface{}{
				"project":     "my-project-from-config",
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "my-project-from-config",
		},
		"project value can be set by environment variable: GOOGLE_PROJECT is used first": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-GOOGLE_PROJECT",
		},
		"project value can be set by environment variable: GOOGLE_CLOUD_PROJECT is used second": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				"GOOGLE_CLOUD_PROJECT":  "project-from-GOOGLE_CLOUD_PROJECT",
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-GOOGLE_CLOUD_PROJECT",
		},
		"project value can be set by environment variable: GCLOUD_PROJECT is used third": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				// GOOGLE_CLOUD_PROJECT unset
				"GCLOUD_PROJECT":        "project-from-GCLOUD_PROJECT",
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-GCLOUD_PROJECT",
		},
		"project value can be set by environment variable: CLOUDSDK_CORE_PROJECT is used fourth": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_PROJECT unset
				// GOOGLE_CLOUD_PROJECT unset
				// GCLOUD_PROJECT unset
				"CLOUDSDK_CORE_PROJECT": "project-from-CLOUDSDK_CORE_PROJECT",
			},
			ExpectedValue: "project-from-CLOUDSDK_CORE_PROJECT",
		},
		"when no project values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// project unset
				"credentials": transport_tpg.TestFakeCredentialsPath,
			},
			ExpectedValue: "",
		},
		// Handling empty strings in config
		"when project is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"project":     "",
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when project is set as an empty string an environment variable will be used": {
			ConfigValues: map[string]interface{}{
				"project":     "",
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_PROJECT":        "project-from-GOOGLE_PROJECT",
			}
			ExpectedValue:    "project-from-GOOGLE_PROJECT",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx := context.Background()
			unsetTestProviderConfigEnvs(t)
			setupTestEnvs(t, tc.EnvVariables)
			p := Provider()
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, p.Schema, tc.ConfigValues)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("project")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected project value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("project")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected project value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.Project != tc.ExpectedValue {
				t.Fatalf("expected project value set in Config struct to be to be %s, got %s", tc.ExpectedValue, config.Project)
			}
		})
	}
}

func TestProvider_providerConfigure_billingProject(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"billing_project value set in the provider config is not overridden by ENVs": {
			ConfigValues: map[string]interface{}{
				"billing_project": "my-billing-project-from-config",
				"credentials":     testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "my-billing-project-from-env",
			},
			ExpectedValue: "my-billing-project-from-config",
		},
		"billing project can be set by environment variable, when no value supplied via the config": {
			ConfigValues: map[string]interface{}{
				// billing_project unset
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "my-billing-project-from-env",
			},
			ExpectedValue: "my-billing-project-from-env",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// billing_project unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		// Handling empty strings in config
		"when billing_project is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"billing_project": "",
				"credentials":     testFakeCredentialsPath,
			},
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when billing_project is set as an empty string an environment variable will be used": {
			ConfigValues: map[string]interface{}{
				"billing_project":     "",
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_BILLING_PROJECT": "my-billing-project-from-env",
			},
			ExpectedValue: "my-billing-project-from-env",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx, p, d := setupSDKProviderConfigTest(t, tc.ConfigValues, tc.EnvVariables)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("billing_project")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected billing_project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected billing_project value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("billing_project")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected billing_project value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected billing_project value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.BillingProject != tc.ExpectedValue {
				t.Fatalf("expected billing_project value set in Config struct to be to be %s, got %s", tc.ExpectedValue, config.BillingProject)
			}
		})
	}
}

func TestProvider_providerConfigure_region(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"region value set in the provider config is not overridden by ENVs": {
			ConfigValues: map[string]interface{}{
				"region":      "my-region-from-config",
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedValue: "my-region-from-config",
		},
		"region value can be set by environment variable: GOOGLE_REGION is used": {
			ConfigValues: map[string]interface{}{
				// region unset
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedValue: "region-from-env",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// region unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		// Handling empty strings in config
		"when region is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"region":      "",
				"credentials": testFakeCredentialsPath,
			},
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when region is set as an empty string an environment variable will be used": {
			ConfigValues: map[string]interface{}{
				"region":     "",
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_REGION": "region-from-env",
			},
			ExpectedValue: "region-from-env",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx, p, d := setupSDKProviderConfigTest(t, tc.ConfigValues, tc.EnvVariables)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("region")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected region value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected region value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("region")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected region value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected region value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.Region != tc.ExpectedValue {
				t.Fatalf("expected region value set in Config struct to be to be %s, got %s", tc.ExpectedValue, config.Region)
			}
		})
	}
}

func TestProvider_providerConfigure_zone(t *testing.T) {

	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    string
		ExpectError      bool
		ExpectFieldUnset bool
	}{
		"zone value set in the provider config is not overridden by ENVs": {
			ConfigValues: map[string]interface{}{
				"zone":        "zone-from-config",
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE": "zone-from-env",
			},
			ExpectedValue: "zone-from-config",
		},
		"when multiple zone environment variables are provided, `GOOGLE_ZONE` is used first": {
			ConfigValues: map[string]interface{}{
				// zone unset,
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE":           "zone-from-GOOGLE_ZONE",
				"GCLOUD_ZONE":           "zone-from-GCLOUD_ZONE",
				"CLOUDSDK_COMPUTE_ZONE": "zone-from-CLOUDSDK_COMPUTE_ZONE",
			},
			ExpectedValue: "zone-from-GOOGLE_ZONE",
		},
		"when multiple zone environment variables are provided, `GCLOUD_ZONE` is used second": {
			ConfigValues: map[string]interface{}{
				// zone unset,
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_ZONE unset
				"GCLOUD_ZONE":           "zone-from-GCLOUD_ZONE",
				"CLOUDSDK_COMPUTE_ZONE": "zone-from-CLOUDSDK_COMPUTE_ZONE",
			},
			ExpectedValue: "zone-from-GCLOUD_ZONE",
		},
		"when multiple zone environment variables are provided, `CLOUDSDK_COMPUTE_ZONE` is used third": {
			ConfigValues: map[string]interface{}{
				// zone unset,
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				// GOOGLE_ZONE unset
				// GCLOUD_ZONE unset
				"CLOUDSDK_COMPUTE_ZONE": "zone-from-CLOUDSDK_COMPUTE_ZONE",
			},
			ExpectedValue: "zone-from-CLOUDSDK_COMPUTE_ZONE",
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// zone unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		// Handling empty strings in config
		"when zone is set as an empty string the field is treated as if it's unset, without error": {
			ConfigValues: map[string]interface{}{
				"zone":        "",
				"credentials": testFakeCredentialsPath,
			},
			ExpectFieldUnset: true,
			ExpectedValue:    "",
		},
		"when zone is set as an empty string an environment variable will be used": {
			ConfigValues: map[string]interface{}{
				"zone":     "",
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"GOOGLE_ZONE": "zone-from-env",
			},
			ExpectedValue: "zone-from-env",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx, p, d := setupSDKProviderConfigTest(t, tc.ConfigValues, tc.EnvVariables)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("zone")
				if ok {
					val := v.(string)
					if val != tc.ExpectedValue {
						t.Fatalf("expected zone value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
					}
					if tc.ExpectFieldUnset {
						t.Fatalf("expected zone value to not be set in provider config data, got %s", val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("zone")
			val := v.(string)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected zone value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected zone value set in provider config data to be %s, got %s", tc.ExpectedValue, val)
			}
			if config.Zone != tc.ExpectedValue {
				t.Fatalf("expected zone value set in Config struct to be to be %s, got %s", tc.ExpectedValue, config.Zone)
			}
		})
	}
}

func TestProvider_providerConfigure_userProjectOverride(t *testing.T) {
	cases := map[string]struct {
		ConfigValues     map[string]interface{}
		EnvVariables     map[string]string
		ExpectedValue    bool
		ExpectFieldUnset bool
		ExpectError      bool
	}{
		"user_project_override value set in the provider schema is not overridden by ENVs": {
			ConfigValues: map[string]interface{}{
				"user_project_override": false,
				"credentials":           testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedValue: false,
		},
		"user_project_override can be set by environment variable: true": {
			ConfigValues: map[string]interface{}{
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "true",
			},
			ExpectedValue: true,
		},
		"user_project_override can be set by environment variable: false": {
			ConfigValues: map[string]interface{}{
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "false",
			},
			ExpectedValue: false,
		},
		"user_project_override can be set by environment variable: 1": {
			ConfigValues: map[string]interface{}{
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "1",
			},
			ExpectedValue: true,
		},
		"user_project_override can be set by environment variable: 0": {
			ConfigValues: map[string]interface{}{
				"credentials": testFakeCredentialsPath,
			},
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "0",
			},
			ExpectedValue: false,
		},
		"error returned due to non-boolean environment variables": {
			EnvVariables: map[string]string{
				"USER_PROJECT_OVERRIDE": "I'm not a boolean",
			},
			ExpectError: true,
		},
		"when no values are provided via config or environment variables, the field remains unset without error": {
			ConfigValues: map[string]interface{}{
				// user_project_override unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:      false,
			ExpectFieldUnset: true,
			ExpectedValue:    false,
		},
		// There isn't an equivalent test case for 'user sets value as empty string' because user_project_override is a boolean; true/false both valid.
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx, p, d := setupSDKProviderConfigTest(t, tc.ConfigValues, tc.EnvVariables)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("user_project_override")
				if ok {
					val := v.(bool)
					if tc.ExpectFieldUnset {
						t.Fatalf("expected user_project_override value to not be set in provider config data, got %v", val)
					}
					if val != tc.ExpectedValue {
						t.Fatalf("expected user_project_override value set in provider config data to be %v, got %v", tc.ExpectedValue, val)
					}
				}
				// Return early in tests where errors expected
				return
			}

			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors

			v, ok := d.GetOk("user_project_override")
			val := v.(bool)
			if ok && tc.ExpectFieldUnset {
				t.Fatal("expected user_project_override value to be unset in provider config data")
			}
			if val != tc.ExpectedValue {
				t.Fatalf("expected user_project_override value set in provider config data to be %v, got %v", tc.ExpectedValue, val)
			}
			if config.UserProjectOverride != tc.ExpectedValue {
				t.Fatalf("expected user_project_override value set in Config struct to be to be %v, got %v", tc.ExpectedValue, config.UserProjectOverride)
			}
		})
	}
}

func TestProvider_providerConfigure_scopes(t *testing.T) {
	cases := map[string]struct {
		ConfigValues        map[string]interface{}
		EnvVariables        map[string]string
		ExpectedSchemaValue []string
		ExpectedConfigValue []string
		ExpectFieldUnset    bool
		ExpectError         bool
	}{
		"scopes are set in the provider config as a list": {
			ConfigValues: map[string]interface{}{
				"credentials": testFakeCredentialsPath,
				"scopes": []string{
					"fizz",
					"buzz",
					"fizzbuzz",
				},
			},
			ExpectedSchemaValue: []string{
				"fizz",
				"buzz",
				"fizzbuzz",
			},
			ExpectedConfigValue: []string{
				"fizz",
				"buzz",
				"fizzbuzz",
			},
		},
		"scopes can be left unset in the provider config without any issues, and a default value is used": {
			ConfigValues: map[string]interface{}{
				// scopes unset
				"credentials": testFakeCredentialsPath,
			},
			ExpectedSchemaValue: nil,
			ExpectedConfigValue: transport_tpg.DefaultClientScopes,
		},
		// Handling empty values in config
		"scopes set as an empty list the field is treated as if it's unset and a default value is used without errors": {
			ConfigValues: map[string]interface{}{
				"scopes":      []string{},
				"credentials": testFakeCredentialsPath,
			},
			ExpectError:         false,
			ExpectFieldUnset:    true, //unset in provider config data, not the subsequent Config struct
			ExpectedSchemaValue: nil,
			ExpectedConfigValue: transport_tpg.DefaultClientScopes,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Arrange
			ctx, p, d := setupSDKProviderConfigTest(t, tc.ConfigValues, tc.EnvVariables)

			// Act
			c, diags := providerConfigure(ctx, d, p)

			// Assert
			if diags.HasError() && !tc.ExpectError {
				t.Fatalf("unexpected error(s): %#v", diags)
			}
			if !diags.HasError() && tc.ExpectError {
				t.Fatal("expected error(s) but got none")
			}
			if diags.HasError() && tc.ExpectError {
				v, ok := d.GetOk("scopes")
				if ok {
					val := v.([]interface{})
					if tc.ExpectFieldUnset {
						t.Fatalf("expected scopes value to not be set in provider config data, got %#v", val)
					}
					if len(val) != len(tc.ExpectedSchemaValue) {
						t.Fatalf("expected scopes value set in provider config data to be %#v, got %#v", tc.ExpectedSchemaValue, val)
					}
					for i := 0; i < len(val); i++ {
						if val[i].(string) != tc.ExpectedSchemaValue[i] {
							t.Fatalf("expected scopes value set in provider config data to be %#v, got %#v", tc.ExpectedSchemaValue, val)
						}
					}
				}
				// Return early in tests where errors expected
				return
			}
			config := c.(*transport_tpg.Config) // Should be non-nil value, as test cases reaching this point experienced no errors
			v, ok := d.GetOk("scopes")
			if ok {
				val := v.([]interface{})

				if len(val) != len(tc.ExpectedSchemaValue) {
					t.Fatalf("expected %v scopes set in provider config data, got %v", len(tc.ExpectedSchemaValue), len(val))
				}
				for i, el := range val {
					scope := el.(string)
					if scope != tc.ExpectedSchemaValue[i] {
						t.Fatalf("expected scopes value set in provider config data to be %v, got %v", tc.ExpectedSchemaValue, val)
					}
				}
			}
			if !ok && (len(tc.ExpectedSchemaValue) > 0) {
				t.Fatalf("expected %v scopes to be set in the provider data, but is unset", tc.ExpectedSchemaValue)
			}

			if len(config.Scopes) != len(tc.ExpectedConfigValue) {
				t.Fatalf("expected %v scopes set in the config struct, got %v", len(tc.ExpectedConfigValue), len(config.Scopes))
			}
			for i, el := range config.Scopes {
				if el != tc.ExpectedConfigValue[i] {
					t.Fatalf("expected scopes value set in provider config data to be %v, got %v", tc.ExpectedConfigValue, config.Scopes)
				}
			}
		})
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("https://www.googleapis.com/compute/beta/", RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", RandString(t, 10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccProviderMeta_setModuleName(t *testing.T) {
	t.Parallel()

	moduleName := "my-module"
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderMeta_setModuleName(moduleName, RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billing := acctest.GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + RandString(t, 10)
	topicName := "tf-test-topic-" + RandString(t, 10)

	config := BootstrapConfig(t)
	accessToken, err := setupProjectsAndGetAccessToken(org, billing, pid, "pubsub", config)
	if err != nil {
		t.Error(err)
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderUserProjectOverride_step2(accessToken, pid, false, topicName),
				ExpectError: regexp.MustCompile("Cloud Pub/Sub API has not been used"),
			},
			{
				Config: testAccProviderUserProjectOverride_step2(accessToken, pid, true, topicName),
			},
			{
				ResourceName:      "google_pubsub_topic.project-2-topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderUserProjectOverride_step3(accessToken, true),
			},
		},
	})
}

// Do the same thing as TestAccProviderUserProjectOverride, but using a resource that gets its project via
// a reference to a different resource instead of a project field.
func TestAccProviderIndirectUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := acctest.GetTestOrgFromEnv(t)
	billing := acctest.GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + RandString(t, 10)

	config := BootstrapConfig(t)
	accessToken, err := setupProjectsAndGetAccessToken(org, billing, pid, "cloudkms", config)
	if err != nil {
		t.Error(err)
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderIndirectUserProjectOverride_step2(pid, accessToken, false),
				ExpectError: regexp.MustCompile(`Cloud Key Management Service \(KMS\) API has not been used`),
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step2(pid, accessToken, true),
			},
			{
				ResourceName:      "google_kms_crypto_key.project-2-key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step3(accessToken, true),
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                   = "compute_custom_endpoint"
  compute_custom_endpoint = "%s"
}

resource "google_compute_address" "default" {
  provider = google.compute_custom_endpoint
  name     = "tf-test-address-%s"
}`, endpoint, name)
}

func testAccProviderMeta_setModuleName(key, name string) string {
	return fmt.Sprintf(`
terraform {
  provider_meta "google" {
    module_name = "%s"
  }
}

resource "google_compute_address" "default" {
	name = "tf-test-address-%s"
}`, key, name)
}

// Set up two projects. Project 1 has a service account that is used to create a
// pubsub topic in project 2. The pubsub API is only enabled in project 2,
// which causes the create to fail unless user_project_override is set to true.

func testAccProviderUserProjectOverride_step2(accessToken, pid string, override bool, topicName string) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the pubsub topic.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// pubsub topic so the whole config can be deleted.
%s

resource "google_pubsub_topic" "project-2-topic" {
	provider = google.project-1-token
	project  = "%s-2"

	name = "%s"
	labels = {
	  foo = "bar"
	}
}
`, testAccProviderUserProjectOverride_step3(accessToken, override), pid, topicName)
}

func testAccProviderUserProjectOverride_step3(accessToken string, override bool) string {
	return fmt.Sprintf(`
provider "google" {
	alias  = "project-1-token"
	access_token = "%s"
	user_project_override = %v
}
`, accessToken, override)
}

func testAccProviderIndirectUserProjectOverride_step2(pid, accessToken string, override bool) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the kms resources.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// kms resources so the whole config can be deleted.
%s

resource "google_kms_key_ring" "project-2-keyring" {
	provider = google.project-1-token
	project  = "%s-2"

	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "project-2-key" {
	provider = google.project-1-token
	name     = "%s"
	key_ring = google_kms_key_ring.project-2-keyring.id
}

data "google_kms_secret_ciphertext" "project-2-ciphertext" {
	provider   = google.project-1-token
	crypto_key = google_kms_crypto_key.project-2-key.id
	plaintext  = "my-secret"
}
`, testAccProviderIndirectUserProjectOverride_step3(accessToken, override), pid, pid, pid)
}

func testAccProviderIndirectUserProjectOverride_step3(accessToken string, override bool) string {
	return fmt.Sprintf(`
provider "google" {
	alias = "project-1-token"

	access_token          = "%s"
	user_project_override = %v
}
`, accessToken, override)
}
