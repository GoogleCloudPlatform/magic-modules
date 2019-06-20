#!/usr/bin/env python

import get_downstream_prs
from github import Github
import os
import re
import operator
import itertools
import sys


def get_unmerged_prs(g, dependencies):
  parsed_dependencies = [re.match(r'https://github.com/([\w-]+/[\w-]+)/pull/(\d+)', d).groups()
          for d in dependencies]
  parsed_dependencies.sort(key=operator.itemgetter(0))
  unmerged_dependencies = []
  # group those dependencies by repo - e.g. [("terraform-provider-google", ["123", "456"]), ...]
  for r, pulls in itertools.groupby(parsed_dependencies, key=operator.itemgetter(0)):
    repo = g.get_repo(r)
    for pull in pulls:
      # check whether the PR is merged - if it is, add it to the list.
      pr = repo.get_pull(int(pull[1]))
      if not pr.is_merged() and not pr.state == "closed":
        unmerged_dependencies.append(pull)
  return unmerged_dependencies


if __name__ == '__main__':
  g = Github(os.environ.get('GH_TOKEN'))
  assert len(sys.argv) == 2
  id_filename = sys.argv[1]
  unmerged = get_unmerged_prs(
          g, get_downstream_prs.get_github_dependencies(
              g, int(open(id_filename).read())))
  if unmerged:
    raise ValueError("some PRs are unmerged", unmerged)
