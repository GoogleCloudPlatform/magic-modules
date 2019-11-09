package google

import (
	"google.golang.org/api/compute/v1"
)

func computeSharedOperationWait(client *compute.Service, op interface{}, project, activity string) error {
	return computeSharedOperationWaitTime(client, op, project, activity, 4)
}

// This is a shell around computeOperationWaitTime. It was originall meant to type switch between Beta and GA wait
// operations but it now serves to differentiate handwritten resource calls to computeWait from generated. This method
// should be eventually removed when the distinction is no longer needed.
func computeSharedOperationWaitTime(client *compute.Service, op interface{}, project, activity string, minutes int) error {
	if op == nil {
		panic("Attempted to wait on an Operation that was nil.")
	}

	return computeOperationWaitTime(client, op, project, activity, minutes)
}
