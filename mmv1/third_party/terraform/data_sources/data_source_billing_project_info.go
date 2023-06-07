package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCoreBillingProjectInfo() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCoreBillingProjectInfo().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCoreBillingProjectInfoRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleCoreBillingProjectInfoRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/billingInfo")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)
	return resourceCoreBillingProjectInfoRead(d, meta)
}
