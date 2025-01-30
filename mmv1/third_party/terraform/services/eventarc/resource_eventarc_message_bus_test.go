package eventarc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccEventarcMessageBus(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic": testAccEventarcMessageBus_basic,
		"full":  testAccEventarcMessageBus_full,
	}

	for name, tc := range testCases {
		// shadow the tc variable into scope so that when
		// the loop continues, if t.Run hasn't executed tc(t)
		// yet, we don't have a race condition
		// see https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		tc := tc
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccEventarcMessageBus_basic(t *testing.T) {
	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMessageBusDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_basicCfg(context),
			},
			{
				ResourceName:      "google_eventarc_message_bus.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcMessageBus_basicCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_message_bus" "primary" {
  location       = "%{region}"
  message_bus_id = "tf-test-messagebus%{random_suffix}"
}
`, context)
}

func testAccEventarcMessageBus_full(t *testing.T) {
	region := envvar.GetTestRegionFromEnv()

	context := map[string]interface{}{
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"region":         region,
		"key":            acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-messagebus-key").CryptoKey.Name,
		"random_suffix":  acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMessageBusDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcMessageBus_fullCfg(context),
			},
			{
				ResourceName:      "google_eventarc_message_bus.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcMessageBus_fullCfg(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_kms_crypto_key_iam_member" "key_member" {
  crypto_key_id = "%{key}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-%{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com"
}

resource "google_eventarc_message_bus" "primary" {
  location        = "%{region}"
  message_bus_id  = "tf-test-messagebus%{random_suffix}"
  crypto_key_name = "%{key}"
  logging_config {
    log_severity = "ALERT"
  }
  depends_on = [google_kms_crypto_key_iam_member.key_member]
}
`, context)
}

func testAccCheckMessageBusDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_eventarc_message_bus" {
				continue
			}

			name := rs.Primary.Attributes["name"]

			url := fmt.Sprintf("https://eventarc.googleapis.com/v1/%s", name)
			_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})

			if err == nil {
				return fmt.Errorf("Error, message bus %s still exists", name)
			}
		}

		return nil
	}
}
