load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "metadata",
    srcs = [
        "field.go",
        "metadata.go",
    ],
    deps = [
        "//third_party/magic_modules/v0/mmv1/api",
        "//third_party/magic_modules/v0/mmv1/google",
    ],
)

go_test(
    name = "metadata_test",
    srcs = [
        "field_test.go",
        "metadata_test.go",
    ],
    library = ":metadata",
    deps = [
        "//third_party/golang/cmp",
        "//third_party/magic_modules/v0/mmv1/api",
        "//third_party/magic_modules/v0/mmv1/api/product",
        "//third_party/magic_modules/v0/mmv1/google",
    ],
)
