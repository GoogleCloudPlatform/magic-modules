package resourcemanager

import (
	"context"
	"reflect"
	"testing"
)

func testGoogleProjectStateDataV1() map[string]any {
	return map[string]any{}
}

func testGoogleProjectStateDataV2() map[string]any {
	return map[string]any{
		"deletion_policy": "PREVENT",
	}
}

func TestGoogleProjectStateUpgradeV1(t *testing.T) {
	expected := testGoogleProjectStateDataV2()
	actual, err := resourceGoogleProjectStateUpgradeV1(context.Background(), testGoogleProjectStateDataV1(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
