// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

func init() {
	registry.FrameworkListResource{
		Name:        "google_sql_user",
		ProductName: "sql",
		Func:        NewGoogleSqlUserListResource,
	}.Register()
}

var _ list.ListResource = &GoogleSqlUserListResource{}

type GoogleSqlUserListResource struct {
	tpgresource.ListResourceMetadata
}

func NewGoogleSqlUserListResource() list.ListResource {
	listR := &GoogleSqlUserListResource{}
	listR.TypeName = "google_sql_user"
	listR.SDKv2Resource = ResourceSqlUser()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
		{Name: "instance", Kind: tpgresource.ListConfigKindString, Optional: false},
	}
	return listR
}

// GoogleSqlUserListModel matches ListResourceMetadata.ListConfigFields.
type GoogleSqlUserListModel struct {
	Project  types.String `tfsdk:"project"`
	Instance types.String `tfsdk:"instance"`
}

func (listR *GoogleSqlUserListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleSqlUserListModel
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
	instance := data.Instance.ValueString()

	errStreamClosed := errors.New("stream closed")
	stream.Results = func(push func(list.ListResult) bool) {
		err := ListSqlUsers(listR.Client, project, instance, func(rd *schema.ResourceData) error {
			result := listReq.NewListResult(ctx)

			if err := listR.SetResult(ctx, listReq.IncludeResource, &result, rd, "name", "instance", "project", "host"); err != nil {
				return err
			}

			if !push(result) {
				return errStreamClosed
			}
			return nil
		})
		if err == nil || errors.Is(err, errStreamClosed) {
			return
		}
		diags.AddError("API Error", err.Error())
		result := listReq.NewListResult(ctx)
		result.Diagnostics = diags
		push(result)
	}
}

func flattenSqlUserListItem(res map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config, project string) error {
	var user sqladmin.User
	if err := tpgresource.Convert(res, &user); err != nil {
		return fmt.Errorf("error converting SQL user list response: %w", err)
	}

	if err := d.Set("host", user.Host); err != nil {
		return fmt.Errorf("error setting host: %w", err)
	}
	if err := d.Set("instance", user.Instance); err != nil {
		return fmt.Errorf("error setting instance: %w", err)
	}
	if err := d.Set("name", user.Name); err != nil {
		return fmt.Errorf("error setting name: %w", err)
	}
	if err := d.Set("type", user.Type); err != nil {
		return fmt.Errorf("error setting type: %w", err)
	}
	if err := d.Set("iam_email", user.IamEmail); err != nil {
		return fmt.Errorf("error setting iam_email: %w", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %w", err)
	}
	if err := d.Set("sql_server_user_details", flattenSqlServerUserDetails(user.SqlserverUserDetails)); err != nil {
		return fmt.Errorf("error setting sql server user details: %w", err)
	}
	if user.PasswordPolicy != nil {
		passwordPolicy := flattenPasswordPolicy(user.PasswordPolicy)
		if len(passwordPolicy.([]map[string]interface{})[0]) != 0 {
			if err := d.Set("password_policy", passwordPolicy); err != nil {
				return fmt.Errorf("error setting password_policy: %w", err)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))

	return nil
}

func ListSqlUsers(config *transport_tpg.Config, project, instance string, callback func(rd *schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}

	d := ResourceSqlUser().Data(&terraform.InstanceState{})
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}
	if instance != "" {
		if err := d.Set("instance", instance); err != nil {
			return fmt.Errorf("error setting instance on temporary resource data: %w", err)
		}
	}

	url, err := tpgresource.ReplaceVars(d, config, transport_tpg.BaseUrl(Product, config)+"projects/{{project}}/instances/{{instance}}/users")
	if err != nil {
		return err
	}

	billingProject := project
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
		Resource:       ResourceSqlUser(),
		ListURL:        url,
		BillingProject: billingProject,
		UserAgent:      userAgent,
		ItemName:       "items",
		Flattener: func(item map[string]interface{}, itemD *schema.ResourceData, c *transport_tpg.Config) error {
			return flattenSqlUserListItem(item, itemD, c, project)
		},
		Callback: callback,
	})
}
