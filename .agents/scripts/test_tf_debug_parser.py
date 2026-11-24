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

    def test_sample_log_populates_new_summary_and_json_fields(self):
        parser = tf_debug_parser.TFDebugParser(self.log_path)
        parser.parse()
        payload = parser.to_dict()

        self.assertEqual(payload["endpoint_summary"], {"example.googleapis.com/v1": 4})
        self.assertEqual(
            payload["test_results"]["counts"], {"pass": 0, "fail": 1, "skip": 0}
        )
        self.assertEqual(len(payload["test_results"]["tests"]), 1)
        self.assertEqual(
            payload["test_results"]["tests"][0]["name"],
            "TestAccExampleResource_update",
        )
        self.assertEqual(payload["test_results"]["tests"][0]["status"], "FAIL")
        # Expected 404 during CheckDestroy is excluded from destroy_errors
        self.assertEqual(payload["destroy_errors"], [])

        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            parser.print_summary()
        output = stdout.getvalue()
        self.assertIn("====== ENDPOINT SUMMARY ======", output)
        self.assertIn("example.googleapis.com/v1: 4", output)
        self.assertIn("====== TEST RESULTS ======", output)
        self.assertIn("Summary: 0 passed, 1 failed, 0 skipped", output)
        self.assertIn("[FAIL] TestAccExampleResource_update (13.42s)", output)
        self.assertIn("====== DESTROY ERRORS ======", output)
        self.assertIn("No non-404 destroy errors found.", output)

    def test_endpoint_summary_hosts_versions_and_url_fallback(self):
        log_content = """=== RUN   TestAccMultiEndpoint
2026/10/28 08:00:01 [DEBUG] sdk.helper_resource: Starting TestStep: test_name=TestAccMultiEndpoint, test_step_number=1
---[ REQUEST ]---------------------------------------
POST /v1beta/projects/p/locations/us-central1/apps?alt=json&prettyPrint=false HTTP/1.1
Host: ces.googleapis.com

{"note": "Host: spoofed.googleapis.com"}
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET https://ces.googleapis.com/v1beta/projects/p/locations/us-central1/apps/a1?alt=json HTTP/1.1
Content-Type: application/json
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET /compute/v1/projects/p/zones/us-central1-a/instances/i1?alt=json HTTP/1.1
Host: compute.googleapis.com
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET /compute/beta/projects/p/zones/us-central1-a/instances/i1 HTTP/1.1
host: COMPUTE.googleapis.com
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET /storage/v1/b/my-bucket?projection=full HTTP/1.1
Host: storage.googleapis.com
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
POST /v1beta1/projects/p/locations/us-central1/endpoints?alt=json HTTP/1.1
Host: aiplatform.googleapis.com
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET /v2/projects/p/locations/global/keys?alt=json HTTP/1.1
Host: apikeys.googleapis.com
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET /v1/projects/p/resources?alt=json HTTP/1.1
Content-Type: application/json

{"Host": "ignored.googleapis.com"}
-----------------------------------------------------
"""
        log_path = os.path.join(self.temp_dir.name, "endpoints.log")
        with open(log_path, "w", encoding="utf-8") as f:
            f.write(log_content)

        parser = tf_debug_parser.TFDebugParser(log_path)
        parser.parse()
        expected = {
            "aiplatform.googleapis.com/v1beta1": 1,
            "apikeys.googleapis.com/v2": 1,
            "ces.googleapis.com/v1beta": 2,
            "compute.googleapis.com/compute/beta": 1,
            "compute.googleapis.com/compute/v1": 1,
            "storage.googleapis.com/storage/v1": 1,
            "unknown/v1": 1,
        }
        self.assertEqual(parser.endpoint_summary(), expected)
        self.assertEqual(list(parser.endpoint_summary().keys()), sorted(expected.keys()))

        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            parser.print_summary()
        output = stdout.getvalue()
        self.assertIn("ces.googleapis.com/v1beta: 2", output)
        self.assertIn("compute.googleapis.com/compute/v1: 1", output)
        self.assertIn("compute.googleapis.com/compute/beta: 1", output)
        self.assertIn("storage.googleapis.com/storage/v1: 1", output)

    def test_test_results_pass_fail_skip_and_subtests(self):
        log_content = """=== RUN   TestAccSuite
=== RUN   TestAccSuite/create
    --- PASS: TestAccSuite/create (4.12s)
=== RUN   TestAccSuite/beta_only
    test_suite_test.go:22: skipping in GA mode
    --- SKIP: TestAccSuite/beta_only (0.00s)
=== RUN   TestAccSuite/nested/broken
        --- FAIL: TestAccSuite/nested/broken (1.05s)
---[ RESPONSE ]--------------------------------------
HTTP/1.1 200 OK

{"note": "--- PASS: FakeInsidePayload (0.00s)"}
-----------------------------------------------------
2026/10/28 08:01:00 --- PASS: TestAccOther (2.50s)
--- FAIL: TestAccSuite (5.17s)
"""
        log_path = os.path.join(self.temp_dir.name, "results.log")
        with open(log_path, "w", encoding="utf-8") as f:
            f.write(log_content)

        parser = tf_debug_parser.TFDebugParser(log_path)
        parser.parse()
        payload = parser.to_dict()

        self.assertEqual(
            payload["test_results"]["counts"], {"pass": 2, "fail": 2, "skip": 1}
        )
        names_statuses = [
            (t["name"], t["status"], t["duration"])
            for t in payload["test_results"]["tests"]
        ]
        self.assertEqual(
            names_statuses,
            [
                ("TestAccSuite/create", "PASS", "4.12s"),
                ("TestAccSuite/beta_only", "SKIP", "0.00s"),
                ("TestAccSuite/nested/broken", "FAIL", "1.05s"),
                ("TestAccOther", "PASS", "2.50s"),
                ("TestAccSuite", "FAIL", "5.17s"),
            ],
        )
        # Backward compatibility: failures still records FAIL lines and excludes SKIP/PASS
        self.assertEqual(len(payload["failures"]), 2)
        self.assertTrue(all("FAIL:" in f["message"] for f in payload["failures"]))

        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            parser.print_summary()
        output = stdout.getvalue()
        self.assertIn("Summary: 2 passed, 2 failed, 1 skipped", output)
        self.assertIn("[SKIP] TestAccSuite/beta_only (0.00s)", output)
        self.assertNotIn("FakeInsidePayload", output)

    def test_destroy_errors_filters_404_and_non_destroy_errors(self):
        log_content = """=== RUN   TestAccDestroyErrors
2026/10/28 08:00:01 [DEBUG] sdk.helper_resource: Starting TestStep: test_name=TestAccDestroyErrors, test_step_number=1
---[ REQUEST ]---------------------------------------
POST /v1/projects/p/resources?alt=json HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
---[ RESPONSE ]--------------------------------------
HTTP/1.1 400 Bad Request
Content-Type: application/json

{"error": {"code": 400, "message": "Step 1 validation error"}}
-----------------------------------------------------
2026/10/28 08:00:02 [DEBUG] sdk.helper_resource: Finished TestStep: test_name=TestAccDestroyErrors, test_step_number=1
---[ REQUEST ]---------------------------------------
DELETE /v1/projects/p/resources/r1?alt=json&force=true HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
---[ RESPONSE ]--------------------------------------
HTTP/1.1 500 Internal Server Error
Content-Type: application/json

{"error": {"code": 500, "message": "Backend failure during delete"}}
-----------------------------------------------------
2026/10/28 08:00:03 [DEBUG] sdk.helper_resource: Calling TestCase CheckDestroy
---[ REQUEST ]---------------------------------------
GET /v1/projects/p/resources/r1?alt=json HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
---[ RESPONSE ]--------------------------------------
HTTP/1.1 404 Not Found
Content-Type: application/json

{"error": {"code": 404, "message": "Resource not found"}}
-----------------------------------------------------
---[ REQUEST ]---------------------------------------
GET /v1/projects/p/resources/r2?alt=json HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
---[ RESPONSE ]--------------------------------------
HTTP/1.1 403 Forbidden
Content-Type: application/json

{"error": {"code": 403, "message": "Permission denied on CheckDestroy"}}
-----------------------------------------------------
--- PASS: TestAccDestroyErrors (3.00s)
=== RUN   TestAccNextTest
---[ REQUEST ]---------------------------------------
GET /v1/projects/p/setup-check?alt=json HTTP/1.1
Host: example.googleapis.com
-----------------------------------------------------
---[ RESPONSE ]--------------------------------------
HTTP/1.1 503 Service Unavailable
-----------------------------------------------------
--- PASS: TestAccNextTest (1.00s)
"""
        log_path = os.path.join(self.temp_dir.name, "destroy_errors.log")
        with open(log_path, "w", encoding="utf-8") as f:
            f.write(log_content)

        parser = tf_debug_parser.TFDebugParser(log_path)
        parser.parse()
        errors = parser.to_dict()["destroy_errors"]

        # Only the 500 DELETE and 403 CheckDestroy GET should be reported;
        # Step 1 400, CheckDestroy 404, and TestAccNextTest setup 503 are excluded.
        self.assertEqual(len(errors), 2)
        self.assertEqual(errors[0]["status_code"], 500)
        self.assertEqual(errors[0]["method"], "DELETE")
        self.assertEqual(errors[0]["host"], "example.googleapis.com")
        self.assertEqual(errors[0]["url"], "/v1/projects/p/resources/r1?force=true")
        self.assertEqual(errors[0]["error_message"], "Backend failure during delete")
        self.assertEqual(errors[0]["phase"], "destroy")

        self.assertEqual(errors[1]["status_code"], 403)
        self.assertEqual(errors[1]["method"], "GET")
        self.assertEqual(errors[1]["url"], "/v1/projects/p/resources/r2")
        self.assertEqual(errors[1]["error_message"], "Permission denied on CheckDestroy")

        stdout = io.StringIO()
        with contextlib.redirect_stdout(stdout):
            parser.print_summary()
        output = stdout.getvalue()
        self.assertIn("500 Internal Server Error (DELETE example.googleapis.com/v1/projects/p/resources/r1?force=true) - Backend failure during delete", output)
        self.assertIn("403 Forbidden (GET example.googleapis.com/v1/projects/p/resources/r2) - Permission denied on CheckDestroy", output)
        self.assertNotIn("Step 1 validation error", output)
        self.assertNotIn("503 Service Unavailable", output.split("====== DESTROY ERRORS ======")[1])
        self.assertNotIn("alt=json", output)


if __name__ == "__main__":
    unittest.main()
