package google

import (
	"fmt"

	"google.golang.org/api/cloudresourcemanager/v1"
	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

type ResourceManagerClientOperationWaiter struct {
	Service *cloudresourcemanager.Service
	CommonOperationWaiter
}

func (w *ResourceManagerClientOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	return w.Service.Operations.Get(w.Op.Name).Do()
}

func resourceManagerClientOperationWaitTime(service *cloudresourcemanager.Service, op *cloudresourcemanager.Operation, activity string, timeoutMin int) error {
	w := &ResourceManagerClientOperationWaiter{
		Service: service,
	}
	if err := w.SetOp(op); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMin)
}

func resourceManagerClientOperationWait(service *cloudresourcemanager.Service, op *cloudresourcemanager.Operation, activity string) error {
	return resourceManagerClientOperationWaitTime(service, op, activity, 4)
}

func resourceManagerV2Beta1OperationWait(service *cloudresourcemanager.Service, op *resourceManagerV2Beta1.Operation, activity string) error {
	return resourceManagerV2Beta1OperationWaitTime(service, op, activity, 4)
}

func resourceManagerV2Beta1OperationWaitTime(service *cloudresourcemanager.Service, op *resourceManagerV2Beta1.Operation, activity string, timeoutMin int) error {
	opV1 := &cloudresourcemanager.Operation{}
	err := Convert(op, opV1)
	if err != nil {
		return err
	}

	return resourceManagerClientOperationWaitTime(service, opV1, activity, timeoutMin)
}
