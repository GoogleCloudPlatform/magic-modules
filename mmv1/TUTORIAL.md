<!---
Note: This tutorial is meant for Google Cloud Shell, and can be opened by going to
https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/magic-modules&tutorial=TUTORIAL.md
--->
# Magic Modules Tutorial

<!-- TODO: analytics id? -->
<walkthrough-author name="danahoffman@google.com" tutorialName="Magic Modules Tutorial" repositoryUrl="https://github.com/GoogleCloudPlatform/magic-modules"></walkthrough-author>

## Intro

This tutorial will walk you through the components that make up Magic Modules.

## api.yaml

Each product's api definition is stored in the magic-modules repo.

Let's open
<walkthrough-editor-open-file filePath="products/pubsub/api.yaml"
                              text="products/pubsub/api.yaml">
</walkthrough-editor-open-file>.

### Product Metadata

The
<walkthrough-editor-select-regex filePath="products/pubsub/api.yaml"
                                 regex="!ruby/object:Api::Product"
                                 text="top section">
</walkthrough-editor-select-regex>
provides metadata about the API, such as name, scopes, and versions.

### Resources

Each `api.yaml` file contains a list of resources. A resource is an item in that product,
such as a PubSub Subscription, a Compute Instance, or a GKE Cluster.
Let's
<walkthrough-editor-select-regex filePath="products/pubsub/api.yaml"
                                 regex="!ruby/object:Api::Resource"
                                 text="look at">
</walkthrough-editor-select-regex>
the first one.

This section contains data about the resource, such as its name, description, and URLs.

### Properties

Each resource contains a list of
<walkthrough-editor-select-regex filePath="products/pubsub/api.yaml"
                                 regex="properties:"
                                 text="properties">
</walkthrough-editor-select-regex>
on the resource that a user might set when creating the resource, or access when reading it.

See the [property type fields](https://github.com/GoogleCloudPlatform/magic-modules/blob/main/mmv1/api/resource.rb#L22)
for more information about the values that can be set on properties.

All of this information comes from the PubSub Subscription [REST API docs](https://cloud.google.com/pubsub/docs/reference/rest/v1/projects.subscriptions)

## [provider].yaml

Within each product directory, each provider has its own `[provider].yaml` file to set information
specific to that provider.

Let's look at
<walkthrough-editor-open-file filePath="products/pubsub/ansible.yaml"
                              text="products/pubsub/ansible.yaml">
</walkthrough-editor-open-file>.

This file consists of information that is specific to Ansible, like Ansible version numbers,
helper code, and additional files to include.

## Making Changes

To add a new API or resource, the only files that need to be modified are `api.yaml`, each
`[provider].yaml`, and any custom code or provider-specific extras.

Let's actually make a change. Go back to
<walkthrough-editor-open-file filePath="products/pubsub/api.yaml"
                              text="products/pubsub/api.yaml">
</walkthrough-editor-open-file>

We're going to add in the Topic Resource now.

## Compiling magic-modules

Now, let's compile those changes.

Since we're running in Cloud Shell, this command will make sure we connect to GitHub via HTTPS
instead of SSH. You will probably not have to do this in your typical development environment.

Run the compiler:
```bash
ruby compiler.rb -p products/pubsub -e ansible -o build/ansible
```

This command tells us to run the compiler for the pubsub API, and generate Ansible into the
`build/ansible/plugins/modules` directory.

Let's see our changes! Navigate to the Ansible folder
```bash
cd build/ansible/plugins/modules/
// view gcp_pubsub_topic.py in editor of your choice.
```

## Congratulations!

<walkthrough-conclusion-trophy></walkthrough-conclusion-trophy>

You've successfully made a change to a resource in Magic Modules.
