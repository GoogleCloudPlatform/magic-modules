package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSqlDatabaseInstance() *schema.Resource {

	dsSchema := DatasourceSchemaFromResourceSchema(resourceSqlDatabaseInstance().Schema)
	AddRequiredFieldsToSchema(dsSchema, "name")
	AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceSqlDatabaseInstanceRead,
		Schema: dsSchema,
	}
}

func dataSourceSqlDatabaseInstanceRead(d *schema.ResourceData, meta interface{}) error {

	return resourceSqlDatabaseInstanceRead(d, meta)

}
