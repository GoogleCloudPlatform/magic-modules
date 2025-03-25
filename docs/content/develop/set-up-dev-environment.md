---
title: "Set up your development environment"
weight: 10
---

# Set up your development environment

Before you start adding or updating a Terraform resource using
[magic-modules]({{< ref "/" >}}), you must first set up your environment by
installing the necessary tools. This page explains the steps for setting up your
development environment.

1. [Install the gcloud CLI.](https://cloud.google.com/sdk/docs/install)
1. In the Google Cloud console, on the project selector page, select or
   [create a Google Cloud project](https://cloud.google.com/resource-manager/docs/creating-managing-projects).
   {{< hint info >}}

   **Note:** If you don't already have a project to use for testing changes to
   the Terraform providers, create a project instead of selecting an existing
   poject. After you finish these steps, you can delete the project, removing
   all resources associated with the project.

   {{< /hint >}}
   {{< button href="https://console.cloud.google.com/projectselector2/home/dashboard" >}}Go to project selector{{< /button >}}
1. Make sure that billing is enabled for your Google Cloud project. Learn how to
   [check if billing is enabled on a project](https://cloud.google.com/billing/docs/how-to/verify-billing-enabled).


1. [Install git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
1. [Install go](https://go.dev/doc/install)
1. Add the following values to your environment settings such as `.bashrc`:
   ```bash
   # Add GOPATH variable for convenience
   export GOPATH=$(go env GOPATH)
   # Add Go binaries to PATH
   export PATH=$PATH:$(go env GOPATH)/bin
   ```
1. Install goimports
   ```bash
   go install golang.org/x/tools/cmd/goimports@latest
   ```
1. [Install terraform](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli)
1. Clone the `magic-modules` repository
   ```bash
   cd ~
   git clone https://github.com/GoogleCloudPlatform/magic-modules.git
   ```

1. Run the following command from the root of your cloned `magic-modules` repository.
  
   ```bash
   cd magic-modules
   ./scripts/doctor
   ```
 
   Expected output if everything is installed properly:
 
   ```
   Check for go in path...
      found!
   Check for goimports in path...
      found!
   Check for git in path...
      found!
   Check for terraform in path...
      found!
   Check for make in path...
      found!
   ```

## What's next

+ [Learn how to add a resource]({{< ref "/develop/add-resource" >}})
+ [Learn how to add custom resource code]({{< ref "/develop/custom-code" >}})
+ [Learn how to add a datasource]({{< ref "/develop/add-handwritten-datasource" >}})
+ [Learn how to promote a resource to GA]({{< ref "/develop/promote-to-ga" >}})

