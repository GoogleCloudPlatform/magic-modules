package iambeta

import (
	"encoding/json"
	"fmt"
	"io"
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
			"resource_name": {
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
				Description: "The list of public keys in the JSON Web Key Set.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kty": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key type.",
						},
						"use": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The intended use of the public key.",
						},
						"kid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for the key.",
						},
						"n": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The modulus for the RSA public key.",
						},
						"e": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The exponent for the RSA public key.",
						},
						"alg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The algorithm intended for use with the key.",
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

	url := d.Get("resource_name").(string)

	// We cannot use standard provider transport (transport_tpg.SendRequest) here because
	// the JWKS endpoint (/openid/jwks) is a public, unauthenticated API.
	// If the provider transport is used, it attaches an OAuth Bearer token to the request which
	// causes Google STS to reject the request with HTTP 400 Bad Request during cross-org testing (e.g. VCR).
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading JWKS response: %s", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return fmt.Errorf("Error parsing JWKS response: %s", err)
	}

	if err := d.Set("jwks_json", string(bodyBytes)); err != nil {
		return fmt.Errorf("Error setting jwks_json: %s", err)
	}

	if keysRaw, ok := res["keys"].([]interface{}); ok {
		keysList := make([]map[string]interface{}, 0, len(keysRaw))
		for _, keyRaw := range keysRaw {
			if keyMap, ok := keyRaw.(map[string]interface{}); ok {
				k := map[string]interface{}{
					"kty": keyMap["kty"],
					"use": keyMap["use"],
					"kid": keyMap["kid"],
					"n":   keyMap["n"],
					"e":   keyMap["e"],
					"alg": keyMap["alg"],
				}
				keysList = append(keysList, k)
			}
		}
		if err := d.Set("keys", keysList); err != nil {
			return fmt.Errorf("Error setting keys: %s", err)
		}
	}

	d.SetId(url)

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_iam_workload_identity_pool_jwks",
		ProductName: "iambeta",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceIAMBetaWorkloadIdentityPoolJwks(),
	}.Register()
}
