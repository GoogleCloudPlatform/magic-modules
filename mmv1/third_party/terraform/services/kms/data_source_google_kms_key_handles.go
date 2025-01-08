// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google-beta/google-beta/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func DataSourceGoogleKmsKeyHandles() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsKeyHandlesRead,
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
					The filter argument is used to add a filter query parameter that limits which key handles are retrieved by the data source: ?filter={{filter}}.
					Example values:
					
					* resourceTypeSelector="{SERVICE}.googleapis.com/{TYPE}".
					[See the documentation about using filters](https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations.keyHandles/list)
				`,
			},
			"key_handles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all the retrieved key handles",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kms_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type_selector": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}

}

func dataSourceGoogleKmsKeyHandlesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	params := make(map[string]string)
	if filter, ok := d.GetOk("filter"); ok {
		fmt.Printf("[DEBUG] Search for key handles using filter ?filter=%s", filter)
		params["filter"] = strings.Replace(filter.(string), "\"", "\\\"", -1)
		fmt.Printf("[DEBUG] Search for key handles using filter ?filter=%s", params["filter"])

		if err != nil {
			return err
		}
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	url, err := tpgresource.ReplaceVars(d, config, "https://cloudkms.googleapis.com/v1/projects/{{project}}/locations/{{location}}/keyHandles")
	if err != nil {
		return err
	}
	url = fmt.Sprintf("%s?filter=%s", url, params["filter"])
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	var keyHandles []interface{}
	for {

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              billingProject,
			RawURL:               url,
			UserAgent:            userAgent,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return fmt.Errorf("Error retrieving keyhandles: %s", err)
		}

		if res["keyHandles"] == nil {
			break
		}
		pageKeyHandles, err := flattenKMSKeyHandlesList(config, res["keyHandles"])
		if err != nil {
			return fmt.Errorf("error flattening key handle list: %s", err)
		}
		keyHandles = append(keyHandles, pageKeyHandles...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}
	log.Printf("[DEBUG] Found %d key handles", len(keyHandles))
	if err := d.Set("key_handles", keyHandles); err != nil {
		return fmt.Errorf("error setting key handles: %s", err)
	}
	return nil
}

// flattenKMSKeyHandlesList flattens a list of key handles
func flattenKMSKeyHandlesList(config *transport_tpg.Config, keyHandlesList interface{}) ([]interface{}, error) {
	var keyHandles []interface{}
	for _, k := range keyHandlesList.([]interface{}) {
		keyHandle := k.(map[string]interface{})

		data := map[string]interface{}{}
		// The google_kms_key_handles resource and dataset set
		// id as the value of name (projects/{{project}}/locations/{{location}}/keyHandles/{{name}})
		// and set name is set as just {{name}}.
		fmt.Printf("Info keyhandle %s", keyHandle)
		data["name"] = data["name"]
		data["kms_key"] = keyHandle["kmsKey"]
		data["resource_type_selector"] = keyHandle["resourceTypeSelector"]
		keyHandles = append(keyHandles, data)
	}

	return keyHandles, nil
}
