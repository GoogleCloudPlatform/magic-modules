package compute

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"regexp"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
	"google.golang.org/api/googleapi"
)

// suppress changes on sample_rate if log_config is set to disabled.
func suppressWhenDisabled(k, old, new string, d *schema.ResourceData) bool {
	_, n := d.GetChange("log_config.0.enable")
	if tpgresource.IsEmptyValue(reflect.ValueOf(n)) {
		return true
	}
	return false
}

// Whether the backend is a global or regional NEG
func isNegBackend(backend map[string]interface{}) bool {
	backendGroup, ok := backend["group"]
	if !ok {
		return false
	}

	match, err := regexp.MatchString("(?:global|regions/[^/]+)/networkEndpointGroups", backendGroup.(string))
	if err != nil {
		// should not happen as long as the regexp pattern compiled correctly
		return false
	}
	return match
}

func resourceGoogleComputeBackendServiceBackendHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})
	log.Printf("[DEBUG] hashing %v", m)

	if group, err := tpgresource.GetRelativePath(m["group"].(string)); err != nil {
		log.Printf("[WARN] Error on retrieving relative path of instance group: %s", err)
		buf.WriteString(fmt.Sprintf("%s-", m["group"].(string)))
	} else {
		buf.WriteString(fmt.Sprintf("%s-", group))
	}

	if v, ok := m["balancing_mode"]; ok {
		if v == nil {
			v = ""
		}

		buf.WriteString(fmt.Sprintf("%v-", v))
	}
	if v, ok := m["capacity_scaler"]; ok {
		if v == nil {
			v = 0.0
		}

		// floats can't be added to the hash with %v as the other values are because
		// %v and %f are not equivalent strings so this must remain as a float so that
		// the hash function doesn't return something else.
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["description"]; ok {
		if v == nil {
			v = ""
		}

		log.Printf("[DEBUG] writing description %s", v)
		buf.WriteString(fmt.Sprintf("%v-", v))
	}
	if v, ok := m["max_rate"]; ok {
		if v == nil {
			v = 0
		}

		buf.WriteString(fmt.Sprintf("%v-", v))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		if v == nil {
			v = 0.0
		}

		// floats can't be added to the hash with %v as the other values are because
		// %v and %f are not equivalent strings so this must remain as a float so that
		// the hash function doesn't return something else.
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["max_connections"]; ok {
		if v == nil {
			v = 0
		}

		buf.WriteString(fmt.Sprintf("%v-", v))
	}
	if v, ok := m["max_connections_per_instance"]; ok {
		if v == nil {
			v = 0
		}

		buf.WriteString(fmt.Sprintf("%v-", v))
	}
	if v, ok := m["max_rate_per_instance"]; ok {
		if v == nil {
			v = 0.0
		}

		// floats can't be added to the hash with %v as the other values are because
		// %v and %f are not equivalent strings so this must remain as a float so that
		// the hash function doesn't return something else.
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["max_connections_per_endpoint"]; ok {
		if v == nil {
			v = 0
		}

		buf.WriteString(fmt.Sprintf("%v-", v))
	}
	if v, ok := m["max_rate_per_endpoint"]; ok {
		if v == nil {
			v = 0.0
		}

		// floats can't be added to the hash with %v as the other values are because
		// %v and %f are not equivalent strings so this must remain as a float so that
		// the hash function doesn't return something else.
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}
	if v, ok := m["max_utilization"]; ok && !isNegBackend(m) {
		if v == nil {
			v = 0.0
		}

		// floats can't be added to the hash with %v as the other values are because
		// %v and %f are not equivalent strings so this must remain as a float so that
		// the hash function doesn't return something else.
		buf.WriteString(fmt.Sprintf("%f-", v.(float64)))
	}

	// This is in region backend service, but not in backend service.  Should be a no-op
	// if it's not present.
	if v, ok := m["failover"]; ok {
		if v == nil {
			v = false
		}
		buf.WriteString(fmt.Sprintf("%v-", v.(bool)))
	}

	log.Printf("[DEBUG] computed hash value of %v from %v", tpgresource.Hashcode(buf.String()), buf.String())
	return tpgresource.Hashcode(buf.String())
}

// ComputeBackendServiceAssetType is the CAI asset type name.
const ComputeBackendServiceAssetType string = "compute.googleapis.com/BackendService"

// ComputeBackendServiceSchemaName is a TF resource schema name.
const ComputeBackendServiceSchemaName string = "google_compute_backend_service"

type ComputeBackendServiceConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewComputeBackendServiceConverter returns an HCL converter for compute backend service.
func NewComputeBackendServiceConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[ComputeBackendServiceSchemaName].Schema

	return &ComputeBackendServiceConverter{
		name:   ComputeBackendServiceSchemaName,
		schema: schema,
	}
}

func (c *ComputeBackendServiceConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
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

func (c *ComputeBackendServiceConverter) convertResourceData(asset *caiasset.Asset, config *transport_tpg.Config) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	assetResourceData := asset.Resource.Data

	hcl, _ := flattenComputeBackendService(assetResourceData, config)

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

func flattenComputeBackendService(resource map[string]interface{}, config *transport_tpg.Config) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var resource_data *schema.ResourceData = nil

	result["affinity_cookie_ttl_sec"] = flattenComputeBackendServiceAffinityCookieTtlSec(resource["affinityCookieTtlSec"], resource_data, config)
	result["backend"] = flattenComputeBackendServiceBackend(resource["backends"], resource_data, config)
	result["circuit_breakers"] = flattenComputeBackendServiceCircuitBreakers(resource["circuitBreakers"], resource_data, config)
	result["compression_mode"] = flattenComputeBackendServiceCompressionMode(resource["compressionMode"], resource_data, config)
	result["consistent_hash"] = flattenComputeBackendServiceConsistentHash(resource["consistentHash"], resource_data, config)
	result["cdn_policy"] = flattenComputeBackendServiceCdnPolicy(resource["cdnPolicy"], resource_data, config)
	if flattenedProp := flattenComputeBackendServiceConnectionDraining(resource["connectionDraining"], resource_data, config); flattenedProp != nil {
		if gerr, ok := flattenedProp.(*googleapi.Error); ok {
			return nil, fmt.Errorf("Error reading BackendService: %s", gerr)
		}
		casted := flattenedProp.([]interface{})[0]
		if casted != nil {
			for k, v := range casted.(map[string]interface{}) {
				result[k] = v
			}
		}
	}
	result["creation_timestamp"] = flattenComputeBackendServiceCreationTimestamp(resource["creationTimestamp"], resource_data, config)
	result["custom_request_headers"] = flattenComputeBackendServiceCustomRequestHeaders(resource["customRequestHeaders"], resource_data, config)
	result["custom_response_headers"] = flattenComputeBackendServiceCustomResponseHeaders(resource["customResponseHeaders"], resource_data, config)
	result["fingerprint"] = flattenComputeBackendServiceFingerprint(resource["fingerprint"], resource_data, config)
	result["description"] = flattenComputeBackendServiceDescription(resource["description"], resource_data, config)
	result["enable_cdn"] = flattenComputeBackendServiceEnableCDN(resource["enableCDN"], resource_data, config)
	result["health_checks"] = flattenComputeBackendServiceHealthChecks(resource["healthChecks"], resource_data, config)
	result["generated_id"] = flattenComputeBackendServiceGeneratedId(resource["id"], resource_data, config)
	result["iap"] = flattenComputeBackendServiceIap(resource["iap"], resource_data, config)
	result["load_balancing_scheme"] = flattenComputeBackendServiceLoadBalancingScheme(resource["loadBalancingScheme"], resource_data, config)
	result["locality_lb_policy"] = flattenComputeBackendServiceLocalityLbPolicy(resource["localityLbPolicy"], resource_data, config)
	result["locality_lb_policies"] = flattenComputeBackendServiceLocalityLbPolicies(resource["localityLbPolicies"], resource_data, config)
	result["name"] = flattenComputeBackendServiceName(resource["name"], resource_data, config)
	result["outlier_detection"] = flattenComputeBackendServiceOutlierDetection(resource["outlierDetection"], resource_data, config)
	result["port_name"] = flattenComputeBackendServicePortName(resource["portName"], resource_data, config)
	result["protocol"] = flattenComputeBackendServiceProtocol(resource["protocol"], resource_data, config)
	result["security_policy"] = flattenComputeBackendServiceSecurityPolicy(resource["securityPolicy"], resource_data, config)
	result["edge_security_policy"] = flattenComputeBackendServiceEdgeSecurityPolicy(resource["edgeSecurityPolicy"], resource_data, config)
	result["security_settings"] = flattenComputeBackendServiceSecuritySettings(resource["securitySettings"], resource_data, config)
	result["session_affinity"] = flattenComputeBackendServiceSessionAffinity(resource["sessionAffinity"], resource_data, config)
	result["timeout_sec"] = flattenComputeBackendServiceTimeoutSec(resource["timeoutSec"], resource_data, config)
	result["log_config"] = flattenComputeBackendServiceLogConfig(resource["logConfig"], resource_data, config)

	return result, nil
}

func flattenComputeBackendServiceAffinityCookieTtlSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceBackend(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"balancing_mode":               flattenComputeBackendServiceBackendBalancingMode(original["balancingMode"], d, config),
			"capacity_scaler":              flattenComputeBackendServiceBackendCapacityScaler(original["capacityScaler"], d, config),
			"description":                  flattenComputeBackendServiceBackendDescription(original["description"], d, config),
			"group":                        flattenComputeBackendServiceBackendGroup(original["group"], d, config),
			"max_connections":              flattenComputeBackendServiceBackendMaxConnections(original["maxConnections"], d, config),
			"max_connections_per_instance": flattenComputeBackendServiceBackendMaxConnectionsPerInstance(original["maxConnectionsPerInstance"], d, config),
			"max_connections_per_endpoint": flattenComputeBackendServiceBackendMaxConnectionsPerEndpoint(original["maxConnectionsPerEndpoint"], d, config),
			"max_rate":                     flattenComputeBackendServiceBackendMaxRate(original["maxRate"], d, config),
			"max_rate_per_instance":        flattenComputeBackendServiceBackendMaxRatePerInstance(original["maxRatePerInstance"], d, config),
			"max_rate_per_endpoint":        flattenComputeBackendServiceBackendMaxRatePerEndpoint(original["maxRatePerEndpoint"], d, config),
			"max_utilization":              flattenComputeBackendServiceBackendMaxUtilization(original["maxUtilization"], d, config),
		})
	}
	return transformed
}
func flattenComputeBackendServiceBackendBalancingMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceBackendCapacityScaler(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceBackendDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceBackendGroup(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeBackendServiceBackendMaxConnections(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceBackendMaxConnectionsPerInstance(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceBackendMaxConnectionsPerEndpoint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceBackendMaxRate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceBackendMaxRatePerInstance(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceBackendMaxRatePerEndpoint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceBackendMaxUtilization(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCircuitBreakers(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["connect_timeout"] =
		flattenComputeBackendServiceCircuitBreakersConnectTimeout(original["connectTimeout"], d, config)
	transformed["max_requests_per_connection"] =
		flattenComputeBackendServiceCircuitBreakersMaxRequestsPerConnection(original["maxRequestsPerConnection"], d, config)
	transformed["max_connections"] =
		flattenComputeBackendServiceCircuitBreakersMaxConnections(original["maxConnections"], d, config)
	transformed["max_pending_requests"] =
		flattenComputeBackendServiceCircuitBreakersMaxPendingRequests(original["maxPendingRequests"], d, config)
	transformed["max_requests"] =
		flattenComputeBackendServiceCircuitBreakersMaxRequests(original["maxRequests"], d, config)
	transformed["max_retries"] =
		flattenComputeBackendServiceCircuitBreakersMaxRetries(original["maxRetries"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceCircuitBreakersConnectTimeout(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeBackendServiceCircuitBreakersConnectTimeoutSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeBackendServiceCircuitBreakersConnectTimeoutNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceCircuitBreakersConnectTimeoutSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCircuitBreakersConnectTimeoutNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCircuitBreakersMaxRequestsPerConnection(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCircuitBreakersMaxConnections(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCircuitBreakersMaxPendingRequests(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCircuitBreakersMaxRequests(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCircuitBreakersMaxRetries(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCompressionMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceConsistentHash(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["http_cookie"] =
		flattenComputeBackendServiceConsistentHashHttpCookie(original["httpCookie"], d, config)
	transformed["http_header_name"] =
		flattenComputeBackendServiceConsistentHashHttpHeaderName(original["httpHeaderName"], d, config)
	transformed["minimum_ring_size"] =
		flattenComputeBackendServiceConsistentHashMinimumRingSize(original["minimumRingSize"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceConsistentHashHttpCookie(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["ttl"] =
		flattenComputeBackendServiceConsistentHashHttpCookieTtl(original["ttl"], d, config)
	transformed["name"] =
		flattenComputeBackendServiceConsistentHashHttpCookieName(original["name"], d, config)
	transformed["path"] =
		flattenComputeBackendServiceConsistentHashHttpCookiePath(original["path"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceConsistentHashHttpCookieTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeBackendServiceConsistentHashHttpCookieTtlSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeBackendServiceConsistentHashHttpCookieTtlNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceConsistentHashHttpCookieTtlSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceConsistentHashHttpCookieTtlNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceConsistentHashHttpCookieName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceConsistentHashHttpCookiePath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceConsistentHashHttpHeaderName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceConsistentHashMinimumRingSize(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["cache_key_policy"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicy(original["cacheKeyPolicy"], d, config)
	transformed["signed_url_cache_max_age_sec"] =
		flattenComputeBackendServiceCdnPolicySignedUrlCacheMaxAgeSec(original["signedUrlCacheMaxAgeSec"], d, config)
	transformed["default_ttl"] =
		flattenComputeBackendServiceCdnPolicyDefaultTtl(original["defaultTtl"], d, config)
	transformed["max_ttl"] =
		flattenComputeBackendServiceCdnPolicyMaxTtl(original["maxTtl"], d, config)
	transformed["client_ttl"] =
		flattenComputeBackendServiceCdnPolicyClientTtl(original["clientTtl"], d, config)
	transformed["negative_caching"] =
		flattenComputeBackendServiceCdnPolicyNegativeCaching(original["negativeCaching"], d, config)
	transformed["negative_caching_policy"] =
		flattenComputeBackendServiceCdnPolicyNegativeCachingPolicy(original["negativeCachingPolicy"], d, config)
	transformed["cache_mode"] =
		flattenComputeBackendServiceCdnPolicyCacheMode(original["cacheMode"], d, config)
	transformed["serve_while_stale"] =
		flattenComputeBackendServiceCdnPolicyServeWhileStale(original["serveWhileStale"], d, config)
	transformed["bypass_cache_on_request_headers"] =
		flattenComputeBackendServiceCdnPolicyBypassCacheOnRequestHeaders(original["bypassCacheOnRequestHeaders"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceCdnPolicyCacheKeyPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["include_host"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeHost(original["includeHost"], d, config)
	transformed["include_protocol"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeProtocol(original["includeProtocol"], d, config)
	transformed["include_query_string"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeQueryString(original["includeQueryString"], d, config)
	transformed["query_string_blacklist"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyQueryStringBlacklist(original["queryStringBlacklist"], d, config)
	transformed["query_string_whitelist"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyQueryStringWhitelist(original["queryStringWhitelist"], d, config)
	transformed["include_http_headers"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeHttpHeaders(original["includeHttpHeaders"], d, config)
	transformed["include_named_cookies"] =
		flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeNamedCookies(original["includeNamedCookies"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeHost(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeProtocol(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeQueryString(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyQueryStringBlacklist(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyQueryStringWhitelist(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeHttpHeaders(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicyCacheKeyPolicyIncludeNamedCookies(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicySignedUrlCacheMaxAgeSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyDefaultTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyMaxTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyClientTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyNegativeCaching(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicyNegativeCachingPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"code": flattenComputeBackendServiceCdnPolicyNegativeCachingPolicyCode(original["code"], d, config),
			"ttl":  flattenComputeBackendServiceCdnPolicyNegativeCachingPolicyTtl(original["ttl"], d, config),
		})
	}
	return transformed
}
func flattenComputeBackendServiceCdnPolicyNegativeCachingPolicyCode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyNegativeCachingPolicyTtl(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyCacheMode(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCdnPolicyServeWhileStale(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCdnPolicyBypassCacheOnRequestHeaders(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"header_name": flattenComputeBackendServiceCdnPolicyBypassCacheOnRequestHeadersHeaderName(original["headerName"], d, config),
		})
	}
	return transformed
}
func flattenComputeBackendServiceCdnPolicyBypassCacheOnRequestHeadersHeaderName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceConnectionDraining(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["connection_draining_timeout_sec"] =
		flattenComputeBackendServiceConnectionDrainingConnectionDrainingTimeoutSec(original["drainingTimeoutSec"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceConnectionDrainingConnectionDrainingTimeoutSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceCreationTimestamp(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceCustomRequestHeaders(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeBackendServiceCustomResponseHeaders(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return schema.NewSet(schema.HashString, v.([]interface{}))
}

func flattenComputeBackendServiceFingerprint(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceEnableCDN(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceHealthChecks(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertAndMapStringArr(v.([]interface{}), tpgresource.ConvertSelfLinkToV1)
}

func flattenComputeBackendServiceGeneratedId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceIap(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["oauth2_client_id"] =
		flattenComputeBackendServiceIapOauth2ClientId(original["oauth2ClientId"], d, config)
	transformed["oauth2_client_secret"] =
		flattenComputeBackendServiceIapOauth2ClientSecret(original["oauth2ClientSecret"], d, config)
	transformed["oauth2_client_secret_sha256"] =
		flattenComputeBackendServiceIapOauth2ClientSecretSha256(original["oauth2ClientSecretSha256"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceIapOauth2ClientId(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceIapOauth2ClientSecret(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return d.Get("iap.0.oauth2_client_secret")
}

func flattenComputeBackendServiceIapOauth2ClientSecretSha256(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceLoadBalancingScheme(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceLocalityLbPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceLocalityLbPolicies(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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
			"policy":        flattenComputeBackendServiceLocalityLbPoliciesPolicy(original["policy"], d, config),
			"custom_policy": flattenComputeBackendServiceLocalityLbPoliciesCustomPolicy(original["customPolicy"], d, config),
		})
	}
	return transformed
}
func flattenComputeBackendServiceLocalityLbPoliciesPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenComputeBackendServiceLocalityLbPoliciesPolicyName(original["name"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceLocalityLbPoliciesPolicyName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceLocalityLbPoliciesCustomPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["name"] =
		flattenComputeBackendServiceLocalityLbPoliciesCustomPolicyName(original["name"], d, config)
	transformed["data"] =
		flattenComputeBackendServiceLocalityLbPoliciesCustomPolicyData(original["data"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceLocalityLbPoliciesCustomPolicyName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceLocalityLbPoliciesCustomPolicyData(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceOutlierDetection(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["base_ejection_time"] =
		flattenComputeBackendServiceOutlierDetectionBaseEjectionTime(original["baseEjectionTime"], d, config)
	transformed["consecutive_errors"] =
		flattenComputeBackendServiceOutlierDetectionConsecutiveErrors(original["consecutiveErrors"], d, config)
	transformed["consecutive_gateway_failure"] =
		flattenComputeBackendServiceOutlierDetectionConsecutiveGatewayFailure(original["consecutiveGatewayFailure"], d, config)
	transformed["enforcing_consecutive_errors"] =
		flattenComputeBackendServiceOutlierDetectionEnforcingConsecutiveErrors(original["enforcingConsecutiveErrors"], d, config)
	transformed["enforcing_consecutive_gateway_failure"] =
		flattenComputeBackendServiceOutlierDetectionEnforcingConsecutiveGatewayFailure(original["enforcingConsecutiveGatewayFailure"], d, config)
	transformed["enforcing_success_rate"] =
		flattenComputeBackendServiceOutlierDetectionEnforcingSuccessRate(original["enforcingSuccessRate"], d, config)
	transformed["interval"] =
		flattenComputeBackendServiceOutlierDetectionInterval(original["interval"], d, config)
	transformed["max_ejection_percent"] =
		flattenComputeBackendServiceOutlierDetectionMaxEjectionPercent(original["maxEjectionPercent"], d, config)
	transformed["success_rate_minimum_hosts"] =
		flattenComputeBackendServiceOutlierDetectionSuccessRateMinimumHosts(original["successRateMinimumHosts"], d, config)
	transformed["success_rate_request_volume"] =
		flattenComputeBackendServiceOutlierDetectionSuccessRateRequestVolume(original["successRateRequestVolume"], d, config)
	transformed["success_rate_stdev_factor"] =
		flattenComputeBackendServiceOutlierDetectionSuccessRateStdevFactor(original["successRateStdevFactor"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceOutlierDetectionBaseEjectionTime(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeBackendServiceOutlierDetectionBaseEjectionTimeSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeBackendServiceOutlierDetectionBaseEjectionTimeNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceOutlierDetectionBaseEjectionTimeSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionBaseEjectionTimeNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionConsecutiveErrors(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionConsecutiveGatewayFailure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionEnforcingConsecutiveErrors(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionEnforcingConsecutiveGatewayFailure(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionEnforcingSuccessRate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionInterval(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["seconds"] =
		flattenComputeBackendServiceOutlierDetectionIntervalSeconds(original["seconds"], d, config)
	transformed["nanos"] =
		flattenComputeBackendServiceOutlierDetectionIntervalNanos(original["nanos"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceOutlierDetectionIntervalSeconds(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionIntervalNanos(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionMaxEjectionPercent(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionSuccessRateMinimumHosts(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionSuccessRateRequestVolume(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceOutlierDetectionSuccessRateStdevFactor(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServicePortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceProtocol(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceSecurityPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceEdgeSecurityPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceSecuritySettings(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["client_tls_policy"] =
		flattenComputeBackendServiceSecuritySettingsClientTlsPolicy(original["clientTlsPolicy"], d, config)
	transformed["subject_alt_names"] =
		flattenComputeBackendServiceSecuritySettingsSubjectAltNames(original["subjectAltNames"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceSecuritySettingsClientTlsPolicy(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.ConvertSelfLinkToV1(v.(string))
}

func flattenComputeBackendServiceSecuritySettingsSubjectAltNames(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceSessionAffinity(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceTimeoutSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeBackendServiceLogConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["enable"] =
		flattenComputeBackendServiceLogConfigEnable(original["enable"], d, config)
	transformed["sample_rate"] =
		flattenComputeBackendServiceLogConfigSampleRate(original["sampleRate"], d, config)
	return []interface{}{transformed}
}
func flattenComputeBackendServiceLogConfigEnable(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeBackendServiceLogConfigSampleRate(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}
