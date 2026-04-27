package compute

import (
	"fmt"
	"log"
	"sort"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-provider-google/google/registry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceGoogleComputeRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleComputeRegionsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"UP", "DOWN"}, false),
			},
		},
	}
}

func dataSourceGoogleComputeRegionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	filter := ""
	if s, ok := d.GetOk("status"); ok {
		filter = fmt.Sprintf(" (status eq %s)", s)
	}

	url := fmt.Sprintf("%sprojects/%s/regions", config.ComputeBasePath, project)
	url, err = transport_tpg.AddQueryParams(url, map[string]string{"filter": filter, "prettyPrint": "false"})
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return err
	}

	regions := flattenRegions(res["items"].([]interface{}))
	log.Printf("[DEBUG] Received Google Compute Regions: %q", regions)

	if err := d.Set("names", regions); err != nil {
		return fmt.Errorf("Error setting names: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s", project))

	return nil
}

func flattenRegions(regions []interface{}) []string {
	result := make([]string, len(regions))
	for i, region := range regions {
		regionMap := region.(map[string]interface{})
		result[i] = regionMap["name"].(string)
	}
	sort.Strings(result)
	return result
}

func init() {
	registry.Schema{
		Name:        "google_compute_regions",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeRegions(),
	}.Register()
}
