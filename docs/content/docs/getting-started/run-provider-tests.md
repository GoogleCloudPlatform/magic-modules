---
title: "Run provider tests"
weight: 30
---

# Run provider tests locally

{{< hint info >}}
**Note:** If you want to test changes you've made in Magic Modules, you need to first [generate](/magic-modules/docs/getting-started/generate-providers/) the provider you want to test.
{{< /hint >}}

## Setup

Tests generally assume the following environment variables must be set in order to run tests:

```
GOOGLE_PROJECT
GOOGLE_CREDENTIALS|GOOGLE_CLOUD_KEYFILE_JSON|GCLOUD_KEYFILE_JSON|GOOGLE_USE_DEFAULT_CREDENTIALS
GOOGLE_REGION
GOOGLE_ZONE
```

Note that the credentials you provide must be granted wide permissions on the specified project. These tests provision real resources, and require permission in order to do so. Most developers on the team grant their test service account `roles/editor` or `roles/owner` on their project. Additionally, to ensure that your tests are performed in a region and zone with wide support for GCP features, `GOOGLE_REGION` should be set to `us-central1` and `GOOGLE_ZONE` to `us-central1-a`.

Additional variable may be required for other tests, and should get flagged when running them by Go skipping the test and flagging in the output it was skipped, with a skip message explaining why. The most typical extra values required are those required for project creation:

```
GOOGLE_ORG
GOOGLE_BILLING_ACCOUNT
```

## Run unit tests

Unit tests (that is, tests that do not interact with the GCP API) are very fast and you can generally run them all if you have changed any of them:

```bash
make test
```

## Run acceptance tests

You can run tests against the provider you generated in the `OUTPUT_PATH` location. When running tests, specify which to run using `TESTARGS`, such as:

```bash
# for ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic'

# for beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool_basic'
```

TESTARGS allows you to pass [testing flags](https://pkg.go.dev/cmd/go#hdr-Testing_flags) to `go test`. The most important is `-run`, which allows you to limit the tests that get run. There are 2000+ tests, and running all of them takes over 9 hours and requires a lot of GCP quota.

`-run` is regexp-like, so multiple tests can be run in parallel by specifying a common substring of those tests (for example, `TestAccContainerNodePool` to run all node pool tests).

## Debugging tests

You can [increase your test verbosity](https://www.terraform.io/docs/internals/debugging.html)  and redirect the output to a log file for analysis. This is often helpful in debugging issues.

```bash
# for ga provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
TF_LOG=TRACE make testacc TEST=./google TESTARGS='-run=TestAccContainerNodePool_basic' > output.log

# for beta provider
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google-beta
TF_LOG=TRACE make testacc TEST=./google-beta TESTARGS='-run=TestAccContainerNodePool_basic' > output.log
```

You can also debug tests with [Delve](https://github.com/go-delve/delve):

```bash
# Navigate to the google package within your local GCP Terraform provider Git clone.
cd $GOPATH/src/github.com/terraform-providers/terraform-provider-google/google

# Execute the dlv command to launch the test.
# Note that the --test.run flag uses the same regexp matching as go test --run.
TF_ACC=1 dlv test -- --test.v --test.run TestAccComputeRegionBackendService_withCdnPolicy
Type 'help' for list of commands.
(dlv) b google.TestAccComputeRegionBackendService_withCdnPolicy
Breakpoint 1 set at 0x1de072b for github.com/terraform-providers/terraform-provider-google/google.TestAccComputeRegionBackendService_withCdnPolicy() ./resource_compute_region_backend_service_test.go:540
(dlv) c
=== RUN   TestAccComputeRegionBackendService_withCdnPolicy
> github.com/terraform-providers/terraform-provider-google/google.TestAccComputeRegionBackendService_withCdnPolicy() ./resource_compute_region_backend_service_test.go:540 (hits goroutine(7):1 total:1) (PC: 0x1de072b)
   535:                         },
   536:                 },
   537:         })
   538: }
   539:
=> 540: func TestAccComputeRegionBackendService_withCdnPolicy(t *testing.T) {
   541:         t.Parallel()
   542:
   543:         var svc compute.BackendService
   544:         resource.Test(t, resource.TestCase{
   545:                 PreCheck:     func() { testAccPreCheck(t) },
(dlv)
```

## Testing with different `terraform` versions

Tests will use whatever version of the `terraform` binary is found on your path. To test with multiple versions of `terraform` core, you must run the tests multiple times with different versions. You can use [`tfenv`](https://github.com/tfutils/tfenv) to manage your system `terraform` versions.