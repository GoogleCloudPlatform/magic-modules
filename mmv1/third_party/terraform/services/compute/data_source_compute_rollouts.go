package compute

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	compute "google.golang.org/api/compute/v0.beta"
)

func DataSourceGoogleComputeRollouts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGoogleComputeRolloutsRead,

		Schema: map[string]*schema.Schema{
			"rollouts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rollout_plan": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_wave_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceGoogleComputeRolloutsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return diag.FromErr(err)
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	client := config.NewComputeClient(userAgent).Rollouts.List(project)
	if filter, ok := d.GetOk("filter"); ok {
		client.Filter(filter.(string))
	}

	var rollouts []map[string]interface{}

	for {
		resp, err := client.Do()
		if err != nil {
			return diag.FromErr(err)
		}
		for _, rollout := range resp.Items {
			rollouts = append(rollouts, flattenRolloutResource(rollout, project))
		}
		if resp.NextPageToken == "" {
			break
		}
		client.PageToken(resp.NextPageToken)
	}

	if err := d.Set("rollouts", rollouts); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project", project); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("projects/%s/global/rollouts", project))

	return nil
}

func flattenRolloutResource(rollout *compute.Rollout, project string) map[string]interface{} {
	return map[string]interface{}{
		"name":                rollout.Name,
		"description":         rollout.Description,
		"rollout_plan":        rollout.RolloutPlan,
		"state":               rollout.State,
		"current_wave_number": rollout.CurrentWaveNumber,
		"self_link":           rollout.SelfLink,
		"project":             project,
	}
}
