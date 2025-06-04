load("//tools/build_defs/go:go_library.bzl", "go_library")
load("//tools/build_defs/go:go_test.bzl", "go_test")

go_library(
    name = "pubsub",
    srcs = [
        "data_source_pubsub_subscription.go",
        "data_source_pubsub_topic.go",
        "iam_pubsub_subscription.go",
        "pubsub_utils.go",
    ],
)

go_test(
    name = "pubsub_test",
    srcs = [
        "data_source_pubsub_subscription_test.go",
        "data_source_pubsub_topic_test.go",
        "iam_pubsub_subscription_test.go",
        "iam_pubsub_topic_test.go",
        "resource_pubsub_schema_test.go",
        "resource_pubsub_subscription_test.go",
        "resource_pubsub_topic_test.go",
    ],
)
