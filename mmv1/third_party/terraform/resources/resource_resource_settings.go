package google

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	resourceSettingsV1 "google.golang.org/api/resourcesettings/v1"
)

func resourceGoogleOrganizationResourceSetting() *schema.Resource {
	return resourceGoogleXResourceSetting("organization")
}

/*
// TODO: Release once testable.
func resourceGoogleFolderResourceSetting() *schema.Resource {
	return resourceGoogleXResourceSetting("folder")
}

func resourceGoogleProjectResourceSetting() *schema.Resource {
	return resourceGoogleXResourceSetting("project")
}
*/

const (
	resourceSettingKeyBooleanValue  = "local_value.0.boolean_value"
	resourceSettingKeyStringValue   = "local_value.0.string_value"
	resourceSettingKeyEnumValue     = "local_value.0.enum_value"
	resourceSettingKeyDurationValue = "local_value.0.duration_value"
)

func resourceGoogleXResourceSetting(parentType string) *schema.Resource {
	localValueKeys := []string{
		resourceSettingKeyBooleanValue,
		resourceSettingKeyStringValue,
		resourceSettingKeyEnumValue,
		resourceSettingKeyDurationValue,
	}

	return &schema.Resource{
		Create: resourceGoogleResourceSettingCreate(parentType),
		Read:   resourceGoogleResourceSettingRead(parentType),
		Update: resourceGoogleResourceSettingUpdate(parentType),
		Delete: resourceGoogleResourceSettingDelete(parentType),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Read:   schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Example: "folder_id".
			resourceSettingParentKey(parentType): {
				Type:        schema.TypeString,
				Description: fmt.Sprintf(`The %s id the resource setting with be applied to.`, parentType),
				Required:    true,
				ForceNew:    true,
			},
			"setting_name": {
				Type:        schema.TypeString,
				Description: `The resource settings name. For example, "gcp-enableMyFeature".`,
				Required:    true,
				ForceNew:    true,
			},
			"local_value": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: fmt.Sprintf(`The configured value of the setting at the %s.`, parentType),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"boolean_value": {
							Type:         schema.TypeBool,
							Optional:     true,
							Description:  `Holds the value for a local value field with boolean type.`,
							AtLeastOneOf: localValueKeys,
						},
						"string_value": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  `Holds the value for a local value field with string type.`,
							AtLeastOneOf: localValueKeys,
						},
						"enum_value": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  `The display name of the enum value.`,
							AtLeastOneOf: localValueKeys,
						},
						"duration_value": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `Defines this value as being a Duration.

A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".`,
							AtLeastOneOf: localValueKeys,
						},
						// TODO: String set.
						// TODO: String map
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func resourceSettingFullName(parentType, parentIdentifier, settingName string) string {
	return fmt.Sprintf("%ss/%s/settings/%s", parentType, parentIdentifier, settingName)
}

func resourceSettingShortName(fullName string) string {
	split := strings.Split(fullName, "/")
	return split[len(split)-1]
}

func resourceSettingParentKey(parentType string) string {
	return parentType + "_id"
}

func resourceGoogleResourceSettingCreate(parentType string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		settingName := d.Get("setting_name").(string)
		parentIdentifier := d.Get(resourceSettingParentKey(parentType)).(string)
		id := resourceSettingFullName(parentType, parentIdentifier, settingName)

		if err := patchResourceSetting(d, meta, false, id, parentType); err != nil {
			return fmt.Errorf("Error creating: %s", err)
		}

		d.SetId(id)

		return resourceGoogleResourceSettingRead(parentType)(d, meta)
	}
}

func resourceGoogleResourceSettingRead(parentType string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		get := getResourceSettingFunc(config, userAgent, parentType)

		// Fetch metadata about the setting.
		settingBasic, err := get(d.Id(), "SETTING_VIEW_BASIC")
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("ResourceSetting Not Found : %s", d.Id()))
		}

		// Fetch the localValue field.
		settingLocal, err := get(d.Id(), "SETTING_VIEW_LOCAL_VALUE")
		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("ResourceSetting Not Found : %s", d.Id()))
		}

		if err := d.Set("setting_name", resourceSettingShortName(settingLocal.Name)); err != nil {
			return fmt.Errorf("Error setting setting_name: %s", err)
		}

		localValue := []map[string]interface{}{
			{},
		}
		if lv := settingLocal.LocalValue; lv != nil {
			switch settingBasic.Metadata.DataType {
			// Basic types //
			case "BOOLEAN":
				localValue[0]["boolean_value"] = settingLocal.LocalValue.BooleanValue
			case "STRING":
				localValue[0]["string_value"] = settingLocal.LocalValue.StringValue
			case "ENUM_VALUE":
				localValue[0]["enum_value"] = settingLocal.LocalValue.EnumValue.Value
			case "DURATION_VALUE":
				localValue[0]["duration_value"] = settingLocal.LocalValue.DurationValue

			// Complex types //
			case "STRING_SET":
				// TODO
			case "STRING_MAP":
				// TODO
			}
		}

		if err := d.Set("local_value", localValue); err != nil {
			return fmt.Errorf("Error setting local_value: %s", err)
		}

		return nil
	}
}

func resourceGoogleResourceSettingUpdate(parentType string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		if err := patchResourceSetting(d, meta, false, d.Id(), parentType); err != nil {
			return fmt.Errorf("Error updating: %s", err)
		}

		return nil
	}
}

func resourceGoogleResourceSettingDelete(parentType string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		if err := patchResourceSetting(d, meta, true, d.Id(), parentType); err != nil {
			return fmt.Errorf("Error deleting: %s", err)
		}

		return nil
	}
}

// patchResourceSetting is used by Create/Update/Delete.
// Delete is implemented by setting localValue to nil/null.
func patchResourceSetting(d *schema.ResourceData, meta interface{}, unset bool, id, parentType string) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	var localValue *resourceSettingsV1.GoogleCloudResourcesettingsV1Value
	if !unset {
		localValue = &resourceSettingsV1.GoogleCloudResourcesettingsV1Value{}
		if val, ok := d.GetOk(resourceSettingKeyBooleanValue); ok {
			localValue.BooleanValue = val.(bool)
		} else if val, ok := d.GetOk(resourceSettingKeyStringValue); ok {
			localValue.StringValue = val.(string)
		} else if val, ok := d.GetOk(resourceSettingKeyEnumValue); ok {
			localValue.EnumValue = &resourceSettingsV1.GoogleCloudResourcesettingsV1ValueEnumValue{Value: val.(string)}
		} else if val, ok := d.GetOk(resourceSettingKeyDurationValue); ok {
			localValue.DurationValue = val.(string)
		}
	}

	if _, err := patchResourceSettingFunc(config, userAgent, parentType)(id, &resourceSettingsV1.GoogleCloudResourcesettingsV1Setting{
		Name:       id,
		LocalValue: localValue,
	}); err != nil {
		return fmt.Errorf("patching: %s", err)
	}

	return nil
}

func patchResourceSettingFunc(config *Config, userAgent, parentType string) func(id string, setting *resourceSettingsV1.GoogleCloudResourcesettingsV1Setting) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
	switch parentType {
	case "organization":
		return func(id string, setting *resourceSettingsV1.GoogleCloudResourcesettingsV1Setting) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
			return config.NewResourceSettingsClient(userAgent).Organizations.Settings.Patch(id, setting).Do()
		}
	case "folder":
		return func(id string, setting *resourceSettingsV1.GoogleCloudResourcesettingsV1Setting) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
			return config.NewResourceSettingsClient(userAgent).Folders.Settings.Patch(id, setting).Do()
		}
	case "project":
		return func(id string, setting *resourceSettingsV1.GoogleCloudResourcesettingsV1Setting) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
			return config.NewResourceSettingsClient(userAgent).Projects.Settings.Patch(id, setting).Do()
		}
	default:
		panic("unknown parentType: " + parentType)
	}
}

func getResourceSettingFunc(config *Config, userAgent, parentType string) func(id, view string) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
	switch parentType {
	case "organization":
		return func(id, view string) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
			return config.NewResourceSettingsClient(userAgent).Organizations.Settings.Get(id).View(view).Do()
		}
	case "folder":
		return func(id, view string) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
			return config.NewResourceSettingsClient(userAgent).Folders.Settings.Get(id).View(view).Do()
		}
	case "project":
		return func(id, view string) (*resourceSettingsV1.GoogleCloudResourcesettingsV1Setting, error) {
			return config.NewResourceSettingsClient(userAgent).Projects.Settings.Get(id).View(view).Do()
		}
	default:
		panic("unknown parentType: " + parentType)
	}
}
