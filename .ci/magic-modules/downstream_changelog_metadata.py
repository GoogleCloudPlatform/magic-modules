#!/usr/bin/env python
"""
Script to edit downstream PRs with CHANGELOG release note and label metadata.

Usage:
  ./downstream_changelog_info.py path/to/.git/.id
  python /downstream_changelog_info.py

Note that release_note/labels are authoritative - if empty or not set in the MM
upstream PR, release notes will be removed from downstreams and labels
unset.
"""
from pyutils import strutils, downstreams
from github import Github
import os
import argparse

CHANGELOG_LABEL_PREFIX = "changelog: "

def downstream_changelog_info(gh, upstream_pr_num, changelog_repos):
  """Edit downstream PRs with CHANGELOG info.

  Args:
    gh: github.Github client
    upstream_pr_num: Upstream PR number
    changelog_repos: List of repo names to downstream changelog metadata for
  """
  # Parse CHANGELOG info from upstream
  upstream_pr = gh.get_repo(downstreams.UPSTREAM_REPO).get_pull(pr_num)
  release_note = changelog.get_release_note(upstream_pr.body)
  labels_to_add = changelog.find_prefixed_labels(
    [l.name for l in upstream_pr.labels],
    CHANGELOG_LABEL_PREFIX)

  print "Applying changelog info to downstreams for upstream PR %d:" % (
    upstream_pr.number)
  print "Release Note: \"%s\"" % release_note
  print "Labels: [%s]" % changelog_labels

  for repo_name, pulls in downstreams.get_parsed_downstream_urls(gh, pr_num):
    if repo_name not in changelog_repos:
      print "[DEBUG] skipping repo %s" % repo_name
      continue

    ghrepo = gh.get_repo(repo_name)
    for _r, pr_num in pulls:
      pr = ghrepo.get_pull(int(pr_num))
      set_changelog_info(pr, release_note, changelog_labels)

def set_changelog_info(gh_pull, release_note, labels_to_add):
  """Set release note and labels on a downstream PR in Github.

  Args:
    gh_pull: A github.PullRequest.PullRequest
    release_note: String of release note text to set
    changelog_labels: List of strings changelog labels to set
  """
  print "Setting changelog info for downstream PR %s" % downstream_pull.html_url
  edited_body = strutils.set_release_note(release_note, downstream_pull.body)
  downstream_pull.edit(body=edited_body)

  # Get all non-changelog-related labels
  original_labels = [l.name for l in downstream_pull.get_labels()]
  new_labels = [l for l in original if not l.startswith(CHANGELOG_LABEL_PREFIX)]
  new_labels += labels_to_add
  downstream_pull.set_labels(*new_labels)

if __name__ == '__main__':
  assert len(sys.argv) == 2, "expected id filename as argument"
  with open(sys.argv[1]) as f:
    pr_num = int(f.read())

    gh = Github(os.environ.get('GITHUB_TOKEN'))
    downstream_urls = os.environ.get('DOWNSTREAM_REPOS').split(',')
    downstream_changelog_info(gh, downstream_urls)
