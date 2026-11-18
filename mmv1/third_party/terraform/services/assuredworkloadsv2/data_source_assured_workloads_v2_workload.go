package assuredworkloadsv2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func DataSourceAssuredWorkloadsV2Workload() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceAssuredWorkloadsV2Workload().Schema)
	tpgresource.DeleteFieldsFromSchema(dsSchema, "resource_config")

	dsSchema["name"].Optional = true
	dsSchema["name"].Computed = true
	dsSchema["name"].ExactlyOneOf = []string{"name", "workload_id"}
	dsSchema["name"].ConflictsWith = []string{"organization", "location", "workload_id"}
	dsSchema["name"].ValidateFunc = validation.StringMatch(
		assuredWorkloadsV2WorkloadNameRegexp,
		"must be in the format organizations/{organization}/locations/{location}/workloads/{workload_id}",
	)

	dsSchema["organization"].Optional = true
	dsSchema["organization"].Computed = true
	dsSchema["organization"].RequiredWith = []string{"location", "workload_id"}
	dsSchema["organization"].ConflictsWith = []string{"name"}

	dsSchema["location"].Optional = true
	dsSchema["location"].Computed = true
	dsSchema["location"].RequiredWith = []string{"organization", "workload_id"}
	dsSchema["location"].ConflictsWith = []string{"name"}

	dsSchema["workload_id"].Optional = true
	dsSchema["workload_id"].Computed = true
	dsSchema["workload_id"].ExactlyOneOf = []string{"name", "workload_id"}
	dsSchema["workload_id"].RequiredWith = []string{"organization", "location"}
	dsSchema["workload_id"].ConflictsWith = []string{"name"}

	return &schema.Resource{
		Read:   dataSourceAssuredWorkloadsV2WorkloadRead,
		Schema: dsSchema,
	}
}

func resolveDataSourceAssuredWorkloadsV2WorkloadId(d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("name"); ok && v.(string) != "" {
		name := v.(string)
		parts := assuredWorkloadsV2WorkloadNameRegexp.FindStringSubmatch(name)
		if len(parts) != 4 {
			return "", fmt.Errorf("invalid workload name %q: must be in the format organizations/{organization}/locations/{location}/workloads/{workload_id}", name)
		}
		if err := d.Set("organization", parts[1]); err != nil {
			return "", fmt.Errorf("Error setting organization: %s", err)
		}
		if err := d.Set("location", parts[2]); err != nil {
			return "", fmt.Errorf("Error setting location: %s", err)
		}
		if err := d.Set("workload_id", parts[3]); err != nil {
			return "", fmt.Errorf("Error setting workload_id: %s", err)
		}
		return name, nil
	}

	org, okOrg := d.GetOk("organization")
	loc, okLoc := d.GetOk("location")
	workloadID, okID := d.GetOk("workload_id")
	if !okOrg || !okLoc || !okID {
		return "", fmt.Errorf("either name or all of organization, location, and workload_id must be specified")
	}

	name := fmt.Sprintf("organizations/%s/locations/%s/workloads/%s", org.(string), loc.(string), workloadID.(string))
	if err := d.Set("name", name); err != nil {
		return "", fmt.Errorf("Error setting name: %s", err)
	}
	return name, nil
}

func dataSourceAssuredWorkloadsV2WorkloadRead(d *schema.ResourceData, meta interface{}) error {
	id, err := resolveDataSourceAssuredWorkloadsV2WorkloadId(d)
	if err != nil {
		return err
	}
	d.SetId(id)

	err = resourceAssuredWorkloadsV2WorkloadRead(d, meta)
	if err != nil {
		return err
	}
	if d.Id() == "" {
		return fmt.Errorf("%s not found", id)
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_assured_workloads_v2_workload",
		ProductName: "assuredworkloadsv2",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceAssuredWorkloadsV2Workload(),
	}.Register()
}
