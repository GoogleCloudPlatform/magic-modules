package memorystore_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	tpgmemorystore "github.com/hashicorp/terraform-provider-google/google/services/memorystore"
)

func TestUnitMemorystoreAuthToken_latestMemorystoreAuthToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input    interface{}
		wantName string
		wantErr  bool
	}{
		"nil slice": {
			input:   nil,
			wantErr: true,
		},
		"empty slice": {
			input:   []interface{}{},
			wantErr: true,
		},
		"missing name": {
			input: []interface{}{
				map[string]interface{}{"token": "abc"},
				map[string]interface{}{"name": ""},
			},
			wantErr: true,
		},
		"single version-1": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-1"},
			},
			wantName: "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-1",
		},
		"multiple versions ascending": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-1"},
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-2"},
			},
			wantName: "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-2",
		},
		"numeric ordering past version-9": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-10"},
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-9"},
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-2"},
			},
			wantName: "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-10",
		},
		"prefers versioned over non-versioned": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/custom-id"},
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-1"},
			},
			wantName: "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/version-1",
		},
		"falls back to non-versioned": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/custom-id"},
			},
			wantName: "projects/p/locations/l/instances/i/tokenAuthUsers/u/authTokens/custom-id",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got, err := tpgmemorystore.LatestMemorystoreAuthToken(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %#v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got["name"] != tc.wantName {
				t.Fatalf("got name %v, want %q", got["name"], tc.wantName)
			}
		})
	}
}

func TestAccMemorystoreAuthToken_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreAuthTokenDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemorystoreAuthToken_basic(context),
			},
			{
				ResourceName:      "google_memorystore_auth_token.token-default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMemorystoreAuthToken_rotation(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreAuthTokenDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemorystoreAuthToken_rotationStep1(context),
			},
			{
				ResourceName:      "google_memorystore_auth_token.token1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            "google_memorystore_auth_token.token2",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_policy"},
			},
			{
				Config: testAccMemorystoreAuthToken_rotationStep2(context),
			},
			{
				Config:      testAccMemorystoreAuthToken_acquireNotFound(context),
				ExpectError: regexp.MustCompile(`not found`),
			},
		},
	})
}

func testAccMemorystoreAuthToken_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance" {
  instance_id                 = "tf-test-ms-instance-%{random_suffix}"
  shard_count                 = 1
  node_type                   = "SHARED_CORE_NANO"
  location                    = "europe-west4"
  authorization_mode          = "TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_memorystore_auth_token" "token-default" {
  instance                    = google_memorystore_instance.instance.id
  user_id                     = "${google_memorystore_instance.instance.id}/tokenAuthUsers/default"
  location                    = "europe-west4"
}
`, context)
}

func testAccMemorystoreAuthToken_rotationStep1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance" {
  instance_id                 = "tf-test-ms-instance-%{random_suffix}"
  shard_count                 = 1
  node_type                   = "SHARED_CORE_NANO"
  location                    = "europe-west4"
  authorization_mode          = "TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_memorystore_auth_token" "token1" {
  instance                    = google_memorystore_instance.instance.instance_id
  user_id                     = "default"
  token_id                    = "version-1"
  location                    = "europe-west4"
}

resource "google_memorystore_auth_token" "token2" {
  instance                    = google_memorystore_instance.instance.instance_id
  user_id                     = "default"
  deletion_policy             = "ABANDON"
  location                    = "europe-west4"
  depends_on                  = [google_memorystore_auth_token.token1]
}
`, context)
}

func testAccMemorystoreAuthToken_rotationStep2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance" {
  instance_id                 = "tf-test-ms-instance-%{random_suffix}"
  shard_count                 = 1
  node_type                   = "SHARED_CORE_NANO"
  location                    = "europe-west4"
  authorization_mode          = "TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_memorystore_auth_token" "token2" {
  instance                    = google_memorystore_instance.instance.instance_id
  user_id                     = "default"
  deletion_policy             = "ABANDON"
  location                    = "europe-west4"
}
`, context)
}

func testAccMemorystoreAuthToken_acquireNotFound(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance" {
  instance_id                 = "tf-test-ms-instance-%{random_suffix}"
  shard_count                 = 1
  node_type                   = "SHARED_CORE_NANO"
  location                    = "europe-west4"
  authorization_mode          = "TOKEN_AUTH"
  deletion_protection_enabled = false
}

resource "google_memorystore_auth_token" "token2" {
  instance                    = google_memorystore_instance.instance.instance_id
  user_id                     = "default"
  deletion_policy             = "ABANDON"
  location                    = "europe-west4"
}

resource "google_memorystore_auth_token" "token_missing" {
  instance                    = google_memorystore_instance.instance.instance_id
  user_id                     = "default"
  token_id                    = "version-99"
  location                    = "europe-west4"
}
`, context)
}
