package google

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/google/compute/beta"
)

func deleteComputeNetworkDefaultRoutes(d *schema.ResourceData, config *Config, res *compute.Network) error {
	if d.Get("delete_default_routes_on_create").(bool) {
		url, err := replaceVars(d, config, "{{ComputeBasePath}}projects/{{project}}/global/networks")
		networkLink := fmt.Sprintf("%s/%s", url, d.Get("name").(string))
		filter := fmt.Sprintf("(network=\"%s\") AND (destRange=\"0.0.0.0/0\")", networkLink)
		log.Printf("[DEBUG] Getting routes for network %q with filter '%q'", d.Get("name").(string), filter)
		routes, err := config.clientComputeDCL.ListRoute(context.Background(), *res.Project)
		if err != nil {
			return fmt.Errorf("Error listing routes in proj: %s", err)
		}
		log.Printf("[DEBUG] Found %d routes rules in %q network", len(routes.Items), d.Get("name").(string))
		for _, route := range routes.Items {
			if *route.Network == networkLink && *route.DestRange == "0.0.0.0/0" {
				err := config.clientComputeDCL.DeleteRoute(context.Background(), route)
				if err != nil {
					return fmt.Errorf("Error deleting route: %s", err)
				}
			}
		}
	}

	return nil
}

func getVpnTunnelLink(config *Config, project string, region string, tunnel string) (string, error) {
	if !strings.Contains(tunnel, "/") {
		// Tunnel value provided is just the name, lookup the tunnel SelfLink
		tunnelData, err := config.clientCompute.VpnTunnels.Get(
			project, region, tunnel).Do()
		if err != nil {
			return "", fmt.Errorf("Error reading tunnel: %s", err)
		}
		tunnel = tunnelData.SelfLink
	}

	return tunnel, nil

}
