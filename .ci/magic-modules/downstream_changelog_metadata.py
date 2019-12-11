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
import os
import sys
import github
from pyutils import strutils, downstreams

CHANGELOG_LABEL_PREFIX = "changelog: "

def downstream_changelog_info(gh, upstream_pr_num, changelog_repos):
  """Edit downstream PRs with CHANGELOG info.

  Args:
    gh: github.Github client
    upstream_pr_num: Upstream PR number
    changelog_repos: List of repo names to downstream changelog metadata for
  """
  # Parse CHANGELOG info from upstream
  print "Fetching upstream PR '%s'..." % upstream_pr_num
  upstream_pr = gh.get_repo(downstreams.UPSTREAM_REPO)\
                  .get_pull(upstream_pr_num)
  release_notes = strutils.get_release_notes(upstream_pr.body)
  labels_to_add = strutils.find_prefixed_labels(
    [l.name for l in upstream_pr.labels],
    CHANGELOG_LABEL_PREFIX)

  if not labels_to_add and not release_notes:
    print "No release note or labels found, skipping PR %d" % (
      upstream_pr_num)
    return

  print "Found changelog info on upstream PR %d:" % (
    upstream_pr.number)
  print "Release Note: \"%s\"" % release_notes
  print "Labels: %s" % labels_to_add

  parsed_urls = downstreams.get_parsed_downstream_urls(gh, upstream_pr.number)
  found = False

  for repo_name, pulls in parsed_urls:
    found = True
    print "Found downstream PR for repo %s" % repo_name

    if repo_name not in changelog_repos:
      print "[DEBUG] skipping repo %s with no CHANGELOG" % repo_name
      continue

    print "Generating changelog for pull requests in %s" % repo_name

    print "Fetching repo %s" % repo_name
    ghrepo = gh.get_repo(repo_name)

    for _r, prnum in pulls:
      print "Fetching %s PR %d" % (repo_name, prnum)
      pr = ghrepo.get_pull(int(prnum))
      set_changelog_info(pr, release_notes, labels_to_add)

  if not found:
    print "No downstreams found for upstream PR %d, returning!" % upstream_pr.number

def set_changelog_info(gh_pull, release_notes, labels_to_add):
  """Set release note and labels on a downstream PR in Github.

  Args:
    gh_pull: A github.PullRequest.PullRequest handle
    release_note: String of release note text to set
    labels_to_add: List of strings. Changelog-related labels to add/replace.
  """
  print "Setting changelog info for downstream PR %s" % gh_pull.url
  edited_body = strutils.set_release_notes(release_notes, gh_pull.body)
  gh_pull.edit(body=edited_body)

  # Get all non-changelog-related labels
  labels_to_set = []
  for l in gh_pull.get_labels():
    if not l.name.startswith(CHANGELOG_LABEL_PREFIX):
      labels_to_set.append(l.name)
  labels_to_set += labels_to_add
  gh_pull.set_labels(*labels_to_set)


if __name__ == '__main__':
  downstream_repos = os.environ.get('DOWNSTREAM_REPOS').split(',')
  if len(downstream_repos) == 0:
    print "Skipping, no downstreams repos given to downstream changelog info for"
    sys.exit(0)

  assert len(sys.argv) == 2, "expected id filename as argument"
  with open(sys.argv[1]) as f:
    pr_num = int(f.read())
    downstream_changelog_info(
      github.Github(os.environ.get('GITHUB_TOKEN')),
      pr_num, downstream_repos)
