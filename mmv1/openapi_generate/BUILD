load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "openapi_generate",
    srcs = ["parser.go"],
    embedsrcs = ["header.txt"],
    deps = [
        "//third_party/golang/kin_openapi/current/openapi3",
        "//third_party/golang/yamlv3:yaml",
        "//third_party/magic_modules/v0/mmv1/api",
        "//third_party/magic_modules/v0/mmv1/api/product",
        "//third_party/magic_modules/v0/mmv1/api/resource",
        "//third_party/magic_modules/v0/mmv1/google",
    ],
)

go_test(
    name = "openapi_generate_test",
    srcs = ["parser_test.go"],
    embedsrcs = ["test_data/test_api.yaml"],
    library = ":openapi_generate",
    deps = ["//third_party/golang/kin_openapi/current/openapi3"],
)
