package compute

import (
	"context"
	"fmt"
	"reflect"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	"google.golang.org/api/googleapi"
)

// Fields in "backends" that are not allowed for non-managed backend services
// (loadBalancingScheme) - the API returns an error if they are set at all
// in the request.
var backendServiceOnlyManagedFieldNames = []string{
	"capacity_scaler",
	"max_connections",
	"max_connections_per_instance",
	"max_connections_per_endpoint",
	"max_rate",
	"max_rate_per_instance",
	"max_rate_per_endpoint",
	"max_utilization",
}

// validateManagedBackendServiceBackends ensures capacity_scaler is set for each backend in a managed
// backend service. To prevent a permadiff, we decided to override the API behavior and require the
// capacity_scaler value in this case.
//
// The API:
// - requires the sum of the backends' capacity_scalers be > 0
// - defaults to 1 if capacity_scaler is omitted from the request
//
// However, the schema.Set Hash function defaults to 0 if not given, which we chose because it's the default
// float and because non-managed backends can't have the value set, so there will be a permadiff for a
// situational non-zero default returned from the API. We can't diff suppress or customdiff a
// field inside a set object in ResourceDiff, since the value also determines the hash for that set object.
func validateManagedBackendServiceBackends(backends []interface{}, d *schema.ResourceDiff) error {
	sum := 0.0

	for _, b := range backends {
		if b == nil {
			continue
		}
		backend := b.(map[string]interface{})
		if v, ok := backend["capacity_scaler"]; ok && v != nil {
			sum += v.(float64)
		} else {
			return fmt.Errorf("capacity_scaler is required for each backend in managed backend service")
		}
	}
	if sum == 0.0 {
		return fmt.Errorf("managed backend service must have at least one non-zero capacity_scaler for backends")
	}
	return nil
}

// If INTERNAL or EXTERNAL, make sure the user did not provide values for any of the fields that cannot be sent.
// We ignore these values regardless when sent to the API, but this adds plan-time validation if a
// user sets the value to non-zero. We can't validate for empty but set because
// of how the SDK handles set objects (on any read, nil fields will get set to empty values)
func validateNonManagedBackendServiceBackends(backends []interface{}, d *schema.ResourceDiff) error {
	for _, b := range backends {
		if b == nil {
			continue
		}
		backend := b.(map[string]interface{})
		for _, fn := range backendServiceOnlyManagedFieldNames {
			if v, ok := backend[fn]; ok && !tpgresource.IsEmptyValue(reflect.ValueOf(v)) {
				return fmt.Errorf("%q cannot be set for non-managed backend service, found value %v", fn, v)
			}
		}
	}
	return nil
}

func customDiffRegionBackendService(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	v, ok := d.GetOk("backend")
	if !ok {
		return nil
	}
	if v == nil {
		return nil
	}

	backends := v.(*schema.Set).List()
	if len(backends) == 0 {
		return nil
	}

	switch d.Get("load_balancing_scheme").(string) {
	case "INTERNAL", "EXTERNAL":
		return validateNonManagedBackendServiceBackends(backends, d)
	case "INTERNAL_MANAGED":
		return nil
	default:
		return validateManagedBackendServiceBackends(backends, d)
	}
}

// ComputeRegionBackendServiceAssetType is the CAI asset type name.
const ComputeRegionBackendServiceAssetType string = "compute.googleapis.com/RegionBackendService"

// ComputeRegionBackendServiceSchemaName is a TF resource schema name.
const ComputeRegionBackendServiceSchemaName string = "google_compute_region_backend_service"

type ComputeRegionBackendServiceConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewComputeRegionBackendServiceConverter returns an HCL converter for compute backend service.
func NewComputeRegionBackendServiceConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[ComputeRegionBackendServiceSchemaName].Schema

	return &ComputeRegionBackendServiceConverter{
		name:   ComputeRegionBackendServiceSchemaName,
		schema: schema,
	}
}

func (c *ComputeRegionBackendServiceConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
	var blocks []*common.HCLResourceBlock
	config := common.NewConfig()

	for _, asset := range assets {
		if asset == nil {
			continue
		}
		if asset.Resource != nil && asset.Resource.Data != nil {
			block, err := c.convertResourceData(asset, config)
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, block)
		}
	}
	return blocks, nil
}

func (c *ComputeRegionBackendServiceConverter) convertResourceData(asset *caiasset.Asset, config *transport_tpg.Config) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	assetResourceData := asset.Resource.Data

	hcl, _ := resourceComputeRegionBackendServiceRead(assetResourceData, config)

	ctyVal, err := common.MapToCtyValWithSchema(hcl, c.schema)
	if err != nil {
		return nil, err
	}

	resourceName := assetResourceData["name"].(string)

	return &common.HCLResourceBlock{
		Labels: []string{c.name, resourceName},
		Value:  ctyVal,
	}, nil
}

func resourceComputeRegionBackendServiceRead(resource map[string]interface{}, config *transport_tpg.Config) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var resource_data *schema.ResourceData = nil

	result["affinity_cookie_ttl_sec"] = flattenComputeRegionBackendServiceAffinityCookieTtlSec(resource["affinityCookieTtlSec"], resource_data, config)
	result["backend"] = flattenComputeRegionBackendServiceBackend(resource["backends"], resource_data, config)
	result["circuit_breakers"] = flattenComputeRegionBackendServiceCircuitBreakers(resource["circuitBreakers"], resource_data, config)
	result["consistent_hash"] = flattenComputeRegionBackendServiceConsistentHash(resource["consistentHash"], resource_data, config)
	result["cdn_policy"] = flattenComputeRegionBackendServiceCdnPolicy(resource["cdnPolicy"], resource_data, config)
	if flattenedProp := flattenComputeRegionBackendServiceConnectionDraining(resource["connectionDraining"], resource_data, config); flattenedProp != nil {
		if gerr, ok := flattenedProp.(*googleapi.Error); ok {
			return nil, fmt.Errorf("Error reading RegionBackendService: %s", gerr)
		}
		casted := flattenedProp.([]interface{})[0]
		if casted != nil {
			for k, v := range casted.(map[string]interface{}) {
				result[k] = v
			}
		}
	}
	result["creation_timestamp"] = flattenComputeRegionBackendServiceCreationTimestamp(resource["creationTimestamp"], resource_data, config)
	result["description"] = flattenComputeRegionBackendServiceDescription(resource["description"], resource_data, config)
	result["failover_policy"] = flattenComputeRegionBackendServiceFailoverPolicy(resource["failoverPolicy"], resource_data, config)
	result["enable_cdn"] = flattenComputeRegionBackendServiceEnableCDN(resource["enableCDN"], resource_data, config)
	result["fingerprint"] = flattenComputeRegionBackendServiceFingerprint(resource["fingerprint"], resource_data, config)
	result["health_checks"] = flattenComputeRegionBackendServiceHealthChecks(resource["healthChecks"], resource_data, config)
	result["iap"] = flattenComputeRegionBackendServiceIap(resource["iap"], resource_data, config)
	result["load_balancing_scheme"] = flattenComputeRegionBackendServiceLoadBalancingScheme(resource["loadBalancingScheme"], resource_data, config)
	result["locality_lb_policy"] = flattenComputeRegionBackendServiceLocalityLbPolicy(resource["localityLbPolicy"], resource_data, config)
	result["name"] = flattenComputeRegionBackendServiceName(resource["name"], resource_data, config)
	result["outlier_detection"] = flattenComputeRegionBackendServiceOutlierDetection(resource["outlierDetection"], resource_data, config)
	result["port_name"] = flattenComputeRegionBackendServicePortName(resource["portName"], resource_data, config)
	result["protocol"] = flattenComputeRegionBackendServiceProtocol(resource["protocol"], resource_data, config)
	result["security_policy"] = flattenComputeRegionBackendServiceSecurityPolicy(resource["securityPolicy"], resource_data, config)
	result["session_affinity"] = flattenComputeRegionBackendServiceSessionAffinity(resource["sessionAffinity"], resource_data, config)
	result["connection_tracking_policy"] = flattenComputeRegionBackendServiceConnectionTrackingPolicy(resource["connectionTrackingPolicy"], resource_data, config)
	result["timeout_sec"] = flattenComputeRegionBackendServiceTimeoutSec(resource["timeoutSec"], resource_data, config)
	result["log_config"] = flattenComputeRegionBackendServiceLogConfig(resource["logConfig"], resource_data, config)
	result["network"] = flattenComputeRegionBackendServiceNetwork(resource["network"], resource_data, config)
	result["subsetting"] = flattenComputeRegionBackendServiceSubsetting(resource["subsetting"], resource_data, config)
	result["region"] = flattenComputeRegionBackendServiceRegion(resource["region"], resource_data, config)

	return result, nil
}

func flattenComputeRegionBackendServiceAffinityCookieTtlSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceBackend(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := schema.NewSet(resourceGoogleComputeBackendServiceBackendHash, []interface{}{})
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed.Add(map[string]interface{}{
			"balancing_mode":               flattenComputeRegionBackendServiceBackendBalancingMode(original["balancingMode"], d, config),
			"capacity_scaler":              flattenComputeRegionBackendServiceBackendCapacityScaler(original["capacityScaler"], d, config),
			"description":                  flattenComputeRegionBackendServiceBackendDescription(original["description"], d, config),
			"failover":                     flattenComputeRegionBackendServiceBackendFailover(original["failover"], d, config),
			"group":                        flattenComputeRegionBackendServiceBackendGroup(original["group"], d, config),
			"max_connections":              flattenComputeRegionBackendServiceBackendMaxConnections(original["maxConnections"], d, config),
			"max_connections_per_instance": flattenComputeRegionBackendServiceBackendMaxConnectionsPerInstance(original["maxConnectionsPerInstance"], d, config),
			"max_connections_per_endpoint": flattenComputeRegionBackendServiceBackendMaxConnectionsPerEndpoint(original["maxConnectionsPerEndpoint"], d, config),
			"max_rate":                     flattenComputeRegionBackendServiceBackendMaxRate(original["maxRate"], d, config),
			"max_rate_per_instance":        flattenComputeRegionBackendServiceBackendMaxRatePerInstance(original["maxRatePerInstance"], d, config),
			"max_rate_per_endpoint":        flattenComputeRegionBackendServiceBackendMaxRatePerEndpoint(original["maxRatePerEndpoint"], d, config),
			"max_utilization":              flattenComputeRegionBackendServiceBackendMaxUtilization(original["maxUtilization"], d, config),
		})
	}
	return transformed
}
func flattenComputeRegionBackendServiceBackendBalancingMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceBackendCapacityScaler(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceBackendDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceBackendFailover(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceBackendGroup(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeRegionBackendServiceBackendMaxConnections(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceBackendMaxConnectionsPerInstance(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceBackendMaxConnectionsPerEndpoint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceBackendMaxRate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceBackendMaxRatePerInstance(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceBackendMaxRatePerEndpoint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceBackendMaxUtilization(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCircuitBreakers(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["connect_timeout"] =
		flattenComputeRegionBackendServiceCircuitBreakersConnectTimeout(original["connectTimeout"], d, config)
	transformed["max_requests_per_connection"] =
		flattenComputeRegionBackendServiceCircuitBreakersMaxRequestsPerConnection(original["maxRequestsPerConnection"], d, config)
	transformed["max_connections"] =
		flattenComputeRegionBackendServiceCircuitBreakersMaxConnections(original["maxConnections"], d, config)
	transformed["max_pending_requests"] =
		flattenComputeRegionBackendServiceCircuitBreakersMaxPendingRequests(original["maxPendingRequests"], d, config)
	transformed["max_requests"] =
		flattenComputeRegionBackendServiceCircuitBreakersMaxRequests(original["maxRequests"], d, config)
	transformed["max_retries"] =
		flattenComputeRegionBackendServiceCircuitBreakersMaxRetries(original["maxRetries"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceCircuitBreakersConnectTimeout(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeRegionBackendServiceCircuitBreakersConnectTimeoutSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeRegionBackendServiceCircuitBreakersConnectTimeoutNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceCircuitBreakersConnectTimeoutSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCircuitBreakersConnectTimeoutNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCircuitBreakersMaxRequestsPerConnection(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCircuitBreakersMaxConnections(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCircuitBreakersMaxPendingRequests(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCircuitBreakersMaxRequests(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCircuitBreakersMaxRetries(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceConsistentHash(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["http_cookie"] =
		flattenComputeRegionBackendServiceConsistentHashHttpCookie(original["httpCookie"], d, config)
	transformed["http_header_name"] =
		flattenComputeRegionBackendServiceConsistentHashHttpHeaderName(original["httpHeaderName"], d, config)
	transformed["minimum_ring_size"] =
		flattenComputeRegionBackendServiceConsistentHashMinimumRingSize(original["minimumRingSize"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceConsistentHashHttpCookie(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["ttl"] =
		flattenComputeRegionBackendServiceConsistentHashHttpCookieTtl(original["ttl"], d, config)
	transformed["name"] =
		flattenComputeRegionBackendServiceConsistentHashHttpCookieName(original["name"], d, config)
	transformed["path"] =
		flattenComputeRegionBackendServiceConsistentHashHttpCookiePath(original["path"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceConsistentHashHttpCookieTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeRegionBackendServiceConsistentHashHttpCookieTtlSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeRegionBackendServiceConsistentHashHttpCookieTtlNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceConsistentHashHttpCookieTtlSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceConsistentHashHttpCookieTtlNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceConsistentHashHttpCookieName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceConsistentHashHttpCookiePath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceConsistentHashHttpHeaderName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceConsistentHashMinimumRingSize(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["cache_key_policy"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicy(original["cacheKeyPolicy"], d, config)
	transformed["signed_url_cache_max_age_sec"] =
		flattenComputeRegionBackendServiceCdnPolicySignedUrlCacheMaxAgeSec(original["signedUrlCacheMaxAgeSec"], d, config)
	transformed["default_ttl"] =
		flattenComputeRegionBackendServiceCdnPolicyDefaultTtl(original["defaultTtl"], d, config)
	transformed["max_ttl"] =
		flattenComputeRegionBackendServiceCdnPolicyMaxTtl(original["maxTtl"], d, config)
	transformed["client_ttl"] =
		flattenComputeRegionBackendServiceCdnPolicyClientTtl(original["clientTtl"], d, config)
	transformed["negative_caching"] =
		flattenComputeRegionBackendServiceCdnPolicyNegativeCaching(original["negativeCaching"], d, config)
	transformed["negative_caching_policy"] =
		flattenComputeRegionBackendServiceCdnPolicyNegativeCachingPolicy(original["negativeCachingPolicy"], d, config)
	transformed["cache_mode"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheMode(original["cacheMode"], d, config)
	transformed["serve_while_stale"] =
		flattenComputeRegionBackendServiceCdnPolicyServeWhileStale(original["serveWhileStale"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["include_host"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeHost(original["includeHost"], d, config)
	transformed["include_protocol"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeProtocol(original["includeProtocol"], d, config)
	transformed["include_query_string"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeQueryString(original["includeQueryString"], d, config)
	transformed["query_string_blacklist"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyQueryStringBlacklist(original["queryStringBlacklist"], d, config)
	transformed["query_string_whitelist"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyQueryStringWhitelist(original["queryStringWhitelist"], d, config)
	transformed["include_named_cookies"] =
		flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeNamedCookies(original["includeNamedCookies"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeHost(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeProtocol(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeQueryString(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyQueryStringBlacklist(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyQueryStringWhitelist(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeRegionBackendServiceCdnPolicyCacheKeyPolicyIncludeNamedCookies(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCdnPolicySignedUrlCacheMaxAgeSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicyDefaultTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicyMaxTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicyClientTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicyNegativeCaching(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCdnPolicyNegativeCachingPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := make([]interface{}, 0, len(l))
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed = append(transformed, map[string]interface{}{
			"code": flattenComputeRegionBackendServiceCdnPolicyNegativeCachingPolicyCode(original["code"], d, config),
			"ttl":  flattenComputeRegionBackendServiceCdnPolicyNegativeCachingPolicyTtl(original["ttl"], d, config),
		})
	}
	return transformed
}
func flattenComputeRegionBackendServiceCdnPolicyNegativeCachingPolicyCode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicyNegativeCachingPolicyTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCdnPolicyCacheMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceCdnPolicyServeWhileStale(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceConnectionDraining(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["connection_draining_timeout_sec"] =
		flattenComputeRegionBackendServiceConnectionDrainingConnectionDrainingTimeoutSec(original["drainingTimeoutSec"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceConnectionDrainingConnectionDrainingTimeoutSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceCreationTimestamp(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceFailoverPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["disable_connection_drain_on_failover"] =
		flattenComputeRegionBackendServiceFailoverPolicyDisableConnectionDrainOnFailover(original["disableConnectionDrainOnFailover"], d, config)
	transformed["drop_traffic_if_unhealthy"] =
		flattenComputeRegionBackendServiceFailoverPolicyDropTrafficIfUnhealthy(original["dropTrafficIfUnhealthy"], d, config)
	transformed["failover_ratio"] =
		flattenComputeRegionBackendServiceFailoverPolicyFailoverRatio(original["failoverRatio"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceFailoverPolicyDisableConnectionDrainOnFailover(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceFailoverPolicyDropTrafficIfUnhealthy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceFailoverPolicyFailoverRatio(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceEnableCDN(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceFingerprint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceHealthChecks(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertAndMapStringArr(v.([]interface{}), tpgresource.ConvertSelfLinkToV1)
}

func flattenComputeRegionBackendServiceIap(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["oauth2_client_id"] =
		flattenComputeRegionBackendServiceIapOauth2ClientId(original["oauth2ClientId"], d, config)
	transformed["oauth2_client_secret"] =
		flattenComputeRegionBackendServiceIapOauth2ClientSecret(original["oauth2ClientSecret"], d, config)
	transformed["oauth2_client_secret_sha256"] =
		flattenComputeRegionBackendServiceIapOauth2ClientSecretSha256(original["oauth2ClientSecretSha256"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceIapOauth2ClientId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceIapOauth2ClientSecret(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return d.Get("iap.0.oauth2_client_secret")
}

func flattenComputeRegionBackendServiceIapOauth2ClientSecretSha256(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceLoadBalancingScheme(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceLocalityLbPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceOutlierDetection(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["base_ejection_time"] =
		flattenComputeRegionBackendServiceOutlierDetectionBaseEjectionTime(original["baseEjectionTime"], d, config)
	transformed["consecutive_errors"] =
		flattenComputeRegionBackendServiceOutlierDetectionConsecutiveErrors(original["consecutiveErrors"], d, config)
	transformed["consecutive_gateway_failure"] =
		flattenComputeRegionBackendServiceOutlierDetectionConsecutiveGatewayFailure(original["consecutiveGatewayFailure"], d, config)
	transformed["enforcing_consecutive_errors"] =
		flattenComputeRegionBackendServiceOutlierDetectionEnforcingConsecutiveErrors(original["enforcingConsecutiveErrors"], d, config)
	transformed["enforcing_consecutive_gateway_failure"] =
		flattenComputeRegionBackendServiceOutlierDetectionEnforcingConsecutiveGatewayFailure(original["enforcingConsecutiveGatewayFailure"], d, config)
	transformed["enforcing_success_rate"] =
		flattenComputeRegionBackendServiceOutlierDetectionEnforcingSuccessRate(original["enforcingSuccessRate"], d, config)
	transformed["interval"] =
		flattenComputeRegionBackendServiceOutlierDetectionInterval(original["interval"], d, config)
	transformed["max_ejection_percent"] =
		flattenComputeRegionBackendServiceOutlierDetectionMaxEjectionPercent(original["maxEjectionPercent"], d, config)
	transformed["success_rate_minimum_hosts"] =
		flattenComputeRegionBackendServiceOutlierDetectionSuccessRateMinimumHosts(original["successRateMinimumHosts"], d, config)
	transformed["success_rate_request_volume"] =
		flattenComputeRegionBackendServiceOutlierDetectionSuccessRateRequestVolume(original["successRateRequestVolume"], d, config)
	transformed["success_rate_stdev_factor"] =
		flattenComputeRegionBackendServiceOutlierDetectionSuccessRateStdevFactor(original["successRateStdevFactor"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceOutlierDetectionBaseEjectionTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeRegionBackendServiceOutlierDetectionBaseEjectionTimeSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeRegionBackendServiceOutlierDetectionBaseEjectionTimeNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceOutlierDetectionBaseEjectionTimeSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionBaseEjectionTimeNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionConsecutiveErrors(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionConsecutiveGatewayFailure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionEnforcingConsecutiveErrors(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionEnforcingConsecutiveGatewayFailure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionEnforcingSuccessRate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionInterval(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeRegionBackendServiceOutlierDetectionIntervalSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeRegionBackendServiceOutlierDetectionIntervalNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceOutlierDetectionIntervalSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionIntervalNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionMaxEjectionPercent(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionSuccessRateMinimumHosts(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionSuccessRateRequestVolume(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceOutlierDetectionSuccessRateStdevFactor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServicePortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceProtocol(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceSecurityPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceSessionAffinity(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceConnectionTrackingPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["idle_timeout_sec"] =
		flattenComputeRegionBackendServiceConnectionTrackingPolicyIdleTimeoutSec(original["idleTimeoutSec"], d, config)
	transformed["tracking_mode"] =
		flattenComputeRegionBackendServiceConnectionTrackingPolicyTrackingMode(original["trackingMode"], d, config)
	transformed["connection_persistence_on_unhealthy_backends"] =
		flattenComputeRegionBackendServiceConnectionTrackingPolicyConnectionPersistenceOnUnhealthyBackends(original["connectionPersistenceOnUnhealthyBackends"], d, config)
	transformed["enable_strong_affinity"] =
		flattenComputeRegionBackendServiceConnectionTrackingPolicyEnableStrongAffinity(original["enableStrongAffinity"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceConnectionTrackingPolicyIdleTimeoutSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceConnectionTrackingPolicyTrackingMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceConnectionTrackingPolicyConnectionPersistenceOnUnhealthyBackends(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceConnectionTrackingPolicyEnableStrongAffinity(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceTimeoutSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := tpgresource.StringToFixed64(strVal); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenComputeRegionBackendServiceLogConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["enable"] =
		flattenComputeRegionBackendServiceLogConfigEnable(original["enable"], d, config)
	transformed["sample_rate"] =
		flattenComputeRegionBackendServiceLogConfigSampleRate(original["sampleRate"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceLogConfigEnable(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceLogConfigSampleRate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceNetwork(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeRegionBackendServiceSubsetting(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["policy"] =
		flattenComputeRegionBackendServiceSubsettingPolicy(original["policy"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionBackendServiceSubsettingPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionBackendServiceRegion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}
