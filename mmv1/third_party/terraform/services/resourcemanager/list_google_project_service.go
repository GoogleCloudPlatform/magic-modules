// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package resourcemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

type GoogleProjectServiceListResource struct {
	tpgresource.ListResourceMetadata
}

// GoogleProjectServiceListModel matches [ListResourceMetadata.ListConfigFields] (tfsdk names and types).
type GoogleProjectServiceListModel struct {
	Project types.String `tfsdk:"project"`
}

func NewGoogleProjectServiceListResource() list.ListResource {
	listR := &GoogleProjectServiceListResource{}
	listR.TypeName = "google_project_service"
	listR.SDKv2Resource = ResourceGoogleProjectService()
	listR.ListConfigFields = []tpgresource.ListConfigField{{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true}}
	return listR
}

func (listR *GoogleProjectServiceListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleProjectServiceListModel
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
	project := tpgresource.GetResourceNameFromSelfLink(listR.GetProject(data.Project))

	tempData := ResourceGoogleProjectService().Data(&terraform.InstanceState{})
	if err := tempData.Set("project", project); err != nil {
		diags.AddError("Config Error", fmt.Sprintf("Error setting project: %s", err))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	userAgent, err := tpgresource.GenerateUserAgentString(tempData, listR.Client.UserAgent)
	if err != nil {
		diags.AddError("Config Error", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	billingProject := project
	if bp, err := tpgresource.GetBillingProject(tempData, listR.Client); err == nil {
		billingProject = bp
	}

	// Service Usage's Services.List endpoint is eventually consistent across
	// replicas, and paginated calls can also span replicas, so a single read
	// may occasionally miss recently-enabled services. Call the endpoint
	// directly (bypassing the request batcher) multiple times, spaced apart in
	// time, and union the results so we're unlikely to miss a service that
	// only a subset of replicas know about yet.
	const (
		maxAttempts  = 12
		sleepBetween = 5 * time.Second
	)
	servicesList := make(map[string]struct{})
	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				diags.AddError("Context Error", ctx.Err().Error())
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			case <-time.After(sleepBetween):
			}
		}
		page, err := ListCurrentlyEnabledServices(project, billingProject, userAgent, listR.Client, time.Minute)
		if err != nil {
			diags.AddError("API Error", err.Error())
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}
		for s := range page {
			servicesList[s] = struct{}{}
		}
	}

	stream.Results = func(push func(list.ListResult) bool) {
		for serviceName := range servicesList {
			rd := ResourceGoogleProjectService().Data(&terraform.InstanceState{})
			if err := rd.Set("project", project); err != nil {
				diags.AddError("Config Error", fmt.Sprintf("Error setting project: %s", err))
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if err := rd.Set("service", serviceName); err != nil {
				diags.AddError("Config Error", fmt.Sprintf("Error setting service: %s", err))
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			rd.SetId(fmt.Sprintf("%s/%s", project, serviceName))

			result := listReq.NewListResult(ctx)
			if err := listR.SetResult(ctx, listReq.IncludeResource, &result, rd, "service"); err != nil {
				diags.AddError("Schema Error", err.Error())
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
			if !push(result) {
				stream.Results = list.ListResultsStreamDiagnostics(diags)
				return
			}
		}
		stream.Results = list.ListResultsStreamDiagnostics(diags)
	}
}
