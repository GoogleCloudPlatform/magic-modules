package iap

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleIapBrand() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceIapBrand().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "project")

	dsSchema["brand"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: "The ID of the brand. If not provided, the project's brand will be returned.",
	}

	return &schema.Resource{
		Read:   dataSourceGoogleIapBrandRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleIapBrandRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	// If brand ID is not provided, we must fetch it using the project ID,
	// as the brand ID is the project number.
	if v, ok := d.GetOk("brand"); !ok || v.(string) == "" {
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return err
		}

		projectID, err := tpgresource.GetProject(d, config)
		if err != nil {
			return err
		}

		rmClient := config.NewResourceManagerClient(userAgent)
		getProjectCall := rmClient.Projects.Get(projectID)

		// Handle UserProjectOverride if it's enabled
		if config.UserProjectOverride {
			billingProject := projectID
			if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
				billingProject = bp
			}
			getProjectCall.Header().Add("X-Goog-User-Project", billingProject)
		}

		project, err := getProjectCall.Do()
		if err != nil {
			return fmt.Errorf("Error fetching project details for '%s': %w", projectID, err)
		}

		if err := d.Set("brand", strconv.FormatInt(project.ProjectNumber, 10)); err != nil {
			return fmt.Errorf("Error setting brand: %s", err)
		}
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/brands/{{brand}}")

	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	err = resourceIapBrandRead(d, meta)
	if err != nil {
		return err
	}

	if d.Id() == "" {
		return fmt.Errorf("IAP Brand not found for project %s", d.Get("project"))
	}
	return nil
}
