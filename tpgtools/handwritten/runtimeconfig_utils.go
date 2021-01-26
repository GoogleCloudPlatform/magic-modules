package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	runtimeconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/runtimeconfig"
)

func runtimeconfigVariableValidateTextOrValueSet(d *schema.ResourceData, config *Config, res *runtimeconfig.Variable) error {
	// Validate that both text and value are not set
	_, textSet := d.GetOk("text")
	_, valueSet := d.GetOk("value")

	if !textSet && !valueSet {
		return fmt.Errorf("You must specify one of value or text.")
	}

	return nil
}
