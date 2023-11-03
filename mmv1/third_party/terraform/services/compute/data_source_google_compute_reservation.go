package compute

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	computeReservationIdTemplate = "projects/%s/zones/%s/reservations/%s"
	computeReservationLinkRegex  = regexp.MustCompile("projects/(.+)/zones/(.+)/reservations/(.+)$")
)

func DataSourceGoogleComputeReservation() *schema.Resource {
	// Generate datasource schema from resource
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeReservation().Schema)

	// Set 'Required' schema elements
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "name")
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "zone")

	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")

	return &schema.Resource{
		Read:   dataSourceGoogleComputeReservationRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleComputeReservationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	err := resourceComputeReservationRead(d, meta)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	d.SetId(fmt.Sprintf("projects/%s/zones/%s/reservations/%s", project, zone, name))
	return nil
}

type ComputeReservationId struct {
	Project string
	Zone    string
	Name    string
}

func (s ComputeReservationId) CanonicalId() string {
	return fmt.Sprintf(computeReservationIdTemplate, s.Project, s.Zone, s.Name)
}

func ParseComputeReservationId(id string, config *transport_tpg.Config) (*ComputeReservationId, error) {
	var parts []string
	if computeReservationLinkRegex.MatchString(id) {
		parts = computeReservationLinkRegex.FindStringSubmatch(id)

		return &ComputeReservationId{
			Project: parts[1],
			Zone:    parts[2],
			Name:    parts[3],
		}, nil
	} else {
		parts = strings.Split(id, "/")
	}

	if len(parts) == 3 {
		return &ComputeReservationId{
			Project: parts[0],
			Zone:    parts[1],
			Name:    parts[2],
		}, nil
	} else if len(parts) == 2 {
		// Project is optional.
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{zone}/{name}` id format.")
		}

		return &ComputeReservationId{
			Project: config.Project,
			Zone:    parts[0],
			Name:    parts[1],
		}, nil
	} else if len(parts) == 1 {
		// Project and zone is optional
		if config.Project == "" {
			return nil, fmt.Errorf("The default project for the provider must be set when using the `{name}` id format.")
		}
		if config.Zone == "" {
			return nil, fmt.Errorf("The default zone for the provider must be set when using the `{name}` id format.")
		}

		return &ComputeReservationId{
			Project: config.Project,
			Zone:    config.Zone,
			Name:    parts[0],
		}, nil
	}

	return nil, fmt.Errorf("Invalid compute reservation id. Expecting resource link, `{project}/{zone}/{name}`, `{zone}/{name}` or `{name}` format.")
}
