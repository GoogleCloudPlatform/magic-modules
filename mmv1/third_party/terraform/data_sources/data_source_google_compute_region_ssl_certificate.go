package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleRegionComputeSslCertificate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceComputeRegionSslCertificate().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")
	AddOptionalFieldsToSchema(dsSchema, "region")

	return &schema.Resource{
		Read:   dataSourceComputeRegionSslCertificateRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeRegionSslCertificateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, region, name, err := GetRegionalResourcePropertiesFromSelfLinkOrSchema(d, config)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("projects/%s/regions/%s/sslCertificates/%s", project, region, name))

	return resourceComputeRegionSslCertificateRead(d, meta)
}
