package acctest

import (
	"github.com/hashicorp/terraform-provider-google/google/sweep"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// SharedConfigForRegion returns a common config setup needed for the sweeper
// functions for a given region
func SharedConfigForRegion(region string) (*transport_tpg.Config, error) {
	return sweep.SharedConfigForRegion(region)
}

func IsSweepableTestResource(resourceName string) bool {
	return sweep.IsSweepableTestResource(resourceName)
}
