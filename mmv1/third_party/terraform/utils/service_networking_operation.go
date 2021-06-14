package google

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/api/servicenetworking/v1"
)

type ServiceNetworkingOperationWaiter struct {
	Service *servicenetworking.APIService
	CommonOperationWaiter
	d *schema.ResourceData
}

func (w *ServiceNetworkingOperationWaiter) QueryOp() (interface{}, error) {
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func serviceNetworkingOperationWaitTime(d *schema.ResourceData, config *Config, op *servicenetworking.Operation, activity, userAgent string, timeout time.Duration) error {
	w := &ServiceNetworkingOperationWaiter{
		Service: config.NewServiceNetworkingClient(userAgent),
		d:       d,
	}

	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeout, config.PollInterval)
}
