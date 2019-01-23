package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

var SharedKeyRing = "tftest-shared-keyring-1"
var SharedCryptoKey = "tftest-shared-key-1"

type bootstrappedKMS struct {
	*kmspb.KeyRing
	*kmspb.CryptoKey
}

/**
* BootstrapKMSkey will return a KMS key that can be used in tests that are
* testing KMS integration with other resources.
*
* This will either return an existing key or create one if it hasn't been created
* in the project yet. The motivation is because keyrings don't get deleted and we
* don't want a linear growth of disabled keyrings in a project. We also don't want
* to incur the overhead of creating a new project for each test that needs to use
* a KMS key.
**/
func BootstrapKMSKey(t *testing.T) bootstrappedKMS {
	if v := os.Getenv("TF_ACC"); v == "" {
		log.Println("Acceptance tests and bootstrapping skipped unless env 'TF_ACC' set")

		// If not running acceptance tests, return an empty object
		return bootstrappedKMS{
			&kmspb.KeyRing{},
			&kmspb.CryptoKey{},
		}
	}

	projectID := getTestProjectFromEnv()
	locationID := "global"
	keyRingParent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)
	keyRingName := fmt.Sprintf("%s/keyRings/%s", keyRingParent, SharedKeyRing)
	keyParent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", projectID, locationID, SharedKeyRing)
	keyName := fmt.Sprintf("%s/cryptoKeys/%s", keyParent, SharedCryptoKey)

	config := Config{
		Credentials: getTestCredsFromEnv(),
		Project:     getTestProjectFromEnv(),
		Region:      getTestRegionFromEnv(),
		Zone:        getTestZoneFromEnv(),
	}

	if err := config.loadAndValidate(); err != nil {
		t.Errorf("Unable to bootstrap KMS key: %s", err)
	}

	// Get or Create the hard coded shared keyring for testing
	ctx := context.Background()
	kmsClient := config.clientKms
	keyRing, err := kmsClient.GetKeyRing(ctx, &kmspb.GetKeyRingRequest{
		Name: keyRingName,
	})
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			keyRing, err = kmsClient.CreateKeyRing(ctx, &kmspb.CreateKeyRingRequest{
				Parent:    keyRingParent,
				KeyRingId: SharedKeyRing,
			})
			if err != nil {
				t.Errorf("Unable to bootstrap KMS key. Cannot create keyRing: %s", err)
			}
		} else {
			t.Errorf("Unable to bootstrap KMS key. Cannot retrieve keyRing: %s", err)
		}
	}

	if keyRing == nil {
		t.Fatalf("Unable to bootstrap KMS key. keyRing is nil!")
	}

	// Get or Create the hard coded, shared crypto key for testing
	cryptoKey, err := kmsClient.GetCryptoKey(ctx, &kmspb.GetCryptoKeyRequest{
		Name: keyName,
	})
	if err != nil {
		if isGoogleApiErrorWithCode(err, 404) {
			cryptoKey, err = kmsClient.CreateCryptoKey(ctx, &kmspb.CreateCryptoKeyRequest{
				Parent:      keyParent,
				CryptoKeyId: SharedCryptoKey,
				CryptoKey: &kmspb.CryptoKey{
					Purpose: kmspb.CryptoKey_ENCRYPT_DECRYPT,
				},
			})
			if err != nil {
				t.Errorf("Unable to bootstrap KMS key. Cannot create new CryptoKey: %s", err)
			}

		} else {
			t.Errorf("Unable to bootstrap KMS key. Cannot call CryptoKey service: %s", err)
		}
	}

	if cryptoKey == nil {
		t.Fatalf("Unable to bootstrap KMS key. CryptoKey is nil!")
	}

	return bootstrappedKMS{
		keyRing,
		cryptoKey,
	}
}
