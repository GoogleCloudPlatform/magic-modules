package apikeys

import (
	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// The Client is the base struct of all operations.  This will receive the
// Get, Delete, List, and Apply operations on all resources.
type Client struct {
	Config *dcl.Config
}

// NewClient creates a client that retries all operations a few times each.
func NewClient(c *dcl.Config) *Client {
	return &Client{
		Config: c,
	}
}
