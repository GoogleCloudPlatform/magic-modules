package fwvalidators

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	googleoauth "golang.org/x/oauth2/google"
)

// Credentials Validator
var _ validator.String = credentialsValidator{}

// credentialsValidator validates that a string Attribute's is valid JSON credentials.
type credentialsValidator struct {
}

// Description describes the validation in plain text formatting.
func (v credentialsValidator) Description(_ context.Context) string {
	return "value must be a path to valid JSON credentials or valid, raw, JSON credentials"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v credentialsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v credentialsValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// if this is a path and we can stat it, assume it's ok
	if _, err := os.Stat(value); err == nil {
		return
	}
	if _, err := googleoauth.CredentialsFromJSON(context.Background(), []byte(value)); err != nil {
		response.Diagnostics.AddError("JSON credentials are not valid", err.Error())
	}
}

func CredentialsValidator() validator.String {
	return credentialsValidator{}
}

// Non Negative Duration Validator
type nonnegativedurationValidator struct {
}

// Description describes the validation in plain text formatting.
func (v nonnegativedurationValidator) Description(_ context.Context) string {
	return "value expected to be a string representing a non-negative duration"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v nonnegativedurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v nonnegativedurationValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()
	dur, err := time.ParseDuration(value)
	if err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("expected %s to be a duration", value), err.Error())
		return
	}

	if dur < 0 {
		response.Diagnostics.AddError("duration must be non-negative", fmt.Sprintf("duration provided: %d", dur))
	}
}

func NonNegativeDurationValidator() validator.String {
	return nonnegativedurationValidator{}
}

// Non Empty String Validator
type nonEmptyStringValidator struct {
}

// Description describes the validation in plain text formatting.
func (v nonEmptyStringValidator) Description(_ context.Context) string {
	return "value expected to be a string that isn't an empty string"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (v nonEmptyStringValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation.
func (v nonEmptyStringValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	if value == "" {
		response.Diagnostics.AddError("expected a non-empty string", fmt.Sprintf("%s was set to `%s`", request.Path, value))
	}
}

func NonEmptyStringValidator() validator.String {
	return nonEmptyStringValidator{}
}

func StringSet(d basetypes.SetValue) []string {

	StringSlice := make([]string, 0)
	for _, v := range d.Elements() {
		StringSlice = append(StringSlice, v.(basetypes.StringValue).ValueString())
	}
	return StringSlice
}

// Define the possible service account name patterns
var serviceAccountNamePatterns = []string{
	`^.+@.+\.iam\.gserviceaccount\.com$`,                     // Standard IAM service account
	`^.+@developer\.gserviceaccount\.com$`,                   // Legacy developer service account
	`^.+@appspot\.gserviceaccount\.com$`,                     // App Engine service account
	`^.+@cloudservices\.gserviceaccount\.com$`,               // Google Cloud services service account
	`^.+@cloudbuild\.gserviceaccount\.com$`,                  // Cloud Build service account
	`^service-[0-9]+@.+-compute\.iam\.gserviceaccount\.com$`, // Compute Engine service account
}

// Create a custom validator for service account names
type ServiceAccountNameValidator struct{}

func (v ServiceAccountNameValidator) Description(ctx context.Context) string {
	return "value must be a valid service account email address"
}

func (v ServiceAccountNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v ServiceAccountNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	valid := false
	for _, pattern := range serviceAccountNamePatterns {
		if matched, _ := regexp.MatchString(pattern, value); matched {
			valid = true
			break
		}
	}

	if !valid {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Service Account Name",
			"Service account name must match one of the expected patterns for Google service accounts",
		)
	}
}

// Create a custom validator for duration
type DurationValidator struct {
	MaxDuration time.Duration
}

func (v DurationValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("value must be a valid duration string less than or equal to %v", v.MaxDuration)
}

func (v DurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v DurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()
	duration, err := time.ParseDuration(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Duration Format",
			"Duration must be a valid duration string (e.g., '3600s', '1h')",
		)
		return
	}

	if duration > v.MaxDuration {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Duration Too Long",
			fmt.Sprintf("Duration must be less than or equal to %v", v.MaxDuration),
		)
	}
}

// ServiceScopeValidator validates that a service scope is in canonical form
var _ validator.String = &ServiceScopeValidator{}

// ServiceScopeValidator validates service scope strings
type ServiceScopeValidator struct {
}

// Description returns a plain text description of the validator's behavior
func (v ServiceScopeValidator) Description(ctx context.Context) string {
	return "service scope must be in canonical form"
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior
func (v ServiceScopeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// ValidateString performs the validation
func (v ServiceScopeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	canonicalized := CanonicalizeServiceScope(req.ConfigValue.ValueString())
	if req.ConfigValue.ValueString() != canonicalized {
		resp.Diagnostics.AddAttributeWarning(
			req.Path,
			"Non-canonical service scope",
			fmt.Sprintf("Service scope %q will be canonicalized to %q",
				req.ConfigValue.ValueString(),
				canonicalized,
			),
		)
	}
}

func CanonicalizeServiceScope(scope string) string {
	// This is a convenience map of short names used by the gcloud tool
	// to the GCE auth endpoints they alias to.
	scopeMap := map[string]string{
		"bigquery":              "https://www.googleapis.com/auth/bigquery",
		"cloud-platform":        "https://www.googleapis.com/auth/cloud-platform",
		"cloud-source-repos":    "https://www.googleapis.com/auth/source.full_control",
		"cloud-source-repos-ro": "https://www.googleapis.com/auth/source.read_only",
		"compute-ro":            "https://www.googleapis.com/auth/compute.readonly",
		"compute-rw":            "https://www.googleapis.com/auth/compute",
		"datastore":             "https://www.googleapis.com/auth/datastore",
		"logging-write":         "https://www.googleapis.com/auth/logging.write",
		"monitoring":            "https://www.googleapis.com/auth/monitoring",
		"monitoring-read":       "https://www.googleapis.com/auth/monitoring.read",
		"monitoring-write":      "https://www.googleapis.com/auth/monitoring.write",
		"pubsub":                "https://www.googleapis.com/auth/pubsub",
		"service-control":       "https://www.googleapis.com/auth/servicecontrol",
		"service-management":    "https://www.googleapis.com/auth/service.management.readonly",
		"sql":                   "https://www.googleapis.com/auth/sqlservice",
		"sql-admin":             "https://www.googleapis.com/auth/sqlservice.admin",
		"storage-full":          "https://www.googleapis.com/auth/devstorage.full_control",
		"storage-ro":            "https://www.googleapis.com/auth/devstorage.read_only",
		"storage-rw":            "https://www.googleapis.com/auth/devstorage.read_write",
		"taskqueue":             "https://www.googleapis.com/auth/taskqueue",
		"trace":                 "https://www.googleapis.com/auth/trace.append",
		"useraccounts-ro":       "https://www.googleapis.com/auth/cloud.useraccounts.readonly",
		"useraccounts-rw":       "https://www.googleapis.com/auth/cloud.useraccounts",
		"userinfo-email":        "https://www.googleapis.com/auth/userinfo.email",
	}

	if matchedURL, ok := scopeMap[scope]; ok {
		return matchedURL
	}

	return scope
}
