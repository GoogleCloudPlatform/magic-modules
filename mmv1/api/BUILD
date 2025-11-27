load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "api",
    srcs = [
        "async.go",
        "compiler.go",
        "product.go",
        "resource.go",
        "timeouts.go",
        "type.go",
    ],
    deps = [
        "//base/go:log",
        "//third_party/golang/yamlv3:yaml",
        "//third_party/magic_modules/v0/mmv1/api/product",
        "//third_party/magic_modules/v0/mmv1/api/resource",
        "//third_party/magic_modules/v0/mmv1/api/utils",
        "//third_party/magic_modules/v0/mmv1/google",
    ],
)

go_test(
    name = "api_test",
    srcs = [
        "product_test.go",
        "resource_test.go",
        "type_test.go",
    ],
    library = ":api",
    deps = ["//third_party/magic_modules/v0/mmv1/api/product"],
)
