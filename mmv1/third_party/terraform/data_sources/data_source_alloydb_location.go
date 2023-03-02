package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceAlloydbLocation() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceAlloydbLocationRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The canonical id for the location. For example: "us-east1".`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `Resource name for the location, which may vary between implementations. For example: "projects/example-project/locations/us-east1`,
			},
			"location_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `The canonical id for this location. For example: "us-east1".`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: `The friendly name for this location, typically a nearby city name. For example, "Tokyo".`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: `Cross-service attributes for the location. For example {"cloud.googleapis.com/region": "us-east1"}`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Optional:    true,
				Description: `Service-specific metadata. For example the available capacity at the given location.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlloydbLocationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	location := d.Get("location").(string)

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	url, err := replaceVars(d, config, "{{AlloydbBasePath}}projects/{{project}}/locations/{{location}}")
	if err != nil {
		return fmt.Errorf("Error setting api endpoint")
	}
	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("location %q", d.Id()))
	}
	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("location_id", res["locationId"]); err != nil {
		return fmt.Errorf("Error setting location_id: %s", err)
	}
	if err := d.Set("display_name", res["displayName"]); err != nil {
		return fmt.Errorf("Error setting displayName: %s", err)
	}
	if res["labels"] != nil {
		labels := make(map[string]string)
		for k, v := range res["labels"].(map[string]interface{}) {
			labels[k] = v.(string)
		}
		if err := d.Set("labels", labels); err != nil {
			return fmt.Errorf("Error setting labels: %s", err)
		}
	}
	if res["metadata"] != nil {
		metadata := make(map[string]string)
		for k, v := range res["metadata"].(map[interface{}]interface{}) {
			metadata[k.(string)] = v.(string)
		}
		if err := d.Set("metadata", metadata); err != nil {
			return fmt.Errorf("Error setting metadata: %s", err)
		}
	}
	d.SetId(fmt.Sprintf("projects/%s/locations/%s", project, location))
	return nil
}
