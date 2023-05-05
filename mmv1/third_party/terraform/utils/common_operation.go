package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func CommonRefreshFunc(w Waiter) resource.StateRefreshFunc {
	return tpgresource.CommonRefreshFunc(w)
}

func OperationWait(w Waiter, activity string, timeout time.Duration, pollInterval time.Duration) error {
	return tpgresource.OperationWait(w, activity, timeout, pollInterval)
}
