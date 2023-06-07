package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/provider"
)

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	return provider.Provider()
}
