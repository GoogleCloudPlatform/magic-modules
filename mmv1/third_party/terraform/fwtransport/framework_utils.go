package fwtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

const uaEnvVar = "TF_APPEND_USER_AGENT"

func CompileUserAgentString(ctx context.Context, name, tfVersion, provVersion string) string {
	ua := fmt.Sprintf("Terraform/%s (+https://www.terraform.io) Terraform-Plugin-SDK/%s %s/%s", tfVersion, "terraform-plugin-framework", name, provVersion)

	if add := os.Getenv(uaEnvVar); add != "" {
		add = strings.TrimSpace(add)
		if len(add) > 0 {
			ua += " " + add
			tflog.Debug(ctx, fmt.Sprintf("Using modified User-Agent: %s", ua))
		}
	}

	return ua
}

func GenerateFrameworkUserAgentString(metaData *fwmodels.ProviderMetaModel, currUserAgent string) string {
	if metaData != nil && !metaData.ModuleName.IsNull() && metaData.ModuleName.ValueString() != "" {
		return strings.Join([]string{currUserAgent, metaData.ModuleName.ValueString()}, " ")
	}

	return currUserAgent
}

func HandleDatasourceNotFoundError(ctx context.Context, err error, state *tfsdk.State, resource string, diags *diag.Diagnostics) {
	if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		tflog.Warn(ctx, fmt.Sprintf("Removing %s because it's gone", resource))
		// The resource doesn't exist anymore
		state.RemoveResource(ctx)
	}

	diags.AddError(fmt.Sprintf("Error when reading or editing %s", resource), err.Error())
}

var DefaultRequestTimeout = 5 * time.Minute

type SendRequestOptions struct {
	Config               *transport_tpg.Config
	Method               string
	Project              string
	RawURL               string
	UserAgent            string
	Body                 map[string]any
	Timeout              time.Duration
	Headers              http.Header
	ErrorRetryPredicates []transport_tpg.RetryErrorPredicateFunc
	ErrorAbortPredicates []transport_tpg.RetryErrorPredicateFunc
}

func SendRequest(opt SendRequestOptions, diags *diag.Diagnostics) (map[string]interface{}) {
	reqHeaders := opt.Headers
	if reqHeaders == nil {
		reqHeaders = make(http.Header)
	}
	reqHeaders.Set("User-Agent", opt.UserAgent)
	reqHeaders.Set("Content-Type", "application/json")

	if opt.Config.UserProjectOverride && opt.Project != "" {
		// When opt.Project is "NO_BILLING_PROJECT_OVERRIDE" in the function GetCurrentUserEmail,
		// set the header X-Goog-User-Project to be empty string.
		if opt.Project == "NO_BILLING_PROJECT_OVERRIDE" {
			reqHeaders.Set("X-Goog-User-Project", "")
		} else {
			// Pass the project into this fn instead of parsing it from the URL because
			// both project names and URLs can have colons in them.
			reqHeaders.Set("X-Goog-User-Project", opt.Project)
		}
	}

	if opt.Timeout == 0 {
		opt.Timeout = DefaultRequestTimeout
	}

	var res *http.Response
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			var buf bytes.Buffer
			if opt.Body != nil {
				err := json.NewEncoder(&buf).Encode(opt.Body)
				if err != nil {
					diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
					return nil
				}
			}

			u, err := transport_tpg.AddQueryParams(opt.RawURL, map[string]string{"alt": "json"})
			if err != nil {
				diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
				return nil
			}
			req, err := http.NewRequest(opt.Method, u, &buf)
			if err != nil {
				diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
				return nil
			}

			req.Header = reqHeaders
			res, err = opt.Config.Client.Do(req)
			if err != nil {
				diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
				return nil
			}

			if err := googleapi.CheckResponse(res); err != nil {
				googleapi.CloseBody(res)
				diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
				return nil
			}

			return nil
		},
		Timeout:              opt.Timeout,
		ErrorRetryPredicates: opt.ErrorRetryPredicates,
		ErrorAbortPredicates: opt.ErrorAbortPredicates,
	})
	if err != nil {
		diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
		return nil
	}

	if res == nil {
		diags.AddError("Unable to parse server response. This is most likely a terraform problem, please file a bug at https://github.com/hashicorp/terraform-provider-google/issues.")
		return nil
	}

	// The defer call must be made outside of the retryFunc otherwise it's closed too soon.
	defer googleapi.CloseBody(res)

	// 204 responses will have no body, so we're going to error with "EOF" if we
	// try to parse it. Instead, we can just return nil.
	if res.StatusCode == 204 {
		return nil
	}
	result := make(map[string]interface{})
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		diags.AddError(fmt.Sprintf("Error when sending HTTP request %s", resource), err.Error())
		return nil
	}

	return result
}

func ReplaceVars(ctx context.Context, req interface{}, diags *diag.Diagnostics, data interface{}, config *transport_tpg.Config, linkTmpl string) (string) {
	return ReplaceVarsRecursive(ctx, req, resp, data, config, linkTmpl, false, 0)
}

// relaceVarsForId shortens variables by running them through GetResourceNameFromSelfLink
// this allows us to use long forms of variables from configs without needing
// custom id formats. For instance:
// accessPolicies/{{access_policy}}/accessLevels/{{access_level}}
// with values:
// access_policy: accessPolicies/foo
// access_level: accessPolicies/foo/accessLevels/bar
// becomes accessPolicies/foo/accessLevels/bar
func ReplaceVarsForId(ctx context.Context, req interface{}, diags *diag.Diagnostics, data interface{}, config *transport_tpg.Config, linkTmpl string) (string) {
	return ReplaceVarsRecursive(ctx, req, resp, data, config, linkTmpl, true, 0)
}

// ReplaceVars must be done recursively because there are baseUrls that can contain references to regions
// (eg cloudrun service) there aren't any cases known for 2+ recursion but we will track a run away
// substitution as 10+ calls to allow for future use cases.
func ReplaceVarsRecursive(ctx context.Context, req interface{}, diags *diag.Diagnostics, data interface{}, config *transport_tpg.Config, linkTmpl string, shorten bool, depth int) (string) {
	if depth > 10 {
		diags.AddError("url building error", "Recursive substitution detected.")
	}

	// https://github.com/google/re2/wiki/Syntax
	re := regexp.MustCompile("{{([%[:word:]]+)}}")
	f := BuildReplacementFunc(ctx, req, diags, data, config, linkTmpl, shorten)
	if resp.Diagnostics.HasError() {
		return
	}
	final := re.ReplaceAllStringFunc(linkTmpl, f)

	if re.Match([]byte(final)) {
		return ReplaceVarsRecursive(ctx, req, diags, data, config, final, shorten, depth+1)
	}

	return final
}

// This function replaces references to Terraform properties (in the form of {{var}}) with their value in Terraform
// It also replaces {{project}}, {{project_id_or_project}}, {{region}}, and {{zone}} with their appropriate values
// This function supports URL-encoding the result by prepending '%' to the field name e.g. {{%var}}
func BuildReplacementFunc(ctx context.Context, re *regexp.Regexp, req interface{}, diags *diag.Diagnostics, data interface{}, config *transport_tpg.Config, linkTmpl string, shorten bool) (func(string) string, error) {
	var project, projectID, region, zone string
	var err error

	if strings.Contains(linkTmpl, "{{project}}") {
		project, err = fwresource.GetProjectFramework(data.Project, types.StringValue(config.Project), diags)
		if diags.HasError() {
			return
		}
		if shorten {
			project = strings.TrimPrefix(project, "projects/")
		}
	}

	if strings.Contains(linkTmpl, "{{project_id_or_project}}") {
		switch req.(type) {
			case resource.CreateRequest || resource.UpdateRequest:
		 		diagInfo := req.Plan.GetAttribute(ctx, path.Root("project_id"), projectID)
			    diags.Append(diagsInfo...)
			    if diags.HasError() {
			        return
			    }
			    if !pid {
	    			project = fwresource.GetProjectFramework(data.Project, types.StringValue(config.Project), diags)
					if diags.HasError() {
						return
					}
			    }
				if shorten {
					project = strings.TrimPrefix(project, "projects/")
					projectID = strings.TrimPrefix(projectID, "projects/")
				}
			case resource.ReadRequest || resource.DeleteRequest:
		 		diagInfo := req.State.GetAttribute(ctx, path.Root("project_id"), projectID)
			    diags.Append(diagsInfo...)
			    if diags.HasError() {
			        return
			    }
			    if !pid {
	    			project = fwresource.GetProjectFramework(data.Project, types.StringValue(config.Project), diags)
					if diags.HasError() {
						return
					}
			    }
				if shorten {
					project = strings.TrimPrefix(project, "projects/")
					projectID = strings.TrimPrefix(projectID, "projects/")
				}
		}
	}

	if strings.Contains(linkTmpl, "{{region}}") {
 		region = fwresource.GetRegionFramework(data.Region, types.StringValue(config.Region), diags)
		if diags.HasError() {
			return
	    }
		if shorten {
			region = strings.TrimPrefix(region, "regions/")
		}
	}

	if strings.Contains(linkTmpl, "{{zone}}") {
 		zone = fwresource.GetRegionFramework(data.Zone, types.StringValue(config.Zone), diags)
		if diags.HasError() {
			return
	    }
		if shorten {
			zone = strings.TrimPrefix(region, "zones/")
		}
	}

	f := func(s string) string {

		m := re.FindStringSubmatch(s)[1]
		if m == "project" {
			return project
		}
		if m == "project_id_or_project" {
			if projectID != "" {
				return projectID
			}
			return project
		}
		if m == "region" {
			return region
		}
		if m == "zone" {
			return zone
		}
		if string(m[0]) == "%" {
			v, ok := d.GetOkExists(m[1:])
			if ok {
				return url.PathEscape(fmt.Sprintf("%v", v))
			}
		} else {
			v, ok := d.GetOkExists(m)
			if ok {
				if shorten {
					return GetResourceNameFromSelfLink(fmt.Sprintf("%v", v))
				} else {
					return fmt.Sprintf("%v", v)
				}
			}
		}

		// terraform-google-conversion doesn't provide a provider config in tests.
		if config != nil {
			// Attempt to draw values from the provider config if it's present.
			if f := reflect.Indirect(reflect.ValueOf(config)).FieldByName(m); f.IsValid() {
				return f.String()
			}
		}
		return ""
	}

	return f, nil
}