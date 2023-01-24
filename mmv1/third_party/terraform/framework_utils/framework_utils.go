package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

// MultiEnvDefaultFunc is a helper function that returns the value of the first
// environment variable in the given list that returns a non-empty value. If
// none of the environment variables return a value, the default value is
// returned.
func MultiEnvDefault(ks []string, dv interface{}) interface{} {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return dv
}

func CompileUserAgentString(name, tfVersion, provVersion string) string {
	ua := fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-SDK/%s %s/%s", tfVersion, "terraform-plugin-framework", name, provVersion)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			log.Printf("[DEBUG] Using modified User-Agent: %s", ua)
		}
	}

	return ua
}

func GetCurrUserEmail(p *frameworkProvider, userAgent string, diags *diag.Diagnostics) string {
	// When environment variables UserProjectOverride and BillingProject are set for the provider,
	// the header X-Goog-User-Project is set for the API requests.
	// But it causes an error when calling GetCurrUserEmail. Set the project to be "NO_BILLING_PROJECT_OVERRIDE".
	// And then it triggers the header X-Goog-User-Project to be set to empty string.

	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, d := sendFrameworkRequest(p, "GET", "NO_BILLING_PROJECT_OVERRIDE", "https://openidconnect.googleapis.com/v1/userinfo", userAgent, nil)
	diags.Append(d...)

	if diags.HasError() {
		log.Printf("[INFO] error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
		return ""
	}
	if res["email"] == nil {
		diags.AddError("error retrieving email from userinfo.", "email was nil in the response.")
		return ""
	}
	return res["email"].(string)
}

// getProject reads the "project" field from the given resource and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProjectFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return getProjectFromSchemaFramework("project", rVal, pVal, diags)
}

func getProjectFromSchemaFramework(projectSchemaField string, rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	if !rVal.IsNull() && rVal.String() != "" {
		return rVal
	}

	if !pVal.IsNull() && pVal.String() != "" {
		return pVal
	}

	diags.AddError("required field is not set", fmt.Sprintf("%s is not set", projectSchemaField))
	return types.String{}
}

func handleDatasourceNotFoundError(ctx context.Context, err error, state *tfsdk.State, resource string, diags *diag.Diagnostics) {
	if isGoogleApiErrorWithCode(err, 404) {
		log.Printf("[WARN] Removing %s because it's gone", resource)
		// The resource doesn't exist anymore
		state.RemoveResource(ctx)
	}

	diags.AddError(fmt.Sprintf("Error when reading or editing %s", resource), err.Error())
}
