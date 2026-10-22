package dataflow

import (
	"testing"

	dataflow "google.golang.org/api/dataflow/v1b3"
)

func TestDataflowFlexJobWorkerMachineType(t *testing.T) {
	t.Parallel()

	workerPool := &dataflow.WorkerPool{MachineType: "n2d-standard-2"}

	cases := map[string]struct {
		options map[string]interface{}
		want    string
	}{
		"java pipeline options": {
			options: map[string]interface{}{"workerMachineType": "n2d-standard-2"},
			want:    "n2d-standard-2",
		},
		"python pipeline options": {
			options: map[string]interface{}{"machine_type": "n2d-standard-2"},
			want:    "n2d-standard-2",
		},
		"empty option": {
			options: map[string]interface{}{"workerMachineType": ""},
			want:    "n2d-standard-2",
		},
		"no options": {
			options: map[string]interface{}{},
			want:    "n2d-standard-2",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := dataflowFlexJobWorkerMachineType(tc.options, workerPool); got != tc.want {
				t.Errorf("dataflowFlexJobWorkerMachineType(%v) = %q, want %q", tc.options, got, tc.want)
			}
		})
	}
}
