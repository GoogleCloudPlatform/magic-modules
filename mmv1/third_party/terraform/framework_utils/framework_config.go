package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	googleoauth "golang.org/x/oauth2/google"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Provider methods

// HandleDefaults will handle all the defaults necessary in the provider
func (p *frameworkProvider) HandleDefaults(ctx context.Context, data *ProviderModel, diags *diag.Diagnostics) {
	if data.AccessToken.IsNull() && data.Credentials.IsNull() {
		credentials := MultiEnvDefault([]string{
			"GOOGLE_CREDENTIALS",
			"GOOGLE_CLOUD_KEYFILE_JSON",
			"GCLOUD_KEYFILE_JSON",
		}, nil)

		data.Credentials = types.StringValue(credentials.(string))

		accessToken := MultiEnvDefault([]string{
			"GOOGLE_OAUTH_ACCESS_TOKEN",
		}, nil)

		data.AccessToken = types.StringValue(accessToken.(string))
	}

	if data.ImpersonateServiceAccount.IsNull() {
		data.ImpersonateServiceAccount = types.StringValue(os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT"))
	}

	if data.Project.IsNull() {
		project := MultiEnvDefault([]string{
			"GOOGLE_PROJECT",
			"GOOGLE_CLOUD_PROJECT",
			"GCLOUD_PROJECT",
			"CLOUDSDK_CORE_PROJECT",
		}, nil)
		if project != nil {
			data.Project = types.StringValue(project.(string))
		}
	}

	if data.BillingProject.IsNull() {
		data.BillingProject = types.StringValue(os.Getenv("GOOGLE_BILLING_PROJECT"))
	}

	if data.Region.IsNull() {
		region := MultiEnvDefault([]string{
			"GOOGLE_REGION",
			"GCLOUD_REGION",
			"CLOUDSDK_COMPUTE_REGION",
		}, nil)

		if region != nil {
			data.Region = types.StringValue(region.(string))
		}
	}

	if data.Zone.IsNull() {
		zone := MultiEnvDefault([]string{
			"GOOGLE_ZONE",
			"GCLOUD_ZONE",
			"CLOUDSDK_COMPUTE_ZONE",
		}, nil)

		if zone != nil {
			data.Zone = types.StringValue(zone.(string))
		}
	}

	if len(data.Scopes.Elements()) == 0 {
		var d diag.Diagnostics
		data.Scopes, d = types.ListValueFrom(ctx, types.StringType, defaultClientScopes)
		diags.Append(d...)
		if diags.HasError() {
			return
		}
	}

	var pbConfig ProviderBatching
	d := data.Batching.As(ctx, &pbConfig, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if diags.HasError() {
		return bc
	}

	if pbConfig.SendAfter.IsNull() {
		pbConfig.SendAfter = types.StringValue("10s")
	}

	if pbConfig.EnableBatching.IsNull() {
		pbConfig.EnableBatching = types.BoolValue(true)
	}

	data.Batching.

	if data.UserProjectOverride.IsNull() {
		override, err := strconv.ParseBool(os.Getenv("USER_PROJECT_OVERRIDE"))
		if err != nil {
			diags.AddError(
				"error parsing environment variable `USER_PROJECT_OVERRIDE` into bool", err.Error())
		}
		data.UserProjectOverride = types.BoolValue(override)
	}

	if data.RequestReason.IsNull() {
		data.RequestReason = types.StringValue(os.Getenv("CLOUDSDK_CORE_REQUEST_REASON"))
	}

	if data.RequestTimeout.IsNull() {
		data.RequestTimeout = types.StringValue("120")
	}
}

func (p *frameworkProvider) SetupClient(ctx context.Context, data ProviderModel, diags *diag.Diagnostics) {
	tokenSource := GetTokenSource(ctx, data, false, diags)
	if diags.HasError() {
		return
	}

	cleanCtx := context.WithValue(ctx, oauth2.HTTPClient, cleanhttp.DefaultClient())

	// 1. MTLS TRANSPORT/CLIENT - sets up proper auth headers
	client, _, err := transport.NewHTTPClient(cleanCtx, option.WithTokenSource(tokenSource))
	if err != nil {
		diags.AddError("error creating new http client", err.Error())
		return
	}

	// Userinfo is fetched before request logging is enabled to reduce additional noise.
	p.logGoogleIdentities(ctx, data, diags)
	if diags.HasError() {
		return
	}

	// 2. Logging Transport - ensure we log HTTP requests to GCP APIs.
	loggingTransport := logging.NewTransport("Google", client.Transport)

	// 3. Retry Transport - retries common temporary errors
	// Keep order for wrapping logging so we log each retried request as well.
	// This value should be used if needed to create shallow copies with additional retry predicates.
	// See ClientWithAdditionalRetries
	retryTransport := NewTransportWithDefaultRetries(loggingTransport)

	// 4. Header Transport - outer wrapper to inject additional headers we want to apply
	// before making requests
	headerTransport := newTransportWithHeaders(retryTransport)
	if !data.RequestReason.IsNull() {
		headerTransport.Set("X-Goog-Request-Reason", data.RequestReason.String())
	}

	// Ensure $userProject is set for all HTTP requests using the client if specified by the provider config
	// See https://cloud.google.com/apis/docs/system-parameters
	if data.UserProjectOverride.ValueBool() && !data.BillingProject.IsNull() {
		headerTransport.Set("X-Goog-User-Project", data.BillingProject.String())
	}

	// Set final transport value.
	client.Transport = headerTransport

	// This timeout is a timeout per HTTP request, not per logical operation.
	timeout, err := time.ParseDuration(data.RequestTimeout.String())
	if err != nil {
		diags.AddError("error parsing request timeout", err.Error())
	}
	client.Timeout = timeout

	p.client = client
}

func (p *frameworkProvider) logGoogleIdentities(ctx context.Context, data ProviderModel, diags *diag.Diagnostics) {
	if data.ImpersonateServiceAccount.IsNull() {

		tokenSource := GetTokenSource(ctx, data, true, diags)
		if diags.HasError() {
			return
		}

		p.client = oauth2.NewClient(ctx, tokenSource) // p.client isn't initialised fully when this code is called.

		email := GetCurrUserEmail(p, p.userAgent, diags)
		if diags.HasError() {
			log.Printf("[INFO] error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
			return
		}

		log.Printf("[INFO] Terraform is using this identity: %s", email)
		return
	}

	// Drop Impersonated ClientOption from OAuth2 TokenSource to infer original identity
	tokenSource := GetTokenSource(ctx, data, true, diags)
	if diags.HasError() {
		return
	}

	p.client = oauth2.NewClient(ctx, tokenSource) // p.client isn't initialised fully when this code is called.
	email := GetCurrUserEmail(p, p.userAgent, diags)
	if diags.HasError() {
		log.Printf("[INFO] error retrieving userinfo for your provider credentials. have you enabled the 'https://www.googleapis.com/auth/userinfo.email' scope?")
		return
	}

	log.Printf("[INFO] Terraform is configured with service account impersonation, original identity: %s, impersonated identity: %s", email, data.ImpersonateServiceAccount.String())

	// Add the Impersonated ClientOption back in to the OAuth2 TokenSource
	tokenSource = GetTokenSource(ctx, data, false, diags)
	if diags.HasError() {
		return
	}

	p.client = oauth2.NewClient(ctx, tokenSource) // p.client isn't initialised fully when this code is called.

	return
}

// Configuration helpers

// GetTokenSource gets token source based on the Google Credentials configured.
// If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds.
func GetTokenSource(ctx context.Context, data ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) oauth2.TokenSource {
	creds := GetCredentials(ctx, data, initialCredentialsOnly, diags)

	return creds.TokenSource
}

// GetCredentials gets credentials with a given scope (clientScopes).
// If initialCredentialsOnly is true, don't follow the impersonation
// settings and return the initial set of creds instead.
func GetCredentials(ctx context.Context, data ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) googleoauth.Credentials {
	var clientScopes []string
	var delegates []string

	d := data.Scopes.ElementsAs(ctx, &clientScopes, false)
	diags.Append(d...)
	if diags.HasError() {
		return googleoauth.Credentials{}
	}

	d = data.ImpersonateServiceAccountDelegates.ElementsAs(ctx, &delegates, false)
	diags.Append(d...)
	if diags.HasError() {
		return googleoauth.Credentials{}
	}

	if !data.AccessToken.IsNull() {
		contents, _, err := pathOrContents(data.AccessToken.String())
		if err != nil {
			diags.AddError("error loading access token", err.Error())
			return googleoauth.Credentials{}
		}

		token := &oauth2.Token{AccessToken: contents}
		if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
			opts := []option.ClientOption{option.WithTokenSource(oauth2.StaticTokenSource(token)), option.ImpersonateCredentials(data.ImpersonateServiceAccount.String(), delegates...), option.WithScopes(clientScopes...)}
			creds, err := transport.Creds(context.TODO(), opts...)
			if err != nil {
				diags.AddError("error impersonating credentials", err.Error())
				return googleoauth.Credentials{}
			}
			return *creds
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'access_token'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return googleoauth.Credentials{
			TokenSource: staticTokenSource{oauth2.StaticTokenSource(token)},
		}
	}

	if !data.Credentials.IsNull() {
		contents, _, err := pathOrContents(data.Credentials.String())
		if err != nil {
			diags.AddError(fmt.Sprintf("error loading credentials: %s", err), err.Error())
			return googleoauth.Credentials{}
		}

		if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
			opts := []option.ClientOption{option.WithCredentialsJSON([]byte(contents)), option.ImpersonateCredentials(data.ImpersonateServiceAccount.String(), delegates...), option.WithScopes(clientScopes...)}
			creds, err := transport.Creds(context.TODO(), opts...)
			if err != nil {
				diags.AddError("error impersonating credentials", err.Error())
				return googleoauth.Credentials{}
			}
			return *creds
		}

		creds, err := googleoauth.CredentialsFromJSON(ctx, []byte(contents), clientScopes...)
		if err != nil {
			diags.AddError(fmt.Sprintf("unable to parse credentials from '%s': %s", contents, err), err.Error())
			return googleoauth.Credentials{}
		}

		log.Printf("[INFO] Authenticating using configured Google JSON 'credentials'...")
		log.Printf("[INFO]   -- Scopes: %s", clientScopes)
		return *creds
	}

	if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
		opts := option.ImpersonateCredentials(data.ImpersonateServiceAccount.String(), delegates...)
		creds, err := transport.Creds(context.TODO(), opts, option.WithScopes(clientScopes...))
		if err != nil {
			diags.AddError("error impersonating credentials", err.Error())
			return googleoauth.Credentials{}
		}

		return *creds
	}

	log.Printf("[INFO] Authenticating using DefaultClient...")
	log.Printf("[INFO]   -- Scopes: %s", clientScopes)
	defaultTS, err := googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
	if err != nil {
		diags.AddError(fmt.Sprintf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  "+
			"No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'"), err.Error())
		return googleoauth.Credentials{}
	}

	return googleoauth.Credentials{
		TokenSource: defaultTS,
	}
}

// GetBatchingConfig returns the batching config object given the
// provider configuration set for batching
func GetBatchingConfig(ctx context.Context, data types.Object, diags *diag.Diagnostics) *batchingConfig {
	bc := &batchingConfig{
		sendAfter:      time.Second * defaultBatchSendIntervalSec,
		enableBatching: true,
	}

	if data.IsNull() {
		return bc
	}

	var pbConfig ProviderBatching
	d := data.As(ctx, &pbConfig, basetypes.ObjectAsOptions{})
	diags.Append(d...)
	if diags.HasError() {
		return bc
	}

	sendAfter, err := time.ParseDuration(pbConfig.SendAfter.String())
	if err != nil {
		diags.AddError("error parsing send after time duration", err.Error())
		return bc
	}

	bc.sendAfter = sendAfter

	if !pbConfig.EnableBatching.IsNull() {
		bc.enableBatching = pbConfig.EnableBatching.ValueBool()
	}

	return bc
}
