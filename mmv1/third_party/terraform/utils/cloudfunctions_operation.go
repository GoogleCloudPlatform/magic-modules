package google

import (
	"time"

	cloudfunctions_tpg "github.com/hashicorp/terraform-provider-google/google/services/cloudfunctions"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudfunctions/v1"
)

// Deprecated: For backward compatibility cloudFunctionsOperationWait is still working,
// but all new code should use CloudFunctionsOperationWait in the cloudfunctions package instead.
func cloudFunctionsOperationWait(config *transport_tpg.Config, op *cloudfunctions.Operation, activity, userAgent string, timeout time.Duration) error {
	return cloudfunctions_tpg.CloudFunctionsOperationWait(config, op, activity, userAgent, timeout)
}

// Deprecated: For backward compatibility IsCloudFunctionsSourceCodeError is still working,
// but all new code should use IsCloudFunctionsSourceCodeError in the cloudfunctions package instead.
func IsCloudFunctionsSourceCodeError(err error) (bool, string) {
	return cloudfunctions_tpg.IsCloudFunctionsSourceCodeError(err)
}
