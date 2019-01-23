package google

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func dataSourceGoogleKmsSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleKmsSecretRead,
		Schema: map[string]*schema.Schema{
			"crypto_key": {
				Type: schema.TypeString,
				Description: "Full ID of the crypto key to use for decryption in the format" +
					"(`projects/{project}/locations/{location}/keyRings/{keyRing}/cryptoKeys/{cryptoKey}`",
				Required: true,
			},
			"ciphertext": {
				Type:        schema.TypeString,
				Description: "Base64-encoded ciphertext",
				Required:    true,
			},
			"additional_authenticated_data": {
				Type:        schema.TypeString,
				Description: "Base64-encoded optional data originally supplied during encryption",
				Optional:    true,
			},
			"plaintext": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceGoogleKmsSecretRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	tfCryptoKeyId, err := parseKmsCryptoKeyId(d.Get("crypto_key").(string), config)
	if err != nil {
		return err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(d.Get("ciphertext").(string))
	if err != nil {
		return fmt.Errorf("Failed to base64 decode ciphertext: %s", err)
	}

	additionalAuthenticatedData, err := base64.StdEncoding.DecodeString(d.Get("additional_authenticated_data").(string))
	if err != nil {
		return fmt.Errorf("failed to base64 decode additional_authenticated_data: %s", err)
	}

	ctx := context.Background()
	decryptResp, err := config.clientKms.Decrypt(ctx, &kmspb.DecryptRequest{
		Name:                        tfCryptoKeyId.cryptoKeyId(),
		Ciphertext:                  ciphertext,
		AdditionalAuthenticatedData: additionalAuthenticatedData,
	})
	if err != nil {
		return fmt.Errorf("Error decrypting ciphertext: %s", err)
	}

	log.Printf("[INFO] Successfully decrypted ciphertext: %s", ciphertext)

	d.Set("plaintext", string(decryptResp.Plaintext))
	d.SetId(time.Now().UTC().String())

	return nil
}
