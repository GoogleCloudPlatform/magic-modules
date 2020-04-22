package google

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var loggingBucketConfigSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"location": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"bucket_id": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
	"description": {
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	},
	"retention_days": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  30,
	},
	"lifecycle_state": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

type loggingBucketConfigIDFunc func(d *schema.ResourceData, config *Config) (string, error)

// ResourceLoggingBucketConfig creates a resource definition by merging a unique field (eg: folder) to a generic logging bucket
// config resource. In practice the only difference between these resources is the url location.
func ResourceLoggingBucketConfig(parentSpecificSchema map[string]*schema.Schema, iDFunc loggingBucketConfigIDFunc) *schema.Resource {
	// func ResourceLoggingBucketConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingBucketConfigAcquire(iDFunc),
		Read:   resourceLoggingBucketConfigRead(),
		Update: resourceLoggingBucketConfigUpdate(),
		Delete: resourceLoggingBucketConfigDelete(),
		Schema: mergeSchemas(loggingBucketConfigSchema, parentSpecificSchema),
	}
}

func resourceLoggingBucketConfigAcquire(iDFunc loggingBucketConfigIDFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		id, err := iDFunc(d, config)
		if err != nil {
			return err
		}

		d.SetId(id)

		return resourceLoggingBucketConfigUpdate()(d, meta)
	}
}

func resourceLoggingBucketConfigRead() schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		log.Printf("[DEBUG] Fetching logging bucket config: %#v", d.Id())

		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
		if err != nil {
			return err
		}

		res, err := sendRequest(config, "GET", "", url, nil)
		if err != nil {
			log.Printf("[WARN] Unable to acquire logging bucket config at %s", d.Id())

			d.SetId("")
			return err
		}

		d.Set("name", res["name"])
		d.Set("description", res["description"])
		d.Set("lifecycle_state", res["lifecycleState"])
		d.Set("retention_days", res["retentionDays"])

		return nil
	}
}

func resourceLoggingBucketConfigUpdate() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)

		obj := make(map[string]interface{})

		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
		obj["retentionDays"] = d.Get("retention_days")
		obj["description"] = d.Get("description")

		updateMask := []string{}
		if d.HasChange("retention_days") {
			updateMask = append(updateMask, "retentionDays")
		}
		if d.HasChange("description") {
			updateMask = append(updateMask, "description")
		}
		url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
		if err != nil {
			return err
		}

		_, err = sendRequestWithTimeout(config, "PATCH", "", url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Logging Bucket Config %q: %s", d.Id(), err)
		}

		return resourceLoggingBucketConfigRead()(d, meta)
	}
}

func resourceLoggingBucketConfigDelete() schema.DeleteFunc {
	return func(d *schema.ResourceData, meta interface{}) error {

		log.Printf("[WARN] Logging bucket configs cannot be deleted. Removing logging bucket config from state: %#v", d.Id())
		d.SetId("")

		return nil
	}
}
