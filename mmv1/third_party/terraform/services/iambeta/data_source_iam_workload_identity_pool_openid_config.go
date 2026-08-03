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

func DataSourceIAMBetaWorkloadIdentityPoolOpenIdConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIAMBetaWorkloadIdentityPoolOpenIdConfigRead,
		Schema: map[string]*schema.Schema{
			"workload_identity_pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the pool whose OpenID provider configuration to retrieve.",
			},
			"issuer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL using the https scheme with no query or fragment components that the OP asserts as its issuer identifier.",
			},
			"jwks_uri": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the OP's JWK Set [JWK] document, which MUST use the https scheme.",
			},
			"authorization_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL pointing to an authorization endpoint under this issuer.",
			},
			"token_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL pointing to a token endpoint under this issuer.",
			},
			"response_types_supported": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "JSON array containing a list of the OAuth 2.0 response_type values that this OP supports.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"subject_types_supported": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "JSON array containing a list of the subject identifier types that this OP supports.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"id_token_signing_alg_values_supported": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "JSON array containing a list of the JWS signing algorithms (alg values) supported by the OP for the ID token to encode the claims in a JWT [JWT].",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIAMBetaWorkloadIdentityPoolOpenIdConfigRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	poolName := d.Get("workload_identity_pool_id").(string)
	url := fmt.Sprintf("https://sts.googleapis.com/v1/%s/.well-known/openid-configuration?alt=json", poolName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error fetching OpenID configuration: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error fetching OpenID configuration: HTTP %d", resp.StatusCode)
	}

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("Error parsing OpenID configuration: %s", err)
	}

	if err := d.Set("issuer", res["issuer"]); err != nil {
		return fmt.Errorf("Error setting issuer: %s", err)
	}
	if err := d.Set("jwks_uri", res["jwks_uri"]); err != nil {
		return fmt.Errorf("Error setting jwks_uri: %s", err)
	}
	if err := d.Set("authorization_endpoint", res["authorization_endpoint"]); err != nil {
		return fmt.Errorf("Error setting authorization_endpoint: %s", err)
	}
	if err := d.Set("token_endpoint", res["token_endpoint"]); err != nil {
		return fmt.Errorf("Error setting token_endpoint: %s", err)
	}
	if err := d.Set("response_types_supported", flattenStringList(res["response_types_supported"])); err != nil {
		return fmt.Errorf("Error setting response_types_supported: %s", err)
	}
	if err := d.Set("subject_types_supported", flattenStringList(res["subject_types_supported"])); err != nil {
		return fmt.Errorf("Error setting subject_types_supported: %s", err)
	}
	if err := d.Set("id_token_signing_alg_values_supported", flattenStringList(res["id_token_signing_alg_values_supported"])); err != nil {
		return fmt.Errorf("Error setting id_token_signing_alg_values_supported: %s", err)
	}

	d.SetId(poolName)

	return nil
}

func flattenStringList(val interface{}) []string {
	if val == nil {
		return nil
	}
	if list, ok := val.([]interface{}); ok {
		result := make([]string, 0, len(list))
		for _, v := range list {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}

func init() {
	registry.Schema{
		Name:        "google_iam_workload_identity_pool_openid_config",
		ProductName: "iambeta",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceIAMBetaWorkloadIdentityPoolOpenIdConfig(),
	}.Register()
}
