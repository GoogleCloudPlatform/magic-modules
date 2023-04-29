package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// datasourceSchemaFromResourceSchema is a recursive func that
// converts an existing Resource schema to a Datasource schema.
// All schema elements are copied, but certain attributes are ignored or changed:
// - all attributes have Computed = true
// - all attributes have ForceNew, Required = false
// - Validation funcs and attributes (e.g. MaxItems) are not copied
func datasourceSchemaFromResourceSchema(rs map[string]*schema.Schema) map[string]*schema.Schema {
	return tpgresource.DatasourceSchemaFromResourceSchema(rs)
}

// fixDatasourceSchemaFlags is a convenience func that toggles the Computed,
// Optional + Required flags on a schema element. This is useful when the schema
// has been generated (using `datasourceSchemaFromResourceSchema` above for
// example) and therefore the attribute flags were not set appropriately when
// first added to the schema definition. Currently only supports top-level
// schema elements.
func fixDatasourceSchemaFlags(schema map[string]*schema.Schema, required bool, keys ...string) {
	tpgresource.FixDatasourceSchemaFlags(schema, required, keys...)
}

func addRequiredFieldsToSchema(schema map[string]*schema.Schema, keys ...string) {
	tpgresource.AddRequiredFieldsToSchema(schema, keys...)
}

func addOptionalFieldsToSchema(schema map[string]*schema.Schema, keys ...string) {
	tpgresource.AddOptionalFieldsToSchema(schema, keys...)
}
