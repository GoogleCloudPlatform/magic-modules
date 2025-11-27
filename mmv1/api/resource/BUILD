load("//tools/build_defs/go:go_library.bzl", "go_library")

go_library(
    name = "resource",
    srcs = [
        "custom_code.go",
        "datasource.go",
        "docs.go",
        "examples.go",
        "iam_policy.go",
        "nested_query.go",
        "reference_links.go",
        "sample.go",
        "step.go",
        "sweeper.go",
        "tgc.go",
        "validation.go",
    ],
    deps = [
        "//base/go:log",
        "//third_party/golang/yamlv3:yaml",
        "//third_party/magic_modules/v0/mmv1/api/product",
        "//third_party/magic_modules/v0/mmv1/api/utils",
        "//third_party/magic_modules/v0/mmv1/google",
    ],
)
