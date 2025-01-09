#!/usr/bin/env python3

import json
import jsonref
import requests
import os
import time

from diff_config import DISCOVERY_DOC_URL


class DiffApiParser:
    def _get_api_schemas(self):
        """
        Retrieves and processes the API schemas from the discovery document.

        Returns:
            bool:
                - `True` if the API schemas are successfully retrieved,
                  decoded, and dereferenced.
                - `False` if there is an error at any step of the process
                  (missing logger, failed JSON decoding, or dereferencing
                  issues).
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        self.log.debug(
            f"Trying to get discovery doc from: {DISCOVERY_DOC_URL}"
        )
        discovery_response = requests.get(DISCOVERY_DOC_URL)
        try:
            self.log.debug("Trying to decode JSON file")
            ref_api_schemas = discovery_response.json()
        except json.decoder.JSONDecodeError:
            self.log.error("Response does not contain the JSON file!")
            return False

        if not ref_api_schemas:
            self.log.error("Unknown error during parsing discovery doc!")
            return False

        try:
            self.log.debug("Trying to dereference API schemas")
            self._api_schemas = jsonref.JsonRef.replace_refs(
                ref_api_schemas,
                base_uri=DISCOVERY_DOC_URL,
                jsonschema=True
            ).get("schemas", {})
        except jsonref.JsonRefError:
            self.log.error("Dereferencing API schema has failed!")
            return False

        if not self._api_schemas:
            self.log.error("Unknown error during dereferencing API schema!")

        return True

    def get_api_component_schema(self, component, save_file=False):
        """
        Retrieves the API schema for a specified component and optionally
        saves it to a file.

        Args:
            component (str): The name of the component for which the API schema
                             is to be retrieved.
            save_file (bool, optional): If `True`, the schema is saved to
                                        a JSON file. Defaults to `False`.

        Returns:
            bool:
                - `True` if the API schema is successfully retrieved and
                  optionally saved to a file.
                - `False` if any required attributes are missing, or the
                  component schema cannot be found.
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        if not hasattr(self, 'cwd'):
            print("Error: Current directory not set!")
            return False

        if not self._get_api_schemas():
            self.log.error("Cannot get GCP API schemas!")
            return False

        self.log.debug(f"Getting {component} schema")
        self.component_api_schema = self._api_schemas.get(component)
        if not self.component_api_schema:
            self.log.error("The specified component not found in the schema!")
            return False

        if save_file:
            file_name = os.path.join(
                self.cwd,
                f"{component}_api_schema_{round(time.time())}.json"
            )
            self.log.debug(
                f"Saving {component} schema to json file {file_name}"
            )
            with open(file_name, "w") as f:
                json.dump(self.component_api_schema, f, indent=2)
        return True

    def _get_api_field(self, key_origin, value_origin):
        """
        Recursively extracts API field keys from a given schema and appends
        them to `self.api_field_list`.

        Args:
            key_origin (str): The key path of the current field being processed
                              (used for nested fields).
            value_origin (dict): The current level of the schema to process,
                                 typically a dictionary containing
                                 "properties" or "items" keys with their
                                 respective values.
        """
        nested = False
        key_appendix = ''
        if key_origin != '':
            key_appendix = key_origin + '.'

        try:
            for key, value in value_origin["properties"].items():
                if key != "properties":
                    self._get_api_field(key_appendix+key, value)
                else:
                    self._get_api_field(key_appendix, value)
            nested = True
        except KeyError:
            pass

        try:
            for key, value in value_origin["items"]["properties"].items():
                self._get_api_field(key_appendix+key, value)
            nested = True
        except KeyError:
            pass

        if not nested and key_origin not in self.api_field_list:
            self.api_field_list.append(key_origin)

    def get_api_fields(self):
        """
        Retrieves the API fields from the component API schema.

        Returns:
            bool:
                - `True` if the API fields are successfully retrieved and
                  `self.api_field_list` is populated.
                - `False` if an error occurs (e.g., missing attributes or
                  empty field list).
        """
        if not hasattr(self, 'log'):
            print("Error: Logger not found!")
            return False

        if not hasattr(self, 'component_api_schema'):
            self.log.error("API component schema not found!")
            return False

        self.api_field_list = []
        self._get_api_field('', self.component_api_schema)

        if not self.api_field_list:
            self.log.error("Failed to get API component fields!")
            return False

        return True
