package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"google.golang.org/api/dataproc/v1"
)

type DataprocClusterOperationWaiter struct {
	Service *dataproc.Service
	Op      *dataproc.Operation
}

func (w *DataprocClusterOperationWaiter) Conf(timeoutMinutes int) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{"false"},
		Target:     []string{"true"},
		Refresh:    w.RefreshFunc(),
		Timeout:    time.Duration(timeoutMinutes) * time.Minute,
		MinTimeout: 2 * time.Second,
	}
}

func (w *DataprocClusterOperationWaiter) RefreshFunc() resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		op, err := w.Service.Projects.Regions.Operations.Get(w.Op.Name).Do()

		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] Got %v while polling for operation %s's 'done' status", op.Done, w.Op.Name)

		return op, fmt.Sprint(op.Done), nil
	}
}

func dataprocClusterOperationWait(config *Config, op *dataproc.Operation, activity string, timeoutMinutes int) error {
	if op.Done {
		if op.Error != nil {
			return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
		}
		return nil
	}

	w := &DataprocClusterOperationWaiter{
		Service: config.clientDataproc,
		Op:      op,
	}

	opRaw, err := w.Conf(timeoutMinutes).WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for %s: %s", activity, err)
	}

	op = opRaw.(*dataproc.Operation)
	if op.Error != nil {
		return fmt.Errorf("Error code %v, message: %s", op.Error.Code, op.Error.Message)
	}

	return nil
}
