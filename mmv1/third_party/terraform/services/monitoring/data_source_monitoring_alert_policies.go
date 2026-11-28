package monitoring

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceMonitoringAlertPolicies returns the schema of the google_monitoring_alert_policies data source.
func DataSourceMonitoringAlertPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMonitoringAlertPoliciesRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose alert policies are listed. If it is not provided, the provider project is used.`,
			},
			"alert_policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The alert policies of the project.`,
				Elem:        &schema.Resource{Schema: dataSourceMonitoringAlertPoliciesItemSchema()},
			},
		},
	}
}

// dataSourceMonitoringAlertPoliciesItemSchema returns the schema of one policy of the list, built from google_monitoring_alert_policy.
func dataSourceMonitoringAlertPoliciesItemSchema() map[string]*schema.Schema {
	policySchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMonitoringAlertPolicy().Schema)
	// These fields are left out: project repeats the argument, and deletion_policy is a setting of the resource.
	tpgresource.DeleteFieldsFromSchema(policySchema, "project", "deletion_policy")
	// A nested block gets no id from the SDK, so the id of the resource is added by hand, as in google_kms_crypto_keys.
	policySchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: `An identifier for the policy with format projects/{{project}}/alertPolicies/{{policy_id}}, the same value as name.`,
	}
	return policySchema
}

// dataSourceMonitoringAlertPoliciesRead lists every policy of the project with the list call of the google_monitoring_alert_policy list resource.
func dataSourceMonitoringAlertPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	policySchema := dataSourceMonitoringAlertPoliciesItemSchema()
	policies := make([]map[string]interface{}, 0)
	// ListMonitoringAlertPolicys follows the pages and hands over each policy as the data of a google_monitoring_alert_policy resource.
	err = ListMonitoringAlertPolicys(config, project, func(rd *schema.ResourceData) error {
		policy := make(map[string]interface{}, len(policySchema))
		for key := range policySchema {
			if key == "id" {
				policy[key] = rd.Id()
			} else {
				policy[key] = rd.Get(key)
			}
		}
		policies = append(policies, policy)
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error retrieving list of alert policies: %w", err)
	}

	if err := d.Set("alert_policies", policies); err != nil {
		return fmt.Errorf("Error setting alert_policies: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(fmt.Sprintf("projects/%s/alertPolicies", project))

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_monitoring_alert_policies",
		ProductName: "monitoring",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceMonitoringAlertPolicies(),
	}.Register()
}
