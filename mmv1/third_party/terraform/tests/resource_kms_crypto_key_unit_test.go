package google

import (
	"context"
	"testing"
	"time"
)

func TestCryptoKeyIdParsing(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		ImportId            string
		ExpectedError       bool
		ExpectedTerraformId string
		ExpectedCryptoKeyId string
		Config              *Config
	}{
		"id is in project/location/keyRingName/cryptoKeyName format": {
			ImportId:            "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
		},
		"id is in domain:project/location/keyRingName/cryptoKeyName format": {
			ImportId:            "example.com:test-project/us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "example.com:test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/example.com:test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
		},
		"id contains name that is longer than 63 characters": {
			ImportId:      "test-project/us-central1/test-key-ring/can-you-believe-that-this-cryptokey-name-is-this-extravagantly-long",
			ExpectedError: true,
		},
		"id is in location/keyRingName/cryptoKeyName format": {
			ImportId:            "us-central1/test-key-ring/test-key-name",
			ExpectedError:       false,
			ExpectedTerraformId: "test-project/us-central1/test-key-ring/test-key-name",
			ExpectedCryptoKeyId: "projects/test-project/locations/us-central1/keyRings/test-key-ring/cryptoKeys/test-key-name",
			Config:              &Config{Project: "test-project"},
		},
		"id is in location/keyRingName/cryptoKeyName format without project in config": {
			ImportId:      "us-central1/test-key-ring/test-key-name",
			ExpectedError: true,
			Config:        &Config{Project: ""},
		},
	}

	for tn, tc := range cases {
		cryptoKeyId, err := ParseKmsCryptoKeyId(tc.ImportId, tc.Config)

		if tc.ExpectedError && err == nil {
			t.Fatalf("bad: %s, expected an error", tn)
		}

		if err != nil {
			if tc.ExpectedError {
				continue
			}
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		if cryptoKeyId.TerraformId() != tc.ExpectedTerraformId {
			t.Fatalf("bad: %s, expected Terraform ID to be `%s` but is `%s`", tn, tc.ExpectedTerraformId, cryptoKeyId.TerraformId())
		}

		if cryptoKeyId.CryptoKeyId() != tc.ExpectedCryptoKeyId {
			t.Fatalf("bad: %s, expected CryptoKey ID to be `%s` but is `%s`", tn, tc.ExpectedCryptoKeyId, cryptoKeyId.CryptoKeyId())
		}
	}
}

func TestCryptoKeyNextRotationCalculation(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	period, _ := time.ParseDuration("1000000s")

	expected := now.Add(period).Format(time.RFC3339Nano)

	timestamp, err := kmsCryptoKeyNextRotation(now, "1000000s")

	if err != nil {
		t.Fatalf("unexpected failure parsing time %s and duration 1000s: %s", now, err.Error())
	}

	if expected != timestamp {
		t.Fatalf("expected %s to equal %s", timestamp, expected)
	}
}

func TestCryptoKeyNextRotationCalculation_validation(t *testing.T) {
	t.Parallel()

	_, errs := validateKmsCryptoKeyRotationPeriod("86399s", "rotation_period")

	if len(errs) == 0 {
		t.Fatalf("Periods of less than a day should be invalid")
	}

	_, errs = validateKmsCryptoKeyRotationPeriod("100000.0000000001s", "rotation_period")

	if len(errs) == 0 {
		t.Fatalf("Numbers with more than 9 fractional digits are invalid")
	}
}

func TestCryptoKeyStateUpgradeV0(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Attributes map[string]interface{}
		Expected   map[string]string
		Meta       interface{}
	}{
		"change key_ring from terraform id fmt to link fmt": {
			Attributes: map[string]interface{}{
				"key_ring": "my-project/my-location/my-key-ring",
			},
			Expected: map[string]string{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Meta: &Config{},
		},
		"key_ring link fmt stays as link fmt": {
			Attributes: map[string]interface{}{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Expected: map[string]string{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Meta: &Config{},
		},
		"key_ring without project to link fmt": {
			Attributes: map[string]interface{}{
				"key_ring": "my-location/my-key-ring",
			},
			Expected: map[string]string{
				"key_ring": "projects/my-project/locations/my-location/keyRings/my-key-ring",
			},
			Meta: &Config{
				Project: "my-project",
			},
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			actual, err := ResourceKMSCryptoKeyUpgradeV0(context.Background(), tc.Attributes, tc.Meta)

			if err != nil {
				t.Error(err)
			}

			for k, v := range tc.Expected {
				if actual[k] != v {
					t.Errorf("expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
						k, v, k, actual[k], actual)
				}
			}
		})
	}
}
