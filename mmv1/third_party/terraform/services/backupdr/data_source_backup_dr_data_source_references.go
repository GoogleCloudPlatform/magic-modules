package backupdr

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRDataSourceReference() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBackupDRDataSourceReferences().Schema)
	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "data_source_reference_id", "location")
	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	return &schema.Resource{
		Read:   dataSourceGoogleCloudBackupDRDataSourceReferenceRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudBackupDRDataSourceReferenceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	data_source_reference_id := d.Get("data_source_reference_id").(string)
	id := fmt.Sprintf("projects/%s/locations/%s/dataSourceReferences/%s", project, location, data_source_reference_id)
	d.SetId(id)
	err = resourceBackupDRDataSourceReferencesRead(d, meta)
	if err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
