package google

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func redisInstanceGetRegionFromLocationID(d *schema.ResourceData, config *Config) (string, error) {
	region, err := getRegionFromSchema("region", "location_id", d, config)
	d.Set("region", region)

	return region, err
}
