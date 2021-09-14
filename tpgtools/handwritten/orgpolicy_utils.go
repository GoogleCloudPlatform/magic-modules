package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrgPolicyPolicyCustomImport(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"^(?P<parent>[^/]+/?[^/]*)/policies/(?P<name>[^/]+)",
		"^(?P<parent>[^/]+/?[^/]*)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return err
	}

	// Replace import id for the resource id
	id, err := replaceVarsRecursive(d, config, "{{parent}}/policies/{{name}}", false, 0)
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return nil
}
