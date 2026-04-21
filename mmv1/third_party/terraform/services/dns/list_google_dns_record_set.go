package dns

import (
	"context"
	"errors"
	"fmt"

	frameworkdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	googledns "google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// var _ tpgresource.ListResourceWithRawV5Schemas = &GoogleDnsRecordSetResource{}

type GoogleDnsRecordSetResource struct {
	tpgresource.ListResourceMetadata
}

type GoogleDnsRecordSetListModel struct {
	Project     types.String `tfsdk:"project"`
	ManagedZone types.String `tfsdk:"managed_zone"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
}

func NewGoogleDnsRecordSetListResource() list.ListResource {
	listR := &GoogleDnsRecordSetResource{}
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

func (listR *GoogleDnsRecordSetResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleDnsRecordSetListModel
	diags := req.Config.Get(ctx, &data)

	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if listR.Client == nil {
		diags = append(diags, frameworkdiag.NewErrorDiagnostic(
			"provider not configured",
			"The Google provider client is not available; ensure the prtovider is configured.",
		))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	project := listR.GetProject(data.Project)
	managedZone := data.ManagedZone.ValueString()
	recordName := ""
	recordType := ""

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		recordName = data.Name.ValueString()
	}

	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		recordType = data.Type.ValueString()
	}

	stream.Results = func(push func(list.ListResult) bool) {
		err := ListDnsRecordSets(listR.Client, project, managedZone, recordName, recordType, func(rd *schema.ResourceData) error {
			result := req.NewListResult(ctx)

			if err := listR.SetResult(ctx, req.IncludeResource, &result, rd); err != nil {
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

func flattenGoogleDNSRecordSetListItem(
	rrset *googledns.ResourceRecordSet,
	d *schema.ResourceData,
	project string,
	managedZone string,
) error {
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %w", err)
	}

	if err := d.Set("managed_zone", managedZone); err != nil {
		return fmt.Errorf("error setting managed_zone: %w", err)
	}

	if err := d.Set("name", rrset.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}

	if err := d.Set("type", rrset.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}

	if err := d.Set("ttl", rrset.Ttl); err != nil {
		return fmt.Errorf("error setting ttl: %w", err)
	}

	if err := d.Set("rrdatas", rrset.Rrdatas); err != nil {
		return fmt.Errorf("error setting rrdatas: %w", err)
	}

	if err := d.Set("routing_policy", flattenDnsRecordSetRoutingPolicy(rrset.RoutingPolicy)); err != nil {
		return fmt.Errorf("error setting routing_policy: %w", err)
	}

	d.SetId(fmt.Sprintf(
		"project/%s/managedZones/%s/rrsets/%s/%s",
		project,
		managedZone,
		rrset.Name,
		rrset.Type,
	))

	return nil
}

func ListDnsRecordSets(
	config *transport_tpg.Config,
	project string,
	managedZone string,
	recordName string,
	recordType string,
	callback func(*schema.ResourceData) error,
) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}

	tempData := ResourceDnsRecordSet().Data(&terraform.InstanceState{})
	if project != "" {
		if err := tempData.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}
	if err := tempData.Set("managed_zone", managedZone); err != nil {
		return fmt.Errorf("error setting managed_zone on temporary resource data: %w", err)
	}

	userAgent, err := tpgresource.GenerateUserAgentString(tempData, config.UserAgent)
	if err != nil {
		return err
	}

	req := config.NewDnsClient(userAgent).ResourceRecordSets.List(project, managedZone)

	if recordName != "" {
		req = req.Name(recordName)
		if recordType != "" {
			req = req.Type(recordType)
		}
	}

	return req.Pages(context.Background(), func(resp *googledns.ResourceRecordSetsListResponse) error {
		for _, rrset := range resp.Rrsets {
			rd := ResourceDnsRecordSet().Data(&terraform.InstanceState{})

			if err := flattenGoogleDNSRecordSetListItem(rrset, rd, project, managedZone); err != nil {
				return err
			}

			if err := callback(rd); err != nil {
				return err
			}
		}
		return nil
	})
}
