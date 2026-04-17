// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	storagapi "google.golang.org/api/storage/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type GoogleStorageBucketListResource struct {
	tpgresource.ListResourceMetadata
}

// GoogleStorageBucketListModel matches [ListResourceMetadata.ListConfigFields] (tfsdk names and types).
// Project and prefix align with [DataSourceGoogleStorageBuckets] list filtering.
type GoogleStorageBucketListModel struct {
	Project types.String `tfsdk:"project"`
	Prefix  types.String `tfsdk:"prefix"`
}

func NewGoogleStorageBucketListResource() list.ListResource {
	listR := &GoogleStorageBucketListResource{}
	listR.TypeName = "google_storage_bucket"
	listR.SDKv2Resource = ResourceStorageBucket()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "prefix", Kind: tpgresource.ListConfigKindString, Optional: true},
	}
	return listR
}

func (listR *GoogleStorageBucketListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleStorageBucketListModel
	diags := listReq.Config.Get(ctx, &data)
	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	if listR.Client == nil {
		diags = append(diags, diag.NewErrorDiagnostic(
			"Provider not configured",
			"The Google provider client is not available; ensure the provider is configured (e.g. credentials and default project).",
		))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	project := listR.GetProject(data.Project)
	prefix := ""
	if !data.Prefix.IsNull() && !data.Prefix.IsUnknown() {
		prefix = data.Prefix.ValueString()
	}

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListStorageBuckets(listR.Client, project, prefix, func(rd *schema.ResourceData) error {
			result := listReq.NewListResult(ctx)

			if err := listR.SetResult(ctx, listReq.IncludeResource, &result, rd, "name"); err != nil {
				return err
			}

			if !push(result) {
				return errors.New("stream closed")
			}
			return nil
		})
		if err != nil {
			diags.AddError("API Error", err.Error())
			result := listReq.NewListResult(ctx)
			result.Diagnostics = diags
			push(result)
		}
	}
}

func flattenStorageBucketListItem(item map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) error {
	var b storagapi.Bucket
	if err := tpgresource.Convert(item, &b); err != nil {
		return fmt.Errorf("converting bucket list item: %w", err)
	}
	if b.Name == "" {
		return fmt.Errorf("bucket list item missing name")
	}
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	d.Set("name", b.Name)

	return setStorageBucket(d, config, &b, b.Name, userAgent)
}

func ListStorageBuckets(config *transport_tpg.Config, project string, prefix string, callback func(*schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}

	d := ResourceStorageBucket().Data(&terraform.InstanceState{})
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}

	proj, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("error resolving project: %w", err)
	}
	if err := d.Set("project", proj); err != nil {
		return fmt.Errorf("error setting project on temporary resource data: %w", err)
	}

	listParams := map[string]string{
		"project": proj,
	}
	if prefix != "" {
		listParams["prefix"] = prefix
	}

	listURL, err := transport_tpg.AddQueryParams("https://storage.googleapis.com/storage/v1/b", listParams)
	if err != nil {
		return err
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil && bp != "" {
		billingProject = bp
	}

	return transport_tpg.ListPages(transport_tpg.ListPagesOptions{
		Config:         config,
		TempData:       d,
		ListURL:        listURL,
		BillingProject: billingProject,
		UserAgent:      userAgent,
		ItemName:       "items",
		Flattener:      flattenStorageBucketListItem,
		Callback:       callback,
	})
}
