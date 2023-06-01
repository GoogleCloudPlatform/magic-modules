package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/kms"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Deprecated: For backward compatibility ParseKmsCryptoKeyId is still working,
// but all new code should use ParseKmsCryptoKeyId in the kms package instead.
func ParseKmsCryptoKeyId(id string, config *transport_tpg.Config) (*KmsCryptoKeyId, error) {
	return kms.ParseKmsCryptoKeyId(id, config)
}
