package cmd

import (
	"strings"
	"testing"
)

func TestCheckLegacyApiaryDiff_ForbiddenImports(t *testing.T) {
	testCases := []struct {
		name string
		diff string
	}{
		{
			name: "single line import",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -10,0 +11,1 @@
+import "google.golang.org/api/compute/v1"
`,
		},
		{
			name: "single line import with alias",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -10,0 +11,1 @@
+import compute "google.golang.org/api/compute/v1"
`,
		},
		{
			name: "block import line",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -10,0 +11,1 @@
+	"google.golang.org/api/compute/v1"
`,
		},
		{
			name: "block import line with alias",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -10,0 +11,1 @@
+	compute "google.golang.org/api/compute/v0.alpha"
`,
		},
		{
			name: "blank import",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -10,0 +11,1 @@
+	_ "google.golang.org/api/compute/v1"
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matches, err := CheckLegacyApiaryDiff(strings.NewReader(tc.diff))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(matches) != 1 {
				t.Fatalf("expected 1 match, got %d", len(matches))
			}
			if matches[0].File != "mmv1/third_party/terraform/services/compute/resource_instance.go" {
				t.Errorf("expected file mmv1/third_party/terraform/services/compute/resource_instance.go, got %s", matches[0].File)
			}
			if matches[0].Pattern != "legacy compute import" {
				t.Errorf("expected pattern 'legacy compute import', got %s", matches[0].Pattern)
			}
		})
	}
}

func TestCheckLegacyApiaryDiff_ForbiddenCalls(t *testing.T) {
	testCases := []struct {
		name string
		diff string
	}{
		{
			name: "tpgcompute.NewClient",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -50,0 +51,1 @@
+	c, err := tpgcompute.NewClient(config)
`,
		},
		{
			name: "compute_tpg.NewClient",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -50,0 +51,1 @@
+	c, err := compute_tpg.NewClient(config)
`,
		},
		{
			name: "DEPRECATED_LegacyApiaryClient property",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -50,0 +51,1 @@
+	client := config.DEPRECATED_LegacyApiaryClient
`,
		},
		{
			name: "tpgcompute.DEPRECATED_LegacyApiaryClient func",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -50,0 +51,1 @@
+	client := tpgcompute.DEPRECATED_LegacyApiaryClient(config)
`,
		},
		{
			name: "compute_tpg.DEPRECATED_LegacyApiaryClient func",
			diff: `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -50,0 +51,1 @@
+	client := compute_tpg.DEPRECATED_LegacyApiaryClient(config)
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			matches, err := CheckLegacyApiaryDiff(strings.NewReader(tc.diff))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(matches) != 1 {
				t.Fatalf("expected 1 match, got %d", len(matches))
			}
			if matches[0].Pattern != "legacy compute invocation" {
				t.Errorf("expected pattern 'legacy compute invocation', got %s", matches[0].Pattern)
			}
		})
	}
}

func TestCheckLegacyApiaryDiff_AllowedChanges(t *testing.T) {
	diff := `diff --git a/mmv1/third_party/terraform/services/compute/resource_instance.go b/mmv1/third_party/terraform/services/compute/resource_instance.go
--- a/mmv1/third_party/terraform/services/compute/resource_instance.go
+++ b/mmv1/third_party/terraform/services/compute/resource_instance.go
@@ -10,3 +10,3 @@
-	"google.golang.org/api/compute/v1"
-	c, err := tpgcompute.NewClient(config)
 	existingContextLine := config.DEPRECATED_LegacyApiaryClient
+	// Modern replacement
+	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
+		Config: config,
+		Method: "GET",
+		RawURL: url,
+	})
+	"google.golang.org/api/container/v1"
diff --git a/mmv1/third_party/terraform/services/compute/DEPRECATED_LegacyApiaryClient.go b/mmv1/third_party/terraform/services/compute/DEPRECATED_LegacyApiaryClient.go
--- a/mmv1/third_party/terraform/services/compute/DEPRECATED_LegacyApiaryClient.go
+++ b/mmv1/third_party/terraform/services/compute/DEPRECATED_LegacyApiaryClient.go
@@ -1,1 +1,1 @@
-old
+new
`

	matches, err := CheckLegacyApiaryDiff(strings.NewReader(diff))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches, got %d: %+v", len(matches), matches)
	}
}
