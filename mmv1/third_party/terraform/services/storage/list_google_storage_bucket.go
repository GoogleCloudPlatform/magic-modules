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

// ListStorageBuckets lists buckets in a project (optional prefix), then loads each bucket with
// Buckets.Get and setStorageBucket — same read path as [dataSourceGoogleStorageBucketRead].
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

	params := map[string]string{
		"project": proj,
	}
	if prefix != "" {
		params["prefix"] = prefix
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	for {
		baseURL := "https://storage.googleapis.com/storage/v1/b"
		url, err := transport_tpg.AddQueryParams(baseURL, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{
				transport_tpg.Is429RetryableQuotaError,
			},
		})
		if err != nil {
			return err
		}

		rawItems, ok := res["items"].([]interface{})
		if ok {
			for _, raw := range rawItems {
				item, ok := raw.(map[string]interface{})
				if !ok {
					return fmt.Errorf("expected bucket item map, got %T", raw)
				}
				name, ok := item["name"].(string)
				if !ok || name == "" {
					return fmt.Errorf("bucket list item missing name")
				}

				if err := d.Set("name", name); err != nil {
					return fmt.Errorf("error setting name on temporary resource data: %w", err)
				}

				bucketRes, err := config.NewStorageClient(userAgent).Buckets.Get(name).Do()
				if err != nil {
					return err
				}

				if err := setStorageBucket(d, config, bucketRes, name, userAgent); err != nil {
					return err
				}

				if err := callback(d); err != nil {
					return err
				}
			}
		}

		nextTok, ok := res["nextPageToken"].(string)
		if !ok || nextTok == "" {
			break
		}
		params["pageToken"] = nextTok
	}

	return nil
}
