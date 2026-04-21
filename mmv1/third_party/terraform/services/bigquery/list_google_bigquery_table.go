// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package bigquery

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

type GoogleBigQueryTableListResource struct {
	tpgresource.ListResourceMetadata
}

type GoogleBigQueryTableListModel struct {
	DatasetID types.String `tfsdk:"dataset_id"`
	Project   types.String `tfsdk:"project"`
}

func NewGoogleBigQueryTableListResource() list.ListResource {
	listR := &GoogleBigQueryTableListResource{}
	listR.TypeName = "google_bigquery_table"
	listR.SDKv2Resource = ResourceBigQueryTable()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "dataset_id", Kind: tpgresource.ListConfigKindString, Optional: false},
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
	}
	return listR
}

func (listR *GoogleBigQueryTableListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleBigQueryTableListModel
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
	datasetID := data.DatasetID.ValueString()

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListBigQueryTables(listR.Client, project, datasetID, func(rd *schema.ResourceData) error {
			result := listReq.NewListResult(ctx)

			if err := listR.SetResult(ctx, listReq.IncludeResource, &result, rd, "table_id"); err != nil {
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

func flattenGoogleBigQueryTableListItem(res map[string]interface{}, d *schema.ResourceData, _ *transport_tpg.Config) error {
	tableRef, ok := res["tableReference"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing tableReference in BigQuery tables list response")
	}

	tableID, ok := tableRef["tableId"].(string)
	if !ok || tableID == "" {
		return fmt.Errorf("missing tableReference.tableId in BigQuery tables list response")
	}

	project := d.Get("project").(string)
	datasetID := d.Get("dataset_id").(string)

	if v, ok := tableRef["projectId"].(string); ok && v != "" {
		project = v
	}
	if v, ok := tableRef["datasetId"].(string); ok && v != "" {
		datasetID = v
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %w", err)
	}
	if err := d.Set("dataset_id", datasetID); err != nil {
		return fmt.Errorf("error setting dataset_id: %w", err)
	}
	if err := d.Set("table_id", tableID); err != nil {
		return fmt.Errorf("error setting table_id: %w", err)
	}

	if labels, ok := res["labels"].(map[string]interface{}); ok {
		if err := d.Set("labels", labels); err != nil {
			return fmt.Errorf("error setting labels: %w", err)
		}
	}

	if tableType, ok := res["type"].(string); ok && tableType != "" {
		if err := d.Set("type", tableType); err != nil {
			return fmt.Errorf("error setting type: %w", err)
		}
	}

	d.SetId(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", project, datasetID, tableID))
	return nil
}

func ListBigQueryTables(config *transport_tpg.Config, project, datasetID string, callback func(rd *schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}
	d := ResourceBigQueryTable().Data(&terraform.InstanceState{})
	if err := d.Set("dataset_id", datasetID); err != nil {
		return fmt.Errorf("error setting dataset_id on temporary resource data: %w", err)
	}
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}
	url, err := tpgresource.ReplaceVars(d, config, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}/tables")
	if err != nil {
		return err
	}

	billingProject := ""
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	return transport_tpg.ListPages(transport_tpg.ListPagesOptions{
		Config:         config,
		TempData:       d,
		ListURL:        url,
		BillingProject: billingProject,
		UserAgent:      userAgent,
		ItemName:       "tables",
		Flattener:      flattenGoogleBigQueryTableListItem,
		Callback:       callback,
	})
}
