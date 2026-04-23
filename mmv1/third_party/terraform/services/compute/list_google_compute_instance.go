// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package compute

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

type GoogleComputeInstanceListResource struct {
	tpgresource.ListResourceMetadata
}

type GoogleComputeInstanceListModel struct {
	Project types.String `tfsdk:"project"`
	Zone    types.String `tfsdk:"zone"`
}

func NewGoogleComputeInstanceListResource() list.ListResource {
	listR := &GoogleComputeInstanceListResource{}
	listR.TypeName = "google_compute_instance"
	listR.SDKv2Resource = ResourceComputeInstance()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "zone", Kind: tpgresource.ListConfigKindString, Optional: true},
	}
	return listR
}

func (listR *GoogleComputeInstanceListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleComputeInstanceListModel
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
	zone := listR.GetZone(data.Zone)

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListComputeInstances(listR.Client, project, zone, func(rd *schema.ResourceData) error {
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

func flattenComputeInstanceListItem(res map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) error {
	name, _ := res["name"].(string)
	if name == "" {
		return fmt.Errorf("missing name in compute instance list response")
	}

	zoneUrl, _ := res["zone"].(string)
	zone := tpgresource.GetResourceNameFromSelfLink(zoneUrl)

	project := d.Get("project").(string)

	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("error setting zone: %w", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, name))
	return nil
}

func ListComputeInstances(config *transport_tpg.Config, project, zone string, callback func(rd *schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}
	d := ResourceComputeInstance().Data(&terraform.InstanceState{})
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("error setting zone on temporary resource data: %w", err)
	}
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}
	url, err := tpgresource.ReplaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/zones/{{zone}}/instances")
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
		ItemName:       "items",
		Flattener:      flattenComputeInstanceListItem,
		Callback:       callback,
	})
}
