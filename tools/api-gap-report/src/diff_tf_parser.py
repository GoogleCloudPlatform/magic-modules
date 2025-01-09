#!/usr/bin/env python3

import json
import os
import re
import subprocess
import time


class DiffTfParser:
    def _terraform_check(self):
        """
        Checks if Terraform is installed and accessible by running the
        `terraform --version` command.

        Returns:
            bool: `True` if Terraform is installed and the version command
                  runs successfully, `False` otherwise.
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        cmd_version = ["terraform", "--version"]
        self.log.debug("CMD: " + " ".join(cmd_version))

        try:
            p = subprocess.Popen(cmd_version, stdout=subprocess.PIPE)
            p.wait()
            if p.returncode != 0:
                self.log.error("Terraform version check failed!")
                return False
        except FileNotFoundError:
            self.log.error("Terraform command not available!")
            return False
        return True

    def _get_tf_schemas(self):
        """
        Retrieves the Terraform schemas using the
        `terraform providers schema -json` command.

        Returns:
            bool: `True` if the Terraform schemas are successfully retrieved
                  and parsed, `False` otherwise.
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        self.log.debug("Checking if Terraform is available")
        if not self._terraform_check():
            return False

        self.log.debug("Trying to get Terraform schemas")
        cmd_get_schemas = ["terraform", "providers", "schema", "-json"]
        p = subprocess.Popen(
            cmd_get_schemas,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        stdout, __ = p.communicate()
        terraform_stdout = stdout.decode("utf-8")

        if p.returncode != 0:
            self.log.error("Getting Terraform schemas failed!")
            return False

        if terraform_stdout == '{"format_version":"1.0"}\n':
            self.log.error(
                "No info about Terraform schemas! "
                "Check if Terraform configuration is available."
            )
            return False

        self.terraform_schemas = json.loads(terraform_stdout)

        return True

    def _camel_to_snake_string(self, camel):
        """
        Method converts camel string into snake case

        Args:
            camel (str): String that needs to be converted

        Returns:
            Snake case string
        """
        return re.sub(r'(?<!^)(?=[A-Z])', '_', camel).lower()

    def _snake_to_camel_string(self, snake):
        """
        Converts a snake_case string to camelCase.

        Args:
        snake (str): The string in snake_case format to be converted.

        Returns:
        str: The string converted to camelCase format.
        """
        parts = snake.split('_')
        return parts[0] + ''.join(word.capitalize() for word in parts[1:])

    def _snake_to_camel_schema(self, schema):
        """
        Recursively converts all dictionary keys in a nested structure from
        snake_case to camelCase.

        Args:
            schema (dict, list, any): The input schema which may be
                                      a dictionary, list, or other data types.
                                      Nested dictionaries or lists will also
                                      be processed recursively.

        Returns:
            dict or list: A new dictionary or list with all dictionary keys
                          converted from snake_case to camelCase.
                          Non-dictionary and non-list values remain unchanged.
        """
        if isinstance(schema, dict):
            return {
                self._snake_to_camel_string(key):
                self._snake_to_camel_schema(value)
                for key, value in schema.items()
            }
        elif isinstance(schema, list):
            return [self._snake_to_camel_schema(item) for item in schema]
        else:
            return schema

    def get_tf_component_schema(self, component, save_file=False):
        """
        Retrieves and processes the Terraform schema for a specific component.

        Args:
            component (str): The name of the Terraform component
                             (e.g., "google_compute_instance") for which the
                             schema is retrieved.
            save_file (bool, optional): If `True`, the retrieved schema will
                                        be saved to a JSON file. Defaults to
                                        `False`.

        Returns:
            bool: Returns `True` if the schema retrieval and processing were
                  successful, otherwise `False`.
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        if not hasattr(self, 'cwd'):
            print("Error: Current directory not set!")
            return False

        if not self._get_tf_schemas():
            self.log.error("Cannot get Terraform schemas!")
            return False

        try:
            self.component_tf_schema = self._snake_to_camel_schema(
                self.terraform_schemas[
                    "provider_schemas"][
                    "registry.terraform.io/hashicorp/google"][
                    "resource_schemas"][
                    f"google_compute_{self._camel_to_snake_string(component)}"]
            )
        except KeyError:
            self.log.error("The specified component not found in the schema.")
            return False

        if save_file:
            file_name = os.path.join(
                self.cwd,
                f"{component}_terraform_schema_{round(time.time())}.json"
            )
            self.log.debug(
                f"Saving {component} schema to json file {file_name}"
            )
            with open(file_name, "w") as f:
                json.dump(self.component_tf_schema, f, indent=2)
        return True

    def _get_nested_attributes(self, key, type_list: list):
        """
        Extracts nested attributes from a list of types and appends them
        to the `tf_field_list`.

        Args:
            key (str): The base key to which nested keys (if any) will
                       be appended.
            type_list (list): A list of types, where some elements may
                              contain nested dictionaries.

        Returns:
            None: The method directly modifies the `tf_field_list`
                  by appending keys or nested keys.
        """
        nested = False

        for type in type_list:
            if isinstance(type, list):
                nested = True
                for subkey in type[1].keys():
                    self.tf_field_list.append(key + "." + subkey)
                continue

        if not nested:
            self.tf_field_list.append(key)

    def _get_tf_field(self, key_origin, value_origin):
        """
        Recursively extracts and appends Terraform field keys to the
        `tf_field_list`.

        Args:
            key_origin (str): The base key (prefix) to prepend to the extracted
                              keys.
            value_origin (dict): The dictionary containing the Terraform block
                                 schema, including attributes and nested block
                                 types.

        Returns:
            None: The method directly modifies `tf_field_list` by appending
                  the extracted keys.
        """
        key_appendix = ''
        if key_origin != '':
            key_appendix = key_origin + '.'

        try:
            for key, value in value_origin["block"]["attributes"].items():
                if not isinstance(value["type"], list):
                    self.tf_field_list.append(key_appendix+key)
                    continue
                self._get_nested_attributes(key_appendix+key, value["type"])
        except KeyError:
            pass

        try:
            for key, value in value_origin["block"]["blockTypes"].items():
                self._get_tf_field(key_appendix+key, value)
        except KeyError:
            pass

    def get_tf_fields(self):
        """
        Extracts Terraform field keys from the component schema and populates
        `tf_field_list`.

        Returns:
            bool: True if Terraform fields were successfully extracted and
                  populated in `tf_field_list`;
                  False if an error occurred (e.g., missing schema or logger).
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        if not hasattr(self, 'component_tf_schema'):
            self.log.error("Terraform component schema not found!")
            return False

        self.tf_field_list = []
        self._get_tf_field('', self.component_tf_schema)

        if not self.tf_field_list:
            self.log.error("Failed to get Terraform component fields!")
            return False

        return True
