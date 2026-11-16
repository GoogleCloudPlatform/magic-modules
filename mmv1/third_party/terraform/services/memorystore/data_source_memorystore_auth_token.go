package memorystore

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceMemorystoreAuthToken() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMemorystoreAuthToken().Schema)

	// Set 'Required' and 'Optional' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "instance", "user_id", "token_id")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project", "location")

	return &schema.Resource{
		Read:   dataSourceMemorystoreAuthTokenRead,
		Schema: dsSchema,
	}
}

func dataSourceMemorystoreAuthTokenRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	instance := d.Get("instance").(string)
	instanceRegex := regexp.MustCompile("projects/([^/]+)/locations/([^/]+)/instances/([^/]+)")
	if instanceRegex.MatchString(instance) {
		parts := instanceRegex.FindStringSubmatch(instance)
		if d.Get("project") == "" {
			d.Set("project", parts[1])
		}
		if d.Get("location") == "" {
			d.Set("location", parts[2])
		}
		d.Set("instance", parts[3])
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	d.Set("project", project)

	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return err
	}
	d.Set("location", location)

	tokenId := tpgresource.GetResourceNameFromSelfLink(d.Get("token_id").(string))
	if err := d.Set("token_id", tokenId); err != nil {
		return fmt.Errorf("Error setting token_id: %s", err)
	}

	userId := tpgresource.GetResourceNameFromSelfLink(d.Get("user_id").(string))
	if err := d.Set("user_id", userId); err != nil {
		return fmt.Errorf("Error setting user_id: %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/instances/{{instance}}/tokenAuthUsers/{{user_id}}/authTokens/{{token_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceMemorystoreAuthTokenRead(d, meta)
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
		Name:        "google_memorystore_auth_token",
		ProductName: "memorystore",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceMemorystoreAuthToken(),
	}.Register()
}
