load("//tools/build_defs/go:go_binary.bzl", "go_binary")

go_binary(
    name = "main",
    srcs = ["main.go"],
    embedsrcs = [
        "//third_party/magic_modules/v0/mmv1/templates",  # keep
        "//third_party/magic_modules/v0/mmv1/third_party",  # keep
    ],
    visibility = ["//third_party/magic_modules/v0:users"],
    deps = [
        "//third_party/magic_modules/v0/mmv1/api",
        "//third_party/magic_modules/v0/mmv1/loader",
        "//third_party/magic_modules/v0/mmv1/openapi_generate",
        "//third_party/magic_modules/v0/mmv1/provider",
    ],
)
