package backupdr

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudBackupDRBackupPlanAssociation() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBackupDRBackupPlanAssociation().Schema)
	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "backup_plan_association_id", "location")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	return &schema.Resource{
		Read:   dataSourceGoogleCloudBackupDRBackupPlanAssociationRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCloudBackupDRBackupPlanAssociationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	backup_plan_association_id := d.Get("backup_plan_association_id").(string)
	id := fmt.Sprintf("projects/%s/locations/%s/backupPlanAssociations/%s", project, location, backup_plan_association_id)
	d.SetId(id)
	err = resourceBackupDRBackupPlanAssociationRead(d, meta)
	if err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
