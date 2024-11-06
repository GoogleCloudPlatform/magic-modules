package resourcemanager

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"google.golang.org/api/iamcredentials/v1"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountIdToken{}

func GoogleEphemeralServiceAccountIdToken() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountIdToken{}
}

type googleEphemeralServiceAccountIdToken struct {
	providerConfig *fwtransport.FrameworkProviderConfig
}

func (p *googleEphemeralServiceAccountIdToken) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_id_token"
}

type ephemeralServiceAccountIdTokenModel struct {
	TargetAudience       types.String `tfsdk:"target_audience"`
	TargetServiceAccount types.String `tfsdk:"target_service_account"`
	Delegates            types.Set    `tfsdk:"delegates"`
	IncludeEmail         types.Bool   `tfsdk:"include_email"`
	IdToken              types.String `tfsdk:"id_token"`
}

func (p *googleEphemeralServiceAccountIdToken) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"target_audience": schema.StringAttribute{
				Required: true,
			},
			"target_service_account": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					serviceAccountNameValidator{},
				},
			},
			"delegates": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					serviceAccountNameSetValidator{},
				},
			},
			"include_email": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"id_token": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountIdToken) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*fwtransport.FrameworkProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *fwtransport.FrameworkProviderConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// Required for accessing userAgent and passing as an argument into a util function
	p.providerConfig = pd
}

func (p *googleEphemeralServiceAccountIdToken) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data ephemeralServiceAccountIdTokenModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	targetAudience := data.TargetAudience.ValueString()
	creds := fwtransport.GetCredentials(ctx, fwmodels.ProviderModel{}, false, &resp.Diagnostics)

	targetServiceAccount := data.TargetServiceAccount.ValueString()
	// If a target service account is provided, use the API to generate the idToken
	if targetServiceAccount != "" {
		service := p.providerConfig.NewIamCredentialsClient(p.providerConfig.UserAgent)
		name := fmt.Sprintf("projects/-/serviceAccounts/%s", targetServiceAccount)
		DelegatesSetValue, _ := data.Delegates.ToSetValue(ctx)
		tokenRequest := &iamcredentials.GenerateIdTokenRequest{
			Audience:     targetAudience,
			IncludeEmail: data.IncludeEmail.ValueBool(),
			Delegates:    StringSet(DelegatesSetValue),
		}
		at, err := service.Projects.ServiceAccounts.GenerateIdToken(name, tokenRequest).Do()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error calling iamcredentials.GenerateIdToken",
				err.Error(),
			)
			return
		}

		data.IdToken = types.StringValue(at.Token)
		resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
		return
	}

	// If no target service account, use the default credentials
	ctx = context.Background()
	co := []option.ClientOption{}
	if creds.JSON != nil {
		co = append(co, idtoken.WithCredentialsJSON(creds.JSON))
	}

	idTokenSource, err := idtoken.NewTokenSource(ctx, targetAudience, co...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve TokenSource",
			err.Error(),
		)
		return
	}
	idToken, err := idTokenSource.Token()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to retrieve Token",
			err.Error(),
		)
		return
	}

	data.IdToken = types.StringValue(idToken.AccessToken)
	resp.Diagnostics.Append(resp.Result.Set(ctx, data)...)
}

func StringSet(d basetypes.SetValue) []string {

	StringSlice := make([]string, 0)
	for _, v := range d.Elements() {
		StringSlice = append(StringSlice, v.(basetypes.StringValue).ValueString())
	}
	return StringSlice
}

var serviceAccountNamePatterns = []string{
	`^.+@.+\.iam\.gserviceaccount\.com$`,                     // Standard IAM service account
	`^.+@developer\.gserviceaccount\.com$`,                   // Legacy developer service account
	`^.+@appspot\.gserviceaccount\.com$`,                     // App Engine service account
	`^.+@cloudservices\.gserviceaccount\.com$`,               // Google Cloud services service account
	`^.+@cloudbuild\.gserviceaccount\.com$`,                  // Cloud Build service account
	`^service-[0-9]+@.+-compute\.iam\.gserviceaccount\.com$`, // Compute Engine service account
}

// Create a custom validator for service account names
type serviceAccountNameValidator struct{}

func (v serviceAccountNameValidator) Description(ctx context.Context) string {
	return "value must be a valid service account email address"
}

func (v serviceAccountNameValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v serviceAccountNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
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

// Create a custom validator for sets of service account names
type serviceAccountNameSetValidator struct{}

func (v serviceAccountNameSetValidator) Description(ctx context.Context) string {
	return "all values must be valid service account email addresses"
}

func (v serviceAccountNameSetValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v serviceAccountNameSetValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	for _, element := range elements {
		stringValue, ok := element.(basetypes.StringValue)
		if !ok {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Element Type",
				"Set element must be a string",
			)
			return
		}

		value := stringValue.ValueString()
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
				fmt.Sprintf("Service account name %q must match one of the expected patterns for Google service accounts", value),
			)
		}
	}
}
