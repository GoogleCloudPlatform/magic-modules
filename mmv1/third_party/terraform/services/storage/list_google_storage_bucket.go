// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

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
// Fields map to Cloud Storage JSON API buckets.list query parameters where applicable; see
// https://cloud.google.com/storage/docs/json_api/v1/buckets/list
type GoogleStorageBucketListModel struct {
	Project types.String `tfsdk:"project"`
	Prefix  types.String `tfsdk:"prefix"`

	MaxResults           types.Int64  `tfsdk:"max_results"`
	Projection           types.String `tfsdk:"projection"`
	ReturnPartialSuccess types.Bool   `tfsdk:"return_partial_success"`
	SoftDeleted          types.Bool   `tfsdk:"soft_deleted"`
}

func NewGoogleStorageBucketListResource() list.ListResource {
	listR := &GoogleStorageBucketListResource{}
	listR.TypeName = "google_storage_bucket"
	listR.SDKv2Resource = ResourceStorageBucket()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "prefix", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "max_results", Kind: tpgresource.ListConfigKindInt64, Optional: true},
		{Name: "projection", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "return_partial_success", Kind: tpgresource.ListConfigKindBool, Optional: true},
		{Name: "soft_deleted", Kind: tpgresource.ListConfigKindBool, Optional: true},
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

	var listOpts storageBucketListOpts
	listOpts.prefix = prefix
	if !data.MaxResults.IsNull() && !data.MaxResults.IsUnknown() {
		n := data.MaxResults.ValueInt64()
		if n < 0 {
			diags.AddError(
				"Invalid list configuration",
				"max_results must be a non-negative integer.",
			)
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}
		if n > 0 {
			listOpts.maxResults = n
		}
	}
	if !data.Projection.IsNull() && !data.Projection.IsUnknown() {
		p := strings.TrimSpace(data.Projection.ValueString())
		if p != "" && p != "full" && p != "noAcl" {
			diags.AddError(
				"Invalid list configuration",
				`projection must be "full", "noAcl", or unset (API default is noAcl).`,
			)
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}
		listOpts.projection = p
	}
	if !data.ReturnPartialSuccess.IsNull() && !data.ReturnPartialSuccess.IsUnknown() {
		listOpts.returnPartialSuccess = data.ReturnPartialSuccess.ValueBool()
	}
	if !data.SoftDeleted.IsNull() && !data.SoftDeleted.IsUnknown() {
		listOpts.softDeleted = data.SoftDeleted.ValueBool()
	}

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListStorageBuckets(listR.Client, project, listOpts, func(rd *schema.ResourceData) error {
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

	return setStorageBucket(d, config, &b, b.Name, userAgent)
}

// storageBucketListOpts carries optional buckets.list query parameters (zero value = omit).
type storageBucketListOpts struct {
	prefix                 string
	maxResults             int64
	projection             string
	returnPartialSuccess   bool
	softDeleted            bool
}

func ListStorageBuckets(config *transport_tpg.Config, project string, opts storageBucketListOpts, callback func(*schema.ResourceData) error) error {
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
	if opts.prefix != "" {
		listParams["prefix"] = opts.prefix
	}
	if opts.maxResults > 0 {
		listParams["maxResults"] = strconv.FormatInt(opts.maxResults, 10)
	}
	if opts.projection != "" {
		listParams["projection"] = opts.projection
	}
	if opts.returnPartialSuccess {
		listParams["returnPartialSuccess"] = "true"
	}
	if opts.softDeleted {
		listParams["softDeleted"] = "true"
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
