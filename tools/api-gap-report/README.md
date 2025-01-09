# API GAP Report Tool

This tool compares the fields between Google Cloud API resources and their
corresponding resources in the `terraform-provider-google`
and `terraform-provider-google-beta`. It generates a detailed report showing
the differences and can optionally save the API and Terraform component schemas
as JSON files.

## Features

- **Compare GCP API fields** with the corresponding Terraform fields.
- **Generate a detailed diff report** of the comparison.
- Option to **save API and Terraform component schemas** as JSON files.
- **Compare the newest report with an old one** to track changes over time.

## Installation

1. Go to API gap report tool path
    ```bash
    cd ./tools/api-gap-report
    ```

2. Create python virtual environment
   ```bash
   python -m venv venv
   ```

3. Activate python venv and install the necessary dependencies:
    ```bash
    source venv/bin/activate
    pip install -r requirements.txt
    ```

2. Ensure you have `terraform` installed and properly configured on your system.

## Usage

To use the tool, run the following command:

```bash
src/diff_report.py -h
```

### Required arguments:

* `-t TERRAFORM_CONFIG`, `--terraform_config` TERRAFORM_CONFIG
                        Path to the terraform config main.tf file
* `-c COMPONENT`, `--component COMPONENT`: The Terraform component that will be
                                          compared with the GCP API (e.g., Instance).

### Optional arguments

* `-d DIFF_REPORT`, `--diff_report DIFF_REPORT`: Path to the old report file that
                                                 will be compared with the newest
                                                 report.
* `-s, --save_file`: Save the API and Terraform component schemas as JSON files.
* `-v, --verbose`: Increase the log verbosity level.
* `-h, --help`: Show the help message and exit.

### Examples

Create simple report for google_compute_instance:

```bash
src/diff_report.py -c Instance -t /path/to/terraform/config
```

Create simple report for google_compute_instance and save schemas as JSON files
and compare with an old report:

```bash
src/diff_report.py -v -c Instance -t /path/to/terraform/config -s -d /path/to/old_report.yaml
```
