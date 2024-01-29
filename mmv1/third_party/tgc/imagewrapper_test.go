package compute

import (
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestResolveImageWrapper_NotPanicWhenClientIsNil(t *testing.T) {
	got, err := ResolveImageWrapper(&transport_tpg.Config{}, "project", "name", "useragent")
	if err != nil {
		t.Fatal(err)
	}
	if got != "name" {
		t.Errorf("ResolveImageWrapper() = %s, want = %s", got, "name")
	}
}
