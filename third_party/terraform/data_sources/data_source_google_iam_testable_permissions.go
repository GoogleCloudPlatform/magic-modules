package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGoogleIamTestablePermissions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleIamTestablePermissionsRead,
		Schema: map[string]*schema.Schema{
			"full_resource_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_support_level": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SUPPORTED",
			},
			"stage": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GA",
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"custom_support_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_disabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleIamTestablePermissionsRead(d *schema.ResourceData, meta interface{}) (err error) {
	config := meta.(*Config)
	body := make(map[string]interface{}, 0)
	body["pageSize"] = 500
	permissions := make([]map[string]interface{}, 0)
	custom_support_level := d.Get("custom_support_level").(string)
	if err = validatePermissionCustomSupport(custom_support_level); err != nil {
		return err
	}

	stage := d.Get("stage").(string)
	if err = validatePermissionStage(stage); err != nil {
		return err
	}

	for {
		url := "https://iam.googleapis.com/v1/permissions:queryTestablePermissions"
		body["fullResourceName"] = d.Get("full_resource_name").(string)
		res, err := sendRequest(config, "POST", "", url, body)
		if err != nil {
			return fmt.Errorf("Error retrieving permissions: %s", err)
		}

		pagePermissions := flattenTestablePermissionsList(res["permissions"], custom_support_level, stage)
		permissions = append(permissions, pagePermissions...)
		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			body["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("permissions", permissions); err != nil {
		return fmt.Errorf("Error retrieving permissions: %s", err)
	}

	d.SetId(d.Get("full_resource_name").(string))
	return nil
}

func flattenTestablePermissionsList(v interface{}, custom_support_level string, stage string) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	permissions := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		p := raw.(map[string]interface{})

		if _, ok := p["name"]; ok {
			csl := true
			if custom_support_level == "SUPPORTED" {
				csl = p["customRolesSupportLevel"] == nil || p["customRolesSupportLevel"] == "SUPPORTED"
			} else {
				csl = p["customRolesSupportLevel"] == custom_support_level
			}

			if csl && p["stage"] == stage {
				permissions = append(permissions, map[string]interface{}{
					"name":                 p["name"],
					"title":                p["title"],
					"stage":                p["stage"],
					"api_disabled":         p["apiDisabled"],
					"custom_support_level": p["customRolesSupportLevel"],
				})
			}
		}
	}

	return permissions
}

func validatePermissionCustomSupport(val string) error {
	allowed := []string{"NOT_SUPPORTED", "SUPPORTED", "TESTING"}
	if !sliceContainsString(allowed, val) {
		return errors.New("custom_support_level must be one of \"NOT_SUPPORTED\", \"SUPPORTED\", \"TESTING\"")
	}
	return nil
}

func validatePermissionStage(val string) error {
	allowed := []string{"ALPHA", "BETA", "GA", "DEPRECATED"}
	if !sliceContainsString(allowed, val) {
		return errors.New("stage must be one of \"ALPHA\", \"BETA\", \"GA\", \"DEPRECATED\"")
	}
	return nil
}

func sliceContainsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
