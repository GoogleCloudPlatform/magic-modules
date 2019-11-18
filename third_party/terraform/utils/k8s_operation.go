package google

import (
	"fmt"
)

// K8sOperation is a struct that can contain a Kubernetes resource's Status block. It is not
// intended to be used for anything other than polling for the success of the given resource.
type K8sOperation struct {
	Metadata struct {
		Name      string
		Namespace string
		SelfLink  string
	}
	Status struct {
		Conditions []struct {
			Type    string
			Status  string
			Reason  string
			Message string
		}
	} `json:"status"`
}

// K8sOperationWaiter allows for polling against an arbitrary resource that implements the
// Kubernetes status schema.
type K8sOperationWaiter struct {
	Config   *Config
	Project  string
	Location string
	WaitURL  string
	Op       K8sOperation
}

// State will return a string representing the status of the Ready condition.
// No other conditions are currently returned as part of the state.
func (w *K8sOperationWaiter) State() string {
	for _, condition := range w.Op.Status.Conditions {
		if condition.Type == "Ready" {
			return fmt.Sprintf("%s:%s", condition.Type, condition.Status)
		}
	}
	return "NotImplemented"
}

func (w *K8sOperationWaiter) Error() error {
	if len(w.Op.Status.Conditions) == 0 {
		// When initially created the returned object doesn't have a
		// status block yet so continue polling until it does.
		return nil
	}

	for _, condition := range w.Op.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "False" {
			return fmt.Errorf("%s - %s", condition.Reason, condition.Message)
		}
	}
	return nil
}

func (w *K8sOperationWaiter) IsRetryable(error) bool {
	return false
}

func (w *K8sOperationWaiter) OpName() string {
	if w == nil {
		return "waiter:<nil>"
	}

	return w.Op.Metadata.SelfLink
}
func (w *K8sOperationWaiter) PendingStates() []string {
	return []string{"Ready:Unknown"}
}
func (w *K8sOperationWaiter) TargetStates() []string {
	return []string{"Ready:True"}
}

func (w *K8sOperationWaiter) SetOp(op interface{}) error {

	err := Convert(op, &w.Op)
	if err != nil {
		return err
	}

	return nil
}

func (w *K8sOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}

	url := fmt.Sprintf("%s", w.WaitURL)
	return sendRequest(w.Config, "GET", w.Project, url, nil)
}

func k8sOperationWaitTime(config *Config, res map[string]interface{}, project, url, activity string, timeoutMinutes int) error {
	op := K8sOperation{}
	err := Convert(res, &op)
	if err != nil {
		return err
	}

	w := &K8sOperationWaiter{
		Config:  config,
		WaitURL: url,
	}
	if err := w.SetOp(res); err != nil {
		return err
	}
	return OperationWait(w, activity, timeoutMinutes)
}
