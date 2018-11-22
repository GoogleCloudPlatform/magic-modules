<!---
Note: This tutorial is meant for Google Cloud Shell, and can be opened by going to
https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/magic-modules&tutorial=tutorial.md
--->
# Magic Modules Tutorial

## Intro

This tutorial will walk you through the components that make up Magic Modules.

## api.yaml

Each product's api definition is stored in the magic-modules repo.

Let's open
<walkthrough-editor-open-file filePath="magic-modules/products/compute/api.yaml"
                              text="products/compute/api.yaml">
</walkthrough-editor-open-file>.

### Product metadata

The
<walkthrough-editor-select-regex filePath="magic-modules/products/compute/api.yaml"
                                 regex="!ruby/object:Api::Product"
                                 text="top section">
</walkthrough-editor-select-regex>
provides metadata about the API, such as name, scopes, and versions.

### Resources

Each `api.yaml` file contains a list of resources.
Let's
<walkthrough-editor-select-regex filePath="magic-modules/products/compute/api.yaml"
                                 regex="!ruby/object:Api::Resource"
                                 text="look at">
</walkthrough-editor-select-regex>
the first one.

This section contains data about the resource, such as its name, description, and URLs.

### Parameters + Properties

Each resource contains a list of
<walkthrough-editor-select-regex filePath="magic-modules/products/compute/api.yaml"
                                 regex="parameters:"
                                 text="URL parameters">
</walkthrough-editor-select-regex>
and
<walkthrough-editor-select-regex filePath="magic-modules/products/compute/api.yaml"
                                 regex="properties:"
                                 text="properties">
</walkthrough-editor-select-regex>
on the resource that a user might set when creating the resource, or access when reading it.

See the
<walkthrough-editor-open-file filePath="magic-modules/DEVELOPER.md"
                              text="Developer Guide">
</walkthrough-editor-open-file>
for more information about the values that can be set on parameters/properties.

## [provider].yaml

Within each product directory, each provider has its own `[provider].yaml` file to set information
specific to that provider.

Let's look at
<walkthrough-editor-open-file filePath="magic-modules/products/compute/terraform.yaml"
                              text="products/compute/terraform.yaml">
</walkthrough-editor-open-file>.

This file consists of information that is specific to Terraform, such as examples, modified
descriptions, custom code, and validation functions.

## Making changes

To add a new API or resource, the only files that need to be modified are `api.yaml`, each
`[provider].yaml`, and any custom code or provider-specific extras (such as Terraform example templates).

Let's actually make a change. Go back to
<walkthrough-editor-open-file filePath="magic-modules/products/compute/api.yaml"
                              text="products/compute/api.yaml">
</walkthrough-editor-open-file>
and change the description on the `Address` resource.

## Compiling magic-modules

Now, let's compile those changes.

Since we're running in cloud shell, this command will make sure we connect to GitHub via HTTPS
instead of SSH. You will probably not have to do this in your typical development environment.
```bash
git config --file=.gitmodules submodule.build/ansible.url https://github.com/modular-magician/ansible.git && git submodule sync
```

Now, initialize the submodules in order to get an up-to-date version of each provider.
Since we only changed the URL for Ansible, we'll only initialize that submodule.
```bash
git submodule update --init build/ansible
```

If you haven't already, run `bundle install` to make sure all ruby dependencies are available:
```bash
bundle install
```

Next, run the compiler:
```bash
bundle exec compiler -p products/compute -e ansible -o build/ansible
```

This command tells us to run the compiler for the compute API, and generate Ansible into the
`build/ansible` directory (where the submodule is).

Let's see our changes! Navigate to the Ansible submodule and run `git diff` to see what changed:
```bash
cd build/ansible && git diff
```

## Congratulations!

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

You've successfully made a change to a resource in Magic Modules.

When submitting PRs to Magic Modules, we rely on the [Magician](https://github.com/modular-magician)
to generate the PRs in each eligible repository and update the submodules.

To clear your submodules directory, run:
```bash
git submodule deinit --force --all
```
