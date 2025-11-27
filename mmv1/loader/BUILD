load("//tools/build_defs/go:go_library.bzl", "go_library")

go_library(
    name = "loader",
    srcs = [
        "custom_errors.go",
        "file_ops.go",
        "loader.go",
    ],
    deps = [
        "//base/go:log",
        "//third_party/magic_modules/v0/mmv1/api",
    ],
)
