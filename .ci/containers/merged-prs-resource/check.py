#! /usr/local/bin/python
from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

from absl import app
import json
import collections
import sys
import get_downstream_prs
import re
import itertools
import operator
from github import Github

def main(argv):
  in_json = json.load(sys.stdin)
  out_version = {}
  g = Github(in_json['source']['token'])
  open_pulls = g.get_repo(in_json['source']['repo']).get_pulls(state='open')
  # For each open pull request, get all the dependencies.
  depends = itertools.chain.from_iterable(
          [get_downstream_prs.get_github_dependencies(g, open_pull.number)
              for open_pull in open_pulls])
  # for each dependency, generate a tuple - (repo, pr_number)
  parsed_dependencies = [re.match(r'https://github.com/([\w-]+/[\w-]+)/pull/(\d+)', d).groups()
      for d in depends]
  parsed_dependencies.sort(key=operator.itemgetter(0))
  # group those dependencies by repo - e.g. [("terraform-provider-google", ["123", "456"]), ...]
  for r, pulls in itertools.groupby(parsed_dependencies, key=operator.itemgetter(0)):
    repo = g.get_repo(r)
    out_version[r] = []
    for pull in pulls:
      # check whether the PR is merged - if it is, add it to the version.
      pr = repo.get_pull(int(pull[1]))
      if pr.is_merged():
        out_version[r].append(pull[1])
  for k, v in out_version.iteritems():
    out_version[k] = ','.join(v)
  print(json.dumps([out_version]))
  # version dict:
  # {
  #   "terraform-providers/terraform-provider-google": "1514,1931",
  #   "terraform-providers/terraform-provider-google-beta": "121,220",
  #   "modular-magician/ansible": "",
  # }

if __name__ == '__main__':
  app.run(main)

