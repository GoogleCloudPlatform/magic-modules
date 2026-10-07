package cloudbuild

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudbuildWorkerPool() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceCloudbuildWorkerPool().Schema)

	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name", "location")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudbuildWorkerPoolRead,
		Schema: dsSchema,
	}

}

func dataSourceGoogleCloudbuildWorkerPoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/workerPools/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}

	id = strings.ReplaceAll(id, "/locations/global/", "/")
	d.SetId(id)

	err = resourceCloudbuildWorkerPoolRead(d, meta)
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
		Name:        "google_cloudbuild_worker_pool",
		ProductName: "cloudbuild",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleCloudbuildWorkerPool(),
	}.Register()
}
