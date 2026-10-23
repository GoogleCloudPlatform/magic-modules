package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeServiceAttachment() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeServiceAttachment().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "region")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeServiceAttachmentRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeServiceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	name := d.Get("name").(string)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/regions/%s/serviceAttachments/%s", project, region, name)
	d.SetId(id)

	err = resourceComputeServiceAttachmentRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_service_attachment",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeServiceAttachment(),
	}.Register()
}
