package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeNetworkAttachment() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeNetworkAttachment().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name", "region")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceComputeNetworkAttachmentRead,
		Schema: dsSchema,
	}
}

func dataSourceComputeNetworkAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}

	name := d.Get("name").(string)
	region := d.Get("region").(string)

	id := fmt.Sprintf("projects/%s/regions/%s/networkAttachments/%s", project, region, name)
	d.SetId(id)

	err = resourceComputeNetworkAttachmentRead(d, meta)
	if err != nil {
		return fmt.Errorf("Error reading Network Attachment %q: %s", id, err)
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}

	return nil
}
