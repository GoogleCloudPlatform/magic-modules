// Contains functions that don't really belong anywhere else.

package google

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	fwDiags "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/googleapi"
)

type TerraformResourceDataChange interface {
	GetChange(string) (interface{}, interface{})
}

type TerraformResourceData interface {
	HasChange(string) bool
	GetOkExists(string) (interface{}, bool)
	GetOk(string) (interface{}, bool)
	Get(string) interface{}
	Set(string, interface{}) error
	SetId(string)
	Id() string
	GetProviderMeta(interface{}) error
	Timeout(key string) time.Duration
}

type TerraformResourceDiff interface {
	HasChange(string) bool
	GetChange(string) (interface{}, interface{})
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
	Clear(string) error
	ForceNew(string) error
}

// getRegionFromZone returns the region from a zone for Google cloud.
// This is by removing the last two chars from the zone name to leave the region
// If there aren't enough characters in the input string, an empty string is returned
// e.g. southamerica-west1-a => southamerica-west1
func getRegionFromZone(zone string) string {
	return tpgresource.GetRegionFromZone(zone)
}

// Infers the region based on the following (in order of priority):
// - `region` field in resource schema
// - region extracted from the `zone` field in resource schema
// - provider-level region
// - region extracted from the provider-level zone
func getRegion(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return getRegionFromSchema("region", "zone", d, config)
}

// getProject reads the "project" field from the given resource data and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProject(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return getProjectFromSchema("project", d, config)
}

// getBillingProject reads the "billing_project" field from the given resource data and falls
// back to the provider's value if not given. If no value is found, an error is returned.
func getBillingProject(d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return getBillingProjectFromSchema("billing_project", d, config)
}

// getProjectFromDiff reads the "project" field from the given diff and falls
// back to the provider's value if not given. If the provider's value is not
// given, an error is returned.
func getProjectFromDiff(d *schema.ResourceDiff, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetProjectFromDiff(d, config)
}

func getRouterLockName(region string, router string) string {
	return tpgresource.GetRouterLockName(region, router)
}

func isFailedPreconditionError(err error) bool {
	return tpgresource.IsFailedPreconditionError(err)
}

func isConflictError(err error) bool {
	return tpgresource.IsConflictError(err)
}

// gRPC does not return errors of type *googleapi.Error. Instead the errors returned are *status.Error.
// See the types of codes returned here (https://pkg.go.dev/google.golang.org/grpc/codes#Code).
func isNotFoundGrpcError(err error) bool {
	return tpgresource.IsNotFoundGrpcError(err)
}

// expandLabels pulls the value of "labels" out of a TerraformResourceData as a map[string]string.
func expandLabels(d TerraformResourceData) map[string]string {
	return expandStringMap(d, "labels")
}

// expandEnvironmentVariables pulls the value of "environment_variables" out of a schema.ResourceData as a map[string]string.
func expandEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "environment_variables")
}

// expandBuildEnvironmentVariables pulls the value of "build_environment_variables" out of a schema.ResourceData as a map[string]string.
func expandBuildEnvironmentVariables(d *schema.ResourceData) map[string]string {
	return expandStringMap(d, "build_environment_variables")
}

// expandStringMap pulls the value of key out of a TerraformResourceData as a map[string]string.
func expandStringMap(d TerraformResourceData, key string) map[string]string {
	v, ok := d.GetOk(key)

	if !ok {
		return map[string]string{}
	}

	return convertStringMap(v.(map[string]interface{}))
}

func convertStringMap(v map[string]interface{}) map[string]string {
	return tpgresource.ConvertStringMap(v)
}

func convertStringArr(ifaceArr []interface{}) []string {
	return tpgresource.ConvertStringArr(ifaceArr)
}

func convertAndMapStringArr(ifaceArr []interface{}, f func(string) string) []string {
	return tpgresource.ConvertAndMapStringArr(ifaceArr, f)
}

func mapStringArr(original []string, f func(string) string) []string {
	return tpgresource.MapStringArr(original, f)
}

func convertStringArrToInterface(strs []string) []interface{} {
	return tpgresource.ConvertStringArrToInterface(strs)
}

func convertStringSet(set *schema.Set) []string {
	return tpgresource.ConvertStringSet(set)
}

func golangSetFromStringSlice(strings []string) map[string]struct{} {
	return tpgresource.GolangSetFromStringSlice(strings)
}

func stringSliceFromGolangSet(sset map[string]struct{}) []string {
	return tpgresource.StringSliceFromGolangSet(sset)
}

func reverseStringMap(m map[string]string) map[string]string {
	return tpgresource.ReverseStringMap(m)
}

func mergeStringMaps(a, b map[string]string) map[string]string {
	return tpgresource.MergeStringMaps(a, b)
}

func mergeSchemas(a, b map[string]*schema.Schema) map[string]*schema.Schema {
	return tpgresource.MergeSchemas(a, b)
}

func StringToFixed64(v string) (int64, error) {
	return tpgresource.StringToFixed64(v)
}

func extractFirstMapConfig(m []interface{}) map[string]interface{} {
	return tpgresource.ExtractFirstMapConfig(m)
}

func lockedCall(lockKey string, f func() error) error {
	mutexKV.Lock(lockKey)
	defer mutexKV.Unlock(lockKey)

	return f()
}

// This is a Printf sibling (Nprintf; Named Printf), which handles strings like
// Nprintf("Hello %{target}!", map[string]interface{}{"target":"world"}) == "Hello world!".
// This is particularly useful for generated tests, where we don't want to use Printf,
// since that would require us to generate a very particular ordering of arguments.
func Nprintf(format string, params map[string]interface{}) string {
	return tpgresource.Nprintf(format, params)
}

// serviceAccountFQN will attempt to generate the fully qualified name in the format of:
// "projects/(-|<project>)/serviceAccounts/<service_account_id>@<project>.iam.gserviceaccount.com"
// A project is required if we are trying to build the FQN from a service account id and
// and error will be returned in this case if no project is set in the resource or the
// provider-level config
func serviceAccountFQN(serviceAccount string, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	// If the service account id is already the fully qualified name
	if strings.HasPrefix(serviceAccount, "projects/") {
		return serviceAccount, nil
	}

	// If the service account id is an email
	if strings.Contains(serviceAccount, "@") {
		return "projects/-/serviceAccounts/" + serviceAccount, nil
	}

	// Get the project from the resource or fallback to the project
	// in the provider configuration
	project, err := getProject(d, config)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("projects/-/serviceAccounts/%s@%s.iam.gserviceaccount.com", serviceAccount, project), nil
}

func paginatedListRequest(project, baseUrl, userAgent string, config *transport_tpg.Config, flattener func(map[string]interface{}) []interface{}) ([]interface{}, error) {
	return tpgresource.PaginatedListRequest(project, baseUrl, userAgent, config, flattener)
}

func getInterconnectAttachmentLink(config *transport_tpg.Config, project, region, ic, userAgent string) (string, error) {
	return tpgresource.GetInterconnectAttachmentLink(config, project, region, ic, userAgent)
}

// Given two sets of references (with "from" values in self link form),
// determine which need to be added or removed // during an update using
// addX/removeX APIs.
func calcAddRemove(from []string, to []string) (add, remove []string) {
	add = make([]string, 0)
	remove = make([]string, 0)
	for _, u := range to {
		found := false
		for _, v := range from {
			if tpgresource.CompareSelfLinkOrResourceName("", v, u, nil) {
				found = true
				break
			}
		}
		if !found {
			add = append(add, u)
		}
	}
	for _, u := range from {
		found := false
		for _, v := range to {
			if tpgresource.CompareSelfLinkOrResourceName("", u, v, nil) {
				found = true
				break
			}
		}
		if !found {
			remove = append(remove, u)
		}
	}
	return add, remove
}

func stringInSlice(arr []string, str string) bool {
	return tpgresource.StringInSlice(arr, str)
}

func migrateStateNoop(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return tpgresource.MigrateStateNoop(v, is, meta)
}

func expandString(v interface{}, d TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return v.(string), nil
}

func changeFieldSchemaToForceNew(sch *schema.Schema) {
	tpgresource.ChangeFieldSchemaToForceNew(sch)
}

func generateUserAgentString(d TerraformResourceData, currentUserAgent string) (string, error) {
	var m transport_tpg.ProviderMeta

	err := d.GetProviderMeta(&m)
	if err != nil {
		return currentUserAgent, err
	}

	if m.ModuleName != "" {
		return strings.Join([]string{currentUserAgent, m.ModuleName}, " "), nil
	}

	return currentUserAgent, nil
}

func SnakeToPascalCase(s string) string {
	split := strings.Split(s, "_")
	for i := range split {
		split[i] = strings.Title(split[i])
	}
	return tpgresource.SnakeToPascalCase(s)
}

func checkStringMap(v interface{}) map[string]string {
	return tpgresource.CheckStringMap(v)
}

// return a fake 404 so requests get retried or nested objects are considered deleted
func fake404(reasonResourceType, resourceName string) *googleapi.Error {
	return tpgresource.Fake404(reasonResourceType, resourceName)
}

// validate name of the gcs bucket. Guidelines are located at https://cloud.google.com/storage/docs/naming-buckets
// this does not attempt to check for IP addresses or close misspellings of "google"
func checkGCSName(name string) error {
	return tpgresource.CheckGCSName(name)
}

// checkGoogleIamPolicy makes assertions about the contents of a google_iam_policy data source's policy_data attribute
func checkGoogleIamPolicy(value string) error {
	return tpgresource.CheckGoogleIamPolicy(value)
}

// Retries an operation while the canonical error code is FAILED_PRECONDTION
// which indicates there is an incompatible operation already running on the
// cluster. This error can be safely retried until the incompatible operation
// completes, and the newly requested operation can begin.
func retryWhileIncompatibleOperation(timeout time.Duration, lockKey string, f func() error) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		if err := lockedCall(lockKey, f); err != nil {
			if isFailedPreconditionError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
}

func frameworkDiagsToSdkDiags(fwD fwDiags.Diagnostics) *diag.Diagnostics {
	return tpgresource.FrameworkDiagsToSdkDiags(fwD)
}

// Deprecated: For backward compatibility isEmptyValue is still working,
// but all new code should use IsEmptyValue in the verify package instead.
func isEmptyValue(v reflect.Value) bool {
	return tpgresource.IsEmptyValue(v)
}
