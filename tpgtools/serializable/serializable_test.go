package serializable

import (
	"testing"
)

func TestListOfResources(t *testing.T) {	
	services, err := ListOfResources("test_specs")
	if err != nil {
		t.Errorf("received error: %v", err)
	}
	if len(services) != 1 {
		t.Errorf("expected 1 service, got: %v", len(services))
	}
	if len(services[0].Resources) != 1 {
		t.Errorf("expected 1 resource, got: %v", len(services[0].Resources))
	}
}
