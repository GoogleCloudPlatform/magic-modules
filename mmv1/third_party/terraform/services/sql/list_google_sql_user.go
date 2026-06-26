// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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

			if err := listR.SetResult(ctx, listReq.IncludeResource, &result, rd, "name"); err != nil {
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

func flattenSqlUserListItem(user *sqladmin.User, d *schema.ResourceData, project string) error {
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
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %w", err)
	}

	iamEmail, databaseRoles, err := sqlUserOptionalFields(user)
	if err != nil {
		return err
	}

	if _, ok := ResourceSqlUser().Schema["iam_email"]; ok {
		if err := d.Set("iam_email", iamEmail); err != nil {
			return fmt.Errorf("error setting iam_email: %w", err)
		}
	}
	if err := d.Set("sql_server_user_details", flattenSqlServerUserDetails(user.SqlserverUserDetails)); err != nil {
		return fmt.Errorf("error setting sql_server_user_details: %w", err)
	}
	if user.PasswordPolicy != nil {
		passwordPolicy := flattenPasswordPolicy(user.PasswordPolicy)
		if len(passwordPolicy.([]map[string]interface{})[0]) != 0 {
			if err := d.Set("password_policy", passwordPolicy); err != nil {
				return fmt.Errorf("error setting password_policy: %w", err)
			}
		}
	}
	if _, ok := ResourceSqlUser().Schema["database_roles"]; ok {
		if err := d.Set("database_roles", databaseRoles); err != nil {
			return fmt.Errorf("error setting database_roles: %w", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))
	return tpgresource.SetResourceIdentityAttributes(d, map[string]interface{}{
		"project":  project,
		"instance": user.Instance,
		"name":     user.Name,
		"host":     user.Host,
	})
}

func sqlUserOptionalFields(user *sqladmin.User) (string, []string, error) {
	rawBytes, err := json.Marshal(user)
	if err != nil {
		return "", nil, fmt.Errorf("error marshalling sql user: %w", err)
	}

	raw := map[string]interface{}{}
	if err := json.Unmarshal(rawBytes, &raw); err != nil {
		return "", nil, fmt.Errorf("error unmarshalling sql user: %w", err)
	}

	iamEmail, _ := raw["iamEmail"].(string)
	databaseRoles := make([]string, 0)
	if roles, ok := raw["databaseRoles"].([]interface{}); ok {
		for _, role := range roles {
			if s, ok := role.(string); ok {
				databaseRoles = append(databaseRoles, s)
			}
		}
	}

	return iamEmail, databaseRoles, nil
}

func ListSqlUsers(config *transport_tpg.Config, project, instance string, callback func(*schema.ResourceData) error) error {
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

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	var users *sqladmin.UsersListResponse
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			var rerr error
			users, rerr = NewClient(config, userAgent).Users.List(project, instance).Do()
			return rerr
		},
		Timeout: 5 * time.Minute,
	})
	if err != nil {
		return fmt.Errorf("error listing SQL users for instance %q in project %q: %w", instance, project, err)
	}

	for _, user := range users.Items {
		if user == nil {
			continue
		}
		if err := flattenSqlUserListItem(user, d, project); err != nil {
			return err
		}
		if err := callback(d); err != nil {
			return err
		}
	}

	return nil
}
