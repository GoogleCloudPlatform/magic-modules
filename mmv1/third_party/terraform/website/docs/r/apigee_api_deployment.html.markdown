---
subcategory: "Apigee"
description: |-
  Deploys a revision of a sharedflow.
---

# google_apigee_api_deployment

Deploys a revision of an api proxy.


To get more information about ApiDeployment, see:

* [API documentation](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.apis.revisions.deployments)
* How-to Guides
    * [apis.revisions.deployments](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.apis.revisions.deployments)

## Example Usage

```hcl
resource "google_apigee_api" "api_proxy" {
  name          = "proxy1"
  org_id        = var.org_id
  config_bundle = data.archive_file.bundle.output_path
}

resource "google_apigee_api_deployment" "api_proxy_deployment" {
  org_id        = var.org_id
  environment   = "my-environment"
  proxy_id      = google_apigee_api.api_proxy.name
  revision      = google_apigee_api.api_proxy.latest_revision_id
}
```

## Argument Reference

The following arguments are supported:


* `org_id` -
  (Required)
  The Apigee Organization associated with the API proxy

* `environment` -
  (Required)
  The resource ID of the environment.

* `proxy_id` -
  (Required)
  Name of the proxy to be deployed.

* `revision` -
  (Required)
  Revision of the proxy to be deployed.


- - -

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{org_id}}/environments/{{environment}}/apis/{{proxy_id}}/revisions/{{revision}}/deployments`


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


API Proxy Deployments can be imported using any of these accepted formats:

* `organizations/{{org_id}}/environments/{{environment}}/apis/{{proxy_id}}/revisions/{{revision}}/deployments`
* `{{org_id}}/{{environment}}/{{proxy_id}}/{{revision}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import ApiDeployment using one of the formats above. For example:

```tf
import {
  id = "organizations/{{org_id}}/environments/{{environment}}/apis/{{proxy_id}}/revisions/{{revision}}/deployments"
  to = google_apigee_api_deployment.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), ApiDeployment can be imported using one of the formats above. For example:

```
$ terraform import google_apigee_api_deployment.default organizations/{{org_id}}/environments/{{environment}}/apis/{{proxy_id}}/revisions/{{revision}}/deployments
$ terraform import google_apigee_api_deployment.default organizations/{{org_id}}/environments/{{environment}}/apis/{{proxy_id}}/revisions/{{revision}}
$ terraform import google_apigee_api_deployment.default {{org_id}}/{{environment}}/{{proxy_id}}/{{revision}}
```
