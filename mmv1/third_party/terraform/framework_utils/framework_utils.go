package google

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

func CompileUserAgentString(ctx context.Context, name, tfVersion, provVersion string) string {
	ua := fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-SDK/%s %s/%s", tfVersion, "terraform-plugin-framework", name, provVersion)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			tflog.Debug(ctx, fmt.Sprintf("Using modified User-Agent: %s", ua))
		}
	}

	return ua
}

func GetCurrentUserEmailFramework(p *frameworkProvider, userAgent string, diags *diag.Diagnostics) string {
	// When environment variables UserProjectOverride and BillingProject are set for the provider,
	// the header X-Goog-User-Project is set for the API requests.
	// But it causes an error when calling GetCurrUserEmail. Set the project to be "NO_BILLING_PROJECT_OVERRIDE".
	// And then it triggers the header X-Goog-User-Project to be set to empty string.

	// See https://github.com/golang/oauth2/issues/306 for a recommendation to do this from a Go maintainer
	// URL retrieved from https://accounts.google.com/.well-known/openid-configuration
	res, d := sendFrameworkRequest(p, "GET", "NO_BILLING_PROJECT_OVERRIDE", "https://openidconnect.googleapis.com/v1/userinfo", userAgent, nil)
	diags.Append(d...)

	if diags.HasError() {
		tflog.Info(p.context, "error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
		return ""
	}
	if res["email"] == nil {
		diags.AddError("error retrieving email from userinfo.", "email was nil in the response.")
		return ""
	}
	return res["email"].(string)
}

func generateFrameworkUserAgentString(metaData *ProviderMetaModel, currUserAgent string) string {
	if metaData != nil && !metaData.ModuleName.IsNull() && metaData.ModuleName.ValueString() != "" {
		return strings.Join([]string{currUserAgent, metaData.ModuleName.ValueString()}, " ")
	}

	return currUserAgent
}

// getProject reads the "project" field from the given resource and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProjectFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	return getProjectFromFrameworkSchema("project", rVal, pVal, diags)
}

func getProjectFromFrameworkSchema(projectSchemaField string, rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	if !rVal.IsNull() && rVal.ValueString() != "" {
		return rVal
	}

	if !pVal.IsNull() && pVal.ValueString() != "" {
		return pVal
	}

	diags.AddError("required field is not set", fmt.Sprintf("%s is not set", projectSchemaField))
	return types.String{}
}

func getRegionFramework(rRegion, rZone, pRegion, pZone types.String, diags *diag.Diagnostics) types.String {
	return getRegionFromFrameworkSchema("region", rRegion, rZone, pRegion, pZone, diags)
}

// Infers the region based on the following (in order of priority):
// - `regionSchemaField` in resource schema
// - region extracted from the `zoneSchemaField` in resource schema
// - provider-level region
// - region extracted from the provider-level zone
func getRegionFromFrameworkSchema(regionSchemaField string, rRegion, rZone, pRegion, pZone types.String, diags *diag.Diagnostics) types.String {

	// Value search #1: use region value from resource/data source config
	if !rRegion.IsNull() && rRegion.ValueString() != "" {
		return rRegion
	}

	// Value search #2: extract rergion from zone value from resource/data source config
	if !rZone.IsNull() && rZone.ValueString() != "" {
		v := getRegionFromZone(rZone.ValueString())
		if v != "" {
			return types.StringValue(v)
		}
	}

	// Value search #3: use region from provider config
	if !pRegion.IsNull() && pRegion.ValueString() != "" {
		return pRegion
	}

	// Value search #4: extract region from zone value in provider config
	if !pZone.IsNull() && pZone.ValueString() != "" {
		v := getRegionFromZone(pZone.ValueString())
		if v != "" {
			return types.StringValue(v)
		}
	}

	diags.AddError("required field is not set", fmt.Sprintf("%s is not set", regionSchemaField))
	return types.StringNull()
}

func getZoneFramework(rVal, pVal types.String, diags *diag.Diagnostics) types.String {
	if !rVal.IsNull() && rVal.ValueString() != "" {
		return rVal
	}

	if !pVal.IsNull() && pVal.ValueString() != "" {
		return pVal
	}

	diags.AddError("required field is not set", "cannot determine zone: set in this resource, or set provider-level zone.")
	return types.StringNull()
}

// getLocationFramework (mimicking getLocation) attempts to get values in this order (if they exist):
// - location argument in the resource config
// - region argument in the resource config
// - zone argument in the resource config
// - zone argument set in the provider config
func getLocationFramework(rLocation, rRegion, rZone, pZone types.String, diags *diag.Diagnostics) types.String {
	if !rLocation.IsNull() && rLocation.ValueString() != "" {
		return rLocation
	}

	if !rRegion.IsNull() && rRegion.ValueString() != "" {
		return rRegion
	}

	if !rZone.IsNull() && rZone.ValueString() != "" {
		return rZone
	}

	if !pZone.IsNull() && pZone.ValueString() != "" {
		return pZone
	}

	diags.AddError("required field is not set", "cannot determine location: set location or region in this resource, or set provider-level zone.")
	return types.StringNull()
}

func handleDatasourceNotFoundError(ctx context.Context, err error, state *tfsdk.State, resource string, diags *diag.Diagnostics) {
	if IsGoogleApiErrorWithCode(err, 404) {
		tflog.Warn(ctx, fmt.Sprintf("Removing %s because it's gone", resource))
		// The resource doesn't exist anymore
		state.RemoveResource(ctx)
	}

	diags.AddError(fmt.Sprintf("Error when reading or editing %s", resource), err.Error())
}

// field helpers

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
func parseProjectFieldValueFramework(resourceType, fieldValue, projectSchemaField string, rVal, pVal types.String, isEmptyValid bool, diags *diag.Diagnostics) *ProjectFieldValue {
	if len(fieldValue) == 0 {
		if isEmptyValid {
			return &ProjectFieldValue{resourceType: resourceType}
		}
		diags.AddError("field can not be empty", fmt.Sprintf("The project field for resource %s cannot be empty", resourceType))
		return nil
	}

	r := regexp.MustCompile(fmt.Sprintf(projectBasePattern, resourceType))
	if parts := r.FindStringSubmatch(fieldValue); parts != nil {
		return &ProjectFieldValue{
			Project: parts[1],
			Name:    parts[2],

			resourceType: resourceType,
		}
	}

	project := getProjectFromFrameworkSchema(projectSchemaField, rVal, pVal, diags)
	if diags.HasError() {
		return nil
	}

	return &ProjectFieldValue{
		Project: project.ValueString(),
		Name:    GetResourceNameFromSelfLink(fieldValue),

		resourceType: resourceType,
	}
}
