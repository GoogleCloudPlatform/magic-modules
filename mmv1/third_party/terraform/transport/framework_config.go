package transport

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// // Configuration helpers

// // GetTokenSource gets token source based on the Google Credentials configured.
// // If initialCredentialsOnly is true, don't follow the impersonation settings and return the initial set of creds.
// func GetTokenSource(ctx context.Context, data ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) oauth2.TokenSource {
// 	creds := GetCredentials(ctx, data, initialCredentialsOnly, diags)

// 	return creds.TokenSource
// }

// // GetCredentials gets credentials with a given scope (clientScopes).
// // If initialCredentialsOnly is true, don't follow the impersonation
// // settings and return the initial set of creds instead.
// func GetCredentials(ctx context.Context, data ProviderModel, initialCredentialsOnly bool, diags *diag.Diagnostics) googleoauth.Credentials {
// 	var clientScopes []string
// 	var delegates []string

// 	d := data.Scopes.ElementsAs(ctx, &clientScopes, false)
// 	diags.Append(d...)
// 	if diags.HasError() {
// 		return googleoauth.Credentials{}
// 	}

// 	d = data.ImpersonateServiceAccountDelegates.ElementsAs(ctx, &delegates, false)
// 	diags.Append(d...)
// 	if diags.HasError() {
// 		return googleoauth.Credentials{}
// 	}

// 	if !data.AccessToken.IsNull() {
// 		contents, _, err := pathOrContents(data.AccessToken.ValueString())
// 		if err != nil {
// 			diags.AddError("error loading access token", err.Error())
// 			return googleoauth.Credentials{}
// 		}

// 		token := &oauth2.Token{AccessToken: contents}
// 		if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
// 			opts := []option.ClientOption{option.WithTokenSource(oauth2.StaticTokenSource(token)), option.ImpersonateCredentials(data.ImpersonateServiceAccount.ValueString(), delegates...), option.WithScopes(clientScopes...)}
// 			creds, err := transport.Creds(context.TODO(), opts...)
// 			if err != nil {
// 				diags.AddError("error impersonating credentials", err.Error())
// 				return googleoauth.Credentials{}
// 			}
// 			return *creds
// 		}

// 		tflog.Info(ctx, "Authenticating using configured Google JSON 'access_token'...")
// 		tflog.Info(ctx, fmt.Sprintf("  -- Scopes: %s", clientScopes))
// 		return googleoauth.Credentials{
// 			TokenSource: transport_tpg.StaticTokenSource{oauth2.StaticTokenSource(token)},
// 		}
// 	}

// 	if !data.Credentials.IsNull() {
// 		contents, _, err := pathOrContents(data.Credentials.ValueString())
// 		if err != nil {
// 			diags.AddError(fmt.Sprintf("error loading credentials: %s", err), err.Error())
// 			return googleoauth.Credentials{}
// 		}

// 		if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
// 			opts := []option.ClientOption{option.WithCredentialsJSON([]byte(contents)), option.ImpersonateCredentials(data.ImpersonateServiceAccount.ValueString(), delegates...), option.WithScopes(clientScopes...)}
// 			creds, err := transport.Creds(context.TODO(), opts...)
// 			if err != nil {
// 				diags.AddError("error impersonating credentials", err.Error())
// 				return googleoauth.Credentials{}
// 			}
// 			return *creds
// 		}

// 		creds, err := googleoauth.CredentialsFromJSON(ctx, []byte(contents), clientScopes...)
// 		if err != nil {
// 			diags.AddError("unable to parse credentials", err.Error())
// 			return googleoauth.Credentials{}
// 		}

// 		tflog.Info(ctx, "Authenticating using configured Google JSON 'credentials'...")
// 		tflog.Info(ctx, fmt.Sprintf("  -- Scopes: %s", clientScopes))
// 		return *creds
// 	}

// 	if !data.ImpersonateServiceAccount.IsNull() && !initialCredentialsOnly {
// 		opts := option.ImpersonateCredentials(data.ImpersonateServiceAccount.ValueString(), delegates...)
// 		creds, err := transport.Creds(context.TODO(), opts, option.WithScopes(clientScopes...))
// 		if err != nil {
// 			diags.AddError("error impersonating credentials", err.Error())
// 			return googleoauth.Credentials{}
// 		}

// 		return *creds
// 	}

// 	tflog.Info(ctx, "Authenticating using DefaultClient...")
// 	tflog.Info(ctx, fmt.Sprintf("  -- Scopes: %s", clientScopes))
// 	defaultTS, err := googleoauth.DefaultTokenSource(context.Background(), clientScopes...)
// 	if err != nil {
// 		diags.AddError(fmt.Sprintf("Attempted to load application default credentials since neither `credentials` nor `access_token` was set in the provider block.  "+
// 			"No credentials loaded. To use your gcloud credentials, run 'gcloud auth application-default login'"), err.Error())
// 		return googleoauth.Credentials{}
// 	}

// 	return googleoauth.Credentials{
// 		TokenSource: defaultTS,
// 	}
// }

// GetBatchingConfig returns the batching config object given the
// provider configuration set for batching
func GetBatchingConfig(ctx context.Context, data types.List, diags *diag.Diagnostics) *batchingConfig {
	bc := &batchingConfig{
		SendAfter:      time.Second * DefaultBatchSendIntervalSec,
		EnableBatching: true,
	}

	if data.IsNull() {
		return bc
	}

	var pbConfigs []ProviderBatching
	d := data.ElementsAs(ctx, &pbConfigs, true)
	diags.Append(d...)
	if diags.HasError() {
		return bc
	}

	sendAfter, err := time.ParseDuration(pbConfigs[0].SendAfter.ValueString())
	if err != nil {
		diags.AddError("error parsing send after time duration", err.Error())
		return bc
	}

	bc.SendAfter = sendAfter

	if !pbConfigs[0].EnableBatching.IsNull() {
		bc.EnableBatching = pbConfigs[0].EnableBatching.ValueBool()
	}

	return bc
}
