load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_test(
    name = "auditmanager_test",
    srcs = ["data_source_auditmanager_framework_audit_scope_report_test.go"],
)

go_library(
    name = "auditmanager",
    srcs = ["data_source_auditmanager_framework_audit_scope_report.go"],
    deps = [
        "//third_party/golang/github.com/hashicorp/terraform-plugin-sdk/v2:diag",
        "//third_party/golang/github.com/hashicorp/terraform-plugin-sdk/v2/helper:schema",
        "//third_party/golang/terraform_providers/google/google:tpgresource",
        "//third_party/golang/terraform_providers/google/google:transport_tpg",
    ],
)
