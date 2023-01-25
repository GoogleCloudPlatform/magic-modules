package google

import(
	"context"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	googleoauth "golang.org/x/oauth2/google"
)

func CustomEndpointValidator() validator.String {
	return stringvalidator.RegexMatches(regexp.MustCompile(`.*/[^/]+/$`), ""),
}

var _ validator.String = credentialsValidator{}

// credentialsValidator validates that a string Attribute's is valid JSON credentials.
type credentialsValidator struct {
}

// Description describes the validation in plain text formatting.
func (validator credentialsValidator) Description(_ context.Context) string {
	return "value must be a path to valid JSON credentials or valid, raw, JSON credentials"
}

// MarkdownDescription describes the validation in Markdown formatting.
func (validator credentialsValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

// Validate performs the validation.
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
