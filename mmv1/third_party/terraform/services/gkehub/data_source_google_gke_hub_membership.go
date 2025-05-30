package gkehub

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleGkeHubMembership() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceGKEHubMembership().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "membership_id")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "location")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleGkeHubMembershipRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleGkeHubMembershipRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/memberships/{{membership_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceGKEHubMembershipRead(d, meta)
	if err != nil {
		return err
	}

	// No labels or annotations for Membership datasource
	return nil
}
