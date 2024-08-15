// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package kms

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleKmsCryptoKeyVersions() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceKMSCryptoKeyVersion().Schema)

	dsSchema["id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Read: dataSourceGoogleKmsCryptoKeyVersionsRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all the retrieved cryptoKeyVersions from the provided crypto key",
				Elem: &schema.Resource{
					Schema: dsSchema,
				},
			},
		},
	}
}

func dataSourceGoogleKmsCryptoKeyVersionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	cryptoKeyId, err := ParseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s/cryptoKeyVersions", cryptoKeyId.CryptoKeyId())
	d.SetId(id)

	log.Printf("[DEBUG] Searching for cryptoKeyVersions in crypto key %s", cryptoKeyId.CryptoKeyId())
	versions, err := dataSourceKMSCryptoKeyVersionsList(d, meta, cryptoKeyId.CryptoKeyId())
	if err != nil {
		return err
	}

	if len(versions) > 0 {
		log.Printf("[DEBUG] Found %d cryptoKeyVersions in crypto key %s", len(versions), cryptoKeyId.CryptoKeyId())
		value, err := flattenKMSCryptoKeyVersionsList(d, config, versions, cryptoKeyId.CryptoKeyId())
		if err != nil {
			return fmt.Errorf("error flattening cryptoKeyVersions list: %s", err)
		}
		if err := d.Set("versions", value); err != nil {
			return fmt.Errorf("error setting versions: %s", err)
		}
	} else {
		log.Printf("[DEBUG] Found 0 versions in crypto key %s", cryptoKeyId.CryptoKeyId())
	}
	// this is meant for latest_version
	// url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}")
	// if err != nil {
	// 	return err
	// }

	// log.Printf("[DEBUG] Getting purpose of CryptoKey: %#v", url)
	// res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
	// 	Config:    config,
	// 	Method:    "GET",
	// 	Project:   cryptoKeyId.KeyRingId.Project,
	// 	RawURL:    url,
	// 	UserAgent: userAgent,
	// })
	// if err != nil {
	// 	return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("KmsCryptoKey %q", d.Id()), url)
	// }

	// if res["purpose"] == "ASYMMETRIC_SIGN" || res["purpose"] == "ASYMMETRIC_DECRYPT" {
	// 	url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}/publicKey")
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Printf("[DEBUG] Getting public key of CryptoKeyVersion: %#v", url)

	// 	res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
	// 		Config:               config,
	// 		Method:               "GET",
	// 		Project:              cryptoKeyId.KeyRingId.Project,
	// 		RawURL:               url,
	// 		UserAgent:            userAgent,
	// 		Timeout:              d.Timeout(schema.TimeoutRead),
	// 		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsCryptoKeyVersionsPendingGeneration},
	// 	})

	// 	if err != nil {
	// 		log.Printf("Error generating public key: %s", err)
	// 		return err
	// 	}

	// 	if err := d.Set("public_key", flattenKmsCryptoKeyVersionPublicKey(res, d)); err != nil {
	// 		return fmt.Errorf("Error setting CryptoKeyVersion public key: %s", err)
	// 	}
	// }
	// d.SetId(fmt.Sprintf("//cloudkms.googleapis.com/v1/%s/cryptoKeyVersions/%d", d.Get("crypto_key"), d.Get("version")))

	return nil
}

func dataSourceKMSCryptoKeyVersionsList(d *schema.ResourceData, meta interface{}, cryptoKeyId string) ([]interface{}, error) {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return nil, err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions?filter=state=ENABLED")
	if err != nil {
		return nil, err
	}

	billingProject := ""

	if parts := regexp.MustCompile(`projects\/([^\/]+)\/`).FindStringSubmatch(url); parts != nil {
		billingProject = parts[1]
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// Always include the filter param, and optionally include the pageToken parameter for subsequent requests
	var params = make(map[string]string, 0)
	if filter, ok := d.GetOk("filter"); ok {
		log.Printf("[DEBUG] Search for cryptoKeyVersions in crypto key %s is using filter ?filter=%s", cryptoKeyId, filter.(string))
		params["filter"] = filter.(string)
	}

	cryptoKeyVersions := make([]interface{}, 0)
	for {
		// Depending on previous iterations, params might contain a pageToken param
		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return nil, err
		}

		headers := make(http.Header)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
			Headers:   headers,
			// ErrorRetryPredicates used to allow retrying if rate limits are hit when requesting multiple pages in a row
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429RetryableQuotaError},
		})
		if err != nil {
			return nil, transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("KMSCryptoKeyVersions %q", d.Id()))
		}

		if res == nil {
			// Decoding the object has resulted in it being gone. It may be marked deleted
			log.Printf("[DEBUG] Removing KMSCryptoKeyVersion because it no longer exists.")
			d.SetId("")
			return nil, nil
		}

		// Store info from this page
		if v, ok := res["cryptoKeyVersions"].([]interface{}); ok {
			cryptoKeyVersions = append(cryptoKeyVersions, v...)
		}

		// Handle pagination for next loop, or break loop
		v, ok := res["nextPageToken"]
		if ok {
			params["pageToken"] = v.(string)
		}
		if !ok {
			break
		}
	}

	return cryptoKeyVersions, nil
}

func flattenKMSCryptoKeyVersionsList(d *schema.ResourceData, meta interface{}, versionsList []interface{}, cryptoKeyId string) ([]interface{}, error) {
	var versions []interface{}
	for _, v := range versionsList {
		version := v.(map[string]interface{})

		data := map[string]interface{}{}
		// The google_kms_crypto_key resource and dataset set
		// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
		// and set name is set as just {{name}}.
		data["id"] = version["name"]
		data["name"] = flattenKmsCryptoKeyVersionName(version["name"], d)
		data["crypto_key"] = cryptoKeyId

		data["state"] = flattenKmsCryptoKeyVersionState(version["state"], d)
		data["protection_level"] = flattenKmsCryptoKeyVersionProtectionLevel(version["protectionLevel"], d)
		data["algorithm"] = flattenKmsCryptoKeyVersionAlgorithm(version["algorithm"], d)

		versions = append(versions, data)
	}

	return versions, nil
}

// func flattenKmsCryptoKeyVersionNameLatest(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionStateLatest(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionProtectionLevelLatest(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionAlgorithmLatest(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionLatest(versionsList []interface{}, d *schema.ResourceData, config *transport_tpg.Config, cryptoKeyId string) (interface{}, error) {
// 	latestVersion := versionsList[len(versionsList)-1].(map[string]interface{})
// 	parsedId, err := parseKmsCryptoKeyVersionId(latestVersion["name"].(string), config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	data := map[string]interface{}{}
// 	// The google_kms_crypto_key resource and dataset set
// 	// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
// 	// and set name is set as just {{name}}.
// 	data["id"] = latestVersion["name"]
// 	data["name"] = parsedId.Name
// 	data["crypto_key"] = cryptoKeyId

// 	// fields can be found in `data_source_google_kms_crypto_key_version.go`
// 	data["version"] = flattenKmsCryptoKeyVersionVersion(latestVersion["name"], d)
// 	data["algorithm"] = flattenKmsCryptoKeyVersionAlgorithm(latestVersion["algorithm"], d)
// 	data["protection_level"] = flattenKmsCryptoKeyVersionProtectionLevel(latestVersion["protectionLevel"], d)
// 	data["state"] = flattenKmsCryptoKeyVersionState(latestVersion["state"], d)
// 	data["public_key"] = flattenKmsCryptoKeyVersionPublicKey(latestVersion["publicKey"], d)

// 	return data, nil
// }
// func flattenKmsCryptoKeyVersionPublicKeyPemLatest(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionPublicKeyAlgorithmLatest(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }
