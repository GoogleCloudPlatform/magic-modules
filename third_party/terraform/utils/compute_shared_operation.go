package google

import (
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

func computeSharedOperationWait(client *compute.Service, op interface{}, project, activity string) error {
	return computeSharedOperationWaitTime(client, op, project, activity, 4)
}

func computeSharedOperationWaitTime(client *compute.Service, op interface{}, project, activity string, minutes int) error {
	if op == nil {
		panic("Attempted to wait on an Operation that was nil.")
	}

	switch op.(type) {
	case *compute.Operation:
		return computeOperationWaitTime(client, op.(*compute.Operation), project, activity, minutes)
	case *computeBeta.Operation:
		return computeBetaOperationWaitTime(client, op.(*computeBeta.Operation), project, activity, minutes)
	default:
		panic("Attempted to wait on an Operation of unknown type.")
	}
}
