# Assured Workloads V2 (DB-V2) Magic Modules Product

This directory defines the Terraform resources and data sources for Google Cloud Assured Workloads V2 (`assuredworkloadsv2`).

## Resources

- `google_assured_workloads_v2_settings`: Singleton resource managing organization-level settings (such as Organization-Level Personnel Controls, `olpc_mode`).
- `google_assured_workloads_v2_workload`: Core resource managing compliance workloads across GCP projects and folders with compliance frameworks, cloud control configurations, and CMEK bindings.

## Data Sources

- `google_assured_workloads_v2_settings`: Read-only inspection of organization settings.
- `google_assured_workloads_v2_workload`: Read-only inspection of an existing workload.

## V1 and V2 Coexistence Guide

Assured Workloads V1 (`google_assured_workloads_workload`) and Assured Workloads V2 (`google_assured_workloads_v2_workload`) interact with distinct control planes:
- V1 manages legacy workloads with fixed regulatory regimes and automatic project/folder creation.
- V2 decouples workload orchestration from underlying resource creation, supporting flexible binding to existing projects or folders (`resource_config`), granular `cloud_control_configs`, and organization singleton settings (`google_assured_workloads_v2_settings`).

Both V1 and V2 resources can coexist safely in the same GCP organization and Terraform root module.
