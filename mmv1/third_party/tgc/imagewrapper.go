package compute

import (
	tpg_services_compute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// ResolveImageWrapper is a wrapper around compute.ResolveImage and recover
// from panic if client is il.
func ResolveImageWrapper(c *transport_tpg.Config, project, name, userAgent string) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			if c.Client == nil {
				ret = name
				err = nil
				return
			}
			panic(r)
		}
	}()

	return tpg_services_compute.ResolveImage(c, project, name, userAgent)
}
