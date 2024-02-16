package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeForwardingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeForwardingRulesRead,

		Schema: map[string]*schema.Schema{

			// "self_link": {
			// 	Type:     schema.TypeString,
			// 	Computed: true,
			// },

			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func dataSourceGoogleComputeForwardingRulesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/regions/%s/forwardingRules", project, region)
	d.SetId(id)

	forwardingRulesAggregatedList, err := config.NewComputeClient(userAgent).ForwardingRules.List(project, region).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Forwarding Rules Not Found : %s", project))
	}

	forwardingRules := []map[string]interface{}{}

	for i := 0; i < len(forwardingRulesAggregatedList.Items); i++ {
		rule := forwardingRulesAggregatedList.Items[i]
		name := rule.Name
		network := rule.Network
		subnet := rule.Subnetwork
		backend := rule.BackendService
		ip := rule.IPAddress
		serviceName := rule.ServiceName
		serviceLabel := rule.ServiceLabel
		description := rule.Description
		selfLink := rule.SelfLink
		mappedData := map[string]interface{}{
			"name":          name,
			"network":       network,
			"subnet":        subnet,
			"backend":       backend,
			"ip":            ip,
			"service name":  serviceName,
			"service label": serviceLabel,
			"description":   description,
			"self link":     selfLink,
		}
		forwardingRules = append(forwardingRules, mappedData)
	}

	if err := d.Set("rules", forwardingRules); err != nil {
		return fmt.Errorf("Error setting the forwarding rules names: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting the network names: %s", err)
	}

	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("Error setting the region: %s", err)
	}

	return nil
}
