package accesscontextmanager

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAccessContextManagerAccessPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAccessContextManagerAccessPolicyRead,
		Schema: map[string]*schema.Schema{
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAccessContextManagerAccessPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}accessPolicies?parent={{parent}}")
	if err != nil {
		return err
	}

	billingProject := ""

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("AccessContextManagerAccessPolicy %q", d.Id()), url)
	}

	if res == nil {
		return fmt.Errorf("Error fetching policies: %s", err)
	}

	target_policy_parent := d.Get("parent").(string)
	fetchedPolicies := res["accessPolicies"].([]interface{})
	for _, policy := range fetchedPolicies {
		fetched_policy_parent := policy.(map[string]interface{})["parent"]

		if fetched_policy_parent == target_policy_parent {
			fetched_policy_name := policy.(map[string]interface{})["name"].(string)

			name_without_prefix := strings.Split(fetched_policy_name, "accessPolicies/")[1]
			d.SetId(name_without_prefix)
			if err := d.Set("name", name_without_prefix); err != nil {
				return fmt.Errorf("Error setting name: %s", err)
			}
			return nil
		}
	}

	return nil
}
