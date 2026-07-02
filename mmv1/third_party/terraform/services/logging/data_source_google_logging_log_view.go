package logging

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleLoggingLogView() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceLoggingLogView().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "parent", "location", "bucket", "name")

	return &schema.Resource{
		Read:   dataSourceGoogleLoggingLogViewRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleLoggingLogViewRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "{{parent}}/locations/{{location}}/buckets/{{bucket}}/views/{{name}}")
	if err != nil {
		return err
	}

	d.SetId(id)

	err = resourceLoggingLogViewRead(d, meta)
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
		Name:        "google_logging_log_view",
		ProductName: "logging",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleLoggingLogView(),
	}.Register()
}
