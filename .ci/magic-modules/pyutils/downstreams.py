"""Helper class for obtaining information about upstream PR and its downstreams.

  Typical usage example:

  import upstream_pull_request

  client = github.Github(github_token)
  downstreams = upstream_pull_request.downstream_urls(client, 100)

"""

import os
import re
import sys
import itertools
import operator
from strutils import *

UPSTREAM_REPO = 'GoogleCloudPlatform/magic-modules'

def find_unmerged_downstreams(client, pr_num):
  """Returns list of urls for unmerged, open downstreams.

  For each downstream PR URL found from get_parsed_downstream_urls(),
  fetches the status of each downstream PR to determine which PRs are still
  unmerged (i.e. not closed and not merged).

  Args:
    client: github.Github client
    pr_num: PR Number for upstream PR
  Returns:
    All unmerged downstreams found for a PR.
  """
  unmerged_dependencies = []
  for r, pulls in get_parsed_downstream_urls(client, pr_num):
    repo = client.get_repo(r)
    for _repo, pr_num in pulls:
      pr = repo.get_pull(int(pr_num))
      # Disregard merged or closed PRs.
      if not pr.is_merged() and not pr.state == "closed":
        unmerged_dependencies.append(pr.html_url)

  return unmerged_dependencies

def get_parsed_downstream_urls(client, pr_num):
  """Get parsed URLs for downstream PRs grouped by repo.

  For each downstream PR URL referenced by the upstream PR, this method
  parses the downstream repo name
  (i.e. "terraform-providers/terraform-providers-google") and PR number
  (e.g. 100) and groups them by repo name so calling code only needs to fetch
  each repo once.

  Example:
    parsed = UpstreamPullRequest(pr_num).parsed_downstream_urls
    for repo, repo_pulls in parsed:
      for _repo, pr in repo_pulls:
        print "Downstream is https://github.com/%s/pull/%d" % (repo, pr)

  Args:
    client: github.Github client
    pr_num: PR Number for upstream PR

  Returns:
    Iterator over $repo and sub-iterators of ($repo, $pr_num) parsed tuples
  """
  parsed = [parse_github_url(u) for u in get_downstream_urls(client, pr_num)]
  return itertools.groupby(parsed, key=operator.itemgetter(0))

def get_downstream_urls(client, pr_num):
  """Get list of URLs for downstream PRs.

  This fetches the upstream PR and finds its downstream PR URLs by
  searching for references in its comments.

  Args:
    client: github.Github client
    pr_num: PR Number for upstream PR

  Returns:
    List of downstream PR URLs.
  """
  urls = []
  print "Getting downstream URLs for PR %d..." % pr_num
  pr = client.get_repo(UPSTREAM_REPO).get_pull(pr_num)
  for comment in pr.get_issue_comments():
    urls = urls + find_dependency_urls_in_comment(comment.body)
  print "Found downstream URLs: %s" % urls
  return urls
