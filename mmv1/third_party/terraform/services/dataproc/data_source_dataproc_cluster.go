package dataproc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceArtifactRegistryRepository() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceDataprocCluster().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name", "region")

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceDataprocClusterRead,
		Schema: dsSchema,
	}
}

func dataSourceDataprocClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	d.SetId(fmt.Sprintf("projects/%s/regions/%s/clusters/%s", project, region, name))

	err = resourceDataprocClusterRead(d, meta)
	if err != nil {
		return err
	}

	return nil
}
