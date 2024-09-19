package secretmanagerregional

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSecretManagerRegionalRegionalSecretVersionAccess() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretManagerRegionalRegionalSecretVersionAccessRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"secret": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_data": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}
func dataSourceSecretManagerRegionalRegionalSecretVersionAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	secretRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/secrets/(.+)$")
	parts := secretRegex.FindStringSubmatch(d.Get("secret").(string))

	var project string

	// if reference of the secret is provided in the secret field
	if len(parts) == 4 {
		// Stores value of project to set in state
		project = parts[1]
		if d.Get("project").(string) != "" && d.Get("project").(string) != parts[1] {
			return fmt.Errorf("The project set on this secret version (%s) is not equal to the project where this secret exists (%s).", d.Get("project").(string), project)
		}
		if d.Get("location").(string) != "" && d.Get("location").(string) != parts[2] {
			return fmt.Errorf("The location set on this secret version (%s) is not equal to the location where this secret exists (%s).", d.Get("location").(string), parts[2])
		}
		if err := d.Set("location", parts[2]); err != nil {
			return fmt.Errorf("Error setting location: %s", err)
		}
		if err := d.Set("secret", parts[3]); err != nil {
			return fmt.Errorf("Error setting secret: %s", err)
		}
	} else { // if secret name is provided in the secret field
		// Stores value of project to set in state
		project, err = tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("Error fetching project for Secret: %s", err)
		}
		if d.Get("location").(string) == "" {
			return fmt.Errorf("Location must be set when providing only secret name")
		}
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	var url string
	versionNum := d.Get("version")

	// set version if provided, else set version to latest
	if versionNum != "" {
		url, err = tpgresource.ReplaceVars(d, config, "{{SecretManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/{{version}}")
		if err != nil {
			return err
		}
	} else {
		url, err = tpgresource.ReplaceVars(d, config, "{{SecretManagerRegionalBasePath}}projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/latest")
		if err != nil {
			return err
		}
	}

	url = fmt.Sprintf("%s:access", url)
	resp, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})

	if err != nil {
		return fmt.Errorf("Error retrieving available secret manager regional secret version access: %s", err.Error())
	}

	if err := d.Set("name", resp["name"].(string)); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	secretVersionRegex := regexp.MustCompile("projects/(.+)/locations/(.+)/secrets/(.+)/versions/(.+)$")
	parts = secretVersionRegex.FindStringSubmatch(resp["name"].(string))
	if len(parts) != 5 {
		return fmt.Errorf("secret name, %s, does not match format, projects/{{project}}/locations/{{location}}/secrets/{{secret}}/versions/{{version}}", resp["name"].(string))
	}

	log.Printf("[DEBUG] Received Google SecretManager Version: %q", parts[3])

	if err := d.Set("version", parts[4]); err != nil {
		return fmt.Errorf("Error setting version: %s", err)
	}

	data := resp["payload"].(map[string]interface{})
	secretData, err := base64.StdEncoding.DecodeString(data["data"].(string))
	if err != nil {
		return fmt.Errorf("Error decoding secret manager regional secret version data: %s", err.Error())
	}
	if err := d.Set("secret_data", string(secretData)); err != nil {
		return fmt.Errorf("Error setting secret_data: %s", err)
	}

	d.SetId(resp["name"].(string))
	return nil
}
