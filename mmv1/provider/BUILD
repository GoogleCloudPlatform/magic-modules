load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "provider",
    srcs = [
        "provider.go",
        "template_data.go",
        "terraform.go",
        "terraform_oics.go",
        "terraform_tgc.go",
        "terraform_tgc_cai2hcl.go",
        "terraform_tgc_next.go",
    ],
    deps = [
        "//base/go:log",
        "//third_party/golang/dircopy:copy",
        "//third_party/golang/yamlv3:yaml",
        "//third_party/magic_modules/v0/mmv1/api",
        "//third_party/magic_modules/v0/mmv1/api/metadata",
        "//third_party/magic_modules/v0/mmv1/api/product",
        "//third_party/magic_modules/v0/mmv1/api/resource",
        "//third_party/magic_modules/v0/mmv1/google",
    ],
)

go_test(
    name = "provider_test",
    srcs = ["terraform_tgc_next_test.go"],
    library = ":provider",
)
