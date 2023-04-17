package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func redisInstanceGetRegionFromLocationID(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	region, err := getRegionFromSchema("region", "location_id", d, config)
	d.Set("region", region)

	return region, err
}
