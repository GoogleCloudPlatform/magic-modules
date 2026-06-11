package servicenetworking

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
)

func DataSourceGoogleServiceNetworkingPeeredDNSDomain() *schema.Resource {
	return &schema.Resource{
		Read: resourceGoogleServiceNetworkingPeeredDNSDomainRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			//UDP schema start
			"deletion_policy": tpgresource.DeletionPolicySchemaEntry("DELETE"),
			//UDP schema end
		},
	}
}

func init() {
	registry.Schema{
		Name:        "google_service_networking_peered_dns_domain",
		ProductName: "servicenetworking",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleServiceNetworkingPeeredDNSDomain(),
	}.Register()
}
