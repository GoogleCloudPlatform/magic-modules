package privilegedaccessmanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGooglePrivilegedAccessManagerEntitlement() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourcePrivilegedAccessManagerEntitlement().Schema)
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "entitlement_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "parent")

	return &schema.Resource{
		Read:   dataSourceGooglePrivilegedAccessManagerEntitlementRead,
		Schema: dsSchema,
	}
}

func dataSourceGooglePrivilegedAccessManagerEntitlementRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/entitlements/{{entitlement_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	err = resourcePrivilegedAccessManagerEntitlementRead(d, meta)
	if err != nil {
		return err
	}

	if err := tpgresource.SetDataSourceLabels(d); err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}
