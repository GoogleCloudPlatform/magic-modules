#!/usr/bin/env python3
"""Unit tests for intake_failure_helper.py."""

import json
import os
import subprocess
import tempfile
import unittest
from unittest import mock

import intake_failure_helper as helper


class TestIntakeFailureHelper(unittest.TestCase):

    def test_is_permission_error_word_boundaries(self):
        self.assertTrue(helper.is_permission_error("ERROR: 403 Forbidden"))
        self.assertTrue(helper.is_permission_error("HTTP 401 Unauthorized"))
        self.assertTrue(
            helper.is_permission_error("User does not have storage.objects.get access")
        )
        # Should NOT match 403 or 401 embedded inside project numbers or IDs
        self.assertFalse(
            helper.is_permission_error("Not found in project-1234039 or id-94012")
        )

    def test_validate_test_name(self):
        self.assertEqual(
            helper.validate_test_name("TestAccRedisCluster_basic"),
            "TestAccRedisCluster_basic",
        )
        for invalid in [
            "",
            "NotTestAcc_basic",
            "TestAccFoo; rm -rf /",
            "TestAccFoo$(id)",
            "TestAccFoo/../bar",
        ]:
            with self.assertRaises(ValueError, msg=f"Should reject {invalid!r}"):
                helper.validate_test_name(invalid)

    def test_validate_provider(self):
        self.assertEqual(helper.validate_provider("ga"), "ga")
        self.assertEqual(helper.validate_provider("BETA"), "beta")
        self.assertEqual(helper.validate_provider(" both "), "both")
        with self.assertRaises(ValueError):
            helper.validate_provider("alpha")

    def test_normalize_and_validate_gcs_uri(self):
        self.assertEqual(
            helper.normalize_and_validate_gcs_uri(
                "gs://nightly-test-data/test-errors/ga/2026-09-23/TestAccFoo.txt"
            ),
            "gs://nightly-test-data/test-errors/ga/2026-09-23/TestAccFoo.txt",
        )
        self.assertEqual(
            helper.normalize_and_validate_gcs_uri(
                "https://storage.googleapis.com/nightly-test-data/test-errors/beta/TestAccFoo.txt"
            ),
            "gs://nightly-test-data/test-errors/beta/TestAccFoo.txt",
        )
        self.assertEqual(
            helper.normalize_and_validate_gcs_uri(
                "https://storage.cloud.google.com/nightly-test-data/debug/ga/TestAccFoo.log"
            ),
            "gs://nightly-test-data/debug/ga/TestAccFoo.log",
        )
        for invalid in [
            "gs://untrusted-bucket/test-errors/ga/TestAccFoo.txt",
            "https://storage.googleapis.com/untrusted-bucket/TestAccFoo.txt",
            "gs://nightly-test-data/../secret.txt",
            "gs://nightly-test-data/test-errors;id/foo.txt",
            "gs://nightly-test-data/test-errors/$(whoami).txt",
        ]:
            with self.assertRaises(ValueError, msg=f"Should reject {invalid!r}"):
                helper.normalize_and_validate_gcs_uri(invalid)

    def test_validate_github_issue_url(self):
        url = "https://github.com/hashicorp/terraform-provider-google/issues/28244"
        self.assertEqual(helper.validate_github_issue_url(url), url)
        for invalid in [
            "https://github.com/hashicorp/terraform-provider-google-beta/issues/28244",
            "https://github.com/evil/repo/issues/1",
            "https://github.com/hashicorp/terraform-provider-google/issues/28244;id",
        ]:
            with self.assertRaises(ValueError, msg=f"Should reject {invalid!r}"):
                helper.validate_github_issue_url(invalid)

    def test_validate_local_path_enforces_base_dir(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            valid_file = os.path.join(tmpdir, "test_output.log")
            with open(valid_file, "w", encoding="utf-8") as f:
                f.write("sample error")

            resolved = helper.validate_local_path(valid_file, base_dir=tmpdir)
            self.assertEqual(resolved, os.path.realpath(valid_file))

            with tempfile.NamedTemporaryFile() as outside_file:
                with self.assertRaises(ValueError):
                    helper.validate_local_path(outside_file.name, base_dir=tmpdir)

    def test_extract_gcs_links_from_issue_beta_in_test_name(self):
        body = """
        ga error message: https://storage.googleapis.com/nightly-test-data/test-errors/ga/TestAccContainerCluster_withEnableKubernetesBetaAPIs.txt
        ga debug log: https://storage.googleapis.com/nightly-test-data/logs/ga/TestAccContainerCluster_withEnableKubernetesBetaAPIs.log
        beta debug log: https://storage.googleapis.com/nightly-test-data/logs/beta/TestAccContainerCluster_withEnableKubernetesBetaAPIs.log
        """
        error_links, debug_links = helper.extract_gcs_links_from_issue(body)
        self.assertEqual(
            error_links,
            {
                "ga": "gs://nightly-test-data/test-errors/ga/TestAccContainerCluster_withEnableKubernetesBetaAPIs.txt"
            },
        )
        self.assertEqual(
            debug_links,
            {
                "ga": "gs://nightly-test-data/logs/ga/TestAccContainerCluster_withEnableKubernetesBetaAPIs.log",
                "beta": "gs://nightly-test-data/logs/beta/TestAccContainerCluster_withEnableKubernetesBetaAPIs.log",
            },
        )

    def test_validate_and_extract_issue_requires_test_failure_label(self):
        non_test_issue = {
            "title": "Failing test(s): TestAccRedisCluster_basic",
            "body": "GA: 100%",
            "labels": [{"name": "bug"}],
        }
        with self.assertRaisesRegex(ValueError, "test-failure"):
            helper.validate_and_extract_issue(non_test_issue)

        valid_beta_issue = {
            "title": "Failing test(s): TestAccRedisCluster_basic",
            "body": (
                "Failure rates: GA: 0%, Beta: 100%\n"
                "https://storage.googleapis.com/nightly-test-data/test-errors/beta/TestAccRedisCluster_basic.txt"
            ),
            "labels": [{"name": "test-failure-100"}],
        }
        extracted = helper.validate_and_extract_issue(valid_beta_issue)
        self.assertEqual(extracted["test_name"], "TestAccRedisCluster_basic")
        self.assertEqual(extracted["target_provider"], "beta")

    def test_run_cli_rejects_conflicting_issue_url_and_local_or_gcs_flags(self):
        with self.assertRaisesRegex(ValueError, "cannot be combined"):
            helper.run_cli(
                [
                    "--issue-url",
                    "https://github.com/hashicorp/terraform-provider-google/issues/28244",
                    "--local-log",
                    "some.log",
                ]
            )
        with self.assertRaisesRegex(ValueError, "cannot be combined"):
            helper.run_cli(
                [
                    "--issue-url",
                    "https://github.com/hashicorp/terraform-provider-google/issues/28244",
                    "--gcs-error-uri",
                    "gs://nightly-test-data/test-errors/ga/foo.txt",
                ]
            )

    @mock.patch("intake_failure_helper.fetch_gcs_content")
    @mock.patch("intake_failure_helper.fetch_issue_payload")
    def test_run_cli_issue_url_preserves_extracted_beta_provider_with_test_name(
        self, mock_fetch_issue, mock_fetch_gcs
    ):
        mock_fetch_issue.return_value = {
            "title": "Failing test(s): TestAccRedisCluster_basic",
            "body": (
                "GA: 0%\nBeta: 100%\n"
                "gs://nightly-test-data/test-errors/beta/TestAccRedisCluster_basic.txt"
            ),
            "labels": [{"name": "test-failure"}],
        }
        mock_fetch_gcs.return_value = "Error: beta test failure"

        with tempfile.TemporaryDirectory() as tmpdir:
            res = helper.run_cli(
                [
                    "--issue-url",
                    "https://github.com/hashicorp/terraform-provider-google/issues/28244",
                    "--test-name",
                    "TestAccRedisCluster_basic",
                    "--output-dir",
                    tmpdir,
                ],
                base_dir=tmpdir,
            )
            payload = res["normalized_failure_payload"]
            self.assertEqual(payload["target_provider"], "beta")
            with open(payload["error_log_file"], "r", encoding="utf-8") as f:
                self.assertIn("Error: beta test failure", f.read())

    @mock.patch("intake_failure_helper.subprocess.run")
    def test_run_cli_local_log_with_parse_debug_log_writes_error_and_nested_dir(
        self, mock_run
    ):
        with tempfile.TemporaryDirectory() as tmpdir:
            local_log = os.path.join(tmpdir, "local_debug.log")
            with open(local_log, "w", encoding="utf-8") as f:
                f.write("=== RUN   TestAccFoo_basic\n--- FAIL: TestAccFoo_basic\n")

            nested_out = os.path.join(
                tmpdir, "TestAccFoo_basic", "TestAccFoo_basic_1700000000"
            )
            mock_run.return_value = subprocess.CompletedProcess(
                args=[],
                returncode=0,
                stdout=f"Extracted API timeline and errors to {nested_out}/\n",
                stderr="",
            )

            res = helper.run_cli(
                [
                    "--test-name",
                    "TestAccFoo_basic",
                    "--local-log",
                    local_log,
                    "--parse-debug-log",
                    "--output-dir",
                    tmpdir,
                ],
                base_dir=tmpdir,
            )
            payload = res["normalized_failure_payload"]
            self.assertEqual(payload["parsed_logs_dir"], f"{nested_out}/")
            with open(payload["error_log_file"], "r", encoding="utf-8") as f:
                self.assertIn("--- FAIL: TestAccFoo_basic", f.read())


if __name__ == "__main__":
    unittest.main()
