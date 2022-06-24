package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDataFusionInstanceVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDataFusionInstanceVersionsRead,
		Schema: map[string]*schema.Schema{
			"location": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_number": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDataFusionInstanceVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	location := d.Get("location")
	url, err := replaceVars(d, config, "{{DataFusionBasePath}}projects/{{project}}/locations/{{location}}/versions")
	if err != nil {
		return err
	}

	versions, err := paginatedListRequest(project, url, userAgent, config, flattenGoogleDataFusionInstanceVersions)
	if err != nil {
		return fmt.Errorf("Error listing Data Fusion instance versions: %s", err)
	}

	log.Printf("[DEBUG] Received Data Fusion Instance Versions: %q", versions)

	if err := d.Set("instance_versions", versions); err != nil {
		return fmt.Errorf("Error setting instance_versions: %s", err)
	}
	if err := d.Set("location", location); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations/%s", project, location))

	return nil
}

func flattenGoogleDataFusionInstanceVersions(resp map[string]interface{}) []interface{} {
	verObjList := resp["availableVersions"].([]interface{})
	versions := make([]interface{}, len(verObjList))
	for i, v := range verObjList {
		verObj := v.(map[string]interface{})
		versions[i] = map[string]interface{}{
			"version_number": verObj["versionNumber"],
		}
	}
	return versions
}
