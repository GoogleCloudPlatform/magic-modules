package cai2hcl

import (
	transport_tpg "github.com/hashicorp/terraform-provider-google-beta/google-beta/transport"
)

func NewConfig() *transport_tpg.Config {
	return &transport_tpg.Config{
		Project:   "",
		Zone:      "",
		Region:    "",
		UserAgent: "",
	}
}
