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

// GoogleServiceAccountListModel matches [ListResourceMetadata.ListConfigFields] (tfsdk names and types).
type GoogleServiceAccountListModel struct {
	Project types.String `tfsdk:"project"`
}

func NewGoogleServiceAccountListResource() list.ListResource {
	r := &GoogleServiceAccountListResource{}
	r.TypeName = "google_service_account"
	r.ResourceSchema = ResourceGoogleServiceAccount()
	r.IdentityAttributes = []string{"email", "project"}
	r.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
	}
	return r
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

			if err := r.SetResult(ctx, req.IncludeResource, &result, rd); err != nil {
				return err
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
