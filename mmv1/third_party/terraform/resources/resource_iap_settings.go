package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIapSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceIapSettingsCreate,
		Read:   resourceIapSettingsRead,
		Delete: resourceIapSettingsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceIapSettingsImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"backend_service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Identifies the backend service that this IAP settings should apply to.`,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: `Access related settings for IAP protected apps.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gcip_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `GCIP claims and endpoint configurations for 3p identity providers.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tenant_ids": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: `GCIP tenant ids that are linked to the IAP resource. tenantIds could be a string beginning with a number character to indicate authenticating with GCIP tenant flow, or in the format of _ to indicate authenticating with GCIP agent flow. If agent flow is used, tenantIds should only contain one single element, while for tenant flow, tenantIds can contain multiple elements.`,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											// TODO validate https URL
										},
									},
									"login_page_uri": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Login page URI associated with the GCIP tenants. Typically, all resources within the same project share the same login page, though it could be overridden at the sub resource level`,
									},
								},
							},
						},
						"cors_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Allows customers to configure HTTP request paths that'll allow HTTP OPTIONS call to bypass authentication and authorization.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_http_options": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Configuration to allow HTTP OPTIONS calls to skip authorization. If undefined, IAP will not apply any special logic to OPTIONS requests.`,
									},
								},
							},
						},
						"oauth_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Settings to configure IAP's OAuth behavior.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"login_hint": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Domain hint to send as hd=? parameter in OAuth request flow. Enables redirect to primary IDP by skipping Google's login screen. https://developers.google.com/identity/protocols/OpenIDConnect#hd-param Note: IAP does not verify that the id token's hd claim matches this value since access behavior is managed by IAM policies.`,
									},
								},
							},
						},
						"reauth_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Settings to configure reauthentication policies in IAP.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"method": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  `Reauth method required by the policy. Possible values: ["METHOD_UNSPECIFIED", "LOGIN", "PASSWORD", "SECURE_KEY"]`,
										ValidateFunc: validateEnum([]string{"METHOD_UNSPECIFIED", "LOGIN", "PASSWORD", "SECURE_KEY"}),
									},
									"max_age": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Reauth session lifetime, how long before a user has to reauthenticate again. A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".`,
									},
									"policy_type": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  `Reauth method required by the policy. Possible values: ["POLICY_TYPE_UNSPECIFIED", "MINIMUM", "DEFAULT"]`,
										ValidateFunc: validateEnum([]string{"POLICY_TYPE_UNSPECIFIED", "MINIMUM", "DEFAULT"}),
									},
								},
							},
						},
					},
				},
			},
			"application_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: `Wrapper over application specific settings for IAP.`,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asm_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Configuration for RCTokens generated for ASM workloads protected by IAP. RCTokens are IAP generated JWTs that can be verified at the application. The RCToken is primarily used for ISTIO deployments, and can be scoped to a single mesh by configuring the audience field accordingly`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rctoken_aud": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Audience claim set in the generated RCToken. This value is not validated by IAP.`,
									},
								},
							},
						},
						"access_denied_page_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Custom content configuration for access denied page. IAP allows customers to define a custom URI to use as the error page when access is denied to users. If IAP prevents access to this page, the default IAP error page will be displayed instead.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_denied_page_uri": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The URI to be redirected to when access is denied.`,
									},
									"generate_troubleshooting_uri": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Whether to generate a troubleshooting URL on access denied events to this application.`,
									},
								},
							},
						},
						"cookie_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The Domain value to set for cookies generated by IAP. This value is not validated by the API, but will be ignored at runtime if invalid.`,
						},
					},
				},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceIapSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	accessSettingsProp, err := expandIapSettingsAccessSettings(d.Get("access_settings"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("access_settings"); !isEmptyValue(reflect.ValueOf(accessSettingsProp)) && (ok || !reflect.DeepEqual(v, accessSettingsProp)) {
		obj["accessSettings"] = accessSettingsProp
	}

	applicationSettingsProp, err := expandIapSettingsApplicationSettings(d.Get("application_settings"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("application_settings"); !isEmptyValue(reflect.ValueOf(applicationSettingsProp)) && (ok || !reflect.DeepEqual(v, applicationSettingsProp)) {
		obj["applicationSettings"] = applicationSettingsProp
	}

	url, err := replaceVars(d, config, "{{IapBasePath}}projects/{{project}}/iap_web/compute/services/{{backend_service_id}}:iapSettings")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new IAP Settings: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for IAP Settings: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating IAP Settings: %s", err)
	}
	if err := d.Set("name", flattenIapSettingsName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// `name` is autogenerated from the api so needs to be set post-create
	name, ok := res["name"]
	if !ok {
		respBody, ok := res["response"]
		if !ok {
			return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
		}

		name, ok = respBody.(map[string]interface{})["name"]
		if !ok {
			return fmt.Errorf("Create response didn't contain critical fields. Create may not have succeeded.")
		}
	}
	if err := d.Set("name", name.(string)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	d.SetId(name.(string))

	err = PollingWaitTime(resourceIapSettingsPollRead(d, meta), PollCheckForExistence, "Creating Settings", d.Timeout(schema.TimeoutCreate), 5)
	if err != nil {
		return fmt.Errorf("Error waiting to create Settings: %s", err)
	}

	log.Printf("[DEBUG] Finished creating Settings %q: %#v", d.Id(), res)

	return resourceIapSettingsRead(d, meta)
}

func resourceIapSettingsPollRead(d *schema.ResourceData, meta interface{}) PollReadFunc {
	return func() (map[string]interface{}, error) {
		config := meta.(*Config)

		url, err := replaceVars(d, config, "{{IapBasePath}}{{name}}:iapSettings")
		if err != nil {
			return nil, err
		}

		billingProject := ""

		project, err := getProject(d, config)
		if err != nil {
			return nil, fmt.Errorf("Error fetching project for IAP Settings: %s", err)
		}
		billingProject = project

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return nil, err
		}

		res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
		if err != nil {
			return res, err
		}
		return res, nil
	}
}

func resourceIapSettingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{IapBasePath}}{{name}}:iapSettings")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for IAP Settings: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := sendRequest(config, "GET", billingProject, url, userAgent, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("IapSettings %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading IAP Settings: %s", err)
	}

	if err := d.Set("access_settings", flattenIapSettingsAccessSettings(res["accessSettings"], d, config)); err != nil {
		return fmt.Errorf("Error reading IAP Access Settings: %s", err)
	}

	if err := d.Set("application_settings", flattenIapSettingsApplicationSettings(res["applicationSettings"], d, config)); err != nil {
		return fmt.Errorf("Error reading IAP Application Settings: %s", err)
	}

	if err := d.Set("name", flattenIapSettingsName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading IAP Settings: %s", err)
	}

	return nil
}

func resourceIapSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARNING] Iap Settings resources"+
		" cannot be deleted from Google Cloud. The resource %s will be removed from Terraform"+
		" state, but will still be present on Google Cloud.", d.Id())
	d.SetId("")

	return nil
}

func resourceIapSettingsImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	// current import_formats can't import fields with forward slashes in their value
	if err := parseImportId([]string{"(?P<name>.+)"}, d, config); err != nil {
		return nil, err
	}

	nameParts := strings.Split(d.Get("name").(string), "/")
	if len(nameParts) != 6 {
		return nil, fmt.Errorf(
			"Saw %s when the name is expected to have shape %s",
			d.Get("name"),
			"projects/{{project}}/iap_web/compute/services/{{backend_service_id}}",
		)
	}

	if err := d.Set("project", nameParts[1]); err != nil {
		return nil, fmt.Errorf("Error setting project: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func flattenIapSettingsAccessSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["gcip_settings"] =
		flattenIapSettingsAccessSettingsGcipSettings(original["gcipSettings"], d, config)
	transformed["cors_settings"] =
		flattenIapSettingsAccessSettingsCorsSettings(original["corsSettings"], d, config)
	transformed["oauth_settings"] =
		flattenIapSettingsAccessSettingsOauthSettings(original["oauthSettings"], d, config)
	transformed["reauth_settings"] =
		flattenIapSettingsAccessSettingsReauthSettings(original["reauthSettings"], d, config)
	return []interface{}{transformed}
}

func flattenIapSettingsApplicationSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["csm_settings"] =
		flattenIapSettingsApplicationSettingsCsmSettings(original["csmSettings"], d, config)
	transformed["access_denied_page_settings"] =
		flattenIapSettingsApplicationSettingsAccessDeniedPageSettings(original["accessDeniedPageSettings"], d, config)
	transformed["cookie_domain"] =
		flattenIapSettingsApplicationSettingsCookieDomain(original["cookieDomain"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsApplicationSettingsCsmSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["rctoken_aud"] =
		flattenIapSettingsApplicationSettingsCsmSettingsRctokenAud(original["rctokenAud"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsApplicationSettingsCsmSettingsRctokenAud(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsApplicationSettingsAccessDeniedPageSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["access_denied_page_uri"] =
		flattenIapSettingsApplicationSettingsAccessDeniedPageSettingsAccessDeniedPageUri(original["accessDeniedPageUri"], d, config)
	transformed["generate_troubleshooting_uri"] =
		flattenIapSettingsApplicationSettingsAccessDeniedPageSettingsGenerateTroubleshootingUri(original["generateTroubleshootingUri"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsApplicationSettingsAccessDeniedPageSettingsAccessDeniedPageUri(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsApplicationSettingsAccessDeniedPageSettingsGenerateTroubleshootingUri(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsApplicationSettingsCookieDomain(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsReauthSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["method"] =
		flattenIapSettingsAccessSettingsReauthSettingsMethod(original["method"], d, config)
	transformed["max_age"] =
		flattenIapSettingsAccessSettingsReauthSettingsMaxAge(original["maxAge"], d, config)
	transformed["policy_type"] =
		flattenIapSettingsAccessSettingsReauthSettingsPolicyType(original["policyType"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsAccessSettingsReauthSettingsMethod(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsReauthSettingsMaxAge(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsReauthSettingsPolicyType(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsOauthSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["login_hint"] =
		flattenIapSettingsAccessSettingsOauthSettingsLoginHint(original["loginHint"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsAccessSettingsOauthSettingsLoginHint(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsCorsSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["allow_http_options"] =
		flattenIapSettingsAccessSettingsGcipSettingsAllowHttpOptions(original["allowHttpOptions"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsAccessSettingsGcipSettingsAllowHttpOptions(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsGcipSettings(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["tenant_ids"] =
		flattenIapSettingsAccessSettingsGcipSettingsTenantIds(original["tenantIds"], d, config)
	transformed["login_page_uri"] =
		flattenIapSettingsAccessSettingsGcipSettingsLoginPageUri(original["loginPageUri"], d, config)

	return []interface{}{transformed}
}

func flattenIapSettingsAccessSettingsGcipSettingsTenantIds(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsAccessSettingsGcipSettingsLoginPageUri(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenIapSettingsName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandIapSettingsAccessSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedGcipSettings, err := expandIapSettingsGcipSettings(original["gcip_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGcipSettings); val.IsValid() && !isEmptyValue(val) {
		transformed["gcipSettings"] = transformedGcipSettings
	}

	transformedCorsSettings, err := expandIapSettingsCorsSettings(original["cors_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCorsSettings); val.IsValid() && !isEmptyValue(val) {
		transformed["corsSettings"] = transformedCorsSettings
	}

	transformedOauthSettings, err := expandIapSettingsOauthSettings(original["oauth_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedOauthSettings); val.IsValid() && !isEmptyValue(val) {
		transformed["oauthSettings"] = transformedOauthSettings
	}

	transformedReauthSettings, err := expandIapSettingsReauthSettings(original["reauth_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedReauthSettings); val.IsValid() && !isEmptyValue(val) {
		transformed["reauthSettings"] = transformedReauthSettings
	}

	return transformed, nil
}

func expandIapSettingsGcipSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedTenantIds, err := expandIapSettingsTenantIds(original["tenant_ids"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTenantIds); val.IsValid() && !isEmptyValue(val) {
		transformed["tenantIds"] = transformedTenantIds
	}

	transformedLoginPageUri, err := expandIapSettingsLoginPageUri(original["login_page_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedLoginPageUri); val.IsValid() && !isEmptyValue(val) {
		transformed["loginPageUri"] = transformedLoginPageUri
	}

	return transformed, nil
}

func expandIapSettingsTenantIds(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsLoginPageUri(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsCorsSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowHttpOptions, err := expandIapSettingsAllowHttpOptions(original["allow_http_options"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowHttpOptions); val.IsValid() && !isEmptyValue(val) {
		transformed["allowHttpOptions"] = transformedAllowHttpOptions
	}

	return transformed, nil
}

func expandIapSettingsAllowHttpOptions(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsOauthSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowLoginHint, err := expandIapSettingsAllowLoginHint(original["login_hint"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowLoginHint); val.IsValid() && !isEmptyValue(val) {
		transformed["loginHint"] = transformedAllowLoginHint
	}

	return transformed, nil
}

func expandIapSettingsAllowLoginHint(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsReauthSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAllowMethod, err := expandIapSettingsAllowMethod(original["method"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAllowMethod); val.IsValid() && !isEmptyValue(val) {
		transformed["method"] = transformedAllowMethod
	}

	transformedMaxAge, err := expandIapSettingsAllowMaxAge(original["max_age"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedMaxAge); val.IsValid() && !isEmptyValue(val) {
		transformed["maxAge"] = transformedMaxAge
	}

	transformedPolicyType, err := expandIapSettingsPolicyType(original["policy_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedPolicyType); val.IsValid() && !isEmptyValue(val) {
		transformed["policyType"] = transformedPolicyType
	}

	return transformed, nil
}

func expandIapSettingsAllowMethod(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsAllowMaxAge(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsPolicyType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsApplicationSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedCsmSettings, err := expandIapSettingsCsmSettings(original["csm_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCsmSettings); val.IsValid() && !isEmptyValue(val) {
		transformed["csmSettings"] = transformedCsmSettings
	}

	transformedCookieDomain, err := expandIapSettingsCookieDomain(original["cookie_domain"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedCookieDomain); val.IsValid() && !isEmptyValue(val) {
		transformed["cookieDomain"] = transformedCookieDomain
	}

	transformedAccessDeniedPageSettings, err := expandIapSettingsAccessDeniedPageSettings(original["access_denied_page_settings"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAccessDeniedPageSettings); val.IsValid() && !isEmptyValue(val) {
		transformed["accessDeniedPageSettings"] = transformedAccessDeniedPageSettings
	}

	return transformed, nil
}

func expandIapSettingsCsmSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedRctokenAud, err := expandIapSettingsRctokenAud(original["rctoken_aud"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedRctokenAud); val.IsValid() && !isEmptyValue(val) {
		transformed["rctokenAud"] = transformedRctokenAud
	}

	return transformed, nil
}

func expandIapSettingsRctokenAud(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsAccessDeniedPageSettings(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedAccessDeniedPageUri, err := expandIapSettingsAccessDeniedPageUri(original["access_denied_page_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedAccessDeniedPageUri); val.IsValid() && !isEmptyValue(val) {
		transformed["accessDeniedPageUri"] = transformedAccessDeniedPageUri
	}

	transformedGenerateTroubleshootingUri, err := expandIapSettingsGenerateTroubleshootingUri(original["generate_troubleshooting_uri"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedGenerateTroubleshootingUri); val.IsValid() && !isEmptyValue(val) {
		transformed["generateTroubleshootingUri"] = transformedGenerateTroubleshootingUri
	}

	return transformed, nil
}

func expandIapSettingsAccessDeniedPageUri(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsGenerateTroubleshootingUri(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandIapSettingsCookieDomain(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
