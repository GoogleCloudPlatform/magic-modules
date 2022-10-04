package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceComputeAddresses(t *testing.T) {
	t.Parallel()

	addressName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	region := "europe-west8"
	region_bis := "asia-east1"
	dsName := "regional_addresses"
	dsFullName := fmt.Sprintf("data.google_compute_addresses.%s", dsName)
	dsAllName := "all_addresses"
	dsAllFullName := fmt.Sprintf("data.google_compute_addresses.%s", dsAllName)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeAddressesConfig(addressName, region, region_bis),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceComputeAddressesRegionSpecificCheck(t, addressName, dsFullName, region),
					testAccDataSourceComputeAddressesAllRegionsCheck(t, addressName, dsAllFullName, region, region_bis),
				),
			},
		},
	})
}

func testAccDataSourceComputeAddressesAllRegionsCheck(t *testing.T, address_name string, data_source_name string, expected_region string, expected_region_bis string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		ds_attr := ds.Primary.Attributes

		if ds_attr["addresses.#"] != "6" {
			return fmt.Errorf("addresses.# is not equal to 6")
		}

		var expected_addresses []expectedAddress
		for i := 0; i < 3; i++ {
			expected_addresses = append(expected_addresses, expectedAddress{
				name:   fmt.Sprintf("%s-%s-%d", address_name, expected_region, i),
				region: expected_region,
			})
		}
		for i := 0; i < 3; i++ {
			expected_addresses = append(expected_addresses, expectedAddress{
				name:   fmt.Sprintf("%s-%s-%d", address_name, expected_region_bis, i),
				region: expected_region_bis,
			})
		}

		for address_index := 0; address_index < 6; address_index++ {
			has_match := false
			for j := 0; j < len(expected_addresses); j++ {
				match, err := expected_addresses[j].checkAddressMatch(address_index, ds_attr)
				if err != nil {
					return err
				} else {
					if match {
						has_match = true
						expected_addresses = removeExpectedAddress(expected_addresses, j)
						break
					}
				}
			}
			if !has_match {
				return fmt.Errorf("unexpected address at index %d", address_index) // TODO improve
			}
		}

		if len(expected_addresses) != 0 {
			return fmt.Errorf("Addresses not found: %+v", expected_addresses)
		}

		return nil
	}
}

func removeExpectedAddress(s []expectedAddress, i int) []expectedAddress {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func testAccDataSourceComputeAddressesRegionSpecificCheck(t *testing.T, address_name string, data_source_name string, expected_region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		ds_attr := ds.Primary.Attributes

		if ds_attr["addresses.#"] != "3" {
			return fmt.Errorf("addresses.# is not equal to 3")
		}

		for n := 0; n < 3; n++ {
			map_key := fmt.Sprintf("addresses.%d.name", n)
			address_prefix := fmt.Sprintf("%s-%s", address_name, expected_region)

			if !strings.HasPrefix(ds_attr[map_key], address_prefix) {
				return fmt.Errorf("%s dont start with %s, got %s", map_key, address_name, ds_attr[map_key])
			}
		}

		for n := 0; n < 3; n++ {
			address_name := ds_attr[fmt.Sprintf("addresses.%d.name", n)]
			map_key := fmt.Sprintf("addresses.%d.labels.mykey", n)

			v, found := ds_attr[map_key]
			if !found {
				return fmt.Errorf("label with key 'mykey' not found for %s", address_name)
			}

			if v != "myvalue" {
				return fmt.Errorf("label value of 'mykey' not equal to 'myvalue' for %s, got %s", address_name, v)
			}
		}

		for n := 0; n < 3; n++ {
			map_key := fmt.Sprintf("addresses.%d.region", n)
			region, found := ds_attr[map_key]
			if !found {
				return fmt.Errorf("%s doesn't exists", map_key)
			}
			if region != expected_region {
				return fmt.Errorf("Unexpected region: got %s expected %s", region, expected_region)
			}
		}

		return nil
	}
}

func testAccDataSourceComputeAddressesConfig(addressName, region, region_bis string) string {
	return fmt.Sprintf(`
locals { 
	region = "%s"
	region_bis  = "%s"
	address_name = "%s"
}

resource "google_compute_address" "address" {
  count = 3

  region = local.region
  name = "${local.address_name}-${local.region}-${count.index}"
  labels = {
	mykey = "myvalue"
  }
}

resource "google_compute_address" "address_region_bis" {
  count = 3

  region = local.region_bis
  name = "${local.address_name}-${local.region_bis}-${count.index}"
  labels = {
	mykey = "myvalue"
  }
}

data "google_compute_addresses" "regional_addresses" {
	filter = "name:${local.address_name}-*"
	depends_on = [google_compute_address.address]
	region = local.region
}
data "google_compute_addresses" "all_addresses" {
	filter = "name:${local.address_name}-*"
	depends_on = [google_compute_address.address, google_compute_address.address_region_bis]
}
`, region, region_bis, addressName)
}

type expectedAddress struct {
	name   string
	region string
}

func (r expectedAddress) checkAddressMatch(index int, attrs map[string]string) (bool, error) {
	map_name := fmt.Sprintf("addresses.%d.name", index)

	if !strings.HasPrefix(attrs[map_name], r.name) {
		return false, nil
	}

	map_region := fmt.Sprintf("addresses.%d.region", index)
	region, found := attrs[map_region]
	if !found {
		return true, fmt.Errorf("%s doesn't exists", map_region)
	}
	if region != r.region {
		return true, fmt.Errorf("Unexpected region: got %s expected %s", region, r.region)
	}

	return true, nil
}
