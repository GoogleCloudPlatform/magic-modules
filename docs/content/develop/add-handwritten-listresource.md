---
title: "Add a list resource"
summary: "List resources let users run terraform query against existing Google Cloud resources of a given managed type."
weight: 65
---

# Add a list resource

List resources plug into Terraform’s plugin-framework [list-resource API](https://developer.hashicorp.com/terraform/plugin/framework/list-resources) so users can run [`terraform query`](https://developer.hashicorp.com/terraform/cli/commands/query) and use **`.tfquery.hcl`** files against resources that already exist in Google Cloud. Each list resource is tied to exactly one **managed resource** type: results expose that type’s [resource identity](https://developer.hashicorp.com/terraform/language/resources/identities), and optionally full state when `include_resource` is set on the `list` block.

For end-user behavior, file layout, and Terraform version requirements, see the registry guide [Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).

## Helpers in [`tpgresource/list_resource.go`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/tpgresource/list_resource.go)

Most shared behavior for handwritten list resources lives in **`mmv1/third_party/terraform/tpgresource/list_resource.go`**. Use it instead of reimplementing wiring:

| Symbol | Role |
|--------|------|
| **`ListResourceMetadata`** | Embed this struct. It implements plugin-framework hooks: **`Metadata`**, **`Configure`** (stores `*transport_tpg.Config` and default project/region/zone from the provider), **`ListResourceConfigSchema`** (builds the `list` block schema from **`ListConfigFields`** via **`NewListConfigSchema`**), and **`RawV5Schemas`** (exposes the managed resource’s SDK v2 schema and identity schema for query results). |
| **`ListConfigField`** / **`ListConfigKindString`** (etc.) | Declares each attribute allowed inside the list block’s `config { ... }` (for example optional `project`). Must match a separate **config model** struct with `tfsdk` tags for use with `listReq.Config.Get`. |
| **`GetProject`**, **`GetRegion`**, **`GetZone`**, **`GetLocation`** | Resolve a config field or fall back to the provider default (`types.String` from the model). |
| **`SetResult`** | After you populate a `*schema.ResourceData` for one listed item, call this to fill **`ListResult`**: identity state, optional full resource state when `listReq.IncludeResource` is true, and **`DisplayName`** using **`ListResultDisplayName`** and the `displayNameKeys` you pass (for example `"display_name", "email"`). |
| **`ListResultDisplayName`** | Standalone helper if you need a label from `ResourceData` without going through **`SetResult`**. |

Lower-level helpers in the same file include **`SetResourceIdentityAttributes`** (identity field writes on `ResourceData`). Identity copying for list rows uses the embedded managed resource’s identity schema via **`setResourceIdentity`** (see comments on **`SetResult`**).

## Prerequisites

1. A working managed resource implementation (SDK v2 `*schema.Resource`) for the type you want to list, including a correct **identity** schema for query results.
1. Familiarity with the Terraform plugin framework **`list`** package (`list.ListResource`, `list.ListRequest`, `list.ListResultsStream`) and Terraform **1.14+** for acceptance tests that exercise the list API.

## Add the list resource implementation

### 1. Create `list_google_<RESOURCE>.go`

Add the file under the appropriate service folder in [`magic-modules/mmv1/third_party/terraform/services`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/services). Pattern: embed **`tpgresource.ListResourceMetadata`**, define a config model whose `tfsdk` tags match **`ListConfigFields`**, implement **`List`**, and add a package-level **`List<ResourcePlural>`** helper plus a **flattener** so the actual REST enumeration goes through **`transport_tpg.ListPages`** (see the next subsection).

Skeleton for **`List`** (the method on your list resource type): it reads config, then delegates paging to **`List<…>`**, whose **callback** receives a fully populated **`ResourceData`** per row and calls **`SetResult`**. This mirrors [`list_google_service_account.go`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/services/resourcemanager/list_google_service_account.go):

```go
type GoogleExampleListResource struct {
	tpgresource.ListResourceMetadata
}

// Config model: field names and types must match ListConfigFields and listschema built by NewListConfigSchema.
type GoogleExampleListModel struct {
	Project types.String `tfsdk:"project"`
}

func NewGoogleExampleListResource() list.ListResource {
	listR := &GoogleExampleListResource{}
	listR.TypeName = "google_example" // must match managed resource type string
	listR.SDKv2Resource = ResourceGoogleExample()
	listR.ListConfigFields = []tpgresource.ListConfigField{
		{Name: "project", Kind: tpgresource.ListConfigKindString, Optional: true},
	}
	return listR
}

func (listR *GoogleExampleListResource) List(ctx context.Context, listReq list.ListRequest, stream *list.ListResultsStream) {
	var data GoogleExampleListModel
	if diags := listReq.Config.Get(ctx, &data); diags.HasError() {
		stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	if listR.Client == nil {
		// diagnostics: provider not configured, then stream.Results = list.ListResultsStreamDiagnostics(diags)
		return
	}
	project := listR.GetProject(data.Project)

	stream.Results = func(push func(list.ListResult) bool) {
		// callback runs once per item returned by the list API (ListPages invokes it after the flattener).
		err := ListExamples(listR.Client, project, func(rd *schema.ResourceData) error {
			result := listReq.NewListResult(ctx)
			if err := listR.SetResult(ctx, listReq.IncludeResource, &result, rd, "display_name", "email"); err != nil {
				return err
			}
			if !push(result) {
				return errors.New("stream closed")
			}
			return nil
		})
		if err != nil {
			// attach diagnostics and/or push a result carrying diagnostics (see list_google_service_account.go)
		}
	}
}
```

Add the usual imports (`context`, `errors`, `list`, `schema`, `tpgresource`, etc.). The **`ListExamples`** name is illustrative—implement **`List<YourApiCollection>`** in the same file as in the following section.

**Important:** **`SetResult`** expects `ResourceData` consistent with **`SDKv2Resource`** (same schema and identity).

### 2. Add `List<Resource>` and `flatten<Resource>ListItem` (ListPages)

Enumeration is **not** implemented inline inside **`List`**. Instead, add a dedicated function (naming convention **`List<ServiceAccounts>`**, **`ListExamples`**, etc.) that contains everything needed for the Google API **list** call: temporary **`ResourceData`** for URL templates, **`tpgresource.ReplaceVars`** for the list URL, billing project and user agent (**`tpgresource.GetBillingProject`**, **`tpgresource.GenerateUserAgentString`**), and a call to **[`transport_tpg.ListPages`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/transport/transport.go)** with:

- **`ListURL`** — resolved list endpoint (often uses the same path patterns as the managed resource’s read URL).
- **`ItemName`** — JSON field name for the array of items in the list response (for example `"accounts"`).
- **`Flattener`** — typically named **`flattenGoogle<Resource>ListItem`**: converts one **`map[string]interface{}`** item from the API into a **`ResourceData`** for the managed resource (for example **`tpgresource.Convert`** into a typed struct, **`d.SetId`**, then **`populateResourceData`** or the same helpers the resource’s **Read** uses).
- **`Callback`** — see the callout below.

> **Per-item callback:** The function parameter **`callback`** (`func(rd *schema.ResourceData) error`) is the logic that runs **for each element** of the list API response (after paging). **`ListPages`** walks the HTTP response, unmarshals each item, runs the **flattener** to produce a row’s **`ResourceData`**, then calls **`callback(rd)`** with that value. Your list resource’s **`List`** method passes a closure as **`callback`**: that closure is where each row becomes a **`ListResult`** (**`NewListResult`**, **`SetResult`**, **`push`**). Returning an error from **`callback`** stops iteration; the outer **`List<ServiceAccounts>`**-style wrapper forwards that error to **`List`**.

Full example (service accounts — same file as the list resource):

```go
func flattenGoogleServiceAccountListItem(res map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) error {
	var sa iam.ServiceAccount
	if err := tpgresource.Convert(res, &sa); err != nil {
		return err
	}
	d.SetId(sa.Name)
	return populateResourceData(d, &sa)
}

func ListServiceAccounts(config *transport_tpg.Config, project string, callback func(rd *schema.ResourceData) error) error {
	if config == nil {
		return fmt.Errorf("provider client is not configured")
	}
	d := ResourceGoogleServiceAccount().Data(&terraform.InstanceState{})
	if project != "" {
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("error setting project on temporary resource data: %w", err)
		}
	}
	url, err := tpgresource.ReplaceVars(d, config, "{{IAMBasePath}}projects/{{project}}/serviceAccounts")
	if err != nil {
		return err
	}

	billingProject := ""
	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}
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
		ItemName:       "accounts",
		Flattener:      flattenGoogleServiceAccountListItem,
		Callback:       callback,
	})
}
```

If your API already has a **`ListPages`**-style helper used by data sources or sweepers, reuse it from **`List<Resource>`** instead of duplicating HTTP logic.

### 3. Register the list resource

Append your constructor to **`handwrittenListResources`** in [`mmv1/third_party/terraform/fwprovider/framework_provider_mmv1_resources.go`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/fwprovider/framework_provider_mmv1_resources.go):

```go
var handwrittenListResources = []func() list.ListResource{
	listResourceFunc(resourcemanager.NewGoogleServiceAccountListResource()),
	listResourceFunc(resourcemanager.NewGoogleProjectServiceListResource()),
	listResourceFunc(yourpkg.NewGoogleExampleListResource()),
}
```

The provider merges **`generatedListResources`** and **`handwrittenListResources`** when building the framework provider (see **`framework_provider.go.tmpl`**). MMv1-generated list resources will eventually populate **`generatedListResources`**.

## Add tests

Use **`Query: true`** on a test step and **`querycheck`** expectations so the test exercises the list resource API (Terraform **1.14+**). Example pattern from [`list_google_service_account_test.go`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/services/resourcemanager/list_google_service_account_test.go):

```go
func TestAccExampleListResource_queryIdentity(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_14_0),
		},
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `...`, // create a known managed resource
			},
			{
				Query:  true,
				Config: testAccExampleListQuery(),
				QueryResultChecks: []querycheck.QueryResultCheck{
					querycheck.ExpectIdentity("google_example.all", map[string]knownvalue.Check{
						"id": knownvalue.StringExact("expected-id"),
					}),
					querycheck.ExpectLengthAtLeast("google_example.all", 1),
				},
			},
		},
	})
}
```

## Add documentation

1. Add a page under [`magic-modules/mmv1/third_party/terraform/website/docs/list-resources/`](https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/list-resources) named after the list type (for example `google_service_account.html.markdown`).
1. Document the `config` block, identity fields in results, and `include_resource`. Link to the managed resource docs for full attributes.
1. Link to [Use list resources with terraform query (Google Cloud provider)](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/using_list_resources_with_terraform_query).
1. Follow the [Handwritten documentation style guide]({{< ref "/document/handwritten-docs-style-guide" >}}) where it applies.
1. [Generate the providers]({{< ref "/develop/generate-providers" >}}) and optionally validate with the Registry [Doc Preview Tool](https://registry.terraform.io/tools/doc-preview).
