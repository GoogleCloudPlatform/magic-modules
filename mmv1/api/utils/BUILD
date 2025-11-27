load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "utils",
    srcs = [
        "omit_defaults.go",
        "utils.go",
    ],
)

go_test(
    name = "utils_test",
    srcs = [
        "omit_defaults_test.go",
        "utils_test.go",
    ],
    library = ":utils",
)
