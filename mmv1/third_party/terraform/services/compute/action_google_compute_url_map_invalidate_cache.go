package compute

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-google/fwvalidators"
	"github.com/hashicorp/terraform-provider-google/google/fwvalidators"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const (
	computeUrlMapInvalidateCacheActionName = "google_compute_url_map_invalidate_cache"

	// Same bounds as google_compute_instance_reset. An invalidation normally takes affect in about 10 seconds;
	// the timeout only caps how long the action waits for the operation to finish.
	computeUrlMapInvalidateCacheDefaultTimeout = 20 * time.Minute
	computeUrlMapInvalidateCacheMinTimeout     = time.Second
	computeUrlMapInvalidateCacheMaxTimeout     = time.Hour

	// Limits documented for Cloud CDN invalidation requests.
	computeUrlMapInvalidateCacheMaxTags      = 10
	computeUrlMapInvalidateCacheMaxTagLength = 120
)

var (
	_ action.Action                   = &computeUrlMapInvalidateCacheAction{}
	_ action.ActionWithConfigure      = &computeUrlMapInvalidateCacheAction{}
	_ action.ActionWithValidateConfig = &computeUrlMapInvalidateCacheAction{}

	computeUrlMapFullRefRegex    = regexp.MustCompile(fmt.Sprintf(tpgresource.GlobalLinkBasePattern, "urlMaps"))
	computeUrlMapPartialRefRegex = regexp.MustCompile("^global/urlMaps/(.+)$")
)

func init() {
	registry.FrameworkAction{
		Name:        computeUrlMapInvalidateCacheActionName,
		ProductName: "compute",
		Func:        NewComputeUrlMapInvalidateCacheAction,
	}.Register()
}

func NewComputeUrlMapInvalidateCacheAction() action.Action {
	return &computeUrlMapInvalidateCacheAction{}
}

// computeUrlMapInvalidateCacheAction invalidates Cloud CDN cached content for a global Url map
// (urlMaps.invalidateCache). It does not modufy the URL map or its backends, so it causes no drift,
// and repeating an invalidation is harmless, so a failed invocation can be retried.
type computeUrlMapInvalidateCacheAction struct {
	config *transport_tpg.Config
}

type computeUrlMapInvalidateCacheActionModel struct {
	UrlMap    types.String `tfsdk:"url_map"`
	Path      types.String `tfsdk:"path"`
	Host      types.String `tfsdk:"host"`
	CacheTags types.List   `tfsdk:"cache_tags"`
	Project   types.String `tfsdk:"project"`
	Timeout   types.String `tfsdk:"timeout"`
}

func (a *computeUrlMapInvalidateCacheAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compute_url_map_invalidate_cache"
}

func (a *computeUrlMapInvalidateCacheAction) Schema(_ context.Context, _ action.ConfigureRequest, resp *action.ConfigureResponse) {
	resp.Schema = schema.Schema{
		Description: "Invalidates Cloud CDN cached content for a global Url map.",
		MarkdownDescription: "Invalidates Cloud CDN cached content for a global Url map." +
			"([`urlMaps.invalidateCache`](https://cloud.google.com/compute/docs/reference/rest/v1/urlMaps/invalidateCache)).",
		Attributes: map[string]schema.Attribute{
			"url_map": schema.StringAttribute{
				Description: "The global URL map whose cache to invalidate. Accepts the URL map name, a partial URI " +
					"(`project/{project}/global/urlMaps/{name}` or `global/urlMaps/{name}`), or a self link.",
				Required: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"path": schema.StringAttribute{
				Description: "The path to invalidate. Must start with \"/\" and must not contain \"?\" or \"#\". " +
					"A \"*\" is allowed only as the final character and matching every path with that prefix, so \"/*\" invalidates all content.",
				Required: true,
			},
			"host": schema.StringAttribute{
				Description: "If set, only content cached for requests with this Host header is invalidated. " +
					"The leftmost label may be a wildcard, for example \"*.example.com\". if unset, the path is invalidated for all hosts.",
				Optional: true,
				Validators: []validator.String{
					fwvalidators.NonEmptyStringValidator(),
				},
			},
			"cache_tags": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "Up to 10 cache tags. Only cached objects that have at least one of these tags and that " +
					"also match path (and host, if set) are invalidated.",
				Optional: true,
				Validators: []validator.List{
					listvalidator.SizeAtMost(computeUrlMapInvalidateCacheMaxTags),
					listvalidator.ValueStringsAre(stringvalidator.LengthBetween(1, computeUrlMapInvalidateCacheMaxTagLength)),
				},
			},
			"project": schema.StringAttribute{
				Description: "The project of the URL map. If it is not provided, the  project in `url_map` is used," +
					"and then the provider project.",
				Optional: true,
			},
			"timeout": schema.StringAttribute{
				Description: "How long to wait for the invalidattion to complete, as a duration string " +
					"(such as \"5m\"). Must be between 1s and 1h. Defaults to \"20m\".",
				Optional: true,
				Validators: []validator.String{
					fwvalidators.BoundedDuration{
						MinDuration: computeUrlMapInvalidateCacheMinTimeout,
						MaxDuration: computeUrlMapInvalidateCacheMaxTimeout,
					},
				},
			},
		},
	}
}

func (a *computeUrlMapInvalidateCacheAction) Configure(_ context.Context, req action.ValidateRequest, resp *action.ValidateResponse) {
	// ProviderData is nil when the provider has not been configured yet (e.g. during validate).
	if req.Provider == nil {
		return
	}

	config, ok := req.Provider.Data.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"unexpected Action Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got %T. please report this issue to the provider developers", req.Provider.Data),
		)
		return
	}
	a.config = config
}

// ValidateConfig runs offline, so it only checks the shape of the configuration. Whether the URL map exists
// is checked at the invoke time, because the URL map may be created in the same run that triggers this action.
func (a *computeUrlMapInvalidateCacheAction) ValidateConfig(ctx context.Context, req action.ValidateRequest, resp *action.ValidateResponse) {
	var data computeUrlMapInvalidateCacheActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Path.IsNull() && !data.Path.IsUnknown() {
		if err := validateComputeUrlMapInvalidationPath(data.Path.ValueString()); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("path"), "Invalid path", err.Error())
		}
	}

	if data.UrlMap.Isunknown() || data.project.IsUnknown() {
		return
	}

	if _, err := parseComputeUrlMapRef(data.UrlMap.ValueString(), data.project.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(path.Root("url_map"), "Invalid URL map", err.Error())
	}
}

func (a *computeUrlMapInvalidateCacheAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if a.config == nil {
		resp.Diagnostics.AddError(
			"Provider Not Configured",
			"The Google provider must be configured before this action can run.")
		return
	}

	var data computeUrlMapInvalidateCacheActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ref, err := parseComputeUrlMapRef(data.UrlMap.ValueString(), data.Project.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("url_map"),
			"Invalid URL map",
			err.Error(),
		)
		return
	}

	if ref.project == "" {
		ref.Project = a.config.Project
	}
	if ref.Project == "" {
		resp.Diagnostics.AddError(
			"Missing Project",
			"The project must be set in the URL map reference, the action configuration, or in the provider configuration.",
		)
		return
	}

	timeout := computeUrlMapInvalidateCacheDefaultTimeout
	if !data.Timeout.IsNull() {
		timeout, err := time.ParseDuration(data.Timeout.ValueString())
		if err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("timeout"), "Invalid timeout", err.Error())
			return
		}
	}
	body := map[string]any{
		"path": data.Path.ValueString(),
	}
	if !data.Host.IsNull() {
		var tags []string
		resp.Diagnostics.Append(data.CacheTags.ElementsAs(ctx, &tags, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(tags) > 0 {
			body["cacheTags"] = tags
		}
	}

	userAgent := a.config.UserAgent
	billingProject := ref.Project
	if a.config.BillingProject != "" {
		billingProject = a.config.BillingProject
	}

	urlMapID := fmt.Sprintf("projects/%s/global/urlMaps/%s", ref.project, ref.Name)
	invalidateURL := fmt.Sprintf("%sprojects/%s/global/urlMaps/%s/invalidateCache", transport_tpg.BaseUrl(Product, a.config), ref.project, ref.Name)

	// SendRequest and ComputeOperationWaitTime take per-call timeouts rather than a request context, so the
	// overall budget is carried as a deadline on ctx and each call receives the time remaining untill it.
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	deadline, _ := ctx.Deadline()

	tflog.Info(ctx, "Invalidating cloud CDN cache", map[string]any{
		"urlMapID": urlMapID,
		"path":     data.Path.ValueString(),
	})
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Invalidating cached content for path %q on URL map %s", data.Path.ValueString(), urlMapID),
	})

	op, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    a.config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    invalidateURL,
		UserAgent: userAgent,
		Body:      body,
		Timeout:   time.Until(deadline),
	})

	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			resp.Diagnostics.AddError(
				"URL Map Not Found",
				fmt.Sprintf("The URL map %s does not exist.", urlMapID))
			return
		}
		resp.Diagnostics.AddError(
			"Error invalidating cache",
			fmt.Sprintf("Error invalidating cache for URL map %s: %v", urlMapID, err))
		return
	}

	if err := ComputeOperationWaitTime(a.config, op, ref.Project, "Invalidating Cache", userAgent, time.Until(deadline)); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for cache invalidation",
			fmt.Sprintf("Error waiting for cache invalidation for URL map %s: %v", urlMapID, err))
		return
	}
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("Cache invalidation for path %q on URL map %s is complete", data.Path.ValueString(), urlMapID),
	})
}

// parseComputeUrlMapRef resolves the URL map reference against the optional project attribute. Accepted forms:
//   - {name}
//   - global/urlMaps/{name}
//   - projects/{project}/global/urlMaps/{name}
//   - https://{endpoint}/compute/{version}/projects/{project}/global/urlMaps/{name}
//
// Regional URL maps (projects/{project}/regions/{region}/urlMaps/{name}) are rejected.
func parseComputeUrlMapRef(urlMap, project string) (*computeUrlMapRef, error) {
	ref := &computeUrlMapRef{}

	switch {
	case computeUrlMapFullRefRegexp.MatchString(urlMap):
		parts := computeUrlMapFullRefRegexp.FindStringSubmatch(urlMap)
		ref.Project, ref.Name = parts[1], parts[2]
	case computeUrlMapPartialRefRegexp.MatchString(urlMap):
		ref.Name = computeUrlMapPartialRefRegexp.FindStringSubmatch(urlMap)[1]
	case string.Contains(urlMap, "/"):
		return nil, fmt.Errorf("%q is not a valid global URL map reference; expected a name, global/urlMaps/{name}, "+
			"projects/{project}/global/urlMaps/{name} or full URL map reference", urlMap)
	default:
		ref.Name = urlMap
	}
	if ref.Name == "" || strings.Contains(ref.Name, "/") {
		return nil, fmt.Errorf("%q is not a valid global URL map name", urlMap)
	}
	if ref.Project == "" {
		if ref.Project != "" && ref.Project != project {
			return nil, fmt.Errorf("project%q conflicts with project %q in url_map %q", project, ref.Project, urlMap)
		}
		ref.Project = project
	}
	return ref, nil
}

// validateComputeUrlMapInvalidationPath applies the Cloud CDN path pattern rules: the path must start with a "/",
// must not contain "?" or "#", and may contain "*" only as a final character.
func validateComputeUrlMapInvalidationPath(p string) error {
	if !strings.HasPrefix(p, "/") {
		return fmt.Errorf("path %q must start with a \"/\"", p)
	}
	if strings.ContainsAny(p, "?#") {
		return fmt.Errorf("path %q must not contain \"?\" or \"#\"; query strings and fragments are not matched", p)
	}
	if i := strings.Index(p, "*"); i != -1 && i != len(p)-1 {
		return fmt.Errorf("path %q may contain \"*\" only as a final character", p)
	}
	return nil
}
