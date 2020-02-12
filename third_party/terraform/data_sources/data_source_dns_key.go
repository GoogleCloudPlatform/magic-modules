package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/dns/v1"
)

func dataSourceDNSKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSKeyRead,

		Schema: map[string]*schema.Schema{
			"managed_zone": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"key_signing_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dnsKeySchema(),
			},
			"zone_signing_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dnsKeySchema(),
			},
		},
	}
}

func dnsKeySchema() *schema.Resource {
  return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"digests": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"digest": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"key_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_tag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func flattenSigningKeys(signingKeys []*dns.DnsKey, keyType string) []map[string]interface{} {
	var keys []map[string]interface{}

	for _, signingKey := range signingKeys {
		if signingKey != nil && signingKey.Type == keyType {
			data := map[string]interface{}{
				"algorithm":     signingKey.Algorithm,
				"creation_time": signingKey.CreationTime,
				"description":   signingKey.Description,
				"digests":       flattenDigests(signingKey.Digests),
				"id":            signingKey.Id,
				"is_active":     signingKey.IsActive,
				"key_length":    signingKey.KeyLength,
				"key_tag":       signingKey.KeyTag,
				"public_key":    signingKey.PublicKey,
			}

			keys = append(keys, data)
		}
	}

	return keys
}

func flattenDigests(dnsKeyDigests []*dns.DnsKeyDigest) []map[string]interface{} {
	var digests []map[string]interface{}

	for _, dnsKeyDigest := range dnsKeyDigests {
		if dnsKeyDigest != nil {
			data := map[string]interface{}{
				"digest": dnsKeyDigest.Digest,
				"type":   dnsKeyDigest.Type,
			}

			digests = append(digests, data)
		}
	}

	return digests
}

func dataSourceDNSKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	fv, err := parseProjectFieldValue("managedZones", d.Get("managed_zone").(string), "project", d, config, false)
	if err != nil {
		return err
	}
	project := fv.Project
	managedZone := fv.Name

	log.Printf("[DEBUG] Fetching DNS keys from managed zone %s", managedZone)

	response, err := config.clientDns.DnsKeys.List(project, managedZone).Do()
	if err != nil && !isGoogleApiErrorWithCode(err, 404) {
		return fmt.Errorf("error retrieving DNS keys: %s", err)
	}

	log.Printf("[DEBUG] Fetched DNS keys from managed zone %s", managedZone)

	d.Set("project", project)
	d.Set("key_signing_keys", flattenSigningKeys(response.DnsKeys, "keySigning"))
	d.Set("zone_signing_keys", flattenSigningKeys(response.DnsKeys, "zoneSigning"))
	d.SetId(fmt.Sprintf("projects/%s/managedZones/%s", project, managedZone))

	return nil
}
