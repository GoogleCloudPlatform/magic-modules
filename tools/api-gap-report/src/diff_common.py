#!/usr/bin/env python3
import argparse
import logging

BOLD = "\033[1m"
RED = "\033[31m"
GREEN = "\033[32m"
YELLOW = "\033[33m"
BLUE = "\033[34m"
ENDC = '\033[0m'


class DiffCommon:
    def diff_cmdline(self):
        """
        Method parses commandline input and shows help message if needed.
        """
        description = (
            "Tool creates report describing differences between Google Cloud"
            " APIs fields and resources integrated in"
            " terraform-provider-google and terraform-provider-google-beta."
        )
        parser = argparse.ArgumentParser(description=description)
        parser.add_argument(
            "-t",
            "--terraform_config",
            help="Path to the terraform config main.tf file",
            required=True
        )
        parser.add_argument(
            "-c",
            "--component",
            help="Terraform component that will be compared with GCP API",
            required=True
        )
        parser.add_argument(
            "-d",
            "--diff_report",
            help=(
                "Old report file path that will be compared with the newest"
                " report"
            )
        )
        parser.add_argument(
            "-s",
            "--save_file",
            action="store_true",
            help="Save API and Terraform component schemas as a JSON files"
        )
        parser.add_argument(
            "-v",
            "--verbose",
            action="store_true",
            help="Increase logs verbosity level"
        )
        return parser

    def diff_log(self, verbose=False):
        """
        Method creates logging system for the tool.

                Create logger based on the verbosity level.

        Args:
            verbose (bool): If True, logging level is set to DEBUG; otherwise
            set it to INFO.
        """
        logging.basicConfig(
            level=logging.DEBUG if verbose else logging.INFO,
            format="%(asctime)s - %(levelname)s - %(message)s",
            datefmt="%Y-%m-%d %H:%M:%S"
        )

        self.log = logging.getLogger(__name__)
