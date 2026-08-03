package iambeta

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceIAMBetaWorkloadIdentityPoolJwks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMBetaWorkloadIdentityPoolJwksRead,
		Schema: map[string]*schema.Schema{
			"jwks_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The JWKS URI to retrieve the public keys from (e.g. from google_iam_workload_identity_pool_openid_config.jwks_uri).",
			},
			"jwks_json": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The raw JSON string representation of the JWKS response.",
			},
			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The JWKS for this OP.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kty": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key type. Currently \"RSA\".",
						},
						"use": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public key use. Currently \"sig\".",
						},
						"kid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key ID.",
						},
						"n": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modulus value for kty=\"RSA\".",
						},
						"e": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Exponent value for kty=\"RSA\".",
						},
						"alg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Algorithm intended for use with the key. Currently \"RS256\".",
						},
					},
				},
			},
		},
	}
}

func dataSourceIAMBetaWorkloadIdentityPoolJwksRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url := d.Get("jwks_uri").(string)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error fetching JWKS: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error fetching JWKS: HTTP %d", resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("Error parsing JWKS: %s", err)
	}

	rawJson, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("Error marshaling JWKS JSON: %s", err)
	}
	if err := d.Set("jwks_json", string(rawJson)); err != nil {
		return fmt.Errorf("Error setting jwks_json: %s", err)
	}

	var keys []interface{}
	if rawKeys, ok := res["keys"].([]interface{}); ok {
		keys = rawKeys
	}
	if err := d.Set("keys", flattenJwksKeys(keys)); err != nil {
		return fmt.Errorf("Error setting keys: %s", err)
	}

	d.SetId(url)

	return nil
}

func flattenJwksKeys(keys []interface{}) []map[string]interface{} {
	if keys == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0, len(keys))
	for _, rawKey := range keys {
		key, ok := rawKey.(map[string]interface{})
		if !ok {
			continue
		}
		item := make(map[string]interface{})
		if val, ok := key["kty"]; ok {
			item["kty"] = val
		}
		if val, ok := key["alg"]; ok {
			item["alg"] = val
		}
		if val, ok := key["use"]; ok {
			item["use"] = val
		}
		if val, ok := key["kid"]; ok {
			item["kid"] = val
		}
		if val, ok := key["n"]; ok {
			item["n"] = val
		}
		if val, ok := key["e"]; ok {
			item["e"] = val
		}
		result = append(result, item)
	}
	return result
}

func init() {
	registry.Register("google_iam_workload_identity_pool_jwks", DataSourceIAMBetaWorkloadIdentityPoolJwks)
}
