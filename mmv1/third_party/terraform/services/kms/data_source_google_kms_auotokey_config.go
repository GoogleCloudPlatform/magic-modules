package kms

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceGoogleKmsAutokeyConfig() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceKMSAutokeyConfig().Schema)
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "folder", "project")
	if folderSchema, ok := dsSchema["folder"]; ok {
		folderSchema.ExactlyOneOf = []string{"folder", "project"}
	}
	if projectSchema, ok := dsSchema["project"]; ok {
		projectSchema.ExactlyOneOf = []string{"folder", "project"}
	}

	return &schema.Resource{
		Read:   dataSourceGoogleKmsAutokeyConfigRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleKmsAutokeyConfigRead(d *schema.ResourceData, meta interface{}) error {
	folder, folderOk := d.GetOk("folder")
	project, projectOk := d.GetOk("project")
	if !folderOk && !projectOk {
		return fmt.Errorf("one of folder or project must be set")
	}

	folderVal := ""
	projectVal := ""
	if folderOk {
		folderVal = normalizeParent(folder.(string), "folder")
	}
	if projectOk {
		projectVal = normalizeParent(project.(string), "project")
	}
	parent := folderVal
	if parent == "" {
		parent = projectVal
	}
	id := parent + "/autokeyConfig"
	d.SetId(id)
	d.Set("name", id)
	err := resourceKMSAutokeyConfigRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
