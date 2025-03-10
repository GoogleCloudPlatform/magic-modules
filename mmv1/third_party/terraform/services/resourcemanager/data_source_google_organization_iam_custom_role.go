package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceGoogleOrganizationIamCustomRole() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceGoogleOrganizationIamCustomRole().Schema)

	dsSchema["org_id"].Computed = false
	dsSchema["org_id"].Required = true
	dsSchema["role_id"].Computed = false
	dsSchema["role_id"].Required = true

	return &schema.Resource{
		Read:   dataSourceOrganizationIamCustomRoleRead,
		Schema: dsSchema,
	}
}

func dataSourceOrganizationIamCustomRoleRead(d *schema.ResourceData, meta interface{}) error {
	orgId := d.Get("org_id").(string)
	roleId := d.Get("role_id").(string)
	d.SetId(fmt.Sprintf("organizations/%s/roles/%s", orgId, roleId))

	id := d.Id()

	if err := resourceGoogleOrganizationIamCustomRoleRead(d, meta); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("Role %s not found!", id)
	}

	return nil
}
