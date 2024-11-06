package accesscontextmanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAccessContextManagerAccessPolicyMap() *schema.Resource {
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

func dataSourceAccessContextManagerAccessPolicyMapRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AccessContextManagerBasePath}}accessPolicies")
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

	// if err := d.Set("name", res["name"]); err != nil {
	// 	return fmt.Errorf("Error setting name: %s", err)
	// }
	// if err := d.Set("account_email", res["accountEmail"]); err != nil {
	// 	return fmt.Errorf("Error setting account_email: %s", err)
	// }
	// d.SetId(res["name"].(string))

	if err := d.Set("name", "foobar"); err != nil && res != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	return nil
}
