package servicenetworking

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func resourceServiceNetworkingVPCServiceControlsCustomCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	return resourceServiceNetworkingVPCServiceControlsSet(d, meta, config)
}
