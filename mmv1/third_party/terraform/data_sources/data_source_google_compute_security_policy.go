package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeSecurityPolicy() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceComputeSecurityPolicy().Schema)
	addRequiredFieldsToSchema(dsSchema, "name", "project")

	return &schema.Resource{
		Read:   dataSourceComputSecurityPolicyRead,
		Schema: dsSchema,
	}
}

func dataSourceComputSecurityPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	projectGetCall := config.NewResourceManagerClient(userAgent).Projects.Get(project)

	if config.UserProjectOverride {
		billingProject := project

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}
		projectGetCall.Header().Add("X-Goog-User-Project", billingProject)
	}

	id, err := replaceVars(d, config, "projects/{{project}}/global/securityPolicies/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	sp := d.Get("name").(string)

	client := config.NewComputeClient(userAgent)

	securityPolicy, err := client.SecurityPolicies.Get(project, sp).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SecurityPolicy %q", d.Id()))
	}

	if err := d.Set("name", securityPolicy.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", securityPolicy.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("type", securityPolicy.Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}
	if err := d.Set("rule", flattenSecurityPolicyRules(securityPolicy.Rules)); err != nil {
		return err
	}
	if err := d.Set("fingerprint", securityPolicy.Fingerprint); err != nil {
		return fmt.Errorf("Error setting fingerprint: %s", err)
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("self_link", ConvertSelfLinkToV1(securityPolicy.SelfLink)); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("advanced_options_config", flattenSecurityPolicyAdvancedOptionsConfig(securityPolicy.AdvancedOptionsConfig)); err != nil {
		return fmt.Errorf("Error setting advanced_options_config: %s", err)
	}

	if err := d.Set("adaptive_protection_config", flattenSecurityPolicyAdaptiveProtectionConfig(securityPolicy.AdaptiveProtectionConfig)); err != nil {
		return fmt.Errorf("Error setting adaptive_protection_config: %s", err)
	}

	if err := d.Set("recaptcha_options_config", flattenSecurityPolicyRecaptchaOptionConfig(securityPolicy.RecaptchaOptionsConfig)); err != nil {
		return fmt.Errorf("Error setting recaptcha_options_config: %s", err)
	}

	return nil
}
