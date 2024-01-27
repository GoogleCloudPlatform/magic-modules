package compute

import (
	"testing"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestResolveImage_NotPanicWhenClientIsNil(t *testing.T) {
	got, err := ResolveImage(&transport_tpg.Config{}, "project", "name", "useragent")
	if err != nil {
		t.Fatal(err)
	}
	if got != "name" {
		t.Errorf("ResolveImage() = %s, want = %s", got, "name")
	}
}
