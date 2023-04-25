package google

import (
	"fmt"

	runtimeconfig "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/runtimeconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func runtimeconfigVariableValidateTextOrValueSet(d *schema.ResourceData, config *transport_tpg.Config, res *runtimeconfig.Variable) error {
	// Validate that both text and value are not set
	_, textSet := d.GetOk("text")
	_, valueSet := d.GetOk("value")

	if !textSet && !valueSet {
		return fmt.Errorf("You must specify one of value or text.")
	}

	return nil
}
