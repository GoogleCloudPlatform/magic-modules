package google

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
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
