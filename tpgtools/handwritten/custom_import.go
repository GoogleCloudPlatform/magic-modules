package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func sourceRepoImport(d *schema.ResourceData, config *Config) error {
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/repos/(?P<name>.+)",
		"(?P<name>.+)",
	}, d, config); err != nil {
		return err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/repos/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}

func runtimeconfigVariableImport(d *schema.ResourceData, config *Config) error {
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/configs/(?P<parent>[^/]+)/variables/(?P<name>.+)",
		"(?P<parent>[^/]+)/(?P<name>.+)",
	}, d, config); err != nil {
		return err
	}


	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/configs/{{parent}}/variables/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}
