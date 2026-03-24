// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tpgiamresource

import "testing"

func TestFormatIAMResourceCanonicalID(t *testing.T) {
	t.Parallel()
	attrs := map[string]string{
		"project": "p1",
		"zone":    "us-central1-a",
		"name":    "d1",
	}
	keys := []string{"project", "zone", "name"}
	got, err := FormatIAMResourceCanonicalID("projects/%s/zones/%s/disks/%s", keys, attrs)
	if err != nil {
		t.Fatal(err)
	}
	want := "projects/p1/zones/us-central1-a/disks/d1"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestFormatIAMResourceCanonicalID_mismatchPlaceholders(t *testing.T) {
	t.Parallel()
	_, err := FormatIAMResourceCanonicalID("projects/%s/zones/%s/disks/%s", []string{"project", "zone"}, map[string]string{"project": "p", "zone": "z"})
	if err == nil {
		t.Fatal("expected error for placeholder count mismatch")
	}
}
