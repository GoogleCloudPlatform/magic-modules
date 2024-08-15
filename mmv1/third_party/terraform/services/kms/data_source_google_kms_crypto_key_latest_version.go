package kms

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleKmsLatestCryptoKeyVersion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsLatestCryptoKeyVersionRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protection_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pem": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleKmsLatestCryptoKeyVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)

	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	cryptoKeyId, err := ParseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s/latestCryptoKeyVersion", cryptoKeyId.CryptoKeyId())
	d.SetId(id)

	versions, err := dataSourceKMSCryptoKeyVersionsList(d, meta, cryptoKeyId.CryptoKeyId(), userAgent)
	if err != nil {
		return err
	}

	// grab latest version
	lv := len(versions) - 1
	if lv < 0 {
		return fmt.Errorf("No CryptoVersions found in crypto key %s", cryptoKeyId.CryptoKeyId())
	}

	latestVersion := versions[lv].(map[string]interface{})

	// The google_kms_crypto_key resource and dataset set
	// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
	// and set name is set as just {{name}}.

	if err := d.Set("name", flattenKmsCryptoKeyVersionName(latestVersion["name"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}
	if err := d.Set("version", flattenKmsCryptoKeyVersionVersion(latestVersion["name"], d)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("state", flattenKmsCryptoKeyVersionState(latestVersion["state"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}
	if err := d.Set("protection_level", flattenKmsCryptoKeyVersionProtectionLevel(latestVersion["protectionLevel"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}
	if err := d.Set("algorithm", flattenKmsCryptoKeyVersionAlgorithm(latestVersion["algorithm"], d)); err != nil {
		return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting attributes for CryptoKeyVersion: %#v", url)

	url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting purpose of CryptoKey: %#v", url)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   cryptoKeyId.KeyRingId.Project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("KmsCryptoKey %q", d.Id()), url)
	}

	if res["purpose"] == "ASYMMETRIC_SIGN" || res["purpose"] == "ASYMMETRIC_DECRYPT" {
		url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions/{{version}}/publicKey")
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Getting public key of CryptoKeyVersion: %#v", url)

		res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:               config,
			Method:               "GET",
			Project:              cryptoKeyId.KeyRingId.Project,
			RawURL:               url,
			UserAgent:            userAgent,
			Timeout:              d.Timeout(schema.TimeoutRead),
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsCryptoKeyVersionsPendingGeneration},
		})

		if err != nil {
			log.Printf("Error generating public key: %s", err)
			return err
		}

		if err := d.Set("public_key", flattenKmsCryptoKeyVersionPublicKey(res, d)); err != nil {
			return fmt.Errorf("Error setting CryptoKeyVersion public key: %s", err)
		}
	}
	d.SetId(fmt.Sprintf("//cloudkms.googleapis.com/v1/%s/cryptoKeyVersions/%d", d.Get("crypto_key"), d.Get("version")))

	return nil
}

// func flattenKMSCryptoKeyLatestVersion(d *schema.ResourceData, meta interface{}, version interface{}, cryptoKeyId string) (map[string]interface{}, error) {

// 		latestVersion := version.(map[string]interface{})

// 		// The google_kms_crypto_key resource and dataset set
// 		// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
// 		// and set name is set as just {{name}}.

// 		if err := d.Set("name", flattenKmsCryptoKeyVersionName(latestVersion["name"], d)){
// 			return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
// 		}
// 		if err := d.Set("version", flattenKmsCryptoKeyVersionVersion(latestVersion["name"], d)); err != nil {
// 			return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
// 		}
// 		if err := d.Set("state",  flattenKmsCryptoKeyVersionState(latestVersion["state"], d)); err != nil {
// 			return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
// 		}
// 		if err := d.Set("protection_level", flattenKmsCryptoKeyVersionProtectionLevel(latestVersion["protectionLevel"], d)); err != nil {
// 			return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
// 		}
// 		if err := d.Set("algorithm", flattenKmsCryptoKeyVersionAlgorithm(latestVersion["algorithm"], d)); err != nil {
// 			return fmt.Errorf("Error setting LatestCryptoKeyVersion: %s", err)
// 		}

// 	return versions, nil
// }

// func flattenKmsCryptoKeyVersionVersion(v interface{}, d *schema.ResourceData) interface{} {
// 	parts := strings.Split(v.(string), "/")
// 	version := parts[len(parts)-1]
// 	// Handles the string fixed64 format
// 	if intVal, err := tpgresource.StringToFixed64(version); err == nil {
// 		return intVal
// 	} // let terraform core handle it if we can't convert the string to an int.
// 	return v
// }

// func flattenKmsCryptoKeyVersionName(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionState(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionProtectionLevel(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionAlgorithm(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionPublicKey(v interface{}, d *schema.ResourceData) interface{} {
// 	if v == nil {
// 		return nil
// 	}
// 	original := v.(map[string]interface{})
// 	if len(original) == 0 {
// 		return nil
// 	}
// 	transformed := make(map[string]interface{})
// 	transformed["pem"] =
// 		flattenKmsCryptoKeyVersionPublicKeyPem(original["pem"], d)
// 	transformed["algorithm"] =
// 		flattenKmsCryptoKeyVersionPublicKeyAlgorithm(original["algorithm"], d)
// 	return []interface{}{transformed}
// }
// func flattenKmsCryptoKeyVersionPublicKeyPem(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }

// func flattenKmsCryptoKeyVersionPublicKeyAlgorithm(v interface{}, d *schema.ResourceData) interface{} {
// 	return v
// }
