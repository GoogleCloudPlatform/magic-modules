package transport

import (
	"reflect"
	"testing"
)

func TestListItemAsMap(t *testing.T) {
	t.Parallel()

	got, err := listItemAsMap(map[string]interface{}{"name": "ts1", "port": 443})
	if err != nil {
		t.Fatalf("map item: %v", err)
	}
	if got["name"] != "ts1" {
		t.Fatalf("map item name = %v", got["name"])
	}

	got, err = listItemAsMap("kvm1")
	if err != nil {
		t.Fatalf("string item: %v", err)
	}
	want := map[string]interface{}{"name": "kvm1"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("string item = %#v, want %#v", got, want)
	}

	if _, err := listItemAsMap(1); err == nil {
		t.Fatal("expected error for int item")
	}
}
