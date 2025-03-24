package logging

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var loggingProjectBucketConfigSchema = map[string]*schema.Schema{
	"project": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The parent project that contains the logging bucket.`,
	},
	"name": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The resource name of the bucket`,
	},
	"location": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The location of the bucket.`,
	},
	"bucket_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The name of the logging bucket. Logging automatically creates two log buckets: _Required and _Default.`,
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		Description: `An optional description for this bucket.`,
	},
	"locked": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: `Whether the bucket is locked. The retention period on a locked bucket cannot be changed. Locked buckets may only be deleted if they are empty.`,
	},
	"retention_days": {
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     30,
		Description: `Logs will be retained by default for this amount of time, after which they will automatically be deleted. The minimum retention period is 1 day. If this value is set to zero at bucket creation time, the default time of 30 days will be used.`,
	},
	"enable_analytics": {
		Type:             schema.TypeBool,
		Optional:         true,
		Description:      `Enable log analytics for the bucket. Cannot be disabled once enabled.`,
		DiffSuppressFunc: enableAnalyticsBackwardsChangeDiffSuppress,
	},
	"lifecycle_state": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: `The bucket's lifecycle such as active or deleted.`,
	},
	"cmek_settings": {
		Type:        schema.TypeList,
		MaxItems:    1,
		Optional:    true,
		Description: `The CMEK settings of the log bucket. If present, new log entries written to this log bucket are encrypted using the CMEK key provided in this configuration. If a log bucket has CMEK settings, the CMEK settings cannot be disabled later by updating the log bucket. Changing the KMS key is allowed.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: `The resource name of the CMEK settings.`,
				},
				"kms_key_name": {
					Type:     schema.TypeString,
					Required: true,
					Description: `The resource name for the configured Cloud KMS key.
KMS key name format:
"projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]"
To enable CMEK for the bucket, set this field to a valid kmsKeyName for which the associated service account has the required cloudkms.cryptoKeyEncrypterDecrypter roles assigned for the key.
The Cloud KMS key used by the bucket can be updated by changing the kmsKeyName to a new valid key name. Encryption operations that are in progress will be completed with the key that was in use when they started. Decryption operations will be completed using the key that was used at the time of encryption unless access to that key has been revoked.
See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.`,
				},
				"kms_key_version_name": {
					Type:     schema.TypeString,
					Computed: true,
					Description: `The CryptoKeyVersion resource name for the configured Cloud KMS key.
KMS key name format:
"projects/[PROJECT_ID]/locations/[LOCATION]/keyRings/[KEYRING]/cryptoKeys/[KEY]/cryptoKeyVersions/[VERSION]"
For example:
"projects/my-project/locations/us-central1/keyRings/my-ring/cryptoKeys/my-key/cryptoKeyVersions/1"
This is a read-only field used to convey the specific configured CryptoKeyVersion of kms_key that has been configured. It will be populated in cases where the CMEK settings are bound to a single key version.`,
				},
				"service_account_id": {
					Type:     schema.TypeString,
					Computed: true,
					Description: `The service account associated with a project for which CMEK will apply.
Before enabling CMEK for a logging bucket, you must first assign the cloudkms.cryptoKeyEncrypterDecrypter role to the service account associated with the project for which CMEK will apply. Use [v2.getCmekSettings](https://cloud.google.com/logging/docs/reference/v2/rest/v2/TopLevel/getCmekSettings#google.logging.v2.ConfigServiceV2.GetCmekSettings) to obtain the service account ID.
See [Enabling CMEK for Logging Buckets](https://cloud.google.com/logging/docs/routing/managed-encryption-storage) for more information.`,
				},
			},
		},
	},
	"index_configs": {
		Type:        schema.TypeSet,
		MaxItems:    20,
		Optional:    true,
		Description: `A list of indexed fields and related configuration data.`,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"field_path": {
					Type:        schema.TypeString,
					Required:    true,
					Description: `The LogEntry field path to index.`,
				},
				"type": {
					Type:     schema.TypeString,
					Required: true,
					Description: `The type of data in this index
Note that some paths are automatically indexed, and other paths are not eligible for indexing. See [indexing documentation]( https://cloud.google.com/logging/docs/view/advanced-queries#indexed-fields) for details.
For example: jsonPayload.request.status`,
				},
			},
		},
	},
}

func projectBucketConfigID(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	projectID := d.Get("project").(string)
	location := d.Get("location").(string)
	bucketID := d.Get("bucket_id").(string)

	if strings.HasPrefix(projectID, "projects/") {
		// Remove "projects/" prefix if it exists
		projectID = strings.TrimPrefix(projectID, "projects/")
	}

	id := fmt.Sprintf("projects/%s/locations/%s/buckets/%s", projectID, location, bucketID)
	return id, nil
}

// Create Logging Bucket config
func ResourceLoggingProjectBucketConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingProjectBucketConfigAcquireOrCreate("project", projectBucketConfigID),
		Read:   resourceLoggingProjectBucketConfigRead,
		Update: resourceLoggingProjectBucketConfigUpdate,
		Delete: resourceLoggingBucketConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingBucketConfigImportState("project"),
		},
		Schema:        loggingProjectBucketConfigSchema,
		UseJSONNumber: true,
	}
}

func resourceLoggingProjectBucketConfigAcquireOrCreate(parentType string, iDFunc loggingBucketConfigIDFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*transport_tpg.Config)
		userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
		if err != nil {
			return err
		}

		id, err := iDFunc(d, config)
		if err != nil {
			return err
		}

		if parentType == "project" {
			//logging bucket can be created only at the project level, in future api may allow for folder, org and other parent resources

			log.Printf("[DEBUG] Fetching logging bucket config: %#v", id)
			url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", id))
			if err != nil {
				return err
			}

			res, _ := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: userAgent,
			})
			if res == nil {
				log.Printf("[DEBUG] Logging Bucket does not exist %s", id)
				// we need to pass the id in here because we don't want to set it in state
				// until we know there won't be any errors on create
				return resourceLoggingProjectBucketConfigCreate(d, meta, id)
			}
		}

		d.SetId(id)

		return resourceLoggingProjectBucketConfigUpdate(d, meta)
	}
}

func resourceLoggingProjectBucketConfigCreate(d *schema.ResourceData, meta interface{}, id string) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	obj["name"] = d.Get("name")
	obj["description"] = d.Get("description")
	obj["locked"] = d.Get("locked")
	obj["retentionDays"] = d.Get("retention_days")
	// Only set analyticsEnabled if it has been explicitly preferenced.
	analyticsRawValue := d.GetRawConfig().GetAttr("enable_analytics")
	if !analyticsRawValue.IsNull() {
		obj["analyticsEnabled"] = analyticsRawValue.True()
	}
	obj["cmekSettings"] = expandCmekSettings(d.Get("cmek_settings"))
	obj["indexConfigs"] = expandIndexConfigs(d.Get("index_configs"))

	url, err := tpgresource.ReplaceVars(d, config, "{{LoggingBasePath}}projects/{{project}}/locations/{{location}}/buckets:createAsync?bucketId={{bucket_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Bucket: %#v", obj)
	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Body:      obj,
		Timeout:   d.Timeout(schema.TimeoutCreate),
	})
	if err != nil {
		return fmt.Errorf("Error creating Bucket: %s", err)
	}

	d.SetId(id)
	var opRes map[string]interface{}
	// Wait for the operation to complete
	waitErr := LoggingOperationWaitTimeWithResponse(config, res, &opRes, "Bucket to create", userAgent, d.Timeout(schema.TimeoutCreate))
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}

	log.Printf("[DEBUG] Finished creating Bucket %q: %#v", d.Id(), res)

	return resourceLoggingProjectBucketConfigRead(d, meta)
}

func resourceLoggingProjectBucketConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Fetching logging bucket config: %#v", d.Id())

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		log.Printf("[WARN] Unable to acquire logging bucket config at %s", d.Id())

		d.SetId("")
		return err
	}

	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", res["description"]); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("locked", res["locked"]); err != nil {
		return fmt.Errorf("Error setting locked: %s", err)
	}
	if err := d.Set("lifecycle_state", res["lifecycleState"]); err != nil {
		return fmt.Errorf("Error setting lifecycle_state: %s", err)
	}
	if err := d.Set("retention_days", res["retentionDays"]); err != nil {
		return fmt.Errorf("Error setting retention_days: %s", err)
	}
	if err := d.Set("enable_analytics", res["analyticsEnabled"]); err != nil {
		return fmt.Errorf("Error setting enable_analytics: %s", err)
	}

	if err := d.Set("cmek_settings", flattenCmekSettings(res["cmekSettings"])); err != nil {
		return fmt.Errorf("Error setting cmek_settings: %s", err)
	}

	if err := d.Set("index_configs", flattenIndexConfigs(res["indexConfigs"])); err != nil {
		return fmt.Errorf("Error setting index_configs: %s", err)
	}

	return nil
}

func resourceLoggingProjectBucketConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})

	url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
	if err != nil {
		return err
	}

	updateMaskAnalytics := []string{}
	// Check if analytics is being enabled. Analytics enablement is an atomic operation and can not be performed while other fields
	// are being updated, so we enable analytics before updating the rest of the fields.
	if d.HasChange("enable_analytics") {
		obj["analyticsEnabled"] = d.Get("enable_analytics")
		updateMaskAnalytics = append(updateMaskAnalytics, "analyticsEnabled")
		asyncUrl := url + ":updateAsync"
		asyncUrl, err = transport_tpg.AddQueryParams(asyncUrl, map[string]string{"updateMask": strings.Join(updateMaskAnalytics, ",")})
		if err != nil {
			return err
		}
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			RawURL:    asyncUrl,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
		})

		var opRes map[string]interface{}
		// Wait for the operation to complete
		waitErr := LoggingOperationWaitTimeWithResponse(config, res, &opRes, "Updating Bucket with analytics", userAgent, d.Timeout(schema.TimeoutCreate))
		if waitErr != nil {
			return fmt.Errorf("Error updating Logging Bucket Config %q: %s", d.Id(), err)
		}
	}

	obj["retentionDays"] = d.Get("retention_days")
	obj["description"] = d.Get("description")
	obj["cmekSettings"] = expandCmekSettings(d.Get("cmek_settings"))
	obj["indexConfigs"] = expandIndexConfigs(d.Get("index_configs"))

	updateMask := []string{}
	if d.HasChange("retention_days") {
		updateMask = append(updateMask, "retentionDays")
	}
	if d.HasChange("description") {
		updateMask = append(updateMask, "description")
	}
	if d.HasChange("cmek_settings") {
		updateMask = append(updateMask, "cmekSettings")
	}
	if d.HasChange("index_configs") {
		updateMask = append(updateMask, "indexConfigs")
	}

	url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}
	if len(updateMask) > 0 {
		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			RawURL:    url,
			UserAgent: userAgent,
			Body:      obj,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
		})
	}
	if err != nil {
		return fmt.Errorf("Error updating Logging Bucket Config %q: %s", d.Id(), err)
	}

	// Check if locked is being changed (although removal will fail). Locking is
	// an atomic operation and can not be performed while other fields.
	// update locked last so that we lock *after* setting the right settings
	if d.HasChange("locked") {
		updateMaskLocked := []string{"locked"}
		objLocked := map[string]interface{}{
			"locked": d.Get("locked"),
		}
		url, err = transport_tpg.AddQueryParams(url, map[string]string{"updateMask": strings.Join(updateMaskLocked, ",")})
		if err != nil {
			return err
		}

		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "PATCH",
			RawURL:    url,
			UserAgent: userAgent,
			Body:      objLocked,
			Timeout:   d.Timeout(schema.TimeoutUpdate),
		})
		if err != nil {
			return fmt.Errorf("Error updating Logging Bucket Config %q: %s", d.Id(), err)
		}
	}

	return resourceLoggingProjectBucketConfigRead(d, meta)
}

func enableAnalyticsBackwardsChangeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldValue, _ := strconv.ParseBool(old)
	if oldValue {
		return true
	}
	return false
}
