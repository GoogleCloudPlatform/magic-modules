package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/servicenetworking/v1beta"
)

type ServiceNetworkingOperationWaiter struct {
	Service *servicenetworking.APIService
	Op      *servicenetworking.Operation
}

func (w *ServiceNetworkingOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.Service.Operations.Get(w.Op.Name).Do()

		if e, ok := err.(*googleapi.Error); ok && (e.Code == 429 || e.Code == 503) {
			return w.Op, fmt.Sprintf("%v", op.Done), nil
		} else if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got done: %v  when asking for operation %q", op.Done, w.Op.Name)

		return op, fmt.Sprintf("%v", op.Done), nil
	}
}

func (w *ServiceNetworkingOperationWaiter) Conf() *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending: []string{"false"},
		Target:  []string{"true"},
		Refresh: w.RefreshFunc(),
	}
}

// ServiceNetworkingOperationError wraps servicenetworking.Status and implements
// the error interface so it can be returned.
type ServiceNetworkingOperationError servicenetworking.Status

func (e ServiceNetworkingOperationError) Error() string {
	return e.Message
}

func serviceNetworkingOperationWait(config *Config, op *servicenetworking.Operation, activity string) error {
	return serviceNetworkingOperationWaitTime(config, op, activity, 10)
}

func serviceNetworkingOperationWaitTime(config *Config, op *servicenetworking.Operation, activity string, timeoutMinutes int) error {
	if op.Done {
		if op.Error != nil {
			return ServiceNetworkingOperationError(*op.Error)
		}
		return nil
	}

	w := &ServiceNetworkingOperationWaiter{
		Service: config.clientServiceNetworking,
		Op:      op,
	}

	state := w.Conf()
	state.Timeout = time.Duration(timeoutMinutes) * time.Minute
	state.MinTimeout = 2 * time.Second
	state.Delay = 5 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s (op %s): %s", activity, op.Name, err)
	}

	op = opRaw.(*servicenetworking.Operation)
	if op.Error != nil {
		return ServiceNetworkingOperationError(*op.Error)
	}

	return nil
}
