---
title: "Update dependencies"
weight: 300
aliases:
  - /docs/update-dependencies
---

# Update provider dependencies

The Magic Modules repository does not contain a complete Go module, preventing the use of automated tooling like `go get` from that repository. To add or update provider dependencies, use standard Go tooling to update the dependency in the `terraform-provider-google` provider and copy the results to the upstream files in Magic Modules.

Below are the steps you can follow to make the change:

1. Navigate to the local `google` provider directory:
```bash
cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
``` 
2. Open the [`go.mod`](https://github.com/hashicorp/terraform-provider-google/blob/main/go.mod) file and add the new entries or modify the versions of existing entries as needed
3. Update dependencies using either of the following methods
   * run the following commands to update all dependencies: 
   ```bash
   go get
   go mod tidy
   ```
   * Alternatively, update a specific package to a desired version:
   ```bash
   go get google.golang.org/api@v0.105.0 
   go mod tidy
   ```
4. Copy the contents of the updated `go.mod` and `go.sum` file into [`mmv1/third_party/terraform/go.mod.erb`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/go.mod.erb) and [`mmv1/third_party/terraform/go.sum`](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/third_party/terraform/go.sum) in the `magic-modules` respectively. Ensure `<% autogen_exception -%>` is still at the top of the file afterwards
5. Create a pull request containing _only_ the dependency changes.
6. Resolve any errors caused by the dependency upgrade.
7. After the dependency pull request is merged, continue with a separate pull request for any additional changes that were blocked by the upgrade.