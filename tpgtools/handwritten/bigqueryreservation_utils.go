package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getProjectField(d TerraformResourceData, config *Config) (string, error) {
	res, ok := d.GetOk("project")
	if !ok {
		return "", nil
	}

	return res.(string), nil
}

func resourceBigqueryReservationAssignmentCustomImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<location>[^/]+)/reservations/(?P<reservation>[^/]+)/assignments/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<location>[^/]+)/(?P<reservation>[^/]+)/(?P<name>[^/]+)",
		"(?P<location>[^/]+)/(?P<reservation>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	if err := d.Set("assignee", d.Get("project")); err != nil {
		return nil, fmt.Errorf("error setting assignee in state: %s", err)
	}

	// Replace import id for the resource id
	id, err := replaceVarsForId(d, config, "projects/{{project}}/locations/{{location}}/reservations/{{reservation}}/assignments/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
