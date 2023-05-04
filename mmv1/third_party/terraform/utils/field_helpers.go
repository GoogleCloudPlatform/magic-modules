package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

const (
	globalLinkTemplate             = "projects/%s/global/%s/%s"
	globalLinkBasePattern          = "projects/(.+)/global/%s/(.+)"
	zonalLinkTemplate              = "projects/%s/zones/%s/%s/%s"
	zonalLinkBasePattern           = "projects/(.+)/zones/(.+)/%s/(.+)"
	zonalPartialLinkBasePattern    = "zones/(.+)/%s/(.+)"
	regionalLinkTemplate           = "projects/%s/regions/%s/%s/%s"
	regionalLinkBasePattern        = "projects/(.+)/regions/(.+)/%s/(.+)"
	regionalPartialLinkBasePattern = "regions/(.+)/%s/(.+)"
	projectLinkTemplate            = "projects/%s/%s/%s"
	projectBasePattern             = "projects/(.+)/%s/(.+)"
	organizationLinkTemplate       = "organizations/%s/%s/%s"
	organizationBasePattern        = "organizations/(.+)/%s/(.+)"
)

// ------------------------------------------------------------
// Field helpers
// ------------------------------------------------------------

func ParseNetworkFieldValue(network string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseNetworkFieldValue(network, d, config)
}

func ParseSubnetworkFieldValue(subnetwork string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseSubnetworkFieldValue(subnetwork, d, config)
}

func ParseSubnetworkFieldValueWithProjectField(subnetwork, projectField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseSubnetworkFieldValueWithProjectField(subnetwork, projectField, d, config)
}

func ParseSslCertificateFieldValue(sslCertificate string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseSslCertificateFieldValue(sslCertificate, d, config)
}

func ParseHttpHealthCheckFieldValue(healthCheck string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseHttpHealthCheckFieldValue(healthCheck, d, config)
}

func ParseDiskFieldValue(disk string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseDiskFieldValue(disk, d, config)
}

func ParseRegionDiskFieldValue(disk string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseRegionDiskFieldValue(disk, d, config)
}

func ParseOrganizationCustomRoleName(role string) (*tpgresource.OrganizationFieldValue, error) {
	return tpgresource.ParseOrganizationCustomRoleName(role)
}

func ParseAcceleratorFieldValue(accelerator string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseAcceleratorFieldValue(accelerator, d, config)
}

func ParseMachineTypesFieldValue(machineType string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseMachineTypesFieldValue(machineType, d, config)
}

func ParseInstanceFieldValue(instance string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseInstanceFieldValue(instance, d, config)
}

func ParseInstanceGroupFieldValue(instanceGroup string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseInstanceGroupFieldValue(instanceGroup, d, config)
}

func ParseInstanceTemplateFieldValue(instanceTemplate string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseInstanceTemplateFieldValue(instanceTemplate, d, config)
}

func ParseMachineImageFieldValue(machineImage string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseMachineImageFieldValue(machineImage, d, config)
}

func ParseSecurityPolicyFieldValue(securityPolicy string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseSecurityPolicyFieldValue(securityPolicy, d, config)
}

func ParseNetworkEndpointGroupFieldValue(networkEndpointGroup string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseNetworkEndpointGroupFieldValue(networkEndpointGroup, d, config)
}

func ParseNetworkEndpointGroupRegionalFieldValue(networkEndpointGroup string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseNetworkEndpointGroupRegionalFieldValue(networkEndpointGroup, d, config)
}

// ------------------------------------------------------------
// Base helpers used to create helpers for specific fields.
// ------------------------------------------------------------

// Parses a global field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/global/{resource_type}/{resource_name}
// - projects/{my_project}/global/{resource_type}/{resource_name}
// - global/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
//
// Deprecated: For backward compatibility parseGlobalFieldValue is still working,
// but all new code should use ParseGlobalFieldValue in the tpgresource package instead.
func parseGlobalFieldValue(resourceType, fieldValue, projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.GlobalFieldValue, error) {
	return tpgresource.ParseGlobalFieldValue(resourceType, fieldValue, projectSchemaField, d, config, isEmptyValid)
}

// Parses a zonal field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/zones/{zone}/{resource_type}/{resource_name}
// - projects/{my_project}/zones/{zone}/{resource_type}/{resource_name}
// - zones/{zone}/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
// If the zone is not specified, it takes the value of `zoneSchemaField`.
//
// Deprecated: For backward compatibility parseZonalFieldValue is still working,
// but all new code should use ParseZonalFieldValue in the tpgresource package instead.
func parseZonalFieldValue(resourceType, fieldValue, projectSchemaField, zoneSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.ZonalFieldValue, error) {
	return tpgresource.ParseZonalFieldValue(resourceType, fieldValue, projectSchemaField, zoneSchemaField, d, config, isEmptyValid)
}

// Deprecated: For backward compatibility getProjectFromSchema is still working,
// but all new code should use GetProjectFromSchema in the tpgresource package instead.
func getProjectFromSchema(projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetProjectFromSchema(projectSchemaField, d, config)
}

// Deprecated: For backward compatibility getBillingProjectFromSchema is still working,
// but all new code should use GetBillingProjectFromSchema in the tpgresource package instead.
func getBillingProjectFromSchema(billingProjectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetBillingProjectFromSchema(billingProjectSchemaField, d, config)
}

// Parses an organization field with the following formats:
// - organizations/{my_organizations}/{resource_type}/{resource_name}
//
// Deprecated: For backward compatibility parseOrganizationFieldValue is still working,
// but all new code should use ParseOrganizationFieldValue in the tpgresource package instead.
func parseOrganizationFieldValue(resourceType, fieldValue string, isEmptyValid bool) (*tpgresource.OrganizationFieldValue, error) {
	return tpgresource.ParseOrganizationFieldValue(resourceType, fieldValue, isEmptyValid)
}

// Parses a regional field supporting 5 different formats:
// - https://www.googleapis.com/compute/ANY_VERSION/projects/{my_project}/regions/{region}/{resource_type}/{resource_name}
// - projects/{my_project}/regions/{region}/{resource_type}/{resource_name}
// - regions/{region}/{resource_type}/{resource_name}
// - resource_name
// - "" (empty string). RelativeLink() returns empty if isEmptyValid is true.
//
// If the project is not specified, it first tries to get the project from the `projectSchemaField` and then fallback on the default project.
// If the region is not specified, see function documentation for `getRegionFromSchema`.
//
// Deprecated: For backward compatibility parseRegionalFieldValue is still working,
// but all new code should use ParseRegionalFieldValue in the tpgresource package instead.
func parseRegionalFieldValue(resourceType, fieldValue, projectSchemaField, regionSchemaField, zoneSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.RegionalFieldValue, error) {
	return tpgresource.ParseRegionalFieldValue(resourceType, fieldValue, projectSchemaField, regionSchemaField, zoneSchemaField, d, config, isEmptyValid)
}

// Infers the region based on the following (in order of priority):
// - `regionSchemaField` in resource schema
// - region extracted from the `zoneSchemaField` in resource schema
// - provider-level region
// - region extracted from the provider-level zone
//
// Deprecated: For backward compatibility getRegionFromSchema is still working,
// but all new code should use GetRegionFromSchema in the tpgresource package instead.
func getRegionFromSchema(regionSchemaField, zoneSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (string, error) {
	return tpgresource.GetRegionFromSchema(regionSchemaField, zoneSchemaField, d, config)
}

// Parses a project field with the following formats:
// - projects/{my_projects}/{resource_type}/{resource_name}
//
// Deprecated: For backward compatibility parseProjectFieldValue is still working,
// but all new code should use ParseProjectFieldValue in the tpgresource package instead.
func parseProjectFieldValue(resourceType, fieldValue, projectSchemaField string, d tpgresource.TerraformResourceData, config *transport_tpg.Config, isEmptyValid bool) (*tpgresource.ProjectFieldValue, error) {
	return tpgresource.ParseProjectFieldValue(resourceType, fieldValue, projectSchemaField, d, config, isEmptyValid)
}
