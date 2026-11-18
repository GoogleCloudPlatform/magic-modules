#!/usr/bin/env python3
"""Unit tests for tf_debug_parser.py."""

import contextlib
import io
import json
import os
import subprocess
import sys
import tempfile
import unittest

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import tf_debug_parser  # noqa: E402


SAMPLE_DEBUG_LOG = """=== RUN   TestAccExampleResource_update
2026/10/28 08:00:00 [INFO] Terraform is using this identity: test-runner@example.iam.gserviceaccount.com
2026/10/28 08:00:01 [DEBUG] sdk.helper_resource: Starting TestStep: test_name=TestAccExampleResource_update, test_step_number=1
2026/10/28 08:00:02 [DEBUG] Google API Request Details:
---[ REQUEST ]---------------------------------------
POST /v1/projects/my-proj/locations/us-central1/examples?alt=json&exampleId=tf-test-ex1&prettyPrint=false HTTP/1.1
Host: example.googleapis.com
Content-Type: application/json

{
  "displayName": "initial",
  "enabled": true
}
-----------------------------------------------------
2026/10/28 08:00:03 [DEBUG] Google API Response Details:
---[ RESPONSE ]--------------------------------------
HTTP/1.1 200 OK
Content-Type: application/json

{
  "name": "projects/my-proj/locations/us-central1/examples/tf-test-ex1",
  "displayName": "initial",
  "enabled": true
}
-----------------------------------------------------
2026/10/28 08:00:04 [DEBUG] sdk.helper_resource: Finished TestStep: test_name=TestAccExampleResource_update, test_step_number=1
2026/10/28 08:00:05 [DEBUG] sdk.helper_resource: Starting TestStep: test_name=TestAccExampleResource_update, test_step_number=2
Terraform will perform the following actions:

  # google_example_resource.primary will be updated in-place
  ~ resource "google_example_resource" "primary" {
      ~ display_name = "initial" -> "updated"
      ~ enabled      = true -> false
    }

Plan: 0 to add, 1 to change, 0 to destroy.
2026/10/28 08:00:06 [DEBUG] Google API Request Details:
---[ REQUEST ]---------------------------------------
PATCH /v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?alt=json&updateMask=displayName%2Cenabled HTTP/1.1
Host: example.googleapis.com
Content-Type: application/json

{
  "displayName": "updated"
}
-----------------------------------------------------
2026/10/28 08:00:07 [DEBUG] Google API Response Details:
---[ RESPONSE ]--------------------------------------
HTTP/1.1 200 OK
Content-Type: application/json

{
  "name": "projects/my-proj/locations/us-central1/examples/tf-test-ex1",
  "displayName": "updated",
  "enabled": true
}
-----------------------------------------------------
    resource_example_test.go:42: Step 2/2 error: After applying this test step, the plan was not empty.
2026/10/28 08:00:08 [DEBUG] sdk.helper_resource: Finished TestStep: test_name=TestAccExampleResource_update, test_step_number=2
2026/10/28 08:00:09 [DEBUG] Google API Request Details:
---[ REQUEST ]---------------------------------------
DELETE /v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?alt=json&force=true HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
2026/10/28 08:00:10 [DEBUG] Google API Response Details:
---[ RESPONSE ]--------------------------------------
HTTP/1.1 200 OK
Content-Type: application/json

{}
-----------------------------------------------------
2026/10/28 08:00:11 [DEBUG] sdk.helper_resource: Calling TestCase CheckDestroy
2026/10/28 08:00:12 [DEBUG] Google API Request Details:
---[ REQUEST ]---------------------------------------
GET /v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?alt=json HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
2026/10/28 08:00:13 [DEBUG] Google API Response Details:
---[ RESPONSE ]--------------------------------------
HTTP/1.1 404 Not Found
Content-Type: application/json

{
  "error": {
    "code": 404,
    "message": "Resource not found"
  }
}
-----------------------------------------------------
--- FAIL: TestAccExampleResource_update (13.42s)
"""


class TestTFDebugParser(unittest.TestCase):

    def setUp(self):
        self.temp_dir = tempfile.TemporaryDirectory()
        self.log_path = os.path.join(self.temp_dir.name, "test_debug.log")
        with open(self.log_path, "w", encoding="utf-8") as f:
            f.write(SAMPLE_DEBUG_LOG)

    def tearDown(self):
        self.temp_dir.cleanup()

    def test_json_cli_and_export_json(self):
        script_path = os.path.join(
            os.path.dirname(os.path.abspath(__file__)), "tf_debug_parser.py"
        )
        proc = subprocess.run(
            [sys.executable, script_path, self.log_path, "--json"],
            capture_output=True,
            text=True,
            check=True,
        )
        payload = json.loads(proc.stdout)

        self.assertEqual(payload["test_name"], "TestAccExampleResource_update")
        self.assertEqual(
            [s["step_number"] for s in payload["steps"]],
            [1, 2],
        )
        self.assertEqual(len(payload["diffs"]), 1)
        self.assertGreaterEqual(len(payload["failures"]), 1)

        requests = [e for e in payload["timeline"] if e["type"] == "request"]
        self.assertEqual([r["method"] for r in requests], ["POST", "PATCH", "DELETE", "GET"])

        # Step 1 POST preserves url_param_only query params
        post_req = requests[0]
        self.assertEqual(post_req["test_step"], 1)
        self.assertEqual(post_req["phase"], "step_1")
        self.assertEqual(
            post_req["url_path"],
            "/v1/projects/my-proj/locations/us-central1/examples",
        )
        self.assertIn("exampleId=tf-test-ex1", post_req["url_full"])
        self.assertEqual(post_req["query_params"]["exampleId"], ["tf-test-ex1"])

        # Step 2 PATCH preserves URL-encoded ?updateMask=displayName%2Cenabled
        patch_req = requests[1]
        self.assertEqual(patch_req["test_step"], 2)
        self.assertEqual(patch_req["phase"], "step_2")
        self.assertEqual(
            patch_req["url_path"],
            "/v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1",
        )
        self.assertEqual(
            patch_req["url_full"],
            "/v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?alt=json&updateMask=displayName%2Cenabled",
        )
        self.assertEqual(
            patch_req["query_params"]["updateMask"],
            ["displayName,enabled"],
        )
        self.assertEqual(patch_req["update_mask"], ["displayName", "enabled"])

        # Post-test DELETE and CheckDestroy GET are tagged in destroy phase
        delete_req = requests[2]
        self.assertIsNone(delete_req["test_step"])
        self.assertEqual(delete_req["phase"], "destroy")
        self.assertEqual(delete_req["query_params"]["force"], ["true"])

        check_destroy_get = requests[3]
        self.assertIsNone(check_destroy_get["test_step"])
        self.assertEqual(check_destroy_get["phase"], "destroy")

    def test_summary_shows_steps_and_update_mask(self):
        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            rc = tf_debug_parser.main([self.log_path, "--summary"])
        self.assertEqual(rc, 0)
        output = stdout.getvalue()

        self.assertIn(
            "[Step 1] REQ: POST /v1/projects/my-proj/locations/us-central1/examples?exampleId=tf-test-ex1",
            output,
        )
        self.assertIn(
            "[Step 2] REQ: PATCH /v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?updateMask=displayName,enabled",
            output,
        )
        self.assertIn(
            "[Destroy] REQ: DELETE /v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?force=true",
            output,
        )
        # Boilerplate ?alt=json should not clutter human summary
        self.assertNotIn("alt=json", output)

    def test_extract_dir_writes_outline_payloads_and_failures(self):
        extract_root = os.path.join(self.temp_dir.name, "extracted")
        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            rc = tf_debug_parser.main(
                [self.log_path, "--extract-dir", extract_root]
            )
        self.assertEqual(rc, 0)

        subdirs = os.listdir(extract_root)
        self.assertEqual(len(subdirs), 1)
        out_dir = os.path.join(extract_root, subdirs[0])

        outline_path = os.path.join(out_dir, "outline.txt")
        self.assertTrue(os.path.exists(outline_path))
        with open(outline_path, "r", encoding="utf-8") as f:
            outline = f.read()

        self.assertIn(">>> STARTING TEST STEP 1 <<<", outline)
        self.assertIn(">>> STARTING TEST STEP 2 <<<", outline)
        self.assertIn(">>> POST-TEST DESTROY <<<", outline)
        self.assertIn("?updateMask=displayName,enabled", outline)
        self.assertIn("?exampleId=tf-test-ex1", outline)
        self.assertIn("?force=true", outline)

        patch_json_path = os.path.join(out_dir, "04_REQUEST_PATCH.json")
        self.assertTrue(os.path.exists(patch_json_path))
        with open(patch_json_path, "r", encoding="utf-8") as f:
            patch_event = json.load(f)
        self.assertEqual(patch_event["update_mask"], ["displayName", "enabled"])
        self.assertEqual(patch_event["test_step"], 2)

        failures_path = os.path.join(out_dir, "failures.json")
        self.assertTrue(os.path.exists(failures_path))
        with open(failures_path, "r", encoding="utf-8") as f:
            failures = json.load(f)
        self.assertTrue(any(item["test_step"] == 2 for item in failures))

    def test_missing_file_returns_nonzero(self):
        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            rc = tf_debug_parser.main(
                [os.path.join(self.temp_dir.name, "nonexistent.log")]
            )
        self.assertEqual(rc, 1)

    def test_aborted_step_transitions_to_destroy_without_finished_marker(self):
        aborted_log = """=== RUN   TestAccExampleResource_abort
2026/10/28 08:00:01 [DEBUG] sdk.helper_resource: Starting TestStep: test_name=TestAccExampleResource_abort, test_step_number=1
2026/10/28 08:00:02 [ERROR] sdk.helper_resource: Unexpected error: Error: invalid value, test_step_number=1
---[ REQUEST ]---------------------------------------
DELETE /v1/projects/my-proj/locations/us-central1/examples/tf-test-ex1?alt=json HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
"""
        aborted_path = os.path.join(self.temp_dir.name, "aborted.log")
        with open(aborted_path, "w", encoding="utf-8") as f:
            f.write(aborted_log)

        parser = tf_debug_parser.TFDebugParser(aborted_path)
        parser.parse()
        self.assertEqual(len(parser.failures), 1)
        self.assertEqual(parser.failures[0]["test_step"], 1)
        self.assertEqual(parser.failures[0]["phase"], "step_1")
        self.assertEqual(len(parser.events), 1)
        self.assertIsNone(parser.events[0]["test_step"])
        self.assertEqual(parser.events[0]["phase"], "destroy")


if __name__ == "__main__":
    unittest.main()
