package compute

import (
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeRollout() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeRolloutRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
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
			"rollout_entity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"orchestrated_entity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"orchestration_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"orchestration_source": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"conflict_behavior": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"wave_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"wave_display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wave_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"orchestrated_wave_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"estimated_total_resources_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"completed_resources_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"failed_resources_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"failed_locations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleComputeRolloutRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	name := d.Get("name").(string)

	id := fmt.Sprintf("projects/%s/global/rollouts/%s", project, name)
	d.SetId(id)

	log.Printf("[DEBUG] Reading Rollout %q", id)

	rollout, err := config.NewComputeClient(userAgent).Rollouts.Get(project, name).Do()
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("Rollout %q", name), id)
	}

	if err := d.Set("description", rollout.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("rollout_plan", rollout.RolloutPlan); err != nil {
		return fmt.Errorf("Error setting rollout_plan: %s", err)
	}
	if err := d.Set("state", rollout.State); err != nil {
		return fmt.Errorf("Error setting state: %s", err)
	}
	if err := d.Set("current_wave_number", rollout.CurrentWaveNumber); err != nil {
		return fmt.Errorf("Error setting current_wave_number: %s", err)
	}
	if err := d.Set("self_link", rollout.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if rollout.RolloutEntity != nil {
		if err := d.Set("rollout_entity", flattenRolloutEntity(rollout.RolloutEntity)); err != nil {
			return fmt.Errorf("Error setting rollout_entity: %s", err)
		}
	}
	if err := d.Set("wave_details", flattenRolloutWaveDetails(rollout.WaveDetails)); err != nil {
		return fmt.Errorf("Error setting wave_details: %s", err)
	}

	return nil
}

func flattenRolloutEntity(entity interface{}) []interface{} {
	if entity == nil || reflect.ValueOf(entity).IsNil() {
		return nil
	}
	m := make(map[string]interface{})

	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	orchestratedEntity := val.FieldByName("OrchestratedEntity")
	if orchestratedEntity.IsValid() && !orchestratedEntity.IsNil() {
		m["orchestrated_entity"] = flattenOrchestratedEntity(orchestratedEntity.Interface())
	}
	return []interface{}{m}
}

func flattenOrchestratedEntity(entity interface{}) []interface{} {
	if entity == nil || reflect.ValueOf(entity).IsNil() {
		return nil
	}
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return []interface{}{
		map[string]interface{}{
			"orchestration_action": getStringField(val, "OrchestrationAction"),
			"orchestration_source": getStringField(val, "OrchestrationSource"),
			"conflict_behavior":    getStringField(val, "ConflictBehavior"),
		},
	}
}

func flattenRolloutWaveDetails(details interface{}) []interface{} {
	if details == nil || reflect.ValueOf(details).IsNil() {
		return nil
	}
	val := reflect.ValueOf(details)
	if val.Kind() != reflect.Slice {
		return nil
	}

	var res []interface{}
	for i := 0; i < val.Len(); i++ {
		v := val.Index(i)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		m := map[string]interface{}{
			"wave_display_name": getStringField(v, "WaveDisplayName"),
			"wave_number":       getIntField(v, "WaveNumber"),
		}

		orchestratedWaveDetails := v.FieldByName("OrchestratedWaveDetails")
		if orchestratedWaveDetails.IsValid() && !orchestratedWaveDetails.IsNil() {
			m["orchestrated_wave_details"] = flattenOrchestratedWaveDetails(orchestratedWaveDetails.Interface())
		}
		res = append(res, m)
	}
	return res
}

func flattenOrchestratedWaveDetails(details interface{}) []interface{} {
	if details == nil || reflect.ValueOf(details).IsNil() {
		return nil
	}
	val := reflect.ValueOf(details)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return []interface{}{
		map[string]interface{}{
			"estimated_total_resources_count": getIntField(val, "EstimatedTotalResourcesCount"),
			"completed_resources_count":       getIntField(val, "CompletedResourcesCount"),
			"failed_resources_count":          getIntField(val, "FailedResourcesCount"),
			"failed_locations":                getStringSliceField(val, "FailedLocations"),
		},
	}
}

func getStringField(v reflect.Value, name string) string {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

func getIntField(v reflect.Value, name string) int {
	f := v.FieldByName(name)
	// Compute API uses int64 usually, TF uses int
	if f.IsValid() && (f.Kind() == reflect.Int || f.Kind() == reflect.Int64) {
		return int(f.Int())
	}
	return 0
}

func getStringSliceField(v reflect.Value, name string) []string {
	f := v.FieldByName(name)
	if f.IsValid() && f.Kind() == reflect.Slice {
		var res []string
		for i := 0; i < f.Len(); i++ {
			res = append(res, f.Index(i).String())
		}
		return res
	}
	return nil
}
