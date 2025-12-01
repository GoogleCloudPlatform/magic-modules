// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/ephemeralvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
	"google.golang.org/api/iam/v1"
)

var _ ephemeral.EphemeralResource = &googleEphemeralServiceAccountKey{}

func GoogleEphemeralServiceAccountKey() ephemeral.EphemeralResource {
	return &googleEphemeralServiceAccountKey{}
}

type googleEphemeralServiceAccountKey struct {
	providerConfig *transport_tpg.Config
}

func (p *googleEphemeralServiceAccountKey) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_account_key"
}

type ephemeralServiceAccountKeyModel struct {
	FetchKey         types.Bool   `tfsdk:"fetch_key"`
	ServiceAccountId types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	PublicKeyType    types.String `tfsdk:"public_key_type"`
	KeyAlgorithm     types.String `tfsdk:"key_algorithm"`
	PublicKeyData    types.String `tfsdk:"public_key_data"`
	PrivateKey       types.String `tfsdk:"private_key"`
	PrivateKeyType   types.String `tfsdk:"private_key_type"`
}

func (p *googleEphemeralServiceAccountKey) ConfigValidators(ctx context.Context) []ephemeral.ConfigValidator {
	return []ephemeral.ConfigValidator{
		ephemeralvalidator.Conflicting(
			path.MatchRoot("public_key_data"),
			path.MatchRoot("private_key_type"),
		),
		ephemeralvalidator.Conflicting(
			path.MatchRoot("public_key_data"),
			path.MatchRoot("key_algorithm"),
		),
		ephemeralvalidator.AtLeastOneOf(
			path.MatchRoot("service_account_id"),
			path.MatchRoot("name"),
		),
	}
}

func (p *googleEphemeralServiceAccountKey) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get an ephemeral service account public key.",
		Attributes: map[string]schema.Attribute{
			"service_account_id": schema.StringAttribute{
				Description: `The ID of the parent service account of the key. This can be a string in the format {ACCOUNT} or projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}, where {ACCOUNT} is the email address or unique id of the service account. If the {ACCOUNT} syntax is used, the project will be inferred from the provider's configuration.`,
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the service account key. This must have format `projects/{PROJECT_ID}/serviceAccounts/{ACCOUNT}/keys/{KEYID}`, where `{ACCOUNT}` is the email address or unique id of the service account.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(verify.ServiceAccountKeyNameRegex),
						"must match regex: "+verify.ServiceAccountKeyNameRegex,
					),
				}},
			"public_key_type": schema.StringAttribute{
				Description: "The output format of the public key requested. TYPE_X509_PEM_FILE is the default output format.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"TYPE_X509_PEM_FILE",
						"TYPE_RAW_PUBLIC_KEY",
					),
				},
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The algorithm used to generate the key.",
				Optional:    true,
				Computed:    true,
			},
			"public_key_data": schema.StringAttribute{
				Description: "The public key, base64 encoded.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The private key, base64 encoded.",
				Optional:    true,
				Computed:    true,
			},
			"private_key_type": schema.StringAttribute{
				Description: "The type of the private key.",
				Optional:    true,
				Computed:    true,
			},
			"fetch_key": schema.BoolAttribute{
				Description: "Whether to fetch the public key.",
				Optional:    true,
			},
		},
	}
}

func (p *googleEphemeralServiceAccountKey) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	pd, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	p.providerConfig = pd
}

type ServiceAccountKeyPrivateData struct {
	Name       string `json:"name"`
	FetchedKey bool   `json:"fetched_key"`
}

var createdServiceAccountKey bool

func (p *googleEphemeralServiceAccountKey) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data = ephemeralServiceAccountKeyModel{}
	var err error
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var serviceAccountKey *iam.ServiceAccountKey
	var saName string
	if data.ServiceAccountId.ValueString() != "" {
		saName = data.ServiceAccountId.ValueString()
	} else {
		saName = data.Name.ValueString()
	}

	var publicKeyType string
	if data.PublicKeyType.ValueString() == "" {
		publicKeyType = "TYPE_X509_PEM_FILE"
	} else {
		publicKeyType = data.PublicKeyType.ValueString()
	}

	saName, err = tpgresource.ServiceAccountFQN(saName, nil, p.providerConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting service account name",
			fmt.Sprintf("Error getting service account name: %s", err),
		)
		return
	}

	createdServiceAccountKey = false
	if !data.FetchKey.ValueBool() {
		if data.PublicKeyData.ValueString() != "" {
			ru := &iam.UploadServiceAccountKeyRequest{
				PublicKeyData: data.PublicKeyData.ValueString(),
			}
			serviceAccountKey, err = p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys.Upload(saName, ru).Do()
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating service account key [Upload]",
					fmt.Sprintf("%s: %s", saName, err),
				)
				return
			}
			createdServiceAccountKey = true
		} else {
			var keyAlgorithm, privateKeyType string
			if data.PrivateKeyType.ValueString() == "" {
				privateKeyType = "TYPE_GOOGLE_CREDENTIALS_FILE"
			} else {
				privateKeyType = data.PrivateKeyType.ValueString()
			}
			if data.KeyAlgorithm.ValueString() == "" {
				keyAlgorithm = "KEY_ALG_RSA_2048"
			} else {
				keyAlgorithm = data.KeyAlgorithm.ValueString()
			}
			rc := &iam.CreateServiceAccountKeyRequest{
				KeyAlgorithm:   keyAlgorithm,
				PrivateKeyType: privateKeyType,
			}
			serviceAccountKey, err = p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys.Create(saName, rc).Do()
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating service account key [Create]",
					fmt.Sprintf("%s: %s", saName, err),
				)
				return
			}
			createdServiceAccountKey = true
		}

		log.Printf("[DEBUG] Retrieving Service Account Key %q\n", serviceAccountKey.Name)
		err = ServiceAccountKeyWaitTime(p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys, serviceAccountKey.Name, publicKeyType, "Retrieving Service account key", 4*time.Minute)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving Service Account Key",
				fmt.Sprintf("Error retrieving Service Account Key %q: %s", serviceAccountKey.Name, err),
			)
			return
		}

		marshalledName, _ := json.Marshal(ServiceAccountKeyPrivateData{Name: serviceAccountKey.Name})

		resp.Private.SetKey(ctx, "name", marshalledName)

		data.Name = types.StringValue(serviceAccountKey.Name)
		data.KeyAlgorithm = types.StringValue(serviceAccountKey.KeyAlgorithm)
		data.PrivateKey = types.StringValue(serviceAccountKey.PrivateKeyData)
		data.PrivateKeyType = types.StringValue(serviceAccountKey.PrivateKeyType)
	}
	data.PublicKeyType = types.StringValue(publicKeyType)
	if serviceAccountKey != nil {
		log.Printf("[DEBUG] Retrieving Service Account Key %q\n", serviceAccountKey.Name)
		err = ServiceAccountKeyWaitTime(p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys, serviceAccountKey.Name, publicKeyType, "Retrieving Service account key", 4*time.Minute)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving Service Account Key",
				fmt.Sprintf("Error retrieving Service Account Key %q: %s", serviceAccountKey.Name, err),
			)
			return
		}
	} else {
		err = ServiceAccountKeyWaitTime(p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys, saName, publicKeyType, "Retrieving Service account key", 4*time.Minute)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving Service Account Key",
				fmt.Sprintf("Error retrieving Service Account Key %q: %s", saName, err),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

func (p *googleEphemeralServiceAccountKey) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	if !createdServiceAccountKey {
		return
	}
	dataBytes, diags := req.Private.GetKey(ctx, "name")
	if diags.HasError() {
		resp.Diagnostics.AddError(
			"Error getting name data",
			fmt.Sprintf("Error getting name data: %s", diags.Errors()),
		)
		return
	}
	var data ServiceAccountKeyPrivateData
	json.Unmarshal(dataBytes, &data)
	if data.Name != "" {
		log.Printf("[DEBUG] Deleting Service Account Key %q\n", data.Name)
		_, err := p.providerConfig.NewIamClient(p.providerConfig.UserAgent).Projects.ServiceAccounts.Keys.Delete(data.Name).Do()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting Service Account Key",
				fmt.Sprintf("Error deleting Service Account Key %q: %s", data.Name, err.Error()),
			)
			return
		}
	}
}
