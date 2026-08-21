package networkconnectivity

import (
	"slices"
	"testing"
)

// TestUnitNetworkConnectivityTransport_networkHubExactlyOneOf guards against
// regressions of https://github.com/hashicorp/terraform-provider-google/issues/28425
// where `network` was incorrectly Required, blocking NCC hub-only Transports.
func TestUnitNetworkConnectivityTransport_networkHubExactlyOneOf(t *testing.T) {
	t.Parallel()

	r := ResourceNetworkConnectivityTransport()

	network, ok := r.Schema["network"]
	if !ok {
		t.Fatal("expected network in schema")
	}
	if network.Required {
		t.Error("network must not be Required; exactly one of network or hub must be set")
	}
	if !network.Optional {
		t.Error("network should be Optional")
	}
	if !slices.Contains(network.ExactlyOneOf, "network") {
		t.Errorf("network ExactlyOneOf = %v, want to include network", network.ExactlyOneOf)
	}

	hub, ok := r.Schema["hub"]
	if !ok {
		// hub is beta-only; on GA it is omitted and ExactlyOneOf is filtered to network alone.
		if len(network.ExactlyOneOf) != 1 || network.ExactlyOneOf[0] != "network" {
			t.Errorf("GA network ExactlyOneOf = %v, want [network]", network.ExactlyOneOf)
		}
		return
	}

	if hub.Required {
		t.Error("hub must not be Required")
	}
	if !hub.Optional {
		t.Error("hub should be Optional")
	}
	if !slices.Contains(network.ExactlyOneOf, "hub") {
		t.Errorf("network ExactlyOneOf = %v, want to include hub", network.ExactlyOneOf)
	}
	if !slices.Contains(hub.ExactlyOneOf, "hub") || !slices.Contains(hub.ExactlyOneOf, "network") {
		t.Errorf("hub ExactlyOneOf = %v, want hub and network", hub.ExactlyOneOf)
	}
}
