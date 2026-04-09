// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/iam/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var _ tpgresource.ListResourceWithRawV5Schemas = &GoogleServiceAccountListResource{}

type GoogleServiceAccountListResource struct {
	tpgresource.ListResourceMetadata
}

func NewGoogleServiceAccountListResource() list.ListResource {
	return &GoogleServiceAccountListResource{}
}

func (r *GoogleServiceAccountListResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "google_service_account"
}

func (r *GoogleServiceAccountListResource) RawV5Schemas(ctx context.Context, _ list.RawV5SchemaRequest, resp *list.RawV5SchemaResponse) {
	sa := ResourceGoogleServiceAccount()
	resp.ProtoV5Schema = sa.ProtoSchema(ctx)()
	resp.ProtoV5IdentitySchema = sa.ProtoIdentitySchema(ctx)()
}

func (r *GoogleServiceAccountListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Defaults(req, resp)
}

func (r *GoogleServiceAccountListResource) ListResourceConfigSchema(ctx context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		Attributes: map[string]listschema.Attribute{
			"project": listschema.StringAttribute{
				Optional: true,
			},
		},
	}
}

type GoogleServiceAccountListModel struct {
	Project types.String `tfsdk:"project"`
}

func (r *GoogleServiceAccountListResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleServiceAccountListModel
	diags := req.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	if r.Client == nil {
		diags = append(diags, diag.NewErrorDiagnostic(
			"Provider not configured",
			"The Google provider client is not available; ensure the provider is configured (e.g. credentials and default project).",
		))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	project := r.GetProject(data.Project)

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListServiceAccounts(r.Client, project, func(rd *schema.ResourceData) error {
			result := req.NewListResult(ctx)

			if err := tpgresource.SetIdentityFields(ctx, result, rd, map[string]string{
				"email":   rd.Get("email").(string),
				"project": project,
			}); err != nil {
				return err
			}

			if req.IncludeResource {
				tfTypeResource, err := rd.TfTypeResourceState()
				if err != nil {
					return err
				}
				if err := result.Resource.Set(ctx, *tfTypeResource); err != nil {
					return errors.New("error setting resource")
				}
			}

			if !push(result) {
				return errors.New("stream closed")
			}
			return nil
		})
		if err != nil {
			diags.AddError("API Error", err.Error())
			result := req.NewListResult(ctx)
			result.Diagnostics = diags
			push(result)
		}
		stream.Results = list.ListResultsStreamDiagnostics(diags)
	}
}

func flattenGoogleServiceAccountListItem(res map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) error {
	var sa iam.ServiceAccount
	if err := tpgresource.Convert(res, &sa); err != nil {
		return err
	}
	d.SetId(sa.Name)
	return populateResourceData(d, &sa)
}

func ListServiceAccounts(config *transport_tpg.Config, project string, callback func(rd *schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}
	resourceData := ResourceGoogleServiceAccount().Data(&terraform.InstanceState{})
	if project != "" {
		if err := resourceData.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}
	url, err := tpgresource.ReplaceVars(resourceData, config, "{{IAMBasePath}}projects/{{project}}/serviceAccounts")
	if err != nil {
		return err
	}

	billingProject := ""
	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}
	if bp, err := tpgresource.GetBillingProject(resourceData, config); err == nil {
		billingProject = bp
	}

	userAgent, err := tpgresource.GenerateUserAgentString(resourceData, config.UserAgent)
	if err != nil {
		return err
	}

	return transport_tpg.ListPages(
		config,
		resourceData,
		url,
		billingProject,
		userAgent,
		"accounts",
		"",
		flattenGoogleServiceAccountListItem,
		callback,
	)
}
