#!/usr/bin/env python
"""
This script takes the name of a file containing an upstream PR number
and returns an error if not all of its downstreams have been merged.

Required env vars:
  GH_TOKEN: Github token
"""
import os
import sys
from github import Github
from pyutils import downstreams

if __name__ == '__main__':
  assert len(sys.argv) == 2, "expected id filename as argument"
  with open(sys.argv[1]) as f:
    pr_num = int(f.read())

  client = Github(os.environ.get('GH_TOKEN'))
  unmerged = downstreams.find_unmerged_downstreams(client, pr_num)
  if unmerged:
    raise ValueError("some PRs are unmerged", unmerged)
