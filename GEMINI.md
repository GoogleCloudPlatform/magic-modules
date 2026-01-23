# Magic Modules Source Code

This directory contains the source code for the Google Cloud Terraform provider generator.

## Structure

- `mmv1/`: Resource definitions and custom code templates.
- `tpgtools/`: Tools for generating the provider.
- `scripts/`: Internal development scripts.

## Local Build & Test

Please refer to the root [GEMINI.md](../GEMINI.md) for detailed instructions on:
- Setting up the environment (`.agent/local.env`)
- Building the provider (`make provider`)
- Running acceptance tests (`go test`)
