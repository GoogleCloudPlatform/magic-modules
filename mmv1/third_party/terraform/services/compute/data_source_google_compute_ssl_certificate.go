package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeSslCertificate() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeSslCertificate().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeSslCertificateRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeSslCertificateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	certificateName := d.Get("name").(string)

	id := fmt.Sprintf("projects/%s/global/sslCertificates/%s", project, certificateName)
	d.SetId(id)

	err = resourceComputeSslCertificateRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
