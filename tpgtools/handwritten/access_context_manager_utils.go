package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func servicePerimeterImport(d *schema.ResourceData, config *Config) error {
	// current import_formats can't import fields with forward slashes in their value
	if err := parseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return err
	}
	stringParts := strings.Split(d.Get("name").(string), "/")
	if len(stringParts) < 2 {
		return fmt.Errorf("Error parsing parent name. Should be in form accessPolicies/{{policy_id}}/servicePerimeters/{{short_name}}")
	}
	if err := d.Set("parent", fmt.Sprintf("%s/%s", stringParts[0], stringParts[1])); err != nil {
		return fmt.Errorf("Error setting parent, %s", err)
	}
	return nil
}

func accessLevelImport(d *schema.ResourceData, config *Config) error {
	// current import_formats can't import fields with forward slashes in their value
	if err := parseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return err
	}
	stringParts := strings.Split(d.Get("name").(string), "/")
	if len(stringParts) < 2 {
		return fmt.Errorf("Error parsing parent name. Should be in form accessPolicies/{{policy_id}}/accessLevels/{{short_name}}")
	}
	if err := d.Set("parent", fmt.Sprintf("%s/%s", stringParts[0], stringParts[1])); err != nil {
		return fmt.Errorf("Error setting parent, %s", err)
	}
	return nil
}
