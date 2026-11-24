#!/usr/bin/env python3
"""Validate PR metadata and build a pre-filled GitHub compare URL.

Enforces Magic Modules PR title length limits, release-note block syntax
matching tools/go-changelog/note.go, and GitHub URL length ceilings so agents
cannot emit invalid PR links.
"""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
import urllib.parse

DEFAULT_REPO = "GoogleCloudPlatform/magic-modules"
DEFAULT_BASE = "main"
TITLE_TARGET_LEN = 60
TITLE_MAX_LEN = 63
MAX_URL_BYTES = 8000

VALID_RELEASE_NOTE_TYPES = (
    "enhancement",
    "bug",
    "note",
    "none",
    "new-resource",
    "new-list-resource",
    "new-datasource",
    "new-action",
    "deprecation",
    "breaking-change",
)
VALID_RELEASE_NOTE_TYPE_SET = frozenset(VALID_RELEASE_NOTE_TYPES)

_RELEASE_NOTE_OPEN_RE = re.compile(r"^```release-note(?::(.*))?$")
_ENHANCEMENT_OR_BUG_RE = re.compile(r"^[a-z0-9_]+: .+$")
_NEW_RESOURCE_OR_DATASOURCE_RE = re.compile(
    r"^`google_[a-z0-9_]+`( \((beta|ga)\))?$"
)
_GITHUB_USER_RE = re.compile(
    r"github\.com[:/]([^/\s]+)/[^/\s]+?(?:\.git)?$", re.IGNORECASE
)
_FENCE_OPEN_RE = re.compile(r"^(`{3,}|~{3,})")
_BOILERPLATE_HEADING_RE = re.compile(
    r"^#{1,6}\s+(?:summary(?:\s+of\s+changes)?|motivation(?:\s*&\s*proof)?|verification)\s*$",
    re.IGNORECASE,
)


def extract_github_user(remote_url: str) -> str | None:
    """Extract the GitHub owner/username from an SSH or HTTPS remote URL."""
    url = remote_url.strip()
    if not url:
        return None
    match = _GITHUB_USER_RE.search(url)
    if not match:
        return None
    return match.group(1)


def detect_fork_user_from_remotes(
    remotes_output: str, upstream_repo: str = DEFAULT_REPO
) -> str | None:
    """Parse `git remote -v` output to find a non-upstream GitHub fork user."""
    upstream_lower = upstream_repo.lower()
    candidates: list[tuple[str, str]] = []
    seen_remotes: set[str] = set()

    for raw_line in remotes_output.splitlines():
        parts = raw_line.strip().split()
        if len(parts) < 2:
            continue
        remote_name, remote_url = parts[0], parts[1]
        if remote_name in seen_remotes:
            continue
        if upstream_lower in remote_url.lower():
            continue
        user = extract_github_user(remote_url)
        if user:
            seen_remotes.add(remote_name)
            candidates.append((remote_name, user))

    for remote_name, user in candidates:
        if remote_name == "origin":
            return user
    return candidates[0][1] if candidates else None


def detect_git_branch(cwd: str | None = None) -> str:
    """Return the current git branch name, or empty string if unavailable."""
    try:
        result = subprocess.run(
            ["git", "rev-parse", "--abbrev-ref", "HEAD"],
            cwd=cwd,
            capture_output=True,
            text=True,
            check=True,
        )
        return result.stdout.strip()
    except (OSError, subprocess.CalledProcessError):
        return ""


def detect_fork_user(
    upstream_repo: str = DEFAULT_REPO, cwd: str | None = None
) -> str:
    """Detect the personal fork GitHub username from git remotes."""
    try:
        result = subprocess.run(
            ["git", "remote", "-v"],
            cwd=cwd,
            capture_output=True,
            text=True,
            check=True,
        )
    except (OSError, subprocess.CalledProcessError):
        return ""
    return detect_fork_user_from_remotes(result.stdout, upstream_repo) or ""


def parse_release_note_blocks(
    body: str,
) -> tuple[list[tuple[str | None, str]], list[str]]:
    """Extract fenced ```release-note:<type> blocks and any parse errors."""
    normalized = body.replace("\r\n", "\n").replace("\r", "\n")
    lines = normalized.split("\n")

    blocks: list[tuple[str | None, str]] = []
    errors: list[str] = []

    in_block = False
    current_type: str | None = None
    current_lines: list[str] = []

    for line in lines:
        if in_block:
            if line.rstrip() == "```":
                blocks.append((current_type, "\n".join(current_lines)))
                in_block = False
                current_type = None
                current_lines = []
                continue
            nested_open = _RELEASE_NOTE_OPEN_RE.match(line.rstrip())
            if nested_open:
                label = (
                    f"release-note:{current_type}"
                    if current_type is not None
                    else "release-note"
                )
                errors.append(
                    f"Unclosed fenced ```{label} block (missing closing ```)."
                )
                raw_type = nested_open.group(1)
                current_type = raw_type.strip() if raw_type is not None else None
                current_lines = []
                continue
            current_lines.append(line)
        else:
            open_match = _RELEASE_NOTE_OPEN_RE.match(line.rstrip())
            if open_match:
                in_block = True
                raw_type = open_match.group(1)
                current_type = raw_type.strip() if raw_type is not None else None
                current_lines = []

    if in_block:
        label = (
            f"release-note:{current_type}"
            if current_type is not None
            else "release-note"
        )
        errors.append(
            f"Unclosed fenced ```{label} block (missing closing ```)."
        )

    return blocks, errors


def validate_title(title: str) -> tuple[list[str], list[str]]:
    """Validate PR title constraints and return (errors, warnings)."""
    errors: list[str] = []
    warnings: list[str] = []

    if not title or not title.strip():
        errors.append("PR title must not be empty.")
        return errors, warnings

    if "\n" in title or "\r" in title:
        errors.append("PR title must be a single line (no newline characters).")

    length = len(title)
    if length > TITLE_MAX_LEN:
        overflow = length - TITLE_MAX_LEN
        errors.append(
            f"PR title is {length} chars ({overflow} over the {TITLE_MAX_LEN}-char "
            f"hard limit; target <= {TITLE_TARGET_LEN}). Shorten it by dropping "
            "redundant 'google_<service>_' resource prefixes after '<service>:' "
            "and summarizing the theme instead of listing every field or check."
        )
    elif length > TITLE_TARGET_LEN:
        warnings.append(
            f"PR title is {length} chars (above the <= {TITLE_TARGET_LEN}-char "
            f"target, within the {TITLE_MAX_LEN}-char hard limit). Consider "
            "shortening it."
        )

    return errors, warnings


def validate_body(body: str) -> tuple[list[str], list[str]]:
    """Validate PR body release notes and style checks; return (errors, warnings)."""
    errors: list[str] = []
    warnings: list[str] = []

    if not body or not body.strip():
        errors.append("PR body must not be empty.")
        return errors, warnings

    blocks, parse_errors = parse_release_note_blocks(body)
    errors.extend(parse_errors)

    if not blocks and not parse_errors:
        errors.append(
            "PR body is missing a fenced ```release-note:<type> block. "
            f"Valid types: {', '.join(VALID_RELEASE_NOTE_TYPES)}."
        )

    for note_type, content in blocks:
        if note_type is None or note_type not in VALID_RELEASE_NOTE_TYPE_SET:
            shown = "<missing>" if not note_type else note_type
            errors.append(
                f"Invalid release-note type '{shown}'. "
                f"Valid types: {', '.join(VALID_RELEASE_NOTE_TYPES)}."
            )
            continue

        if "\n" in content:
            errors.append(
                f"Release-note block 'release-note:{note_type}' contains multiple "
                "lines. Use a single line per block (or multiple blocks for "
                "separate changes)."
            )
            continue

        if note_type == "none":
            if content.strip():
                errors.append(
                    "Release-note block 'release-note:none' must be empty, "
                    f"got: {content!r}."
                )
            continue

        if not content.strip():
            errors.append(
                f"Release-note block 'release-note:{note_type}' must not be empty."
            )
            continue

        if note_type in ("enhancement", "bug"):
            if not _ENHANCEMENT_OR_BUG_RE.match(content):
                errors.append(
                    f"Release-note block 'release-note:{note_type}' must match "
                    "'<product>: <description>' (regex ^[a-z0-9_]+: .+$), "
                    f"got: {content!r}."
                )
        elif note_type in ("new-resource", "new-datasource"):
            if not _NEW_RESOURCE_OR_DATASOURCE_RE.match(content):
                errors.append(
                    f"Release-note block 'release-note:{note_type}' must match "
                    "'`google_<name>`' (with optional ' (beta)' or ' (ga)' suffix), "
                    f"got: {content!r}."
                )

    # Soft warnings on non-fenced markdown body lines.
    normalized = body.replace("\r\n", "\n").replace("\r", "\n")
    active_fence_char = ""
    active_fence_len = 0
    found_headings: list[str] = []
    found_html_comment = False

    for line in normalized.split("\n"):
        stripped = line.strip()
        if active_fence_char:
            if (
                len(stripped) >= active_fence_len
                and set(stripped) == {active_fence_char}
            ):
                active_fence_char = ""
                active_fence_len = 0
            continue

        fence_match = _FENCE_OPEN_RE.match(line)
        if fence_match:
            token = fence_match.group(1)
            active_fence_char = token[0]
            active_fence_len = len(token)
            continue

        if _BOILERPLATE_HEADING_RE.match(stripped):
            found_headings.append(stripped)
        if "<!--" in line:
            found_html_comment = True

    if found_headings:
        headings_list = ", ".join(repr(h) for h in found_headings)
        warnings.append(
            f"PR body contains AI-boilerplate heading(s) ({headings_list}). "
            "Prefer plain reviewer-focused prose without boilerplate headings."
        )
    if found_html_comment:
        warnings.append(
            "PR body contains an HTML comment ('<!--'). Remove leftover PR "
            "template comments before generating the link."
        )

    return errors, warnings


def build_compare_url(
    repo: str,
    base: str,
    user: str,
    branch: str,
    title: str,
    body: str,
) -> str:
    """Build a percent-encoded GitHub compare URL for opening a PR."""
    encoded_base = urllib.parse.quote(base, safe="/")
    encoded_user = urllib.parse.quote(user, safe="")
    encoded_branch = urllib.parse.quote(branch, safe="/")
    query = urllib.parse.urlencode(
        {"expand": "1", "title": title, "body": body},
        quote_via=urllib.parse.quote,
    )
    return (
        f"https://github.com/{repo}/compare/"
        f"{encoded_base}...{encoded_user}:{encoded_branch}?{query}"
    )


def validate_and_build(
    title: str,
    body: str,
    branch: str,
    base: str = DEFAULT_BASE,
    user: str = "",
    repo: str = DEFAULT_REPO,
) -> tuple[str | None, list[str], list[str]]:
    """Validate all inputs and return (url_or_none, errors, warnings)."""
    errors: list[str] = []
    warnings: list[str] = []

    title_errors, title_warnings = validate_title(title)
    errors.extend(title_errors)
    warnings.extend(title_warnings)

    body_errors, body_warnings = validate_body(body)
    errors.extend(body_errors)
    warnings.extend(body_warnings)

    if not repo or not repo.strip():
        errors.append("Target repository (--repo) must not be empty.")
    if not base or not base.strip():
        errors.append("Base branch (--base) must not be empty.")
    if not branch or not branch.strip():
        errors.append(
            "Head branch could not be determined; pass --branch explicitly."
        )
    elif branch.strip() == "HEAD":
        errors.append(
            "Detached HEAD detected; check out a feature branch or pass --branch."
        )
    elif base and branch.strip() == base.strip():
        errors.append(
            f"Head branch '{branch.strip()}' cannot be the same as base branch "
            f"'{base.strip()}'."
        )

    if not user or not user.strip():
        errors.append(
            "Fork GitHub username could not be determined from git remotes; "
            "pass --user explicitly."
        )

    url: str | None = None
    if (
        repo
        and repo.strip()
        and base
        and base.strip()
        and user
        and user.strip()
        and branch
        and branch.strip()
    ):
        url = build_compare_url(
            repo=repo.strip(),
            base=base.strip(),
            user=user.strip(),
            branch=branch.strip(),
            title=title,
            body=body,
        )
        url_bytes = len(url.encode("utf-8"))
        if url_bytes > MAX_URL_BYTES:
            errors.append(
                f"Generated PR URL is {url_bytes} bytes (exceeds {MAX_URL_BYTES}-byte "
                "limit and may trigger GitHub HTTP 414 URI Too Long). Trim the PR body."
            )

    if errors:
        return None, errors, warnings
    return url, [], warnings


def format_markdown_link(url: str) -> str:
    """Format the compare URL as a clickable Markdown callout link."""
    return f"👉 **[Create Pull Request on GitHub]({url})**"


def build_arg_parser() -> argparse.ArgumentParser:
    """Build CLI argument parser."""
    parser = argparse.ArgumentParser(
        description=(
            "Validate PR title/body and generate a pre-filled GitHub compare URL."
        )
    )
    parser.add_argument("--title", required=True, help="PR title (<= 63 chars; target <= 60).")
    body_group = parser.add_mutually_exclusive_group(required=True)
    body_group.add_argument("--body", help="PR body text.")
    body_group.add_argument(
        "--body-file",
        help="Path to file containing PR body, or '-' to read from stdin.",
    )
    parser.add_argument(
        "--branch",
        help="Head branch name (defaults to current git branch).",
    )
    parser.add_argument(
        "--base",
        default=DEFAULT_BASE,
        help=f"Base branch name (default: {DEFAULT_BASE}).",
    )
    parser.add_argument(
        "--user",
        help="Fork GitHub username (defaults to non-upstream git remote owner).",
    )
    parser.add_argument(
        "--repo",
        default=DEFAULT_REPO,
        help=f"Upstream GitHub repository (default: {DEFAULT_REPO}).",
    )
    parser.add_argument(
        "--markdown",
        action="store_true",
        help="Output a Markdown link instead of the raw URL.",
    )
    return parser


def main(argv: list[str] | None = None) -> int:
    parser = build_arg_parser()
    args = parser.parse_args(argv)

    if args.body is not None:
        body = args.body
    elif args.body_file == "-":
        body = sys.stdin.read()
    else:
        try:
            with open(args.body_file, encoding="utf-8") as fh:
                body = fh.read()
        except OSError as exc:
            print(f"Error: failed to read --body-file: {exc}", file=sys.stderr)
            return 1

    branch = args.branch if args.branch is not None else detect_git_branch()
    user = args.user if args.user is not None else detect_fork_user(args.repo)

    url, errors, warnings = validate_and_build(
        title=args.title,
        body=body,
        branch=branch,
        base=args.base,
        user=user,
        repo=args.repo,
    )

    for warning in warnings:
        print(f"Warning: {warning}", file=sys.stderr)

    if errors:
        for error in errors:
            print(f"Error: {error}", file=sys.stderr)
        return 1

    assert url is not None
    print(format_markdown_link(url) if args.markdown else url)
    return 0


if __name__ == "__main__":
    sys.exit(main())
