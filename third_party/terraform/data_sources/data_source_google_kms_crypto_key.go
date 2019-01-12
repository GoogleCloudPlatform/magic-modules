package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/validation"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGoogleKmsCryptoKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsCryptoKeyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_ring": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: kmsCryptoKeyRingsEquivalent,
			},
			"rotation_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_template": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:     schema.TypeString,
							Required: true,
						},
						"protection_level": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "SOFTWARE",
							ValidateFunc: validation.StringInSlice([]string{"SOFTWARE", "HSM", ""}, false),
						},
					},
				},
			},
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func dataSourceGoogleKmsCryptoKeyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	keyRingId, err := parseKmsKeyRingId(d.Get("key_ring").(string), config)
	if err != nil {
		return err
	}

	cryptoKeyId := &kmsCryptoKeyId{
		KeyRingId: *keyRingId,
		Name:      d.Get("name").(string),
	}
	log.Printf("[DEBUG] Executing read for KMS CryptoKey %s", cryptoKeyId.cryptoKeyId())

	cryptoKey, err := config.clientKms.Projects.Locations.KeyRings.CryptoKeys.Get(cryptoKeyId.cryptoKeyId()).Do()
	if err != nil {
		return fmt.Errorf("Error reading CryptoKey: %s", err)
	}
	d.Set("key_ring", cryptoKeyId.KeyRingId.terraformId())
	d.Set("name", cryptoKeyId.Name)
	d.Set("rotation_period", cryptoKey.RotationPeriod)
	d.Set("self_link", cryptoKey.Name)

	if err = d.Set("version_template", flattenVersionTemplate(cryptoKey.VersionTemplate)); err != nil {
		return fmt.Errorf("Error setting version_template in state: %s", err.Error())
	}

	d.SetId(cryptoKeyId.cryptoKeyId())

	return nil
}
