package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestFilestoreInstanceMigrateState(t *testing.T) {
	cases := map[string]struct {
		StateVersion int
		Attributes   map[string]string
		Expected     map[string]string
		Meta         interface{}
	}{
		"migrate zone to location": {
			StateVersion: 0,
			Attributes: map[string]string{
				"zone": "us-central1-a",
			},
			Expected: map[string]string{
				"zone":     "us-central1-a",
				"location": "us-central1-a",
			},
		},
	}
	for tn, tc := range cases {
		is := &terraform.InstanceState{
			ID:         "i-abc123",
			Attributes: tc.Attributes,
		}
		is, err := resourceFilestoreInstanceMigrateState(
			tc.StateVersion, is, tc.Meta)

		if err != nil {
			t.Fatalf("bad: %s, err: %#v", tn, err)
		}

		for k, v := range tc.Expected {
			if is.Attributes[k] != v {
				t.Fatalf(
					"bad: %s\n\n expected: %#v -> %#v\n got: %#v -> %#v\n in: %#v",
					tn, k, v, k, is.Attributes[k], is.Attributes)
			}
		}
	}
}

func TestAccFilestoreInstance_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_update(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
			{
				Config: testAccFilestoreInstance_update2(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  file_shares {
    capacity_gb = 2660
    name        = "share"
  }
  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
  labels = {
    baz = "qux"
  }
  tier        = "PREMIUM"
  description = "An instance created during testing."
}
`, name)
}

func testAccFilestoreInstance_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  file_shares {
    capacity_gb = 2760
    name        = "share"
  }
  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
  tier        = "PREMIUM"
  description = "A modified instance created during testing."
}
`, name)
}

func TestAccFilestoreInstance_reservedIpRange_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstance_reservedIpRange_update(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
			{
				Config: testAccFilestoreInstance_reservedIpRange_update2(name),
			},
			{
				ResourceName:            "google_filestore_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "location"},
			},
		},
	})
}

func testAccFilestoreInstance_reservedIpRange_update(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier    = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.30.0/29"
  }
}
`, name)
}

func testAccFilestoreInstance_reservedIpRange_update2(name string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "tf-instance-%s"
  zone = "us-central1-b"
  tier    = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network           = "default"
    modes             = ["MODE_IPV4"]
    reserved_ip_range = "172.19.31.0/29"
  }
}
`, name)
}
