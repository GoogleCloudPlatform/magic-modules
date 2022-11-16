<img src="docs/static/magic-modules.svg" alt="Magic Modules Logo" width="300" align="right" />

# Magic Modules

Magic Modules is a code generator and CI system that's used to develop the Terraform providers
for Google Platform, [`google`](https://github.com/hashicorp/terraform-provider-google) (or TPG) and
[`google-beta`](https://github.com/hashicorp/terraform-provider-google-beta) (or TPGB).

Magic Modules allows contributors to make changes against a single codebase and develop both
provider versions simultaneously. After sending a pull request against this repository, the
`modular-magician` robot user will manage (most of) the heavy lifting from generating a
complete output, running presubmit tests, and updating the providers following your
change.

For information on how to use or contribute to Magic Modules, see [the documentation](https://googlecloudplatform.github.io/magic-modules).
