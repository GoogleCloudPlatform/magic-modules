package google

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute/beta"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func deleteComputeNetworkDefaultRoutes(d *schema.ResourceData, client *compute.Client, config *Config, res *compute.Network) error {
	if d.Get("delete_default_routes_on_create").(bool) {
		routes, err := client.ListRoute(context.Background(), &compute.Route{Project: res.Project})
		if err != nil {
			return fmt.Errorf("Error listing routes in proj: %s", err)
		}
		log.Printf("[DEBUG] Found %d routes rules in %q network", len(routes.Items), d.Get("name").(string))
		for _, route := range routes.Items {
			log.Printf("[DEBUG] route.Network: %s, network selfLink: %s", *route.Network, *res.SelfLink)
			if dcl.SelfLinkToSelfLink(route.Network, res.SelfLink) && *route.DestRange == "0.0.0.0/0" {
				err := client.DeleteRoute(context.Background(), route)
				if err != nil {
					return fmt.Errorf("Error deleting route: %s", err)
				}
			}
		}
	}

	return nil
}

