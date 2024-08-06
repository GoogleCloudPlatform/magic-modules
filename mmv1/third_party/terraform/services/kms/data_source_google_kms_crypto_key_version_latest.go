package kms

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceGoogleLatestKmsCryptoKeyVersionLatest() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleLatestKmsCryptoKeyVersionLatestRead,
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

func dataSourceGoogleLatestKmsCryptoKeyVersionLatestRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}/cryptoKeyVersions")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting list of cryptoKeyVersions from cryptoKey: {{crypto_key}}")

	cryptoKeyId, err := ParseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   cryptoKeyId.KeyRingId.Project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleDataSourceNotFoundError(err, d, fmt.Sprintf("KmsCryptoKeyVersion %q", d.Id()), url)
	}

	v, err := flattenKmsCryptoKeyVersionLatest(res["cryptoKeyVersions"].([]interface{}), d, config, cryptoKeyId.CryptoKeyId())
	if err != nil {
		return fmt.Errorf("Error getting latest CryptoKeyVersion from crypto key: %s", cryptoKeyId)
	}
	latestVersion := v.(map[string]interface{})

	if err := d.Set("version", latestVersion["version"].(int64)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("name", latestVersion["name"].(string)); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("state", latestVersion["state"]); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("protection_level", latestVersion["protectionLevel"]); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}
	if err := d.Set("algorithm", latestVersion["algorithm"]); err != nil {
		return fmt.Errorf("Error setting CryptoKeyVersion: %s", err)
	}

	url, err = tpgresource.ReplaceVars(d, config, "{{KMSBasePath}}{{crypto_key}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Getting purpose of CryptoKey: %#v", url)
	res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
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

func flattenKmsCryptoKeyVersionVersionLatest(v interface{}, d *schema.ResourceData) interface{} {
	parts := strings.Split(v.(string), "/")
	version := parts[len(parts)-1]
	// Handles the string fixed64 format
	if intVal, err := tpgresource.StringToFixed64(version); err == nil {
		return intVal
	} // let terraform core handle it if we can't convert the string to an int.
	return v
}

func flattenKmsCryptoKeyVersionNameLatest(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionStateLatest(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionProtectionLevelLatest(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionAlgorithmLatest(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionLatest(versionsList []interface{}, d *schema.ResourceData, config *transport_tpg.Config, cryptoKeyId string) (interface{}, error) {
	latestVersion := versionsList[len(versionsList)-1].(map[string]interface{})
	parsedId, err := parseKmsCryptoKeyVersionId(latestVersion["name"].(string), config)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	// The google_kms_crypto_key resource and dataset set
	// id as the value of name (projects/{{project}}/locations/{{location}}/keyRings/{{keyRing}}/cryptoKeys/{{name}})
	// and set name is set as just {{name}}.
	data["id"] = latestVersion["name"]
	data["name"] = parsedId.Name
	data["crypto_key"] = cryptoKeyId

	// fields can be found in `data_source_google_kms_crypto_key_version.go`
	data["version"] = flattenKmsCryptoKeyVersionVersion(latestVersion["name"], d)
	data["algorithm"] = flattenKmsCryptoKeyVersionAlgorithm(latestVersion["algorithm"], d)
	data["protection_level"] = flattenKmsCryptoKeyVersionProtectionLevel(latestVersion["protectionLevel"], d)
	data["state"] = flattenKmsCryptoKeyVersionState(latestVersion["state"], d)
	data["public_key"] = flattenKmsCryptoKeyVersionPublicKey(latestVersion["publicKey"], d)

	return data, nil
}
func flattenKmsCryptoKeyVersionPublicKeyPemLatest(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenKmsCryptoKeyVersionPublicKeyAlgorithmLatest(v interface{}, d *schema.ResourceData) interface{} {
	return v
}
