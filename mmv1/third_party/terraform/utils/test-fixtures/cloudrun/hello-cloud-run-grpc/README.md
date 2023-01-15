# About text fixtures in this directory

The file contents in this directory are mainly for preparing a container image of a gRPC server to be used in the tests of Cloud Run services located in [`mmv1/third_party/terraform/tests/resource_cloud_run_v2_service_test.go`](./mmv1/third_party/terraform/tests/resource_cloud_run_v2_service_test.go).

Google Cloud provides convenience container images for Cloud Run Services (`gcr.io/cloudrun/hello` / `us-docker.pkg.dev/cloudrun/container/hello`) and Cloud Run Jobs (`gcr.io/cloudrun/job` / `us-docker.pkg.dev/cloudrun/container/job`). These prebuilt sample container images are used in Google Cloud Run tutorials.
These public container images, which require little configuration to pass health checks, are also useful when initializing Cloud Run Services and Jobs and bootstrapping Terraform projects.
However, the sample container images for Cloud Run Services only support REST API and Cloud Event as the invocation types, and does not support gRPC invocation.
Therefore we need to prepare container images of a gRPC server for our tests by ourselves.
