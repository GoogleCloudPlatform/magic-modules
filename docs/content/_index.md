---
title: "Home"
weight: 0
type: "docs"
date: 2022-11-14T09:50:49-08:00
---


# Magic Modules

Magic Modules is a code generator and CI system that's used to develop the Terraform providers
for Google Platform, [`google`](https://github.com/hashicorp/terraform-provider-google) (or TPG) and
[`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) (or TPGB).

Magic Modules allows contributors to make changes against a single codebase and develop both
provider versions simultaneously. After sending a pull request against this repository, the
`modular-magician` robot user will manage (most of) the heavy lifting from generating a
complete output, running presubmit tests, and updating the providers following your
change.

## Getting started

Check out the [setup guide](/magic-modules/docs/getting-started/setup/) for information on how to set up your environment.

## Other Resources

* [Extending Terraform](https://www.terraform.io/plugin)
   * [How Terraform Works](https://www.terraform.io/plugin/how-terraform-works)
   * [Writing Custom Providers / Calling APIs with Terraform Providers](https://learn.hashicorp.com/collections/terraform/providers?utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS)
* [Terraform Glossary](https://www.terraform.io/docs/glossary)
