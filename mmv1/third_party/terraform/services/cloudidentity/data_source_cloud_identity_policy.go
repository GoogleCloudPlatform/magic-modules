package cloudidentity

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudIdentityPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleCloudIdentityPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource name of the policy to retrieve.`,
			},
			"customer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The customer that the policy belongs to.`,
			},
			"policy_query": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The CEL query that defines which entities the policy applies to.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The query that defines which entities the policy applies to.",
						},
						"group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The group that the policy applies to.",
						},
						"org_unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The org unit that the policy applies to.",
						},
						"sort_order": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The sort order of the policy.",
						},
					},
				},
			},
			"setting": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The setting configured by this policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the policy.`,
			},
		},
	}
}

func dataSourceGoogleCloudIdentityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	name, ok := d.GetOk("name")
	if !ok {
		return fmt.Errorf("error getting policy name")
	}

	policiesGetCall := config.NewCloudIdentityClient(userAgent).Policies.Get(name.(string))

	if config.UserProjectOverride {
		billingProject := ""
		// err may be nil - project isn't required for this resource
		if project, err := tpgresource.GetProject(d, config); err == nil {
			billingProject = project
		}

		// err == nil indicates that the billing_project value was found
		if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
			billingProject = bp
		}

		if billingProject != "" {
			policiesGetCall.Header().Set("X-Goog-User-Project", billingProject)
		}
	}

	resp, err := policiesGetCall.Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("CloudIdentityPolicy %q", d.Id()), "Policies")
	}

	if err := d.Set("customer", resp.Customer); err != nil {
		return fmt.Errorf("error setting policy customer: %s", err)
	}

	if resp.PolicyQuery != nil {
		pq := map[string]interface{}{
			"query":      resp.PolicyQuery.Query,
			"group":      resp.PolicyQuery.Group,
			"org_unit":   resp.PolicyQuery.OrgUnit,
			"sort_order": resp.PolicyQuery.SortOrder,
		}
		if err := d.Set("policy_query", []interface{}{pq}); err != nil {
			return fmt.Errorf("error setting policy_query: %s", err)
		}
	}

	if resp.Setting != nil {
		settingBytes, err := json.Marshal(resp.Setting)
		if err != nil {
			return fmt.Errorf("error marshalling policy setting: %s", err)
		}
		if err := d.Set("setting", string(settingBytes)); err != nil {
			return fmt.Errorf("error setting policy setting: %s", err)
		}
	}
	if err := d.Set("type", resp.Type); err != nil {
		return fmt.Errorf("error setting policy type: %s", err)
	}

	d.SetId(resp.Name)
	return nil
}
