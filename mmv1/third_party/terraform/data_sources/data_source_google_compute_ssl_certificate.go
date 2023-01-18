package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleComputeSslCertificate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := DatasourceSchemaFromResourceSchema(ResourceComputeSslCertificate().Schema)

	// Set 'Required' schema elements
	AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeSslCertificateRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeSslCertificateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := GetProject(d, config)
	if err != nil {
		return err
	}
	certificateName := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/global/sslCertificates/%s", project, certificateName))

	return resourceComputeSslCertificateRead(d, meta)
}
