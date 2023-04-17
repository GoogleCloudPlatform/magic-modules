package transport

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ProviderBatching struct {
	SendAfter      types.String `tfsdk:"send_after"`
	EnableBatching types.Bool   `tfsdk:"enable_batching"`
}

var ProviderBatchingAttributes = map[string]attr.Type{
	"send_after":      types.StringType,
	"enable_batching": types.BoolType,
}

// ProviderMetaModel describes the provider meta model
type ProviderMetaModel struct {
	ModuleName types.String `tfsdk:"module_name"`
}
