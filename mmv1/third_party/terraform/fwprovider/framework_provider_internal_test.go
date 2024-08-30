package fwprovider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
)

func TestFrameworkProvider_impl(t *testing.T) {
	var _ provider.ProviderWithMetaSchema = New()
}
