#!/usr/bin/env python3
"""Unit tests for build_pr_link.py."""

from __future__ import annotations

import contextlib
import io
import os
import subprocess
import sys
import tempfile
import unittest
from unittest import mock

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
import build_pr_link  # noqa: E402


VALID_BODY_NONE = """Explain why this change exists.

```release-note:none
```"""


class TestTitleValidation(unittest.TestCase):
    def test_empty_and_whitespace_title(self):
        for title in ("", "   "):
            errors, warnings = build_pr_link.validate_title(title)
            self.assertEqual(len(errors), 1)
            self.assertIn("empty", errors[0].lower())
            self.assertEqual(warnings, [])

    def test_multiline_title_rejected(self):
        for title in ("agents: first line\nsecond line", "agents: line1\rline2"):
            errors, _ = build_pr_link.validate_title(title)
            self.assertTrue(any("single line" in e for e in errors))

    def test_title_length_boundaries(self):
        # 60 chars: target met, no error or warning
        t60 = "a" * 60
        errors, warnings = build_pr_link.validate_title(t60)
        self.assertEqual(errors, [])
        self.assertEqual(warnings, [])

        # 61 chars: soft warning, no error
        t61 = "a" * 61
        errors, warnings = build_pr_link.validate_title(t61)
        self.assertEqual(errors, [])
        self.assertEqual(len(warnings), 1)
        self.assertIn("61 chars", warnings[0])
        self.assertIn("<= 60", warnings[0])

        # 63 chars: hard cap boundary, soft warning, no error
        t63 = "a" * 63
        errors, warnings = build_pr_link.validate_title(t63)
        self.assertEqual(errors, [])
        self.assertEqual(len(warnings), 1)
        self.assertIn("63 chars", warnings[0])

        # 64 chars: hard error with exact count, overflow, target, and hints
        t64 = "a" * 64
        errors, warnings = build_pr_link.validate_title(t64)
        self.assertEqual(warnings, [])
        self.assertEqual(len(errors), 1)
        self.assertIn("64 chars", errors[0])
        self.assertIn("1 over", errors[0])
        self.assertIn("<= 60", errors[0])
        self.assertIn("google_<service>_", errors[0])
        self.assertIn("summarizing", errors[0])

        # 75 chars: verify overflow arithmetic (75 - 63 = 12)
        t75 = "a" * 75
        errors, _ = build_pr_link.validate_title(t75)
        self.assertEqual(len(errors), 1)
        self.assertIn("75 chars", errors[0])
        self.assertIn("12 over", errors[0])


class TestReleaseNotesAndBodyValidation(unittest.TestCase):
    def test_empty_body(self):
        for body in ("", "   \n\t "):
            errors, warnings = build_pr_link.validate_body(body)
            self.assertEqual(len(errors), 1)
            self.assertIn("empty", errors[0].lower())
            self.assertEqual(warnings, [])

    def test_missing_release_note_block(self):
        errors, _ = build_pr_link.validate_body("Just a description without notes.")
        self.assertEqual(len(errors), 1)
        self.assertIn("missing", errors[0].lower())

    def test_unclosed_release_note_block_at_eof(self):
        body = "Description.\n\n```release-note:none\n"
        errors, _ = build_pr_link.validate_body(body)
        self.assertEqual(len(errors), 1)
        self.assertIn("Unclosed", errors[0])

    def test_unclosed_release_note_block_before_second_block(self):
        body = (
            "Description.\n\n"
            "```release-note:bug\n"
            "compute: fixed foo\n"
            "```release-note:enhancement\n"
            "compute: added bar\n"
            "```"
        )
        errors, _ = build_pr_link.validate_body(body)
        self.assertEqual(len(errors), 1)
        self.assertIn("Unclosed", errors[0])

    def test_invalid_release_note_types(self):
        cases = [
            "```release-note\ncompute: fixed foo\n```",
            "```release-note:\ncompute: fixed foo\n```",
            "```release-note:REPLACEME\n\n```",
            "```release-note:feature\ncompute: added foo\n```",
        ]
        for block in cases:
            body = f"Desc.\n\n{block}"
            errors, _ = build_pr_link.validate_body(body)
            self.assertTrue(
                any("Invalid release-note type" in e for e in errors),
                f"Expected invalid type error for {block!r}, got {errors}",
            )

    def test_all_ten_valid_release_note_types(self):
        valid_blocks = {
            "enhancement": "compute: added `foo` field to `google_compute_instance`",
            "bug": "backup_dr: fixed drift on `google_backup_dr_backup_vault`",
            "note": "provider: upgraded Go toolchain to 1.24",
            "none": "",
            "new-resource": "`google_compute_example`",
            "new-list-resource": "`google_compute_example`",
            "new-datasource": "`google_compute_example`",
            "new-action": "`google_compute_instance_reset`",
            "deprecation": "compute: deprecated `old_field` on `google_compute_instance`",
            "breaking-change": "compute: removed `old_field` from `google_compute_instance`",
        }
        self.assertEqual(
            set(valid_blocks.keys()),
            set(build_pr_link.VALID_RELEASE_NOTE_TYPES),
        )
        for note_type, content in valid_blocks.items():
            if content:
                body = f"Desc.\n\n```release-note:{note_type}\n{content}\n```"
            else:
                body = f"Desc.\n\n```release-note:{note_type}\n```"
            errors, warnings = build_pr_link.validate_body(body)
            self.assertEqual(errors, [], f"Unexpected errors for {note_type}: {errors}")
            self.assertEqual(warnings, [])

    def test_release_note_none_variants(self):
        # Single blank line inside release-note:none is valid
        body_blank = "Desc.\n\n```release-note:none\n\n```"
        errors, _ = build_pr_link.validate_body(body_blank)
        self.assertEqual(errors, [])

        # Two blank lines inside release-note:none -> multi-line error
        body_two_blanks = "Desc.\n\n```release-note:none\n\n\n```"
        errors, _ = build_pr_link.validate_body(body_two_blanks)
        self.assertTrue(any("multiple lines" in e for e in errors))

        # Non-empty content inside release-note:none -> error
        body_non_empty = "Desc.\n\n```release-note:none\nno user impact\n```"
        errors, _ = build_pr_link.validate_body(body_non_empty)
        self.assertTrue(any("must be empty" in e for e in errors))

    def test_non_none_empty_content_rejected(self):
        for note_type in ("enhancement", "bug", "note", "deprecation"):
            for inner in ("", "\n", "   \n"):
                body = f"Desc.\n\n```release-note:{note_type}\n{inner}```"
                errors, _ = build_pr_link.validate_body(body)
                self.assertTrue(
                    any("must not be empty" in e for e in errors),
                    f"Expected empty error for {note_type} with {inner!r}, got {errors}",
                )

    def test_multiline_release_note_rejected(self):
        body = (
            "Desc.\n\n"
            "```release-note:enhancement\n"
            "compute: added `foo`\n"
            "compute: added `bar`\n"
            "```"
        )
        errors, _ = build_pr_link.validate_body(body)
        self.assertTrue(any("multiple lines" in e for e in errors))

        # Trailing blank line inside non-none block also fails go-changelog
        body_trailing = (
            "Desc.\n\n"
            "```release-note:enhancement\n"
            "compute: added `foo`\n"
            "\n"
            "```"
        )
        errors, _ = build_pr_link.validate_body(body_trailing)
        self.assertTrue(any("multiple lines" in e for e in errors))

    def test_multiple_valid_release_note_blocks(self):
        body = (
            "Desc.\n\n"
            "```release-note:bug\n"
            "compute: fixed `foo` drift\n"
            "```\n\n"
            "```release-note:enhancement\n"
            "compute: added `bar` field\n"
            "```"
        )
        errors, warnings = build_pr_link.validate_body(body)
        self.assertEqual(errors, [])
        self.assertEqual(warnings, [])

    def test_enhancement_and_bug_format(self):
        invalid_contents = [
            "added `foo` field",
            "Compute: added `foo` field",
            "compute engine: added `foo` field",
            "`compute`: added `foo` field",
            "compute:",
        ]
        for note_type in ("enhancement", "bug"):
            for content in invalid_contents:
                body = f"Desc.\n\n```release-note:{note_type}\n{content}\n```"
                errors, _ = build_pr_link.validate_body(body)
                self.assertTrue(
                    any("^[a-z0-9_]+: .+$" in e for e in errors),
                    f"Expected format error for {note_type} {content!r}, got {errors}",
                )

    def test_new_resource_and_datasource_format(self):
        valid_contents = [
            "`google_compute_instance`",
            "`google_backup_dr_vault` (beta)",
            "`google_compute_instance` (ga)",
        ]
        invalid_contents = [
            "google_compute_instance",
            "`compute_instance`",
            "`google compute instance`",
            "compute: added `google_compute_instance`",
            "`google_compute_instance` (alpha)",
        ]
        for note_type in ("new-resource", "new-datasource"):
            for content in valid_contents:
                body = f"Desc.\n\n```release-note:{note_type}\n{content}\n```"
                errors, _ = build_pr_link.validate_body(body)
                self.assertEqual(errors, [], f"Unexpected error for {content!r}: {errors}")
            for content in invalid_contents:
                body = f"Desc.\n\n```release-note:{note_type}\n{content}\n```"
                errors, _ = build_pr_link.validate_body(body)
                self.assertTrue(
                    any("`google_<name>`" in e for e in errors),
                    f"Expected format error for {note_type} {content!r}, got {errors}",
                )

    def test_crlf_line_endings_supported(self):
        body = "Fix bug.\r\n\r\n```release-note:bug\r\ncompute: fixed foo\r\n```"
        errors, warnings = build_pr_link.validate_body(body)
        self.assertEqual(errors, [])
        self.assertEqual(warnings, [])

    def test_soft_warnings_boilerplate_headings_and_html_comments(self):
        for heading in (
            "## Summary",
            "### Summary of Changes",
            "## Motivation",
            "## Motivation & Proof",
            "### Verification",
        ):
            body = f"{heading}\n\nDetails.\n\n```release-note:none\n```"
            errors, warnings = build_pr_link.validate_body(body)
            self.assertEqual(errors, [])
            self.assertTrue(
                any("AI-boilerplate heading" in w for w in warnings),
                f"Expected warning for {heading!r}, got {warnings}",
            )

        # Preferred reviewer-oriented headings must not warn
        good_body = (
            "Short problem statement.\n\n"
            "### Notable implementation details\n\n"
            "- Kept validation pure.\n\n"
            "```release-note:none\n```"
        )
        errors, warnings = build_pr_link.validate_body(good_body)
        self.assertEqual(errors, [])
        self.assertEqual(warnings, [])

        # Leftover HTML comment triggers soft warning
        html_body = "<!-- template comment -->\nDesc.\n\n```release-note:none\n```"
        errors, warnings = build_pr_link.validate_body(html_body)
        self.assertEqual(errors, [])
        self.assertTrue(any("HTML comment" in w for w in warnings))

        # Boilerplate heading or HTML comment inside fenced code block does not warn
        fenced_body = (
            "Example:\n\n"
            "```markdown\n"
            "## Summary\n"
            "<!-- comment -->\n"
            "```\n\n"
            "```release-note:none\n```"
        )
        errors, warnings = build_pr_link.validate_body(fenced_body)
        self.assertEqual(errors, [])
        self.assertEqual(warnings, [])


class TestGitRemoteParsing(unittest.TestCase):
    def test_extract_github_user_formats(self):
        cases = {
            "https://github.com/c2thorn/magic-modules.git": "c2thorn",
            "https://github.com/c2thorn/magic-modules": "c2thorn",
            "https://token@github.com/alice-dev/magic-modules.git": "alice-dev",
            "git@github.com:c2thorn/magic-modules.git": "c2thorn",
            "git@github.com:bob/magic-modules": "bob",
            "ssh://git@github.com/carol/magic-modules.git": "carol",
            "git://github.com/dave/magic-modules.git": "dave",
            "https://gitlab.com/not-github/repo.git": None,
            "": None,
        }
        for url, expected in cases.items():
            self.assertEqual(build_pr_link.extract_github_user(url), expected, url)

    def test_detect_fork_user_from_remotes(self):
        standard = (
            "origin\thttps://github.com/c2thorn/magic-modules.git (fetch)\n"
            "origin\thttps://github.com/c2thorn/magic-modules.git (push)\n"
            "upstream\thttps://github.com/GoogleCloudPlatform/magic-modules.git (fetch)\n"
            "upstream\thttps://github.com/GoogleCloudPlatform/magic-modules.git (push)\n"
        )
        self.assertEqual(
            build_pr_link.detect_fork_user_from_remotes(standard), "c2thorn"
        )

        # origin is upstream, fork is SSH under a custom remote name
        swapped = (
            "origin\tgit@github.com:googlecloudplatform/magic-modules.git (fetch)\n"
            "fork\tgit@github.com:rileykarson/magic-modules.git (fetch)\n"
        )
        self.assertEqual(
            build_pr_link.detect_fork_user_from_remotes(swapped), "rileykarson"
        )

        # Prefers origin when multiple non-upstream remotes exist
        multi = (
            "backup\thttps://github.com/other/magic-modules.git (fetch)\n"
            "origin\tgit@github.com:primary/magic-modules.git (fetch)\n"
        )
        self.assertEqual(
            build_pr_link.detect_fork_user_from_remotes(multi), "primary"
        )

        # Only upstream configured -> None
        upstream_only = (
            "upstream\thttps://github.com/GoogleCloudPlatform/magic-modules.git (fetch)\n"
        )
        self.assertIsNone(build_pr_link.detect_fork_user_from_remotes(upstream_only))


class TestUrlBuildingAndEndToEndValidation(unittest.TestCase):
    def test_url_encoding_and_markdown_safety(self):
        url, errors, warnings = build_pr_link.validate_and_build(
            title="agents: validate PR links (63-char cap)",
            body="Fix `a & b` (see #1).\n\n```release-note:none\n```",
            branch="feat/pr-link",
            base="main",
            user="c2thorn",
        )
        self.assertEqual(errors, [])
        self.assertEqual(warnings, [])
        assert url is not None
        self.assertTrue(
            url.startswith(
                "https://github.com/GoogleCloudPlatform/magic-modules/compare/"
                "main...c2thorn:feat/pr-link?expand=1&"
            )
        )
        # Parentheses must be percent-encoded so Markdown link syntax never breaks
        self.assertNotIn("(", url)
        self.assertNotIn(")", url)
        self.assertIn("%2863-char%20cap%29", url)
        self.assertIn("%26", url)
        self.assertIn("%231", url)

    def test_branch_and_user_errors(self):
        # branch == base
        url, errors, _ = build_pr_link.validate_and_build(
            title="agents: valid title",
            body=VALID_BODY_NONE,
            branch="main",
            base="main",
            user="c2thorn",
        )
        self.assertIsNone(url)
        self.assertTrue(any("same as base branch" in e for e in errors))

        # detached HEAD
        url, errors, _ = build_pr_link.validate_and_build(
            title="agents: valid title",
            body=VALID_BODY_NONE,
            branch="HEAD",
            base="main",
            user="c2thorn",
        )
        self.assertIsNone(url)
        self.assertTrue(any("Detached HEAD" in e for e in errors))

        # empty branch and user
        url, errors, _ = build_pr_link.validate_and_build(
            title="agents: valid title",
            body=VALID_BODY_NONE,
            branch="",
            base="main",
            user="",
        )
        self.assertIsNone(url)
        self.assertEqual(len(errors), 2)

    def test_collects_multiple_errors_together(self):
        url, errors, _ = build_pr_link.validate_and_build(
            title="x" * 70,
            body="No release note here.",
            branch="main",
            base="main",
            user="",
        )
        self.assertIsNone(url)
        self.assertGreaterEqual(len(errors), 4)

    def test_url_length_limit_enforced(self):
        huge_body = ("a" * 8000) + "\n\n```release-note:none\n```"
        url, errors, _ = build_pr_link.validate_and_build(
            title="agents: test huge body",
            body=huge_body,
            branch="my-branch",
            base="main",
            user="c2thorn",
        )
        self.assertIsNone(url)
        self.assertTrue(any("8000-byte" in e for e in errors))


class TestCLI(unittest.TestCase):
    def _run_main(self, argv: list[str], stdin_text: str = "") -> tuple[int, str, str]:
        stdout = io.StringIO()
        stderr = io.StringIO()
        with (
            contextlib.redirect_stdout(stdout),
            contextlib.redirect_stderr(stderr),
            mock.patch.object(sys, "stdin", io.StringIO(stdin_text)),
        ):
            code = build_pr_link.main(argv)
        return code, stdout.getvalue(), stderr.getvalue()

    def test_cli_body_and_markdown_flag(self):
        code, out, err = self._run_main(
            [
                "--title",
                "agents: enforce PR title limit",
                "--body",
                VALID_BODY_NONE,
                "--branch",
                "fix-pr-link",
                "--user",
                "c2thorn",
                "--markdown",
            ]
        )
        self.assertEqual(code, 0)
        self.assertEqual(err, "")
        self.assertTrue(
            out.strip().startswith("👉 **[Create Pull Request on GitHub](https://github.com/")
        )
        self.assertTrue(out.strip().endswith(")**"))

    def test_cli_body_file_and_stdin(self):
        with tempfile.NamedTemporaryFile("w", encoding="utf-8", delete=False) as tf:
            tf.write(VALID_BODY_NONE)
            body_path = tf.name
        try:
            code, out, err = self._run_main(
                [
                    "--title",
                    "agents: test body file",
                    "--body-file",
                    body_path,
                    "--branch",
                    "feat-file",
                    "--user",
                    "c2thorn",
                ]
            )
            self.assertEqual(code, 0)
            self.assertEqual(err, "")
            self.assertIn("compare/main...c2thorn:feat-file", out)
        finally:
            os.unlink(body_path)

        # Read from stdin with --body-file -
        code, out, err = self._run_main(
            [
                "--title",
                "agents: test stdin body",
                "--body-file",
                "-",
                "--branch",
                "feat-stdin",
                "--user",
                "c2thorn",
            ],
            stdin_text=VALID_BODY_NONE,
        )
        self.assertEqual(code, 0)
        self.assertEqual(err, "")
        self.assertIn("compare/main...c2thorn:feat-stdin", out)

    def test_cli_missing_body_file_errors_cleanly(self):
        code, out, err = self._run_main(
            [
                "--title",
                "agents: missing file",
                "--body-file",
                "/nonexistent/path/pr_body.md",
                "--branch",
                "feat-x",
                "--user",
                "c2thorn",
            ]
        )
        self.assertEqual(code, 1)
        self.assertEqual(out, "")
        self.assertIn("failed to read --body-file", err)

    def test_cli_warnings_still_exit_zero(self):
        title_62 = "a" * 62
        body = "## Summary\n\n<!-- leftover -->\n\n```release-note:none\n```"
        code, out, err = self._run_main(
            [
                "--title",
                title_62,
                "--body",
                body,
                "--branch",
                "feat-warn",
                "--user",
                "c2thorn",
            ]
        )
        self.assertEqual(code, 0)
        self.assertIn("https://github.com/", out)
        self.assertIn("Warning:", err)
        self.assertIn("62 chars", err)
        self.assertIn("AI-boilerplate heading", err)
        self.assertIn("HTML comment", err)

    def test_cli_hard_error_exits_nonzero(self):
        code, out, err = self._run_main(
            [
                "--title",
                "a" * 64,
                "--body",
                "No release note",
                "--branch",
                "feat-err",
                "--user",
                "c2thorn",
            ]
        )
        self.assertEqual(code, 1)
        self.assertEqual(out, "")
        self.assertIn("Error:", err)
        self.assertIn("64 chars", err)
        self.assertIn("missing", err.lower())

    def test_cli_auto_detects_branch_and_user(self):
        with (
            mock.patch.object(build_pr_link, "detect_git_branch", return_value="my-feat"),
            mock.patch.object(build_pr_link, "detect_fork_user", return_value="c2thorn"),
        ):
            code, out, err = self._run_main(
                [
                    "--title",
                    "agents: auto detect",
                    "--body",
                    VALID_BODY_NONE,
                ]
            )
        self.assertEqual(code, 0)
        self.assertEqual(err, "")
        self.assertIn("compare/main...c2thorn:my-feat", out)

    def test_executable_script_invocation(self):
        script_path = os.path.join(
            os.path.dirname(os.path.abspath(__file__)), "build_pr_link.py"
        )
        result = subprocess.run(
            [
                sys.executable,
                script_path,
                "--title",
                "agents: subprocess smoke test",
                "--body-file",
                "-",
                "--branch",
                "smoke-branch",
                "--base",
                "FEATURE-BRANCH-major-release-8.0.0",
                "--user",
                "c2thorn",
            ],
            input=VALID_BODY_NONE,
            capture_output=True,
            text=True,
            check=False,
        )
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn(
            "compare/FEATURE-BRANCH-major-release-8.0.0...c2thorn:smoke-branch",
            result.stdout,
        )


if __name__ == "__main__":
    unittest.main()
