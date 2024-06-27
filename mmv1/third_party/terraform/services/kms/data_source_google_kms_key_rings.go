package kms

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleKmsKeyRings() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsKeyRingsRead,
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
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `
					The filter argument is used to add a filter query parameter that limits which keys are retrieved by the data source: ?filter={{filter}}.
					Example values:
					
					* "name:my-key-" will retrieve key rings that contain "my-key-" anywhere in their name. Note: names take the form projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}.
					* "name=projects/my-project/locations/global/keyRings/my-key-ring" will only retrieve a key ring with that exact name.
					
					[See the documentation about using filters](https://cloud.google.com/kms/docs/sorting-and-filtering)
				`,
			},
			"key_rings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all the retrieved key rings",
				Elem: &schema.Resource{
					// schema isn't used from resource_kms_key_ring due to having project and location fields which are empty when grabbed in a list.
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleKmsKeyRingsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/locations/{{location}}/keyRings")
	if err != nil {
		return err
	}
	d.SetId(id)

	log.Printf("[DEBUG] Searching for keyrings")
	res, err := dataSourceGoogleKmsKeyRingsList(d, meta)
	if err != nil {
		return err
	}

	// Check for keyRings field, as empty response lacks keyRings
	// If found, set data in the data source's `keyRing` field
	if keyRings, ok := res["keyRings"].([]interface{}); ok {
		log.Printf("[DEBUG] Found %d key rings", len(keyRings))
		value, err := flattenKMSKeyRingsList(d, config, keyRings)
		if err != nil {
			return fmt.Errorf("error flattening key rings list: %s", err)
		}
		if err := d.Set("key_rings", value); err != nil {
			return fmt.Errorf("error setting key rings: %s", err)
		}
	} else {
		log.Printf("[DEBUG] Found 0 key rings")
	}

	return nil
}

func dataSourceGoogleKmsKeyRingsList(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}projects/{{project}}/locations/{{location}}/keyRings")
	if err != nil {
		return nil, err
	}

	if filter, ok := d.GetOk("filter"); ok {
		log.Printf("[DEBUG] Search for key rings using filter ?filter=%s", filter.(string))
		url, err = transport_tpg.AddQueryParams(url, map[string]string{"filter": filter.(string)})
		if err != nil {
			return nil, err
		}
	}

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, fmt.Errorf("Error fetching project for keyRings: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	headers := make(http.Header)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
		Headers:   headers,
	})
	if err != nil {
		return nil, transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("KMSKeyRing %q", d.Id()))
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing KMSKeyRing because it no longer exists.")
		d.SetId("")
		return nil, nil
	}
	return res, nil
}

// flattenKMSKeyRingsList flattens a list of key rings
func flattenKMSKeyRingsList(d *schema.ResourceData, config *transport_tpg.Config, keyRingsList []interface{}) ([]interface{}, error) {
	var keyRings []interface{}
	for _, k := range keyRingsList {
		keyRing := k.(map[string]interface{})

		parsedId, err := parseKmsKeyRingId(keyRing["name"].(string), config)
		if err != nil {
			return nil, err
		}

		data := map[string]interface{}{}
		// The google_kms_key_rings resource and dataset set
		// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{name}})
		// and set name is set as just {{name}}.
		data["id"] = keyRing["name"]
		data["name"] = parsedId.Name

		keyRings = append(keyRings, data)
	}

	return keyRings, nil
}
