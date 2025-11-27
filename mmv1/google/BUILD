load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "google",
    srcs = [
        "slice_utils.go",
        "string_utils.go",
        "template_utils.go",
        "yaml_validator.go",
    ],
    deps = [
        "//base/go:log",
        "//third_party/golang/yamlv3:yaml",
    ],
)

go_test(
    name = "google_test",
    srcs = [
        "slice_utils_test.go",
        "string_utils_test.go",
    ],
    library = ":google",
)
