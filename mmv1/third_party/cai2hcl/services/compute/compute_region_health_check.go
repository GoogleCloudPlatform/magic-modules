package compute

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/cai2hcl/common"
	"github.com/GoogleCloudPlatform/terraform-google-conversion/v5/caiasset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

const ComputeRegionHealthCheckAssetType string = "compute.googleapis.com/RegionHealthCheck"

const ComputeRegionHealthCheckAssetNameRegex string = "projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/healthChecks"

// ComputeRegionHealthCheckSchemaName is a TF resource schema name.
const ComputeRegionHealthCheckSchemaName string = "google_compute_region_health_check"

type ComputeRegionHealthCheckConverter struct {
	name   string
	schema map[string]*schema.Schema
}

// NewComputeRegionHealthCheckConverter returns an HCL converter for compute backend service.
func NewComputeRegionHealthCheckConverter(provider *schema.Provider) common.Converter {
	schema := provider.ResourcesMap[ComputeRegionHealthCheckSchemaName].Schema

	return &ComputeRegionHealthCheckConverter{
		name:   ComputeRegionHealthCheckSchemaName,
		schema: schema,
	}
}

func (c *ComputeRegionHealthCheckConverter) Convert(assets []*caiasset.Asset) ([]*common.HCLResourceBlock, error) {
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

func (c *ComputeRegionHealthCheckConverter) convertResourceData(asset *caiasset.Asset, config *transport_tpg.Config) (*common.HCLResourceBlock, error) {
	if asset == nil || asset.Resource == nil || asset.Resource.Data == nil {
		return nil, fmt.Errorf("asset resource data is nil")
	}

	assetResourceData := asset.Resource.Data

	hcl, _ := resourceComputeRegionHealthCheckRead(assetResourceData, config)

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

func resourceComputeRegionHealthCheckRead(resource map[string]interface{}, config *transport_tpg.Config) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var resource_data *schema.ResourceData = nil

	result["check_interval_sec"] = flattenComputeRegionHealthCheckCheckIntervalSec(resource["checkIntervalSec"], resource_data, config)
	result["creation_timestamp"] = flattenComputeRegionHealthCheckCreationTimestamp(resource["creationTimestamp"], resource_data, config)
	result["description"] = flattenComputeRegionHealthCheckDescription(resource["description"], resource_data, config)
	result["healthy_threshold"] = flattenComputeRegionHealthCheckHealthyThreshold(resource["healthyThreshold"], resource_data, config)
	result["name"] = flattenComputeRegionHealthCheckName(resource["name"], resource_data, config)
	result["unhealthy_threshold"] = flattenComputeRegionHealthCheckUnhealthyThreshold(resource["unhealthyThreshold"], resource_data, config)
	result["timeout_sec"] = flattenComputeRegionHealthCheckTimeoutSec(resource["timeoutSec"], resource_data, config)
	result["type"] = flattenComputeRegionHealthCheckType(resource["type"], resource_data, config)
	result["http_health_check"] = flattenComputeRegionHealthCheckHttpHealthCheck(resource["httpHealthCheck"], resource_data, config)
	result["https_health_check"] = flattenComputeRegionHealthCheckHttpsHealthCheck(resource["httpsHealthCheck"], resource_data, config)
	result["tcp_health_check"] = flattenComputeRegionHealthCheckTcpHealthCheck(resource["tcpHealthCheck"], resource_data, config)
	result["ssl_health_check"] = flattenComputeRegionHealthCheckSslHealthCheck(resource["sslHealthCheck"], resource_data, config)
	result["http2_health_check"] = flattenComputeRegionHealthCheckHttp2HealthCheck(resource["http2HealthCheck"], resource_data, config)
	result["grpc_health_check"] = flattenComputeRegionHealthCheckGrpcHealthCheck(resource["grpcHealthCheck"], resource_data, config)
	result["log_config"] = flattenComputeRegionHealthCheckLogConfig(resource["logConfig"], resource_data, config)
	result["region"] = flattenComputeRegionHealthCheckRegion(resource["region"], resource_data, config)

	return result, nil
}

func flattenComputeRegionHealthCheckCheckIntervalSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckCreationTimestamp(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckDescription(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHealthyThreshold(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckUnhealthyThreshold(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckTimeoutSec(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckType(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpHealthCheck(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["host"] =
		flattenComputeRegionHealthCheckHttpHealthCheckHost(original["host"], d, config)
	transformed["request_path"] =
		flattenComputeRegionHealthCheckHttpHealthCheckRequestPath(original["requestPath"], d, config)
	transformed["response"] =
		flattenComputeRegionHealthCheckHttpHealthCheckResponse(original["response"], d, config)
	transformed["port"] =
		flattenComputeRegionHealthCheckHttpHealthCheckPort(original["port"], d, config)
	transformed["port_name"] =
		flattenComputeRegionHealthCheckHttpHealthCheckPortName(original["portName"], d, config)
	transformed["proxy_header"] =
		flattenComputeRegionHealthCheckHttpHealthCheckProxyHeader(original["proxyHeader"], d, config)
	transformed["port_specification"] =
		flattenComputeRegionHealthCheckHttpHealthCheckPortSpecification(original["portSpecification"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionHealthCheckHttpHealthCheckHost(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpHealthCheckRequestPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpHealthCheckResponse(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpHealthCheckPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckHttpHealthCheckPortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpHealthCheckProxyHeader(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpHealthCheckPortSpecification(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpsHealthCheck(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["host"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckHost(original["host"], d, config)
	transformed["request_path"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckRequestPath(original["requestPath"], d, config)
	transformed["response"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckResponse(original["response"], d, config)
	transformed["port"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckPort(original["port"], d, config)
	transformed["port_name"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckPortName(original["portName"], d, config)
	transformed["proxy_header"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckProxyHeader(original["proxyHeader"], d, config)
	transformed["port_specification"] =
		flattenComputeRegionHealthCheckHttpsHealthCheckPortSpecification(original["portSpecification"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionHealthCheckHttpsHealthCheckHost(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpsHealthCheckRequestPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpsHealthCheckResponse(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpsHealthCheckPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckHttpsHealthCheckPortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpsHealthCheckProxyHeader(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttpsHealthCheckPortSpecification(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckTcpHealthCheck(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["request"] =
		flattenComputeRegionHealthCheckTcpHealthCheckRequest(original["request"], d, config)
	transformed["response"] =
		flattenComputeRegionHealthCheckTcpHealthCheckResponse(original["response"], d, config)
	transformed["port"] =
		flattenComputeRegionHealthCheckTcpHealthCheckPort(original["port"], d, config)
	transformed["port_name"] =
		flattenComputeRegionHealthCheckTcpHealthCheckPortName(original["portName"], d, config)
	transformed["proxy_header"] =
		flattenComputeRegionHealthCheckTcpHealthCheckProxyHeader(original["proxyHeader"], d, config)
	transformed["port_specification"] =
		flattenComputeRegionHealthCheckTcpHealthCheckPortSpecification(original["portSpecification"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionHealthCheckTcpHealthCheckRequest(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckTcpHealthCheckResponse(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckTcpHealthCheckPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckTcpHealthCheckPortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckTcpHealthCheckProxyHeader(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckTcpHealthCheckPortSpecification(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckSslHealthCheck(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["request"] =
		flattenComputeRegionHealthCheckSslHealthCheckRequest(original["request"], d, config)
	transformed["response"] =
		flattenComputeRegionHealthCheckSslHealthCheckResponse(original["response"], d, config)
	transformed["port"] =
		flattenComputeRegionHealthCheckSslHealthCheckPort(original["port"], d, config)
	transformed["port_name"] =
		flattenComputeRegionHealthCheckSslHealthCheckPortName(original["portName"], d, config)
	transformed["proxy_header"] =
		flattenComputeRegionHealthCheckSslHealthCheckProxyHeader(original["proxyHeader"], d, config)
	transformed["port_specification"] =
		flattenComputeRegionHealthCheckSslHealthCheckPortSpecification(original["portSpecification"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionHealthCheckSslHealthCheckRequest(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckSslHealthCheckResponse(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckSslHealthCheckPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckSslHealthCheckPortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckSslHealthCheckProxyHeader(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckSslHealthCheckPortSpecification(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttp2HealthCheck(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["host"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckHost(original["host"], d, config)
	transformed["request_path"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckRequestPath(original["requestPath"], d, config)
	transformed["response"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckResponse(original["response"], d, config)
	transformed["port"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckPort(original["port"], d, config)
	transformed["port_name"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckPortName(original["portName"], d, config)
	transformed["proxy_header"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckProxyHeader(original["proxyHeader"], d, config)
	transformed["port_specification"] =
		flattenComputeRegionHealthCheckHttp2HealthCheckPortSpecification(original["portSpecification"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionHealthCheckHttp2HealthCheckHost(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttp2HealthCheckRequestPath(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttp2HealthCheckResponse(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttp2HealthCheckPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckHttp2HealthCheckPortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttp2HealthCheckProxyHeader(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckHttp2HealthCheckPortSpecification(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckGrpcHealthCheck(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["port"] =
		flattenComputeRegionHealthCheckGrpcHealthCheckPort(original["port"], d, config)
	transformed["port_name"] =
		flattenComputeRegionHealthCheckGrpcHealthCheckPortName(original["portName"], d, config)
	transformed["port_specification"] =
		flattenComputeRegionHealthCheckGrpcHealthCheckPortSpecification(original["portSpecification"], d, config)
	transformed["grpc_service_name"] =
		flattenComputeRegionHealthCheckGrpcHealthCheckGrpcServiceName(original["grpcServiceName"], d, config)
	return []interface{}{transformed}
}
func flattenComputeRegionHealthCheckGrpcHealthCheckPort(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
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

func flattenComputeRegionHealthCheckGrpcHealthCheckPortName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckGrpcHealthCheckPortSpecification(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckGrpcHealthCheckGrpcServiceName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func flattenComputeRegionHealthCheckLogConfig(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	transformed := make(map[string]interface{})
	if v == nil {
		// Disabled by default, but API will not return object if value is false
		transformed["enable"] = false
		return []interface{}{transformed}
	}

	original := v.(map[string]interface{})
	transformed["enable"] = original["enable"]
	return []interface{}{transformed}
}

func flattenComputeRegionHealthCheckRegion(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	if v == nil {
		return v
	}
	return tpgresource.NameFromSelfLinkStateFunc(v)
}
