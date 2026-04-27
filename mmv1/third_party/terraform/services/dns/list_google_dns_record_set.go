// Copyright (c) IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type GoogleDnsRecordSetListResource struct {
	tpgresource.ListResourceMetadata
}

type GoogleDnsRecordSetListModel struct {
	Project     types.String `tfsdk:"project"`
	ManagedZone types.String `tfsdk:"managed_zone"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
}

func NewGoogleDnsRecordSetListResource() list.ListResource {
	listR := &GoogleDnsRecordSetListResource{}
	listR.TypeName = "google_dns_record_set"
	listR.SDKv2Resource = ResourceDnsRecordSet()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "managed_zone", Kind: tpgresource.ListConfigKindString, Optional: false},
		{Name: "name", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "type", Kind: tpgresource.ListConfigKindString, Optional: true},
	}
	return listR
}

func (listR *GoogleDnsRecordSetListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleDnsRecordSetListModel
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
	managedZone := data.ManagedZone.ValueString()
	name := data.Name.ValueString()
	recordType := data.Type.ValueString()

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListDnsRecordSets(listR.Client, project, managedZone, name, recordType, func(rd *schema.ResourceData) error {
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

func ListDnsRecordSets(config *transport_tpg.Config, project, managedZone, name, recordType string, callback func(rd *schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}
	d := ResourceDnsRecordSet().Data(&terraform.InstanceState{})
	if err := d.Set("managed_zone", managedZone); err != nil {
		return fmt.Errorf("error setting managed_zone on temporary resource data: %w", err)
	}
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	req := config.NewDnsClient(userAgent).ResourceRecordSets.List(project, managedZone)
	if name != "" {
		req.Name(name)
		if recordType != "" {
			req.Type(recordType)
		}
	}

	return req.Pages(context.Background(), func(page *dns.ResourceRecordSetsListResponse) error {
		for _, rrset := range page.Rrsets {
			if recordType != "" && name == "" && rrset.Type != recordType {
				continue
			}

			rd := ResourceDnsRecordSet().Data(&terraform.InstanceState{})
			rd.SetId(fmt.Sprintf("projects/%s/managedZones/%s/rrsets/%s/%s", project, managedZone, rrset.Name, rrset.Type))

			if err := populateDnsRecordSetResourceData(rd, rrset, project, managedZone); err != nil {
				return err
			}
			if err := callback(rd); err != nil {
				return err
			}
		}
		return nil
	})
}
