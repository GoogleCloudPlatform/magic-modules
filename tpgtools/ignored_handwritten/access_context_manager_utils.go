package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func servicePerimeterImport(d *schema.ResourceData, config *transport_tpg.Config) error {
	// current import_formats can't import fields with forward slashes in their value
	if err := tpgresource.ParseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
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

func accessLevelImport(d *schema.ResourceData, config *transport_tpg.Config) error {
	// current import_formats can't import fields with forward slashes in their value
	if err := tpgresource.ParseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
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
