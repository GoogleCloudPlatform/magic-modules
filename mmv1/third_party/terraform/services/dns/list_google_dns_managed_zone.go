package dns

import (
	"context"
	"errors"
	"fmt"

	frameworkdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	googledns "google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func init() {
	registry.FrameworkListResource{
		Name:        "google_dns_managed_zone",
		ProductName: "dns",
		Func:        NewGoogleDnsManagedZoneListResource,
	}.Register()
}

type GoogleDnsManagedZoneResource struct {
	tpgresource.ListResourceMetadata
}

type GoogleDnsManagedZoneListModel struct {
	Project types.String `tfsdk:"project"`
}

func NewGoogleDnsManagedZoneListResource() list.ListResource {
	listR := &GoogleDnsManagedZoneResource{}
	listR.TypeName = "google_dns_managed_zone"
	listR.SDKv2Resource = ResourceDNSManagedZone()
	listR.ListConfigFields = []tpgresource.ListConfigField{{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true}}
	return listR
}

func (listR *GoogleDnsManagedZoneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	listR.ListResourceMetadata.Metadata(ctx, req, resp)

	resp.ResourceBehavior = resource.ResourceBehavior{
		MutableIdentity: true,
	}
}

func (listR *GoogleDnsManagedZoneResource) List(ctx context.Context, req list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleDnsManagedZoneListModel
	diags := req.Config.Get(ctx, &data)

	if diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	if listR.Client == nil {
		diags = append(diags, frameworkdiag.NewErrorDiagnostic(
			"provider not configured",
			"The Google provider client is not available; ensure the provider is configured.",
		))
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	tempData := listR.SDKv2Resource.Data(&terraform.InstanceState{})
	if !data.Project.IsNull() && !data.Project.IsUnknown() {
		if err := tempData.Set("project", data.Project.ValueString()); err != nil {
			diags.AddError("Failed to set project", err.Error())
			stream.Results = list.ListResultsStreamDiagnostics(diags)
			return
		}
	}

	project, err := tpgresource.GetProject(tempData, listR.Client)
	if err != nil {
		diags.AddError("Failed to resolve project", err.Error())
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}

	stream.Results = func(push func(list.ListResult) bool) {
		var streamDiags frameworkdiag.Diagnostics

		err := ListDnsManagedZones(listR.Client, project, func(rd *schema.ResourceData) error {
			result := req.NewListResult(ctx)

			if err := listR.SetResult(ctx, req.IncludeResource, &result, rd); err != nil {
				streamDiags.AddError("Failed to set result", err.Error())
				return err
			}

			if !push(result) {
				return errors.New("stream closed")
			}
			return nil
		})

		if err != nil {
			if err.Error() == "stream closed" {
				return
			}
			streamDiags.AddError("API Error listing DNS managed zones", fmt.Sprintf("Failed to list DNS managed zones in project %q: %v", project, err))
			result := req.NewListResult(ctx)
			result.Diagnostics = streamDiags
			push(result)
			return
		}

		if streamDiags.HasError() {
			stream.Results = list.ListResultsStreamDiagnostics(streamDiags)
		}
	}
}

func ListDnsManagedZones(config *transport_tpg.Config, project string, callback func(*schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}

	managedZoneSchema := ResourceDNSManagedZone()
	tempData := managedZoneSchema.Data(&terraform.InstanceState{})
	if project != "" {
		if err := tempData.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}

	userAgent, err := tpgresource.GenerateUserAgentString(tempData, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(tempData, config, "{{DNSBasePath}}projects/{{project}}/managedZones")
	if err != nil {
		return err
	}

	queryParams := make(map[string]string)
	for {
		reqURL, err := transport_tpg.AddQueryParams(url, queryParams)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    reqURL,
			UserAgent: userAgent,
		})
		if err != nil {
			return err
		}

		var listResp googledns.ManagedZonesListResponse
		if err := tpgresource.Convert(res, &listResp); err != nil {
			return fmt.Errorf("error parsing DNS managed zone list response: %w", err)
		}

		for _, managedZone := range listResp.ManagedZones {
			rd := managedZoneSchema.Data(&terraform.InstanceState{})
			if err := rd.Set("name", managedZone.Name); err != nil {
				return fmt.Errorf("error setting name on temporary resource data: %w", err)
			}
			if err := rd.Set("project", project); err != nil {
				return fmt.Errorf("error setting project on temporary resource data: %w", err)
			}
			rd.SetId(fmt.Sprintf("projects/%s/managedZones/%s", project, managedZone.Name))

			if err := managedZoneSchema.Read(rd, config); err != nil {
				return err
			}
			if err := rd.Set("name", managedZone.Name); err != nil {
				return fmt.Errorf("error setting name on temporary resource data: %w", err)
			}
			if err := rd.Set("project", project); err != nil {
				return fmt.Errorf("error setting project on temporary resource data: %w", err)
			}

			if err := tpgresource.SetResourceIdentityAttributes(rd, map[string]interface{}{
				"name":    managedZone.Name,
				"project": project,
			}); err != nil {
				return fmt.Errorf("error setting identity: %w", err)
			}

			if err := callback(rd); err != nil {
				return err
			}
		}

		if listResp.NextPageToken == "" {
			break
		}

		queryParams["pageToken"] = listResp.NextPageToken
	}

	return nil
}
